#!/usr/bin/env python3

import argparse
import os

from skimage.io import imread, imsave


def parse_args():
    parser = argparse.ArgumentParser()
    parser.add_argument('-t', '--target-file', required=True)
    return parser.parse_args()


def main():
    args = parse_args()
    data = imread(args.target_file)
    imsave(args.target_file, data[1300:1450, 250:1000])


if __name__ == '__main__':
    main()
