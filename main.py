import argparse
from datetime import time
from io import FileIO
import os
from json import decoder
from collections import defaultdict 

def main():
    parser = argparse.ArgumentParser(description="Script with datetime input.")
    parser.add_argument('all')
    args = parser.parse_args()
    file_dir = FileIO(str(os.curdir) + "/data.json")
    decodedJSON = decoder.JSONDecoder().decode(file_dir.read().decode())
    bookMap = defaultdict(tuple)

    #TODO make this better when internet is available  
    for n in decodedJSON:
        bookMap[n['id']] = tuple([[n['cur_time'], n['total_time']], [n['cur_page'], n['total_pages']]])

    try:
        if args.all:
            print(bookMap.items())
            # this commented section will split the time when i get there
            # n = args.time.split(":")
            # if len(n) == 3:
            #     n = [int(x) for x in n]
            #     test = time(n[0], n[1], n[-1])
            #     time_obj = test.strftime("%H:%M:%S")
    except ValueError:
        print("Error: Please use date format YYYY-MM-DD.")

if __name__ == "__main__":
    main()

