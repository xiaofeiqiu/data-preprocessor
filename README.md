# mlstock

Todo
1. ema 20,50,100,200
2. cci 100
3. Aroon 50 up, and Aroon 50 down
4. macd 20,200,200


Raw data
ema20
ema50
ema100
ema200


Normalized calculator
===============
ema20_N ema50_N ema100_N ema200_N CCI_N aroon-up_N aroon-down_N macd_N macd_sig

emaX_N = (emaX - emaMin)/(emaMax-emaMin)  

CCI_N= cci/100

Aroon_N = Aroon / 100

MAC= hist/sig

==============
Trend Calculation
diff = today-yesterday
diffN = (diff - diffMin)/(diffMax - diffMin)

=================================
Final input
Normalized EMAs
Normalized Diff EMAs (14 days avg diff)
Normalized CCI
Normalized Diff CCI (14 days avg diff)
Normalized AroonUp and down (14 days avg diff)
Normalized Diff AroonUp and down (14 days avg diff)
Normalized Macd
Normalized Diff Macd (14 days avg diff)
Normalized Diff macd hist (14 days avg diff)

Out put:
sell -1
hold 0
buy 1

tagging:
change = ((today+14)' price - today' price) / today's price
if ABS(change) > 3%, sell or buy
else hold


ANN
use CategoricalCrossentropy for loss funciton
use softmax for output activation func

https://machinelearningmastery.com/multi-class-classification-tutorial-keras-deep-learning-library/
https://www.youtube.com/watch?v=oOSXQP7C7ck