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

import java.io.PrintStream;
import java.util.TreeMap;

public class Options {
    public enum Action {
        DEFAULT,
        ERROR,
        HELP
    }

    private Action mAction = Action.DEFAULT;

    public Action getAction() {
        return mAction;
    }

    public static void printHelp(PrintStream out) {
        out.println("usage: product_config");
        out.println();
        out.println("OPTIONS");
        out.println("  --hide ERROR_ID...       Suppress these errors.");
        out.println("  --error ERROR_ID...      Make these ERROR_IDs fatal errors.");
        out.println("  --help -h                This message.");
        out.println("  --warning ERROR_ID...    Make these ERROR_IDs warnings.");
        out.println();
        out.println("ERRORS");
        out.println("  The following are the errors that can be controlled on the");
        out.println("  commandline with the --hide --warning --error flags.");

        TreeMap<Integer,Errors.Category> sorted = new TreeMap((new Errors()).getCategories());

        for (final Errors.Category category: sorted.values()) {
            if (category.isLevelSettable()) {
                out.println(String.format("    %-3d      %s", category.getCode(),
                category.getHelp().replace("\n", "\n             ")));
            }
        }
    }

    static class Parser {
        private class ParseException extends Exception {
            public ParseException(String message) {
                super(message);
            }
        }

        private Errors mErrors;
        private String[] mArgs;
        private Options mResult = new Options();
        private int mIndex;

        public Parser(Errors errors, String[] args) {
            mErrors = errors;
            mArgs = args;
        }

        public Options parse() {
            try {
                while (mIndex < mArgs.length) {
                    final String arg = mArgs[mIndex];

                    if ("--hide".equals(arg)) {
                        handleErrorCodes(arg, Errors.Level.HIDDEN);
                    } else if ("--error".equals(arg)) {
                        handleErrorCodes(arg, Errors.Level.ERROR);
                    } else if ("--help".equals(arg) || "-h".equals(arg)) {
                        // Help overrides all other commands if there isn't an error, but
                        // we will stop here.
                        if (mResult.mAction != Action.ERROR) {
                            mResult.mAction = Action.HELP;
                        }
                        return mResult;
                    } else if ("--warning".equals(arg)) {
                        handleErrorCodes(arg, Errors.Level.WARNING);
                    } else {
                        throw new ParseException("Unknown command line argument: " + arg);
                    }

                    mIndex++;
                }
            } catch (ParseException ex) {
                addFatalError(ex.getMessage());
            }

            return mResult;
        }

        private void addFatalError(String message) {
            mErrors.add(mErrors.ERROR_COMMAND_LINE, message);
            mResult.mAction = Action.ERROR;
        }

        private String getNextNonFlagArg() {
            if (mIndex == mArgs.length - 1) {
                return null;
            }
            if (mArgs[mIndex + 1].startsWith("-")) {
                return null;
            }
            mIndex++;
            return mArgs[mIndex];
        }

        private Integer getNextNumberArg(String errorMessage) throws ParseException {
            final String arg = getNextNonFlagArg();
            if (arg == null) {
                return null;
            }
            try {
                return Integer.parseInt(arg);
            } catch (NumberFormatException ex) {
                throw new ParseException(errorMessage + arg);
            }
        }

        private void handleErrorCodes(String arg, Errors.Level level) throws ParseException {
            boolean foundOne = false;
            while (true) {
                final Integer code = getNextNumberArg(arg + " requires numeric arguments: ");
                if (code == null) {
                    break;
                }
                final Errors.Category category = mErrors.getCategories().get(code);
                if (category == null) {
                    mErrors.add(mErrors.WARNING_UNKNOWN_COMMAND_LINE_ERROR,
                            "Unknown error code: " + code);
                    continue;
                }
                if (!category.isLevelSettable()) {
                    addFatalError("Can't set level for error " + code);
                    continue;
                }
                category.setLevel(level);
                foundOne = true;
            }
            if (!foundOne) {
                throw new ParseException(arg + " requires at least one numeric argument");
            }
        }
    }

    /**
     * Parse the arguments and return an options object.
     * <p>
     * Updates errors with the hidden / warning / error levels.
     * <p>
     * Adds errors encountered to Errors object.
     */
    public static Options parse(Errors errors, String[] args) {
        return (new Parser(errors, args)).parse();
    }
}
