package com.android.build.starrun;

import net.starlark.java.eval.EvalException;
import net.starlark.java.syntax.ParserInput;
import org.junit.Assert;
import org.junit.Test;

import java.io.ByteArrayOutputStream;
import java.io.PrintStream;

public class RunnerTest {

  String eval(String program) throws Exception {
    Runner r = Runner.newRunner();
    r.run(ParserInput.fromString(program, "data"));
    ByteArrayOutputStream bs = new ByteArrayOutputStream();
    try (PrintStream sink = new PrintStream(bs)) {
      r.printVariables(sink);
    }
    return bs.toString();
  }

  @Test
  public void testRunSuccess() throws Exception {
    Assert.assertEquals("VAR:=value\n", eval("set('VAR', 'value')"));
    Assert.assertEquals("VAR:=value\n", eval("setFinal('VAR', 'value')"));
    Assert.assertEquals("LVAR:=bar foo\n",
            eval("appendTo('LVAR', 'foo')\nappendTo('LVAR', 'bar')"));
    Assert.assertEquals("V1:=value1\nV2:=value2\n",
            eval("set('V2', 'value2')\nappendTo('V1', 'value1')"));
    Assert.assertEquals("VAR:=value\n",
            eval("loadGenerated('/bin/echo',['set(\"VAR\", \"value\")'])"));
  }

  @Test
  public void testRunErrors() throws Exception {
    try {
      eval("setFinal('VAR', 'value1')\nset('VAR', 'value2')");
      Assert.fail("Expected Evaluation exception");
    } catch (EvalException ee) {
      Assert.assertEquals("VAR has been set to 'value1' and cannot be changed", ee.getMessage());
    }
  }
}