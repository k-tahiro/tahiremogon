#!/usr/bin/env python3
import argparse

from PIL import Image
import torch
import torch.nn as nn
from torchvision import models, transforms


def get_transformer(input_size: int):
    return transforms.Compose([
        transforms.Resize(input_size),
        transforms.CenterCrop(input_size),
        transforms.ToTensor(),
        transforms.Normalize([0.485, 0.456, 0.406], [0.229, 0.224, 0.225])
    ])


def set_parameter_requires_grad(model, feature_extracting):
    if feature_extracting:
        for param in model.parameters():
            param.requires_grad = False


def load_model(model_file):
    device = torch.device("cpu")
    model_ft = models.squeezenet1_0(pretrained=True)
    set_parameter_requires_grad(model_ft, True)
    model_ft.classifier[1] = nn.Conv2d(512, 2, kernel_size=(1,1), stride=(1,1))
    model_ft.num_classes = 2
    model_ft.load_state_dict(torch.load(model_file, map_location=device))
    model_ft.eval()
    return model_ft


def parse_args():
    parser = argparse.ArgumentParser()
    parser.add_argument('input_file')
    parser.add_argument('-i', '--input-size', type=int, default=224)
    parser.add_argument('-m', '--model-file', default='model.pth')
    return parser.parse_args()


def main():
    args = parse_args()

    im = Image.open(args.input_file)
    transform = get_transformer(args.input_size)
    data = torch.stack([transform(im)])

    model = load_model(args.model_file)
    outputs = model(data)
    _, preds = torch.max(outputs, 1)
    
    result = preds.numpy()[0]
    print(result)


if __name__ == '__main__':
    main()
