import numpy

maxSeqLen = 1024

def readFile(filename):
    fin = open(filename, "r")
    a = numpy.fromfile(fin, dtype=numpy.uint16)
    a.shape = (a.size // maxSeqLen, maxSeqLen)
    return a