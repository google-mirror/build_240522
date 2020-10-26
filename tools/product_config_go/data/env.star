def init():
    for n in dir(rblf_env):
        print(n, "=", getattr(rblf_env, n))


init()
