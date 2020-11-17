package com.android.signpayload;


import com.google.protobuf.ByteString;

import java.io.IOException;
import java.io.RandomAccessFile;
import java.nio.charset.Charset;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.security.GeneralSecurityException;
import java.security.KeyFactory;
import java.security.PrivateKey;
import java.security.Signature;
import java.security.spec.PKCS8EncodedKeySpec;
import java.util.Arrays;

import chromeos_update_engine.UpdateMetadata;
import chromeos_update_engine.UpdateMetadata.Extent;
import chromeos_update_engine.UpdateMetadata.DeltaArchiveManifest;
import chromeos_update_engine.UpdateMetadata.Signatures;

import java.util.Base64;

class SignPayload {

    static public final class InvalidPayloadException extends Exception {
        public InvalidPayloadException(String message) {
            super(message);
        }
    }

    static public final class SigningException extends Exception {
        public SigningException(String message) {
            super(message);
        }
    }


    /**
     * Format of the payload:
     *
     * char magic[4] = "CrAU";
     * uint64 payload_version;
     * uint64 manifest_size;  // Size of protobuf DeltaArchiveManifest
     * uint32 metadata_signature_size;
     *
     * // The serialized DeltaArchiveManifest protobuf
     * char manifest[manifest_size];
     * // The signature of the metadata
     * char metadata_signature_message[metadata_signature_size];
     *
     * // Data blobs for files
     * struct {char data[];} blobs[];
     *
     * // The signature of the entire payload, exclude the metadata_signature_message.
     * char payload_signatures_message[payload_signatures_message_size];
     */
    static class PayloadMetadata {
        static final byte[] kMagic = "CrAU".getBytes();
        static final int kManifestOffset = kMagic.length + 8 * 2 + 4;

        private final long mPayloadVersion;
        private final int mManifestSize;

        private final DeltaArchiveManifest mManifest;
        private final int mMetadataSize;


        public PayloadMetadata(RandomAccessFile payload) throws IOException,
                InvalidPayloadException {
            byte[] magic = new byte[kMagic.length];
            payload.readFully(magic, 0, kMagic.length);
            if (!Arrays.equals(kMagic, magic)) {
                throw new InvalidPayloadException("Invalid magic " + new String(magic));
            }

            // TODO error check
            mPayloadVersion = payload.readLong();
            mManifestSize = (int) payload.readLong();
            long metadataSignatureSize = payload.readInt();

            System.err.printf("manifestSize %d, old metadataSignatureSize %d\n", mManifestSize,
                    metadataSignatureSize);

            byte[] manifestBytes = new byte[mManifestSize];
            payload.readFully(manifestBytes, 0, mManifestSize);
            mManifest = DeltaArchiveManifest.parseFrom(manifestBytes);

            System.err.printf("signature offset %d, signature size %d\n",
                    mManifest.getSignaturesOffset(),
                    mManifest.getSignaturesSize());

            if (mManifest.getSignaturesOffset() == 0) {
                throw new InvalidPayloadException("Expect to have signature offset");
            }

            mMetadataSize = kManifestOffset + mManifestSize;
        }

    }

    private static final String SIG_ALGORITHM = "SHA256withRSA";

    private final RandomAccessFile mInputPayload;
    private final PrivateKey mPrivateKey;
    private final int mMaximumSignatureSize;
    private final RandomAccessFile mOutput;


    public SignPayload(String payloadPath, String privateKeyPath,
            int maximumSignatureSize, String outputPath) throws Exception {
        mInputPayload = new RandomAccessFile(payloadPath, "r");
        mPrivateKey = readPrivateKey(privateKeyPath);
        mMaximumSignatureSize = maximumSignatureSize;
        mOutput = new RandomAccessFile(outputPath, "rw");
        mOutput.setLength(0);
    }

    private static String encodeHex(byte[] bytes) {
        StringBuffer hex = new StringBuffer(bytes.length * 2);

        for (int i = 0; i < bytes.length; i++) {
            int byteIntValue = bytes[i] & 0xff;
            if (byteIntValue < 0x10) {
                hex.append("0");
            }
            hex.append(Integer.toString(byteIntValue, 16));
        }

        return hex.toString();
    }

