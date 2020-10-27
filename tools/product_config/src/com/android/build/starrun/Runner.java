package com.android.build.starrun;

import com.google.common.collect.ImmutableMap;
import com.google.common.collect.Lists;
import net.starlark.java.annot.Param;
import net.starlark.java.annot.StarlarkMethod;
import net.starlark.java.eval.Module;
import net.starlark.java.eval.*;
import net.starlark.java.syntax.FileOptions;
import net.starlark.java.syntax.ParserInput;
import net.starlark.java.syntax.SyntaxError;

import java.io.IOException;
import java.io.InputStream;
import java.io.PrintStream;
import java.nio.charset.StandardCharsets;
import java.util.*;

public class Runner {

  private static class VariableValue {
    ArrayList<String> value;
    boolean isFrozen;
    boolean isList;

    public static VariableValue newStringValue(String value, boolean frozen) {
      VariableValue v = new VariableValue();
      v.value = Lists.newArrayList(value);
      v.isList = false;
      v.isFrozen = frozen;
      return v;
    }

    public static VariableValue newListValue(String value) {
      VariableValue v = new VariableValue();
      v.value = Lists.newArrayList(value);
      v.isList = true;
      v.isFrozen = false;
      return v;
    }
  }

  private final Bindings bindings;
  private final Map<String, VariableValue> variables;
  private final Module module;
  private final StarlarkThread.Loader loader;

  /**
   * Constructs new runner.
   *
   * @return Runner instance
   */
  public static Runner newRunner() {
    return new Runner();
  }

  private Runner() {
    this.bindings = new Bindings();
    this.variables = Collections.synchronizedMap(new HashMap<>());
    ImmutableMap.Builder<String, Object> env = ImmutableMap.builder();
    Starlark.addMethods(env, this.bindings, StarlarkSemantics.DEFAULT);
    this.module = Module.withPredeclared(StarlarkSemantics.DEFAULT, env.build());
    this.loader = new Loader();
  }

  /**
   * Reads starlark program and executes it.
   *
   * @param input program source
   * @throws SyntaxError.Exception syntax error
   * @throws EvalException         runtime error
   * @throws InterruptedException  don't ask
   */
  public void run(ParserInput input) throws SyntaxError.Exception, EvalException, InterruptedException {
    run(input, this.module);
  }

  private void run(ParserInput input, Module module)
          throws SyntaxError.Exception, EvalException, InterruptedException {
    StarlarkThread thread = new StarlarkThread(Mutability.create("interpreter"),
            StarlarkSemantics.DEFAULT);
    thread.setPrintHandler((th, msg) -> System.out.println(msg));
    thread.setLoader(loader);
    Starlark.execFile(input, FileOptions.DEFAULT, module, thread);
  }

  /**
   * Prints all variable values as makefile-style assignments, sorted by variable name.
   * List variable values are sorted, too.
   *
   * @param out print sink
   */
  public void printVariables(PrintStream out) {
    List<String> sorted = new ArrayList<>(variables.keySet());
    Collections.sort(sorted);
    for (String name : sorted) {
      VariableValue value = variables.get(name);
      if (value.isList) {
        List<String> sortedValues = value.value;
        Collections.sort(sortedValues);
        out.printf("%s:=%s\n", name, String.join(" ", sortedValues));
      } else {
        out.printf("%s:=%s\n", name, value.value.get(0));
      }
    }
  }

  private class Bindings {
    private void _set(String name, String value, boolean freeze) throws EvalException {
      if (variables.containsKey(name)) {
        VariableValue v = variables.get(name);
        if (v.isList) {
          throw Starlark.errorf("%s is a list, not a string", name);
        }
        if (v.isFrozen) {
          throw Starlark.errorf("%s has been set to '%s' and cannot be changed", name, v.value.get(0));
        }
        v.value.set(0, value);
        v.isFrozen = freeze;
      } else {
        variables.put(name, VariableValue.newStringValue(value, freeze));
      }
    }

    @StarlarkMethod(
            name = "set",
            parameters = {@Param(name = "name"), @Param(name = "value")}
    )
    public Object set(String name, String value) throws EvalException {
      _set(name, value, false);
      return Starlark.NONE;
    }

    @StarlarkMethod(
            name = "setFinal",
            parameters = {@Param(name = "name"), @Param(name = "value")}
    )
    public Object setFinal(String name, String value) throws EvalException {
      _set(name, value, true);
      return Starlark.NONE;
    }

    @StarlarkMethod(
            name = "appendTo",
            parameters = {@Param(name = "name"), @Param(name = "value")}
    )
    public Object appendTo(String name, String value) {
      if (!variables.containsKey(name)) {
        variables.put(name, VariableValue.newListValue(value));
      } else {
        variables.get(name).value.add(value);
      }
      return Starlark.NONE;
    }

    @StarlarkMethod(
            name = "loadGenerated",
            parameters = {@Param(name = "command"), @Param(name = "args")}
    )
    public Object loadGenerated(String cmd, StarlarkList<?> argList)
            throws SyntaxError.Exception, EvalException, InterruptedException {
      String[] args = new String[1 + argList.size()];
      args[0] = cmd;
      for (int i = 0; i < argList.size(); i++) {
        Object v = argList.get(i);
        if (v.getClass() == String.class) {
          args[i + 1] = (String) v;
        } else {
          throw Starlark.errorf("args should be list of strings");
        }
      }
      ProcessBuilder pb = new ProcessBuilder(args);
      String cmdline = String.join(" ", args);
      try {
        Process process = pb.start();
        int rc = process.waitFor();
        if (rc != 0) {
          try (InputStream err = process.getErrorStream()) {
            throw Starlark.errorf("'%s' failed with rc=%d:\n%s\n",
                    cmdline, rc,
                    new String(err.readAllBytes(), StandardCharsets.UTF_8));
          }
        }
        try (InputStream out = process.getInputStream()) {
          run(ParserInput.fromUTF8(out.readAllBytes(), String.format("$(%s)", cmdline)));
        }
      } catch (IOException ioe) {
        throw Starlark.errorf("'%s' failed: %s", cmdline, ioe.getMessage());
      }
      return Starlark.NONE;
    }
  }

  private class Loader implements StarlarkThread.Loader {
    private final Map<String, Module> loadedModules = new HashMap<>();

    @Override
    public Module load(String moduleName) {
      try {
        if (loadedModules.containsKey(moduleName)) {
          Module oldModule = loadedModules.get(moduleName);
          if (oldModule != null) {
            return oldModule;
          }
          System.err.printf("'%s' is called recursively\n", moduleName);
          return null;
        } else {
          loadedModules.put(moduleName, null);
          // TODO(asmundak): is this correct, or should a new module
          // be created and passed to run?
          run(ParserInput.readFile(moduleName), Runner.this.module);
          loadedModules.put(moduleName, Runner.this.module);
          return module;
        }
      } catch (SyntaxError.Exception | EvalException | InterruptedException e) {
        // Nothing
      } catch (IOException e) {
        System.err.println(e.getMessage());
      }
      return null;
    }
  }
}
