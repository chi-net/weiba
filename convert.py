import json
import hashlib

def trans(a):
    strlist = "abcdefghijklmnopqrstuvwxyz!@#$ABCDEFGHIJKLMNOPQRSTUVWXYZ%^&*1234567890()-=_+[]{}|\\:;<>?,./`~"
    resp = ""
    # print(a, end=":")
    while a > 0:
        # print(a % len(strlist), end=" ")
        resp = strlist[a % len(strlist)] + resp
        a //= len(strlist)
    # print("")
    return resp

print("Weiba exported chats convertor")
print("Repo: https://github.com/chi-net/weiba")
print("Powered by chi Network(c)2022-2024.")
print("======================================")
# Open a JSON file and read its contents
print("Opening result.json, please wait for a couple of minutes...")
with open('result.json', 'r') as file:
    data = json.load(file)["chats"]["list"]

print("Found " + str(len(data)) + " Chats...")

print("This is the chat list contained in your result.json:")
i = 0
for instance in data:
    i += 1
    print(instance["name"] + "(" + str(instance["id"]) + ")", end="; ")
    if i % 3 == 0:
        print("")

print()
print("You can input the id of the chat you would like to skip converting.[Please use space to spilt each chats you "
      "want to skip, default(press enter) is none]")

skipids = input()
# the chat ids you want to skip converting
skiplist = skipids.split(" ")

if not skiplist:
    print("You do not select any chats to exclude.")

# the digits you want to transform your sha1-encrypted message, the more, the less possibility to break
# Sometimes, 4 is enough because its possibility is less than 1/10000000
transform_digits = 4

#print(skipids, skiplist)

export = {
    "data": [],
    "lastupdate": 0
}

total = 0
for instance in data:
    ins = {
        "i": instance["id"],
        "d": []
    }
    if str(instance["id"]) in skiplist:
        print("Skip Converting " + instance["name"] + "(id:" + str(instance["id"]) + ",type:" + instance["type"] +
        ");Total " + str(
            len(instance["messages"])) + " Messages...")
        continue
    else:
        print("Now Converting " + instance["name"] + "(id:" + str(instance["id"]) + ",type:" + instance["type"] +
        ");Total " + str(
            len(instance["messages"])) + " Messages...")
    count = 0
    for messages in instance["messages"]:
        if messages["text"] != "" and isinstance(messages["text"], str):
            sha1_hash = hashlib.sha1(messages["text"].encode('utf-8')).hexdigest()
            # print(messages["text"] + ":" + sha1_hash)
            hash_int = int(sha1_hash[:16], 16)
            # print(hash_int)
            mod_result = hash_int % (92 ** transform_digits)
            count += 1
            # print(mod_result)
            ins["d"].append(trans(messages["id"]) + " " + trans(mod_result))
            total += 1
    if count > 0:
        print("Converted " + str(count) + " Messages...")
        export["data"].append(ins)
    else:
        print("Skipped Converting because it did not contain any messages.")
print("Writing json...")
with open('result_converted.json', 'w') as json_file:
    json.dump(export, json_file)
    print("Done!")
print("Successfully converted " + str(total) + " messages.")