    public void sign() throws Exception {
        PayloadMetadata payloadMetadata = new PayloadMetadata(mInputPayload);
        int metadataSize = PayloadMetadata.kManifestOffset + payloadMetadata.mManifestSize;

        System.err.printf("Maximum signature length %d\n", mMaximumSignatureSize);

        mOutput.write(PayloadMetadata.kMagic);
        mOutput.writeLong(payloadMetadata.mPayloadVersion);
        mOutput.writeLong(payloadMetadata.mManifestSize);


        byte[] manifestSignature = signManifest(metadataSize);
        if (manifestSignature.length > mMaximumSignatureSize) {
            throw new SigningException("manifestSignature size too large");
        }
        byte[] manifestSignatureProto = convertSignatureToProtobuf(manifestSignature);
        System.err.printf("Manifest signature proto length %d\n", manifestSignatureProto.length);
        mOutput.writeInt(manifestSignatureProto.length);

        long payloadSignatureSize = manifestSignatureProto.length;
        DeltaArchiveManifest updatedManifest = payloadMetadata.mManifest.toBuilder()
                .setSignaturesSize(payloadSignatureSize).build();
        mOutput.write(updatedManifest.toByteArray());

        Signature sig = Signature.getInstance(SIG_ALGORITHM);
        sig.initSign(mPrivateKey);
        updateSignature(sig, mOutput, 0, metadataSize);
        manifestSignatureProto = convertSignatureToProtobuf(sig.sign());

        mOutput.write(manifestSignatureProto);

        Signature payloadSig = Signature.getInstance(SIG_ALGORITHM);
        payloadSig.initSign(mPrivateKey);

        updateSignature(payloadSig, mOutput, 0, payloadMetadata.mMetadataSize);

        long bytesToCopy = payloadMetadata.mManifest.getSignaturesOffset();
        long consumed = 0;
        byte[] buffer = new byte[1024];
        mInputPayload.seek(payloadMetadata.mMetadataSize);
        mOutput.seek(payloadMetadata.mMetadataSize + manifestSignatureProto.length);
        while (consumed < bytesToCopy) {
            int readLimit = (int) Math.min(buffer.length, bytesToCopy - consumed);
            int bytesRead = mInputPayload.read(buffer, 0, readLimit);

            payloadSig.update(buffer, 0, bytesRead);
            mOutput.write(buffer, 0, bytesRead);
            consumed += bytesRead;
        }

        byte[] payloadSignature = payloadSig.sign();
        byte[] payloadSignatureProto = convertSignatureToProtobuf(payloadSignature);
        System.err.printf("payload signature proto length %d\n", payloadSignatureProto.length);

        mOutput.write(payloadSignatureProto);
    }

    private byte[] convertSignatureToProtobuf(byte[] signature) {
        byte[] paddedSignature = new byte[mMaximumSignatureSize];
        System.arraycopy(signature, 0, paddedSignature, 0, signature.length);
        Signatures sig = Signatures.newBuilder().addSignatures(
                Signatures.Signature.newBuilder()
                        .setUnpaddedSignatureSize(signature.length)
                        .setData(ByteString.copyFrom(paddedSignature))).build();

        return sig.toByteArray();
    }

    private byte[] signManifest(int metadataSize) throws Exception {
        Signature sig = Signature.getInstance(SIG_ALGORITHM);
        sig.initSign(mPrivateKey);
        updateSignature(sig, mInputPayload, 0, metadataSize);
        byte[] manifestSignature = sig.sign();
        System.err.printf("signature, length %d\n %s\n", manifestSignature.length,
                encodeHex(manifestSignature));
        return manifestSignature;
    }

    static private void updateSignature(Signature sig, RandomAccessFile payload, int offset,
            int length) throws SigningException {
        byte[] buffer = new byte[1024];
        try {
            payload.seek(offset);
            int consumed = 0;
            while (consumed < length) {
                int readLimit = Math.min(buffer.length, length - consumed);
                int bytesRead = payload.read(buffer, 0, readLimit);

                sig.update(buffer, 0, bytesRead);
                consumed += bytesRead;
            }
        } catch (IOException | GeneralSecurityException e) {
            throw new SigningException(e.getMessage());
        }
    }

    public static PrivateKey readPrivateKey(String path) throws Exception {
        String key = new String(Files.readAllBytes(Paths.get(path)), Charset.defaultCharset());

        String privateKeyPEM = key
                .replace("-----BEGIN PRIVATE KEY-----", "")
                .replaceAll(System.lineSeparator(), "")
                .replace("-----END PRIVATE KEY-----", "");

        byte[] decoded = Base64.getDecoder().decode(privateKeyPEM);

        KeyFactory keyFactory = KeyFactory.getInstance("RSA");
        PKCS8EncodedKeySpec keySpec = new PKCS8EncodedKeySpec(decoded);
        PrivateKey privateKey = keyFactory.generatePrivate(keySpec);

        System.err.printf("privateKey %s\n", privateKey.getFormat());
        return privateKey;
    }

    public static void main(String[] args) throws Exception {
        String payloadPath = "";
        String privateKeyPath = "";
        int maximumSignatureSize = 0;
        String outputPath = "";

        for (int i = 0; i < args.length; i++) {
            if ("--payload".equals(args[i])) {
                payloadPath = args[i + 1];
            } else if ("--private_key".equals(args[i])) {
                privateKeyPath = args[i + 1];
            } else if ("--maximum_signature_size".equals(args[i])) {
                maximumSignatureSize = Integer.parseInt(args[i + 1]);
            } else if ("--output".equals(args[i])) {
                outputPath = args[i + 1];
            }
        }
        SignPayload payloadSigner = new SignPayload(payloadPath, privateKeyPath,
                maximumSignatureSize, outputPath);
        payloadSigner.sign();

        Extent extent = Extent.newBuilder().setNumBlocks(1).build();
        System.err.printf("signpayload: running %d, %d\n", extent.getStartBlock(),
                extent.getNumBlocks());
    }


}
