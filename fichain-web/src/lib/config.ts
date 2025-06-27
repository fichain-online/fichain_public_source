interface Config {
  wsBaseUrl: string;
  ekycApiBaseUrl: string;
  explorerApiBaseUrl: string;
  bridgeApiBaseUrl: string;
  blockchainVersion: number;
  blockchainId: number;

  // contract addresses
  savingContractAddress: string;
  serviceBillContractAddress: string;

  goldTokenContractAddress: string;
  goldInvestContractAddress: string;

  invoiceContractAddress: string;
}

const version = Number(process.env.NEXT_PUBLIC_BLOCKCHAIN_VERSION);
const id = Number(process.env.NEXT_PUBLIC_BLOCKCHAIN_ID);

const config: Config = {
  wsBaseUrl: process.env.NEXT_PUBLIC_WS_BASE_URL || 'ws://localhost:9001/ws', // Default to localhost
  ekycApiBaseUrl: process.env.NEXT_PUBLIC_EKYC_API_BASE_URL || 'http://127.0.0.1:8082/api',
  explorerApiBaseUrl: process.env.NEXT_PUBLIC_EXPLORER_API_BASE_URL || 'http://127.0.0.1:8080/api',
  bridgeApiBaseUrl: process.env.NEXT_PUBLIC_BRIDGE_API_BASE_URL || 'http://127.0.0.1:8081/api',
  blockchainVersion: !isNaN(version) ? version : 1,
  blockchainId: !isNaN(id) ? id: 2510,

  savingContractAddress: process.env.NEXT_PUBLIC_SAVING_CONTRACT_ADDRESS || '0x278F9C08ba5E2f6554B71bb35FCe3831Cd36edDB',
  serviceBillContractAddress: process.env.NEXT_PUBLIC_SERVICE_BILL_CONTRACT_ADDRESS || '0xDa35eb5cc5203cF287e907fcC6BC7e8b27327De0',
  
  goldTokenContractAddress: process.env.NEXT_PUBLIC_GOLD_TOKEN_CONTRACT_ADDRESS || '0x9BC715aEeBcF34554a5F221848013108d2073d95',
  goldInvestContractAddress: process.env.NEXT_PUBLIC_GOLD_INVEST_CONTRACT_ADDRESS || '0x23573Ff8ECE6c5Eb6d34a0E99DeE1b9d3fBc5e49',

  invoiceContractAddress: process.env.NEXT_PUBLIC_INVOICE_CONTRACT_ADDRESS || '0xc0EEBb77e36e338Ee2b21A521Fc6Bf9E5AF1db8b',
};

export default config;
