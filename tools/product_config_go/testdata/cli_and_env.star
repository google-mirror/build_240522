# Tests rblf_env access
load("assert.star", "assert")


def test():
    assert.eq(rblf_env.TEST_ENVIRONMENT_FOO, "test_environment_foo")
    assert.fails(lambda: rblf_env.FOO_BAR_BAZ, ".*propset has no .FOO_BAR_BAZ field or method$")
    assert.eq(rblf_cli.CLI_FOO, "foo")


test()
