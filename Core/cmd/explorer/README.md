## Command to create db:
```
docker run --name fichain-postgres \
  -p 5432:5432 \
  -e POSTGRES_USER=explorer_admin \
  -e POSTGRES_PASSWORD=super_secret_database_password \
  -e POSTGRES_DB=fichain_explorer \
  -d postgres
```

## Test key:
```
Private key: b185138661f1b075e6ae06789f25f1a8654642e2ed6d3045f63a6b41d37c2cac
Public key: 0x049d240c294ec39f6a00df520b5d14d386b5340db4e141fec73a033a2e393e667012c9a64c1dfae701fa484e409e6db32868736263fbb9cafc59a38f69129902b6
Address: 0x56E2558eb1e16035bf7f290f55997500Df4692EE
```
