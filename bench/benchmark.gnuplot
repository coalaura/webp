# Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# gnuplot <benchmark.gnuplot

reset

# for windows
set encoding utf8

set terminal pngcairo font "simsun,12" size 1200,850 noenhanced

set style data linespoints
set pointsize 0.8

set output "benchmark_decode.png"
set title "WebP Decode Benchmark (Low is Better)"
set xlabel ""
set ylabel "ns/op"
set xtics rotate by -90

set key outside right top
set key box opaque
set key spacing 1.2
set grid ytics

#set yrange [0:50]
plot \
	"benchmark_result_coalaura_webp.txt" using 3:xticlabels(1) title "coalaura/webp" with linespoints, \
	"benchmark_result_coalaura_webp_threaded.txt" using 3:xticlabels(1) title "coalaura/webp (threads)" with linespoints, \
	"benchmark_result_chai2010_webp.txt" using 3:xticlabels(1) title "chai2010/webp" with linespoints, \
	"benchmark_result_x_image_webp.txt" using 3:xticlabels(1) title "x/image/webp" with linespoints

set output "benchmark_encode.png"
set title "WebP Encode Benchmark (Low is Better)"

plot \
	"benchmark_result_encode_coalaura_webp.txt" using 3:xticlabels(1) title "coalaura/webp" with linespoints, \
	"benchmark_result_encode_coalaura_webp_threaded.txt" using 3:xticlabels(1) title "coalaura/webp (threads)" with linespoints, \
	"benchmark_result_encode_chai2010_webp.txt" using 3:xticlabels(1) title "chai2010/webp" with linespoints
