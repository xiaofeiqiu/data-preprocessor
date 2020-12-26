# mlstock

Todo
1. ema 20,50,100,200
2. cci 100
3. Aroon 50 up, and Aroon 50 down
4. macd 20,200,100


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
Normalized Diff EMAs
Normalized CCI
Normalized Diff CCI
Normalized AroonUp and down
Normalized Diff AroonUp and down
Normalized Macd
Normalized Diff Macd
