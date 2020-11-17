package com.android.signpayload;


import com.google.protobuf.ByteString;

import java.io.IOException;
import java.io.RandomAccessFile;
import java.nio.Buffer;
import java.nio.ByteBuffer;
import java.nio.charset.Charset;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.security.GeneralSecurityException;
import java.security.KeyFactory;
import java.security.PrivateKey;
import java.security.Signature;
import java.security.spec.PKCS8EncodedKeySpec;
import java.util.Arrays;

import chromeos_update_engine.UpdateMetadata.DeltaArchiveManifest;
import chromeos_update_engine.UpdateMetadata.Signatures;

import java.util.Base64;

/**
 * The current OTA signing process involve 3 steps:
 * 1. call delta_generator to create 2 hash digest for metadata and payload
 * 2. sign these 2 digests with external signer
 * 3. add them back to the payload
 *
 * This payload signer tries to do the above step in 1 shot. So it's different from the signer in
 * java/com/google/wireless/android/buildtools/signingservice/AbOtaPayloadSigner.java
 *
 * The signing strategy:
 * 1. Calculate the signature protobuf size with the new key
 * 2. Read the payload manifest from the unsigned payload, update its signature size field
 * 3. Write the updated metadata to output payload
 * 4. Sign the output metadata, write the metadata signature to output payload
 * 5. Copy the actual patch blob from unsigned payload to output payload
 * 6. Sign the data that include the output metadata & patch blob. Write the result signature to
 * the end of output payload as the payload signature.
 */
class PayloadSigner {
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
     * // The serialized DeltaArchiveManifest protobuf & signatures protobuf
     * char manifest[manifest_size];
     * char metadata_signature_message[metadata_signature_size];
     *
     * struct {char data[];} blobs[];  // Data blobs for payload
     *
     * // The signature of the entire payload, excluding the metadata_signature_message.
     * char payload_signatures_message[payload_signatures_message_size];
     */
    static class PayloadMetadata {
        static final byte[] MAGIC = "CrAU".getBytes();
        static final int MANIFEST_OFFSET = MAGIC.length + 8 * 2 + 4;

        private final long mPayloadVersion;
        private final int mMetadataSignatureSize;
        private final DeltaArchiveManifest mManifest;
        private final int mMetadataSize;


        public PayloadMetadata(RandomAccessFile payload) throws IOException,
                InvalidPayloadException {
            byte[] magic = new byte[MAGIC.length];
            payload.readFully(magic, 0, MAGIC.length);
            if (!Arrays.equals(MAGIC, magic)) {
                throw new InvalidPayloadException("Invalid magic " + new String(magic));
            }

            // TODO error check
            mPayloadVersion = payload.readLong();
            int manifestSize = (int) payload.readLong();
            mMetadataSignatureSize = payload.readInt();

            System.out.printf("manifestSize %d, old metadataSignatureSize %d\n", manifestSize,
                    mMetadataSignatureSize);

            byte[] manifestBytes = new byte[manifestSize];
            payload.readFully(manifestBytes, 0, manifestSize);
            mManifest = DeltaArchiveManifest.parseFrom(manifestBytes);
            mMetadataSize = MANIFEST_OFFSET + manifestSize;

            System.out.printf("signature offset %d, signature size %d\n",
                    mManifest.getSignaturesOffset(),
                    mManifest.getSignaturesSize());

            // TODO error check with signature offset
        }

        public PayloadMetadata(long payloadVersion, int signatureProtoLength,
                DeltaArchiveManifest manifest) {
            mPayloadVersion = payloadVersion;
            mManifest = manifest;
            mMetadataSignatureSize = signatureProtoLength;
            mMetadataSize = MANIFEST_OFFSET + manifest.toByteArray().length;
        }


        public ByteBuffer outputToBuffer() {
            ByteBuffer outputBuffer = ByteBuffer.allocate(mMetadataSize);
            outputBuffer.put(MAGIC);
            outputBuffer.putLong(mPayloadVersion);

            byte[] manifestBytes = mManifest.toByteArray();
            outputBuffer.putLong(manifestBytes.length);
            outputBuffer.putInt(mMetadataSignatureSize);
            outputBuffer.put(manifestBytes);
            return outputBuffer;
        }
    }

    // The signing algorithm should be provided by the caller.
    // Note the algorithm identifier of the PKCS1-v1_5 padding has already provided by the signer.
    private static final String SIG_ALGORITHM = "SHA256withRSA";

