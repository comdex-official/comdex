package v6

// jq '.delegation_responses | map({address:.delegation.delegator_address,amount:((.balance.amount | tonumber)*0.05*((0.42/365)*13+1) | floor) | tostring})' DAN.JSON > to_mint.json

// Slash was 5%
// Lost APR is 36% for 20 days

var recordsJSONString = `[
  {
      "address": "comdex1qp59k44d9338qqc2vx4cc8rm0tnfwwvh2t70fk",
      "amount": 348630
  },
  {
      "address": "comdex1qrh36gnrgttj42xxqpmrehnm6jx2lgj5q96vqp",
      "amount": 278904
  },
  {
      "address": "comdex1qyyqykjmamxdr2m9lvm3r3470s0hcaarspfv39",
      "amount": 550835
  },
  {
      "address": "comdex1qffz59msmkudvjjzh7v20xtcscnty9yy52em88",
      "amount": 348630
  },
  {
      "address": "comdex1qv4p64e503cymydz0kearj4k3hf43xull0z5cy",
      "amount": 2091
  },
  {
      "address": "comdex1qwra375ylerzzlj5e0w7zezvawp9dz403zkujy",
      "amount": 1136534
  },
  {
      "address": "comdex1qs3emsjfgm0y6xkm9j2jyam7kjtfkhe2tjhy38",
      "amount": 42532
  },
  {
      "address": "comdex1q30y9yqcjt8s7w8c3sta3ajgyk6lyn45ct5ld8",
      "amount": 786509
  },
  {
      "address": "comdex1qerwnwm7gy25d2sp4e70n00qhekltk5cx5nyg4",
      "amount": 26449383
  },
  {
      "address": "comdex1qazqc3wr40y3z5gnmydrpxnz5zvxlqwpx7mar6",
      "amount": 725150
  },
  {
      "address": "comdex1qly96sy3yklkwv6cf4he5ves4qu8yn58a83ztp",
      "amount": 697260
  },
  {
      "address": "comdex1pzn40yxntlnfta7tkgrs7q6fd693em4kk4usap",
      "amount": 69726
  },
  {
      "address": "comdex1pydcjl95vl2c6rvlny977mcw8f9p8c79p7adnh",
      "amount": 46228497
  },
  {
      "address": "comdex1px33xnq4mjsdljxkmyermuxnhatwt62x80pd5s",
      "amount": 2022054
  },
  {
      "address": "comdex1pg3uzk7v0ddhxen8gze5hq0z89far4htx04yjq",
      "amount": 3523811
  },
  {
      "address": "comdex1pfdg20wh4vfcfaxn72qegqxlq7cvllx3nxfw9s",
      "amount": 38311
  },
  {
      "address": "comdex1pd5v3d5mx52634y5vu3fxy8k5380fenf6nr87r",
      "amount": 1743150
  },
  {
      "address": "comdex1pduq0yazfp6yzhp84wgrhxgv8v5mzmkn5acv8l",
      "amount": 446381
  },
  {
      "address": "comdex1p0jye26qt94xy45jtadmaguf4kyrur2s4jlx5k",
      "amount": 1917465
  },
  {
      "address": "comdex1p332nd0skn9ct479k62jraq00qfqsfflpytgzh",
      "amount": 662397
  },
  {
      "address": "comdex1pklmua5j9245z9rn32xdu7xlhysxe6g8q2rj2x",
      "amount": 1255068
  },
  {
      "address": "comdex1zrxtzvufnqn2ansjqllwf70d26ryqzpz6s846u",
      "amount": 7042328
  },
  {
      "address": "comdex1zf89488lx56vt8h6etuhuqand3gczrm5xzme2y",
      "amount": 6972
  },
  {
      "address": "comdex1z2xhjfy3y9adus6t50u4ha02dj7j7cnwkjfmk4",
      "amount": 5717534
  },
  {
      "address": "comdex1zdjgc2zxmq068kqvncptmlghf9rvg6lqjlw75t",
      "amount": 146424
  },
  {
      "address": "comdex1zdek2gepft4mm5m5m56z4sdg9e67ddkcqw4dyg",
      "amount": 12543712
  },
  {
      "address": "comdex1zwkanqgznsge9nrawt5m9wm8w5ngp732mdzpcw",
      "amount": 7739589
  },
  {
      "address": "comdex1zsfljc0lphktsjwry83qdw2ktvkyhq7fx8wrcm",
      "amount": 2025666
  },
  {
      "address": "comdex1z3l6sqf75zw7stcj0mm9ysrm3vn0a2u6q3ah6y",
      "amount": 46019
  },
  {
      "address": "comdex1zcfesl9uzazwzd76ual4lxknd9s5q3c0rgesa0",
      "amount": 6972
  },
  {
      "address": "comdex1ze2ye5u5k3qdlexvt2e0nn0508p040944qwels",
      "amount": 35423
  },
  {
      "address": "comdex1zm7vax6zax6k5t2nsu7au0undz0axqwu8s0c7w",
      "amount": 43927
  },
  {
      "address": "comdex1zat7tefp8exw65cqqmy4xs0tgl3s7ezy0qvewe",
      "amount": 334684
  },
  {
      "address": "comdex1z7trpse5s9d70y5hgnc2rhlhjwy3dkkus8ls0g",
      "amount": 987320
  },
  {
      "address": "comdex1rzm25ul3vmw5djnwzhgnthuqx0zvqyccave5gd",
      "amount": 418356
  },
  {
      "address": "comdex1rrrmtxt8gxaxt5ffrvnusm7hrq0ul2z2hmj5tr",
      "amount": 697260
  },
  {
      "address": "comdex1rrk0hgawmg5lzlkf948f395ksqmhr867f4kgem",
      "amount": 6972
  },
  {
      "address": "comdex1r8ez7ah24ls66ayyd22uz8592yu2x4jktag9kl",
      "amount": 7669863
  },
  {
      "address": "comdex1rgh67ez5vva7sdztqm34dprwxhz7mgg85wvfls",
      "amount": 139452
  },
  {
      "address": "comdex1r2vr672lweqh32ttthptful2zn2ssrs4e9r850",
      "amount": 313767
  },
  {
      "address": "comdex1r2j4xfdcxtrs20lecjhapz5vm5y5u0wq3jzzuy",
      "amount": 1133047
  },
  {
      "address": "comdex1rth8szh4jf88yqfrqp3xmfdacdja6cymu5am4a",
      "amount": 418356
  },
  {
      "address": "comdex1r02hqjuq0hxqw8v3wsl2rfdlgcc23sr6q7r8m7",
      "amount": 223707
  },
  {
      "address": "comdex1rhd7lywrqh6hg9qgmv3xqz3n6np42mj0squu3c",
      "amount": 36271817
  },
  {
      "address": "comdex1rertkh60g8t7j2v6fdsphcdvvfgxj2e4n6f5r7",
      "amount": 43927
  },
  {
      "address": "comdex1r6dgn9q9da65d7c3gv8ga99lpzv9a973k8c8up",
      "amount": 22144986
  },
  {
      "address": "comdex1r6uavvun5hretydv266dj60wkmedyeg4qjkq6f",
      "amount": 2689332
  },
  {
      "address": "comdex1rmgujl9ek0v8ghwp9vve0s5jsrh4evq4kmull5",
      "amount": 22892679
  },
  {
      "address": "comdex1rlv9jdyaurqj8hy77d8wslwfxxlesgmlldqr3r",
      "amount": 97616
  },
  {
      "address": "comdex1yp923424vva3eeprfvmjzztry36dcc30mf0wzl",
      "amount": 20917
  },
  {
      "address": "comdex1yru7uear3j6uxpqa46pvdal7wnvpdekulh7ug9",
      "amount": 2091780
  },
  {
      "address": "comdex1y946vr3nxmxp4ueuunt4hazmvf5srdk9a3uwcv",
      "amount": 4253287
  },
  {
      "address": "comdex1y8ndg3uh0trmnpmyxzpm4jpkcwz6zahtv6dhac",
      "amount": 69726
  },
  {
      "address": "comdex1yfaqm2zu42aqgk6t2tedfx4u2z3flkcgc5lqcq",
      "amount": 6993520
  },
  {
      "address": "comdex1yvpke2wu5ce7y9tahmns48mvl849y7jqkeg2er",
      "amount": 3486301
  },
  {
      "address": "comdex1yvehzsac0gd8vxge3q3m3le8a6h63rl885xqhs",
      "amount": 3486
  },
  {
      "address": "comdex1ydqnfczdpuc8xpk8q77urruztzpghmy8k6kj92",
      "amount": 13945
  },
  {
      "address": "comdex1y5dafscf8xz8qazzl8nchhlxcyr375t9p44v95",
      "amount": 209178
  },
  {
      "address": "comdex1ye0y4vnhzndgptrgjf597dwkcz98u5ye4cvje8",
      "amount": 181287
  },
  {
      "address": "comdex1ym9sllrjs7fs7urv47czxu4ph8dw40j0lcshn3",
      "amount": 29982
  },
  {
      "address": "comdex1yu0qxd3q63kgmpecfh6s7p7lfqw0mc2dnxjev2",
      "amount": 48808
  },
  {
      "address": "comdex1yu63kmfau0dtuv7p8dfarr4asm9jjxqu83j7dy",
      "amount": 139452
  },
  {
      "address": "comdex1y7s6xcqs5dul55rhfygrlq029jea0pfsj7fga2",
      "amount": 2684452
  },
  {
      "address": "comdex1yl2p7f8glqmwq2gpz4tlu6x02vxjksvz20k62v",
      "amount": 16166676
  },
  {
      "address": "comdex19ru6rtn59w9kd2v2687vv8sl98w9wurw2dk4yh",
      "amount": 383828
  },
  {
      "address": "comdex19r7a8v2dh2ygp288jcltq9c90apctmxvkk03z6",
      "amount": 6972
  },
  {
      "address": "comdex19tm2nnw960kwjdd2y48cvce93z8wjz4argl73r",
      "amount": 221905
  },
  {
      "address": "comdex193jzygzen74jt6rkkr485ud7numpunsdn3ch5h",
      "amount": 2370684
  },
  {
      "address": "comdex195lfnvs7jeevj6prn5a9ghr563s6tmzl65fr8n",
      "amount": 3648763
  },
  {
      "address": "comdex19mwsw0a5f2kk5w6pfjwm4y9k4l2h9c4t6c44ef",
      "amount": 1115616
  },
  {
      "address": "comdex19mn3w3zy3q4fgq8fdtdqukfetdrjl52x57x7gh",
      "amount": 278904
  },
  {
      "address": "comdex19mc5kzqfqfn8h2fay90lpaadfchtmexvhqe0qw",
      "amount": 906438
  },
  {
      "address": "comdex19lp0rt6tg2mc4urft32mr5jag2eg7lty7tflkg",
      "amount": 721664
  },
  {
      "address": "comdex1xq58ug0l2mumnx4fepyxlwk6twmx8ac6x5273f",
      "amount": 3242260
  },
  {
      "address": "comdex1xpzw5seg4u28sj5z4dre4h03naagvddx009wm2",
      "amount": 11156164
  },
  {
      "address": "comdex1xpyglnn8m2w270d5ttjd7muk6nx0lpnjnkf520",
      "amount": 1394
  },
  {
      "address": "comdex1xzsu93xry7wze2xfggukprhvpcq8vsrdav5slm",
      "amount": 9517602
  },
  {
      "address": "comdex1xymp5mz5rlyd4z7w5v29q9gfl4h4kjjf0z5cnw",
      "amount": 822767
  },
  {
      "address": "comdex1xxzp5xxj2mdu0ssn52p6xrp0895djsrekpss2t",
      "amount": 43032381
  },
  {
      "address": "comdex1xt5kvt65ctp2dzpfgtz6su5l4pu25jpganw2lg",
      "amount": 299821
  },
  {
      "address": "comdex1x3qede7qlhen0k6wnn4dpq0mjv9ssv75klfr64",
      "amount": 139452
  },
  {
      "address": "comdex1x5x0lnsv5ftfdce7zhdwurhlt327f38e6p7vvu",
      "amount": 766986
  },
  {
      "address": "comdex1x5mtgrwvwejwvmu2jujlr63hn4q2j67d7eg2h5",
      "amount": 1394520
  },
  {
      "address": "comdex1xk5cmldhsz3wlev70m7l6p50e7fjqgvdpk459t",
      "amount": 330709
  },
  {
      "address": "comdex1xhnls8dtd88356g0lgg38kn6r2xgpr850lwvup",
      "amount": 869895
  },
  {
      "address": "comdex1x7hm7ap8nea5qt4q7kvdz0qfwfrmyjhnc6gpdc",
      "amount": 15116602
  },
  {
      "address": "comdex18zph6wtfave8xc6q84r5fxzfqzsh6h2uejs3sy",
      "amount": 7401417
  },
  {
      "address": "comdex18rfdktgpxeepy52n57zmu64y6wz0pxqyekxhvl",
      "amount": 348630
  },
  {
      "address": "comdex18rlnm5tc2zj6whmpu788lzk8fnusc3ppvcnlkc",
      "amount": 383493
  },
  {
      "address": "comdex18gry7apt60sndrw6phjyw736ayquef7waas2hd",
      "amount": 906438
  },
  {
      "address": "comdex18d9yadmu3yz94prcz2unez0u6zcu6frxwzut0s",
      "amount": 69726
  },
  {
      "address": "comdex18w92ch5dryrcq40ns23e3tjlpd2vzrycgxwfcr",
      "amount": 69726
  },
  {
      "address": "comdex1830649ahmsntdtqm4633cpslt3qllrfp0f8de6",
      "amount": 9064383
  },
  {
      "address": "comdex18n7xdhj2l65w3hqdhne0yj9gjwcpqqcxzterr3",
      "amount": 1108643
  },
  {
      "address": "comdex184h2rl5qet4nqywhv62npqnsyur80p7t6yqunr",
      "amount": 766986
  },
  {
      "address": "comdex18ct84ud7q8xqhxvt3q2ddceqg9kq2aknj76yyf",
      "amount": 1394520
  },
  {
      "address": "comdex18mwtg2c2y8scjcthzg0m0gqv73yydh2kja7daf",
      "amount": 61010273972
  },
  {
      "address": "comdex18uzl0hleecwx8p0ys49hzpjt2fu4sz7qdjty2l",
      "amount": 69726
  },
  {
      "address": "comdex18ut040zez44fxejxmfs02ed6vtdmwhpe7f9knn",
      "amount": 69726
  },
  {
      "address": "comdex1872h267jq0hnqzzxkjgm62y44n3vpztnzllpjm",
      "amount": 45321
  },
  {
      "address": "comdex187wy2grtqlxpp6j3gspyaej3c58efargu6qnhf",
      "amount": 299713
  },
  {
      "address": "comdex187llah7hvkdq4atwgkqqjwe52dcs2amweq8uut",
      "amount": 223123
  },
  {
      "address": "comdex1gqcsp5u2j3uew074en3jv8ekcn9y0xssph8cp3",
      "amount": 1331767
  },
  {
      "address": "comdex1g9y9yp732z6g73hgh2zs8uauzf0f0qkd7n6w7x",
      "amount": 1464246
  },
  {
      "address": "comdex1ggqnctdleprqskjsvuyhdszqaqqd9w00fnuvt9",
      "amount": 690287
  },
  {
      "address": "comdex1gf8v08cvd0k486au38rh7f706f2txnnvv5lwz7",
      "amount": 10180000
  },
  {
      "address": "comdex1g2z8ul7xq2754zwesess8x6d5xaay0upv37yd6",
      "amount": 234062
  },
  {
      "address": "comdex1g2mjkrugdjr0h75cwg9ml5r8ynltd4z9jgtjf6",
      "amount": 425328
  },
  {
      "address": "comdex1gvjgvqrvm2uqqlgp5dlzwhx7egpllln8r563ru",
      "amount": 174315
  },
  {
      "address": "comdex1gv58ljzjjkk0p6cv3zlnp4r7dyedrj3za3ty63",
      "amount": 697
  },
  {
      "address": "comdex1gs7has74zl62n3cmqmw2uulgj4u7tv09zucraq",
      "amount": 263558
  },
  {
      "address": "comdex1g38pv08t94ng2a29qqmnlcn050nwhgmc3mtm35",
      "amount": 348630
  },
  {
      "address": "comdex1gjh6atlq7xypygs7ga32r75kjpey09m9ucrcaj",
      "amount": 17431
  },
  {
      "address": "comdex1gkm5auug7fqfy2u9p8ctqa0j4xmxhphkf9c0gp",
      "amount": 139452
  },
  {
      "address": "comdex1gcxkh4z0xhg447ng39u7lfsfhg6cl946mfgrvl",
      "amount": 348630
  },
  {
      "address": "comdex1g602ezrhwl7763w7h5acu4c4n3e06ngahxfsu9",
      "amount": 139452
  },
  {
      "address": "comdex1gu9g0hxrfvcajur6ahja6wha7f84v8ljcydjfv",
      "amount": 2789
  },
  {
      "address": "comdex1gapfted65eq77ltp7ftr9jk6vq5rwzed8cg46a",
      "amount": 34863
  },
  {
      "address": "comdex1gal22l4farrzhh4r2ux33dc6ks2c8gyjlvm6vy",
      "amount": 836712
  },
  {
      "address": "comdex1fr5qv96a0gg8gqc0fdzv9yrwm2jlu9n3df6nyh",
      "amount": 41835
  },
  {
      "address": "comdex1fx6nxq7n0cxzwps5pm87rafplcgj9tk3p9m6ma",
      "amount": 4741369
  },
  {
      "address": "comdex1f83n7erp22g9dm30wftex290vuw3jjrj8j80gg",
      "amount": 6972
  },
  {
      "address": "comdex1f8ln64cqgg8mdyqdzys69wxwkxx88dguyzeuqk",
      "amount": 1045890
  },
  {
      "address": "comdex1f2ptqja5alaxe2d57tn7r82r39geuag804ymkm",
      "amount": 4183561
  },
  {
      "address": "comdex1f0u98wpsq0ss3fnq4ehwnmt4lfxtp5pmtup8k0",
      "amount": 2517109
  },
  {
      "address": "comdex1fskutxldxqwjds3wtgsqh2raha4hyrm8pxmzv9",
      "amount": 530941
  },
  {
      "address": "comdex1fjhtga2hg3tke378fn772r8zyt08wvsjzmckn6",
      "amount": 418356
  },
  {
      "address": "comdex1f5x9sy2hl2044vpkz27xv6pvmpx3jm5kwlhvch",
      "amount": 20917
  },
  {
      "address": "comdex1f5da8843g7r43gknzqdm9zrf3tdmk5dzmygdzd",
      "amount": 1394520
  },
  {
      "address": "comdex1fkh53jzjltugmxc5y50pm727a85nscrnp7z9s9",
      "amount": 12510243
  },
  {
      "address": "comdex1fhpeeq8jfaczhzdjxfxd22ztjft7mdh27x6544",
      "amount": 1394520
  },
  {
      "address": "comdex1f6f27nvugfqz0qj6dzgzl4p57ee4gufa7pk5wn",
      "amount": 14642465
  },
  {
      "address": "comdex1f6feajvfzskwkm3jgz62la3f8l3nr8u2e6qru6",
      "amount": 119928
  },
  {
      "address": "comdex1fm839jv3l6l2kt28kgzcsp9nr9f4vmudeuvw53",
      "amount": 230095
  },
  {
      "address": "comdex1fmk7uuz73cgn8snkugyl6ggxcrfgmnc9z4m2us",
      "amount": 15518081
  },
  {
      "address": "comdex1flp4675dnwjgrd7evyrzgyk6yl8sn0pye68w37",
      "amount": 707801
  },
  {
      "address": "comdex1fl3n0u9ksn6en29e06gcln5kgd5nhp2v4mfvk2",
      "amount": 1111536
  },
  {
      "address": "comdex12qhqqdxv6v46w7a7mtwkrrl8uzggw2e64vt5h4",
      "amount": 209178
  },
  {
      "address": "comdex12rztp7yuy7fgdjku7vxqx8dxdhvlwn97ztasnw",
      "amount": 238463
  },
  {
      "address": "comdex12rlzcc86p83f30phwruj6swchmg3zqs5nm52y0",
      "amount": 742946
  },
  {
      "address": "comdex12y9ynga89chw4meuufppsgxxa6f0qqqc47472d",
      "amount": 2440410
  },
  {
      "address": "comdex12tyknrawl9cfstdxpyjur0j0ys0a8u403rrrg6",
      "amount": 21545342
  },
  {
      "address": "comdex12wxxg4e073hlckps6ht3af2jklm0nneza3t4nq",
      "amount": 488082
  },
  {
      "address": "comdex1204g59lju5lysugsc9d5d686jnnrhavdlyql49",
      "amount": 1031945
  },
  {
      "address": "comdex1206m2g6dnmyvem538u82futqa3tkej9c3tv6kv",
      "amount": 1394520
  },
  {
      "address": "comdex123q27a4y4uw9qqnv0j8fml4hht63kpt9m074g6",
      "amount": 8270204
  },
  {
      "address": "comdex123tnlz8jggrsecpc8qc039ctex78t7uxcegn4y",
      "amount": 69726
  },
  {
      "address": "comdex12jds5as3c3tk4pttkjywwjjkcxs64l4d0gxj2z",
      "amount": 697260
  },
  {
      "address": "comdex124tztqwj5lkf434pq9jjcxykda2mn964gr4972",
      "amount": 2091780
  },
  {
      "address": "comdex124t7yuz42u6k4yzve4zhu2ak20upff2u3verd7",
      "amount": 432301
  },
  {
      "address": "comdex12h6ytg49zf2w02gyj96svzxem93phfavc7ah42",
      "amount": 348630
  },
  {
      "address": "comdex12c6wdy5gp2swlt3aucw00jkv672n6k3zya3jkr",
      "amount": 97616
  },
  {
      "address": "comdex1txexnzgpavgk7ncxh2p7weqthsza6n7c3zp6gv",
      "amount": 278904
  },
  {
      "address": "comdex1tg20eujl62v226h0404thnu4pu42rd8v4u2g2r",
      "amount": 1185342
  },
  {
      "address": "comdex1tg6kav2cu8envfv7qxhalma9wawcfvxgmd4s0n",
      "amount": 3915116
  },
  {
      "address": "comdex1tf3jl53fuctqx8s3p57vkhrc9l7tptfzxs0nvu",
      "amount": 3978496
  },
  {
      "address": "comdex1tsdtq6tm46djrwazk7hg4d6rr0zju4n5ze3t3e",
      "amount": 812308
  },
  {
      "address": "comdex1tkelcm4rpx7cec5hpleyz855hs8fvuvtq9d7m5",
      "amount": 8227671
  },
  {
      "address": "comdex1tkahrlxmsn4qsfp2sawaezzyqyxucymmphpn0s",
      "amount": 334684
  },
  {
      "address": "comdex1t6uptc2grgupudp477ry0dxd6fzu5l2d207h2e",
      "amount": 278904
  },
  {
      "address": "comdex1tuzdqwynm2xtpm2jjpr80cjrf46hcp2nkd7jcf",
      "amount": 10428224
  },
  {
      "address": "comdex1t75lpa0tyt29ka87azcvtwhk0hwup83mvclvja",
      "amount": 6972
  },
  {
      "address": "comdex1tlkhrt874den98e0jed4qsamhf2r85ps7yh27d",
      "amount": 488082
  },
  {
      "address": "comdex1v98gdcgq0dzucgpkjwr76s8gkk8n4lsl70txnd",
      "amount": 209178
  },
  {
      "address": "comdex1v2ytakyw3nlj6yuuc4z2rwcwwszwefnafend2f",
      "amount": 376520
  },
  {
      "address": "comdex1vv8jcjhq7ntazp6kq7zw665hl876ku8a2cvyrd",
      "amount": 634506
  },
  {
      "address": "comdex1vvdqzqq23kvldcvwurkuthnwmgp6wredpppjzc",
      "amount": 3486301
  },
  {
      "address": "comdex1vwpqc8mfxlux4rmrdz4en78f7cxzd35xyj5qdy",
      "amount": 568964
  },
  {
      "address": "comdex1v0d6rz6ew9dgp7e5wdrnrvsyeysqd6f0x28u67",
      "amount": 4880821
  },
  {
      "address": "comdex1v436zed2mmhx7g9rcpvyh7l50l76jjgkgjmtht",
      "amount": 46716
  },
  {
      "address": "comdex1vkxq45an94x7av85mp9kgnhgep5ye3anw4lru2",
      "amount": 10180000
  },
  {
      "address": "comdex1vcr7fvfkpplxhz87fn97yj65vamyv0y5zundlp",
      "amount": 1115616
  },
  {
      "address": "comdex1vcl9hak978xglczrhdgaprgjytkseescpd3uap",
      "amount": 815794
  },
  {
      "address": "comdex1v7skkttgpxawuavlwelem36el30lppvnuk7xhp",
      "amount": 20917
  },
  {
      "address": "comdex1vl3ferckwrzhdfl7j4jemv4ct84vf2q8lg6qkn",
      "amount": 181287
  },
  {
      "address": "comdex1dqsq3duffjk5czwvckem39fq6arga5dz53tquj",
      "amount": 139452
  },
  {
      "address": "comdex1dpfzr9c05gyltmlxjthhcz93a7fafznuud2d5d",
      "amount": 69726
  },
  {
      "address": "comdex1dpuvtmk6u47f3udwuzdftn5kc8n9gytvjhjvqs",
      "amount": 41071
  },
  {
      "address": "comdex1drxq9y3kdps3txt0sdvvfugeunzuclfrpzz3jy",
      "amount": 174315
  },
  {
      "address": "comdex1d9ewgkn4826duzprzpthn8qrupd0m7t78emqf7",
      "amount": 228701
  },
  {
      "address": "comdex1d8v9te0792xwl958zhel30mkydcpmxfv9qu7r6",
      "amount": 209178
  },
  {
      "address": "comdex1df6nprg53qsydl72gjzks2lj9frf6p08j9tfh3",
      "amount": 361808
  },
  {
      "address": "comdex1dvdp8gcajrujvdhvqzhl38ugudmff0r7rlcllw",
      "amount": 2923612
  },
  {
      "address": "comdex1dwzxaessw2fpvfzyh4wp7kgmgcfgeqth4jrefj",
      "amount": 1185342
  },
  {
      "address": "comdex1dshjyglml4twpn8a3f0l8e7w78fv88qc9czmsn",
      "amount": 3486301
  },
  {
      "address": "comdex1dcr2sh4vtnzrqef2t6x39yu3u0uwxj9z0v03n3",
      "amount": 80184
  },
  {
      "address": "comdex1denhtqkkulvk00pja9y2xwa3hwymk4fygtadww",
      "amount": 288665
  },
  {
      "address": "comdex1d7zp2lynypuucv8ecm96q6k8cewcxf3fry590d",
      "amount": 240554
  },
  {
      "address": "comdex1wrzyqgh4u24cwdyx9cfvltr85mu0u80qwzl9m6",
      "amount": 1831005
  },
  {
      "address": "comdex1wyfmrr2awjanngyj27zt7x5mtfrpscz6j3grnv",
      "amount": 22312
  },
  {
      "address": "comdex1w9ajdrazwp85nvgnaz6zn6etjm8wgs3mfjt27s",
      "amount": 1494228
  },
  {
      "address": "comdex1wxjkwhyw7uwzuevw70cn6mlea8mx4ezulfgnpx",
      "amount": 34863013698
  },
  {
      "address": "comdex1wgey85a4ad0y35z4h64s92asvvsngpc6altpsu",
      "amount": 13945902
  },
  {
      "address": "comdex1wfdsemrxg8c3yuqtx8af3lrenje87705lsn0dt",
      "amount": 697260
  },
  {
      "address": "comdex1wfwqp447e2kl5hz5nyvzgurrm02tyrfv7vspfx",
      "amount": 20917
  },
  {
      "address": "comdex1wvjv9ujz5vy25ps8p2cal4f2vf2cmkla30k427",
      "amount": 348630
  },
  {
      "address": "comdex1w0cac5s88hspvufphq5sww9r4ya0kpxma0c3p0",
      "amount": 697260
  },
  {
      "address": "comdex1w4tv7u3kzusl0ze5e5zacn6xnadrgf2mxpdm2n",
      "amount": 29494109
  },
  {
      "address": "comdex1wesdsz9zyq9snz6w7u88845rwhccktv3gvwhpy",
      "amount": 2684452
  },
  {
      "address": "comdex1w65ynudwkynt2e2gnz2rkzewszh4clz4zsa2h5",
      "amount": 1882602
  },
  {
      "address": "comdex10q0vwkwcftr24yde5lpukfry87qkn8gf9ns9ze",
      "amount": 278904
  },
  {
      "address": "comdex100jq4vg2lwc6pm73f9njl7e8wwrjzkf3cvfa5t",
      "amount": 1366630
  },
  {
      "address": "comdex10stqrste3s8ccz2zrv3wktuzcffyjcjqc4jkv9",
      "amount": 53652
  },
  {
      "address": "comdex105evv74egn4epg3ea3gdx92kjmze0arj93v783",
      "amount": 1045890
  },
  {
      "address": "comdex10a8xwps7g9q7v5th6vqsvh6mkgagythrj4tk3d",
      "amount": 4586578
  },
  {
      "address": "comdex107r0ylq95jdantsxr2v376s3mu0ufh6gmaxaam",
      "amount": 697260
  },
  {
      "address": "comdex1spe42m04gq93j7hj793lq7avlmehydm05hs3n4",
      "amount": 9064
  },
  {
      "address": "comdex1sz8t9sm3pwgeqdfkzm4jgwlkuqjaud0f9s2cek",
      "amount": 4253287
  },
  {
      "address": "comdex1sztkmw500fllurt40svw34lkj74yf7ewkn7m0r",
      "amount": 557808
  },
  {
      "address": "comdex1srujzhj2v9fkzhnn635udlczyhdpetuhmp5jaq",
      "amount": 6693698
  },
  {
      "address": "comdex1syk5nlrn2n32tl43v965gvva28y9hs7zz2cgv8",
      "amount": 557808
  },
  {
      "address": "comdex1sxxazdfzcktly8ypeq6mlj0lwyt456pltssmpm",
      "amount": 46019
  },
  {
      "address": "comdex1sgtkl53fvgjnscahnpknaf6dxxe6cutv8d2847",
      "amount": 19090986
  },
  {
      "address": "comdex1sghqztvgjv3ug6ltpvvwwyhpt35u7hu8hc30xw",
      "amount": 4113835
  },
  {
      "address": "comdex1ssmeg0jhyr57kwjlt508lgz8jeqsfx850al2ty",
      "amount": 697260
  },
  {
      "address": "comdex1s4utnveupv6gkcgtfh03ylx8gcu6ht8pp0gvst",
      "amount": 836712
  },
  {
      "address": "comdex1sc9h4dlftwqv82utj52xfgrv0xevuwng6zfuj4",
      "amount": 1900034
  },
  {
      "address": "comdex1scn3q5728nq9f48uleclsvjqjxhm2x2rnp38s7",
      "amount": 11567547
  },
  {
      "address": "comdex1s6xayykxqmf0y939reekq5gnxd336hkm78y5e2",
      "amount": 58500136
  },
  {
      "address": "comdex1salv3vd0hzm9nh9rhdf2xs6fqvzfrwqzzr06lk",
      "amount": 2091780
  },
  {
      "address": "comdex1s7565kamh03cqdvktwphjestnwsdgl29y4q8j5",
      "amount": 2022054
  },
  {
      "address": "comdex13qg8wcye68jy2pzxlft33kp6q5u846c2cm5s0d",
      "amount": 697260
  },
  {
      "address": "comdex13yu3kq4a0edshnk5sf8k7403jg3q7h6qgu0a0a",
      "amount": 348630
  },
  {
      "address": "comdex130698mc92gykruhp045sqlkdwlxqffwmc28rf6",
      "amount": 23512210
  },
  {
      "address": "comdex13s3gy8xpcqennypj6jyttpqz2hcmf3z2ah6gvg",
      "amount": 181287
  },
  {
      "address": "comdex13306hrp8e9y9z8lv5h2eccupffj8kw95c6hc20",
      "amount": 69726
  },
  {
      "address": "comdex134sne96llys4rzxt0jayv0a4wleqk9mycn043u",
      "amount": 69726
  },
  {
      "address": "comdex13hvk652t024wedxyf5e7djny5sxcedvap85h0c",
      "amount": 90643
  },
  {
      "address": "comdex13cpyxcn4qgyq4zgctztla6ssy50xthm03fe57h",
      "amount": 405306
  },
  {
      "address": "comdex13cwnj8y0pxpth8570evysgsh9erq7veph06qe8",
      "amount": 6700671
  },
  {
      "address": "comdex1367lycvgnpsqpc8n9xn6kctvq2fnfy7ezeg4r2",
      "amount": 96221
  },
  {
      "address": "comdex13al2crynaqq9wdeynx5nul9smtf4g428gka4gg",
      "amount": 488082
  },
  {
      "address": "comdex13753xk4kt8503ex794udl2l0fmzlaxn68sg4gg",
      "amount": 6972602
  },
  {
      "address": "comdex1jq0ys3n8z88atkww4qm2p3j73dhgwzf2kjvh5s",
      "amount": 43927
  },
  {
      "address": "comdex1jq5satyq8vhjlr975z3rmnmhgcvxec6nuarm79",
      "amount": 3172534
  },
  {
      "address": "comdex1jpvxqtd5uq5q4nljgu0zr0d25u6605gq4tlxrd",
      "amount": 6972602
  },
  {
      "address": "comdex1jyvxm6h9r8wjml83wfrhh2ycvek58944urtlpl",
      "amount": 348630
  },
  {
      "address": "comdex1j8wunrzrt79gvn0rvtmanghe92ks84qss9eexr",
      "amount": 5578
  },
  {
      "address": "comdex1jgfmpq3c9p2gxt02m7cy8tpen820n9c27d6vsg",
      "amount": 1394520
  },
  {
      "address": "comdex1jgfawj9mjfaglfkqqy5y7sutyzd0vsxwpjrjnt",
      "amount": 3486301
  },
  {
      "address": "comdex1jv0qqczqmx42rzffxqf488rwn8vlnte4xpf4gf",
      "amount": 697260
  },
  {
      "address": "comdex1jdg0m53ssgf2u8fnn5dkngl78vym6y3lgs9tae",
      "amount": 209178
  },
  {
      "address": "comdex1jwexxu70xc7kvwsgrqp2mp56f4lkqpdqqyxw29",
      "amount": 209178
  },
  {
      "address": "comdex1j053fyvnkkxq25dzczt7dsgrgqnq6u7we0xq30",
      "amount": 976164
  },
  {
      "address": "comdex1jnpevf3zg3wc3vyqssxu96lz899g00skew73yd",
      "amount": 120416849
  },
  {
      "address": "comdex1j5l0tg7nxjj84hnkyyywph2avulyzyshgw766n",
      "amount": 766986
  },
  {
      "address": "comdex1jhlme46p9ej5t5nl56zs5n7dmve77n6gzu6scx",
      "amount": 19809164
  },
  {
      "address": "comdex1j65d93fk4ae4nj3gd42n7m64vvayhlm9054ama",
      "amount": 31376712
  },
  {
      "address": "comdex1jmvshjwlu956nh6lhqd22dfyr2rzpvep9zw9ka",
      "amount": 2509718
  },
  {
      "address": "comdex1jmjjk8xy3r6dydlewss4gyhz7cdnxse2zjdslt",
      "amount": 766986
  },
  {
      "address": "comdex1j7538dmeqnuvfz542s85anfrwd0dpytn7592qt",
      "amount": 55799
  },
  {
      "address": "comdex1jl82qy8645dpmrcxe64wfxtq8rygyxlrffy6sh",
      "amount": 20917808
  },
  {
      "address": "comdex1nqruweed0w43p00z8yz27taygk7pqu89w844l9",
      "amount": 1743150
  },
  {
      "address": "comdex1nz0szznjs6xp6y96qprcye5me5u09qpcp8r9kh",
      "amount": 697260
  },
  {
      "address": "comdex1ndslxsucavg3eglqe4mzge74tdx67rcn7fh0sp",
      "amount": 6972602
  },
  {
      "address": "comdex1n5gytpynpg89z39jy7w2l7gur4w58calg2e7a8",
      "amount": 1115616
  },
  {
      "address": "comdex1ncpnpe5n94582kvcn4574dhxl4fk30srprywq3",
      "amount": 705297
  },
  {
      "address": "comdex1nethnlm2j7k0vfreew973fu62s0xn8purd2mxy",
      "amount": 418356
  },
  {
      "address": "comdex1nmphggkmjqsdd8h6d7tc8tl6mun9ncrwekyg7h",
      "amount": 181287
  },
  {
      "address": "comdex1nucauk03hcdsgjskasal49g2x6xs0r0jw660zv",
      "amount": 390465
  },
  {
      "address": "comdex1n7nk0n7uww2y68lpfwzdt6n3u435ry00sqhc0a",
      "amount": 6066164
  },
  {
      "address": "comdex15qa7hj8lds0llsq6gg0yjzqpgg97l2fekcn2ee",
      "amount": 607598
  },
  {
      "address": "comdex15x28vg4pyqhch87yeu9cxr5y3hm5tj5dt86anq",
      "amount": 348630
  },
  {
      "address": "comdex15xjkhlhlnrppfj479zn3enxaj0ay6hz2eatfvr",
      "amount": 69726
  },
  {
      "address": "comdex158ep04lmck7hx7pt5vwkjhshkqs6wrhjljecg9",
      "amount": 2175452
  },
  {
      "address": "comdex15dznakwa8e78waj0mzn2xx9ee07xegc5x9elv6",
      "amount": 26356438
  },
  {
      "address": "comdex15spr5r7f7egwmzfzehahu7cxdlde9yw8vqnkqn",
      "amount": 13359663
  },
  {
      "address": "comdex15sdehu0qksljngzjvtjav2uh8s45ugm79xjy4x",
      "amount": 5724506
  },
  {
      "address": "comdex15nf66lutnk3zcldx9k7u05ew8pp5cd4qkjwxuq",
      "amount": 75033
  },
  {
      "address": "comdex154v4lhj72p630aafk9fusdffm0faqgkn52dyxk",
      "amount": 69726
  },
  {
      "address": "comdex15hthm2jjc7u6zt0fq0u55yezyxrqu67tgrv7lt",
      "amount": 146424
  },
  {
      "address": "comdex15hu8vsgdgjfrnecxrsfgv7qk0pea96qnsamr0s",
      "amount": 697260
  },
  {
      "address": "comdex15hl304vcz8jlke3ntrqm0r3ug7lfpywpgyslte",
      "amount": 826214
  },
  {
      "address": "comdex15csqzwakl57nn3gc877c93aqlfs08dxze3kt0a",
      "amount": 13945
  },
  {
      "address": "comdex156tpw2ks6cxd6xgvfljzdwak0ntskpru3la25q",
      "amount": 488082
  },
  {
      "address": "comdex1565dt8wfm92gacdmhj5rj3ecceteq0pvhgee8k",
      "amount": 2998219
  },
  {
      "address": "comdex15a84pmu5qa9pps40g7vyzqar0p77mg0eeedkjz",
      "amount": 139452
  },
  {
      "address": "comdex14zuqtexs7gyurtspw9ztxdlcmql890hjnxqwdx",
      "amount": 2231232
  },
  {
      "address": "comdex14yaad6w7unfwy7qvav4vkskkqeflda5d9v9xra",
      "amount": 1394520
  },
  {
      "address": "comdex14xqu93tlh4cp2l4jdsjjn0t2xuue09p2wncdje",
      "amount": 14642
  },
  {
      "address": "comdex1489a36405052np9vgf4ajzam4taphu39vy7erh",
      "amount": 287137
  },
  {
      "address": "comdex148jaxj5kc7gvzd6yahed0txezy5ncg3pxxdgdq",
      "amount": 1394520
  },
  {
      "address": "comdex1436sgl4vc9rgxc326xnse38myxnr3sy8kag5hs",
      "amount": 360078
  },
  {
      "address": "comdex1437537gm88mh3xsuhzl33cvc05995d4vf02xpk",
      "amount": 174315
  },
  {
      "address": "comdex144ena0yqqmfya0ksz8qjg2d0rdm86nlqyjaju5",
      "amount": 139452
  },
  {
      "address": "comdex14kr5dzuvnytyatytwx9gs4n4d8pp5jzm5z9sfr",
      "amount": 557808
  },
  {
      "address": "comdex14khvd82d582zhns8nun04p8qje03zj36xt40hw",
      "amount": 1743150
  },
  {
      "address": "comdex1kqfva8npaz9s8p9n5khf3k4vrza7p6xr2awgg3",
      "amount": 36568043
  },
  {
      "address": "comdex1krwl595epnmuy3ca0msdel498clplvljunjjhn",
      "amount": 42532
  },
  {
      "address": "comdex1kr7jcknas4q9t63pk4v8u4w6y4plrng9339vmx",
      "amount": 285876
  },
  {
      "address": "comdex1krl02crdupsy2ljrf8kenswxm5wu0z2xx2khy6",
      "amount": 1778013
  },
  {
      "address": "comdex1kykrr9f2uw3jxcvs59u8g0guglywlre7hqj7zz",
      "amount": 6
  },
  {
      "address": "comdex1kt74j7ny8lu0yzmtsdgntwhd94cq89nd5lmsna",
      "amount": 592671
  },
  {
      "address": "comdex1kscjvmcrcwra47yfugyew0rh4rw8475twty7qp",
      "amount": 4253287
  },
  {
      "address": "comdex1kjxa7524tcjzcmtvld0z8vny53ps084vylpwaz",
      "amount": 3277123
  },
  {
      "address": "comdex1k4zpaqh8cmpw2l7vk5xj0zuxdv52fkfq5669lj",
      "amount": 3486301
  },
  {
      "address": "comdex1khgnev0l94jsnp06kqvn844juxxlw456glyqh8",
      "amount": 766986
  },
  {
      "address": "comdex1k6v7yycttwplr4465ry74qv7ukax3n78zxluqh",
      "amount": 642176
  },
  {
      "address": "comdex1kadqxcyt4nnxpxgdtafd3ckg60m8r79ph844kd",
      "amount": 94130
  },
  {
      "address": "comdex1ka0yduzwes79hn8hctpw7tw8hqkz6pwzdhr7tm",
      "amount": 467164
  },
  {
      "address": "comdex1k7vzggejlw8uwdc99qe8let4sp5sc7updhp3fe",
      "amount": 48808
  },
  {
      "address": "comdex1klp3lqdsdfnus9g4dcl7hffe74wy3utq2fgccl",
      "amount": 5654780
  },
  {
      "address": "comdex1hz8kfk5d5ksd4zg8v50j8nsvg8g5swwxpurfvs",
      "amount": 2579863
  },
  {
      "address": "comdex1hr7tz3qvyhlpghtyfxdlwrfkw5zcjkxnpzdank",
      "amount": 90643
  },
  {
      "address": "comdex1hx82q2jajldrx39txl9va50r4csp92283yqt5q",
      "amount": 1533972
  },
  {
      "address": "comdex1hx0kc7v9r4x637kr32320sqdqgu8u5lausvh5p",
      "amount": 348630
  },
  {
      "address": "comdex1hfql2ue3a5ecxrjypxmzc7u2h8c6ds6ph030df",
      "amount": 1234150
  },
  {
      "address": "comdex1hfmayju8c9vglnwfvx7cknshxev8l289un4nvr",
      "amount": 69726
  },
  {
      "address": "comdex1hvjylnwyhcz0vu0pdyzt462e9x4lr7twtcmx72",
      "amount": 20917
  },
  {
      "address": "comdex1hjt3payed00wp7zh80safvcmz6y64cqh5yqs77",
      "amount": 110167
  },
  {
      "address": "comdex1h4qf7fw3u7h0e5kylcmn5ywxz3cpmhldq6zxcm",
      "amount": 34863013
  },
  {
      "address": "comdex1hey6zup6z04ewsf4apws7fzg7p3309jphprl3f",
      "amount": 279708311
  },
  {
      "address": "comdex1hmzk8ngj5zx4gxt80n8z72r50zxvlpk88up0yc",
      "amount": 3137671
  },
  {
      "address": "comdex1hmx5xqzagk0pyy8wgyqt7qpud5366kn45gws0c",
      "amount": 137604
  },
  {
      "address": "comdex1hls4rdk5c9c8ff2cdgrjnu3x4dg0v75n74lgv9",
      "amount": 1533972
  },
  {
      "address": "comdex1cqchutq347dsxx2dlem48whfzft7ldsddd46la",
      "amount": 4183561
  },
  {
      "address": "comdex1czr90qmtskavyq334t2cwduwlzuw5m4zxddmmr",
      "amount": 7112054
  },
  {
      "address": "comdex1crcxlut83lcey4aewc3uk2nhc724qucsuphf07",
      "amount": 76698
  },
  {
      "address": "comdex1c9yrdc0rs8j7xxedcgyaz3en945hutgfj9c3qp",
      "amount": 400227
  },
  {
      "address": "comdex1c9x6yxwk5tgqwh2st0ar4k6zdjs2wvmaj6jrfd",
      "amount": 488082
  },
  {
      "address": "comdex1c83rhl4c3xgy892jgrucmwv79fdhylpaser6ek",
      "amount": 14650832
  },
  {
      "address": "comdex1cffrhp203p8k6za80hgssxagp0dsugf4g2x2zx",
      "amount": 766986
  },
  {
      "address": "comdex1cw8metw5h4upvxayu8vw7tmghzu4xsth9ay0xj",
      "amount": 278904
  },
  {
      "address": "comdex1cw2mul75wptmfld0svztdse4f4qnfn8vjeqxvz",
      "amount": 836712
  },
  {
      "address": "comdex1c4ux38yh77q8kjmg68f29kxw4rt3wlxadwg6cm",
      "amount": 1045890
  },
  {
      "address": "comdex1ch7xgjptr3e7fx89dwpvcsy973mkcdyms48888",
      "amount": 13945205
  },
  {
      "address": "comdex1chlmc3um9qrs66g3afvcpwtexfnhdf3a8v4384",
      "amount": 20931
  },
  {
      "address": "comdex1ccyn9rd6j2gxd0rpuxpq76n9lzejshwaf8xe0p",
      "amount": 209178
  },
  {
      "address": "comdex1ce80gp9gzgvtzr9cd6k0dagzkl49c9t3xk3z38",
      "amount": 278904
  },
  {
      "address": "comdex1ca4e5r24vhr9q7lft4rdjtkv2gj6jch2pz8tu6",
      "amount": 2719315
  },
  {
      "address": "comdex1ezcy89pdxypt575j242hndqxgzhq2st4wapspe",
      "amount": 1324794
  },
  {
      "address": "comdex1e9xsxlwyc7rqqsmpzs3ykyhcsfazr7m42j2ma7",
      "amount": 3207397
  },
  {
      "address": "comdex1efsvp7apm2utc3zlxh239a4xs5vh5ekms8dwe0",
      "amount": 132175
  },
  {
      "address": "comdex1evh3h67cw0kxm00hqdxssslser0tfhr8j7tuxu",
      "amount": 1348027
  },
  {
      "address": "comdex1edkh8yjpw56hfu8ed02yukllkmapkjsfq9c0ja",
      "amount": 275066
  },
  {
      "address": "comdex1ewjud30a3sjldh3l0h47e2cdn93kdavvxq2swa",
      "amount": 592671
  },
  {
      "address": "comdex1e30klgckzt2tm77pj26hpmyn8ayy2rxsvjtee4",
      "amount": 45321
  },
  {
      "address": "comdex1e3jgcrm5x3gwtzy6vg8fcsh8qkqun7sdr4rtp4",
      "amount": 1533972
  },
  {
      "address": "comdex1ekcp008fydmda9vlv3s2ech7677xwhdcwyxq0a",
      "amount": 1394520
  },
  {
      "address": "comdex1ek6mxuayf030zaew7a3hcvffke7ey5rhau0r5h",
      "amount": 6972602
  },
  {
      "address": "comdex1ehgxht05ymwtw46l0vz8dwpryhf47stlg98z05",
      "amount": 2064120
  },
  {
      "address": "comdex1ecfy0fp389hesxz2q2us5de3ve8ha374hpkgux",
      "amount": 697260
  },
  {
      "address": "comdex1ecfca6270ux57kv84eduy4sl2zp9dpyhllxn2w",
      "amount": 264958
  },
  {
      "address": "comdex1ecvvgp7vspxxtzj05t9y3ys3awlxfwe4wl2uth",
      "amount": 4880821
  },
  {
      "address": "comdex1e6ec2n8zgx604csuzhqdm83utj8x899vkfwenz",
      "amount": 836712
  },
  {
      "address": "comdex1ea4wxmfqtedgvnpx6hju5qcmuqspht9ca057qx",
      "amount": 522
  },
  {
      "address": "comdex16qp4ttkkcgn6fw5dvwt7k8aqh3vq0taw46z6ga",
      "amount": 418356
  },
  {
      "address": "comdex16pyllqyq22wxc5yxms7rzw76fes4f9lx3gtmgh",
      "amount": 1673424
  },
  {
      "address": "comdex16z2xdmay655t00ju8q0x6djy8dnc9ecs624xjx",
      "amount": 488082
  },
  {
      "address": "comdex16ryvjr8zuvxfaycwup3jgvgxthafzq0pv69pp4",
      "amount": 979705
  },
  {
      "address": "comdex16g3l409egd4pfa0fc93ey9vvrexddek7dquwgr",
      "amount": 11156164
  },
  {
      "address": "comdex16f65k3v7cesmqldf865h93qyh6plk2vjy2yvtm",
      "amount": 1282958
  },
  {
      "address": "comdex16tgadds932w5j3utfprdcygy85tmy34dnr5m3d",
      "amount": 662397
  },
  {
      "address": "comdex16wf7njyt5hrda65cjesg3yw5r9s3z64ke64spx",
      "amount": 75722465
  },
  {
      "address": "comdex16kdysctr2q5hlf80sp8vwa66cz8a7z29ah09tr",
      "amount": 801849
  },
  {
      "address": "comdex16c7nw5h7wrfth69p7ulrktv66dxqgs76zvla3g",
      "amount": 1394520547
  },
  {
      "address": "comdex166cc97v0x5dalv3gc8c7q856v9949hkmgwe4fr",
      "amount": 2022054
  },
  {
      "address": "comdex16lnqygvcmmd2z63dzuyx80pgyzdxnrsnn0qqkf",
      "amount": 110864
  },
  {
      "address": "comdex1m8q0z73c4ud5vcnu4tg465hayu4472t4f5q59v",
      "amount": 2649589
  },
  {
      "address": "comdex1mgm84ze8a284lgawrafxdftv6pyrlqd3gm7993",
      "amount": 15339726
  },
  {
      "address": "comdex1mfyfrxnee33kvsz3wgc8fqp5w479wjzyu73n5e",
      "amount": 278904
  },
  {
      "address": "comdex1m0qser0za9raujck6wgvsnmt9zzn03dzqycl2x",
      "amount": 5525787
  },
  {
      "address": "comdex1m3esxfpc2qlk75jyjhfgyhma3std4cmvn6x3pw",
      "amount": 50202
  },
  {
      "address": "comdex1mewapalkpsc0wku98zenx6ecdaw07ez3wyrvwu",
      "amount": 50551369863
  },
  {
      "address": "comdex1memcf6sqd00h6yzf0e8e8e707tjmszpd8632er",
      "amount": 697260
  },
  {
      "address": "comdex1m65ackcyu0j857xetslptvrqu84722d0xch9e0",
      "amount": 509941
  },
  {
      "address": "comdex1mm7z2zanke4vguug84faf9k8f9xfl27wp9tlgr",
      "amount": 251013
  },
  {
      "address": "comdex1up7p4uhgzfzrk6fnqjes0wty9rv7q0yxnp0vus",
      "amount": 348630
  },
  {
      "address": "comdex1uyza5f6l78w2uy68qs25nn9mgyc4zm3a7d9tlc",
      "amount": 1429383
  },
  {
      "address": "comdex1uy5hzhmdxky4l76afes2szhdj3zlt8tu9aahe2",
      "amount": 2565917
  },
  {
      "address": "comdex1ux35kxht4dhxnp43z0q9cdjfcfz05qu0a4h5pm",
      "amount": 2738940
  },
  {
      "address": "comdex1udzenjrhjaehpjj7qzqul3dehqqx9ex5j793jg",
      "amount": 1667986
  },
  {
      "address": "comdex1uswu5spxn994kqxtfn40xxqpktydppvc54n8ql",
      "amount": 278904
  },
  {
      "address": "comdex1us6asg8qkm5d6x6exr08std2cdev9e6kpzj73v",
      "amount": 557808
  },
  {
      "address": "comdex1u3yq3sgmnq36p4c8jj73fmjxgeej69g7u8cqma",
      "amount": 32803230
  },
  {
      "address": "comdex1uj8zszllhzky9mh76tknpfer8lmksrta8p6rsx",
      "amount": 209178
  },
  {
      "address": "comdex1ujv9p0h94dqgyvuss6zcjzq9257ejygxs4j4uu",
      "amount": 20638904
  },
  {
      "address": "comdex1u5zudst3vfmqa0tm4xalnksuywfah0cm3v5wmt",
      "amount": 418356
  },
  {
      "address": "comdex1u5r0j0jytzk2c5z5wswlzy4h385elwwx2x52sa",
      "amount": 70954
  },
  {
      "address": "comdex1u5w88m7glt08pea0c8xxd0gggpnymn9wvp68hr",
      "amount": 46855890
  },
  {
      "address": "comdex1ukx2knxzsh66cnftnrxuaxvpk33kgwaa6afj47",
      "amount": 62753424657
  },
  {
      "address": "comdex1uh3rkv9k929lc3quy7fvy00twxxrla7pfd76y8",
      "amount": 14942964
  },
  {
      "address": "comdex1ucjlwl3vljw7kv6740jz0as6u0jd6dd56ejmjc",
      "amount": 2245178
  },
  {
      "address": "comdex1ue5e5v2wsmkxr87hzlp2h3zkz4fh7cvjt3zwel",
      "amount": 1150479
  },
  {
      "address": "comdex1uaxp2ex423wrz80s39enndqqxyzltvgslq0wmg",
      "amount": 3974383
  },
  {
      "address": "comdex1ultl0u2kp0ftseptesd0sj3m3mvqpwx8z6gnyq",
      "amount": 1665057
  },
  {
      "address": "comdex1apgvcwhn0ukck6vr77rl20zg020x4x0x5ld7t8",
      "amount": 209178
  },
  {
      "address": "comdex1ap6jl9gw5gmqnngk9jxjtcr39vrjdfnkddaj2d",
      "amount": 181287
  },
  {
      "address": "comdex1azvc3wmszj7nu9vaseu0u8s7v7w284ad6p5svp",
      "amount": 18826027
  },
  {
      "address": "comdex1azhlpges0xl4p0mdfmfwnzck3jhfxvh23rj5vf",
      "amount": 983136
  },
  {
      "address": "comdex1ay666skul67prkttj8zcf0ap7rq8jwwjf7jchz",
      "amount": 4312002
  },
  {
      "address": "comdex1a9h9wjwpaxhuu9exfjfrky6rwzhtungdy9w4hh",
      "amount": 14468150
  },
  {
      "address": "comdex1axp55y8y3cptx8wzz0tv6ezqhr86z70s29cm85",
      "amount": 22870136
  },
  {
      "address": "comdex1ax7nplm3jegg4z5ayfjkamrfh4a4e8pwnexfkv",
      "amount": 74606
  },
  {
      "address": "comdex1a20a9ergpqk2dq0jwcdv94c5h4vfpgvsnle5k3",
      "amount": 1045890
  },
  {
      "address": "comdex1atdjt862etcwvr3xnw8sedq07n7nan0prgz55k",
      "amount": 1896547
  },
  {
      "address": "comdex1avwxuq2ye8ut98gz0j3u74fhlx4qmcaqud6ge6",
      "amount": 2409034
  },
  {
      "address": "comdex1awsjg9e7pj232syxf4cf253u5g0l34p87u8v83",
      "amount": 278904
  },
  {
      "address": "comdex17r5zhwcfj57sscjpatgxkz9guac4nuwfp9uakn",
      "amount": 1743150
  },
  {
      "address": "comdex17yzaplnaswn75n0sx3z5lwhf4wzv692fme94r6",
      "amount": 40929178
  },
  {
      "address": "comdex17y9gdgxfv06f67vcterep6s3ug26a5t9c36gwk",
      "amount": 14503013
  },
  {
      "address": "comdex17x50z4t8zsze5amk0m98tg26cscuvzt3uurlxu",
      "amount": 1610671
  },
  {
      "address": "comdex17g9p6qwjac45x7sv26n0zxj3lykwpjltrfvzpg",
      "amount": 1507320
  },
  {
      "address": "comdex17tr4hpzukzk75mawk3tr2f5f88tnsclm3n2av3",
      "amount": 758991
  },
  {
      "address": "comdex17v5m02vs5vap6fwye0l6ecnmj6dg74m7uyysjg",
      "amount": 488082
  },
  {
      "address": "comdex170t47mdjjkgrhzxmyck2x8kgw7rvzu2dv4zclj",
      "amount": 13813994
  },
  {
      "address": "comdex1759qud67da7f5x2nrdkqxgspd9v7eauew3t6fk",
      "amount": 348630
  },
  {
      "address": "comdex17elvh5ptjcfz4yaevfmewhdttwsr9rf3qth0j4",
      "amount": 2091780
  },
  {
      "address": "comdex17m4804kq3yd74p090fcm6dg44vp8xq2d7y8h56",
      "amount": 1251323
  },
  {
      "address": "comdex17ung966cfeymuzmu9vyzc3uqra8d7sxzfw6pya",
      "amount": 57523
  },
  {
      "address": "comdex17ape0gst9pqfmwu52hqdm28ggs067g7yxfswhd",
      "amount": 6972602
  },
  {
      "address": "comdex1lqyfk780f3g25v9ya72l4cpwetg7ky68d4re75",
      "amount": 83671
  },
  {
      "address": "comdex1lz20nzwk9tvpjae6jehws8pn3nhhaww23vc562",
      "amount": 2590422
  },
  {
      "address": "comdex1l9z475x2egvs686j5jtjxm6n8fehlz9ssgrq53",
      "amount": 10458904
  },
  {
      "address": "comdex1lgmalr3y02yl8sjr25uyz430426z4kh6un79fa",
      "amount": 1882602
  },
  {
      "address": "comdex1l238kq8vc2n233lmcxc5urq35yqc5zsa68chtq",
      "amount": 232884
  },
  {
      "address": "comdex1ldmcrk5m88erpaj29yyntqs7q9560nzgn83lt4",
      "amount": 81573
  },
  {
      "address": "comdex1lsdpgfkcvnlx8q9qkr0yc7smtf48untgkxpgm8",
      "amount": 303836
  },
  {
      "address": "comdex1ljwg82hzepj3em96565wd2ej3unkrtp4km343n",
      "amount": 83671
  },
  {
      "address": "comdex1lkd0tc04rchp08vs599klpak7gyry3e5zep9aq",
      "amount": 453219
  },
  {
      "address": "comdex1l66mg6cv79kvgmtcwsv4hzrwy8pj7xf55uq8ul",
      "amount": 2022054
  },
  {
      "address": "comdex1lltxfhkqyhd835ha3j5690wc4s0canjwptv30x",
      "amount": 348630
  }
]`
