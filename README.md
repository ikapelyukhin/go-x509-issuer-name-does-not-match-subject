# failed to verify certificate: x509: issuer name does not match subject from issuing certificate

## The issue

Go program fails with `failed to verify certificate: x509: issuer name does not match subject from issuing certificate`.

## Description

Issuer fields in the certificates of the signing CA and the server obviously have to be the same. However, they must also have the same ASN.1 data type.

Up-to-date versions of `openssl` seem compare the fields properly regardless of the data type. Go implementation, however, raises the error if ASN.1 data types are different.

## Files in this repo

* `print_asn_strings.rb` -- a Ruby script that prints ASN.1 data types of strings in a certificate.
* `gox509verify` -- a minimal command line utility in Go to verify validity of a certificate.
* `rmt-ca.crt` -- CA certificate with all strings in `PrintableString` and `IA5String`.
* `rmt-server-bad.crt` -- server certificate with some strings in `UTF8String`, which cause the verification error described above.
* `rmt-server-good.crt` -- server certificate with all strings `PrintableString` and `IA5String`, successfully passes verification.

## Fixing the error

1. Reproduce the error:

    Running `gox509verify smt-gce.susecloud.net rmt-ca.crt rmt-server-bad.crt` produces:

   ```
    panic: failed to verify certificate: x509: issuer name does not match subject from issuing certificate

    goroutine 1 [running]:
    main.VerifyCert(0x7fff2e002eb5, 0x15, 0x7fff2e002ecb, 0xa, 0x7fff2e002ed6, 0x12)
    ```
2. Examine certificate issuer field:

    2.1. The issuer fields have the same value:
    ```
    $ openssl x509 -in rmt-ca.crt -issuer -noout
    issuer=C = DE, ST = Bavaria, L = Nuremberg, O = SUSE, OU = CSM, CN = SUSE, emailAddress = suse-public-cloud@xxxxxxxxx.net
    $ openssl x509 -in rmt-server-bad.crt -issuer -noout
    issuer=C = DE, ST = Bavaria, L = Nuremberg, O = SUSE, OU = CSM, CN = SUSE, emailAddress = suse-public-cloud@xxxxxxxxx.net
    ```
   2.2. Examine ASN.1 data types in the certificate:
    * CA certificate, running `print_asn_strings.rb rmt-ca.crt` produces:
        ```
        DE                                       OpenSSL::ASN1::PrintableString
        Bavaria                                  OpenSSL::ASN1::PrintableString
        Nuremberg                                OpenSSL::ASN1::PrintableString
        SUSE                                     OpenSSL::ASN1::PrintableString
        CSM                                      OpenSSL::ASN1::PrintableString
        SUSE                                     OpenSSL::ASN1::PrintableString
        suse-public-cloud@xxxxxxxxx.net          OpenSSL::ASN1::IA5String
        ```
    * Server certificate, running `print_asn_strings.rb rmt-server-bad.crt` produces:
        ```
        DE                                       OpenSSL::ASN1::PrintableString
        Bavaria                                  OpenSSL::ASN1::UTF8String
        Nuremberg                                OpenSSL::ASN1::UTF8String
        SUSE                                     OpenSSL::ASN1::UTF8String
        CSM                                      OpenSSL::ASN1::UTF8String
        SUSE                                     OpenSSL::ASN1::UTF8String
        suse-public-cloud@xxxxxxxxx.net          OpenSSL::ASN1::IA5String
        ```
    * Some strings in the certificates have different ASN.1 data types.
3. Re-generate the server certificate with the same ASN.1 data type as the CA (`rmt-server-good.crt`).
4. Validate that the new certificate is valid:

    Running `gox509verify smt-gce.susecloud.net rmt-ca.crt rmt-server-good.crt` now produces `OK`.

## References

* https://github.com/golang/go/issues/31440
* https://github.com/SUSE/container-suseconnect/issues/36
