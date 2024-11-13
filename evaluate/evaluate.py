import tensorflow as tf
import sys
import numpy
import os

device = export_path = sys.argv[2]
if device == "CPU":
  os.environ["CUDA_VISIBLE_DEVICES"] = "-1"

number_of_tokens = int(input())
tokens = []
for x in range(number_of_tokens):
    current_tokens = []
    for i in range(1024):
        current_tokens.append(int(input()))
    tokens.append(current_tokens)
tokens = numpy.array(tokens)

model_path = sys.argv[1]
model = tf.keras.models.load_model(model_path)

result = model.predict(tokens)
for res in result:
    print(res[0])
