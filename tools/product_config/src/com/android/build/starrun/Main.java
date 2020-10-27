package com.android.build.starrun;

import net.starlark.java.eval.EvalException;
import net.starlark.java.syntax.ParserInput;
import net.starlark.java.syntax.SyntaxError;

import java.io.IOException;

public class Main {
  public static void main(String[] args) throws IOException {
    ParserInput input = null;
    // parse flags
    int i;
    for (i = 0; i < args.length; i++) {
      if (!args[i].startsWith("-")) {
        break;
      }
      if (args[i].equals("-c")) {
        if (i + 1 == args.length) {
          throw new IOException("-c <cmd> flag needs an argument");
        }
        input = ParserInput.fromString(args[++i], "<command-line>");
      } else {
        throw new IOException("unknown flag: " + args[i]);
      }
    }

    if (i >= args.length) {
      if (input == null) {
        System.err.println("usage: Starlark [-c cmd | file]");
        System.exit(1);
      }
    } else {
      // positional arguments
      if (input != null) {
        throw new IOException("cannot specify both -c <cmd> and file");
      }
      if (i + 1 < args.length) {
        throw new IOException("too many positional arguments");
      }
      input = ParserInput.readFile(args[i]);
    }

    try {
      Runner runner = Runner.newRunner();
      runner.run(input);
      runner.printVariables(System.out);
      System.exit(0);
    } catch (SyntaxError.Exception ex) {
      for (SyntaxError error : ex.errors()) {
        System.err.println(error);
      }
      System.exit(1);
    } catch (EvalException ex) {
      System.err.println(ex.getMessageWithStack());
      System.exit(1);
    } catch (InterruptedException e) {
      System.err.println("Interrupted");
      System.exit(1);
    }
  }
}