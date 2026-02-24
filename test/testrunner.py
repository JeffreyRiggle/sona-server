import time

allTests = []

def addTest(message, test):
    global allTests
    allTests.append({
        "message": message,
        "test": test
    })

def run():
    global allTests
    passed = 0
    failed = 0

    for t in allTests:
        try:
            print(t.get("message"))
            start = time.time()
            t.get("test")()
            passed = passed + 1
            print(t.get("message") + " PASSED in " + str(time.time() - start))
        except Exception as e:
            failed = failed + 1
            print(t.get("message") + " FAILED with " + str(e))
    
    print("COMPLETE " + str(passed) + " out of " + str(len(allTests)) + " PASSED")

    if failed > 0:
        raise Exception(str(failed) + " Tests failed!")