    private final RandomAccessFile mUnsignedPayload;
    // The private key should be provided by the signing service.
    private final PrivateKey mPrivateKey;
    // The maximum signature size in bytes should be provided by the caller, and the size bundles
    // with the signing algorithm. For example, 2048 bits RSA key has a fixed size of 256 bytes.
    // But signature signed with EC key may have varied length. So the signer will pad the signature
    // to let the signature protobuf have a fixed size.
    private final int mMaximumSignatureSize;
    private final RandomAccessFile mOutputPayload;

    public PayloadSigner(String payloadPath, String privateKeyPath,
            int maximumSignatureSize, String outputPath) throws Exception {
        mUnsignedPayload = new RandomAccessFile(payloadPath, "r");
        mPrivateKey = readPrivateKey(privateKeyPath);
        mMaximumSignatureSize = maximumSignatureSize;
        mOutputPayload = new RandomAccessFile(outputPath, "rw");
        mOutputPayload.setLength(0);
    }

    public void sign() throws Exception {
        // Calculate the size of the signature protobuf with a fake signature. The manifest and
        // payload signature protobuf will have the same length; because we sign both with the same
        // key.
        int signatureProtoLength = convertSignatureToProtobuf(
                new byte[mMaximumSignatureSize]).length;
        System.err.printf("Manifest signature proto length %d\n", signatureProtoLength);

        PayloadMetadata unsignedPayloadMetadata = new PayloadMetadata(mUnsignedPayload);
        // Update the signature size field in the manifest.
        DeltaArchiveManifest updatedManifest = unsignedPayloadMetadata.mManifest.toBuilder()
                .setSignaturesSize(signatureProtoLength).build();
        PayloadMetadata updatedPayloadMetadata = new PayloadMetadata(
                unsignedPayloadMetadata.mPayloadVersion, signatureProtoLength, updatedManifest);
        ByteBuffer updateMetadataBuffer = updatedPayloadMetadata.outputToBuffer();

        // Sign the metadata; then write the metadata and signature to the output payload.
        byte[] manifestSignature = signMetadata(updateMetadataBuffer);
        mOutputPayload.write(updateMetadataBuffer.array());
        mOutputPayload.write(convertSignatureToProtobuf(manifestSignature));

        // Copy the actual patch blobs from the input payload to output, and update the payload
        // signature.
        Signature payloadSig = Signature.getInstance(SIG_ALGORITHM);
        payloadSig.initSign(mPrivateKey);
        payloadSig.update(updateMetadataBuffer.array());
        copyPatchBlobAndUpdateSignature(unsignedPayloadMetadata, payloadSig);

        // Sign the entire payload.
        byte[] payloadSignature = payloadSig.sign();
        if (payloadSignature.length > mMaximumSignatureSize) {
            throw new SigningException("manifestSignature size too large");
        }
        mOutputPayload.write(convertSignatureToProtobuf(payloadSignature));
    }

    private byte[] signMetadata(ByteBuffer metadataBuffer)
            throws SigningException, GeneralSecurityException {
        Signature sig = Signature.getInstance(SIG_ALGORITHM);
        sig.initSign(mPrivateKey);
        sig.update(metadataBuffer.array());

        byte[] metadataSignature = sig.sign();
        if (metadataSignature.length > mMaximumSignatureSize) {
            throw new SigningException("manifestSignature size too large");
        }

        return metadataSignature;
    }

    // Copy the actual patch blobs, and update the signature with the patch bytes.
    void copyPatchBlobAndUpdateSignature(PayloadMetadata unsignedPayloadMetadata,
            Signature payloadSig)
            throws IOException, GeneralSecurityException {
        long patchBlobOffset = unsignedPayloadMetadata.mMetadataSize;
        long patchBlobSize = unsignedPayloadMetadata.mManifest.getSignaturesOffset();

        byte[] buffer = new byte[1024];
        mUnsignedPayload.seek(patchBlobOffset);
        int consumed = 0;
        while (consumed < patchBlobSize) {
            int readLimit = (int) Math.min(buffer.length, patchBlobSize - consumed);
            int bytesRead = mUnsignedPayload.read(buffer, 0, readLimit);

            payloadSig.update(buffer, 0, bytesRead);
            mOutputPayload.write(buffer, 0, bytesRead);
            consumed += bytesRead;
        }
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
        PayloadSigner payloadSigner = new PayloadSigner(payloadPath, privateKeyPath,
                maximumSignatureSize, outputPath);
        payloadSigner.sign();
    }
}
