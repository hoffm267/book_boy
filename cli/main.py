import argparse
from datetime import time
from io import FileIO
import os
from json import decoder
from collections import defaultdict 

def main():
    parser = argparse.ArgumentParser(description="BookBoy CLI")
    parser.add_argument("command", choices=["all", "stats"], help="What to run")
    args = parser.parse_args()
    
    #Change this to not be local file
    script_dir = os.path.dirname(os.path.abspath(__file__))
    data_path = os.path.join(script_dir, "data.json")

    with open(data_path, "r") as f:
        decodedJSON = decoder.JSONDecoder().decode(f.read())

    bookMap = defaultdict(tuple)

    #TODO make this better when internet is available  
    for n in decodedJSON:
        bookMap[n['id']] = tuple([[n['cur_time'], n['total_time']], [n['cur_page'], n['total_pages']]])

    if args.command == "all":
        print(bookMap.items())
        print("TEST")
        # this commented section will split the time when i get there
        # n = args.time.split(":")
        # if len(n) == 3:
        #     n = [int(x) for x in n]
        #     test = time(n[0], n[1], n[-1])
        #     time_obj = test.strftime("%H:%M:%S")

if __name__ == "__main__":
    main()

