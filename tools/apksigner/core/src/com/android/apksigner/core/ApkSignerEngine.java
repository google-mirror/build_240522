/*
 * Copyright (C) 2016 The Android Open Source Project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package com.android.apksigner.core;

import java.io.Closeable;
import java.io.IOException;
import java.security.InvalidKeyException;
import java.security.SignatureException;
import java.util.List;

import com.android.apksigner.core.util.DataSource;

/**
 * APK signing logic which is independent of how input and output APKs are stored, parsed, and
 * generated.
 *
 * <p><h3>Operating Model</h3>
 *
 * The operating model is that there is an input APK which is being signed, thus producing an output
 * APK. In reality, there may be just an output APK being built from scratch, or the input APK and
 * the output APK may be the same file. Because this engine does not deal with reading and writing
 * files, it can handle all of these scenarios.
 *
 * <p>NOTE: The engine is stateful and thus cannot be used for signing multiple APKs.
 *
 * <p>In this engine's operating model, a signed output APK is produced as follows.
 * <ol>
 * <li>JAR entries to be signed are output,</li>
 * <li>JAR archive is signed using JAR signing, thus adding the so-called v1 signature to the
 *     output,</li>
 * <li>JAR archive is signed using APK Signature Scheme v2, thus adding the so-called v2 signature
 *     to the output.</li>
 * </ol>
 * If the client of this engine is in the position to skip outputting some of the JAR entries, those
 * entries should be provided into the engine as input APK JAR entries, thus letting the engine tell
 * its client which of these JAR entries should be output.
 *
 * <p>At each step, the client of {@code ApkSignerEngine} is expected to invoke the engine. Some
 * invocations may provide the client with a task to perform. The client is expected to perform all
 * requested tasks before proceeding to the next stage of signing. See documentation of each method
 * about the deadlines for performing the tasks requested by the method.
 *
 * <p>To use the engine to sign an input APK (or a collection of input JAR entries), follow these
 * steps:
 * <ol>
 * <li>Obtain a new instance of the engine. Engine instances are stateful and thus cannot be used
 *     for signing multiple APKs.</li>
 * <li>Locate the input APK's APK Signing Block and invoke
 *     {@link #inputApkSigningBlock(DataSource)}.</li>
 * <li>For each JAR entry in the input APK, invoke {@link #inputJarEntry(String)} to determine
 *     whether this entry should be output and whether the engine needs to inspect the entry's
 *     uncompressed contents.</li>
 * <li>For each output JAR entry, invoke {@link #outputJarEntry(String)} which may request to
 *     inspect the entry's uncompressed contents.</li>
 * <li>Once all JAR entries have been output, invoke {@link #outputJarEntries()} which may request
 *     that additional JAR entries are output. These entries comprise v1 signature.</li>
 * <li>Locate the ZIP Central Directory and ZIP End of Central Directory sections in the output and
 *     invoke {@link #outputZipSections(DataSource, DataSource, DataSource)} which may request that
 *     an APK Signature Block is inserted before the ZIP Central Directory. The block contains v2
 *     signature.</li>
 * <li>Invoke {@link #outputDone()} to signal that the full APK is output. The engine will confirm
 *     that the APK is signed.</li>
 * <li>Invoke {@link #close()} to signal that the engine will no longer be used.
 * </ol>
 *
 * <p><h3>Incremental Operation</h3>
 *
 * The engine supports incremental operation where the APK is signed, then changed, and then signed
 * again. If, after an APK is signed by the engine, an input APK's entry is added/changed or removed
 * APK, invoke {@link #inputJarEntry(String)} or {@link #inputJarEntryRemoved(String)} respectively,
 * and then run through step 5 onwards to re-sign the APK. Similarly, if an output APK's entry is
 * added/changed or removed, invoke {@link #outputJarEntry(String)} or
 * {@link #outputJarEntryRemoved(String)}, and then run through step 5 onwards to re-sign the APK.
 *
 * <p><h3>Output-only Operation</h3>
 *
 * The engine's general operating model consists of an input APK and an output APK. However, it is
 * possible to use the engine in output-only mode. In this mode, the engine has less control over
 * output, because it cannot request to some JAR entries are not output. Nevertheless, the engine
 * will attempt to make the output APK signed and will report an error if cannot do so. Because
 * there is no input APK, there is no need to invoke any of the {@code input...} methods of the
 * engine.
 */
public interface ApkSignerEngine extends Closeable {

