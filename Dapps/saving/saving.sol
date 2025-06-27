// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/**
 * @title SavingsContract
 * @author B-chain
 * @notice An advanced savings contract with multiple, owner-defined saving tiers.
 * Each tier has a specific lock duration, interest rate, and a lower rate for early withdrawals.
 * Interest is paid from a pre-funded pool managed by the contract owner.
 */
contract SavingsContract {
    //=========== Constants ===========//
    uint256 private constant BASIS_POINTS = 10000; // 100% = 10,000 BPS
    uint256 private constant SECONDS_IN_YEAR = 31536000;

    //=========== State Variables ===========//
    address public owner;
    uint256 public interestPool;

    // Struct for defining a savings plan/tier
    struct SavingTier {
        uint256 duration;                 // Lock duration in seconds
        uint256 interestRateBps;          // APR for completing the full term
        uint256 earlyWithdrawalRateBps;   // APR for withdrawing before the term ends
        bool isAvailable;                 // Whether this tier can be used for new savings
    }

    // Struct for an individual user's saving
    struct Saving {
        address owner;
        uint256 principal;
        uint256 creationTime;
        uint256 unlockTime;
        uint256 tierId;                   // Links the saving to its tier
        bool isActive;
    }

    SavingTier[] public tiers;
    Saving[] public savings;
    mapping(address => uint[]) public userSavings;

    //=========== Events ===========//
    event TierCreated(uint indexed tierId, uint256 duration, uint256 interestRateBps);
    event TierAvailabilityChanged(uint indexed tierId, bool isAvailable);
    event SavingCreated(uint indexed savingId, address indexed owner, uint256 principal, uint256 unlockTime, uint tierId);
    event Withdrawal(uint indexed savingId, uint256 principal, uint256 interest);
    event EarlyWithdrawal(uint indexed savingId, uint256 principal, uint256 interest);
    event InterestPoolFunded(address indexed funder, uint256 amount);

    //=========== Modifiers ===========//
    modifier onlyOwner() {
        require(msg.sender == owner, "Only the contract owner can call this function.");
        _;
    }

    modifier onlySavingOwner(uint _savingId) {
        require(_savingId < savings.length, "Saving ID does not exist.");
        require(msg.sender == savings[_savingId].owner, "You are not the owner of this saving.");
        _;
    }

    //=========== Constructor ===========//
    constructor() {
        owner = msg.sender;
    }

    //=========== Owner Functions ===========//

    /**
     * @notice Adds a new savings tier (e.g., "7 days at 3%").
     * @param _duration The lock duration in seconds.
     * @param _interestRateBps The full-term annual interest rate in BPS.
     * @param _earlyWithdrawalRateBps The early-withdrawal annual interest rate in BPS.
     */
    function addTier(
        uint256 _duration,
        uint256 _interestRateBps,
        uint256 _earlyWithdrawalRateBps
    ) external onlyOwner {
        require(_duration > 0, "Duration must be positive.");
        require(_earlyWithdrawalRateBps <= _interestRateBps, "Early rate cannot exceed full rate.");
        
        uint tierId = tiers.length;
        tiers.push(
            SavingTier({
                duration: _duration,
                interestRateBps: _interestRateBps,
                earlyWithdrawalRateBps: _earlyWithdrawalRateBps,
                isAvailable: true
            })
        );
        emit TierCreated(tierId, _duration, _interestRateBps);
    }

    /**
     * @notice Toggles the availability of a tier for new savings.
     * @param _tierId The ID of the tier to toggle.
     */
    function toggleTierAvailability(uint _tierId) external onlyOwner {
        require(_tierId < tiers.length, "Tier ID does not exist.");
        tiers[_tierId].isAvailable = !tiers[_tierId].isAvailable;
        emit TierAvailabilityChanged(_tierId, tiers[_tierId].isAvailable);
    }

    function fundInterestPool() external payable onlyOwner {
        require(msg.value > 0, "Funding amount must be positive.");
        interestPool += msg.value;
        emit InterestPoolFunded(msg.sender, msg.value);
    }

    //=========== Core User Functions ===========//

    /**
     * @notice Creates a new saving pot by selecting a tier.
     * @param _tierId The ID of the desired saving tier.
     */
    function createSaving(uint256 _tierId) external payable {
        require(msg.value > 0, "Cannot create an empty saving.");
        require(_tierId < tiers.length, "Tier ID does not exist.");
        require(tiers[_tierId].isAvailable, "This tier is not available for new savings.");

        SavingTier storage selectedTier = tiers[_tierId];
        uint256 savingId = savings.length;
        uint256 unlockTime = block.timestamp + selectedTier.duration;

        savings.push(
            Saving({
                owner: msg.sender,
                principal: msg.value,
                creationTime: block.timestamp,
                unlockTime: unlockTime,
                tierId: _tierId,
                isActive: true
            })
        );

        userSavings[msg.sender].push(savingId);
        emit SavingCreated(savingId, msg.sender, msg.value, unlockTime, _tierId);
    }

    /**
     * @notice Withdraws funds AFTER the lock period is over to receive full interest.
     * @param _savingId The ID of the saving to withdraw from.
     */
    function withdraw(uint _savingId) external onlySavingOwner(_savingId) {
        Saving storage currentSaving = savings[_savingId];

        require(currentSaving.isActive, "Saving already withdrawn.");
        require(block.timestamp >= currentSaving.unlockTime, "Lock period has not ended. Use withdrawEarly() instead.");
        
        uint256 interest = _calculateInterest(_savingId, false);
        _processWithdrawal(currentSaving, interest);

        emit Withdrawal(_savingId, currentSaving.principal, interest);
    }

    /**
     * @notice Withdraws funds BEFORE the lock period is over for a lower interest rate.
     * @param _savingId The ID of the saving to withdraw from.
     */
    function withdrawEarly(uint _savingId) external onlySavingOwner(_savingId) {
        Saving storage currentSaving = savings[_savingId];
        
        require(currentSaving.isActive, "Saving already withdrawn.");
        require(block.timestamp < currentSaving.unlockTime, "Lock period has ended. Use withdraw() for full interest.");

        uint256 interest = _calculateInterest(_savingId, true);
        _processWithdrawal(currentSaving, interest);

        emit EarlyWithdrawal(_savingId, currentSaving.principal, interest);
    }

    //=========== Internal & View Functions ===========//

    /**
     * @dev Internal function to handle the shared logic of withdrawals.
     */
    function _processWithdrawal(Saving storage _saving, uint256 _interest) internal {
        require(interestPool >= _interest, "Insufficient funds in interest pool.");
        uint256 totalToWithdraw = _saving.principal + _interest;

        // Effects (State change before interaction)
        _saving.isActive = false;
        interestPool -= _interest;

        // Interaction
        (bool success, ) = msg.sender.call{value: totalToWithdraw}("");
        require(success, "Failed to send Ether.");
    }

    /**
     * @notice Calculates the interest for a saving, either for a full term or early withdrawal.
     * @param _savingId The ID of the saving.
     * @param _isEarly Boolean indicating if it's an early withdrawal.
     * @return The amount of interest earned in wei.
     */
    function calculateInterest(uint256 _savingId, bool _isEarly) public view returns (uint256) {
        return _calculateInterest(_savingId, _isEarly);
    }
    
    function _calculateInterest(uint256 _savingId, bool _isEarly) internal view returns (uint256) {
        Saving storage s = savings[_savingId];
        SavingTier storage t = tiers[s.tierId];
        
        uint256 rateToUse = _isEarly ? t.earlyWithdrawalRateBps : t.interestRateBps;
        
        uint256 timeElapsed;
        if (_isEarly) {
            // Interest is accrued up to the moment of withdrawal
            timeElapsed = block.timestamp - s.creationTime;
        } else {
            // Interest is for the full, completed lock duration
            timeElapsed = s.unlockTime - s.creationTime;
        }

        // Formula: Interest = (Principal * Rate * Time) / (BPS * SecondsInYear)
        return (s.principal * rateToUse * timeElapsed) / (BASIS_POINTS * SECONDS_IN_YEAR);
    }
    
    // --- Other helper view functions ---
    function getSavingDetails(uint _savingId) public view returns (Saving memory) {
        require(_savingId < savings.length, "Saving ID does not exist.");
        return savings[_savingId];
    }
    
    function getUserSavingIds(address _user) public view returns (uint[] memory) {
        return userSavings[_user];
    }

    function getTierDetails(uint _tierId) public view returns (SavingTier memory) {
        require(_tierId < tiers.length, "Tier ID does not exist.");
        return tiers[_tierId];
    }
}
