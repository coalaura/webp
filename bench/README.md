Benchmark
=========

Decode
------

![](benchmark_decode.png)

Encode
------

![](benchmark_encode.png)

I ran all benchmarks (including the plotted ones) on a CachyOS laptop with no other software running to get as clean results as possible.

```
go test -bench=.
goos: linux
goarch: amd64
pkg: github.com/coalaura/webp/bench
cpu: AMD Ryzen 7 7840U w/ Radeon(TM) 780M Graphics
BenchmarkDecode_1_webp_a_coalaura_webp-16                                               	    1026	   1104062 ns/op
BenchmarkDecode_1_webp_a_x_image_webp-16                                                	     499	   2325213 ns/op
BenchmarkDecode_1_webp_a_coalaura_webp_threaded-16                                      	    1011	   1105059 ns/op
BenchmarkEncode_1_webp_a_coalaura_webp-16                                               	      87	  13341902 ns/op
BenchmarkEncode_1_webp_a_coalaura_webp_threaded-16                                      	     122	   9850120 ns/op
BenchmarkDecode_1_webp_ll_coalaura_webp-16                                              	     985	   1147442 ns/op
BenchmarkDecode_1_webp_ll_x_image_webp-16                                               	     422	   2819762 ns/op
BenchmarkDecode_1_webp_ll_coalaura_webp_threaded-16                                     	     985	   1139376 ns/op
BenchmarkEncode_1_webp_ll_coalaura_webp-16                                              	      88	  13327835 ns/op
BenchmarkEncode_1_webp_ll_coalaura_webp_threaded-16                                     	     126	   9493497 ns/op
BenchmarkDecode_2_webp_a_coalaura_webp-16                                               	     918	   1203540 ns/op
BenchmarkDecode_2_webp_a_x_image_webp-16                                                	     486	   2423699 ns/op
BenchmarkDecode_2_webp_a_coalaura_webp_threaded-16                                      	     927	   1205751 ns/op
BenchmarkEncode_2_webp_a_coalaura_webp-16                                               	      51	  21057896 ns/op
BenchmarkEncode_2_webp_a_coalaura_webp_threaded-16                                      	      62	  16709707 ns/op
BenchmarkDecode_2_webp_ll_coalaura_webp-16                                              	    1756	    643857 ns/op
BenchmarkDecode_2_webp_ll_x_image_webp-16                                               	     588	   1998652 ns/op
BenchmarkDecode_2_webp_ll_coalaura_webp_threaded-16                                     	    1731	    637021 ns/op
BenchmarkEncode_2_webp_ll_coalaura_webp-16                                              	      49	  20551112 ns/op
BenchmarkEncode_2_webp_ll_coalaura_webp_threaded-16                                     	      66	  15945768 ns/op
BenchmarkDecode_3_webp_a_coalaura_webp-16                                               	     369	   3213008 ns/op
BenchmarkDecode_3_webp_a_x_image_webp-16                                                	     147	   8004813 ns/op
BenchmarkDecode_3_webp_a_coalaura_webp_threaded-16                                      	     362	   3226063 ns/op
BenchmarkEncode_3_webp_a_coalaura_webp-16                                               	      15	  68789926 ns/op
BenchmarkEncode_3_webp_a_coalaura_webp_threaded-16                                      	      20	  52887881 ns/op
BenchmarkDecode_3_webp_ll_coalaura_webp-16                                              	     390	   3036193 ns/op
BenchmarkDecode_3_webp_ll_x_image_webp-16                                               	     156	   7657944 ns/op
BenchmarkDecode_3_webp_ll_coalaura_webp_threaded-16                                     	     382	   3123139 ns/op
BenchmarkEncode_3_webp_ll_coalaura_webp-16                                              	      16	  65348635 ns/op
BenchmarkEncode_3_webp_ll_coalaura_webp_threaded-16                                     	      21	  50607538 ns/op
BenchmarkDecode_4_webp_a_coalaura_webp-16                                               	    1276	    899823 ns/op
BenchmarkDecode_4_webp_a_x_image_webp-16                                                	     618	   1903738 ns/op
BenchmarkDecode_4_webp_a_coalaura_webp_threaded-16                                      	    1280	    899894 ns/op
BenchmarkEncode_4_webp_a_coalaura_webp-16                                               	     150	   8037875 ns/op
BenchmarkEncode_4_webp_a_coalaura_webp_threaded-16                                      	     196	   6030981 ns/op
BenchmarkDecode_4_webp_ll_coalaura_webp-16                                              	    1672	    663152 ns/op
BenchmarkDecode_4_webp_ll_x_image_webp-16                                               	     646	   1826115 ns/op
BenchmarkDecode_4_webp_ll_coalaura_webp_threaded-16                                     	    1701	    662781 ns/op
BenchmarkEncode_4_webp_ll_coalaura_webp-16                                              	     156	   7642947 ns/op
BenchmarkEncode_4_webp_ll_coalaura_webp_threaded-16                                     	     207	   5759319 ns/op
BenchmarkDecode_5_webp_a_coalaura_webp-16                                               	     513	   2271640 ns/op
BenchmarkDecode_5_webp_a_x_image_webp-16                                                	     250	   4737550 ns/op
BenchmarkDecode_5_webp_a_coalaura_webp_threaded-16                                      	     516	   2275517 ns/op
BenchmarkEncode_5_webp_a_coalaura_webp-16                                               	      40	  26856318 ns/op
BenchmarkEncode_5_webp_a_coalaura_webp_threaded-16                                      	      49	  21986384 ns/op
BenchmarkDecode_5_webp_ll_coalaura_webp-16                                              	    1028	   1115415 ns/op
BenchmarkDecode_5_webp_ll_x_image_webp-16                                               	     411	   2903927 ns/op
BenchmarkDecode_5_webp_ll_coalaura_webp_threaded-16                                     	    1046	   1106766 ns/op
BenchmarkEncode_5_webp_ll_coalaura_webp-16                                              	      40	  26564519 ns/op
BenchmarkEncode_5_webp_ll_coalaura_webp_threaded-16                                     	      54	  21681097 ns/op
BenchmarkDecode_blue_purple_pink_large_lossless_coalaura_webp-16                        	     428	   2750279 ns/op
BenchmarkDecode_blue_purple_pink_large_lossless_x_image_webp-16                         	     176	   6776859 ns/op
BenchmarkDecode_blue_purple_pink_large_lossless_coalaura_webp_threaded-16               	     430	   2753747 ns/op
BenchmarkEncode_blue_purple_pink_large_lossless_coalaura_webp-16                        	       6	 187788772 ns/op
BenchmarkEncode_blue_purple_pink_large_lossless_coalaura_webp_threaded-16               	       6	 188290472 ns/op
BenchmarkDecode_blue_purple_pink_large_no_filter_lossy_coalaura_webp-16                 	    1102	   1051506 ns/op
BenchmarkDecode_blue_purple_pink_large_no_filter_lossy_x_image_webp-16                  	     451	   2635391 ns/op
BenchmarkDecode_blue_purple_pink_large_no_filter_lossy_coalaura_webp_threaded-16        	    1113	   1050897 ns/op
BenchmarkEncode_blue_purple_pink_large_no_filter_lossy_coalaura_webp-16                 	      60	  17682497 ns/op
BenchmarkEncode_blue_purple_pink_large_no_filter_lossy_coalaura_webp_threaded-16        	      62	  17294720 ns/op
BenchmarkDecode_blue_purple_pink_large_normal_filter_lossy_coalaura_webp-16             	    1048	   1135876 ns/op
BenchmarkDecode_blue_purple_pink_large_normal_filter_lossy_x_image_webp-16              	     286	   4107205 ns/op
BenchmarkDecode_blue_purple_pink_large_normal_filter_lossy_coalaura_webp_threaded-16    	    1029	   1154684 ns/op
BenchmarkEncode_blue_purple_pink_large_normal_filter_lossy_coalaura_webp-16             	      60	  18130297 ns/op
BenchmarkEncode_blue_purple_pink_large_normal_filter_lossy_coalaura_webp_threaded-16    	      61	  17829265 ns/op
BenchmarkDecode_blue_purple_pink_large_simple_filter_lossy_coalaura_webp-16             	    1052	   1108180 ns/op
BenchmarkDecode_blue_purple_pink_large_simple_filter_lossy_x_image_webp-16              	     375	   3186629 ns/op
BenchmarkDecode_blue_purple_pink_large_simple_filter_lossy_coalaura_webp_threaded-16    	    1060	   1099017 ns/op
BenchmarkEncode_blue_purple_pink_large_simple_filter_lossy_coalaura_webp-16             	      58	  17792970 ns/op
BenchmarkEncode_blue_purple_pink_large_simple_filter_lossy_coalaura_webp_threaded-16    	      61	  17462247 ns/op
BenchmarkDecode_blue_purple_pink_lossless_coalaura_webp-16                              	    5152	    200536 ns/op
BenchmarkDecode_blue_purple_pink_lossless_x_image_webp-16                               	    1970	    588040 ns/op
BenchmarkDecode_blue_purple_pink_lossless_coalaura_webp_threaded-16                     	    5005	    200872 ns/op
BenchmarkEncode_blue_purple_pink_lossless_coalaura_webp-16                              	      33	  32597917 ns/op
BenchmarkEncode_blue_purple_pink_lossless_coalaura_webp_threaded-16                     	      33	  32685282 ns/op
BenchmarkDecode_blue_purple_pink_lossy_coalaura_webp-16                                 	   14517	     82550 ns/op
BenchmarkDecode_blue_purple_pink_lossy_x_image_webp-16                                  	    5920	    171293 ns/op
BenchmarkDecode_blue_purple_pink_lossy_coalaura_webp_threaded-16                        	   14480	     82460 ns/op
BenchmarkEncode_blue_purple_pink_lossy_coalaura_webp-16                                 	     850	   1387369 ns/op
BenchmarkEncode_blue_purple_pink_lossy_coalaura_webp_threaded-16                        	     819	   1437099 ns/op
BenchmarkDecode_gopher_doc_1bpp_lossless_coalaura_webp-16                               	   70616	     16582 ns/op
BenchmarkDecode_gopher_doc_1bpp_lossless_x_image_webp-16                                	   34231	     35110 ns/op
BenchmarkDecode_gopher_doc_1bpp_lossless_coalaura_webp_threaded-16                      	   70950	     16794 ns/op
BenchmarkEncode_gopher_doc_1bpp_lossless_coalaura_webp-16                               	    2943	    365942 ns/op
BenchmarkEncode_gopher_doc_1bpp_lossless_coalaura_webp_threaded-16                      	    3135	    366029 ns/op
BenchmarkDecode_gopher_doc_2bpp_lossless_coalaura_webp-16                               	   60405	     19267 ns/op
BenchmarkDecode_gopher_doc_2bpp_lossless_x_image_webp-16                                	   26905	     44383 ns/op
BenchmarkDecode_gopher_doc_2bpp_lossless_coalaura_webp_threaded-16                      	   60348	     19316 ns/op
BenchmarkEncode_gopher_doc_2bpp_lossless_coalaura_webp-16                               	    2050	    545483 ns/op
BenchmarkEncode_gopher_doc_2bpp_lossless_coalaura_webp_threaded-16                      	    2082	    536788 ns/op
BenchmarkDecode_gopher_doc_4bpp_lossless_coalaura_webp-16                               	   52784	     22131 ns/op
BenchmarkDecode_gopher_doc_4bpp_lossless_x_image_webp-16                                	   19488	     61094 ns/op
BenchmarkDecode_gopher_doc_4bpp_lossless_coalaura_webp_threaded-16                      	   54132	     22345 ns/op
BenchmarkEncode_gopher_doc_4bpp_lossless_coalaura_webp-16                               	    1447	    784414 ns/op
BenchmarkEncode_gopher_doc_4bpp_lossless_coalaura_webp_threaded-16                      	    1449	    781499 ns/op
BenchmarkDecode_gopher_doc_8bpp_lossless_coalaura_webp-16                               	   38144	     31569 ns/op
BenchmarkDecode_gopher_doc_8bpp_lossless_x_image_webp-16                                	   13437	     89293 ns/op
BenchmarkDecode_gopher_doc_8bpp_lossless_coalaura_webp_threaded-16                      	   38246	     31356 ns/op
BenchmarkEncode_gopher_doc_8bpp_lossless_coalaura_webp-16                               	     680	   1753466 ns/op
BenchmarkEncode_gopher_doc_8bpp_lossless_coalaura_webp_threaded-16                      	     684	   1747623 ns/op
BenchmarkDecode_photo_lossy_coalaura_webp-16                                            	      62	  19258733 ns/op
BenchmarkDecode_photo_lossy_x_image_webp-16                                             	      12	  95948029 ns/op
BenchmarkDecode_photo_lossy_coalaura_webp_threaded-16                                   	      54	  19185629 ns/op
BenchmarkEncode_photo_lossy_coalaura_webp-16                                            	       2	 538554082 ns/op
BenchmarkEncode_photo_lossy_coalaura_webp_threaded-16                                   	       2	 520435850 ns/op
BenchmarkDecode_tux_lossless_coalaura_webp-16                                           	    1598	    746417 ns/op
BenchmarkDecode_tux_lossless_x_image_webp-16                                            	     548	   2182580 ns/op
BenchmarkDecode_tux_lossless_coalaura_webp_threaded-16                                  	    1653	    750282 ns/op
BenchmarkEncode_tux_lossless_coalaura_webp-16                                           	       6	 185406582 ns/op
BenchmarkEncode_tux_lossless_coalaura_webp_threaded-16                                  	       6	 184004584 ns/op
BenchmarkDecode_video_001_lossy_coalaura_webp-16                                        	   11923	    101054 ns/op
BenchmarkDecode_video_001_lossy_x_image_webp-16                                         	    5722	    211738 ns/op
BenchmarkDecode_video_001_lossy_coalaura_webp_threaded-16                               	   11941	    100283 ns/op
BenchmarkEncode_video_001_lossy_coalaura_webp-16                                        	     768	   1522229 ns/op
BenchmarkEncode_video_001_lossy_coalaura_webp_threaded-16                               	     758	   1544305 ns/op
BenchmarkDecode_video_001_coalaura_webp-16                                              	   11872	    100821 ns/op
BenchmarkDecode_video_001_x_image_webp-16                                               	    4911	    211447 ns/op
BenchmarkDecode_video_001_coalaura_webp_threaded-16                                     	   11978	     99518 ns/op
BenchmarkEncode_video_001_coalaura_webp-16                                              	     784	   1520681 ns/op
BenchmarkEncode_video_001_coalaura_webp_threaded-16                                     	     768	   1540773 ns/op
BenchmarkDecode_yellow_rose_lossless_coalaura_webp-16                                   	     894	   1310757 ns/op
BenchmarkDecode_yellow_rose_lossless_x_image_webp-16                                    	     406	   2940964 ns/op
BenchmarkDecode_yellow_rose_lossless_coalaura_webp_threaded-16                          	     912	   1306015 ns/op
BenchmarkEncode_yellow_rose_lossless_coalaura_webp-16                                   	       6	 171205331 ns/op
BenchmarkEncode_yellow_rose_lossless_coalaura_webp_threaded-16                          	       6	 169581062 ns/op
BenchmarkDecode_yellow_rose_lossy_with_alpha_coalaura_webp-16                           	    1315	    850581 ns/op
BenchmarkDecode_yellow_rose_lossy_with_alpha_x_image_webp-16                            	     507	   2332309 ns/op
BenchmarkDecode_yellow_rose_lossy_with_alpha_coalaura_webp_threaded-16                  	    1333	    857808 ns/op
BenchmarkEncode_yellow_rose_lossy_with_alpha_coalaura_webp-16                           	      86	  13347517 ns/op
BenchmarkEncode_yellow_rose_lossy_with_alpha_coalaura_webp_threaded-16                  	     121	   9790344 ns/op
BenchmarkDecode_yellow_rose_lossy_coalaura_webp-16                                      	    1520	    744428 ns/op
BenchmarkDecode_yellow_rose_lossy_x_image_webp-16                                       	     637	   1850292 ns/op
BenchmarkDecode_yellow_rose_lossy_coalaura_webp_threaded-16                             	    1530	    750788 ns/op
BenchmarkEncode_yellow_rose_lossy_coalaura_webp-16                                      	     142	   8329520 ns/op
BenchmarkEncode_yellow_rose_lossy_coalaura_webp_threaded-16                             	     145	   8229983 ns/op
PASS
ok  	github.com/coalaura/webp/bench	193.463s
```