    /**
     * Indicates to this engine that the input APK contains the provided APK Signing Block.
     *
     * @param apkSigningBlock APK signing block of the input APK. The provided data source is
     *        guaranteed to not be used by the engine after this method terminates.
     */
    void inputApkSigningBlock(DataSource apkSigningBlock);

    /**
     * Indicates to this engine that the specified JAR entry was encountered in the input APK.
     *
     * @return instructions about how to proceed with this entry
     */
    InputJarEntryInstructions inputJarEntry(String entryName);

    /**
     * Indicates to this engine that the specified JAR entry was output.
     *
     * <p>It is unnecessary to invoke this method for entries added to output by this engine (e.g.,
     * requested by {@link #outputJarEntries()}) provided the entries were output with exactly the
     * contents requested by the engine.
     *
     * @return request to inspect the entry's contents or {@code null} if the engine does not need
     *         to inspect the entry's contents. The request must be fulfilled before
     *         {@link #outputJarEntries()} is invoked.
     */
    InspectJarEntryContentsRequest outputJarEntry(String entryName);

    /**
     * Indicates to this engine that the specified JAR entry was removed from the input. It's safe
     * to invoke this for entries for which {@link #inputJarEntry(String)} hasn't been invoked.
     *
     * @return output policy of this JAR entry. The policy indicates how this input entry affects
     *         the output APK. The client of this engine should use this information to determine
     *         how the removal of this input APK's JAR entry affects the output APK.
     */
    InputJarEntryInstructions.OutputPolicy inputJarEntryRemoved(String entryName);

    /**
     * Indicates to this engine that the specified JAR entry was removed from the output. It's safe
     * to invoke this for entries for which {@link #outputJarEntry(String)} hasn't been invoked.
     */
    void outputJarEntryRemoved(String entryName);

    /**
     * Indicates to this engine that all JAR entries have been output.
     *
     * @return request to add v1 signature to the output or {@code null} if there is no need to add
     *         a v1 signature. The request will contain additional JAR entries to be output. The
     *         request must be fulfilled before
     *         {@link #outputZipSections(DataSource, DataSource, DataSource)} is invoked.
     *
     * @throws InvalidKeyException if a signature could not be generated because a signing key is
     *         not suitable for generating the signature
     * @throws SignatureException if an error occurred while generating the signature
     * @throws IllegalStateException if there are unfulfilled requests, such as to inspect contents
     *         of some JAR entries
     */
    OutputV1SignatureRequest outputJarEntries() throws InvalidKeyException, SignatureException;

    /**
     * Indicates to this engine that the ZIP sections comprising the output APK have been output.
     *
     * <p>The provided data sources are guaranteed to not be used by the engine after this method
     * terminates.
     *
     * @param zipEntries the section of ZIP archive containing Local File Header records and
     *        contents of the ZIP entries. In a well-formed archive, this section starts at the
     *        start of the archive and extends all the way to the ZIP Central Directory.
     * @param zipCentralDirectory ZIP Central Directory section
     * @param zipEocd ZIP End of Central Directory (EoCD) record
     *
     * @return request to add v2 signature to the output or {@code null} if there is no need to add
     *         a v2 signature. The request must be fulfilled before {@link #outputDone()} is
     *         invoked.
     *
     * @throws IOException if an I/O error occurs while reading the provided ZIP sections
     * @throws InvalidKeyException if a signature could not be generated because a signing key is
     *         not suitable for generating the signature
     * @throws SignatureException if an error occurred while generating the signature
     * @throws IllegalStateException if there are unfulfilled requests, such as to inspect contents
     *         of some JAR entries or to output v1 signature
     */
    OutputV2SignatureRequest outputZipSections(
            DataSource zipEntries,
            DataSource zipCentralDirectory,
            DataSource zipEocd) throws IOException, InvalidKeyException, SignatureException;

    /**
     * Indicates to this engine that the signed APK was output. This gives the engine the
     * opportunity to confirm that output is now signed.
     *
     * @throws IllegalStateException if there are unfulfilled requests, such as to inspect contents
     *         of some JAR entries or to output signatures
     */
    void outputDone();

    /**
     * Indicates to this engine that it will no longer be used.
     */
    @Override
    void close();

