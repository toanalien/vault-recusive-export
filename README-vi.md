# Vault Recursive Copy

Đây là một công cụ CLI để sao chép đệ quy các bí mật từ một Vault instance sang một tệp JSON.

## Lưu đồ

```mermaid
graph TD
    A[Bắt đầu] --> B{Đọc địa chỉ Vault và Token};
    B --> C{Liệt kê tất cả các secret engine};
    C --> D{Lọc các KV engine};
    D --> E{Lặp qua từng KV engine};
    E --> F{Liệt kê tất cả các bí mật đệ quy};
    F --> G{Đọc giá trị bí mật};
    G --> H{Ghi vào tệp JSON};
    H --> I[Kết thúc];
```

## Build

Để build ứng dụng, bạn cần cài đặt Go. Sau đó, bạn có thể chạy lệnh sau:

```bash
go build
```

## Sử dụng

Để sử dụng ứng dụng, bạn có thể chạy lệnh sau:

```bash
./vault-recursive-copy --token <your-vault-token> --addr <your-vault-address> --output secrets.json
```

### Flags

*   `--token`: Vault token của bạn.
*   `--addr`: Địa chỉ của Vault instance.
*   `--output`: Tệp đầu ra để lưu các bí mật. Mặc định là `secrets.json`.
