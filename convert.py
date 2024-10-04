import json

# Open a JSON file and read its contents
with open('result.json', 'r') as file:
    data = json.load(file)["chats"]["list"]

print("Converting " + str(len(data)) + " Chats...")

# the chat ids you want to skip converting
skiplist = [1125107539]

export = {
    "data": [],
    "lastupdate": 0
}

for instance in data:
    ins = {
        "i": instance["id"],
        "d": []
    }
    if int(instance["id"]) in skiplist:
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
            count += 1
            ins["d"].append({
                "u": messages["id"],
                "t": messages["text"],
            })
    if count > 0:
        print("Converted " + str(count) + " Messages...")
        export["data"].append(ins)
    else:
        print("Skipped Converting because it did not contain any messages.")
print("Writing json...")
with open('result_converted.json', 'w') as json_file:
    json.dump(export, json_file)
    print("Done!")