    /**
     * Instructions about how to handle an input APK's JAR entry.
     *
     * <p>The instructions indicate whether to output the entry (see {@link #getOutputPolicy()}) and
     * may contain a request to inspect the entry's contents
     * (see {@link #getInspectContentsRequest()}), in which case the entry's contents must be
     * provided to the engine before {@link ApkSignerEngine#outputJarEntries()}.
     */
    public static class InputJarEntryInstructions {
        private final OutputPolicy mOutputPolicy;
        private final InspectJarEntryContentsRequest mInspectContentsRequest;

        /**
         * Constructs a new {@code InputJarEntryInstructions} instance with the provided entry
         * output policy and without a request to inspect the entry's contents.
         */
        public InputJarEntryInstructions(OutputPolicy outputPolicy) {
            this(outputPolicy, null);
        }

        /**
         * Constructs a new {@code InputJarEntryInstructions} instance with the provided entry
         * output mode and with the provided request to inspect the entry's contents.
         *
         * @param inspectContentsRequest request to inspect the entry's uncompressed contents or
         *        {@code null} if there's no need to inspect the entry's contents.
         */
        public InputJarEntryInstructions(
                OutputPolicy outputPolicy,
                InspectJarEntryContentsRequest inspectContentsRequest) {
            mOutputPolicy = outputPolicy;
            mInspectContentsRequest = inspectContentsRequest;
        }

        /**
         * Returns the output policy for this entry.
         */
        public OutputPolicy getOutputPolicy() {
            return mOutputPolicy;
        }

        /**
         * Returns the request to inspect the entry's uncompressed contents or {@code null} if there
         * is no need to inspect the entry's contents.
         */
        public InspectJarEntryContentsRequest getInspectContentsRequest() {
            return mInspectContentsRequest;
        }

        /**
         * Output policy for an input APK's JAR entry.
         */
        public static enum OutputPolicy {
            /** Entry must not be output. */
            SKIP,

            /** Entry should be output. */
            OUTPUT,

            /** Entry will be output by the engine later, overwriting its current contents. */
            OUTPUT_WITH_ENGINE_PROVIDED_CONTENTS,
        }
    }

    /**
     * Request to inspect the uncompressed contents of a JAR entry.
     *
     * <p>Uncompressed contents of the entry must be provided to the data sink returned by
     * {@link #getContentsSink()}. Once the contents have been provided to sink, {@link #done()}
     * should be invoked.
     */
    interface InspectJarEntryContentsRequest {

        /**
         * Returns the data sink into which the uncompressed contents of the JAR entry should be
         * sent.
         */
        DataSource.Sink getContentsSink();

        /**
         * Indicates that uncompressed contents of the JAR entry have been provided in full.
         */
        void done();

        /**
         * Returns the name of the JAR entry.
         */
        String getEntryName();
    }

    /**
     * Request to add JAR signature (aka v1 signature) to the output APK.
     *
     * <p>Entries listed in {@link #getAdditionalJarEntries()} must be added to the output APK after
     * which {@link #done()} must be invoked.
     */
    interface OutputV1SignatureRequest {

        /**
         * Returns JAR entries that must be added to the output APK.
         */
        List<JarEntry> getAdditionalJarEntries();

        /**
         * Indicates that the JAR entries contained in this request were added to the output APK.
         */
        void done();

        /**
         * JAR entry.
         */
        public static class JarEntry {
            private final String mName;
            private final byte[] mContents;

            /**
             * Constructs a new {@code JarEntry} with the provided name and contents.
             *
             * @param contents contents. Changes to these data will not be reflected in
             *        {@link #getContents()}.
             */
            public JarEntry(String name, byte[] contents) {
                mName = name;
                mContents = contents.clone();
            }

            /**
             * Returns the name of this ZIP entry.
             */
            public String getName() {
                return mName;
            }

            /**
             * Returns the uncompressed contents of this JAR entry.
             */
            public byte[] getContents() {
                return mContents.clone();
            }
        }
    }

    /**
     * Request to add APK Signature Scheme v2 signature (aka v2 signature) to the output APK.
     *
     * <p>The APK Signing Block returned by {@link #getApkSigningBlock()} must be inserted into the
     * output APK immediately before the ZIP Central Directory, the offset of ZIP Central Directory
     * in the ZIP End of Central Directory record must be adjusted accordingly, and then
     * {@link #done()} must be invoked.
     *
     * <p>If the output contains an APK Signing Block, that block must be replaced by the block
     * contained in this request.
     */
    interface OutputV2SignatureRequest {

        /**
         * Returns the APK Signing Block.
         */
        byte[] getApkSigningBlock();

        /**
         * Indicates that the APK Signing Block was output as requested.
         */
        void done();
    }
}
