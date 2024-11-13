import json
import numpy
import random
import sys
import os
device = export_path = sys.argv[7]
if device == "CPU":
  # tf.config.set_visible_devices([], 'GPU')
  os.environ["CUDA_VISIBLE_DEVICES"] = "-1"
os.environ["XLA_FLAGS"] = "--xla_gpu_cuda_data_dir=/usr/lib/cuda"
sys.path.append(os.path.relpath("../"))
sys.path.append(os.path.relpath("./"))
import tensorflow as tf
from readbin import readFile
print("Tensorflow version: ", tf.__version__)
print(device)

from matplotlib import pyplot as plt

train_data_path = sys.argv[1]
test_data_path = sys.argv[2]
validation_data_path = sys.argv[3]
dictionary_path = sys.argv[4]
number_of_epochs = int(sys.argv[5])
export_path = sys.argv[6]

def getVocabsize(path):
    dictFile = open(path, "r")
    dict = json.load(dictFile)
    dictFile.close()
    return len(dict)


def readData(path, dist=None):
    data0 = readFile(os.path.join(path, "data0.bin"))
    data1 = readFile(os.path.join(path, "data1.bin"))

    x_data = data0.copy()
    x_data = numpy.append(x_data, data1, 0)
    y_data = numpy.append(numpy.zeros(len(data0), dtype=bool), numpy.ones(len(data1), dtype=bool))

    return (numpy.array(x_data), numpy.array(y_data), len(data0) / len(data1))

(x_train, y_train, distTrain) = readData(train_data_path)
(x_test, y_test, distTest) = readData(test_data_path)
(x_validate, y_validate, distVal) = readData(validation_data_path)

print("Train data length: ", len(x_train), "dist:", distTrain)
print("Test data length: ", len(x_test), "dist:", distTest)
print("Validation data length: ", len(x_validate), "dist:", distVal)

initial_bias = numpy.log(1/distTrain)
output_bias = tf.keras.initializers.Constant(initial_bias)

vocab_size = getVocabsize(dictionary_path)
print("Vocabulary size: ", vocab_size)

model = tf.keras.models.Sequential([
  tf.keras.layers.Embedding(vocab_size + 1, 8),
  tf.keras.layers.Dropout(0.5),
  tf.keras.layers.Dense(64, activation='relu'),
  tf.keras.layers.Dropout(0.5),
  tf.keras.layers.Dense(64, activation='relu'),
  tf.keras.layers.Dropout(0.5),
  tf.keras.layers.GlobalAveragePooling1D(),
  tf.keras.layers.Dropout(0.5),
#   tf.keras.layers.Dense(1)
  tf.keras.layers.Dense(1, activation='sigmoid')
#   tf.keras.layers.Dense(1, activation='sigmoid', bias_initializer=output_bias)
])

model.compile(loss= tf.keras.losses.BinaryCrossentropy(),
              optimizer='adam',
              metrics=[tf.metrics.BinaryAccuracy(name="bin_acc"),
                    #    tf.keras.metrics.TruePositives(name='tp'),
                        tf.keras.metrics.FalsePositives(name='fp'),
                        # tf.keras.metrics.TrueNegatives(name='tn'),
                        tf.keras.metrics.FalseNegatives(name='fn'), 
                       tf.keras.metrics.Precision(name='precision'),
                       tf.keras.metrics.Recall(name='recall'),
                       tf.keras.metrics.AUC(name='auc'),
                       tf.keras.metrics.AUC(name='prc', curve='PR')])

early_stopping = tf.keras.callbacks.EarlyStopping(
    monitor='val_prc', 
    verbose=1,
    patience=10,
    mode='max',
    restore_best_weights=True)


history = model.fit(x_train, y_train, validation_data=(x_validate, y_validate), epochs=number_of_epochs, callbacks=early_stopping)
model.evaluate(x_test,  y_test, verbose=2)

acc = history.history['bin_acc']
val_acc = history.history['val_bin_acc']
loss = history.history['loss']
val_loss = history.history['val_loss']
epochs = range(1, len(acc) + 1)

plt.plot(epochs, acc, 'bo', label='Training acc')
plt.plot(epochs, val_acc, 'b', label='Validation acc')
plt.title('Training and validation accuracy')
plt.legend()
plt.figure()

plt.plot(epochs, loss, 'bo', label='Training loss')
plt.plot(epochs, val_loss, 'b', label='Validation loss')
plt.title('Training and validation loss')
plt.legend()
plt.show()

model.save(export_path)