import json
import sys


def read_arguments():
    l = len(sys.argv)
    if (l != 4):
        raise Exception("There must be three of arguments \
        \"python FNdwithjson.py input.json output.json words_len\"" +
                        "\n ther are " + str(l) + " of arguments.")

    input_file_name = sys.argv[1]
    output_file_name = sys.argv[2]

    try:
        words_len = int(sys.argv[3])
    except:
        raise Exception("Third argument must be a integer.")

    return input_file_name, output_file_name, words_len


if __name__ == "__main__":
    print("This is core mock program")

    input_file_name, output_file_name, words_len = read_arguments()

    output = {
        'TF': "mock-tf",
        'Score': 13,
        'Words': ["mock-words"],
        'Weights': [1.23, 1.4233]
    }

    with open(output_file_name, "w", encoding="utf-8") as f:
        json.dump(output, f)
