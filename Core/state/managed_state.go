package state

// import (
// 	"math/big"
// 	"sync"
//
// 	"FichainCore/common"
// )
//
// type account struct {
// 	stateObject *stateObject
// 	nstart      *big.Int
// 	nonces      []bool
// }
//
// type ManagedState struct {
// 	*StateDB
//
// 	mu sync.RWMutex
//
// 	accounts map[common.Address]*account
// }
//
// // ManagedState returns a new managed state with the statedb as it's backing layer
// func ManageState(statedb *StateDB) *ManagedState {
// 	return &ManagedState{
// 		StateDB:  statedb.Copy(),
// 		accounts: make(map[common.Address]*account),
// 	}
// }
//
// // SetState sets the backing layer of the managed state
// func (ms *ManagedState) SetState(statedb *StateDB) {
// 	ms.mu.Lock()
// 	defer ms.mu.Unlock()
// 	ms.StateDB = statedb
// }
//
// // RemoveNonce removed the nonce from the managed state and all future pending nonces
// func (ms *ManagedState) RemoveNonce(addr common.Address, n uint64) {
// 	if ms.hasAccount(addr) {
// 		ms.mu.Lock()
// 		defer ms.mu.Unlock()
//
// 		account := ms.getAccount(addr)
// 		if n-account.nstart <= uint64(len(account.nonces)) {
// 			reslice := make([]bool, n-account.nstart)
// 			copy(reslice, account.nonces[:n-account.nstart])
// 			account.nonces = reslice
// 		}
// 	}
// }
//
// // NewNonce returns the new canonical nonce for the managed account
// func (ms *ManagedState) NewNonce(addr common.Address) uint64 {
// 	ms.mu.Lock()
// 	defer ms.mu.Unlock()
//
// 	account := ms.getAccount(addr)
// 	for i, nonce := range account.nonces {
// 		if !nonce {
// 			return account.nstart + uint64(i)
// 		}
// 	}
// 	account.nonces = append(account.nonces, true)
//
// 	return uint64(len(account.nonces)-1) + account.nstart
// }
//
// // GetNonce returns the canonical nonce for the managed or unmanaged account.
// //
// // Because GetNonce mutates the DB, we must take a write lock.
// func (ms *ManagedState) GetNonce(addr common.Address) *big.Int {
// 	ms.mu.Lock()
// 	defer ms.mu.Unlock()
//
// 	if ms.hasAccount(addr) {
// 		account := ms.getAccount(addr)
// 		// return uint64(len(account.nonces)) + account.nstart
// 		return new(
// 			big.Int,
// 		).Add(big.NewInt(int64(len(account.nonces))), big.NewInt(int64(account.nstart)))
// 	} else {
// 		return ms.StateDB.GetNonce(addr)
// 	}
// }
//
// // SetNonce sets the new canonical nonce for the managed state
// func (ms *ManagedState) SetNonce(addr common.Address, nonce *big.Int) {
// 	ms.mu.Lock()
// 	defer ms.mu.Unlock()
//
// 	so := ms.GetOrNewStateObject(addr)
// 	so.SetNonce(nonce)
//
// 	ms.accounts[addr] = newAccount(so)
// }
//
// // HasAccount returns whether the given address is managed or not
// func (ms *ManagedState) HasAccount(addr common.Address) bool {
// 	ms.mu.RLock()
// 	defer ms.mu.RUnlock()
// 	return ms.hasAccount(addr)
// }
//
// func (ms *ManagedState) hasAccount(addr common.Address) bool {
// 	_, ok := ms.accounts[addr]
// 	return ok
// }
//
// // populate the managed state
// func (ms *ManagedState) getAccount(addr common.Address) *account {
// 	if account, ok := ms.accounts[addr]; !ok {
// 		so := ms.GetOrNewStateObject(addr)
// 		ms.accounts[addr] = newAccount(so)
// 	} else {
// 		// Always make sure the state account nonce isn't actually higher
// 		// than the tracked one.
// 		so := ms.StateDB.getStateObject(addr)
// 		if so != nil &&
// 			new(big.Int).Add(
// 				big.NewInt(int64(len(account.nonces))),
// 				big.NewInt(int64(account.nstart)),
// 			).Cmp(so.Nonce()) < 0 {
// 			ms.accounts[addr] = newAccount(so)
// 		}
// 	}
// 	return ms.accounts[addr]
// }
//
// func newAccount(so *stateObject) *account {
// 	return &account{so, so.Nonce(), nil}
// }
