// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";

/**
 * @title GoldInvest
 * @dev A smart contract that acts as a decentralized exchange for a specific
 * "Gold" ERC20 token. The owner can set the buy and sell prices. Users can
 * buy gold tokens with Ether and sell their gold tokens for Ether.
 */
contract GoldInvest is Ownable, ReentrancyGuard {

    // --- State Variables ---

    IERC20 public immutable goldToken; // The instance of the Gold ERC20 token contract

    // Prices are in wei for one full unit of the gold token (1 * 10^18 token units)
    uint256 public buyPricePerUnit;  // How much Ether a user must pay to buy 1 gold unit
    uint256 public sellPricePerUnit; // How much Ether a user receives for selling 1 gold unit

    // The smallest unit of the token (10^18 for tokens with 18 decimals)
    uint256 private constant TOKEN_UNIT = 10**18;

    // --- Events ---

    event PricesUpdated(uint256 newBuyPrice, uint256 newSellPrice);
    event GoldPurchased(address indexed buyer, uint256 etherAmount, uint256 goldAmount);
    event GoldSold(address indexed seller, uint256 goldAmount, uint256 etherAmount);
    event FundsWithdrawn(address indexed to, uint256 etherAmount);
    event TokensWithdrawn(address indexed to, uint256 tokenAmount);

    // --- Constructor ---

    /**
     * @param _goldTokenAddress The address of the deployed Gold ERC20 token contract.
     */
    constructor(address _goldTokenAddress) Ownable(msg.sender) {
        require(_goldTokenAddress != address(0), "Token address cannot be zero");
        goldToken = IERC20(_goldTokenAddress);
    }

    // --- Owner-Only Functions ---

    /**
     * @dev Sets the buy and sell prices for one unit of the gold token.
     * @param _buyPrice The new price in wei for a user to buy 1 gold unit.
     * @param _sellPrice The new price in wei for a user to sell 1 gold unit.
     */
    function setPrices(uint256 _buyPrice, uint256 _sellPrice) public onlyOwner {
        require(_buyPrice > 0, "Buy price must be positive");
        require(_sellPrice > 0, "Sell price must be positive");
        // It's good practice to ensure the buy price is higher than the sell price
        // to create a spread and prevent immediate arbitrage loss for the owner.
        require(_buyPrice > _sellPrice, "Buy price must be greater than sell price");

        buyPricePerUnit = _buyPrice;
        sellPricePerUnit = _sellPrice;
        emit PricesUpdated(_buyPrice, _sellPrice);
    }

    /**
     * @dev Owner can withdraw accumulated Ether (profits) from the contract.
     */
    function withdrawEther() public onlyOwner {
        uint256 balance = address(this).balance;
        require(balance > 0, "No Ether to withdraw");
        (bool success, ) = owner().call{value: balance}("");
        require(success, "Ether withdrawal failed");
        emit FundsWithdrawn(owner(), balance);
    }

    /**
     * @dev Owner can withdraw any gold tokens held by this contract.
     */
    function withdrawTokens() public onlyOwner {
        uint256 balance = goldToken.balanceOf(address(this));
        require(balance > 0, "No tokens to withdraw");
        goldToken.transfer(owner(), balance);
        emit TokensWithdrawn(owner(), balance);
    }

    // --- User-Facing Functions ---

    /**
     * @dev Allows a user to buy gold tokens by sending Ether.
     */
    function buyGold() public payable nonReentrant {
        require(buyPricePerUnit > 0, "Buying is currently disabled");
        require(msg.value > 0, "Must send Ether to buy gold");

        // Calculate how much gold the user gets, maintaining precision
        uint256 goldAmountToReceive = (msg.value) / buyPricePerUnit;
        
        // Check if the contract has enough gold tokens in its reserve to sell
        uint256 contractGoldBalance = goldToken.balanceOf(address(this));
        require(contractGoldBalance >= goldAmountToReceive, "Not enough gold in reserve to fulfill order");

        // Transfer the gold tokens to the buyer
        goldToken.transfer(msg.sender, goldAmountToReceive);

        emit GoldPurchased(msg.sender, msg.value, goldAmountToReceive);
    }

    /**
     * @dev Allows a user to sell their gold tokens to the contract to receive Ether.
     * IMPORTANT: The user must first call `approve()` on the gold token contract,
     * allowing this contract to spend their tokens.
     * @param _goldAmountToSell The amount of gold tokens (in the smallest unit) to sell.
     */
    function sellGold(uint256 _goldAmountToSell) public nonReentrant {
        require(sellPricePerUnit > 0, "Selling is currently disabled");
        require(_goldAmountToSell > 0, "Must sell a positive amount of gold");

        // Calculate how much Ether the user will receive
        uint256 etherAmountToReceive = (_goldAmountToSell * sellPricePerUnit) / TOKEN_UNIT;

        // Check if this contract has enough Ether to pay the user
        require(address(this).balance >= etherAmountToReceive, "Not enough Ether in reserve to fulfill order");
        
        // This is the critical step that requires prior approval from the user.
        // The contract pulls the gold tokens from the user's wallet.
        bool success = goldToken.transferFrom(msg.sender, address(this), _goldAmountToSell);
        require(success, "Token transfer failed. Did you approve first?");

        // Send the Ether to the seller
        (bool sent, ) = msg.sender.call{value: etherAmountToReceive}("");
        require(sent, "Ether payment to seller failed");

        emit GoldSold(msg.sender, _goldAmountToSell, etherAmountToReceive);
    }

    // --- View Functions ---
    /**
     * @dev Returns the current buy and sell prices.
     */
    function getPrices() public view returns (uint256, uint256) {
        return (buyPricePerUnit, sellPricePerUnit);
    }

    /**
     * @dev Returns the contract's current reserves.
     */
    function getReserves() public view returns (uint256 etherBalance, uint256 goldTokenBalance) {
        etherBalance = address(this).balance;
        goldTokenBalance = goldToken.balanceOf(address(this));
    }
}
