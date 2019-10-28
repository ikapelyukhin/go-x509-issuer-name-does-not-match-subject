#!/usr/bin/env ruby

require 'openssl'

def print_asn_strings(obj, depth = 0)
  if obj.respond_to? :each
    obj.each do |item|
      print_asn_strings(item, depth + 1)
    end
  else
    printf("%-40s %s\n", obj.value, obj.class) if (
      obj.class.to_s.match(/String/) &&
      obj.class != OpenSSL::ASN1::BitString
    )
  end
end

raise "Usage: #{$0} cert-file.crt" unless ARGV[0]

certificate = OpenSSL::X509::Certificate.new(File.read(ARGV[0]))
asn1 = OpenSSL::ASN1.decode(certificate.to_der)
print_asn_strings(asn1)
