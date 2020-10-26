def init():
    print("Before error")
    rblf_error("runtime error")
    print("Should not be shown")


init()
