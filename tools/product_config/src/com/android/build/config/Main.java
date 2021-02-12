/*
 * Copyright (C) 2020 The Android Open Source Project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package com.android.build.config;

import java.io.IOException;
import java.util.List;
import java.util.Map;
import java.util.TreeMap;
import java.util.TreeSet;

public class Main {
    private final Errors mErrors;
    private final Options mOptions;

    public Main(Errors errors, Options options) {
        mErrors = errors;
        mOptions = options;
    }

    void run() {
        // Load the make configs from kati.
        Kati kati = new KatiImpl(mErrors, mOptions);
        Map<String, MakeConfig> makeConfigs = kati.loadProductConfig();
        if (makeConfigs == null || mErrors.hadError()) {
            return;
        }
        if (false) {
            for (MakeConfig makeConfig: (new TreeMap<String, MakeConfig>(makeConfigs)).values()) {
                System.out.println();
                System.out.println("=======================================");
                System.out.println("PRODUCT CONFIG FILES : " + makeConfig.getPhase());
                System.out.println("=======================================");
                makeConfig.printToStream(System.out);
            }
        }

        // Convert the make configs to generic configs.
        ConvertMakeToGenericConfig m2g = new ConvertMakeToGenericConfig(mErrors);
        GenericConfig generic = m2g.convert(makeConfigs);
        if (generic == null || mErrors.hadError()) {
            return;
        }
        if (false) {
            System.out.println("======================");
            System.out.println("REGENERATED MAKE FILES");
            System.out.println("======================");
            MakeWriter.write(System.out, generic, 0);
        }

        // TODO: Parse starlark configs.

        // Flatten all of the configs into one final set of variables.
        FlatConfig flat = FlattenConfig.flatten(mErrors, generic);
        if (flat == null || mErrors.hadError()) {
            return;
        }
        if (false) {
            System.out.println("=======================");
            System.out.println("FLATTENED VARIABLE LIST");
            System.out.println("=======================");
            MakeWriter.write(System.out, flat, 0);
        }

        // Check that the flat config is internally consistent.
        // TODO: Right now this compares against the original make based ones. When
        // we introduce starlark ones we should still check the make ones, even
        // though the results won't be 1:1.
        OutputChecker checker = new OutputChecker(flat);
        checker.reportErrors(mErrors);
        if (mErrors.hadError()) {
            return;
        }

        // Output the confiugration file for make.
        final String makeOutputFilename = kati.getWorkDirPath() + "/configuration.mk";
        try {
            MakeWriter.writeToFile(makeOutputFilename, flat,
                    MakeWriter.FLAG_WRITE_HEADER | MakeWriter.FLAG_WRITE_ANNOTATIONS);
        } catch (IOException ex) {
            mErrors.ERROR_OUTPUT.add("Error writing to file " + makeOutputFilename + ": "
                    + ex.getMessage());
        }

        // TODO: Output a soong configuration file so it doesn't have to use dumpvar, maybe
        // the same format as bazel?
        // TODO: Output a bazel configuration file.
    }

    public static void main(String[] args) {
        Errors errors = new Errors();
        int exitCode = 0;

        try {
            Options options = Options.parse(errors, args, System.getenv());
            if (errors.hadError()) {
                Options.printHelp(System.err);
                System.err.println();
                throw new CommandException();
            }

            switch (options.getAction()) {
                case DEFAULT:
                    (new Main(errors, options)).run();
                    return;
                case HELP:
                    Options.printHelp(System.out);
                    return;
            }
        } catch (CommandException | Errors.FatalException ex) {
            // These are user errors, so don't show a stack trace
            exitCode = 1;
        } catch (Throwable ex) {
            // These are programming errors in the code of this tool, so print the exception.
            // We'll try to print this.  If it's something unrecoverable, then we'll hope
            // for the best. We will still print the errors below, because they can be useful
            // for debugging.
            ex.printStackTrace(System.err);
            System.err.println();
            exitCode = 1;
        } finally {
            // Print errors and warnings
            errors.printErrors(System.err);
        }
        System.exit(exitCode);
    }
}
