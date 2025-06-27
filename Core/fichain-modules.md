# 🧱 Fichain – Blockchain Layer 1 Cho Ngân Hàng Việt Nam

Fichain là một blockchain Layer 1 được thiết kế đặc biệt để tích hợp với hệ thống ngân hàng Việt Nam. Blockchain này sử dụng cơ chế đồng thuận **Proof of Authority (PoA)**, hỗ trợ EVM, không có native token mà dùng tiền fiat (như VNĐ), và yêu cầu xác thực KYC từ ngân hàng trước khi tài khoản có thể giao dịch.

---

## 🔧 Danh sách các module chính

### 1. `p2p/` – Mạng ngang hàng
- Kết nối giữa các node validator.
- Gửi/nhận block, transaction và thông điệp đồng thuận.
- Có thể sử dụng thư viện như `libp2p`.

### 2. `consensus/` – Cơ chế đồng thuận
- Thực hiện cơ chế **Proof of Authority**.
- Quản lý lượt tạo block giữa các validator.
- Kiểm tra chữ ký, chống fork, và đảm bảo đồng thuận.

### 3. `block/` – Cấu trúc block
- Định nghĩa block header, nội dung block.
- Kiểm tra tính hợp lệ và tạo block mới.

### 4. `chain/` – Chuỗi khối
- Duy trì chuỗi block và các fork nếu có.
- Tương tác với state, block, và transaction.
- Quản lý chiều cao và block hiện tại.

### 5. `transaction/` – Giao dịch
- Cấu trúc và kiểm tra chữ ký giao dịch.
- Định tuyến đến mempool hoặc executor.
- Hỗ trợ định dạng tương thích với EVM.

### 6. `state/` – Trạng thái chuỗi
- Quản lý account, balance, smart contract state.
- Hỗ trợ snapshot, rollback và state diff.

### 7. `mempool/` – Bộ nhớ đệm giao dịch
- Quản lý danh sách giao dịch chờ xử lý.
- Ưu tiên theo thời gian hoặc loại giao dịch.

### 8. `evm/` – Công cụ thực thi EVM
- Tích hợp với EVM để hỗ trợ smart contract.
- Quản lý log, sự kiện và trạng thái thực thi.

### 9. `account/` – Quản lý tài khoản
- Quản lý người dùng, tài khoản hợp lệ.
- Liên kết tài khoản với KYC từ ngân hàng.

### 10. `validator/` – Trình xác thực
- Danh sách các ngân hàng đóng vai trò validator.
- Thêm/xóa validator, kiểm soát quyền tạo block.

### 11. `kyc/` – Tích hợp xác thực
- Tích hợp dữ liệu eKYC từ ngân hàng.
- Lưu hash eKYC on-chain, dữ liệu đầy đủ lưu off-chain.

### 12. `storage/` – Lưu trữ
- Lưu block mới (hot storage).
- Lưu block cũ (cold storage).
- Cơ chế tự động archive các block cũ.

### 13. `rpc/` – API/RPC Layer
- Hỗ trợ JSON-RPC (giống Geth).
- REST API phục vụ ngân hàng và Dapp.
- Các endpoint: block, tx, account, validator,...

### 14. `smartcontract/` – Quản lý hợp đồng
- Triển khai và thực thi smart contract.
- Ghi log, call event, kiểm tra hợp lệ.

### 15. `monitoring/` – Giám sát & logging
- Tích hợp Prometheus, Grafana.
- Theo dõi trạng thái node, block, tx rate.
- Hệ thống log chi tiết và cảnh báo sự cố.

### 16. `config/` – Cấu hình & khởi tạo mạng
- Tệp `genesis.json` định nghĩa validator ban đầu.
- Các thông số như block time, gas limit, v.v.

---

## 📁 Gợi ý cấu trúc thư mục dự án

```plaintext
fichain/
├── cmd/
│   └── node/                  # Entry point để khởi chạy node
├── account/                  
├── block/                    
├── chain/                    
├── config/                   
├── consensus/                
├── evm/                      
├── kyc/                      
├── mempool/                  
├── monitoring/               
├── p2p/                      
├── rpc/                      
├── smartcontract/            
├── state/                    
├── storage/                  
├── transaction/              
└── validator/
