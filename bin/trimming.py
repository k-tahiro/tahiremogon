#!/usr/bin/env python3

import argparse
import os

import fire
from skimage.io import imread, imsave


def trim(input_file: str, output_file: str):
    data = imread(input_file)
    imsave(output_file, data[1250:1400, 750:1400])


def main():
    fire.Fire(trim)


if __name__ == '__main__':
    main()
