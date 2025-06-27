# Hướng dẫn cài đặt Fichain

## Yêu cầu hệ thống

- Hệ điều hành: Linux, macOS hoặc Windows
- Go >= 1.23
- Node.js >= 18.x và pnpm (hoặc npm/yarn)
- Git
- PostgreSQL

## 1. Cài đặt và chạy node


### 1.1. Khởi tạo DB

```bash
cd Core/cmd/gen_genesis/
go run . --config=config.yaml --genesis=genesis.json
cp -r db ../node
```

### 1.2. Chạy node

```bash
cd Core/cmd/node/
go run . --config=config.yaml
```

## 2. Triển khai các hợp đồng thông minh

### 2.1. Triển khai hợp đồng thông minh gửi tiết kiệm

```bash
cd Dapps/deploy_dapps_script/
go run . --config=config.yaml --data=deploy_data.json
```


### 2.2. Triển khai hợp đồng thông minh đầu tư vàng

```bash
cd Dapps/deploy_dapps_script/
go run . --config=config.yaml --data=deploy_gold_invest.json
```

### 2.3. Triển khai hợp đồng thông minh thanh toán tiên ích

```bash
cd Dapps/deploy_dapps_script/
go run . --config=config.yaml --data=deploy_service_bills.json
```

### 2.4. Triển khai hợp đồng thông minh hóa đơn

```bash
cd Dapps/deploy_dapps_script/
go run . --config=config.yaml --data=deploy_invoice.json
```

### 2.5. Từ các file result lấy đại chỉ các smart contract mới được deploy

```bash
ContractAddress: 0xCfD73870154A10e35f20ab9992f6f4dEf344829D => Đây là địa chỉ mới deploy nằm trong các file result.dat
```


## 3. Cài đặt trình khám phá (Explorer)

### 3.1. Chỉnh sửa config.yaml
```bash
vim Core/cmd/explorer/config.yaml
database:
  # The hostname or IP address of the database server.
  host:

  # The port the database server is listening on.
  port: 

  # The username for connecting to the database.
  user: 

  # The password for the database user.
  password: 

  # The name of the specific database to connect to.
  dbname: 

  # The SSL mode for the database connection.
  sslmode: 
```

### 3.2. Chạy explorer
```bash
go run . --config=config.yaml
```
## 4. Cài đặt cầu nối (Bridge)

### 4.1. Chỉnh sửa config.yaml
```bash
vim Dapps/bridge/server/config.yaml
database:
  # The hostname or IP address of the database server.
  host:

  # The port the database server is listening on.
  port: 

  # The username for connecting to the database.
  user: 

  # The password for the database user.
  password: 

  # The name of the specific database to connect to.
  dbname: 

  # The SSL mode for the database connection.
  sslmode: 

token_map:
  # Address smart contract USDT in BSC
  USDT: 

  # Address smart contract ETH in BSC
  ETH: 
  
  # Address smart contract BTC in BSC
  BTC: 

fichain_token_map:
  # Các địa chỉ dưới đây lấy từ bước 2
  # Address smart contract USDT in Fichain 
  USDT: 

  # Address smart contract ETH in Fichain
  ETH: 
  
  # Address smart contract BTC in Fichain
  BTC: 
```

### 4.2. Chạy bridge
```bash
go run . --config=config.yaml
```

## 5. Chạy web

### 5.1. Chỉnh sửa file .env
```bash
vim fichain-web/.env
# Các địa chỉ dưới đây lấy từ bước 2
NEXT_PUBLIC_SAVING_CONTRACT_ADDRESS=
NEXT_PUBLIC_SERVICE_BILL_CONTRACT_ADDRESS=
NEXT_PUBLIC_GOLD_TOKEN_CONTRACT_ADDRESS=
NEXT_PUBLIC_GOLD_INVEST_CONTRACT_ADDRESS=
NEXT_PUBLIC_INVOICE_CONTRACT_ADDRESS=
```

### 5.2. Cài các package cần thiết
```bash
cd fichain-web
pnpm install
# hoặc dùng npm install nếu không có pnpm
```

### 5.3. Chạy web
```bash
cd fichain-web
pnpm dev
# hoặc npm run dev
```
