// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title Invoicing Contract
 * @author Your Name
 * @notice A smart contract to create and pay invoices with line items and automatic tax distribution.
 */
contract Invoicing is Ownable {

    //=========== Structs and Enums ===========//

    /**
     * @notice A single line item on an invoice.
     * @param description The name or description of the product/service.
     * @param unitPrice The price per single item in wei.
     * @param quantity The number of units for this item.
     */
    struct InvoiceItem {
        string description;
        uint256 unitPrice;
        uint256 quantity;
    }

    /**
     * @notice A structure to hold all information about an invoice.
     * @param id The unique identifier of the invoice.
     * @param payee The person who created the invoice and will receive payment.
     * @param payer The person who is required to pay the invoice.
     * @param totalAmount The total calculated amount of the invoice in wei.
     * @param description A general description for the whole invoice.
     * @param items An array of line items included in the invoice.
     * @param status The current status of the invoice.
     */
    struct Invoice {
        uint256 id;
        address payable payee;
        address payer;
        uint256 totalAmount;
        string description;
        InvoiceItem[] items;
        InvoiceStatus status;
    }

    enum InvoiceStatus {
        Created,
        Paid,
        Cancelled
    }

    //=========== State Variables ===========//

    mapping(uint256 => Invoice) public invoices;
    uint256 public nextInvoiceId;
    address payable public governmentAddress;
    uint256 public taxRateBps; // Tax rate in Basis Points (BPS), e.g., 500 = 5%

    //=========== Events ===========//

    event InvoiceCreated(
        uint256 indexed id,
        address indexed payee,
        address indexed payer,
        uint256 totalAmount
    );
    event InvoicePaid(uint256 indexed id, address indexed payer, uint256 amountPaid, uint256 taxAmount);
    event InvoiceCancelled(uint256 indexed id);
    event TaxInfoUpdated(address indexed newGovernmentAddress, uint256 newTaxRateBps);

    //=========== Constructor ===========//

    constructor(
        address _initialOwner,
        address payable _initialGovernmentAddress,
        uint256 _initialTaxRateBps
    ) Ownable(_initialOwner) {
        require(_initialGovernmentAddress != address(0), "Government address cannot be zero");
        require(_initialTaxRateBps <= 10000, "Tax rate cannot exceed 10000 BPS (100%)");
        
        governmentAddress = _initialGovernmentAddress;
        taxRateBps = _initialTaxRateBps;
        nextInvoiceId = 1;
    }

    //=========== Admin Functions (Owner Only) ===========//

    function setTaxInfo(address payable _newGovernmentAddress, uint256 _newTaxRateBps) external onlyOwner {
        require(_newGovernmentAddress != address(0), "Government address cannot be zero");
        require(_newTaxRateBps <= 10000, "Tax rate cannot exceed 10000 BPS (100%)");
        
        governmentAddress = _newGovernmentAddress;
        taxRateBps = _newTaxRateBps;

        emit TaxInfoUpdated(_newGovernmentAddress, _newTaxRateBps);
    }

    //=========== Core Functions ===========//

    /**
     * @notice Creates a new invoice with multiple line items.
     * @dev The total amount is calculated automatically from the items.
     * @param _payer The address that is expected to pay the invoice.
     * @param _description A general description for the invoice (e.g., "Q4 Services").
     * @param _items An array of InvoiceItem structs representing the products/services.
     * @return The ID of the newly created invoice.
     */
    function createInvoice(
        address _payer,
        string calldata _description,
        InvoiceItem[] calldata _items
    ) external returns (uint256) {
        require(_payer != address(0), "Payer address cannot be zero");
        require(_items.length > 0, "Invoice must have at least one item");

        uint256 invoiceId = nextInvoiceId;
        uint256 calculatedTotalAmount = 0;

        // Create a copy of the items in memory to be stored.
        InvoiceItem[] memory newItems = new InvoiceItem[](_items.length);

        // Calculate total amount and copy items to memory.
        // Solidity 0.8+ automatically checks for arithmetic overflow.
        for (uint i = 0; i < _items.length; i++) {
            require(_items[i].quantity > 0, "Item quantity must be greater than zero");
            require(_items[i].unitPrice > 0, "Item unit price must be greater than zero");
            
            uint256 itemTotal = _items[i].unitPrice * _items[i].quantity;
            calculatedTotalAmount += itemTotal;

            newItems[i] = InvoiceItem({
                description: _items[i].description,
                unitPrice: _items[i].unitPrice,
                quantity: _items[i].quantity
            });
        }
        
        require(calculatedTotalAmount > 0, "Total invoice amount must be greater than zero");

        invoices[invoiceId] = Invoice({
            id: invoiceId,
            payee: payable(msg.sender),
            payer: _payer,
            totalAmount: calculatedTotalAmount,
            description: _description,
            items: newItems,
            status: InvoiceStatus.Created
        });

        nextInvoiceId++;

        emit InvoiceCreated(invoiceId, msg.sender, _payer, calculatedTotalAmount);
        return invoiceId;
    }

    /**
     * @notice Pays a specific invoice.
     * @dev The logic remains the same because we pre-calculate `totalAmount`.
     * @param _invoiceId The ID of the invoice to be paid.
     */
    function payInvoice(uint256 _invoiceId) external payable {
        Invoice storage invoice = invoices[_invoiceId];

        require(msg.sender == invoice.payer, "Only the designated payer can pay this invoice");
        require(invoice.status == InvoiceStatus.Created, "Invoice is not available for payment");
        require(msg.value == invoice.totalAmount, "Incorrect payment amount sent");
        
        invoice.status = InvoiceStatus.Paid;

        uint256 taxAmount = (invoice.totalAmount * taxRateBps) / 10000;
        uint256 payeeAmount = invoice.totalAmount - taxAmount;

        (bool successGov, ) = governmentAddress.call{value: taxAmount}("");
        require(successGov, "Failed to send tax payment");

        (bool successPayee, ) = invoice.payee.call{value: payeeAmount}("");
        require(successPayee, "Failed to send payment to payee");
        
        emit InvoicePaid(_invoiceId, msg.sender, invoice.totalAmount, taxAmount);
    }
    
    function cancelInvoice(uint256 _invoiceId) external {
        Invoice storage invoice = invoices[_invoiceId];
        
        require(msg.sender == invoice.payee, "Only the payee can cancel this invoice");
        require(invoice.status == InvoiceStatus.Created, "Only a created invoice can be cancelled");
        
        invoice.status = InvoiceStatus.Cancelled;
        
        emit InvoiceCancelled(_invoiceId);
    }

    //=========== View Functions ===========//

    /**
     * @notice A view function to get the core details of an invoice.
     * @param _invoiceId The ID of the invoice.
     * @return payee, payer, totalAmount, general description, and status.
     */
    function getInvoiceDetails(uint256 _invoiceId) external view returns (
        address payee,
        address payer,
        uint256 totalAmount,
        string memory description,
        InvoiceStatus status
    ) {
        Invoice storage inv = invoices[_invoiceId];
        return (inv.payee, inv.payer, inv.totalAmount, inv.description, inv.status);
    }

    /**
     * @notice A view function to get all the line items for a specific invoice.
     * @param _invoiceId The ID of the invoice.
     * @return An array of InvoiceItem structs.
     */
    function getInvoiceItems(uint256 _invoiceId) external view returns (InvoiceItem[] memory) {
        return invoices[_invoiceId].items;
    }
}
