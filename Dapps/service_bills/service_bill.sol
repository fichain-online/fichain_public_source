pragma solidity ^0.8.20;

import "@openzeppelin/contracts/utils/Counters.sol";

/**
 * @title ServiceBillManager
 * @dev A contract where an owner (e.g., landlord, service provider) can create bills
 * for specific customers, who can then pay them. Payments are forwarded to the owner.
 * This contract acts as a transparent, on-chain invoicing system.
 */
contract ServiceBillManager {
    using Counters for Counters.Counter;

    // --- State Variables ---
    address public immutable owner; // The service provider who creates bills and receives payments
    Counters.Counter private _billIds;

    enum BillType { Water, Electric, Internet, Rent, ServiceFee, Other }

    struct Bill {
        uint256 id;
        address customer;       // The address of the person who must pay the bill
        string description;
        BillType billType;
        uint256 amount;         // Amount in wei
        uint256 dueDate;
        uint256 paymentDate;
        bool isPaid;
        uint256 usageValue;
        bytes32 usageUnit;
    }

    // --- Mappings and Arrays ---
    mapping(uint256 => Bill) public bills;
    mapping(address => uint256[]) public billsForCustomer; // Find all bills for a specific customer
    uint256[] public allIssuedBillIds; // A list of all bill IDs created by the owner

    // --- Events ---
    event BillCreated(
        uint256 indexed billId,
        address indexed customer,
        uint256 amount,
        string description
    );

    event BillPaid(
        uint256 indexed billId,
        address indexed customer,
        uint256 amount,
        uint256 paymentDate
    );
    
    event BillCancelled(uint256 indexed billId);

    // --- Constructor ---
    constructor() {
        owner = msg.sender;
    }

    // --- Modifiers ---
    modifier onlyOwner() {
        require(msg.sender == owner, "Only the owner can perform this action.");
        _;
    }

    // --- Core Functions ---

    /**
     * @dev Owner creates a new bill for a specific customer.
     * @param _customer The address of the customer who needs to pay.
     * @param _description A short description of the bill.
     * @param _billType The type of the bill.
     * @param _amount The amount due, in wei.
     * @param _dueDate The timestamp when the bill is due.
     * @param _usageValue The numerical value of consumption (e.g., 150).
     * @param _usageUnit The unit of consumption, as bytes32 (e.g., "kWh").
     */
    function createBillForCustomer(
        address _customer,
        string memory _description,
        BillType _billType,
        uint256 _amount,
        uint256 _dueDate,
        uint256 _usageValue,
        bytes32 _usageUnit
    ) public onlyOwner {
        require(_customer != address(0), "Customer address cannot be zero");
        require(_amount > 0, "Bill amount must be greater than zero");

        _billIds.increment();
        uint256 newBillId = _billIds.current();

        bills[newBillId] = Bill({
            id: newBillId,
            customer: _customer,
            description: _description,
            billType: _billType,
            amount: _amount,
            dueDate: _dueDate,
            paymentDate: 0,
            isPaid: false,
            usageValue: _usageValue,
            usageUnit: _usageUnit
        });

        // Add bill ID to relevant lookup tables
        billsForCustomer[_customer].push(newBillId);
        allIssuedBillIds.push(newBillId);

        emit BillCreated(newBillId, _customer, _amount, _description);
    }

    /**
     * @dev A customer pays their specific bill.
     * The ETH sent with the transaction is used for payment.
     * The payment is forwarded to the contract owner.
     * @param _billId The ID of the bill to be paid.
     */
    function payBill(uint256 _billId) public payable {
        Bill storage billToPay = bills[_billId];

        // --- Checks ---
        require(billToPay.id != 0, "Bill does not exist.");
        require(msg.sender == billToPay.customer, "You are not the designated customer for this bill.");
        require(!billToPay.isPaid, "Bill is already paid.");
        require(msg.value >= billToPay.amount, "Insufficient ETH sent for payment.");

        // --- Effects ---
        billToPay.isPaid = true;
        billToPay.paymentDate = block.timestamp;

        // --- Interactions ---
        // Transfer the exact bill amount to the contract owner.
        (bool success, ) = owner.call{value: billToPay.amount}("");
        require(success, "Payment transfer to owner failed.");

        // Refund any overpayment back to the customer (msg.sender).
        uint256 refundAmount = msg.value - billToPay.amount;
        if (refundAmount > 0) {
            (bool refundSuccess, ) = msg.sender.call{value: refundAmount}("");
            require(refundSuccess, "Refund failed.");
        }

        emit BillPaid(_billId, msg.sender, billToPay.amount, block.timestamp);
    }
    
    /**
     * @dev Owner cancels an unpaid bill created in error.
     * @param _billId The ID of the bill to cancel.
     */
    function cancelBill(uint256 _billId) public onlyOwner {
        Bill storage billToCancel = bills[_billId];
        
        require(billToCancel.id != 0, "Bill does not exist.");
        require(!billToCancel.isPaid, "Cannot cancel a bill that is already paid.");
        
        billToCancel.isPaid = true;
        billToCancel.paymentDate = block.timestamp; // Mark as "processed"
        
        emit BillCancelled(_billId);
    }

    // --- View Functions (Read-only) ---

    /**
     * @dev Retrieves all bill IDs issued to a specific customer.
     * @param _customerAddress The address of the customer.
     * @return An array of bill IDs.
     */
    function getBillsForCustomer(address _customerAddress) public view returns (uint256[] memory) {
        return billsForCustomer[_customerAddress];
    }
    
    /**
     * @dev Retrieves all bill IDs ever issued by the owner.
     */
    function getAllIssuedBills() public view onlyOwner returns (uint256[] memory) {
        return allIssuedBillIds;
    }
}
