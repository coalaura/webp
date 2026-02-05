# Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# gnuplot <benchmark.gnuplot

reset

# for windows
set encoding utf8

set terminal pngcairo font "simsun,12" size 1000,800 noenhanced

set output "benchmark_result.png"
set title "WebP Decode Benchmark (Low is Better)"
set xlabel ""
set ylabel "ns/op"
set xtics rotate by -90

set key top right
set grid ytics

#set yrange [0:50]
plot \
	"benchmark_result_coalaura_webp.txt" using 3:xticlabels(1) title "coalaura/webp" with linespoints, \
	"benchmark_result_x_image_webp.txt" using 3:xticlabels(1) title "x/image/webp" with linespoints, \
