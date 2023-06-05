package v11

import (
	"encoding/json"
	"fmt"

	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	rewardskeeper "github.com/comdex-official/comdex/x/rewards/keeper"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Reward struct {
	Address string   `json:"address"`
	Reward  sdk.Coin `json:"reward"`
}

func DistributeRewards(
	ctx sdk.Context,
	accountKeeper authkeeper.AccountKeeper,
	bankKeeper bankkeeper.Keeper,
	rewardsKeeper rewardskeeper.Keeper,
) {
	rewardsStr := `[
		{
			"address": "comdex1qe86f8rt3hlzsspf7efq3h0c3dlsyfqwvv6fep",
			"reward": {
				"denom": "ucmdx",
				"amount": "5124052"
			}
		},
		{
			"address": "comdex1z8zyckttkkgtp9s2vdyc8hd5nu4x6keujx84qs",
			"reward": {
				"denom": "ucmdx",
				"amount": "63696203"
			}
		},
		{
			"address": "comdex1rsgg8v4elqsckvpkrvrdrj590xupvqyv93mx24",
			"reward": {
				"denom": "ucmdx",
				"amount": "31869"
			}
		},
		{
			"address": "comdex1y6jv9m9937ftxzmg4q5zgywngxf70p7a8cv7m9",
			"reward": {
				"denom": "ucmdx",
				"amount": "266106"
			}
		},
		{
			"address": "comdex19xfw39nduwmace77gp3ukuan3z26wyh3s8gesq",
			"reward": {
				"denom": "ucmdx",
				"amount": "4646457"
			}
		},
		{
			"address": "comdex1xfn5uj3fxuemgrjwajgtwynrlds4czcw2ccdu5",
			"reward": {
				"denom": "ucmdx",
				"amount": "149514"
			}
		},
		{
			"address": "comdex1x6vkt9mq4haleuyg8kfdxxa6kdev0hkrum8w62",
			"reward": {
				"denom": "ucmdx",
				"amount": "2983373856"
			}
		},
		{
			"address": "comdex18krzzgkdper7hplh0fredqkdtjzvyjrvt70zl4",
			"reward": {
				"denom": "ucmdx",
				"amount": "39767"
			}
		},
		{
			"address": "comdex1gywh2dmsg8tm6fwtghj22q5f4sdwav36yvcsjr",
			"reward": {
				"denom": "ucmdx",
				"amount": "8153"
			}
		},
		{
			"address": "comdex12ngc8xjqx8adxfl08prgvjgafxfdwsclnudjke",
			"reward": {
				"denom": "ucmdx",
				"amount": "64717014"
			}
		},
		{
			"address": "comdex1tyqutyv520jlj9k96qth8xhfe2xdn6lgaj2796",
			"reward": {
				"denom": "ucmdx",
				"amount": "9836133"
			}
		},
		{
			"address": "comdex1tg8he6zqlust7alpx5tq5fcmwm237xqgqfmnjw",
			"reward": {
				"denom": "ucmdx",
				"amount": "2584529"
			}
		},
		{
			"address": "comdex1tkrcx0u0xkthvf962zjqsce0u36s2swprnzxa7",
			"reward": {
				"denom": "ucmdx",
				"amount": "68056279"
			}
		},
		{
			"address": "comdex1t63mzktv73sz52f8n6pfxrq4dzy379877cfjqx",
			"reward": {
				"denom": "ucmdx",
				"amount": "332600807"
			}
		},
		{
			"address": "comdex1v3fgacyt5f8lmguvrvtyjqqnpjmw2ncj7qv0tt",
			"reward": {
				"denom": "ucmdx",
				"amount": "2784647"
			}
		},
		{
			"address": "comdex1vkjds0kmxqrhcsy2f77r2v8kaywwtdn5u8etkt",
			"reward": {
				"denom": "ucmdx",
				"amount": "4874"
			}
		},
		{
			"address": "comdex1vkurj6pke7fu89chwdv3lqmd7c74ylyrym7rma",
			"reward": {
				"denom": "ucmdx",
				"amount": "42967102"
			}
		},
		{
			"address": "comdex1vlc2wpeqc5fpyfwmxhyrjh45eapfcxp9qkxqe2",
			"reward": {
				"denom": "ucmdx",
				"amount": "100226073"
			}
		},
		{
			"address": "comdex1dp8pnper7ka76d5t987mg7fsqsqslh5qhwsuxj",
			"reward": {
				"denom": "ucmdx",
				"amount": "2826011"
			}
		},
		{
			"address": "comdex10x2jkt4a6ykzurvq9j004jmzhggpjmup505dvu",
			"reward": {
				"denom": "ucmdx",
				"amount": "55305429"
			}
		},
		{
			"address": "comdex1s5hdhhu3xr32nmq7dlxlruu28wztlxcqfumc2h",
			"reward": {
				"denom": "ucmdx",
				"amount": "28008"
			}
		},
		{
			"address": "comdex1nmh039s4fnd29ahfhg077l52cmkw3mjuss0jfv",
			"reward": {
				"denom": "ucmdx",
				"amount": "2291004"
			}
		},
		{
			"address": "comdex15md4wkfqux35ak5rh2409yk3gx8l4y5stjg0qh",
			"reward": {
				"denom": "ucmdx",
				"amount": "24302732"
			}
		},
		{
			"address": "comdex1575cct5fewj2q0kf52mrlvs2xe30knh8ce4q9a",
			"reward": {
				"denom": "ucmdx",
				"amount": "7400677"
			}
		},
		{
			"address": "comdex1ky7m39mszht0rveq9y0ramqnash6jp5kpuul2y",
			"reward": {
				"denom": "ucmdx",
				"amount": "547724403"
			}
		},
		{
			"address": "comdex1cz39juz4ycuhvndp8rk5qhvtuy7ag2q5p7c66e",
			"reward": {
				"denom": "ucmdx",
				"amount": "61084070"
			}
		},
		{
			"address": "comdex1cvfaevfqm02htqqf9qdzqhyq9m4yzdyzrlta6q",
			"reward": {
				"denom": "ucmdx",
				"amount": "9962522"
			}
		},
		{
			"address": "comdex1e7ra39kd43wsr9j7stt3jkyexp53z5tznzdyw8",
			"reward": {
				"denom": "ucmdx",
				"amount": "2941804"
			}
		},
		{
			"address": "comdex16am6h79srmghqx0q5032ugzeuuvz2y08z33znq",
			"reward": {
				"denom": "ucmdx",
				"amount": "56928"
			}
		},
		{
			"address": "comdex1uhgw4s6dfktjcyxyd2gs3yxj364a8t0l7m4vp9",
			"reward": {
				"denom": "ucmdx",
				"amount": "169648"
			}
		},
		{
			"address": "comdex1uc2un2zttpvpdfarcjtyp382sszkv8zjznvxhf",
			"reward": {
				"denom": "ucmdx",
				"amount": "56411244"
			}
		},
		{
			"address": "comdex1azpm7xl3vfj85jnk3dlw334469p7c0aswnp9y4",
			"reward": {
				"denom": "ucmdx",
				"amount": "1208767"
			}
		},
		{
			"address": "comdex1a9h8ffdjrhx5v4q0r5v2zxwdsgfpqv8j353lzq",
			"reward": {
				"denom": "ucmdx",
				"amount": "2760539"
			}
		},
		{
			"address": "comdex1ate9ea9zsne4wtzkh9439ceneg4w9sj9d4r08m",
			"reward": {
				"denom": "ucmdx",
				"amount": "8455762"
			}
		},
		{
			"address": "comdex1as38t6e97eqrvxwgtr76msc7p0tkltelunyf7t",
			"reward": {
				"denom": "ucmdx",
				"amount": "7162193"
			}
		},
		{
			"address": "comdex17xquhkcu2r99zucp5cvs0a3ewdwa92m8jdr3wc",
			"reward": {
				"denom": "ucmdx",
				"amount": "462759"
			}
		},
		{
			"address": "comdex17gg0twymd2zeer3e2fpa09pwguul8hrf4e8x9n",
			"reward": {
				"denom": "ucmdx",
				"amount": "29061"
			}
		},
		{
			"address": "comdex17tdpstcymdc9pg2fravm3f7qkhndd9uydqdcwk",
			"reward": {
				"denom": "ucmdx",
				"amount": "3988"
			}
		},
		{
			"address": "comdex17hgjycfeheh2j9q8t4k5v2xptc66xzs04ymh9z",
			"reward": {
				"denom": "ucmdx",
				"amount": "11465"
			}
		},
		{
			"address": "comdex17lylkfrex48yuewd3dw2j9vned3wgqcgefzru5",
			"reward": {
				"denom": "ucmdx",
				"amount": "7422371"
			}
		},
		{
			"address": "comdex17llw5vaayd4qjplw2t8n4shzex9vw906843tgk",
			"reward": {
				"denom": "ucmdx",
				"amount": "262225"
			}
		},
		{
			"address": "comdex1lqgw9gczxfrnlyqj3wxcjpe72xlwka8l94g5qw",
			"reward": {
				"denom": "ucmdx",
				"amount": "33521175"
			}
		},
		{
			"address": "comdex1znu9japz2ukk9a2zxmcktcqsu80p099mumyx8q",
			"reward": {
				"denom": "ucmdx",
				"amount": "9559742"
			}
		},
		{
			"address": "comdex1yyxqjugck0v0dfkqs8fks7555q7uc0c8my5q5h",
			"reward": {
				"denom": "ucmdx",
				"amount": "14037101"
			}
		},
		{
			"address": "comdex19xeucpasxg08jthuhktw3m82gcnz8myxqvlvcr",
			"reward": {
				"denom": "ucmdx",
				"amount": "193900"
			}
		},
		{
			"address": "comdex19wmd9xzhjnvmr90z8r3cnjlns5kgem9qglt2ll",
			"reward": {
				"denom": "ucmdx",
				"amount": "1584312"
			}
		},
		{
			"address": "comdex1xlwupe9mesp3t9q38dx6xmmwtsccxr9h09fm9q",
			"reward": {
				"denom": "ucmdx",
				"amount": "5263420"
			}
		},
		{
			"address": "comdex18yu5egzzyqhw429fu892d7qvy2syc8n7yyhzqh",
			"reward": {
				"denom": "ucmdx",
				"amount": "1502163"
			}
		},
		{
			"address": "comdex129rw3gnlk6x9annfw93jk86p55tn7wlek8pdcz",
			"reward": {
				"denom": "ucmdx",
				"amount": "59673068"
			}
		},
		{
			"address": "comdex1sr0lw3rnl4u3zz7hxrek0pp26hjhapecy60xa9",
			"reward": {
				"denom": "ucmdx",
				"amount": "468133"
			}
		},
		{
			"address": "comdex1sew77pwgd5ygh3qn7xk44f4qg4nm9hnlv66z4e",
			"reward": {
				"denom": "ucmdx",
				"amount": "814257"
			}
		},
		{
			"address": "comdex1302zcexj6hruy3fpqd5j7m6sh5u0c362dmrqa9",
			"reward": {
				"denom": "ucmdx",
				"amount": "330"
			}
		},
		{
			"address": "comdex1jz7av7cq45gh5hhrugtak7lkps2ga5v0u64nz6",
			"reward": {
				"denom": "ucmdx",
				"amount": "4374807"
			}
		},
		{
			"address": "comdex1564zf4undyu4t5rgmmjwgek4spp3ctyulkhshn",
			"reward": {
				"denom": "ucmdx",
				"amount": "353468"
			}
		},
		{
			"address": "comdex14e04xex0yg635r46u36l2l320t980xdh5mkvf7",
			"reward": {
				"denom": "ucmdx",
				"amount": "157772"
			}
		},
		{
			"address": "comdex1c32l497qja2k47q0mvxxgtx767lh3setj8n7lj",
			"reward": {
				"denom": "ucmdx",
				"amount": "268752"
			}
		},
		{
			"address": "comdex1mgll8lfmpdyadctt2ga5vza0dtw3gtpyt5kln4",
			"reward": {
				"denom": "ucmdx",
				"amount": "4754560"
			}
		},
		{
			"address": "comdex1alp7d7p80pny3dvzpp8a0fr49mmututy7qztlg",
			"reward": {
				"denom": "ucmdx",
				"amount": "36956747"
			}
		},
		{
			"address": "comdex1zsjynu7qum0zjsmf49t5y4quvsa0nc09h8kyjy",
			"reward": {
				"denom": "ucmdx",
				"amount": "380786"
			}
		},
		{
			"address": "comdex1865y49tks72cus005v2zdg8prqf0wyk75zfjj0",
			"reward": {
				"denom": "ucmdx",
				"amount": "2969782"
			}
		},
		{
			"address": "comdex1nmm797vnca753pny7ptudy5vx5uz0tgz0t8swx",
			"reward": {
				"denom": "ucmdx",
				"amount": "152563087"
			}
		},
		{
			"address": "comdex1e8ls9qgznqe4qk27c7lcmrmqtuzz89qcl6kpte",
			"reward": {
				"denom": "ucmdx",
				"amount": "2140167"
			}
		},
		{
			"address": "comdex1enkvrxhlhc8w3ta30j8ur4wgk5zv0fapjrp7ht",
			"reward": {
				"denom": "ucmdx",
				"amount": "3796975"
			}
		},
		{
			"address": "comdex1qqrealw66hljd36lz20pz7z4z40x33mqxhvfq7",
			"reward": {
				"denom": "ucmdx",
				"amount": "1756"
			}
		},
		{
			"address": "comdex1qpg2942rlm5t98d95eeyhlnnezrtcpy6avn56x",
			"reward": {
				"denom": "ucmdx",
				"amount": "106225324"
			}
		},
		{
			"address": "comdex1qrxxezf7gd6ash9ynpjx7ad0tkhjnxdaf5q9ag",
			"reward": {
				"denom": "ucmdx",
				"amount": "1044"
			}
		},
		{
			"address": "comdex1qywsy6nhp5axsxr4qzd28j279q9jzp8h5gx3kl",
			"reward": {
				"denom": "ucmdx",
				"amount": "94"
			}
		},
		{
			"address": "comdex1qsndm44zccyjz4qj9h2e3pmw5r70x07tamcy9u",
			"reward": {
				"denom": "ucmdx",
				"amount": "8767"
			}
		},
		{
			"address": "comdex1qn08klln63pf390z3mvlt4us45qedknryj7sg0",
			"reward": {
				"denom": "ucmdx",
				"amount": "584850"
			}
		},
		{
			"address": "comdex1qclqe7mgty4vvzhk0zzr99f5g5hqpxm8vrc3ry",
			"reward": {
				"denom": "ucmdx",
				"amount": "114300"
			}
		},
		{
			"address": "comdex1pq7mnp9wk5mzeus9gz047ajalzcl2q7prd4xas",
			"reward": {
				"denom": "ucmdx",
				"amount": "957"
			}
		},
		{
			"address": "comdex1prp23c44ycmhg6e9htdkwtuffpn6zz4uw0adtz",
			"reward": {
				"denom": "ucmdx",
				"amount": "9410693"
			}
		},
		{
			"address": "comdex1pyvh2um57twtenjqvehnchtce66wl02mgcsdh6",
			"reward": {
				"denom": "ucmdx",
				"amount": "23699"
			}
		},
		{
			"address": "comdex1p9he4f4n96rgtf23hdvysjvhzqs5uk8gtjkurn",
			"reward": {
				"denom": "ucmdx",
				"amount": "12433"
			}
		},
		{
			"address": "comdex1pvg7zv5evqqxlkl6wtu92gemf49vjvexatcysu",
			"reward": {
				"denom": "ucmdx",
				"amount": "33637"
			}
		},
		{
			"address": "comdex1p0uesu8nf4fgh8xxmper5gqthd4wz6xtesec6a",
			"reward": {
				"denom": "ucmdx",
				"amount": "98600"
			}
		},
		{
			"address": "comdex1p3vq9vkq03q2y2rmjmx3c6gg2zh93e6glcy3tc",
			"reward": {
				"denom": "ucmdx",
				"amount": "15019701"
			}
		},
		{
			"address": "comdex1pjd9frrdmupezfyz4c80p2yg84hmasaf9xql8l",
			"reward": {
				"denom": "ucmdx",
				"amount": "28105"
			}
		},
		{
			"address": "comdex1p4xquku4hcvas662wdqav9us0qh0r5arw70ugp",
			"reward": {
				"denom": "ucmdx",
				"amount": "1302"
			}
		},
		{
			"address": "comdex1p4l8y5ugshz3lg9ln9q9krzjpmmawz2peg8c8x",
			"reward": {
				"denom": "ucmdx",
				"amount": "86622"
			}
		},
		{
			"address": "comdex1paq5h6cp2fcwv2969xs6vr6ljtyaj3fvgzcyh4",
			"reward": {
				"denom": "ucmdx",
				"amount": "190930"
			}
		},
		{
			"address": "comdex1zxqgxxe84g5racwfdcwjyn4hfa98e2g4j0h9xq",
			"reward": {
				"denom": "ucmdx",
				"amount": "237528"
			}
		},
		{
			"address": "comdex1zfhah8kx68xl252aqm29j6wqfuh2vtm6fy86m0",
			"reward": {
				"denom": "ucmdx",
				"amount": "21789"
			}
		},
		{
			"address": "comdex1zn3a84lujzyg8ecdjh5x4qs7p2mk78jpfscuax",
			"reward": {
				"denom": "ucmdx",
				"amount": "1039025"
			}
		},
		{
			"address": "comdex1zuaxewfq3r39r5af5m0z8mmftmz4kl3l3gysk3",
			"reward": {
				"denom": "ucmdx",
				"amount": "6879"
			}
		},
		{
			"address": "comdex1zl093c99c3cfgymz80gh03d9l5dlpd0zrsjkxv",
			"reward": {
				"denom": "ucmdx",
				"amount": "18270982"
			}
		},
		{
			"address": "comdex1rrrq4pgg3sdf72pk9w7sx9y5pyner69vkpfzm6",
			"reward": {
				"denom": "ucmdx",
				"amount": "63806"
			}
		},
		{
			"address": "comdex1rrrlpgf7fnlcs2ehh6dyesgs5v6pxdur7sff7k",
			"reward": {
				"denom": "ucmdx",
				"amount": "2305"
			}
		},
		{
			"address": "comdex1rrns6qd70dl0zuwlnytf90tnyyvux7t5pe36sp",
			"reward": {
				"denom": "ucmdx",
				"amount": "1288337"
			}
		},
		{
			"address": "comdex1r92pra66962r6k7vaynyfs7mwcvjehwzqefdru",
			"reward": {
				"denom": "ucmdx",
				"amount": "11292"
			}
		},
		{
			"address": "comdex1rgf4dznny6fx60yjwmq27xxs42vl34ylmpf2k2",
			"reward": {
				"denom": "ucmdx",
				"amount": "911093"
			}
		},
		{
			"address": "comdex1rdg7ec7a8men3dvxnh5afkpdxcevqkrknv4h0w",
			"reward": {
				"denom": "ucmdx",
				"amount": "1546890"
			}
		},
		{
			"address": "comdex1r55admgr2pnacykfnsng653suqwt352wjsaccw",
			"reward": {
				"denom": "ucmdx",
				"amount": "3104"
			}
		},
		{
			"address": "comdex1reeycz4d4pu4fddzqafzsyh6vvjp3nflp84xpp",
			"reward": {
				"denom": "ucmdx",
				"amount": "322352"
			}
		},
		{
			"address": "comdex1remxeunduu0rs2e6f74jshl4u98q0sce69navs",
			"reward": {
				"denom": "ucmdx",
				"amount": "12557"
			}
		},
		{
			"address": "comdex1ruzff99wnlky2vx02dwdncr42rq0z2n6hyjp5t",
			"reward": {
				"denom": "ucmdx",
				"amount": "5272233"
			}
		},
		{
			"address": "comdex1ypk8p72lqfezjpypxwu6n69de6qzm2sycvspfj",
			"reward": {
				"denom": "ucmdx",
				"amount": "357073"
			}
		},
		{
			"address": "comdex1yzrurgrlqhs60f6amdntrv8rd0hq8rvg7ypq9g",
			"reward": {
				"denom": "ucmdx",
				"amount": "621921"
			}
		},
		{
			"address": "comdex1yxw7a7f5qxnhchzyfu0gx2j5nrhhachmsw33cm",
			"reward": {
				"denom": "ucmdx",
				"amount": "902143"
			}
		},
		{
			"address": "comdex1yg36eryq2zgq23wnhw3fjs4ax45slw7ar3pwuu",
			"reward": {
				"denom": "ucmdx",
				"amount": "1923705"
			}
		},
		{
			"address": "comdex1yngrpun2rggr5rk65y7dvtnkff62n8hrlcltcr",
			"reward": {
				"denom": "ucmdx",
				"amount": "24187"
			}
		},
		{
			"address": "comdex1yhxf6pk7kv9738psnz60vkrvv5jcd39u3relvy",
			"reward": {
				"denom": "ucmdx",
				"amount": "1032144"
			}
		},
		{
			"address": "comdex1yhkzamcwhn7v590x8ezjterj97vl5vtq3vc5c4",
			"reward": {
				"denom": "ucmdx",
				"amount": "103160"
			}
		},
		{
			"address": "comdex1yllsmhc44p2m0dev5wz45wlcryxned50f7xxmy",
			"reward": {
				"denom": "ucmdx",
				"amount": "2771"
			}
		},
		{
			"address": "comdex19qd4s58w5evd94u0cs5j2pzs8g6hs9qults78r",
			"reward": {
				"denom": "ucmdx",
				"amount": "13170"
			}
		},
		{
			"address": "comdex198u54gttp4ecwmg4putuvz379at7vxfne8z333",
			"reward": {
				"denom": "ucmdx",
				"amount": "969"
			}
		},
		{
			"address": "comdex19fddhpk2errrmvd4chp4lmx6rl75ydqjwdpxkp",
			"reward": {
				"denom": "ucmdx",
				"amount": "69422"
			}
		},
		{
			"address": "comdex19wwga4qvwmf2yzrz5vqxycf75dlds92hrx7dcz",
			"reward": {
				"denom": "ucmdx",
				"amount": "78793"
			}
		},
		{
			"address": "comdex195929rxw30rqjrcfvfc8ncsgugnpv49r2p6yjx",
			"reward": {
				"denom": "ucmdx",
				"amount": "2093244"
			}
		},
		{
			"address": "comdex19546utgx69m27uc3wrpx8956f3xyhruqscfnf7",
			"reward": {
				"denom": "ucmdx",
				"amount": "76606"
			}
		},
		{
			"address": "comdex195l5nezqwemae27qshqzsedpyjgpjvpejxqfqc",
			"reward": {
				"denom": "ucmdx",
				"amount": "4176"
			}
		},
		{
			"address": "comdex196tq2xzaphu4d0e7wcr3yccjcyztnzw7m6vczw",
			"reward": {
				"denom": "ucmdx",
				"amount": "21126"
			}
		},
		{
			"address": "comdex1xqgyn7q534lckptdncvhwpzv09tfrsxrx988dx",
			"reward": {
				"denom": "ucmdx",
				"amount": "159044"
			}
		},
		{
			"address": "comdex1xyg6scf4t8qq8karduxhump89hyrvl3v6qtaaf",
			"reward": {
				"denom": "ucmdx",
				"amount": "2436"
			}
		},
		{
			"address": "comdex1xxznmq3a0647kygevplz9ltlu4fzw2uycre26u",
			"reward": {
				"denom": "ucmdx",
				"amount": "28465"
			}
		},
		{
			"address": "comdex1xxzc8awwllqq4waxtej8u8vlyr89lmuyhh8h4n",
			"reward": {
				"denom": "ucmdx",
				"amount": "307250"
			}
		},
		{
			"address": "comdex1xt8rd6pr2tkqqgxkzrdq4yy2tsyejllsdcwd6r",
			"reward": {
				"denom": "ucmdx",
				"amount": "748426927"
			}
		},
		{
			"address": "comdex1x30p6j2gjputppsy43gnac86kn9czn9dq2rny5",
			"reward": {
				"denom": "ucmdx",
				"amount": "3371557"
			}
		},
		{
			"address": "comdex1xjapfntlvpz2uf738lud7yvxcutkvk24wkte24",
			"reward": {
				"denom": "ucmdx",
				"amount": "40904"
			}
		},
		{
			"address": "comdex1xkz6wgg7tfhafc6ler46uuw8m3fv04ttxkjyvd",
			"reward": {
				"denom": "ucmdx",
				"amount": "3045901"
			}
		},
		{
			"address": "comdex1xkua6dnd0067t2wt0agluwqdc4avw2lfkhgyxn",
			"reward": {
				"denom": "ucmdx",
				"amount": "11805"
			}
		},
		{
			"address": "comdex1xhk3gjwle8na9fqh9qpqmsheuvj8qzt7tw6yzc",
			"reward": {
				"denom": "ucmdx",
				"amount": "1037435"
			}
		},
		{
			"address": "comdex1xc7jh0xqk0xeyderel53lkhldygfu0sssnmu8w",
			"reward": {
				"denom": "ucmdx",
				"amount": "97969"
			}
		},
		{
			"address": "comdex18qhfpt0l6hv7ay2atu58z5creekqrxf2v8kc99",
			"reward": {
				"denom": "ucmdx",
				"amount": "143"
			}
		},
		{
			"address": "comdex18xr3rvvx7skkglphes7zrmtdq4ces92dg2pk7q",
			"reward": {
				"denom": "ucmdx",
				"amount": "40973"
			}
		},
		{
			"address": "comdex18gs29dp5c0jvsq8x0z994vwr90yr637asydkas",
			"reward": {
				"denom": "ucmdx",
				"amount": "267783"
			}
		},
		{
			"address": "comdex18g6ghwrtjk90j3l8mc40p2pgzrgs7jv597nsuk",
			"reward": {
				"denom": "ucmdx",
				"amount": "180978"
			}
		},
		{
			"address": "comdex18fk8pzdg36gwnkheuwaekwag0p69ed8ctl4pu9",
			"reward": {
				"denom": "ucmdx",
				"amount": "9508"
			}
		},
		{
			"address": "comdex18dpnjxslnrv6jryek3weqs5ndukussumvczfaz",
			"reward": {
				"denom": "ucmdx",
				"amount": "345"
			}
		},
		{
			"address": "comdex18srx5n9590403f49l72c70v7v2kx8fe77r6zhs",
			"reward": {
				"denom": "ucmdx",
				"amount": "998"
			}
		},
		{
			"address": "comdex18lnd2qfx3recvrmk865twk32jnllgacrla7tn3",
			"reward": {
				"denom": "ucmdx",
				"amount": "440638"
			}
		},
		{
			"address": "comdex1gfysn6yvulu5yekt5nf68nxljh58yeqmk7qnrl",
			"reward": {
				"denom": "ucmdx",
				"amount": "131077"
			}
		},
		{
			"address": "comdex1gnveq6d27e03ndd3rhtqzkz5f2lgmqdxugn8rl",
			"reward": {
				"denom": "ucmdx",
				"amount": "31050"
			}
		},
		{
			"address": "comdex1gkcj3sa7fyjl9qepqjuaskq0qa8gvejqn7yfv9",
			"reward": {
				"denom": "ucmdx",
				"amount": "106682"
			}
		},
		{
			"address": "comdex1gep50gr4jpg7qju2epckc03up0t9mjmldjeayh",
			"reward": {
				"denom": "ucmdx",
				"amount": "16715"
			}
		},
		{
			"address": "comdex1g692jde5e7nqpml4rmxgge22m3qvft8ugnnrz8",
			"reward": {
				"denom": "ucmdx",
				"amount": "4081036"
			}
		},
		{
			"address": "comdex1g7y7999pyv4xj5g6hkhwfuyrc09lvuqw0tv4jy",
			"reward": {
				"denom": "ucmdx",
				"amount": "2352"
			}
		},
		{
			"address": "comdex1g7wuulppm7m7gy83skac6wf9thejaqt2cnusvw",
			"reward": {
				"denom": "ucmdx",
				"amount": "1149827"
			}
		},
		{
			"address": "comdex1fq5ftjcrddtjs4ycqd50849jalpyf95nr7cqs5",
			"reward": {
				"denom": "ucmdx",
				"amount": "957643"
			}
		},
		{
			"address": "comdex1f88kcq3ramnam9937p9zf8l5twujr5jxzjvwl5",
			"reward": {
				"denom": "ucmdx",
				"amount": "631"
			}
		},
		{
			"address": "comdex1f28hyqjfhldzp0mdt7vss3mt85xn7kd35tj8h6",
			"reward": {
				"denom": "ucmdx",
				"amount": "24457"
			}
		},
		{
			"address": "comdex1faqm24yzk0dunnr99uvszes0shzqcsqejz7uct",
			"reward": {
				"denom": "ucmdx",
				"amount": "251611"
			}
		},
		{
			"address": "comdex12qdvcku276ydt6katndy5snpkq02rsra9tvgvw",
			"reward": {
				"denom": "ucmdx",
				"amount": "11493054"
			}
		},
		{
			"address": "comdex122esu76xehp8sq9t88kcn666ejjum5g5ynxu0k",
			"reward": {
				"denom": "ucmdx",
				"amount": "1177535"
			}
		},
		{
			"address": "comdex1203temyrsh29uucra5kxq7dssw86jljmkjheya",
			"reward": {
				"denom": "ucmdx",
				"amount": "303793"
			}
		},
		{
			"address": "comdex12c2zn3mt736lmzsj6a5c8g97fu05q53km27776",
			"reward": {
				"denom": "ucmdx",
				"amount": "378902"
			}
		},
		{
			"address": "comdex12lzvcacrthkmx33dkxfce9pffnuazz5emya9ae",
			"reward": {
				"denom": "ucmdx",
				"amount": "73723"
			}
		},
		{
			"address": "comdex1tp6jmrrf6hlhpwa8e4e590jvyvsdkvv8da0nuc",
			"reward": {
				"denom": "ucmdx",
				"amount": "26328549"
			}
		},
		{
			"address": "comdex1t987w4w68wqa2tgdaxjh9lpfaxaeym85uwy0z3",
			"reward": {
				"denom": "ucmdx",
				"amount": "2052857"
			}
		},
		{
			"address": "comdex1t9dwh03pyhe94u6x9d7a3d48dd2kcl63s7d2h6",
			"reward": {
				"denom": "ucmdx",
				"amount": "506791175"
			}
		},
		{
			"address": "comdex1tfvjatveyl9wwf4shkj0rr4tegxhfv97wlzr9s",
			"reward": {
				"denom": "ucmdx",
				"amount": "2207"
			}
		},
		{
			"address": "comdex1tttqjvg3peg69cd0ln34qycs9agtk3wqknn554",
			"reward": {
				"denom": "ucmdx",
				"amount": "308724"
			}
		},
		{
			"address": "comdex1t5zgnfz0jrvflywjmgs95rey3un57n427gvsds",
			"reward": {
				"denom": "ucmdx",
				"amount": "17648028"
			}
		},
		{
			"address": "comdex1tk8g0x5lzg6drjnnm8tfa8tn7nu6zksk3yz6x9",
			"reward": {
				"denom": "ucmdx",
				"amount": "379145721"
			}
		},
		{
			"address": "comdex1tmttlndvk6adkf6c9ugkda92mwpel3whxetsud",
			"reward": {
				"denom": "ucmdx",
				"amount": "78213"
			}
		},
		{
			"address": "comdex1tlmaprfnafdvczyn3xpgvstd9uz3rtsm9eh27c",
			"reward": {
				"denom": "ucmdx",
				"amount": "1550700"
			}
		},
		{
			"address": "comdex1vqd6hxuqurm7phg586vw6rv28rga4edl74a3ze",
			"reward": {
				"denom": "ucmdx",
				"amount": "75895879"
			}
		},
		{
			"address": "comdex1vxg90vp5pq5rktf8wze8hkynvj0mnykru3jn29",
			"reward": {
				"denom": "ucmdx",
				"amount": "16474"
			}
		},
		{
			"address": "comdex1vvf8vrzhm4x8vajdkv3m6h0wk8u0kzsw34n2u0",
			"reward": {
				"denom": "ucmdx",
				"amount": "120476"
			}
		},
		{
			"address": "comdex1vhmvws8nmfjsfvexmqcdntyzsd3t45wvln4j44",
			"reward": {
				"denom": "ucmdx",
				"amount": "2248672"
			}
		},
		{
			"address": "comdex1vemmttg4c0lsqwsdhrjc24ktaahfnh9hqv8qls",
			"reward": {
				"denom": "ucmdx",
				"amount": "657804"
			}
		},
		{
			"address": "comdex1dp6chezxk8zvpc0rl4xmrsam6wd6xfnwqh5psk",
			"reward": {
				"denom": "ucmdx",
				"amount": "11988522"
			}
		},
		{
			"address": "comdex1d20uz30evswtucdn77dflg4pw5ag8frp2uqkv5",
			"reward": {
				"denom": "ucmdx",
				"amount": "65893"
			}
		},
		{
			"address": "comdex1d24ddxn9lqxvhk04xjlufa8hnzcmk6n83qvg7n",
			"reward": {
				"denom": "ucmdx",
				"amount": "789840"
			}
		},
		{
			"address": "comdex1dnnl0pj4pruwlkkp6gsguhfh88e492jgxjr859",
			"reward": {
				"denom": "ucmdx",
				"amount": "150624"
			}
		},
		{
			"address": "comdex1dcvlnewq7ktr5kd47755d4x7jvzsn3aewj972j",
			"reward": {
				"denom": "ucmdx",
				"amount": "170655"
			}
		},
		{
			"address": "comdex1dmddduuvnaf87yk2j6yf2f7ctjns8n4x755vqq",
			"reward": {
				"denom": "ucmdx",
				"amount": "32458"
			}
		},
		{
			"address": "comdex1dl2gcq66twxydx4tumt9qcrh302u9grt3m3mr9",
			"reward": {
				"denom": "ucmdx",
				"amount": "182175"
			}
		},
		{
			"address": "comdex1wqvfvc7emqkwtkegkn7ugxaqhrfh3z7ndl7caf",
			"reward": {
				"denom": "ucmdx",
				"amount": "27413"
			}
		},
		{
			"address": "comdex1w8985yzl94c7vhdzdn5tduz2qjyr5s36kx2z7x",
			"reward": {
				"denom": "ucmdx",
				"amount": "3037975"
			}
		},
		{
			"address": "comdex1w80g5swy7armj8ls3r433v0lknxyrp53w4n2xv",
			"reward": {
				"denom": "ucmdx",
				"amount": "89665700"
			}
		},
		{
			"address": "comdex1w83rmdhrqudkj5qgs0c8lea6ypacp5534mrv3u",
			"reward": {
				"denom": "ucmdx",
				"amount": "3497307"
			}
		},
		{
			"address": "comdex1w3dgg0tsunvjmgcd2zfzwhhl5nv60h8zzdrkl3",
			"reward": {
				"denom": "ucmdx",
				"amount": "605122"
			}
		},
		{
			"address": "comdex1w5lep3d53p5dtkg37gerq6qxdlagykyryta989",
			"reward": {
				"denom": "ucmdx",
				"amount": "60202473"
			}
		},
		{
			"address": "comdex1w48qsnnes07n9ust93zxk6mh4g0ynz6rnvgxyn",
			"reward": {
				"denom": "ucmdx",
				"amount": "3953786"
			}
		},
		{
			"address": "comdex1wmy300g9dckyynp83u4nf9nnl5g4efy5zxtnsu",
			"reward": {
				"denom": "ucmdx",
				"amount": "368"
			}
		},
		{
			"address": "comdex1wlsqpmtt04hdr80y6mn9khgjrxgrfngw58dkeu",
			"reward": {
				"denom": "ucmdx",
				"amount": "988922"
			}
		},
		{
			"address": "comdex10rsp2tdhk4pst94qx84pvgrl08mkdf620mczwj",
			"reward": {
				"denom": "ucmdx",
				"amount": "111313"
			}
		},
		{
			"address": "comdex10ramtpjevxutxfmw7mq8chpsx6puelcsj25ypj",
			"reward": {
				"denom": "ucmdx",
				"amount": "14369"
			}
		},
		{
			"address": "comdex10y6p0f0srtsgxqu9glu60wr2962jn6f383gxkk",
			"reward": {
				"denom": "ucmdx",
				"amount": "10047"
			}
		},
		{
			"address": "comdex10xz7jwc3m2z3uxtw6zkjusxsdzu7mtleqf5g5x",
			"reward": {
				"denom": "ucmdx",
				"amount": "114895"
			}
		},
		{
			"address": "comdex10hg23wl72adrj5vxxfln63m4pu9vxhewx9smrc",
			"reward": {
				"denom": "ucmdx",
				"amount": "44963"
			}
		},
		{
			"address": "comdex106v0dgp92mwgerjph7aw84wwkutzuq4zukyl4q",
			"reward": {
				"denom": "ucmdx",
				"amount": "297864"
			}
		},
		{
			"address": "comdex10mrevq3wafxu3v7e4rpkvcxskk660v899kh9lt",
			"reward": {
				"denom": "ucmdx",
				"amount": "86237"
			}
		},
		{
			"address": "comdex10mva4dx3kap9r92v5j8g6ea9axqkg96zkqrcrm",
			"reward": {
				"denom": "ucmdx",
				"amount": "1674525"
			}
		},
		{
			"address": "comdex10l2h578tq2fsp8fnnvrwg94s2uqtaz2qc90tr4",
			"reward": {
				"denom": "ucmdx",
				"amount": "600098"
			}
		},
		{
			"address": "comdex10lng7wks2wnfykacuajydhglz23aupygqnksy9",
			"reward": {
				"denom": "ucmdx",
				"amount": "630450"
			}
		},
		{
			"address": "comdex1szv2gzvmk27sgpeajhtjdmexl38kxpkxf5ptz3",
			"reward": {
				"denom": "ucmdx",
				"amount": "1750646"
			}
		},
		{
			"address": "comdex1s9xp57kdmu7t54m7uf4wc4rx5hd5dwjxwc24g0",
			"reward": {
				"denom": "ucmdx",
				"amount": "44920615"
			}
		},
		{
			"address": "comdex1s8nzgwexdz8laux65600tkq6mq8pr6h5ya6fjw",
			"reward": {
				"denom": "ucmdx",
				"amount": "28965"
			}
		},
		{
			"address": "comdex1sgu2npk4s87l2rm0dphnl29htqvvznm5ut8e3j",
			"reward": {
				"denom": "ucmdx",
				"amount": "20813"
			}
		},
		{
			"address": "comdex1svfpddvanq0p4ejp0v9myxt5zrfy6mdy4pc4h4",
			"reward": {
				"denom": "ucmdx",
				"amount": "51225"
			}
		},
		{
			"address": "comdex1sdn6us3ayfzluf4tq3gvvmq9x86f9csgvhhdeq",
			"reward": {
				"denom": "ucmdx",
				"amount": "131366"
			}
		},
		{
			"address": "comdex1snnuntk4xy6zcm2jmtsqefnzu9tykp2p4p7n03",
			"reward": {
				"denom": "ucmdx",
				"amount": "409199"
			}
		},
		{
			"address": "comdex1shk2aw96fd6l7yvdr2td33rstey4j72cst6trm",
			"reward": {
				"denom": "ucmdx",
				"amount": "2584490"
			}
		},
		{
			"address": "comdex13qg8wcye68jy2pzxlft33kp6q5u846c2cm5s0d",
			"reward": {
				"denom": "ucmdx",
				"amount": "861276"
			}
		},
		{
			"address": "comdex13p5dt20sggkxpj9w0lgenk87tuj8ykz0e38dzh",
			"reward": {
				"denom": "ucmdx",
				"amount": "241659"
			}
		},
		{
			"address": "comdex13xjrerteau5qrtr8a53zd0aj969xzw5qxmf4kq",
			"reward": {
				"denom": "ucmdx",
				"amount": "1547372"
			}
		},
		{
			"address": "comdex13fyrv9huq880vkc3unl45grn2klgt5qqvqaeu7",
			"reward": {
				"denom": "ucmdx",
				"amount": "146469"
			}
		},
		{
			"address": "comdex1323uenwvychnv4gplx9895w43tcd66lut9d0kd",
			"reward": {
				"denom": "ucmdx",
				"amount": "117884"
			}
		},
		{
			"address": "comdex1304vt66qjuqqwfqzh5y4tdhfy083w4rmfqzgzy",
			"reward": {
				"denom": "ucmdx",
				"amount": "23462512"
			}
		},
		{
			"address": "comdex134ultcrapjnu8nxttt5rlc05tgeex2janww2xk",
			"reward": {
				"denom": "ucmdx",
				"amount": "357714"
			}
		},
		{
			"address": "comdex1jq36nll00t5a4dvmqzcht56e2qp25w7dxq2pqv",
			"reward": {
				"denom": "ucmdx",
				"amount": "62772"
			}
		},
		{
			"address": "comdex1jyp2hj8znp24uhfrlzjtqr6ftlfs788zmnxey8",
			"reward": {
				"denom": "ucmdx",
				"amount": "6641"
			}
		},
		{
			"address": "comdex1jxquz5w36nwtmuql40l7wkzgwx2m606pn8jtkz",
			"reward": {
				"denom": "ucmdx",
				"amount": "606581"
			}
		},
		{
			"address": "comdex1jfrs6r5q9shj0w27ydm3gn86pa44fh9acfrzxu",
			"reward": {
				"denom": "ucmdx",
				"amount": "28430"
			}
		},
		{
			"address": "comdex1jf3jh0hcnf8d64slauzl5s83zy3yhl09zscqf8",
			"reward": {
				"denom": "ucmdx",
				"amount": "2493"
			}
		},
		{
			"address": "comdex1judt9894tekjl27axpgsczezlmk88egxvnsfss",
			"reward": {
				"denom": "ucmdx",
				"amount": "2109080"
			}
		},
		{
			"address": "comdex1ja8v5zer5ucdhslkdhyfd3j38vhun39k9xhp3r",
			"reward": {
				"denom": "ucmdx",
				"amount": "779837"
			}
		},
		{
			"address": "comdex1nr63devsvvh5furqaw8zsevch8e77nymaza5hu",
			"reward": {
				"denom": "ucmdx",
				"amount": "242180"
			}
		},
		{
			"address": "comdex1nx3chh7ju3l28w8nsu2wpf3hhqfkg7lk08xsuq",
			"reward": {
				"denom": "ucmdx",
				"amount": "18898"
			}
		},
		{
			"address": "comdex1n8gqw8jdh73r7w6nuvr0wy6javsjgwt338xn9x",
			"reward": {
				"denom": "ucmdx",
				"amount": "328341"
			}
		},
		{
			"address": "comdex1n8gr2p44e84x22lqe0q88eg4n37wxcgpsr997c",
			"reward": {
				"denom": "ucmdx",
				"amount": "9835789"
			}
		},
		{
			"address": "comdex1n4c56vddeqg67ktukprkteqdmpph2ck02qfhst",
			"reward": {
				"denom": "ucmdx",
				"amount": "13532"
			}
		},
		{
			"address": "comdex1nc0hm70ujjdc4z24fzl3t27fweerwdtqyyy8vz",
			"reward": {
				"denom": "ucmdx",
				"amount": "5357619"
			}
		},
		{
			"address": "comdex1nen8k454qmeuh98583vyknwt20k09v8wfv69dj",
			"reward": {
				"denom": "ucmdx",
				"amount": "771467"
			}
		},
		{
			"address": "comdex15fqhy5usyr2qsnd4qr4rp86wnpy58jqzaxd970",
			"reward": {
				"denom": "ucmdx",
				"amount": "16871"
			}
		},
		{
			"address": "comdex153zjzc704l5e8l6npdp2lnw8q7y0tjmpdth6y4",
			"reward": {
				"denom": "ucmdx",
				"amount": "1984594"
			}
		},
		{
			"address": "comdex154fewkzd835vgy5d4kk7ts8sxnrm3df7pzw9zw",
			"reward": {
				"denom": "ucmdx",
				"amount": "36038"
			}
		},
		{
			"address": "comdex15ezd9p9ha0fvp5mtuct8rfrttsty4a2erg4vju",
			"reward": {
				"denom": "ucmdx",
				"amount": "16695"
			}
		},
		{
			"address": "comdex15avn2fq0293lmjv6rjkjltj7jfzkqf7h0c89fl",
			"reward": {
				"denom": "ucmdx",
				"amount": "497003"
			}
		},
		{
			"address": "comdex1572f07js6gnnd9p0fer9c87f5shkr05dmfemn0",
			"reward": {
				"denom": "ucmdx",
				"amount": "123431"
			}
		},
		{
			"address": "comdex1489a36405052np9vgf4ajzam4taphu39vy7erh",
			"reward": {
				"denom": "ucmdx",
				"amount": "387207"
			}
		},
		{
			"address": "comdex14v0e9hrth8a6mpyfwdhj27j7r06c0dcuejggly",
			"reward": {
				"denom": "ucmdx",
				"amount": "41462"
			}
		},
		{
			"address": "comdex14sxe4r62gvr9z8s7hgyuqe2lv9kdsncwnus2ds",
			"reward": {
				"denom": "ucmdx",
				"amount": "132354"
			}
		},
		{
			"address": "comdex1434zluxmfj85vwpc3stk03wrcuzrqllcqgewk6",
			"reward": {
				"denom": "ucmdx",
				"amount": "174830"
			}
		},
		{
			"address": "comdex144rnmenu77zphrjmf8ylpk0zr7dq4rt2e25005",
			"reward": {
				"denom": "ucmdx",
				"amount": "48549"
			}
		},
		{
			"address": "comdex14kc7pmt50a3ju3g3kk6t5v3arrznrjfh5pfud4",
			"reward": {
				"denom": "ucmdx",
				"amount": "186069"
			}
		},
		{
			"address": "comdex1kpaa3fwz0z3q00vxlzmkn0p45gkkhv0nraqtj5",
			"reward": {
				"denom": "ucmdx",
				"amount": "187013"
			}
		},
		{
			"address": "comdex1kykrr9f2uw3jxcvs59u8g0guglywlre7hqj7zz",
			"reward": {
				"denom": "ucmdx",
				"amount": "1097"
			}
		},
		{
			"address": "comdex1kycslqsw5z3ysuw2nrrzam266pxplx8j9zun6k",
			"reward": {
				"denom": "ucmdx",
				"amount": "11028869"
			}
		},
		{
			"address": "comdex1kxzcgtrd55zcnajjx8890c82jflvwvzx4turcl",
			"reward": {
				"denom": "ucmdx",
				"amount": "7680"
			}
		},
		{
			"address": "comdex1k0aq3m7reynyeme3fh4fkx5y04q00h8cnrv042",
			"reward": {
				"denom": "ucmdx",
				"amount": "2804026"
			}
		},
		{
			"address": "comdex1k325h8hhqclgfxmexe3crjrkavqn3pqw0xt6e8",
			"reward": {
				"denom": "ucmdx",
				"amount": "25421"
			}
		},
		{
			"address": "comdex1kjtmkagp4fz7tx9nsu5a2xfnz6m50nudnwzum9",
			"reward": {
				"denom": "ucmdx",
				"amount": "119414"
			}
		},
		{
			"address": "comdex1ke6r9t67jg2wnq3z5zany8v93rq06mm4mz4pfx",
			"reward": {
				"denom": "ucmdx",
				"amount": "454288927"
			}
		},
		{
			"address": "comdex1k6r4k7pdc948wmu2rkvl4r7fn0cq859ngrprjw",
			"reward": {
				"denom": "ucmdx",
				"amount": "62862"
			}
		},
		{
			"address": "comdex1h99dx0qu6e9xu6q86gzpgh550geqawjm2nq7p8",
			"reward": {
				"denom": "ucmdx",
				"amount": "1044235"
			}
		},
		{
			"address": "comdex1hx2gy2qnx7mf32cdc47nvmk2gvcft5s2s2k5wh",
			"reward": {
				"denom": "ucmdx",
				"amount": "2546"
			}
		},
		{
			"address": "comdex1h86x3s99aurjzxzetjcgjggrenmf80smaxfgqw",
			"reward": {
				"denom": "ucmdx",
				"amount": "45325"
			}
		},
		{
			"address": "comdex1hgv7e8pjk8jqtfwsypq34y3puqxtald2xswvng",
			"reward": {
				"denom": "ucmdx",
				"amount": "6104523"
			}
		},
		{
			"address": "comdex1hfp05gmced7ks78vnv8hy86ls3jmc7n2fks9me",
			"reward": {
				"denom": "ucmdx",
				"amount": "251925"
			}
		},
		{
			"address": "comdex1htjazgc22yq9fd8xq2q2tesp9gy3lfc3g9p0d2",
			"reward": {
				"denom": "ucmdx",
				"amount": "7906"
			}
		},
		{
			"address": "comdex1h3u87jxu2yrnud2czeyusk9scj2vpknznjrvp3",
			"reward": {
				"denom": "ucmdx",
				"amount": "64604"
			}
		},
		{
			"address": "comdex1h4zffvsfnydkptg2mgpm0hquxv2p9zjyhztccx",
			"reward": {
				"denom": "ucmdx",
				"amount": "1236"
			}
		},
		{
			"address": "comdex1cqku07qz9das0pcnys74cq64lflqen09weunhn",
			"reward": {
				"denom": "ucmdx",
				"amount": "89404232"
			}
		},
		{
			"address": "comdex1c90kyse82th8v79j7dr6yuw4nrkhkp0pgjgdmj",
			"reward": {
				"denom": "ucmdx",
				"amount": "81865"
			}
		},
		{
			"address": "comdex1cggk4nd3jv55ddlwe6xfdjdqsk4y4dvgnztlvt",
			"reward": {
				"denom": "ucmdx",
				"amount": "4234116"
			}
		},
		{
			"address": "comdex1c347zr3z9ck7kwf9csxnn5vj5d06xcng4jd7j6",
			"reward": {
				"denom": "ucmdx",
				"amount": "39200"
			}
		},
		{
			"address": "comdex1c36s76tqahjd959g55v6fzt4x3ctmnsjs2k4z3",
			"reward": {
				"denom": "ucmdx",
				"amount": "4196763"
			}
		},
		{
			"address": "comdex1c5zre8yzq4s5c5htqffjkf6puf2mmlflaruhpn",
			"reward": {
				"denom": "ucmdx",
				"amount": "21274"
			}
		},
		{
			"address": "comdex1c5cd25l2ujusap22umvy2sk6s4nd9gvysm5r0w",
			"reward": {
				"denom": "ucmdx",
				"amount": "89015"
			}
		},
		{
			"address": "comdex1ch5kezhwsr0kap7h86v98exjhme4xc3fr6juwz",
			"reward": {
				"denom": "ucmdx",
				"amount": "9166"
			}
		},
		{
			"address": "comdex1c77nzlhwsqyejev4cezpx83yj6u88fvtfgffeg",
			"reward": {
				"denom": "ucmdx",
				"amount": "2054138"
			}
		},
		{
			"address": "comdex1ezgedsezl60q2q67fwxk9p0y50tuzhhvccuywa",
			"reward": {
				"denom": "ucmdx",
				"amount": "5606863"
			}
		},
		{
			"address": "comdex1eymapsam3l5ymas76avez0adlvw3u293hl7slz",
			"reward": {
				"denom": "ucmdx",
				"amount": "5187695"
			}
		},
		{
			"address": "comdex1e49cvjgtf0anp0gpdqhs3tzg24mcgse076sdu5",
			"reward": {
				"denom": "ucmdx",
				"amount": "19680"
			}
		},
		{
			"address": "comdex1ehgxht05ymwtw46l0vz8dwpryhf47stlg98z05",
			"reward": {
				"denom": "ucmdx",
				"amount": "37862"
			}
		},
		{
			"address": "comdex1e7c4qnthxz2pzkjkm3k5dflqyn6s3522cy3ng6",
			"reward": {
				"denom": "ucmdx",
				"amount": "556731"
			}
		},
		{
			"address": "comdex16zh0ak7vymdu3ahuvl3as78q6p7hdvaqtrsrs3",
			"reward": {
				"denom": "ucmdx",
				"amount": "890"
			}
		},
		{
			"address": "comdex16r88f2zpw5ujkucl0kx7pdgj3l26547jknymr6",
			"reward": {
				"denom": "ucmdx",
				"amount": "26675"
			}
		},
		{
			"address": "comdex168qsl9q8ltxhj48lrpwxh3q4wyue2nk8h9d2kx",
			"reward": {
				"denom": "ucmdx",
				"amount": "22812"
			}
		},
		{
			"address": "comdex16fa9eqd8fnckald372th7wjm46cd3cvgj46nk9",
			"reward": {
				"denom": "ucmdx",
				"amount": "1628500"
			}
		},
		{
			"address": "comdex16tvh6wfyz5t3fskvftxns49suv2avg2syfc22e",
			"reward": {
				"denom": "ucmdx",
				"amount": "18386"
			}
		},
		{
			"address": "comdex16t4xlrg03snrerg5pqx7k6zwlml84vjnhs8aua",
			"reward": {
				"denom": "ucmdx",
				"amount": "219829"
			}
		},
		{
			"address": "comdex16n433589j8duk3xet7njrrsx68aj8thwwws5y8",
			"reward": {
				"denom": "ucmdx",
				"amount": "4878"
			}
		},
		{
			"address": "comdex167wsxxcte77qm95art8z8whmlm0klh3xmurylh",
			"reward": {
				"denom": "ucmdx",
				"amount": "49057"
			}
		},
		{
			"address": "comdex1mrnzl2cf0tp4pp2jg5yr5vetf2gv9645dkc56v",
			"reward": {
				"denom": "ucmdx",
				"amount": "230754"
			}
		},
		{
			"address": "comdex1mt9q7g276qsjhxkuxl73kmv9zea0nscdmt809k",
			"reward": {
				"denom": "ucmdx",
				"amount": "109605"
			}
		},
		{
			"address": "comdex1md56l4zt95f0lve0uur3np49cr58d5u8vd6dfg",
			"reward": {
				"denom": "ucmdx",
				"amount": "28024"
			}
		},
		{
			"address": "comdex1m4qjakxzf9qlgpjpd0jws4qk2jq9gnttzmupk4",
			"reward": {
				"denom": "ucmdx",
				"amount": "2912"
			}
		},
		{
			"address": "comdex1mh4ljwk9zrhzj2xnyw7dlw2e5wwfj46s96xx8j",
			"reward": {
				"denom": "ucmdx",
				"amount": "137582"
			}
		},
		{
			"address": "comdex1memcf6sqd00h6yzf0e8e8e707tjmszpd8632er",
			"reward": {
				"denom": "ucmdx",
				"amount": "59540"
			}
		},
		{
			"address": "comdex1uyphdztvh3r4jv23utrht7qtgf0ref34l2y7ra",
			"reward": {
				"denom": "ucmdx",
				"amount": "782609"
			}
		},
		{
			"address": "comdex1uta2mx7jdxc0tydznryd8a8we6p3xntsgee3j0",
			"reward": {
				"denom": "ucmdx",
				"amount": "190467"
			}
		},
		{
			"address": "comdex1utl04swwlt9xsr9c2vdx2nkea4plp66s9pm8zc",
			"reward": {
				"denom": "ucmdx",
				"amount": "309"
			}
		},
		{
			"address": "comdex1ujpheelct6cpqg5r6ayp9raan8qr09dv6jdpg2",
			"reward": {
				"denom": "ucmdx",
				"amount": "8740"
			}
		},
		{
			"address": "comdex1uh7qk3mxa0hd4f3rr55rq6jty8szf3hev5j3mn",
			"reward": {
				"denom": "ucmdx",
				"amount": "62186"
			}
		},
		{
			"address": "comdex1axqe3vc09fu6908ze7vcz0ns8deyszd28jm9e5",
			"reward": {
				"denom": "ucmdx",
				"amount": "323520"
			}
		},
		{
			"address": "comdex1agprce4jnftapsg2k4vkyqq3ez6dpjdl8uwfsg",
			"reward": {
				"denom": "ucmdx",
				"amount": "24"
			}
		},
		{
			"address": "comdex1a3d39tv7qzmnjqwu7x8frrwphxq6y98jca9ppu",
			"reward": {
				"denom": "ucmdx",
				"amount": "10756"
			}
		},
		{
			"address": "comdex1ajes8srjm2gu4t49pjk5fqtu0f5edxhd2lst4r",
			"reward": {
				"denom": "ucmdx",
				"amount": "2301"
			}
		},
		{
			"address": "comdex1a5u7gudmk7qqpsm3hs3s2w3rha7f2e7y4w9jwj",
			"reward": {
				"denom": "ucmdx",
				"amount": "10119"
			}
		},
		{
			"address": "comdex1a4vddyr55m0rpe92zk73vq5vurlpgyv2u5jxvq",
			"reward": {
				"denom": "ucmdx",
				"amount": "6529"
			}
		},
		{
			"address": "comdex1acszd5jt3ky65zlaej83lm36f6lu9tw70rmst4",
			"reward": {
				"denom": "ucmdx",
				"amount": "1962480341"
			}
		},
		{
			"address": "comdex17pfkzytqf9aafpp4rucg0xex9csv5nzjghtjh2",
			"reward": {
				"denom": "ucmdx",
				"amount": "18361156"
			}
		},
		{
			"address": "comdex17y9gdgxfv06f67vcterep6s3ug26a5t9c36gwk",
			"reward": {
				"denom": "ucmdx",
				"amount": "838041"
			}
		},
		{
			"address": "comdex179u9fjwnv2kwjctcu079mmtk8et0rnspqw9gps",
			"reward": {
				"denom": "ucmdx",
				"amount": "1549895"
			}
		},
		{
			"address": "comdex17xsk7cma5dg54q2lxdcasvgnhxknyfrlq0w69t",
			"reward": {
				"denom": "ucmdx",
				"amount": "337196"
			}
		},
		{
			"address": "comdex174ml5wq2lea29jrg5jv46thm20vdacmp7nq3k9",
			"reward": {
				"denom": "ucmdx",
				"amount": "1100141"
			}
		},
		{
			"address": "comdex17c6syh2twxa3tqwlrg945kj7exvz825xp5ndjg",
			"reward": {
				"denom": "ucmdx",
				"amount": "795"
			}
		},
		{
			"address": "comdex177gd5wayqxrccrp97g257dd4h73ltxrtyjuy9z",
			"reward": {
				"denom": "ucmdx",
				"amount": "272983"
			}
		},
		{
			"address": "comdex177cepqsjgnwyjrzw0kmp7xjvpg55vhl278g7y2",
			"reward": {
				"denom": "ucmdx",
				"amount": "63793"
			}
		},
		{
			"address": "comdex1lyhe9syygu6nwkf3jsg6vs4xneqektltxxmpr3",
			"reward": {
				"denom": "ucmdx",
				"amount": "118618"
			}
		},
		{
			"address": "comdex1lychryh5ym4lhw88kautnarqu743fzjeukq6l2",
			"reward": {
				"denom": "ucmdx",
				"amount": "3091740"
			}
		},
		{
			"address": "comdex1l9z475x2egvs686j5jtjxm6n8fehlz9ssgrq53",
			"reward": {
				"denom": "ucmdx",
				"amount": "1946830"
			}
		},
		{
			"address": "comdex1l9luwxxlggh7a9z9j3tdz4ydfxhm4gz7r5pc9m",
			"reward": {
				"denom": "ucmdx",
				"amount": "344533"
			}
		},
		{
			"address": "comdex1lspy7yu506jz2qa0npwuut22qayde49uljragz",
			"reward": {
				"denom": "ucmdx",
				"amount": "17921"
			}
		},
		{
			"address": "comdex1lsxwwpghn065ft24at8mdd09sc5v59ujrlamzc",
			"reward": {
				"denom": "ucmdx",
				"amount": "12090867"
			}
		},
		{
			"address": "comdex1lc6sv2d5qd7dsjcchyk50hzu0c0s0cdm7fwrkm",
			"reward": {
				"denom": "ucmdx",
				"amount": "76558"
			}
		},
		{
			"address": "comdex1legedfyfdrf7rf24rwqy8fsydrtmda9t50mfs0",
			"reward": {
				"denom": "ucmdx",
				"amount": "164880"
			}
		},
		{
			"address": "comdex1l7s5a7vp9j7welu48g8ply3eg8mce33ccqyaqk",
			"reward": {
				"denom": "ucmdx",
				"amount": "12450"
			}
		},
		{
			"address": "comdex1zw8u7heex2hy3vk25nwtwfngu9jphjvmsach99",
			"reward": {
				"denom": "ucmdx",
				"amount": "4964180"
			}
		},
		{
			"address": "comdex196knuj9suju45dwvzam2a4sn0hda06un63ma5u",
			"reward": {
				"denom": "ucmdx",
				"amount": "125775969"
			}
		},
		{
			"address": "comdex1kwwm2zln6z8lv30szcya5f9wnrdu6lamnrdj6a",
			"reward": {
				"denom": "ucmdx",
				"amount": "23301353"
			}
		},
		{
			"address": "comdex17rx4caclnlkqwlkq4hq3aq0cvj993pnf0xfs7q",
			"reward": {
				"denom": "ucmdx",
				"amount": "144763848"
			}
		},
		{
			"address": "comdex1zu3pf2azugvx9uun4ypd3a73vhsukhjnjrtmpv",
			"reward": {
				"denom": "ucmdx",
				"amount": "53868841"
			}
		},
		{
			"address": "comdex19m2hz8gdwg495cp59xz7fc28rumanmyym544lj",
			"reward": {
				"denom": "ucmdx",
				"amount": "2290771"
			}
		},
		{
			"address": "comdex1thtzdut0c7kv2q2msxr0hh00wker2w3wgpj55c",
			"reward": {
				"denom": "ucmdx",
				"amount": "2096234"
			}
		},
		{
			"address": "comdex1vqy7p65f6hqaa2n7cuur6jf9m07eq84x8603xz",
			"reward": {
				"denom": "ucmdx",
				"amount": "629204"
			}
		},
		{
			"address": "comdex100855lp5s7g2lm86423p44y4kq8xmy9rjsjwwj",
			"reward": {
				"denom": "ucmdx",
				"amount": "2179155"
			}
		},
		{
			"address": "comdex1sxvmwpsunp2h8759zrlqz34qnemj425xfl9z96",
			"reward": {
				"denom": "ucmdx",
				"amount": "1547775"
			}
		},
		{
			"address": "comdex1yrgrvmawrtvyzf7hhu2m3qchqm80k2gl3uehnr",
			"reward": {
				"denom": "ucmdx",
				"amount": "1820339"
			}
		},
		{
			"address": "comdex19489kp7c2patrw4el7vkq4cham5rq5ugs8tm9g",
			"reward": {
				"denom": "ucmdx",
				"amount": "3452191"
			}
		},
		{
			"address": "comdex1x5w3z8q2eczk6tr48v7dah28fjvkm0huzydtm5",
			"reward": {
				"denom": "ucmdx",
				"amount": "3189851"
			}
		},
		{
			"address": "comdex1fxf23rgnee64y76sprf50mfh65zjctd0g5v00r",
			"reward": {
				"denom": "ucmdx",
				"amount": "2203005"
			}
		},
		{
			"address": "comdex12uwj2p82k97ksvs6dcllktxkr68zsm25gqd9z8",
			"reward": {
				"denom": "ucmdx",
				"amount": "23798"
			}
		},
		{
			"address": "comdex1w0d56frz9e30efe45jqeefx85e69jqlwj32ah4",
			"reward": {
				"denom": "ucmdx",
				"amount": "1361543"
			}
		},
		{
			"address": "comdex1wc4skry7qm78gn9ghjfvmr5yvp9s8j94anhj26",
			"reward": {
				"denom": "ucmdx",
				"amount": "10000488"
			}
		},
		{
			"address": "comdex1sm7c84n4fryl39490zf3gr9550cln6g7vgq8ny",
			"reward": {
				"denom": "ucmdx",
				"amount": "7851258"
			}
		},
		{
			"address": "comdex1c7kd05vjvl4krk506352wfalt8cgnez8qmvytd",
			"reward": {
				"denom": "ucmdx",
				"amount": "4135265"
			}
		},
		{
			"address": "comdex1uw3nsdcz3cj720p5zwggaxzy0zleuew8v7cc6t",
			"reward": {
				"denom": "ucmdx",
				"amount": "61417096"
			}
		},
		{
			"address": "comdex1ay372nh5gh0w8tpzltcn9acenk0al3pz6pc3qr",
			"reward": {
				"denom": "ucmdx",
				"amount": "59901130"
			}
		},
		{
			"address": "comdex1pk59mg6te7vusgh6j6kem32f6ex85829hffmar",
			"reward": {
				"denom": "ucmdx",
				"amount": "15554344"
			}
		},
		{
			"address": "comdex192f3epa9wq9v4zl3n2u3ux7v3z57mt4m3g7rmk",
			"reward": {
				"denom": "ucmdx",
				"amount": "13632987"
			}
		},
		{
			"address": "comdex193tk9u56gu6qyeurhwxugxgwxh3vtcnevyqkhq",
			"reward": {
				"denom": "ucmdx",
				"amount": "1103635"
			}
		},
		{
			"address": "comdex1xu5sm7jv47v4rhcs7dnemsvhtcswj5gpwfgfh0",
			"reward": {
				"denom": "ucmdx",
				"amount": "28846081"
			}
		},
		{
			"address": "comdex183nnun3xl7amzskn3apmnx7h6gwr427cmuj8k8",
			"reward": {
				"denom": "ucmdx",
				"amount": "1345533"
			}
		},
		{
			"address": "comdex1grvqe50wv3t67fud5zr3vqp508h2xdcdz0lkh8",
			"reward": {
				"denom": "ucmdx",
				"amount": "23142161"
			}
		},
		{
			"address": "comdex1ggecdhzqwlkxxkeg727a0ejvz4q6z6gxgs4chh",
			"reward": {
				"denom": "ucmdx",
				"amount": "41710962"
			}
		},
		{
			"address": "comdex12xtvcw94zmd4hnea65xj4tjuyz9mnasulclfxj",
			"reward": {
				"denom": "ucmdx",
				"amount": "1991335"
			}
		},
		{
			"address": "comdex12eedqrmxxsx4s26yzeqs36se5dwwg2tcemwegr",
			"reward": {
				"denom": "ucmdx",
				"amount": "612731"
			}
		},
		{
			"address": "comdex1tvgqvwsf2raywnpxrftul06qsrxwkwwz2mfjg4",
			"reward": {
				"denom": "ucmdx",
				"amount": "511798"
			}
		},
		{
			"address": "comdex1tkelcm4rpx7cec5hpleyz855hs8fvuvtq9d7m5",
			"reward": {
				"denom": "ucmdx",
				"amount": "1342821"
			}
		},
		{
			"address": "comdex1vpdehzac9n8lz22xhz09kstttf8kedp26wdl2h",
			"reward": {
				"denom": "ucmdx",
				"amount": "60408"
			}
		},
		{
			"address": "comdex1wsuj6e6n4qj5yrjr3q6kqw4aa5mkn2a2xel2g5",
			"reward": {
				"denom": "ucmdx",
				"amount": "29260749"
			}
		},
		{
			"address": "comdex10zejq4a7l3cnd0uy52587nyrhnnxkhd0jv6865",
			"reward": {
				"denom": "ucmdx",
				"amount": "925848"
			}
		},
		{
			"address": "comdex107vjqq2g0ksant3ga08s6am7j9pk7vwqcx77v2",
			"reward": {
				"denom": "ucmdx",
				"amount": "80907916"
			}
		},
		{
			"address": "comdex1swxpf95tts40mwarm8dq9t4r3z8w8rwgcwzte6",
			"reward": {
				"denom": "ucmdx",
				"amount": "4344999"
			}
		},
		{
			"address": "comdex13y20ud6h642cyp46ysp090djtggn9alh5y7n78",
			"reward": {
				"denom": "ucmdx",
				"amount": "665793"
			}
		},
		{
			"address": "comdex13yt4wyvj7aup24ftgxjjstk3p4um0ppx0dk6qd",
			"reward": {
				"denom": "ucmdx",
				"amount": "46058762"
			}
		},
		{
			"address": "comdex13j3d33etzccfdrkvkncva50hyd9myvd3wluyxe",
			"reward": {
				"denom": "ucmdx",
				"amount": "717729"
			}
		},
		{
			"address": "comdex1ng4q0sfg2yzftnl5mw9l3snh023zh97rmlg9hx",
			"reward": {
				"denom": "ucmdx",
				"amount": "5338880"
			}
		},
		{
			"address": "comdex1nahkrrrzfddr9rl6vxngua7rk4wt4w3lg0ksy0",
			"reward": {
				"denom": "ucmdx",
				"amount": "407584"
			}
		},
		{
			"address": "comdex15v963q2lshf9xyjv48mukuzfsx9snzuk7kqwj8",
			"reward": {
				"denom": "ucmdx",
				"amount": "30500786"
			}
		},
		{
			"address": "comdex1kqjxvmv7l9wk48a6uwugpas43pklt4umql9gc9",
			"reward": {
				"denom": "ucmdx",
				"amount": "8076525"
			}
		},
		{
			"address": "comdex1c9l62lktjnm5kc0xcc7gsendg6flcyczhykvmf",
			"reward": {
				"denom": "ucmdx",
				"amount": "4371824"
			}
		},
		{
			"address": "comdex1cv9jnq4qkcq3lk6szj2nxpwq6ew8dusa7pjh6n",
			"reward": {
				"denom": "ucmdx",
				"amount": "125701030"
			}
		},
		{
			"address": "comdex16t35tqn64t90kuqnuyr83tegqftxzze84653eu",
			"reward": {
				"denom": "ucmdx",
				"amount": "205116"
			}
		},
		{
			"address": "comdex1m33ataln9ecmkkp2cxk3wq0xz4cucvxz7kwcgx",
			"reward": {
				"denom": "ucmdx",
				"amount": "6449751"
			}
		},
		{
			"address": "comdex1mag869x9qf2nvlx4vwmcapnhjc2vpc8qyjyjwm",
			"reward": {
				"denom": "ucmdx",
				"amount": "20277879"
			}
		},
		{
			"address": "comdex1uwkvzf0smw6je8mgydfsex3pv60azczzpylxl5",
			"reward": {
				"denom": "ucmdx",
				"amount": "44251006"
			}
		},
		{
			"address": "comdex1u6v530yt3a089dkmt3zesytpwe2xne9cc7g8pg",
			"reward": {
				"denom": "ucmdx",
				"amount": "49046128"
			}
		},
		{
			"address": "comdex1aq4usrgpgre72zlfp4lzqjtqyrqhpmskwe4yza",
			"reward": {
				"denom": "ucmdx",
				"amount": "9587385"
			}
		},
		{
			"address": "comdex1a0acmzcq03z3pdxyakkckzd85thdvckmjde2qa",
			"reward": {
				"denom": "ucmdx",
				"amount": "1001287"
			}
		},
		{
			"address": "comdex1asg0632vdjfvg9r59uych3wv2lu0v0vsvssjv4",
			"reward": {
				"denom": "ucmdx",
				"amount": "17343001"
			}
		},
		{
			"address": "comdex173pn5d5lct5xkpx9l85ww8p97c8ssge3y2nw23",
			"reward": {
				"denom": "ucmdx",
				"amount": "130054"
			}
		},
		{
			"address": "comdex17jp3jml0txnraw2m9wjwwwnhldf47rw2635t4f",
			"reward": {
				"denom": "ucmdx",
				"amount": "27771"
			}
		},
		{
			"address": "comdex1l6wmcnkme82p8zztryautdp3xs09k2mfqn3uj3",
			"reward": {
				"denom": "ucmdx",
				"amount": "26969603"
			}
		},
		{
			"address": "comdex1qgu227rkd8qwq92mdr0l0cp4e5n6y4yrujvz70",
			"reward": {
				"denom": "ucmdx",
				"amount": "588018"
			}
		},
		{
			"address": "comdex1pl6c8zquzjvp0nwx6n0t0s2s3v9x2378e7j5t4",
			"reward": {
				"denom": "ucmdx",
				"amount": "65810764"
			}
		},
		{
			"address": "comdex1xpayr902xj2mzm3k7j3z6sqt0mqg9hr3rzaskn",
			"reward": {
				"denom": "ucmdx",
				"amount": "2055305"
			}
		},
		{
			"address": "comdex1gw32e9aec4vkju46ddkp5qv6xv6zrk33swyxg6",
			"reward": {
				"denom": "ucmdx",
				"amount": "114"
			}
		},
		{
			"address": "comdex1fj8d32l98legw3l23ag5gl6ulc64akfr9a3r97",
			"reward": {
				"denom": "ucmdx",
				"amount": "841006"
			}
		},
		{
			"address": "comdex1fmxg5edzddftgv5g8dmzycguyncq5dr95yreze",
			"reward": {
				"denom": "ucmdx",
				"amount": "18553"
			}
		},
		{
			"address": "comdex12q0708jnrd6d5ud7ap5lz4tgu3yshppfwd9x28",
			"reward": {
				"denom": "ucmdx",
				"amount": "37104062"
			}
		},
		{
			"address": "comdex12wdjtsdgnyek2jag90fnzmctt3hf2nmf869k98",
			"reward": {
				"denom": "ucmdx",
				"amount": "7383463"
			}
		},
		{
			"address": "comdex1w4cjautlrggjfqq72ca4edsmr0m2d9wex524sa",
			"reward": {
				"denom": "ucmdx",
				"amount": "1237628"
			}
		},
		{
			"address": "comdex106h5juhmegsnr4sr5ukpg23ha3z7m4z6v044n0",
			"reward": {
				"denom": "ucmdx",
				"amount": "29912600"
			}
		},
		{
			"address": "comdex15r63trllw4jlsje55u0lvfd0tf2vj5pn69uqh5",
			"reward": {
				"denom": "ucmdx",
				"amount": "8823368"
			}
		},
		{
			"address": "comdex14p3263d5sppcxkq90vp5aq5mta3njh6x0pp6jr",
			"reward": {
				"denom": "ucmdx",
				"amount": "5270965"
			}
		},
		{
			"address": "comdex1kkaryts8yf99xjdytec467lgm8zjekx0ghvjdc",
			"reward": {
				"denom": "ucmdx",
				"amount": "1460160"
			}
		},
		{
			"address": "comdex1keyf5w20uwzua3f44ekues05ru2zpjj282xnnp",
			"reward": {
				"denom": "ucmdx",
				"amount": "12916896"
			}
		},
		{
			"address": "comdex1hf3ljlcrvvp53hm8aff3wgkwgmg2lgxsj5rrqw",
			"reward": {
				"denom": "ucmdx",
				"amount": "6025858"
			}
		},
		{
			"address": "comdex1h35xrvws8e7yxecncv0qgvyquhcfemz8d6g7aj",
			"reward": {
				"denom": "ucmdx",
				"amount": "332244"
			}
		},
		{
			"address": "comdex1chlmc3um9qrs66g3afvcpwtexfnhdf3a8v4384",
			"reward": {
				"denom": "ucmdx",
				"amount": "36573321"
			}
		},
		{
			"address": "comdex1u8xff6p6p4uunrpha6m7fskv6j5z49qhtdddyp",
			"reward": {
				"denom": "ucmdx",
				"amount": "3835544"
			}
		},
		{
			"address": "comdex1ut8s48pfm5nnmjt3dsnssmrqwktdkkdazwftne",
			"reward": {
				"denom": "ucmdx",
				"amount": "37071"
			}
		},
		{
			"address": "comdex17ydjt24h65tq5ylva3as0r4vvk0a2wjgcye4g5",
			"reward": {
				"denom": "ucmdx",
				"amount": "939"
			}
		},
		{
			"address": "comdex178gfvn2es4fuv9g84zxtwrf2dfhx57h79cv7a2",
			"reward": {
				"denom": "ucmdx",
				"amount": "1869747"
			}
		},
		{
			"address": "comdex1q20t5l6pz3wzjznwq2m4wd827yhr2nrqmvqgvw",
			"reward": {
				"denom": "ucmdx",
				"amount": "30223"
			}
		},
		{
			"address": "comdex1qhf0we0fsyzczljvyct59ajzy745afdwr9fzwz",
			"reward": {
				"denom": "ucmdx",
				"amount": "630"
			}
		},
		{
			"address": "comdex1pwpz0acvw0mc0clr4kknedt94efhwzj8w0ljvn",
			"reward": {
				"denom": "ucmdx",
				"amount": "432250"
			}
		},
		{
			"address": "comdex1psuq9r88vl47ryrwprpkz2f3wj2la8g5rxhvua",
			"reward": {
				"denom": "ucmdx",
				"amount": "117461"
			}
		},
		{
			"address": "comdex1pa9afmmfwhy66sauda5esxrelgue625da2rfhl",
			"reward": {
				"denom": "ucmdx",
				"amount": "129527"
			}
		},
		{
			"address": "comdex1zq8tm7kj4cjjx6h3eyla23qrvq50kekannx05j",
			"reward": {
				"denom": "ucmdx",
				"amount": "14421"
			}
		},
		{
			"address": "comdex1z0p2u7k4rxuqrj9hn2qetmrykz6yat369pr2mt",
			"reward": {
				"denom": "ucmdx",
				"amount": "172905"
			}
		},
		{
			"address": "comdex1z4dc7c3acgslhhy98eut7022t30r2n4r3xms5q",
			"reward": {
				"denom": "ucmdx",
				"amount": "92382"
			}
		},
		{
			"address": "comdex1rqnp7dtmyh4jg55zkk8fthphmyywy5aygh3u5d",
			"reward": {
				"denom": "ucmdx",
				"amount": "2060461"
			}
		},
		{
			"address": "comdex1r9l5e4exvgcxcaxsx6nu90zyvv85y9de35v875",
			"reward": {
				"denom": "ucmdx",
				"amount": "17432"
			}
		},
		{
			"address": "comdex1rjazewzf3ecvlgnx79kw54upg3zxlugr4zlegv",
			"reward": {
				"denom": "ucmdx",
				"amount": "67616079"
			}
		},
		{
			"address": "comdex1r6h7w64333j4zderxk8huumpu2helrkf08tchm",
			"reward": {
				"denom": "ucmdx",
				"amount": "190084"
			}
		},
		{
			"address": "comdex1y7g9eumdeeajzlq2lqcznq9jtgw79jnmkrchn2",
			"reward": {
				"denom": "ucmdx",
				"amount": "2678"
			}
		},
		{
			"address": "comdex198uj3ffetplr9szd774wznk2ssjnjkghw6llx2",
			"reward": {
				"denom": "ucmdx",
				"amount": "1332184"
			}
		},
		{
			"address": "comdex1xfuras3838c8ln3xkmzs7f5kh5m5c6g9mhfkc2",
			"reward": {
				"denom": "ucmdx",
				"amount": "314877"
			}
		},
		{
			"address": "comdex1xvkrdmtj0qw6htydtlzj88qvpqk8h96mdstc00",
			"reward": {
				"denom": "ucmdx",
				"amount": "123118"
			}
		},
		{
			"address": "comdex1x0wa3ffgvz7ypwxw3vn0j9xsdl44f4d6y9gcmv",
			"reward": {
				"denom": "ucmdx",
				"amount": "70730558"
			}
		},
		{
			"address": "comdex1x58ve32e0puzuxd2nyf5e3hp8zz90nfaj08w24",
			"reward": {
				"denom": "ucmdx",
				"amount": "6646428"
			}
		},
		{
			"address": "comdex1x7u0njqav7x444uggrwphmfdkk50ru73lhxuww",
			"reward": {
				"denom": "ucmdx",
				"amount": "114009"
			}
		},
		{
			"address": "comdex189w77d3482h3v7d75yyj970924maak6cemd9cu",
			"reward": {
				"denom": "ucmdx",
				"amount": "831199"
			}
		},
		{
			"address": "comdex182hdh9w3wj77uzjsljdwjktpqwljf4rh5cwa0h",
			"reward": {
				"denom": "ucmdx",
				"amount": "43349"
			}
		},
		{
			"address": "comdex18djaz26cw6skc7x3kefg9c56gnqw4twpu5z2pe",
			"reward": {
				"denom": "ucmdx",
				"amount": "58220"
			}
		},
		{
			"address": "comdex1g6a8mqm5sc3q827vw9nte25txqu6uykjnlexym",
			"reward": {
				"denom": "ucmdx",
				"amount": "60064"
			}
		},
		{
			"address": "comdex1fwcnug7tph0lm7jem86372d3a8q26496u4r0w2",
			"reward": {
				"denom": "ucmdx",
				"amount": "19453948"
			}
		},
		{
			"address": "comdex1f3zda5nvg2a7jn8c8ehxsx8xwh9tzr64t522h0",
			"reward": {
				"denom": "ucmdx",
				"amount": "10385"
			}
		},
		{
			"address": "comdex1f3ase4a9lq2tfe24mhzqxwjtsqc67eypu7zayn",
			"reward": {
				"denom": "ucmdx",
				"amount": "36445255"
			}
		},
		{
			"address": "comdex1fjg9nz9f58a836m6plp0dj2t56as58aluvcsx2",
			"reward": {
				"denom": "ucmdx",
				"amount": "3213576"
			}
		},
		{
			"address": "comdex12gz0mh4zz4dzuhs9m76z46a5mryngwh37769xt",
			"reward": {
				"denom": "ucmdx",
				"amount": "3045"
			}
		},
		{
			"address": "comdex122hjrmj8k2f0lc0m8nepdysgx9f2gan2a8kz7e",
			"reward": {
				"denom": "ucmdx",
				"amount": "4849"
			}
		},
		{
			"address": "comdex12e4cr5w807csadkfytk3fep8kdhcv6x2vpmg7y",
			"reward": {
				"denom": "ucmdx",
				"amount": "2577"
			}
		},
		{
			"address": "comdex1t8r93jrvfu04758dtqq8yqrhcw940edxq3py9y",
			"reward": {
				"denom": "ucmdx",
				"amount": "9298"
			}
		},
		{
			"address": "comdex1vkthfcdhnn4uvz83q29vwkx2qpzjmpuemqzwl9",
			"reward": {
				"denom": "ucmdx",
				"amount": "38551"
			}
		},
		{
			"address": "comdex1va6r62ul3yqqc5xt0v59np8mlr7m5pxe0erx45",
			"reward": {
				"denom": "ucmdx",
				"amount": "17555"
			}
		},
		{
			"address": "comdex1drw8jjhdsrsgwvvkccnnva635m84veeffsc07w",
			"reward": {
				"denom": "ucmdx",
				"amount": "2050247"
			}
		},
		{
			"address": "comdex1djdzvwwjdvhqrjl5ww2cpnmj23th7sjp8fmmmp",
			"reward": {
				"denom": "ucmdx",
				"amount": "601996"
			}
		},
		{
			"address": "comdex1d5e3fppc5xaaql2n9dgeppppjzgsnkha3nn8ep",
			"reward": {
				"denom": "ucmdx",
				"amount": "21015"
			}
		},
		{
			"address": "comdex1dk29mferv4nuk90pn7xca5y24t5m87d2v257qa",
			"reward": {
				"denom": "ucmdx",
				"amount": "354605"
			}
		},
		{
			"address": "comdex1wn0schdngyc8qudg8n4757zytlp5f39l6elcc2",
			"reward": {
				"denom": "ucmdx",
				"amount": "334487"
			}
		},
		{
			"address": "comdex109fhpxjztun2s6plnyzxpsmctcz9cdj4cc7p0w",
			"reward": {
				"denom": "ucmdx",
				"amount": "271143"
			}
		},
		{
			"address": "comdex103nkz0hgnpwmjveafmr57xwmvlluuj4s5qzds2",
			"reward": {
				"denom": "ucmdx",
				"amount": "1254713"
			}
		},
		{
			"address": "comdex10kcns66kdspnywaslcphcknyvlrr8r8da2aal3",
			"reward": {
				"denom": "ucmdx",
				"amount": "606625"
			}
		},
		{
			"address": "comdex10c3cjjvkk6s5m8qzpm8stjce98sw7tsftkak2e",
			"reward": {
				"denom": "ucmdx",
				"amount": "353323"
			}
		},
		{
			"address": "comdex10u24rwqfxl48rshaasextp5a25gejuzarj3ptc",
			"reward": {
				"denom": "ucmdx",
				"amount": "49232"
			}
		},
		{
			"address": "comdex1slg5rjchd5ls2e3avmyvvyw4qj69rkt0u3g4ln",
			"reward": {
				"denom": "ucmdx",
				"amount": "212477"
			}
		},
		{
			"address": "comdex13a6nz8mej9e9gur525x6d2hfw0ql027nq90tk5",
			"reward": {
				"denom": "ucmdx",
				"amount": "1891793"
			}
		},
		{
			"address": "comdex137nm0r0yycejyaajxmswyh9r459xw6gtpszlwe",
			"reward": {
				"denom": "ucmdx",
				"amount": "279051"
			}
		},
		{
			"address": "comdex1jyehgmj4fr2r2us2fq9dj50xretqk70h6yfd9n",
			"reward": {
				"denom": "ucmdx",
				"amount": "22990975"
			}
		},
		{
			"address": "comdex1jcjgggrjk03clt9y3cy94kdj4nynvdufa0vv0h",
			"reward": {
				"denom": "ucmdx",
				"amount": "3543"
			}
		},
		{
			"address": "comdex1j69vyad5kr0zpppluty24ephhw7s82nk6k73fr",
			"reward": {
				"denom": "ucmdx",
				"amount": "2991"
			}
		},
		{
			"address": "comdex1jllhnmv2v9q68fyttfxhcd5xahd4h4y5p5f8fj",
			"reward": {
				"denom": "ucmdx",
				"amount": "74456"
			}
		},
		{
			"address": "comdex1nfu9j53v3dd650usnnatea8vywz5sz84rz7f49",
			"reward": {
				"denom": "ucmdx",
				"amount": "64033"
			}
		},
		{
			"address": "comdex15qqeu3mj4y4f0wu5758ym592n4aghynx8ah8h7",
			"reward": {
				"denom": "ucmdx",
				"amount": "196943"
			}
		},
		{
			"address": "comdex158er92qej8hj7aw7n08yf0e0z8lun7cymgm7s9",
			"reward": {
				"denom": "ucmdx",
				"amount": "10138"
			}
		},
		{
			"address": "comdex15j7g66x5y99d9kjs0zvkvejgyhpgfaelzaq2d9",
			"reward": {
				"denom": "ucmdx",
				"amount": "1982218"
			}
		},
		{
			"address": "comdex15es5dhgm0yf97kcuwxuwf9vyvpy0hddwt53uw0",
			"reward": {
				"denom": "ucmdx",
				"amount": "1451978"
			}
		},
		{
			"address": "comdex1kxg02c97kzjmkjwl9r00jduv62rwsfywsuyuju",
			"reward": {
				"denom": "ucmdx",
				"amount": "331788"
			}
		},
		{
			"address": "comdex1kf5uk0quus8t34ys7p3judrxdqjayl0s2h5xuy",
			"reward": {
				"denom": "ucmdx",
				"amount": "352223"
			}
		},
		{
			"address": "comdex1hw9s8lm9v29sm3t6mw5r5cqrux5cs9usl7nxxa",
			"reward": {
				"denom": "ucmdx",
				"amount": "11571"
			}
		},
		{
			"address": "comdex1hs9rq5r4sn5dlpqerxq7zpu32q2eld8xlcp7wq",
			"reward": {
				"denom": "ucmdx",
				"amount": "29766"
			}
		},
		{
			"address": "comdex1h3mr4kkqy05v9v7a09y2zc80np0t4637v537a7",
			"reward": {
				"denom": "ucmdx",
				"amount": "12270"
			}
		},
		{
			"address": "comdex1cxyjsjj7qjkh9gearccn563hwujrqdewdw4j9j",
			"reward": {
				"denom": "ucmdx",
				"amount": "113766"
			}
		},
		{
			"address": "comdex1cmk6j07uhgcv83g46n9d6p04vmj79p9n72rsq4",
			"reward": {
				"denom": "ucmdx",
				"amount": "97011"
			}
		},
		{
			"address": "comdex1cljvfex2cvth9t7l54meast0lvtudp8xq5p9dz",
			"reward": {
				"denom": "ucmdx",
				"amount": "225596"
			}
		},
		{
			"address": "comdex1ejunsa3khkhz7k6wc7hezxn4umfqnm4cphza27",
			"reward": {
				"denom": "ucmdx",
				"amount": "602017"
			}
		},
		{
			"address": "comdex1e6zdjcj3wry34hns4pw9hvse7xcv9mmtp0sakz",
			"reward": {
				"denom": "ucmdx",
				"amount": "9569890"
			}
		},
		{
			"address": "comdex16phjecw25muxjyxgft4jxmwntvx24xhc5rqdca",
			"reward": {
				"denom": "ucmdx",
				"amount": "1094700019"
			}
		},
		{
			"address": "comdex16937wn3kqensn3gkd7tggaww4m5apsjczzcgyl",
			"reward": {
				"denom": "ucmdx",
				"amount": "8685"
			}
		},
		{
			"address": "comdex1m9yhy6vnxsle3kufr2n2v58u8rnmlpe0pfme3w",
			"reward": {
				"denom": "ucmdx",
				"amount": "903100"
			}
		},
		{
			"address": "comdex1m8h2psgm0jtpcaw89hscwhf6w3h07r2dq4vxjg",
			"reward": {
				"denom": "ucmdx",
				"amount": "63611209"
			}
		},
		{
			"address": "comdex1mdsku4csjmh20u4rfq0c2psqhd5p7ftuqgcx46",
			"reward": {
				"denom": "ucmdx",
				"amount": "23543666"
			}
		},
		{
			"address": "comdex1mjt7j5a34yn2085qx5vdg5fgelu4s973hkcj76",
			"reward": {
				"denom": "ucmdx",
				"amount": "456328"
			}
		},
		{
			"address": "comdex1u9hpqcz90f9q8eqjgp70yca0k4nmeyz2t4z20u",
			"reward": {
				"denom": "ucmdx",
				"amount": "7391319"
			}
		},
		{
			"address": "comdex1ug4p9pyeqnfd2kvtvgl5ulu2h7k428syflwp62",
			"reward": {
				"denom": "ucmdx",
				"amount": "90630219"
			}
		},
		{
			"address": "comdex1u0ljqruvle6rhc3elx309j7e4ukpw506778ypj",
			"reward": {
				"denom": "ucmdx",
				"amount": "19123"
			}
		},
		{
			"address": "comdex1apyptpdnapp2wqw5edutp3lfrj67mulpewq9jf",
			"reward": {
				"denom": "ucmdx",
				"amount": "118169"
			}
		},
		{
			"address": "comdex1afumn0gvtgdafvq4fd9vs78n6dk295qzxnz37e",
			"reward": {
				"denom": "ucmdx",
				"amount": "25215"
			}
		},
		{
			"address": "comdex1a0895y9u9955cegxgdsfgrsfvrgz2ssmcnvtmm",
			"reward": {
				"denom": "ucmdx",
				"amount": "6902623"
			}
		},
		{
			"address": "comdex1a0a35fn6xjztsnz6ttv8n7rj6f0qvr8ecc95f6",
			"reward": {
				"denom": "ucmdx",
				"amount": "675"
			}
		},
		{
			"address": "comdex1asrupfzjyepcjhafz2fjcz63wtslj4ue45fuup",
			"reward": {
				"denom": "ucmdx",
				"amount": "2366121"
			}
		},
		{
			"address": "comdex1ajv29pm0dy0ftrupyvyf2g8ru2yapyljzr6fjp",
			"reward": {
				"denom": "ucmdx",
				"amount": "67931"
			}
		},
		{
			"address": "comdex17ylsw59mw9d06fhlgaf30fcyrcxjc7nq4dsff2",
			"reward": {
				"denom": "ucmdx",
				"amount": "64934"
			}
		},
		{
			"address": "comdex17al5rus4yfr2s0pd8cwy300tvugz90wxk4aryg",
			"reward": {
				"denom": "ucmdx",
				"amount": "59799"
			}
		},
		{
			"address": "comdex1lsevtcws2cg4jhqn4f9uzyvgjyzz0rkl80mk7d",
			"reward": {
				"denom": "ucmdx",
				"amount": "1111502"
			}
		},
		{
			"address": "comdex1lamfr7ufuflu6l7aukgjnrqzcg4ad2776654nf",
			"reward": {
				"denom": "ucmdx",
				"amount": "302771"
			}
		},
		{
			"address": "comdex1z0ld6q0usehh5y4vcks25kay7ruqaxmy4tcvy7",
			"reward": {
				"denom": "ucmdx",
				"amount": "1472672"
			}
		},
		{
			"address": "comdex1rn4jsdp2hr8zgjemn23uf5u3qfwkes2pfy38g4",
			"reward": {
				"denom": "ucmdx",
				"amount": "9328090"
			}
		},
		{
			"address": "comdex195ngs6h55u2pakjrxuk0nmcfdlm968vn82mck2",
			"reward": {
				"denom": "ucmdx",
				"amount": "46257"
			}
		},
		{
			"address": "comdex18xddcf3ymdhq8nkcduqtt8x5t2zpjg7waxhl7t",
			"reward": {
				"denom": "ucmdx",
				"amount": "23616799"
			}
		},
		{
			"address": "comdex1garqed6pawqcadhcecq8gfghn7cqssdvm88kvd",
			"reward": {
				"denom": "ucmdx",
				"amount": "941326"
			}
		},
		{
			"address": "comdex1204g59lju5lysugsc9d5d686jnnrhavdlyql49",
			"reward": {
				"denom": "ucmdx",
				"amount": "2793218"
			}
		},
		{
			"address": "comdex126yssj4fpkngvcxrdspq6j042d7098hm2c87tk",
			"reward": {
				"denom": "ucmdx",
				"amount": "4520836"
			}
		},
		{
			"address": "comdex1t9m89l744kaasaq7w9j3agmwjejnjrrcaxtm2x",
			"reward": {
				"denom": "ucmdx",
				"amount": "65378170"
			}
		},
		{
			"address": "comdex1v89s8e6rd8h0sm5h6k2qd4q8m2yehaughf0cr3",
			"reward": {
				"denom": "ucmdx",
				"amount": "5740996"
			}
		},
		{
			"address": "comdex1v6aj98yxlr3sn8g4xy9l7s74cy8hd8eyj33m09",
			"reward": {
				"denom": "ucmdx",
				"amount": "1616546"
			}
		},
		{
			"address": "comdex1d5lmyv7ktrpjmu94ka0zteuhpef6a0p0vpeq4r",
			"reward": {
				"denom": "ucmdx",
				"amount": "47685"
			}
		},
		{
			"address": "comdex1sefn0da30juv928gj4anarq27vqq73v499m9j8",
			"reward": {
				"denom": "ucmdx",
				"amount": "2739281"
			}
		},
		{
			"address": "comdex14kf23cs9evs66uvl63ef6ggp20apqdc3ykfgwy",
			"reward": {
				"denom": "ucmdx",
				"amount": "14669"
			}
		},
		{
			"address": "comdex1cq27hkg8h6u6fkz2vh0xjqvktal3d23splzgf3",
			"reward": {
				"denom": "ucmdx",
				"amount": "1846962"
			}
		},
		{
			"address": "comdex17nun6cg7a5xtu9fz2p8scxr7u9jmrf5vtmpuwm",
			"reward": {
				"denom": "ucmdx",
				"amount": "120091"
			}
		},
		{
			"address": "comdex1ljuc5cguh2ehalzp8s53zrt6g0u5jmjk4fn8th",
			"reward": {
				"denom": "ucmdx",
				"amount": "78767"
			}
		},
		{
			"address": "comdex1pyeq6crs7g4357wlk0gassflplq3j99jp7j6s8",
			"reward": {
				"denom": "ucmdx",
				"amount": "343761"
			}
		},
		{
			"address": "comdex1xksfw6mjzanaa6tdwpshkpn97mv7g3mn6dm20u",
			"reward": {
				"denom": "ucmdx",
				"amount": "186542659"
			}
		},
		{
			"address": "comdex1gn2saq5vy2qjn9ktn8rl9cx9ajvn8drccvarm3",
			"reward": {
				"denom": "ucmdx",
				"amount": "1405"
			}
		},
		{
			"address": "comdex1jdcvuhd83yv026sa5frheexyf7zhmpd5xlcwkw",
			"reward": {
				"denom": "ucmdx",
				"amount": "137074327"
			}
		},
		{
			"address": "comdex14u2mtgtqw0yuge8f7g89ruunpafhhrgng3v4dw",
			"reward": {
				"denom": "ucmdx",
				"amount": "17980054"
			}
		},
		{
			"address": "comdex1ezfw4qu70rwkg0hx2dh9whyjre5ldah0qaqesh",
			"reward": {
				"denom": "ucmdx",
				"amount": "54693290"
			}
		},
		{
			"address": "comdex1p77sfv0aefkxst0uuc5vvavuvp08jw3gaas5kp",
			"reward": {
				"denom": "ucmdx",
				"amount": "362418"
			}
		},
		{
			"address": "comdex1yjycpl472f6z4sp08eyj0s5fynp4k6r3lejafc",
			"reward": {
				"denom": "ucmdx",
				"amount": "553"
			}
		},
		{
			"address": "comdex1wgec85sy32jgjwmgfrhzna492nqd9qfk07v0v3",
			"reward": {
				"denom": "ucmdx",
				"amount": "46241748"
			}
		},
		{
			"address": "comdex15zv0g0y4652z8vhj0kvx325rea7r2n67j839l9",
			"reward": {
				"denom": "ucmdx",
				"amount": "5619596"
			}
		},
		{
			"address": "comdex14csg6zvxelnxsu3wysxcecwds9dqv3j07l6jwu",
			"reward": {
				"denom": "ucmdx",
				"amount": "100142942"
			}
		},
		{
			"address": "comdex1cq8dwur2xk342yyze54yuckyajyfm42la770uk",
			"reward": {
				"denom": "ucmdx",
				"amount": "13105"
			}
		},
		{
			"address": "comdex1mjvyusevcmqtf64qja7djnquls674w2vuu550d",
			"reward": {
				"denom": "ucmdx",
				"amount": "26436376"
			}
		},
		{
			"address": "comdex1qqm6nam6hse8qpn8y8sy5jwpnmkcpdpmqjakmx",
			"reward": {
				"denom": "ucmdx",
				"amount": "4827594"
			}
		},
		{
			"address": "comdex1qzgl0cmnhwrn9n4874zdgc86czxdepg8y0tpyk",
			"reward": {
				"denom": "ucmdx",
				"amount": "3203266"
			}
		},
		{
			"address": "comdex1qzu34xw9hh5jndqsex73gydjnr5ck6gjdewthn",
			"reward": {
				"denom": "ucmdx",
				"amount": "8706677"
			}
		},
		{
			"address": "comdex1qxh8dfj72vymqnqamatdn5e8zmun5ggxjl0mjt",
			"reward": {
				"denom": "ucmdx",
				"amount": "6105827"
			}
		},
		{
			"address": "comdex1q8rknj62swgv7fj8v0d3ncz57gq60hwv7e6ydt",
			"reward": {
				"denom": "ucmdx",
				"amount": "16093"
			}
		},
		{
			"address": "comdex1qgaxjh75wsy284747zlqr2kdjqrtmwg62dd5p0",
			"reward": {
				"denom": "ucmdx",
				"amount": "121765"
			}
		},
		{
			"address": "comdex1qnqqj4hwzfjsma2ut0j0rqm823tefj3amq3l8g",
			"reward": {
				"denom": "ucmdx",
				"amount": "1795"
			}
		},
		{
			"address": "comdex1qcf56hq8x4dwa380tk0hfmmfdswhucerf4c0nf",
			"reward": {
				"denom": "ucmdx",
				"amount": "287680"
			}
		},
		{
			"address": "comdex1qamfln8u5w8d3vlhp5t9mhmylfkgad4jn828x8",
			"reward": {
				"denom": "ucmdx",
				"amount": "9548"
			}
		},
		{
			"address": "comdex1pxc063fv3795vygudz6k0uf0lew4zv68f4wauh",
			"reward": {
				"denom": "ucmdx",
				"amount": "12266"
			}
		},
		{
			"address": "comdex1pvgdhrlfc5eqkkzlptjrs6u35zd0qcteq2jj7u",
			"reward": {
				"denom": "ucmdx",
				"amount": "14759"
			}
		},
		{
			"address": "comdex1p0thq5p7r0qnn3fgn0eta7wnwypk55frl6crd5",
			"reward": {
				"denom": "ucmdx",
				"amount": "2"
			}
		},
		{
			"address": "comdex1p0s7kwr895wa8780aauj7ksh0guhxfmv3v3cyd",
			"reward": {
				"denom": "ucmdx",
				"amount": "20009"
			}
		},
		{
			"address": "comdex1p3punn2xl9pspg79vwtam7zdgw2e3ua2r0cnhx",
			"reward": {
				"denom": "ucmdx",
				"amount": "21"
			}
		},
		{
			"address": "comdex1pnsplaqssgv9uz5p6zvmg6vj8a3vxd796cl2fy",
			"reward": {
				"denom": "ucmdx",
				"amount": "50523"
			}
		},
		{
			"address": "comdex1p5p0f2nn27kyz6ssfwvqg93jufwr890trrpl9c",
			"reward": {
				"denom": "ucmdx",
				"amount": "23849"
			}
		},
		{
			"address": "comdex1p5v8hrard43g66rp92cwmr8hvvrmn76xjdh923",
			"reward": {
				"denom": "ucmdx",
				"amount": "47576"
			}
		},
		{
			"address": "comdex1pehx7rd5rg6x0q54t9jredfjuue70ehvp09llt",
			"reward": {
				"denom": "ucmdx",
				"amount": "15546"
			}
		},
		{
			"address": "comdex1p6fef9dy2gvr0zpg2nlz3f8d08zkqef9wv9l8g",
			"reward": {
				"denom": "ucmdx",
				"amount": "93503"
			}
		},
		{
			"address": "comdex1pl7kt82s76zdkqcdnk2rzcfalyy9m6w2kxpda0",
			"reward": {
				"denom": "ucmdx",
				"amount": "10688"
			}
		},
		{
			"address": "comdex1zzuttewnk3keg6m3xf2n3wxd2h7wkvj7255mrm",
			"reward": {
				"denom": "ucmdx",
				"amount": "150766"
			}
		},
		{
			"address": "comdex1zfqwy2d93x7krkn8k04ahr5knd49aszyx3fkw8",
			"reward": {
				"denom": "ucmdx",
				"amount": "7905"
			}
		},
		{
			"address": "comdex1zd2j8zjnnjk86funkde8gh6szpcrt5h9twfspw",
			"reward": {
				"denom": "ucmdx",
				"amount": "602"
			}
		},
		{
			"address": "comdex1zwds5s2mkhgl3gev5p9avzrqlw0l4lfu2v48h4",
			"reward": {
				"denom": "ucmdx",
				"amount": "97494"
			}
		},
		{
			"address": "comdex1z024vgu2v80dv0gg64aawnz8kynh3apu0jphap",
			"reward": {
				"denom": "ucmdx",
				"amount": "2029"
			}
		},
		{
			"address": "comdex1znpy8xveyszp7lfmpan9v650kvk25qmw3qctty",
			"reward": {
				"denom": "ucmdx",
				"amount": "1700871"
			}
		},
		{
			"address": "comdex1z4xjq43r7r8jjs46lr9uqeax57hu7g4lh6x2yv",
			"reward": {
				"denom": "ucmdx",
				"amount": "177947"
			}
		},
		{
			"address": "comdex1zmv7m80wztjl7xhf8lwxsz680e6wmuj9nvlwzp",
			"reward": {
				"denom": "ucmdx",
				"amount": "32155674"
			}
		},
		{
			"address": "comdex1zau8qapuptcl7gu824jnr475j95scxn6uwhjus",
			"reward": {
				"denom": "ucmdx",
				"amount": "320549"
			}
		},
		{
			"address": "comdex1z7h2qxlzjmq92x3cmynk8rr9r9gesyz6rflm7a",
			"reward": {
				"denom": "ucmdx",
				"amount": "101"
			}
		},
		{
			"address": "comdex1rrjx6jkk0jhrwazchzeuycdvs43hvdnmuv5mkq",
			"reward": {
				"denom": "ucmdx",
				"amount": "1563188"
			}
		},
		{
			"address": "comdex1ryw2duq6jqv4tl2t48w8qg2ecth45tw5hp0na2",
			"reward": {
				"denom": "ucmdx",
				"amount": "1920"
			}
		},
		{
			"address": "comdex1r88lnegkjtfyd699kh8vk73yqwv8d26du9est8",
			"reward": {
				"denom": "ucmdx",
				"amount": "19"
			}
		},
		{
			"address": "comdex1rftzn4676f5x2zj05hy4m0ew740qejzgde52m5",
			"reward": {
				"denom": "ucmdx",
				"amount": "66099"
			}
		},
		{
			"address": "comdex1rfd08kn78a8eyplg248y77pz9ppf3rrsuznx4e",
			"reward": {
				"denom": "ucmdx",
				"amount": "75717"
			}
		},
		{
			"address": "comdex1r2t22p9vuuhhcauzaaqwe8flv52mlegkvwn6sa",
			"reward": {
				"denom": "ucmdx",
				"amount": "1578768"
			}
		},
		{
			"address": "comdex1r2kht9xt8dv7lu6p2rsm3p5mjfnz6v5jfnd40z",
			"reward": {
				"denom": "ucmdx",
				"amount": "1983639"
			}
		},
		{
			"address": "comdex1rtjqaf4m587n768kve9h3rmkdupghcpx4k9af8",
			"reward": {
				"denom": "ucmdx",
				"amount": "9112"
			}
		},
		{
			"address": "comdex1rk3rt28nayujfnc69rzwx07z9jc8jczu6hk430",
			"reward": {
				"denom": "ucmdx",
				"amount": "6207630"
			}
		},
		{
			"address": "comdex1rmk0ztuhd7txscp7jr2srfjn5vvsdw252xnr4l",
			"reward": {
				"denom": "ucmdx",
				"amount": "1142054"
			}
		},
		{
			"address": "comdex1ypnke0r4uk6u82w4gh73kc5tz0qsn0ah05uhfj",
			"reward": {
				"denom": "ucmdx",
				"amount": "8946539"
			}
		},
		{
			"address": "comdex1y9p49mzleflp2w4z79q2cquhqsjt5w0se7sqhx",
			"reward": {
				"denom": "ucmdx",
				"amount": "39624"
			}
		},
		{
			"address": "comdex1ygx70d7u0waq2g4wrjcp0ce5teg0gwfn5mjfr4",
			"reward": {
				"denom": "ucmdx",
				"amount": "17250312"
			}
		},
		{
			"address": "comdex1yg4ghtvvpgqasy6hggqc58w0yatslktwcvuwa6",
			"reward": {
				"denom": "ucmdx",
				"amount": "58272"
			}
		},
		{
			"address": "comdex1ytagpftcttkjt9y5g0e5e4eln03kjaglanwprp",
			"reward": {
				"denom": "ucmdx",
				"amount": "1248388"
			}
		},
		{
			"address": "comdex1yvk507d2cthqmxvuw4408reuyzd54m3n2tmqtr",
			"reward": {
				"denom": "ucmdx",
				"amount": "22651"
			}
		},
		{
			"address": "comdex1yd47xy28cv05e70fd6d706fsmtkuqu86xmr342",
			"reward": {
				"denom": "ucmdx",
				"amount": "3625603"
			}
		},
		{
			"address": "comdex1y35ygj62m2cr2384tm52nnfqwlvmnh97ul80p6",
			"reward": {
				"denom": "ucmdx",
				"amount": "6201"
			}
		},
		{
			"address": "comdex1yjq7cnjzq4k4498afreeuaxvntvwr830yj6uyc",
			"reward": {
				"denom": "ucmdx",
				"amount": "24067"
			}
		},
		{
			"address": "comdex1ynulqdg0l27f5hhtdwejtrx0p6z7qxkcrjp7rc",
			"reward": {
				"denom": "ucmdx",
				"amount": "178448"
			}
		},
		{
			"address": "comdex1ykyvwkzydq7u05q0yjzue7v7xva3wydq44exnj",
			"reward": {
				"denom": "ucmdx",
				"amount": "57568"
			}
		},
		{
			"address": "comdex1ycu42lxdp4kyyyv74x3pppg666qqgldhnrvakt",
			"reward": {
				"denom": "ucmdx",
				"amount": "17686"
			}
		},
		{
			"address": "comdex1yejx7ke7kh3gudqahk5ky2lyf8u6494wj9s4m3",
			"reward": {
				"denom": "ucmdx",
				"amount": "51"
			}
		},
		{
			"address": "comdex1y74dcm0lwvcgva2qfyw3xdafl852aznge4gs3a",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex19zrtt2uncccycgcm9efgedu0hr27qfdvnqhqgk",
			"reward": {
				"denom": "ucmdx",
				"amount": "358141"
			}
		},
		{
			"address": "comdex1999nzcaapfrzqddj8htlucmfjez43ywtw3dq3q",
			"reward": {
				"denom": "ucmdx",
				"amount": "677"
			}
		},
		{
			"address": "comdex198mp2nyzwshjrytzw6g6tuyulkvqngtw8jyctu",
			"reward": {
				"denom": "ucmdx",
				"amount": "206728"
			}
		},
		{
			"address": "comdex19gqzgxe99gq6du2va5p0ahszklrumpkfm3n3l0",
			"reward": {
				"denom": "ucmdx",
				"amount": "3128542"
			}
		},
		{
			"address": "comdex19fgws22u6m94v2rnk09nvqxfgaga08eyzhknhr",
			"reward": {
				"denom": "ucmdx",
				"amount": "52"
			}
		},
		{
			"address": "comdex190crhcu09fp4492a0wua59hfl97pxnjvdf7ry2",
			"reward": {
				"denom": "ucmdx",
				"amount": "2259831"
			}
		},
		{
			"address": "comdex195khh6nr3zsvaxkm53wavkxv46uucdlayed5m7",
			"reward": {
				"denom": "ucmdx",
				"amount": "79242"
			}
		},
		{
			"address": "comdex19kd6f0yn50qh80xw54jtr7aeesch0ksh47xz8y",
			"reward": {
				"denom": "ucmdx",
				"amount": "1108383"
			}
		},
		{
			"address": "comdex19e2y7w5s7qvldsse0rznh4lh0pvyxhdqh6ad5l",
			"reward": {
				"denom": "ucmdx",
				"amount": "34"
			}
		},
		{
			"address": "comdex19eltgsnnzajpyjkmywd3uvraa4p4ysx0zfel85",
			"reward": {
				"denom": "ucmdx",
				"amount": "727"
			}
		},
		{
			"address": "comdex19umd68kd2w6y9fgh8n4vclm75d8s9fvkuystt3",
			"reward": {
				"denom": "ucmdx",
				"amount": "1384377"
			}
		},
		{
			"address": "comdex1xzvt0t944uervef24y5x9tfy0kj2fe2m79xjjj",
			"reward": {
				"denom": "ucmdx",
				"amount": "149954"
			}
		},
		{
			"address": "comdex1xdv6wl99344kdr3e3k4muhkzwvy3yr3g4awxuk",
			"reward": {
				"denom": "ucmdx",
				"amount": "731867"
			}
		},
		{
			"address": "comdex1xszdnlfvesfq3pcgksjc2pkr6twllga50vs8yp",
			"reward": {
				"denom": "ucmdx",
				"amount": "741915"
			}
		},
		{
			"address": "comdex1x334jvdcav3y75s5eczejp7sd8076gzyeyagh8",
			"reward": {
				"denom": "ucmdx",
				"amount": "299448"
			}
		},
		{
			"address": "comdex1xhdyruz4krz5mtsz3wpawvrp3m03cm0g47yhdu",
			"reward": {
				"denom": "ucmdx",
				"amount": "4951"
			}
		},
		{
			"address": "comdex1xe4jqaxets0yg5hsmpqtkansrycac4gkkwwzat",
			"reward": {
				"denom": "ucmdx",
				"amount": "382291"
			}
		},
		{
			"address": "comdex1x79e0allylw5t7ms53hdccs7wa69rzzj98996e",
			"reward": {
				"denom": "ucmdx",
				"amount": "6619"
			}
		},
		{
			"address": "comdex18qc780dez6hmzhvcrnsmxzkdnr64tlcv5tp9md",
			"reward": {
				"denom": "ucmdx",
				"amount": "3010"
			}
		},
		{
			"address": "comdex18z63wjpqkfw3s2ekv962shqeg55uqs7e4h9nkc",
			"reward": {
				"denom": "ucmdx",
				"amount": "1047"
			}
		},
		{
			"address": "comdex1899zef37g7zkrha970dhryjdf26udqfxh3jzs6",
			"reward": {
				"denom": "ucmdx",
				"amount": "1903159"
			}
		},
		{
			"address": "comdex18xhyg627flh7pmr6rhchngg6ctgk7ch38ysfme",
			"reward": {
				"denom": "ucmdx",
				"amount": "49927"
			}
		},
		{
			"address": "comdex18fljjw9a0u0q4ukye6l0azd5gtkp780lfe54m6",
			"reward": {
				"denom": "ucmdx",
				"amount": "1023"
			}
		},
		{
			"address": "comdex18n24cavw2xmds00jy5yqgkn6w5avrcjfluw5nl",
			"reward": {
				"denom": "ucmdx",
				"amount": "5420"
			}
		},
		{
			"address": "comdex184vnnwy3j4zmsjrml28y8heq47uhq8armz2xu7",
			"reward": {
				"denom": "ucmdx",
				"amount": "1673146"
			}
		},
		{
			"address": "comdex184swp2zc9vnaznlgtjfjnzgq5cfkuz6lvjffqk",
			"reward": {
				"denom": "ucmdx",
				"amount": "496600"
			}
		},
		{
			"address": "comdex18lv5fc9cscn0nxamlv0f0qp0f6h32c6d352duy",
			"reward": {
				"denom": "ucmdx",
				"amount": "450172"
			}
		},
		{
			"address": "comdex1ggn6an4ulcmk6ml8g28jg385maakc02me5xty0",
			"reward": {
				"denom": "ucmdx",
				"amount": "146559"
			}
		},
		{
			"address": "comdex1gg5lr2rrh5fwgjpmxgtmvyqft5dx73sjcq2ec7",
			"reward": {
				"denom": "ucmdx",
				"amount": "239062"
			}
		},
		{
			"address": "comdex1g2d6r276lnjdhlklemgc2agw4gms55pdfs5l7q",
			"reward": {
				"denom": "ucmdx",
				"amount": "19"
			}
		},
		{
			"address": "comdex1gt24kehn0jxgt5udev7lqa3nt4jpg5g55v80dg",
			"reward": {
				"denom": "ucmdx",
				"amount": "87290644"
			}
		},
		{
			"address": "comdex1gtkxagakmuf5tpukayj89kh9kkun77utqa0763",
			"reward": {
				"denom": "ucmdx",
				"amount": "1426"
			}
		},
		{
			"address": "comdex1gtaegw4m2n7ww538fn3chnmnda5mqe733fqg08",
			"reward": {
				"denom": "ucmdx",
				"amount": "22734"
			}
		},
		{
			"address": "comdex1gkc3kwss33xzum3dwdtlhpdv4ntyu9we6shuev",
			"reward": {
				"denom": "ucmdx",
				"amount": "87706"
			}
		},
		{
			"address": "comdex1gk737smj58de7pqhw8f5ud4kknjup2fv8v2h6g",
			"reward": {
				"denom": "ucmdx",
				"amount": "1658"
			}
		},
		{
			"address": "comdex1gu9g0hxrfvcajur6ahja6wha7f84v8ljcydjfv",
			"reward": {
				"denom": "ucmdx",
				"amount": "842"
			}
		},
		{
			"address": "comdex1f9hfkukcc9m4nzy896d3va63r9shxgpturwxl9",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1fgtezhnxtfsu9g3j3d973p76k8ldxgzmlz9sp2",
			"reward": {
				"denom": "ucmdx",
				"amount": "68243"
			}
		},
		{
			"address": "comdex1fwmuw744shcnuu8ctgf3t3ynsxkd3jv67n254v",
			"reward": {
				"denom": "ucmdx",
				"amount": "2783"
			}
		},
		{
			"address": "comdex1fnf4zkruvctm6ugkq4u9qw5t7vm23050tujm63",
			"reward": {
				"denom": "ucmdx",
				"amount": "95832"
			}
		},
		{
			"address": "comdex1f48744gpgjau49v5nvr5qjtwyh7u37meg4g27k",
			"reward": {
				"denom": "ucmdx",
				"amount": "13491"
			}
		},
		{
			"address": "comdex1fchyergst2n7wjlzdglg55xpdmu8ufanl6juuv",
			"reward": {
				"denom": "ucmdx",
				"amount": "38895"
			}
		},
		{
			"address": "comdex12puefzd0m6fnnf6junsxdrjd5mvqft80fa7f0e",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex129e9nqhpupsl9rxmd9t0ehms363ghyye0gq6s9",
			"reward": {
				"denom": "ucmdx",
				"amount": "266266"
			}
		},
		{
			"address": "comdex12xlrsfte5c7x45qd88e8t7pv34calr59aapsnf",
			"reward": {
				"denom": "ucmdx",
				"amount": "62994"
			}
		},
		{
			"address": "comdex128tulkkerslw0d7pnsqx560zqq5gkvmrtpwl4l",
			"reward": {
				"denom": "ucmdx",
				"amount": "3678"
			}
		},
		{
			"address": "comdex12g9tv7kzg06nx28szkvp7ud5jywpqv9sp0vdef",
			"reward": {
				"denom": "ucmdx",
				"amount": "88757"
			}
		},
		{
			"address": "comdex12gtfqzezj55s3ydmslvxsth4c2fnzvmmhkuq8k",
			"reward": {
				"denom": "ucmdx",
				"amount": "31580"
			}
		},
		{
			"address": "comdex12gmpp2fp2axr7r92w2na6ucq04svvka6n6sgh8",
			"reward": {
				"denom": "ucmdx",
				"amount": "461"
			}
		},
		{
			"address": "comdex12tzhxs2gwf5ddt0rewz9tq0w9s4y2xq256wtvx",
			"reward": {
				"denom": "ucmdx",
				"amount": "1396"
			}
		},
		{
			"address": "comdex123kvp3veh8aahvxrhe34m6rc888kj3p4tane76",
			"reward": {
				"denom": "ucmdx",
				"amount": "16985"
			}
		},
		{
			"address": "comdex12k24dpqqv75jqwu6a6jdt83lz79eqtl2dhdt33",
			"reward": {
				"denom": "ucmdx",
				"amount": "595218"
			}
		},
		{
			"address": "comdex12hg9fwu2hacynqez0vqchlu9k8efy9hefg4a2f",
			"reward": {
				"denom": "ucmdx",
				"amount": "122548"
			}
		},
		{
			"address": "comdex12euywa7npkwje3xd9tcvc64ye42vukna2za9jq",
			"reward": {
				"denom": "ucmdx",
				"amount": "210190"
			}
		},
		{
			"address": "comdex12e7x392uux4xx3suzf7xmhnzdmwalzumj26fv6",
			"reward": {
				"denom": "ucmdx",
				"amount": "19652"
			}
		},
		{
			"address": "comdex1tpy0t4kc2va5np7p8t5l62z3uc6708auqu05ye",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1t923w04lv3p5767d8rzs5dempjawaq6htykuuu",
			"reward": {
				"denom": "ucmdx",
				"amount": "30019"
			}
		},
		{
			"address": "comdex1tvz4rtnm3dtcw0y2m5a5e8ynes0kf5jtv4egcd",
			"reward": {
				"denom": "ucmdx",
				"amount": "4678023"
			}
		},
		{
			"address": "comdex1tvwd0mdve7k302xz4casgmwmsj65v4qwrj4jte",
			"reward": {
				"denom": "ucmdx",
				"amount": "16252"
			}
		},
		{
			"address": "comdex1td4y8k23pqyjn7pyz99fhvyp0ghxdzwfmvyfl4",
			"reward": {
				"denom": "ucmdx",
				"amount": "66119"
			}
		},
		{
			"address": "comdex1ts68fx54gf05ypmlstcngantfca8uzv64pkm6e",
			"reward": {
				"denom": "ucmdx",
				"amount": "1695"
			}
		},
		{
			"address": "comdex1tndcw4zgche8ykcgfq4rpy8fsna6027azy8c67",
			"reward": {
				"denom": "ucmdx",
				"amount": "250862"
			}
		},
		{
			"address": "comdex1t4ajsf03xxaluraclwlw3qly55r2fpcg8hykqq",
			"reward": {
				"denom": "ucmdx",
				"amount": "103073"
			}
		},
		{
			"address": "comdex1tc5pt52kjs6d9hk2xgf3c9a32tsnkdf4pdwzah",
			"reward": {
				"denom": "ucmdx",
				"amount": "588"
			}
		},
		{
			"address": "comdex1tm0vtepfzzggznjxxr0dmlwfgwe4ncr9l4n8ma",
			"reward": {
				"denom": "ucmdx",
				"amount": "51"
			}
		},
		{
			"address": "comdex1tackwayj6exxdexpj2qpjeenmdwcyrcyt0la05",
			"reward": {
				"denom": "ucmdx",
				"amount": "256"
			}
		},
		{
			"address": "comdex1t7ypyvmkh7wlnucz4rqs2kuhty5ss34ns39hsk",
			"reward": {
				"denom": "ucmdx",
				"amount": "87479"
			}
		},
		{
			"address": "comdex1vpmnmkrtvzgrdv50qqadp880ckylgc3z8542te",
			"reward": {
				"denom": "ucmdx",
				"amount": "4555"
			}
		},
		{
			"address": "comdex1vzxsfr6mpskfu8t3fafeaqm436t5kw87syk7y9",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1vzj2ke64ms039rgtmq5alz6v4vpqjjy8hnykdy",
			"reward": {
				"denom": "ucmdx",
				"amount": "8958340"
			}
		},
		{
			"address": "comdex1v92dav3kyarp8dsjqve9wqfvyhm5xxqgka29m9",
			"reward": {
				"denom": "ucmdx",
				"amount": "378"
			}
		},
		{
			"address": "comdex1v8qzgpnp0ewq0sda5rnk4eayq07eu7a0uw6rw2",
			"reward": {
				"denom": "ucmdx",
				"amount": "25970608"
			}
		},
		{
			"address": "comdex1v8n6xqzlq58h49uhed8t86y5uvcujqp45studv",
			"reward": {
				"denom": "ucmdx",
				"amount": "70137"
			}
		},
		{
			"address": "comdex1vgryc9wvm6f3lv4m868t08kmhunahnu6c2st7m",
			"reward": {
				"denom": "ucmdx",
				"amount": "5623"
			}
		},
		{
			"address": "comdex1v0gh3wmpc985u7mqmx6am4ep9fcxetqfr7k3dg",
			"reward": {
				"denom": "ucmdx",
				"amount": "16976"
			}
		},
		{
			"address": "comdex1v4442d72zdge6rtpal87nehagp70a0f9fzk5rl",
			"reward": {
				"denom": "ucmdx",
				"amount": "7179275"
			}
		},
		{
			"address": "comdex1ve5k4ul2fj42heup7ac7qvkhxzfzjrf6k4ktl4",
			"reward": {
				"denom": "ucmdx",
				"amount": "81242"
			}
		},
		{
			"address": "comdex1vuu4vaslmt3tmzxlq02wk2uq4mxu6grtzk3ka2",
			"reward": {
				"denom": "ucmdx",
				"amount": "289"
			}
		},
		{
			"address": "comdex1d87kckntmf2zsvh54khhtlce58va9dqdaa0206",
			"reward": {
				"denom": "ucmdx",
				"amount": "440354"
			}
		},
		{
			"address": "comdex1dfzjp5ks0g8duqvyprefy0zu3z54er52uj63pr",
			"reward": {
				"denom": "ucmdx",
				"amount": "854953"
			}
		},
		{
			"address": "comdex1dfmp46n8s4xxp9a638yazuarafvaz4hka2n2f7",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1dd8jqldgzexz4d6shz8rzx2ypn2wz6t5dnpgw5",
			"reward": {
				"denom": "ucmdx",
				"amount": "9984"
			}
		},
		{
			"address": "comdex1dwkcyrzzpt5dvqm9eyhygz7p52u7ce9hcp5jml",
			"reward": {
				"denom": "ucmdx",
				"amount": "1003187"
			}
		},
		{
			"address": "comdex1dhu8grdkd6sw9rw848kkfuf7e0x7a47esqy8jl",
			"reward": {
				"denom": "ucmdx",
				"amount": "190856"
			}
		},
		{
			"address": "comdex1dc6ctlxjdhaku5f36pzmrqtxvvm3w2de2zx3gw",
			"reward": {
				"denom": "ucmdx",
				"amount": "16308"
			}
		},
		{
			"address": "comdex1d6m9523kw0kp5c2auppx37t9vlu8fs8xyl0gvc",
			"reward": {
				"denom": "ucmdx",
				"amount": "8782"
			}
		},
		{
			"address": "comdex1wqga4v4nc77zah2n6nc9y0cdsu9hqcv4cpclgn",
			"reward": {
				"denom": "ucmdx",
				"amount": "189769"
			}
		},
		{
			"address": "comdex1wr8pdqarf933zf6glwt022fhaz8l0tv6wz8ypj",
			"reward": {
				"denom": "ucmdx",
				"amount": "244704"
			}
		},
		{
			"address": "comdex1wx3gcudvtxnsrvvhqzgqkuqxylmsqdc7xfxus0",
			"reward": {
				"denom": "ucmdx",
				"amount": "374002"
			}
		},
		{
			"address": "comdex1wx33kqfj0zs3kpsu90resx5e7zn9v5kh0r9ulg",
			"reward": {
				"denom": "ucmdx",
				"amount": "250143"
			}
		},
		{
			"address": "comdex1w5utzyjtn6a66g4d57gzrueqtn7lt63p8munqx",
			"reward": {
				"denom": "ucmdx",
				"amount": "20970"
			}
		},
		{
			"address": "comdex1wm48qr3d6qyx2ngzmju9tw9wqxeyv2ep796sce",
			"reward": {
				"denom": "ucmdx",
				"amount": "129013"
			}
		},
		{
			"address": "comdex1waywf486zhsetklxe3hze95yr763je282ywqnv",
			"reward": {
				"denom": "ucmdx",
				"amount": "7267"
			}
		},
		{
			"address": "comdex1wav9pnjeeq2e0fjgddlqv7pz2qctnge9n4ey6s",
			"reward": {
				"denom": "ucmdx",
				"amount": "69594"
			}
		},
		{
			"address": "comdex1wlp3s003m8yktac7gza0rnu8tum09v2tr6rtxd",
			"reward": {
				"denom": "ucmdx",
				"amount": "63573"
			}
		},
		{
			"address": "comdex10qsrwygnwrwfnz5npmvjd3namy4dadgk9snglk",
			"reward": {
				"denom": "ucmdx",
				"amount": "867"
			}
		},
		{
			"address": "comdex10rx4suuy4qy7dx4wl4x7kwvsv6f5y6a30wp6jf",
			"reward": {
				"denom": "ucmdx",
				"amount": "109893"
			}
		},
		{
			"address": "comdex109hq6dan6hg8k6nklk6gh0t0ey204rgc3p68dy",
			"reward": {
				"denom": "ucmdx",
				"amount": "9765"
			}
		},
		{
			"address": "comdex10gh63wta3mlccq9ljnjxudp26a95v7luvyseht",
			"reward": {
				"denom": "ucmdx",
				"amount": "26389"
			}
		},
		{
			"address": "comdex10vectgsg5vww7rka8qsp6c7f6yvc9l6mx95p8c",
			"reward": {
				"denom": "ucmdx",
				"amount": "50794"
			}
		},
		{
			"address": "comdex1002cc5zkqvlvennxy5a5jrkywa052gdhmxn3uc",
			"reward": {
				"denom": "ucmdx",
				"amount": "2008"
			}
		},
		{
			"address": "comdex100c867fxeu5dq6ht0flc7vf0rl65sfm7eaef0a",
			"reward": {
				"denom": "ucmdx",
				"amount": "5724"
			}
		},
		{
			"address": "comdex10jxducphn90d5yunvfemgehmfkchumu8ns3xv8",
			"reward": {
				"denom": "ucmdx",
				"amount": "10191"
			}
		},
		{
			"address": "comdex1059qwc07gxwu7fh9c7t6xsr9fgc2e0kssq953m",
			"reward": {
				"denom": "ucmdx",
				"amount": "646960"
			}
		},
		{
			"address": "comdex10m5mgyy0nq6qp9uq9luzv6vtxt3k2896xjppdv",
			"reward": {
				"denom": "ucmdx",
				"amount": "23784"
			}
		},
		{
			"address": "comdex10ufkg6c7qg40n4a5kjg4tf20xxzm7j45lyt5k5",
			"reward": {
				"denom": "ucmdx",
				"amount": "1380101"
			}
		},
		{
			"address": "comdex10awfv0z70w26ayyxgxrjd5n40guygzq4qmts5d",
			"reward": {
				"denom": "ucmdx",
				"amount": "8960"
			}
		},
		{
			"address": "comdex10lwh9n5hsxqhhma4nc26k72mghy0lr4drw8mag",
			"reward": {
				"denom": "ucmdx",
				"amount": "47190"
			}
		},
		{
			"address": "comdex1sztkmw500fllurt40svw34lkj74yf7ewkn7m0r",
			"reward": {
				"denom": "ucmdx",
				"amount": "72927"
			}
		},
		{
			"address": "comdex1szmqsu3rndntlh6fgkknc8hgzg0cazgfrqf4qm",
			"reward": {
				"denom": "ucmdx",
				"amount": "119774"
			}
		},
		{
			"address": "comdex1swxrvqjycece2zg6475rhaz3k4mr3smfemmepx",
			"reward": {
				"denom": "ucmdx",
				"amount": "77676"
			}
		},
		{
			"address": "comdex1ssu65vdunum4q7stec34ff3ftkwg34ztg3md34",
			"reward": {
				"denom": "ucmdx",
				"amount": "218163"
			}
		},
		{
			"address": "comdex1s3xq4zq7w3tl2knr9c2demg3v00zwnzehye4gh",
			"reward": {
				"denom": "ucmdx",
				"amount": "1622"
			}
		},
		{
			"address": "comdex1slvhtsyq3rsuypvzp6mne8qp36re9z3mzja0q8",
			"reward": {
				"denom": "ucmdx",
				"amount": "4551"
			}
		},
		{
			"address": "comdex13ql5gnxtayvzzy9unyjkx3yy6lg693x98pzkcz",
			"reward": {
				"denom": "ucmdx",
				"amount": "7766013"
			}
		},
		{
			"address": "comdex13rkq7eavqy6y23u3n7pu793ddqt2tvnlvpdha7",
			"reward": {
				"denom": "ucmdx",
				"amount": "667817"
			}
		},
		{
			"address": "comdex1385zfq3c3z2d8tqldn4u0g6vlqmdek2d823hh6",
			"reward": {
				"denom": "ucmdx",
				"amount": "3488"
			}
		},
		{
			"address": "comdex13fl5m99ypackm8jlt2l8aejpvj9vv3jl6zdfxa",
			"reward": {
				"denom": "ucmdx",
				"amount": "11559"
			}
		},
		{
			"address": "comdex13tp0aelqjq6wk6pgtxecr8uz54dspegcgnwc7r",
			"reward": {
				"denom": "ucmdx",
				"amount": "42017"
			}
		},
		{
			"address": "comdex133yu7zp8u902auwn2676x80sqpdw9hhd432usl",
			"reward": {
				"denom": "ucmdx",
				"amount": "19"
			}
		},
		{
			"address": "comdex135yvw9ug3qp3mk3g50hzz6z6egwwfh5q0af8ty",
			"reward": {
				"denom": "ucmdx",
				"amount": "761568"
			}
		},
		{
			"address": "comdex135svane9nts78jwhpuxz3r6vvf27nff25md9fh",
			"reward": {
				"denom": "ucmdx",
				"amount": "388929"
			}
		},
		{
			"address": "comdex135up9j7jnsk79cjhz9s0wkss7xtg2pka3ts7x7",
			"reward": {
				"denom": "ucmdx",
				"amount": "304619"
			}
		},
		{
			"address": "comdex13h9whxjw9vujgz7tlagrxa577zsfe0qs7s7rzq",
			"reward": {
				"denom": "ucmdx",
				"amount": "130170"
			}
		},
		{
			"address": "comdex1360878zq3qxhh9eunw4ua7udmvx47kvs879e4l",
			"reward": {
				"denom": "ucmdx",
				"amount": "547363"
			}
		},
		{
			"address": "comdex13a9tgvmnmap2e647xf8uaewwzz5sef7ds8pqf0",
			"reward": {
				"denom": "ucmdx",
				"amount": "39074"
			}
		},
		{
			"address": "comdex13lcazesvwpcrg9wla3gv5el04h9vkj669lmw3e",
			"reward": {
				"denom": "ucmdx",
				"amount": "2051"
			}
		},
		{
			"address": "comdex1jpmxyr33p3rmv9xum3g83frczvwy54yjhg9dzc",
			"reward": {
				"denom": "ucmdx",
				"amount": "9890"
			}
		},
		{
			"address": "comdex1jpmcg9xmjsapmxly830w73gg9hw97frmjtx7k4",
			"reward": {
				"denom": "ucmdx",
				"amount": "678922"
			}
		},
		{
			"address": "comdex1j9vf6gf23sx8ax82xdrmfgu0apyj9zaam5mwse",
			"reward": {
				"denom": "ucmdx",
				"amount": "32"
			}
		},
		{
			"address": "comdex1j9j8e63u225dg4sv7l7yf5fx3vphqy47xsz4rn",
			"reward": {
				"denom": "ucmdx",
				"amount": "7180"
			}
		},
		{
			"address": "comdex1jdnv6e0496ka39p9x5457knrdelvlf05d5sj99",
			"reward": {
				"denom": "ucmdx",
				"amount": "5035"
			}
		},
		{
			"address": "comdex1j0pcvwm9shud7dy6645ckelsfdzysgu8zp8a5h",
			"reward": {
				"denom": "ucmdx",
				"amount": "79373"
			}
		},
		{
			"address": "comdex1j5gumzvrna8hay2vpufaactsetvj2pc75vda0e",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1j5ftntpz7y2sda5t0y27ll26sj0vxm43w9r78x",
			"reward": {
				"denom": "ucmdx",
				"amount": "18629"
			}
		},
		{
			"address": "comdex1jhfn3q9l0sffv90n5egdjsv60nj6zefjn4fpls",
			"reward": {
				"denom": "ucmdx",
				"amount": "2385873"
			}
		},
		{
			"address": "comdex1j6tcsgrp262wfemp3wzn0502h4nha0u5kxtjj2",
			"reward": {
				"denom": "ucmdx",
				"amount": "386355"
			}
		},
		{
			"address": "comdex1j6jn5pu6yuwulalhq5402almr3y42hgjju7q3x",
			"reward": {
				"denom": "ucmdx",
				"amount": "102945"
			}
		},
		{
			"address": "comdex1jaxv0rkv8kp4pzkg8qwk7f2mdptn5xyvfyxasr",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1np0079lc7ga3qpxxtvn8me4jmvz469hlvtuud4",
			"reward": {
				"denom": "ucmdx",
				"amount": "72153"
			}
		},
		{
			"address": "comdex1n93lzalurfdadx43486s46luyph4497ue5z2jf",
			"reward": {
				"denom": "ucmdx",
				"amount": "29787"
			}
		},
		{
			"address": "comdex1ng2dc6d5mlsc8gwv8l5uv6t84va62537j0v0uw",
			"reward": {
				"denom": "ucmdx",
				"amount": "1799"
			}
		},
		{
			"address": "comdex1n2m75zj8d8j6g9x7yz9v5rpp62z4tsdlf8gvke",
			"reward": {
				"denom": "ucmdx",
				"amount": "93089"
			}
		},
		{
			"address": "comdex1nd86lwg7se890hau65eeq6ru9cxunc53pjgwf2",
			"reward": {
				"denom": "ucmdx",
				"amount": "17951"
			}
		},
		{
			"address": "comdex1nstkkxcx0zx0n4teud0a8s3azgxg2zt8uxzq5h",
			"reward": {
				"denom": "ucmdx",
				"amount": "68"
			}
		},
		{
			"address": "comdex1nsempd9ljljetpecuyc8m6agh0up8e43p36hpe",
			"reward": {
				"denom": "ucmdx",
				"amount": "46610"
			}
		},
		{
			"address": "comdex1n4e78ahrx4h9c85pncqfelxfrhhqnagujzus26",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1nchc75uswex4v5sa7fvprqpvr7jelnd3rwg5uv",
			"reward": {
				"denom": "ucmdx",
				"amount": "7200"
			}
		},
		{
			"address": "comdex1n7lyqen3hpkurptju8hec8sqsfzw8gs5pn9pf3",
			"reward": {
				"denom": "ucmdx",
				"amount": "56155"
			}
		},
		{
			"address": "comdex15pghxjl2t42m225pthqwp0qncncefzhw2hkhhn",
			"reward": {
				"denom": "ucmdx",
				"amount": "15324"
			}
		},
		{
			"address": "comdex15pn5hs98f2rt9m3jy2sjg4rvugp844hzgklsgx",
			"reward": {
				"denom": "ucmdx",
				"amount": "238937"
			}
		},
		{
			"address": "comdex15ystxc8h5cef6u3p5uulmleg55t3vvk9dc787g",
			"reward": {
				"denom": "ucmdx",
				"amount": "74654"
			}
		},
		{
			"address": "comdex15xxz5lwf6yqagzqq0z36g4twv5ju6rn645xnmr",
			"reward": {
				"denom": "ucmdx",
				"amount": "336703"
			}
		},
		{
			"address": "comdex158p9rzu3yl4y8ffsdz2m8zzem74tdj8tyfz88c",
			"reward": {
				"denom": "ucmdx",
				"amount": "1606676"
			}
		},
		{
			"address": "comdex158fw0vcpc9knyqeh8ckf6yxlnwgacwh9pra7ap",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex15f46k58ca87cj7844rgz4z47axpksq96lrm6jx",
			"reward": {
				"denom": "ucmdx",
				"amount": "326947"
			}
		},
		{
			"address": "comdex15t59r2kqx6wufdcgsf8mu97465u40pnsxcznvd",
			"reward": {
				"denom": "ucmdx",
				"amount": "48056"
			}
		},
		{
			"address": "comdex15jnuy7m2k2tu03lealcnymmnrerwwwurwjed4h",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex15neml2l2pyykysqwwe5ng5n93zx7gkv09xedug",
			"reward": {
				"denom": "ucmdx",
				"amount": "73086"
			}
		},
		{
			"address": "comdex15nau4akt3xn4t7x5rdq87l2utun496g8ydnl7k",
			"reward": {
				"denom": "ucmdx",
				"amount": "3128"
			}
		},
		{
			"address": "comdex15kczn5w2zpdzp8zkj234q9ldsgapzuzkcygl3a",
			"reward": {
				"denom": "ucmdx",
				"amount": "278387"
			}
		},
		{
			"address": "comdex156u9pnyfrvq2az0xvg6ntrvqmzc539xyw3s40n",
			"reward": {
				"denom": "ucmdx",
				"amount": "1881"
			}
		},
		{
			"address": "comdex14044c4jwgx5psuu6u3up02paccq4pwwcvwx7eg",
			"reward": {
				"denom": "ucmdx",
				"amount": "7052"
			}
		},
		{
			"address": "comdex14jxz7vdsx4r4zgc7jc23lfvna6y9rmezccprjw",
			"reward": {
				"denom": "ucmdx",
				"amount": "30192"
			}
		},
		{
			"address": "comdex1kyrveyaurv0aldd5a8rs92earptqfcr65d2zgl",
			"reward": {
				"denom": "ucmdx",
				"amount": "11228"
			}
		},
		{
			"address": "comdex1k8yv6mpppa5suf5g008g6wn8t23tpx6wspxz6k",
			"reward": {
				"denom": "ucmdx",
				"amount": "223698"
			}
		},
		{
			"address": "comdex1kfvqh8w34mn67658s0ukc8dem2sxwavmxs96e5",
			"reward": {
				"denom": "ucmdx",
				"amount": "498"
			}
		},
		{
			"address": "comdex1k2zuneuuml8tf4egwv4zyfmhr8t9m382d2ss9h",
			"reward": {
				"denom": "ucmdx",
				"amount": "139009"
			}
		},
		{
			"address": "comdex1k297velgrt692klywdjreelvjr97k03u845hm5",
			"reward": {
				"denom": "ucmdx",
				"amount": "41256"
			}
		},
		{
			"address": "comdex1k278k46nsvnc9ahx0n94qwt9qul7vea6putycd",
			"reward": {
				"denom": "ucmdx",
				"amount": "242353"
			}
		},
		{
			"address": "comdex1k0uwnf2fvclhwdtz62ermq9ssjx4qmsfyyx68s",
			"reward": {
				"denom": "ucmdx",
				"amount": "299786"
			}
		},
		{
			"address": "comdex1kscv7wsh5qlpaquegmt6c2dt2ake9gxkc7dyu2",
			"reward": {
				"denom": "ucmdx",
				"amount": "1990108"
			}
		},
		{
			"address": "comdex1k3f8y623r9qr4lg4jckmcxxmsqdjcl34ndzelq",
			"reward": {
				"denom": "ucmdx",
				"amount": "43873"
			}
		},
		{
			"address": "comdex1k6cd4ev5ecvrdneyufld6fz8hcn5vzpk2xe77y",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1kax69ywsyncwph3ll3l0gzv2wc49athz4rw8zq",
			"reward": {
				"denom": "ucmdx",
				"amount": "626708"
			}
		},
		{
			"address": "comdex1hgzdrnelpd7gv0ndt8fn3kjft4fhv8v0txvq9m",
			"reward": {
				"denom": "ucmdx",
				"amount": "2587"
			}
		},
		{
			"address": "comdex1hfpxf58zjc0r7h9z0qnzcm88rlwdyuwddp5lmt",
			"reward": {
				"denom": "ucmdx",
				"amount": "283996"
			}
		},
		{
			"address": "comdex1hfygllaur5879xnt4qgq2653n0y9g5vqt9f9yw",
			"reward": {
				"denom": "ucmdx",
				"amount": "150582"
			}
		},
		{
			"address": "comdex1h2cs4py349m2j3ra2k85wkcffpfztqw8gwgaxt",
			"reward": {
				"denom": "ucmdx",
				"amount": "1554760"
			}
		},
		{
			"address": "comdex1hdz0fz5z20ldt0aw0zxcqc4lt9p4zn0nm8njnz",
			"reward": {
				"denom": "ucmdx",
				"amount": "141"
			}
		},
		{
			"address": "comdex1hdusu20ltnm5086vjjnftk62cvnlmuvqmjmaru",
			"reward": {
				"denom": "ucmdx",
				"amount": "1892545"
			}
		},
		{
			"address": "comdex1h0ucd9pyvc6x98244p8f7jkrud3n9w729t6wfx",
			"reward": {
				"denom": "ucmdx",
				"amount": "4995"
			}
		},
		{
			"address": "comdex1h507jazzrcewj44ez0876pyadv47x7tzsknmp4",
			"reward": {
				"denom": "ucmdx",
				"amount": "3467739"
			}
		},
		{
			"address": "comdex1h57snd9e2f6zpejxamg2w8fzxckdwth6tjc6vh",
			"reward": {
				"denom": "ucmdx",
				"amount": "642649"
			}
		},
		{
			"address": "comdex1h4v8jxrwzex5yj6d49jm92a827sempjtelnak5",
			"reward": {
				"denom": "ucmdx",
				"amount": "267396"
			}
		},
		{
			"address": "comdex1h6h06aqc79fcvju8gpwxurtdt9thgyu577k8zr",
			"reward": {
				"denom": "ucmdx",
				"amount": "25148"
			}
		},
		{
			"address": "comdex1hmll8lejjsvq4ldwt6k5jdmrzmj38yjv3chjlp",
			"reward": {
				"denom": "ucmdx",
				"amount": "5824"
			}
		},
		{
			"address": "comdex1hlhtlry8ymehlffskfasq3kq2vvvlvyy6u4yjk",
			"reward": {
				"denom": "ucmdx",
				"amount": "1271"
			}
		},
		{
			"address": "comdex1cze68e2fu3t4ywsrjncgx8c39ph97hzaftcfvn",
			"reward": {
				"denom": "ucmdx",
				"amount": "81256"
			}
		},
		{
			"address": "comdex1c95w0kvata3pugcsrs99a69pvf8hszvaugtd5w",
			"reward": {
				"denom": "ucmdx",
				"amount": "35079"
			}
		},
		{
			"address": "comdex1cgev9ly5znkg3qgpsrn3ftrtmakukyjq5l3t5z",
			"reward": {
				"denom": "ucmdx",
				"amount": "170155"
			}
		},
		{
			"address": "comdex1cvfwd6j2d2em7ygew5lj379vmrl4ef052ces9s",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1cvwg7gmjawqd7kpqjuk8n60s9223th6lxqntd3",
			"reward": {
				"denom": "ucmdx",
				"amount": "13393"
			}
		},
		{
			"address": "comdex1c0sgcfj6x4qh9wnz0ruzf8nhhfgpy6n87pn6eg",
			"reward": {
				"denom": "ucmdx",
				"amount": "100"
			}
		},
		{
			"address": "comdex1csf3up94qheqz88fvvlwjl7vc5sdymarjh4ax7",
			"reward": {
				"denom": "ucmdx",
				"amount": "619747"
			}
		},
		{
			"address": "comdex1ce6czv3lvs7j2yehu407j2p5xg7z9k9ldsjpev",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1cagl36lucd5sfwzrquhncqyf67w0jtrjp95mfp",
			"reward": {
				"denom": "ucmdx",
				"amount": "64355"
			}
		},
		{
			"address": "comdex1c75rn45dtq5d84rzmzp25fxgm67aq982hlq49p",
			"reward": {
				"denom": "ucmdx",
				"amount": "1014678"
			}
		},
		{
			"address": "comdex1ey7ztz4dvxje56ckcvfhqfkssmu90e0cch8uef",
			"reward": {
				"denom": "ucmdx",
				"amount": "5090"
			}
		},
		{
			"address": "comdex1e84e9gkftuwx8pg9w4qax7f70nhf20gt89gx9c",
			"reward": {
				"denom": "ucmdx",
				"amount": "60480"
			}
		},
		{
			"address": "comdex1e8cxrds3hkkkgv2wuu7knnjgeyfnzpt3eaw4z9",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1efj6emzg9zp7gyxre055s3kdp8j595jww8mn2w",
			"reward": {
				"denom": "ucmdx",
				"amount": "156333"
			}
		},
		{
			"address": "comdex1e26k332qmtc8y4r495356gfchlpclmzszhv9hg",
			"reward": {
				"denom": "ucmdx",
				"amount": "5025"
			}
		},
		{
			"address": "comdex1e2u8lgv98z4uwp5qs8jye9ahtcnkr9pd93elsw",
			"reward": {
				"denom": "ucmdx",
				"amount": "922761"
			}
		},
		{
			"address": "comdex1ev49nxg3e6y23z9ayvkgj3vlmg26h727q7ddft",
			"reward": {
				"denom": "ucmdx",
				"amount": "90106"
			}
		},
		{
			"address": "comdex1ewu290xt5fhfz23jk34rq266ymzp6he6jrd338",
			"reward": {
				"denom": "ucmdx",
				"amount": "3138042"
			}
		},
		{
			"address": "comdex1e0v5qu9hknmtpy3ha0ye9vh2pqvsdn8hc0rp8y",
			"reward": {
				"denom": "ucmdx",
				"amount": "319024"
			}
		},
		{
			"address": "comdex1en4mvfpjsrk3el4xlu568hmt4ke02kz4ze6w63",
			"reward": {
				"denom": "ucmdx",
				"amount": "19"
			}
		},
		{
			"address": "comdex1e5tush5f5nxafgqluvwkc5pf8kkkc0njmp9905",
			"reward": {
				"denom": "ucmdx",
				"amount": "6156"
			}
		},
		{
			"address": "comdex1e4sfeuxv2zqxaswpngm43mp5rsvzta75eexcng",
			"reward": {
				"denom": "ucmdx",
				"amount": "589198"
			}
		},
		{
			"address": "comdex1ek6a6p7sys2hj50ekqdaht2ana82jwz9ylfmjs",
			"reward": {
				"denom": "ucmdx",
				"amount": "252285"
			}
		},
		{
			"address": "comdex1eh4x428myh6z4r58z58cdadaaegx753wd4mrru",
			"reward": {
				"denom": "ucmdx",
				"amount": "36698"
			}
		},
		{
			"address": "comdex1eh6p8n8lksz28nup38muylvw8j4wyhyv6u3e4m",
			"reward": {
				"denom": "ucmdx",
				"amount": "202764"
			}
		},
		{
			"address": "comdex1eejdeyvjp923vam3krdfc4vf7766srewlt2z5m",
			"reward": {
				"denom": "ucmdx",
				"amount": "476"
			}
		},
		{
			"address": "comdex16p3mxsqexfmxkp2pvyde9nr53k5dhh960c2p6j",
			"reward": {
				"denom": "ucmdx",
				"amount": "4535"
			}
		},
		{
			"address": "comdex16zghh2vu2cfxgm83dt8qz6pkqvaumukjx9at4u",
			"reward": {
				"denom": "ucmdx",
				"amount": "1073573"
			}
		},
		{
			"address": "comdex16gld4flx9uy3lf0v02c5mh89juh53pargvhyl7",
			"reward": {
				"denom": "ucmdx",
				"amount": "10246"
			}
		},
		{
			"address": "comdex16s3nqpvl6jypaxg5j2u68kejf5fuegtxydmkcv",
			"reward": {
				"denom": "ucmdx",
				"amount": "96904"
			}
		},
		{
			"address": "comdex16n6klgn79dvey5mc5k73tvtm9nn58r77nya9rz",
			"reward": {
				"denom": "ucmdx",
				"amount": "2500"
			}
		},
		{
			"address": "comdex16kttfdm0q0lwlj9zk4tart7fhxwj64xxhsfp2x",
			"reward": {
				"denom": "ucmdx",
				"amount": "21"
			}
		},
		{
			"address": "comdex16eraz5pxjq0g285gsqgyq9xlsjarqvdsve5y36",
			"reward": {
				"denom": "ucmdx",
				"amount": "21171"
			}
		},
		{
			"address": "comdex16u7us6gvunwquhecdk6x3n43hd36yuq727qewk",
			"reward": {
				"denom": "ucmdx",
				"amount": "267546"
			}
		},
		{
			"address": "comdex1mpc0lk2xnq32r477upqxf8x3k65ke2fmzvsj7x",
			"reward": {
				"denom": "ucmdx",
				"amount": "50387"
			}
		},
		{
			"address": "comdex1m9ux94a89pygdpp3vd4spx6q46gsg77gxq6pl8",
			"reward": {
				"denom": "ucmdx",
				"amount": "99718"
			}
		},
		{
			"address": "comdex1mxquqp8anmjnj6e5n87fuam2cskfpwqd5tvhwl",
			"reward": {
				"denom": "ucmdx",
				"amount": "705"
			}
		},
		{
			"address": "comdex1mtxa2xj0w0vn44vlcyja7jx2kv6dfc5h9csynk",
			"reward": {
				"denom": "ucmdx",
				"amount": "62"
			}
		},
		{
			"address": "comdex1m3r7lr25e6yqx02dqy2e8alrt3hlwmsr8p4t7p",
			"reward": {
				"denom": "ucmdx",
				"amount": "90899"
			}
		},
		{
			"address": "comdex1m5rtxaknkj39veg8u9epa3lqs2k7wpftp6rcpx",
			"reward": {
				"denom": "ucmdx",
				"amount": "4583"
			}
		},
		{
			"address": "comdex1m5863nszf8qqugdsgzdsc2uz7fgrwvt743f4n8",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1mhr7pdf44xsf9n5hta2t0xf5puqgu6gpeu40wz",
			"reward": {
				"denom": "ucmdx",
				"amount": "1469"
			}
		},
		{
			"address": "comdex1mmqslm3xc7vlh254cw07369ps3tna6vhsqrqul",
			"reward": {
				"denom": "ucmdx",
				"amount": "18045"
			}
		},
		{
			"address": "comdex1m799kf6nhve8nr4v6mpv3ke8qwc5e6jft262vg",
			"reward": {
				"denom": "ucmdx",
				"amount": "33985"
			}
		},
		{
			"address": "comdex1m7glw5w5tamprdx5h5rkge2jsjlpm70d7cc55u",
			"reward": {
				"denom": "ucmdx",
				"amount": "345597"
			}
		},
		{
			"address": "comdex1m7lqalv2s5qkh283ur7382e0dhwnrw55e4a3uj",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1uqavktrl3zkn4t0fe0dmwmampta25un2uguc8w",
			"reward": {
				"denom": "ucmdx",
				"amount": "1072705"
			}
		},
		{
			"address": "comdex1uzy4az5323lvv5ceg4r3lyenta8uv8vj7eq3yv",
			"reward": {
				"denom": "ucmdx",
				"amount": "66003"
			}
		},
		{
			"address": "comdex1ux35kxht4dhxnp43z0q9cdjfcfz05qu0a4h5pm",
			"reward": {
				"denom": "ucmdx",
				"amount": "88969"
			}
		},
		{
			"address": "comdex1uggj0h68ed5xy8scaq2vp5vwm4lucguenvszqh",
			"reward": {
				"denom": "ucmdx",
				"amount": "370723"
			}
		},
		{
			"address": "comdex1utxx4enrhnsstwnzkrsrfdufkyxwx7t788xqwq",
			"reward": {
				"denom": "ucmdx",
				"amount": "7239"
			}
		},
		{
			"address": "comdex1uw73zg3qh6f6fvc8dnecl05jgxfx9mkqprx7u2",
			"reward": {
				"denom": "ucmdx",
				"amount": "53977"
			}
		},
		{
			"address": "comdex1u4g5vc0lhdatsa7ntkrhmp8gdeglt53layhwhq",
			"reward": {
				"denom": "ucmdx",
				"amount": "65592"
			}
		},
		{
			"address": "comdex1umztzkqckfu0qtt7e4dwkvp5ryd3dgh0ct4cls",
			"reward": {
				"denom": "ucmdx",
				"amount": "26262"
			}
		},
		{
			"address": "comdex1ax4e7jl8hwnnnsltjtqu04mmp7cpulwz73vpn3",
			"reward": {
				"denom": "ucmdx",
				"amount": "317693"
			}
		},
		{
			"address": "comdex1a82wetv2az49h4w2qke9zjm09s6mk70ce8zcpk",
			"reward": {
				"denom": "ucmdx",
				"amount": "63661"
			}
		},
		{
			"address": "comdex1agyhst9990lykr5w6gdwckh2kzg436pqthtdhv",
			"reward": {
				"denom": "ucmdx",
				"amount": "305838"
			}
		},
		{
			"address": "comdex1agcqex33msxvwr9ghpx37p2xmt2hh0dy5nwcqh",
			"reward": {
				"denom": "ucmdx",
				"amount": "23074"
			}
		},
		{
			"address": "comdex1a2y430s4qqk5cr2pcgsrtzd9fd5q3g40vhxkph",
			"reward": {
				"denom": "ucmdx",
				"amount": "111394"
			}
		},
		{
			"address": "comdex1adgkq0j84y8c70ps7hwsvetcsn8c5zrh75r3r5",
			"reward": {
				"denom": "ucmdx",
				"amount": "1226"
			}
		},
		{
			"address": "comdex1asq7mvh0re6euq5qg3426ayc4uuc3rpwk3q5a0",
			"reward": {
				"denom": "ucmdx",
				"amount": "21694"
			}
		},
		{
			"address": "comdex1asjgj3y9n3xjagrt82965cdswx6tnsy474rn9j",
			"reward": {
				"denom": "ucmdx",
				"amount": "203360"
			}
		},
		{
			"address": "comdex1ajhz8n9yawladzkhlny6cyd03kr49xqumd783a",
			"reward": {
				"denom": "ucmdx",
				"amount": "734239"
			}
		},
		{
			"address": "comdex1ajl8ga754hmems6c33ljxerxvpw9detpa6qa4x",
			"reward": {
				"denom": "ucmdx",
				"amount": "165017"
			}
		},
		{
			"address": "comdex1ajl2kfddkerlp4msrya4wc0r373g4vtndf0fuc",
			"reward": {
				"denom": "ucmdx",
				"amount": "6021"
			}
		},
		{
			"address": "comdex1a4jkhl5f3se7pcfkuvrzslprkz7rmlaqat3cw3",
			"reward": {
				"denom": "ucmdx",
				"amount": "113106"
			}
		},
		{
			"address": "comdex1au7c7avqa7p980ut6qgnudjwhda02efv8r4tm5",
			"reward": {
				"denom": "ucmdx",
				"amount": "80656"
			}
		},
		{
			"address": "comdex1aaprhhcx8c0g5vtq9eyq56dvpwvpzswpsf3vfe",
			"reward": {
				"denom": "ucmdx",
				"amount": "126413"
			}
		},
		{
			"address": "comdex17puvlewlvvg9vamnlrt5exyjkjcckkh8c6eyup",
			"reward": {
				"denom": "ucmdx",
				"amount": "233630"
			}
		},
		{
			"address": "comdex17zl64anzay9vzvyw5tzffgnnd42jjh29wkl4wv",
			"reward": {
				"denom": "ucmdx",
				"amount": "145615"
			}
		},
		{
			"address": "comdex17rzlk2yvhfpvlzmfxwqtawlj3a46cfmge7t654",
			"reward": {
				"denom": "ucmdx",
				"amount": "52147"
			}
		},
		{
			"address": "comdex17tfzw54en5tcaz20gw5stppn857c8dss9tz3ag",
			"reward": {
				"denom": "ucmdx",
				"amount": "4714"
			}
		},
		{
			"address": "comdex17wjj4e4xjlk6qtdk3uw432vs3pmvkv0d74qusd",
			"reward": {
				"denom": "ucmdx",
				"amount": "5205520"
			}
		},
		{
			"address": "comdex17jv2g8y0ckmelgn39272h8nat9fjs20ytdrd4v",
			"reward": {
				"denom": "ucmdx",
				"amount": "70296"
			}
		},
		{
			"address": "comdex1lzpz72s3gz3vhrtap4vmcuz4daf8nf3vuwv083",
			"reward": {
				"denom": "ucmdx",
				"amount": "79817"
			}
		},
		{
			"address": "comdex1lgu3kck7ct9ejarkjdf7at680f30yyzmlnt3te",
			"reward": {
				"denom": "ucmdx",
				"amount": "66230"
			}
		},
		{
			"address": "comdex1lwkr9wpm35e93tnjj9awhhph0p6620n0agk29n",
			"reward": {
				"denom": "ucmdx",
				"amount": "38163"
			}
		},
		{
			"address": "comdex1l3hkzt9tzf59xah5vdwjem3gg7s4cvasz7q24h",
			"reward": {
				"denom": "ucmdx",
				"amount": "44946"
			}
		},
		{
			"address": "comdex1l5dnjn6s904hugw0m6yrz26nr8cd84vfuhujay",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex19tm2nnw960kwjdd2y48cvce93z8wjz4argl73r",
			"reward": {
				"denom": "ucmdx",
				"amount": "924200"
			}
		},
		{
			"address": "comdex1gffcj0fcuzkfqeknlrvvcyzqp5xng9s9flfk9d",
			"reward": {
				"denom": "ucmdx",
				"amount": "45717"
			}
		},
		{
			"address": "comdex12z3uxxzsalafwnlnt0xc2glm8c9h9dmyurz50q",
			"reward": {
				"denom": "ucmdx",
				"amount": "121112"
			}
		},
		{
			"address": "comdex1j6mtmgc2mcc4jvhgnxp87ggr03qhly3lgwedzj",
			"reward": {
				"denom": "ucmdx",
				"amount": "183138"
			}
		},
		{
			"address": "comdex1442ks20qpww9e2rkt4edxj6qcjzhm500vyv9pw",
			"reward": {
				"denom": "ucmdx",
				"amount": "243576"
			}
		},
		{
			"address": "comdex1kf4hnapmyj5je6egugvusdqkrju0rejnxw4kcc",
			"reward": {
				"denom": "ucmdx",
				"amount": "58423092"
			}
		},
		{
			"address": "comdex1cujykzh8z0nx6pjp5yvrmv50jes5shh6a9ht8t",
			"reward": {
				"denom": "ucmdx",
				"amount": "55072"
			}
		},
		{
			"address": "comdex1eqhl73ed0kwccnvjqzq28cngcp5u38y6yeepwd",
			"reward": {
				"denom": "ucmdx",
				"amount": "0"
			}
		},
		{
			"address": "comdex1uj3m5ly0qfwx3tkzgnk55hrektr56d05t892c4",
			"reward": {
				"denom": "ucmdx",
				"amount": "7162181"
			}
		},
		{
			"address": "comdex1l3rrp7dlxfjm233jc6u564zxwuedfk6ysvt69c",
			"reward": {
				"denom": "ucmdx",
				"amount": "14371775"
			}
		},
		{
			"address": "comdex1tx5fu2d7vdyjn9jt7ppm9447n5xzj7m030c5l2",
			"reward": {
				"denom": "ucmdx",
				"amount": "10624079"
			}
		},
		{
			"address": "comdex1vqe538w5066vtswds98equcccdwa7j0mgkz9xg",
			"reward": {
				"denom": "ucmdx",
				"amount": "1861455"
			}
		},
		{
			"address": "comdex1wn6w5xs7n5q0lvxjsap75ytthmn93rsv5hurs9",
			"reward": {
				"denom": "ucmdx",
				"amount": "1335216"
			}
		},
		{
			"address": "comdex10902j8s6hr6xgaxzq372ns0mdcyxg4c0xhyr9f",
			"reward": {
				"denom": "ucmdx",
				"amount": "44179348"
			}
		},
		{
			"address": "comdex1k9gn52902fhhm4pgmku0hw47h59ttldfuvkc5c",
			"reward": {
				"denom": "ucmdx",
				"amount": "3481508"
			}
		},
		{
			"address": "comdex1kjp9jl5xfjrevltffwulr64enlfzkpcvk7rcgd",
			"reward": {
				"denom": "ucmdx",
				"amount": "9370077"
			}
		},
		{
			"address": "comdex1er3ksjsauh6teycrtpj5gzp9ltj08zufvhuxhd",
			"reward": {
				"denom": "ucmdx",
				"amount": "44736"
			}
		},
		{
			"address": "comdex1udtnr7yzyz2nn62t478dsm9kagz7feemrrx3mj",
			"reward": {
				"denom": "ucmdx",
				"amount": "13763732"
			}
		},
		{
			"address": "comdex1a8nw4jhwdqfus5fn8ep9g68sxghdm4jxa2z55j",
			"reward": {
				"denom": "ucmdx",
				"amount": "5006"
			}
		},
		{
			"address": "comdex17rt5ny8earrvpxq9xdvf9z6y67wwm2t90tg29s",
			"reward": {
				"denom": "ucmdx",
				"amount": "578689"
			}
		},
		{
			"address": "comdex178swfpjp3qwes3wj8ffzqla9m56ewzxmqzyhfn",
			"reward": {
				"denom": "ucmdx",
				"amount": "220357"
			}
		},
		{
			"address": "comdex17w55ykyx4nj3vzju2nju42xnfpe50rzxynct3w",
			"reward": {
				"denom": "ucmdx",
				"amount": "4585758"
			}
		},
		{
			"address": "comdex1l5x5tzg9v723nurtvkppst9e0q3lmkncvfpy09",
			"reward": {
				"denom": "ucmdx",
				"amount": "84475"
			}
		},
		{
			"address": "comdex1ym6x2ht9xl6yjkqtqkkudt2ktjzckq0y43jwzq",
			"reward": {
				"denom": "ucmdx",
				"amount": "111246"
			}
		},
		{
			"address": "comdex186kettvmn8chcax63zsuz75ketzdj8tm5jfgrt",
			"reward": {
				"denom": "ucmdx",
				"amount": "2437259"
			}
		},
		{
			"address": "comdex1vuz3smamsfnrk65ktd9r9klnjfqyf44ggkvw27",
			"reward": {
				"denom": "ucmdx",
				"amount": "2967656"
			}
		},
		{
			"address": "comdex1wvh4mga33wshf3qlt5h9q8sf7h2kqlthlnf2tz",
			"reward": {
				"denom": "ucmdx",
				"amount": "36603519"
			}
		},
		{
			"address": "comdex1njn4zjjrhs89nfwskjtzerchdkur5pfxqxupln",
			"reward": {
				"denom": "ucmdx",
				"amount": "15658973"
			}
		},
		{
			"address": "comdex15r675pn9wccz66w66p7aqj82w2etg6m7xq4f58",
			"reward": {
				"denom": "ucmdx",
				"amount": "32483846"
			}
		},
		{
			"address": "comdex1rd5v839pndaxn2kkpj9s9dy79sy4q9dykvsl6t",
			"reward": {
				"denom": "ucmdx",
				"amount": "665973"
			}
		},
		{
			"address": "comdex19226sex0s9x6k6trpadcyhf75qenzuc0c645yp",
			"reward": {
				"denom": "ucmdx",
				"amount": "116498"
			}
		},
		{
			"address": "comdex1xrv89wtat29mkgr3y3hed7pn2a4uxp73q56cq5",
			"reward": {
				"denom": "ucmdx",
				"amount": "2518567"
			}
		},
		{
			"address": "comdex13dkyl4knw325jdd526kfleu5jzavlegtmp0ef7",
			"reward": {
				"denom": "ucmdx",
				"amount": "158"
			}
		},
		{
			"address": "comdex1jtz4sj4zw67t0qedfh4nzdh2fp4h79wwp3rfhe",
			"reward": {
				"denom": "ucmdx",
				"amount": "7387987"
			}
		},
		{
			"address": "comdex153l8ng6puwut2kcjj74fk5wx34jaktckl2cz58",
			"reward": {
				"denom": "ucmdx",
				"amount": "4527830"
			}
		},
		{
			"address": "comdex1kjxa7524tcjzcmtvld0z8vny53ps084vylpwaz",
			"reward": {
				"denom": "ucmdx",
				"amount": "22859"
			}
		},
		{
			"address": "comdex1enz0knm7yl8932946n26pvnenxcr9w3ur3z2gz",
			"reward": {
				"denom": "ucmdx",
				"amount": "69565"
			}
		},
		{
			"address": "comdex1ljwsja7mphyj4ua7cm4d35c0mzevvnpq53tpd6",
			"reward": {
				"denom": "ucmdx",
				"amount": "6277503"
			}
		},
		{
			"address": "comdex1prcmpmwggaxt8v8csvcpz9vc6dnnrr5qsv7rr8",
			"reward": {
				"denom": "ucmdx",
				"amount": "19503"
			}
		},
		{
			"address": "comdex1zkh736hh9mkmpa37mh4qdsm07j90c3za80fpmw",
			"reward": {
				"denom": "ucmdx",
				"amount": "95290"
			}
		},
		{
			"address": "comdex1rc9kyu89k4qpfhj4mhaa66n5n03pkd5mm8r844",
			"reward": {
				"denom": "ucmdx",
				"amount": "12845335"
			}
		},
		{
			"address": "comdex1y67nmrx5ykeuasp3kupvaa59nu2p65zfs88rd8",
			"reward": {
				"denom": "ucmdx",
				"amount": "8210"
			}
		},
		{
			"address": "comdex19jz4t3zu2u6nvc54l3sgyjvh4n562vxncse9r8",
			"reward": {
				"denom": "ucmdx",
				"amount": "6502614"
			}
		},
		{
			"address": "comdex1xpfdf7nfsxygh0gu5ewgaw9gqlf7y2fu6smrn5",
			"reward": {
				"denom": "ucmdx",
				"amount": "166937"
			}
		},
		{
			"address": "comdex12y4gqmav0sums6896s048a8l5nkrj76xlnghpj",
			"reward": {
				"denom": "ucmdx",
				"amount": "4167"
			}
		},
		{
			"address": "comdex1vss9v2rq0ggnnyn4ypu3dpu8djuuhveu7erlsq",
			"reward": {
				"denom": "ucmdx",
				"amount": "1314347"
			}
		},
		{
			"address": "comdex1v5k4st3duzftjsrd4m82gwvz8aygcgz8yj6jm6",
			"reward": {
				"denom": "ucmdx",
				"amount": "11305200"
			}
		},
		{
			"address": "comdex1v4gq9z4hsf5ma6kfs3mvu3jlyrwz3t84pputgu",
			"reward": {
				"denom": "ucmdx",
				"amount": "19956"
			}
		},
		{
			"address": "comdex1ddmmtady24ncs5z8k5tuqyvpcfvxjc6ushq5h6",
			"reward": {
				"denom": "ucmdx",
				"amount": "1556600"
			}
		},
		{
			"address": "comdex1d47cdt2s0hz9z39dg82t5jactszp69ef7rr6k4",
			"reward": {
				"denom": "ucmdx",
				"amount": "20200"
			}
		},
		{
			"address": "comdex1srzz5fckw302usl2w0kresl8qpmqdmmxhhm7nr",
			"reward": {
				"denom": "ucmdx",
				"amount": "7477454"
			}
		},
		{
			"address": "comdex13z2arscmsdttujpdmc9dkrw2hkt043ersn6mca",
			"reward": {
				"denom": "ucmdx",
				"amount": "654"
			}
		},
		{
			"address": "comdex1j9t4sqsys4vkvl2eqr4q5vm90eu8rl9w69t8cq",
			"reward": {
				"denom": "ucmdx",
				"amount": "6877"
			}
		},
		{
			"address": "comdex1neawwhn97nlnpft80gtk7m6feh745rhyjdfsc7",
			"reward": {
				"denom": "ucmdx",
				"amount": "38797"
			}
		},
		{
			"address": "comdex153q8f6xk2tjl48kvm36pr3fqrsyamus8vfdyxa",
			"reward": {
				"denom": "ucmdx",
				"amount": "6918532"
			}
		},
		{
			"address": "comdex15n4jyvyh9grdh27awkpny5wqangm0tmyc0k5h6",
			"reward": {
				"denom": "ucmdx",
				"amount": "2506"
			}
		},
		{
			"address": "comdex1hls4rdk5c9c8ff2cdgrjnu3x4dg0v75n74lgv9",
			"reward": {
				"denom": "ucmdx",
				"amount": "106408"
			}
		},
		{
			"address": "comdex1entd98ckj04wpdn8ce73ma87swvmgrk3kh3z2s",
			"reward": {
				"denom": "ucmdx",
				"amount": "9960"
			}
		},
		{
			"address": "comdex1eeyvsyzpf63wrhu3dwtseku8darlsmraajyce8",
			"reward": {
				"denom": "ucmdx",
				"amount": "3417989"
			}
		},
		{
			"address": "comdex16ars5rflyxqkfamdfs96fj6ncz27xuvw29flx4",
			"reward": {
				"denom": "ucmdx",
				"amount": "8318"
			}
		},
		{
			"address": "comdex1u6ez787csmez2rfx6fsqhkcl294ffp7ntgumc9",
			"reward": {
				"denom": "ucmdx",
				"amount": "263190"
			}
		},
		{
			"address": "comdex1a4yxyhy4c7cvyrqsyddvgrp0g6nxd486mpglqg",
			"reward": {
				"denom": "ucmdx",
				"amount": "876826"
			}
		},
		{
			"address": "comdex1gvfny70eqvu4wmqhzlnxkqkynu3w4rmz35zrjk",
			"reward": {
				"denom": "ucmdx",
				"amount": "2613229"
			}
		},
		{
			"address": "comdex1thccnsyv4qa3r7kqgfgc4k9xu7wf2qc096lfdk",
			"reward": {
				"denom": "ucmdx",
				"amount": "6607821"
			}
		},
		{
			"address": "comdex1lyerd3m23ms8598rpuep0lwl3d6qc259yluk7w",
			"reward": {
				"denom": "ucmdx",
				"amount": "4701"
			}
		},
		{
			"address": "comdex1ej26n0p8fugn3g74l8ggf74md8mc5jry54d6ey",
			"reward": {
				"denom": "ucmdx",
				"amount": "9717020"
			}
		},
		{
			"address": "comdex1mn7mh2fjcampcp2dcygqrhp5gf3ley95eqfres",
			"reward": {
				"denom": "ucmdx",
				"amount": "1491050"
			}
		},
		{
			"address": "comdex1puqcqc0xappgeupw8gfffyztzrdp7l6q3e8lak",
			"reward": {
				"denom": "ucmdx",
				"amount": "838427"
			}
		},
		{
			"address": "comdex1rdm7hnfk8tuwvmvegydrxda3khkvknhs7hpggy",
			"reward": {
				"denom": "ucmdx",
				"amount": "4886765"
			}
		},
		{
			"address": "comdex1ya6ka4532mlw6ca64f2mu9yfcfdqvgxtjg6fjg",
			"reward": {
				"denom": "ucmdx",
				"amount": "2887957"
			}
		},
		{
			"address": "comdex198mhlu3zlv9usa0afkmn2m7ccty5ykvk7gl5zk",
			"reward": {
				"denom": "ucmdx",
				"amount": "43827813"
			}
		},
		{
			"address": "comdex1xvkx9ht2eu6gtprsx9cr3cur0j76efgx2uh9jp",
			"reward": {
				"denom": "ucmdx",
				"amount": "99629"
			}
		},
		{
			"address": "comdex18xts44z9425d4eutvy45jtkgq55hk7wps6m0a9",
			"reward": {
				"denom": "ucmdx",
				"amount": "4023250"
			}
		},
		{
			"address": "comdex1f0m8va6j99rrc88maqxj05d69v8zfuk2h7qcp7",
			"reward": {
				"denom": "ucmdx",
				"amount": "13063808"
			}
		},
		{
			"address": "comdex1thfntksw0d35n2tkr0k8v54fr8wxtxwxwhedyj",
			"reward": {
				"denom": "ucmdx",
				"amount": "3442387"
			}
		},
		{
			"address": "comdex1vxvmljtx0eyt4uk8ywuasvx8j3l2e8llgnm0ys",
			"reward": {
				"denom": "ucmdx",
				"amount": "11924221"
			}
		},
		{
			"address": "comdex1nkcq0d0ph4uqfndw0tpmfau5kjsyf2ad8c8d4f",
			"reward": {
				"denom": "ucmdx",
				"amount": "36824387"
			}
		},
		{
			"address": "comdex1kxayenudz73u8qpmzsn5vnpshzw7urv6fu86ah",
			"reward": {
				"denom": "ucmdx",
				"amount": "819358"
			}
		},
		{
			"address": "comdex1699txsl87agrzxhk9wp2af2gruf5sps90fqpj9",
			"reward": {
				"denom": "ucmdx",
				"amount": "133771"
			}
		},
		{
			"address": "comdex17xfpkp90waa9e3743gye6e0kemxr720ezuy72s",
			"reward": {
				"denom": "ucmdx",
				"amount": "2710264"
			}
		},
		{
			"address": "comdex1q6le7sh8ya0utytzcvhz92kt3v9j2vwmw504qx",
			"reward": {
				"denom": "ucmdx",
				"amount": "105173"
			}
		},
		{
			"address": "comdex1934qwv8eamfpufjucv74wpunw9tduftuzzr7yg",
			"reward": {
				"denom": "ucmdx",
				"amount": "9722"
			}
		},
		{
			"address": "comdex1x2gtwryw06a6q35wprnu22n048sqjy375829hw",
			"reward": {
				"denom": "ucmdx",
				"amount": "1471680"
			}
		},
		{
			"address": "comdex1v0y32nzcpv99edjeyu0usfp73s9k5ga73mp4lv",
			"reward": {
				"denom": "ucmdx",
				"amount": "127909"
			}
		},
		{
			"address": "comdex1w5ey4nf64mzrjzuec69338dtpr3rl7w9zd6jxl",
			"reward": {
				"denom": "ucmdx",
				"amount": "108572"
			}
		},
		{
			"address": "comdex1sdxk2z7jyw0ckt7qek6cpqge68zelfs7pcdyuk",
			"reward": {
				"denom": "ucmdx",
				"amount": "764361"
			}
		},
		{
			"address": "comdex1ss4krazunlfnc8munks4uvsa5c42vlwrdlqux8",
			"reward": {
				"denom": "ucmdx",
				"amount": "6921375"
			}
		},
		{
			"address": "comdex137w8w3cszp6u2u8rfjk5a54ycq0fktpgdxse0k",
			"reward": {
				"denom": "ucmdx",
				"amount": "52687116"
			}
		},
		{
			"address": "comdex1k0zdx92k6d60vszdhm45hk4tpn06f6kgcv0wgs",
			"reward": {
				"denom": "ucmdx",
				"amount": "0"
			}
		},
		{
			"address": "comdex1cf6q3enhv2qhjcc9lv0w3knkctvw7g2eje6ku5",
			"reward": {
				"denom": "ucmdx",
				"amount": "14984244"
			}
		},
		{
			"address": "comdex16jv0fay4mrak6lw2rau3d593l99hv8jq5j6m8p",
			"reward": {
				"denom": "ucmdx",
				"amount": "4337"
			}
		},
		{
			"address": "comdex17yzaplnaswn75n0sx3z5lwhf4wzv692fme94r6",
			"reward": {
				"denom": "ucmdx",
				"amount": "5168692"
			}
		},
		{
			"address": "comdex17k0mksh2atxjtcxdgtrpae6zt4vhuevlq5jtd8",
			"reward": {
				"denom": "ucmdx",
				"amount": "5684480"
			}
		},
		{
			"address": "comdex1yftk8j7zq7gnd5ukhvvcsdyk4dej8e5dg5mj0v",
			"reward": {
				"denom": "ucmdx",
				"amount": "776540"
			}
		},
		{
			"address": "comdex18q85cqsuu97hpu7vn54433k86cxv29k2t4xmff",
			"reward": {
				"denom": "ucmdx",
				"amount": "429457"
			}
		},
		{
			"address": "comdex1t3qp52hf5p4julm52ftssqxdhpvyh9c04vnzk3",
			"reward": {
				"denom": "ucmdx",
				"amount": "39502"
			}
		},
		{
			"address": "comdex1sdyqpt0y29eg2fzse85cd7r0klahqufuxsfc5a",
			"reward": {
				"denom": "ucmdx",
				"amount": "673878"
			}
		},
		{
			"address": "comdex1qmhdz3dlgys30q7e3tvv6gr9r9su3j9et90nj8",
			"reward": {
				"denom": "ucmdx",
				"amount": "18143980"
			}
		},
		{
			"address": "comdex127ju3sedue92nrnv28g8dqsrarlgx3we83ykts",
			"reward": {
				"denom": "ucmdx",
				"amount": "41475252"
			}
		},
		{
			"address": "comdex10qey6aqjf5wlkglp0ljjg3v6h8ws26pmwuzfp2",
			"reward": {
				"denom": "ucmdx",
				"amount": "1729265"
			}
		},
		{
			"address": "comdex13uxtw430xyr2x387j4d3jkwzs0c03qpn6wvhrv",
			"reward": {
				"denom": "ucmdx",
				"amount": "329582"
			}
		},
		{
			"address": "comdex1hrdv45z737fj3vxhskheyl4222x9g5ckfs6lfz",
			"reward": {
				"denom": "ucmdx",
				"amount": "121737231"
			}
		},
		{
			"address": "comdex122atfyf5f6utduzqtn93wu29k935cdzdw3na0v",
			"reward": {
				"denom": "ucmdx",
				"amount": "7417385"
			}
		},
		{
			"address": "comdex1062pqrt6z6w0nsp6t25ekxmgxx9dl03jurxhjl",
			"reward": {
				"denom": "ucmdx",
				"amount": "12821112"
			}
		},
		{
			"address": "comdex1pyzxvfa7f8gr5x20n8kl0lu78jczuevnxw78vy",
			"reward": {
				"denom": "ucmdx",
				"amount": "48616274"
			}
		},
		{
			"address": "comdex1zr6tc009mam808yxf6xz8tq5aq39mptr3smtef",
			"reward": {
				"denom": "ucmdx",
				"amount": "1404401"
			}
		},
		{
			"address": "comdex1fxafs2mfgwzyjn42h7njcn9qng0nskwm8my0ak",
			"reward": {
				"denom": "ucmdx",
				"amount": "1340528"
			}
		},
		{
			"address": "comdex1h03pu8xlrtx40sqnmdnrvxn0m4cjwpy4khyyth",
			"reward": {
				"denom": "ucmdx",
				"amount": "13940814"
			}
		},
		{
			"address": "comdex1eykjz9ekzht9635d3gr9pq9nd4ctgu2974n88f",
			"reward": {
				"denom": "ucmdx",
				"amount": "31127647"
			}
		},
		{
			"address": "comdex1mxz09unulpudwjcc5px6fklknne79xgezvry2y",
			"reward": {
				"denom": "ucmdx",
				"amount": "3427302"
			}
		},
		{
			"address": "comdex1xxxwefdamtsg5kqaufyw8900deptkjmcy0f74n",
			"reward": {
				"denom": "ucmdx",
				"amount": "3714"
			}
		},
		{
			"address": "comdex1x6ghffw4fghzmn3u9xcrjehmanjdrtxd4fn8xl",
			"reward": {
				"denom": "ucmdx",
				"amount": "16789004"
			}
		},
		{
			"address": "comdex18qs6jenjswr72gu8wxcfqsgtcszkz038r85aar",
			"reward": {
				"denom": "ucmdx",
				"amount": "9558"
			}
		},
		{
			"address": "comdex1j96vu4dyf8qvzstn50dgwc2032t7ny7mgn32hq",
			"reward": {
				"denom": "ucmdx",
				"amount": "49579"
			}
		},
		{
			"address": "comdex1n0arcudyx3gtt829qjn2yl3a9mllnrvqyra9ql",
			"reward": {
				"denom": "ucmdx",
				"amount": "1335911"
			}
		},
		{
			"address": "comdex1m78t4rygjmp05sww2x93ectvu66s69t43kjzll",
			"reward": {
				"denom": "ucmdx",
				"amount": "10095"
			}
		},
		{
			"address": "comdex1ph65sftlddjcrvwl4lf4cpycmhhxpnu27t6c2c",
			"reward": {
				"denom": "ucmdx",
				"amount": "35080862"
			}
		},
		{
			"address": "comdex1z0fkt377zq6prpf58lsf404dj35ds55m956hez",
			"reward": {
				"denom": "ucmdx",
				"amount": "494913"
			}
		},
		{
			"address": "comdex123kc3cfvh2epdughwhv9r98y8hl4l3qh7rn6n2",
			"reward": {
				"denom": "ucmdx",
				"amount": "267860"
			}
		},
		{
			"address": "comdex1vzg2x804s9580m9zxq9k9nzek02vq2e28rrt0h",
			"reward": {
				"denom": "ucmdx",
				"amount": "1586083"
			}
		},
		{
			"address": "comdex1w0xf9d28aelpf7fvcan05gs36n7npjey6vf60w",
			"reward": {
				"denom": "ucmdx",
				"amount": "5519601"
			}
		},
		{
			"address": "comdex1w5nmf0hjxumew3w09s9gxutaxmrcm3n0dzlnsm",
			"reward": {
				"denom": "ucmdx",
				"amount": "3298450"
			}
		},
		{
			"address": "comdex16mxvcxrq8tkucvaftcm4g7guqfdkx9je42l9km",
			"reward": {
				"denom": "ucmdx",
				"amount": "128567"
			}
		},
		{
			"address": "comdex1aus5nuynfjr9xq4d9rrgnhhy8mma5mgry0zdvr",
			"reward": {
				"denom": "ucmdx",
				"amount": "4146342"
			}
		},
		{
			"address": "comdex1qqp2aydslhpx4emvqdhsyn8ztrltd4ze9hpgnq",
			"reward": {
				"denom": "ucmdx",
				"amount": "1363"
			}
		},
		{
			"address": "comdex1qqp75uvzj098umvff8v4ada2zk7gnk6xtse2a4",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1qqyl24rxge02cgkqnq4p2340s28kdd89rhzdws",
			"reward": {
				"denom": "ucmdx",
				"amount": "2869"
			}
		},
		{
			"address": "comdex1qq9q07ppx7s88md8capqgkahnna6va9cxewhqm",
			"reward": {
				"denom": "ucmdx",
				"amount": "140"
			}
		},
		{
			"address": "comdex1qqff9szr0def4x9wgt8neujwd7hah9tuesxmc3",
			"reward": {
				"denom": "ucmdx",
				"amount": "5665"
			}
		},
		{
			"address": "comdex1qq29rk7u5fvx4ms02q3u36qexc4p6zmxqz40x0",
			"reward": {
				"denom": "ucmdx",
				"amount": "1173"
			}
		},
		{
			"address": "comdex1qq2j2hfwy0nepnlvttvnzywff95u5v7h9gmxtm",
			"reward": {
				"denom": "ucmdx",
				"amount": "1297970"
			}
		},
		{
			"address": "comdex1qqwwxcylfsncw45jpavh0s6065dnh49jshlzeq",
			"reward": {
				"denom": "ucmdx",
				"amount": "5515"
			}
		},
		{
			"address": "comdex1qqje6vfxluqqfyxjglea8920z2ay04cz3pmg6m",
			"reward": {
				"denom": "ucmdx",
				"amount": "1688"
			}
		},
		{
			"address": "comdex1qqn662s9unp6jn9u6p5eamvdg8ymumw44lh560",
			"reward": {
				"denom": "ucmdx",
				"amount": "7294"
			}
		},
		{
			"address": "comdex1qqnlsrpk5uhf6ktdq08g52tjeqmlmrsyzz6x9p",
			"reward": {
				"denom": "ucmdx",
				"amount": "34"
			}
		},
		{
			"address": "comdex1qqu9qkfvjzqx2u5v4tefuwjdflkadcfq2cnlq7",
			"reward": {
				"denom": "ucmdx",
				"amount": "2014"
			}
		},
		{
			"address": "comdex1qqax75yllyggh9dpfdp3lhh5euehmw7vljsfju",
			"reward": {
				"denom": "ucmdx",
				"amount": "25518"
			}
		},
		{
			"address": "comdex1qq7uz6rzytr3u5e2m6gulv9fpzkemhkyj4zh5n",
			"reward": {
				"denom": "ucmdx",
				"amount": "129282"
			}
		},
		{
			"address": "comdex1qpqyyegsr2rd5ycnqlxk5a2t6dpxxx3xfx6su7",
			"reward": {
				"denom": "ucmdx",
				"amount": "1090"
			}
		},
		{
			"address": "comdex1qpq22tpk0v8qs9e8lq377l850ga73zjsth698m",
			"reward": {
				"denom": "ucmdx",
				"amount": "7246"
			}
		},
		{
			"address": "comdex1qpqnajk7ryumrsff59flf3ypy2pufjkjtaxq5s",
			"reward": {
				"denom": "ucmdx",
				"amount": "129"
			}
		},
		{
			"address": "comdex1qpp5t6hm49mpzx39vqg2lfaja7xld4l343z8vm",
			"reward": {
				"denom": "ucmdx",
				"amount": "3662"
			}
		},
		{
			"address": "comdex1qpy00pkhxyh9pga470rsththcft5ry6qtz90tz",
			"reward": {
				"denom": "ucmdx",
				"amount": "28"
			}
		},
		{
			"address": "comdex1qp9q9ncpttzluft2kzyfytcswvfy6l286e2jz3",
			"reward": {
				"denom": "ucmdx",
				"amount": "7196"
			}
		},
		{
			"address": "comdex1qpxra9yk5dux0q0ca75umejd94ncajtkp3zsut",
			"reward": {
				"denom": "ucmdx",
				"amount": "605"
			}
		},
		{
			"address": "comdex1qpfhgrvyqdcqxjyvst4gqggwv8gykflh5pud99",
			"reward": {
				"denom": "ucmdx",
				"amount": "22"
			}
		},
		{
			"address": "comdex1qptp2hntwecfcyjwxyr8cngmjl9fmfc64p6t23",
			"reward": {
				"denom": "ucmdx",
				"amount": "88"
			}
		},
		{
			"address": "comdex1qpvem7qvpygc87e0aqtjqz6lx7r3z35c8y37yd",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1qpw7vtx2w0yqygee2r9xrtuz2v3n2kcr0ucqg5",
			"reward": {
				"denom": "ucmdx",
				"amount": "3785"
			}
		},
		{
			"address": "comdex1qp00ra9chgnn5q5y3d435cws6plwvj5a8scr34",
			"reward": {
				"denom": "ucmdx",
				"amount": "6882"
			}
		},
		{
			"address": "comdex1qp39re424d78puvywlafhsddx7zt6r28gp3yd3",
			"reward": {
				"denom": "ucmdx",
				"amount": "2877"
			}
		},
		{
			"address": "comdex1qpny9sqwwq6a6auhvq3p6wulh6jnmdffykp47w",
			"reward": {
				"denom": "ucmdx",
				"amount": "175064"
			}
		},
		{
			"address": "comdex1qp59k44d9338qqc2vx4cc8rm0tnfwwvh2t70fk",
			"reward": {
				"denom": "ucmdx",
				"amount": "2648"
			}
		},
		{
			"address": "comdex1qp47ummal6vlmwa9a6yznwhpdhfdea5r9xapwx",
			"reward": {
				"denom": "ucmdx",
				"amount": "11265"
			}
		},
		{
			"address": "comdex1qpkgjd4amzjtfr3480esa82sykafpddhrke6ng",
			"reward": {
				"denom": "ucmdx",
				"amount": "1790"
			}
		},
		{
			"address": "comdex1qpktdl30qqcr5n9wkgkh8py8a25hfzk3dgj6c0",
			"reward": {
				"denom": "ucmdx",
				"amount": "126487"
			}
		},
		{
			"address": "comdex1qpkjncgldyumu3ftjh52lyufycmpxrq5qnncuu",
			"reward": {
				"denom": "ucmdx",
				"amount": "3516"
			}
		},
		{
			"address": "comdex1qphzltgp27frpp9h8s4chjlcxw0lpwy62aarl5",
			"reward": {
				"denom": "ucmdx",
				"amount": "5374"
			}
		},
		{
			"address": "comdex1qphy2vce2w4evtaw6r465cehxmacz87wq7p74t",
			"reward": {
				"denom": "ucmdx",
				"amount": "1734"
			}
		},
		{
			"address": "comdex1qphcrjpzx636ajc4ku4hmgm7vgv8z2vhghswlx",
			"reward": {
				"denom": "ucmdx",
				"amount": "4240466"
			}
		},
		{
			"address": "comdex1qphem9arpjavjynu72arhp8qy8fmna6rva6h03",
			"reward": {
				"denom": "ucmdx",
				"amount": "353"
			}
		},
		{
			"address": "comdex1qpeeygy9xh6uz8t7cadgf30q64vtdku3tf4lgz",
			"reward": {
				"denom": "ucmdx",
				"amount": "317"
			}
		},
		{
			"address": "comdex1qp749jf5nf8cgs374prymzsm8pl8xjgaxx8n0g",
			"reward": {
				"denom": "ucmdx",
				"amount": "17262"
			}
		},
		{
			"address": "comdex1qzyjzhmav3rehqw52vl6um5x36y2anvhf3r9yj",
			"reward": {
				"denom": "ucmdx",
				"amount": "28410"
			}
		},
		{
			"address": "comdex1qz9ml20945genaszrxvk7zxch6cn8jc24svr9p",
			"reward": {
				"denom": "ucmdx",
				"amount": "2454"
			}
		},
		{
			"address": "comdex1qzg9y2swv5z3mg72f9hjdlr7fyjddm8lfxmakr",
			"reward": {
				"denom": "ucmdx",
				"amount": "8228"
			}
		},
		{
			"address": "comdex1qzggkmzvegegdc7nchcnf20j85xmymnvz99psy",
			"reward": {
				"denom": "ucmdx",
				"amount": "173"
			}
		},
		{
			"address": "comdex1qz2lp0m5w24dpr8r8wuj265k0fsk5zpn4qlmh5",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1qztzc3pze6d3w5cpat4ysqhv46gxqf0hwsfjys",
			"reward": {
				"denom": "ucmdx",
				"amount": "107854"
			}
		},
		{
			"address": "comdex1qzvrc782mglck54ryefg4rp25jh45rvjds42nh",
			"reward": {
				"denom": "ucmdx",
				"amount": "74"
			}
		},
		{
			"address": "comdex1qzjr7jek5e6dk3t87dp62p8gy39hcuxtxmcuvq",
			"reward": {
				"denom": "ucmdx",
				"amount": "198"
			}
		},
		{
			"address": "comdex1qzk9v44tf4ygy7rsc4m63xqmyftsp9pu6ce62p",
			"reward": {
				"denom": "ucmdx",
				"amount": "195"
			}
		},
		{
			"address": "comdex1qzkden27can3yfg5de73v7g9cjqxsk32ln29n3",
			"reward": {
				"denom": "ucmdx",
				"amount": "1756"
			}
		},
		{
			"address": "comdex1qz6r4chdr8j9qadh7mswctc6def25grs3z034z",
			"reward": {
				"denom": "ucmdx",
				"amount": "90019"
			}
		},
		{
			"address": "comdex1qz6kn9s8hzvhy9kkaaskvqzjdpl8hfhhumey9v",
			"reward": {
				"denom": "ucmdx",
				"amount": "814"
			}
		},
		{
			"address": "comdex1qzms6cn5p8muv8a6470h20gyccdc8fv96zw06s",
			"reward": {
				"denom": "ucmdx",
				"amount": "1644"
			}
		},
		{
			"address": "comdex1qzukx3ktd3xcnq2j6pshmudkgyke0s5fem94q2",
			"reward": {
				"denom": "ucmdx",
				"amount": "35238"
			}
		},
		{
			"address": "comdex1qzul26v2tnqkjpcd2qxqdamr40zr7aj45qlej7",
			"reward": {
				"denom": "ucmdx",
				"amount": "809"
			}
		},
		{
			"address": "comdex1qrpd94sfdk4nqyk6frhw4qkysqwjqxgyc5clgf",
			"reward": {
				"denom": "ucmdx",
				"amount": "291"
			}
		},
		{
			"address": "comdex1qrr9238a4ka6q8cd0vkczu62425szekcqrx78m",
			"reward": {
				"denom": "ucmdx",
				"amount": "0"
			}
		},
		{
			"address": "comdex1qrytx9vstv2wa3ns0xcuafnqw7w44pkpc537ty",
			"reward": {
				"denom": "ucmdx",
				"amount": "5096"
			}
		},
		{
			"address": "comdex1qr9x6mdfmktv0rcj2t9cgpje32mfsj7um6j470",
			"reward": {
				"denom": "ucmdx",
				"amount": "1876"
			}
		},
		{
			"address": "comdex1qrx26f88x6j5477p08uva4lk6xdvfkra0t3rt4",
			"reward": {
				"denom": "ucmdx",
				"amount": "163"
			}
		},
		{
			"address": "comdex1qrfk2x4aygkensgjsskhtmj8txyn52uy9c6pn2",
			"reward": {
				"denom": "ucmdx",
				"amount": "1429"
			}
		},
		{
			"address": "comdex1qr2zspe3ua3z20r2getxfu278sy292ysrzfe74",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1qrsrnaxudcx9vy4uhvu35fl4m6zhmvgccyt240",
			"reward": {
				"denom": "ucmdx",
				"amount": "13942"
			}
		},
		{
			"address": "comdex1qr3zqsuvuz6ey62j5uh8wy4ms7ey3uz0tmtk2s",
			"reward": {
				"denom": "ucmdx",
				"amount": "24012"
			}
		},
		{
			"address": "comdex1qr3swll95kcam9gfwny0nv0u4cn46sga23zek5",
			"reward": {
				"denom": "ucmdx",
				"amount": "1286"
			}
		},
		{
			"address": "comdex1qr5jcmgx4pm6uzmcmwnh4cq5nmafmnx2zvsrj8",
			"reward": {
				"denom": "ucmdx",
				"amount": "392"
			}
		},
		{
			"address": "comdex1qr554ryfnkvk2h2p2a5zw5jv4u53rrx0ehwn7k",
			"reward": {
				"denom": "ucmdx",
				"amount": "7617"
			}
		},
		{
			"address": "comdex1qrk3546fp8lkgtvehkhjpjkkszr3mhpucrumcl",
			"reward": {
				"denom": "ucmdx",
				"amount": "20386"
			}
		},
		{
			"address": "comdex1qrm2jnqm93ejyzlmsvzhu6s4q047clnjhy7626",
			"reward": {
				"denom": "ucmdx",
				"amount": "169"
			}
		},
		{
			"address": "comdex1qr7eckxg5zgsp40sdwr9qj2arfx5ftgpml9he2",
			"reward": {
				"denom": "ucmdx",
				"amount": "124552"
			}
		},
		{
			"address": "comdex1qyqpyd5kaz4f3np33f696zkn2vfd87ap9l03cw",
			"reward": {
				"denom": "ucmdx",
				"amount": "6432"
			}
		},
		{
			"address": "comdex1qyxqsjle8afcdhjleqgv6kzaw4qx7kh77ncnxw",
			"reward": {
				"denom": "ucmdx",
				"amount": "1424"
			}
		},
		{
			"address": "comdex1qyfllhxpatxmzs3qnrd86cvls35prtvjssd29y",
			"reward": {
				"denom": "ucmdx",
				"amount": "8688"
			}
		},
		{
			"address": "comdex1qyt55k85cegzr9c6z4j8pkuptaygda9hvr3psx",
			"reward": {
				"denom": "ucmdx",
				"amount": "1645"
			}
		},
		{
			"address": "comdex1qywetm2u0rtjzv7mnue9q6hwqkhqasz77q58jl",
			"reward": {
				"denom": "ucmdx",
				"amount": "1426"
			}
		},
		{
			"address": "comdex1qywmwaskh9l4sakxhudqnell4x4c3x0uartsta",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1qy4aw8vtlddzt70d7k0rm3mqlpd3zh94g63hft",
			"reward": {
				"denom": "ucmdx",
				"amount": "11390"
			}
		},
		{
			"address": "comdex1qykpgnx9w9ezw45qqee3tywzzzn2zzlyzlaze3",
			"reward": {
				"denom": "ucmdx",
				"amount": "1438"
			}
		},
		{
			"address": "comdex1qycfaxw2x0hu4cvc2m2023fp7lapaes7fuhdp0",
			"reward": {
				"denom": "ucmdx",
				"amount": "132"
			}
		},
		{
			"address": "comdex1qycn95enjukm2zluvmlg5nkj7mpnskgnmmvu86",
			"reward": {
				"denom": "ucmdx",
				"amount": "468"
			}
		},
		{
			"address": "comdex1qyeczt8yewuc4jd5ahdfjsa05teksh9evqq7uz",
			"reward": {
				"denom": "ucmdx",
				"amount": "13"
			}
		},
		{
			"address": "comdex1qymeqmky7en30mfud9nw7rx03x29y9d0y9228d",
			"reward": {
				"denom": "ucmdx",
				"amount": "8067"
			}
		},
		{
			"address": "comdex1qyar8pxn0q5yela70mv4r77nwuptw6eljl7gpt",
			"reward": {
				"denom": "ucmdx",
				"amount": "131"
			}
		},
		{
			"address": "comdex1qyajtkn39qjhesupnyech925z2ehsft6rfg0y3",
			"reward": {
				"denom": "ucmdx",
				"amount": "186"
			}
		},
		{
			"address": "comdex1qy728jfmkcqqd7p2e4cjffywcgrfxwqn3wg588",
			"reward": {
				"denom": "ucmdx",
				"amount": "1930"
			}
		},
		{
			"address": "comdex1qylmu92gtyqrlr506nnxwfrvsyx6rxp4fjnwrg",
			"reward": {
				"denom": "ucmdx",
				"amount": "251011"
			}
		},
		{
			"address": "comdex1q9pc0lgasemk5yyg6ht4tj0w4e6ehtusc43y79",
			"reward": {
				"denom": "ucmdx",
				"amount": "16789"
			}
		},
		{
			"address": "comdex1q99rnv2tq4hwgmn8aukhasz4yt6zlft0m7ph9c",
			"reward": {
				"denom": "ucmdx",
				"amount": "546"
			}
		},
		{
			"address": "comdex1q9x4e97v9qzqkw9q9hc7n9drzagrj9mqr03ek6",
			"reward": {
				"denom": "ucmdx",
				"amount": "4082"
			}
		},
		{
			"address": "comdex1q9852p2xyrxdst70kwt84ng970hq5qr8nfllf7",
			"reward": {
				"denom": "ucmdx",
				"amount": "50"
			}
		},
		{
			"address": "comdex1q98u2u4y403fvs66400wfs8pr8kl37mc0x3keq",
			"reward": {
				"denom": "ucmdx",
				"amount": "1430"
			}
		},
		{
			"address": "comdex1q9fu9a5qnx9hzc744tue27as68hyy6ktxtvjxu",
			"reward": {
				"denom": "ucmdx",
				"amount": "268"
			}
		},
		{
			"address": "comdex1q9v58lzm2c4se429wz0kzjflmxmhhupns4g3hs",
			"reward": {
				"denom": "ucmdx",
				"amount": "140"
			}
		},
		{
			"address": "comdex1q9dda6vydlyrm7rl3qjw7eaygkf6yg7txshf6r",
			"reward": {
				"denom": "ucmdx",
				"amount": "54773"
			}
		},
		{
			"address": "comdex1q9w6jdl437lgu0lasfyxt9j3mr5fxwrhvxmfpv",
			"reward": {
				"denom": "ucmdx",
				"amount": "12089"
			}
		},
		{
			"address": "comdex1q90ze5xxczlptkqrr8u2u8z7mrtxmf8emelskt",
			"reward": {
				"denom": "ucmdx",
				"amount": "94"
			}
		},
		{
			"address": "comdex1q9spuklw7zdnh7a7gu3a422ma0lerr343z0u33",
			"reward": {
				"denom": "ucmdx",
				"amount": "56082"
			}
		},
		{
			"address": "comdex1q9snn84jfrd9ge8t46kdcggpe58dua82awn8z0",
			"reward": {
				"denom": "ucmdx",
				"amount": "51358"
			}
		},
		{
			"address": "comdex1q9nfsdllk86vyl0jcfjdy0qv6sq8rexw878uuw",
			"reward": {
				"denom": "ucmdx",
				"amount": "1429"
			}
		},
		{
			"address": "comdex1q9cp0e6tgd7xwt6jnf7g0nqdp6uk9us45587ru",
			"reward": {
				"denom": "ucmdx",
				"amount": "289"
			}
		},
		{
			"address": "comdex1q9el8t3czw0ygj4d5nvzwme4vasjjw8z6s8h7d",
			"reward": {
				"denom": "ucmdx",
				"amount": "55797"
			}
		},
		{
			"address": "comdex1q9u3nud7zyjtjgs36m8z8vd50wvs8paevlvfsx",
			"reward": {
				"denom": "ucmdx",
				"amount": "4295"
			}
		},
		{
			"address": "comdex1q9al7fp8ugg7d4jty6n78jkh0cjk9yn9hpcst9",
			"reward": {
				"denom": "ucmdx",
				"amount": "2013"
			}
		},
		{
			"address": "comdex1q97d8nwteaeqkwvxt8z9nluxtax9vnc65547zs",
			"reward": {
				"denom": "ucmdx",
				"amount": "3344"
			}
		},
		{
			"address": "comdex1qxqeze4y4w778685zsunsxe3meuauzcg2e2zqd",
			"reward": {
				"denom": "ucmdx",
				"amount": "179"
			}
		},
		{
			"address": "comdex1qxzlt45zurtrhju9sua0ec6nrp3qfjylw5jr59",
			"reward": {
				"denom": "ucmdx",
				"amount": "32021"
			}
		},
		{
			"address": "comdex1qx9eu2sfnklmf0kmggz40elt88klrgewz9wttw",
			"reward": {
				"denom": "ucmdx",
				"amount": "145"
			}
		},
		{
			"address": "comdex1qxgjzyysg4yvwvynv0psq769rnn6qym5tls899",
			"reward": {
				"denom": "ucmdx",
				"amount": "7186"
			}
		},
		{
			"address": "comdex1qxg4kk59yqz49qkeppk0klnseecpgj67m9sx2h",
			"reward": {
				"denom": "ucmdx",
				"amount": "491"
			}
		},
		{
			"address": "comdex1qxt4ym2nymln0xhdn420pewqjnalr3ejvxkk2q",
			"reward": {
				"denom": "ucmdx",
				"amount": "3164"
			}
		},
		{
			"address": "comdex1qxvqhff8xkuyd8l6cjwlpefammtsku7m5afvmq",
			"reward": {
				"denom": "ucmdx",
				"amount": "2718"
			}
		},
		{
			"address": "comdex1qxd6ge0dxa8fe8qzc4punqramqcq390t4t4wgp",
			"reward": {
				"denom": "ucmdx",
				"amount": "551"
			}
		},
		{
			"address": "comdex1qxnqc24rtrk7njcauqjhqe5zqmpt05re6usjmv",
			"reward": {
				"denom": "ucmdx",
				"amount": "106273"
			}
		},
		{
			"address": "comdex1qxkysemk47lwxeafgz6dadlrzg0wn2yzur7zqx",
			"reward": {
				"denom": "ucmdx",
				"amount": "28636"
			}
		},
		{
			"address": "comdex1qxes636ffjlfxw5gl8gztdgaetswce7ap3kqme",
			"reward": {
				"denom": "ucmdx",
				"amount": "6138"
			}
		},
		{
			"address": "comdex1q8qpw9cww02urxmlkk2nqvmn83h9etgsk3vqxq",
			"reward": {
				"denom": "ucmdx",
				"amount": "250"
			}
		},
		{
			"address": "comdex1q8q4k72sartzqxxvkw0k96kj3apjtdk0a03gjm",
			"reward": {
				"denom": "ucmdx",
				"amount": "4965"
			}
		},
		{
			"address": "comdex1q8qkghenqacqeq5phgmny42pfuxmmuhlfh0gge",
			"reward": {
				"denom": "ucmdx",
				"amount": "7654"
			}
		},
		{
			"address": "comdex1q8qu54hjkxz0fp8qm0fa2rqefz8zcur9cjm8zj",
			"reward": {
				"denom": "ucmdx",
				"amount": "1199"
			}
		},
		{
			"address": "comdex1q8pj2g7d2rcp9qq383mq8fquqn6sp9ky4qgm3k",
			"reward": {
				"denom": "ucmdx",
				"amount": "168"
			}
		},
		{
			"address": "comdex1q8ryuy79ev69d97j8cm8ughw4uela5valmh6hq",
			"reward": {
				"denom": "ucmdx",
				"amount": "732"
			}
		},
		{
			"address": "comdex1q8rh8hjc7fqs9devyaee6wmdp8p4ndanl8zcf0",
			"reward": {
				"denom": "ucmdx",
				"amount": "303"
			}
		},
		{
			"address": "comdex1q8x8f7g9j5x4xa239h5gyfnkdw5yjmze9y2dcl",
			"reward": {
				"denom": "ucmdx",
				"amount": "1252"
			}
		},
		{
			"address": "comdex1q8ga0sgnpn0f792wqr6qnyv4yr35hvldlscspg",
			"reward": {
				"denom": "ucmdx",
				"amount": "90"
			}
		},
		{
			"address": "comdex1q8fy46kgzjfdv7xu7glmlkx0tv9r3azkv7q9r6",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1q8fwp62lesuauhncceu222rsqe2fwx390vcfl0",
			"reward": {
				"denom": "ucmdx",
				"amount": "292"
			}
		},
		{
			"address": "comdex1q82z2eeq7cecw5w9wt2lg798m70uzl36sryrjl",
			"reward": {
				"denom": "ucmdx",
				"amount": "391"
			}
		},
		{
			"address": "comdex1q8tsp0jc2pjksdpp96q5t2la2a4djev2nqryvu",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1q8dzxvvfesw476njhvzcze72p49636umc3rp3k",
			"reward": {
				"denom": "ucmdx",
				"amount": "1139"
			}
		},
		{
			"address": "comdex1q8wshsexd4p2nelpq662su93dz62rc4vca43np",
			"reward": {
				"denom": "ucmdx",
				"amount": "19564"
			}
		},
		{
			"address": "comdex1q83s0wehh8tvn9fc9eqfun8v844csmy80d7y06",
			"reward": {
				"denom": "ucmdx",
				"amount": "57958"
			}
		},
		{
			"address": "comdex1q83h557e77t9peqr9cwc848qwkr887h4kqkn2z",
			"reward": {
				"denom": "ucmdx",
				"amount": "890"
			}
		},
		{
			"address": "comdex1q85c764zhzp4ka4t5htrmacvxwvhak258dx567",
			"reward": {
				"denom": "ucmdx",
				"amount": "20912"
			}
		},
		{
			"address": "comdex1q85eskmkjejtahw8zyvfsu44swvjumj5vmkcln",
			"reward": {
				"denom": "ucmdx",
				"amount": "5078"
			}
		},
		{
			"address": "comdex1q8568awdzem73m0pptfvmga00m37mr5ctvpcvj",
			"reward": {
				"denom": "ucmdx",
				"amount": "1389"
			}
		},
		{
			"address": "comdex1q85ak6agjvsuak600hqpkdtulmywjs6ypv33vg",
			"reward": {
				"denom": "ucmdx",
				"amount": "7053"
			}
		},
		{
			"address": "comdex1q85lyu2feytxpjhdqs4l2whahhz99eyehcfrsh",
			"reward": {
				"denom": "ucmdx",
				"amount": "12006"
			}
		},
		{
			"address": "comdex1q8cxk5pktrhzzcy2axmnejd26phrtwdlpm425d",
			"reward": {
				"denom": "ucmdx",
				"amount": "163"
			}
		},
		{
			"address": "comdex1q86973k9dt66l090yg2vw24ajj98pry5ehpj65",
			"reward": {
				"denom": "ucmdx",
				"amount": "5290"
			}
		},
		{
			"address": "comdex1q8l0d9rk9fwe2ul822adpmk0h94gfz2kmttpqs",
			"reward": {
				"denom": "ucmdx",
				"amount": "594"
			}
		},
		{
			"address": "comdex1qgz2q9tlw3m6e4hkx0h86lxu0qym0mgn8wkxhl",
			"reward": {
				"denom": "ucmdx",
				"amount": "143002"
			}
		},
		{
			"address": "comdex1qgr5wawt5hqc9cjg2plxa7jltd7h69nqqmxwxs",
			"reward": {
				"denom": "ucmdx",
				"amount": "14215"
			}
		},
		{
			"address": "comdex1qgrcu9p46nmxddeuqe6saqhccdzqr6xw752q2m",
			"reward": {
				"denom": "ucmdx",
				"amount": "17944"
			}
		},
		{
			"address": "comdex1qgrmajh0900xe98993hvuf83f2z75g2y7ajc57",
			"reward": {
				"denom": "ucmdx",
				"amount": "3589"
			}
		},
		{
			"address": "comdex1qg8y2hqa37766c3p2ahdujex7y6ux80j7d2zh6",
			"reward": {
				"denom": "ucmdx",
				"amount": "87"
			}
		},
		{
			"address": "comdex1qggq0f0pkk45sqmwme2crqsvt4stuymt28tgw8",
			"reward": {
				"denom": "ucmdx",
				"amount": "44094"
			}
		},
		{
			"address": "comdex1qgf2r9jzap77e4mqqpsm4devef5yetllv3pu0g",
			"reward": {
				"denom": "ucmdx",
				"amount": "11828"
			}
		},
		{
			"address": "comdex1qg2vgahmvrx0p64rl9hrmrg7mptrgnfzwe8syd",
			"reward": {
				"denom": "ucmdx",
				"amount": "130092"
			}
		},
		{
			"address": "comdex1qgdhdhu47hgx8j7qkdcsneu5e2saypmlplme9j",
			"reward": {
				"denom": "ucmdx",
				"amount": "77630"
			}
		},
		{
			"address": "comdex1qg04faedykl4da9p66wk7937zkauv766lwz6sm",
			"reward": {
				"denom": "ucmdx",
				"amount": "13962"
			}
		},
		{
			"address": "comdex1qgsud6zs7raufdcvxmcsqda2ccw0syepj33ezk",
			"reward": {
				"denom": "ucmdx",
				"amount": "14622"
			}
		},
		{
			"address": "comdex1qgs7r0t2djqzvf8u33mhcw7hqq2jfc5htd6a5q",
			"reward": {
				"denom": "ucmdx",
				"amount": "359"
			}
		},
		{
			"address": "comdex1qg3qqrwp7gqyr8s8dj9anfvd2gx4c0w5cyl3y6",
			"reward": {
				"denom": "ucmdx",
				"amount": "5442"
			}
		},
		{
			"address": "comdex1qgk9huqtrx2u7ry8rwuu6lrs0hn52swjee8he2",
			"reward": {
				"denom": "ucmdx",
				"amount": "28110"
			}
		},
		{
			"address": "comdex1qghufceaw35ja57s0j34mgqzlhulz5nm5u4rta",
			"reward": {
				"denom": "ucmdx",
				"amount": "14360"
			}
		},
		{
			"address": "comdex1qgcrc6lgfycffdre09fwp4t35lkcags6xm94vc",
			"reward": {
				"denom": "ucmdx",
				"amount": "1502"
			}
		},
		{
			"address": "comdex1qge7yt898l4xvt72rtlcrjn5pyyhe2q0tqjz4x",
			"reward": {
				"denom": "ucmdx",
				"amount": "3604"
			}
		},
		{
			"address": "comdex1qgme8vlq4ly8tcye6xdgxnz4khzq9es0acc5zp",
			"reward": {
				"denom": "ucmdx",
				"amount": "195"
			}
		},
		{
			"address": "comdex1qgudpxeuh7hnkj7vyvfp06jtcjq0pajuqpf72p",
			"reward": {
				"denom": "ucmdx",
				"amount": "8854"
			}
		},
		{
			"address": "comdex1qgu4l692mx3m9lg0uxnt8f62zxpzn0tspwqq9t",
			"reward": {
				"denom": "ucmdx",
				"amount": "117"
			}
		},
		{
			"address": "comdex1qgaylp3nmqqtmhtsk7s6egw6qyq68meeuc0724",
			"reward": {
				"denom": "ucmdx",
				"amount": "129"
			}
		},
		{
			"address": "comdex1qfpdpa0j702sajr3454uf0crtshts9n8tyhs75",
			"reward": {
				"denom": "ucmdx",
				"amount": "14474"
			}
		},
		{
			"address": "comdex1qfzpwe7nwytfrknqtxgvz0s7t6rn07t3k9hw38",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1qf9nep5hzyv5yywxcr0wjd86lytsc34fu9cy6x",
			"reward": {
				"denom": "ucmdx",
				"amount": "144"
			}
		},
		{
			"address": "comdex1qf8v3ts73uhhqlckpjjt38hxe7338u778t3fc8",
			"reward": {
				"denom": "ucmdx",
				"amount": "74091"
			}
		},
		{
			"address": "comdex1qfv7khxqlpgtx897c0d6mrf26s9ma3kcnntr2m",
			"reward": {
				"denom": "ucmdx",
				"amount": "4033"
			}
		},
		{
			"address": "comdex1qf0480jfhx3ksr3gwu4eruk9qxm25jvp8eze5k",
			"reward": {
				"denom": "ucmdx",
				"amount": "1404"
			}
		},
		{
			"address": "comdex1qf0kxkk60qrcj5qa7v7t439249qfkcd40wa8pa",
			"reward": {
				"denom": "ucmdx",
				"amount": "16672"
			}
		},
		{
			"address": "comdex1qf3fycpz5ahw4dpl275u3y5ngz290lh22ajvn2",
			"reward": {
				"denom": "ucmdx",
				"amount": "38926"
			}
		},
		{
			"address": "comdex1qf36e6wmq9h4twhdvs6pyq9qcaeu7ye0pvve7p",
			"reward": {
				"denom": "ucmdx",
				"amount": "18919"
			}
		},
		{
			"address": "comdex1qfjhfwnur3vagt0hs7fqmmw349hmtffyxw9c0a",
			"reward": {
				"denom": "ucmdx",
				"amount": "180"
			}
		},
		{
			"address": "comdex1qf4afymv9evpfunhtj0fpwllt0k5hlzlzhpsnc",
			"reward": {
				"denom": "ucmdx",
				"amount": "48459"
			}
		},
		{
			"address": "comdex1qfkxws2mgj5hsmttymf6xjagzuvu3945plhy7d",
			"reward": {
				"denom": "ucmdx",
				"amount": "75"
			}
		},
		{
			"address": "comdex1qfh7v60frxc887zkym8zej9nr58werqs2h2u7z",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1qf73ju06s0u0lyk9jppv9j3vfl5q9h4dqhad39",
			"reward": {
				"denom": "ucmdx",
				"amount": "1062"
			}
		},
		{
			"address": "comdex1qflrtckeplg5teh4vj5wekn2d5fxkz4pg3khqr",
			"reward": {
				"denom": "ucmdx",
				"amount": "5282"
			}
		},
		{
			"address": "comdex1q2glzcjnxu8r3usulev7uqr0kq2md6k08q9tn0",
			"reward": {
				"denom": "ucmdx",
				"amount": "13718"
			}
		},
		{
			"address": "comdex1q2glkvf2adkvn4ufu0syvw35x3s2yyrezva8d0",
			"reward": {
				"denom": "ucmdx",
				"amount": "1975"
			}
		},
		{
			"address": "comdex1q2t96wltse9s8ygl2a4e92tccmahrve6qxyc74",
			"reward": {
				"denom": "ucmdx",
				"amount": "2591779"
			}
		},
		{
			"address": "comdex1q2dh3ruvktptpp0usj057vt0zffndmuftuwy5l",
			"reward": {
				"denom": "ucmdx",
				"amount": "75"
			}
		},
		{
			"address": "comdex1q2wdqkjjsay2qc8ze2ef2d32jssm4rrx87cfrc",
			"reward": {
				"denom": "ucmdx",
				"amount": "417"
			}
		},
		{
			"address": "comdex1q2wketg7fcpvk2gmv2yd7qq8f44myupkp65jlm",
			"reward": {
				"denom": "ucmdx",
				"amount": "846"
			}
		},
		{
			"address": "comdex1q2sm52c86w586hqduv6aqrre8ylpnhukuxm8mv",
			"reward": {
				"denom": "ucmdx",
				"amount": "1787"
			}
		},
		{
			"address": "comdex1q2nc3f9pv0xsdrts56uu6lsrjxjhkhht6e55cz",
			"reward": {
				"denom": "ucmdx",
				"amount": "45"
			}
		},
		{
			"address": "comdex1q2cxn0ud3utagllzj2zmk20ss7kwkzp4udn973",
			"reward": {
				"denom": "ucmdx",
				"amount": "1243"
			}
		},
		{
			"address": "comdex1q2cuye4r58c6aerg5e8qh3acs03nhe43fvrt5y",
			"reward": {
				"denom": "ucmdx",
				"amount": "63431"
			}
		},
		{
			"address": "comdex1q26rkt6zj7aau4qspu3w5ps3ps2f2qcd7vemeg",
			"reward": {
				"denom": "ucmdx",
				"amount": "779"
			}
		},
		{
			"address": "comdex1q26y8aqrta0equ94petqyxw4wzkrps73pxj4hg",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1q26avdk4xgqrz0gurled4mw0f4c4ar3ses0djj",
			"reward": {
				"denom": "ucmdx",
				"amount": "14"
			}
		},
		{
			"address": "comdex1q2m399slayvqhqpm36lfcvelw93r9taxypc7dn",
			"reward": {
				"denom": "ucmdx",
				"amount": "166"
			}
		},
		{
			"address": "comdex1q2m4607vdf3rwce9pq8ld8vq3hshtzsc39zv3w",
			"reward": {
				"denom": "ucmdx",
				"amount": "6332"
			}
		},
		{
			"address": "comdex1q2lp9etvcd4kfre6hnxxzrxtwneq5ekc3fzsj8",
			"reward": {
				"denom": "ucmdx",
				"amount": "282"
			}
		},
		{
			"address": "comdex1q2lemtjr8rg9528tjennzgzle2fa35kfuxds23",
			"reward": {
				"denom": "ucmdx",
				"amount": "404"
			}
		},
		{
			"address": "comdex1q2llca8gsuk22j46cu2796hsvmr55nlejt5c2z",
			"reward": {
				"denom": "ucmdx",
				"amount": "1327"
			}
		},
		{
			"address": "comdex1qtqeatrcltvk03y07u0fu6ajc4uty3ewr7zh3d",
			"reward": {
				"denom": "ucmdx",
				"amount": "271"
			}
		},
		{
			"address": "comdex1qtz8r7z80ecn4kd5khzvefmj7cg25zd40449ep",
			"reward": {
				"denom": "ucmdx",
				"amount": "50024"
			}
		},
		{
			"address": "comdex1qtylxwt7vu2dvy0xzcp66lzsttxtj4c8mxdqg0",
			"reward": {
				"denom": "ucmdx",
				"amount": "149"
			}
		},
		{
			"address": "comdex1qtyldmyj8r53v39kq4hgvxnrtcsstxkrvhfhm6",
			"reward": {
				"denom": "ucmdx",
				"amount": "31835"
			}
		},
		{
			"address": "comdex1qtxvrwewz4fxkn829c4lulyzx63jupgpqxlg2w",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1qtg904z3ewk64dqrurwfy3z5ml2yy7lf5y7vx4",
			"reward": {
				"denom": "ucmdx",
				"amount": "39611"
			}
		},
		{
			"address": "comdex1qttmfljc3qcerypthjtxw4a8675xcf9hcug6ed",
			"reward": {
				"denom": "ucmdx",
				"amount": "54123"
			}
		},
		{
			"address": "comdex1qt0nlqrfkfq8de2vqecx5s832y3naj8wx7mcew",
			"reward": {
				"denom": "ucmdx",
				"amount": "1816"
			}
		},
		{
			"address": "comdex1qtsptcq0yt6rz0pa6sly6v0e8geayea58l3kfw",
			"reward": {
				"denom": "ucmdx",
				"amount": "44339"
			}
		},
		{
			"address": "comdex1qtj9z8w2r8p2nc2dasd965zy3tlpe59n8a0kdm",
			"reward": {
				"denom": "ucmdx",
				"amount": "7212"
			}
		},
		{
			"address": "comdex1qt5jnkqnrn7mx9qhq5jgrghrqxjrsxrq42s7ad",
			"reward": {
				"denom": "ucmdx",
				"amount": "693"
			}
		},
		{
			"address": "comdex1qt4jy9v3mhvex4prsnvsef368jymk88meppj56",
			"reward": {
				"denom": "ucmdx",
				"amount": "712596"
			}
		},
		{
			"address": "comdex1qtk9vgdz6t20vzcrsawnvvx7lqgtpz89uxvne4",
			"reward": {
				"denom": "ucmdx",
				"amount": "22526"
			}
		},
		{
			"address": "comdex1qtc2a0rjrjcgg2lrkkpaynyc637ju44qh7ag3z",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1qtepnzqmdjt9e95ezywhvg7f0pqekjwvwskc76",
			"reward": {
				"denom": "ucmdx",
				"amount": "141"
			}
		},
		{
			"address": "comdex1qte2kxnq9zhfty44yym8tnwl7636s43urj40w0",
			"reward": {
				"denom": "ucmdx",
				"amount": "14490"
			}
		},
		{
			"address": "comdex1qtut300t4az2xfd6e4aywvd8u80szzwkyzd60r",
			"reward": {
				"denom": "ucmdx",
				"amount": "46419"
			}
		},
		{
			"address": "comdex1qt7p0sflljh57v2ls53we98lktyd3uyv7px3fs",
			"reward": {
				"denom": "ucmdx",
				"amount": "172"
			}
		},
		{
			"address": "comdex1qtl3p007vd9zz7ws8xj9vyjt5x7elg7chpaqdy",
			"reward": {
				"denom": "ucmdx",
				"amount": "422745"
			}
		},
		{
			"address": "comdex1qvq39e9mjaqner30swwrq6vf59huyxetjjwp6k",
			"reward": {
				"denom": "ucmdx",
				"amount": "20887"
			}
		},
		{
			"address": "comdex1qvqcwlf2sqpz5hnthm6pqxymqxyevmlseqj7mq",
			"reward": {
				"denom": "ucmdx",
				"amount": "1214"
			}
		},
		{
			"address": "comdex1qvq6s6pnk6k3qxtlkcydqrru8uvu98665kr3pw",
			"reward": {
				"denom": "ucmdx",
				"amount": "3759"
			}
		},
		{
			"address": "comdex1qvqu8fp4x8ctxumet87fwv7qkpvtq2s329dyck",
			"reward": {
				"denom": "ucmdx",
				"amount": "616"
			}
		},
		{
			"address": "comdex1qv9x5meurqh8yjflcw7zlax5400422zcpnqrsp",
			"reward": {
				"denom": "ucmdx",
				"amount": "126569"
			}
		},
		{
			"address": "comdex1qv9js24cdq3jsahgt7jwjz4uv5qsc6areu3nht",
			"reward": {
				"denom": "ucmdx",
				"amount": "176"
			}
		},
		{
			"address": "comdex1qvgrrvernfdx2qf5fztn82swzxr37gr3nafa6f",
			"reward": {
				"denom": "ucmdx",
				"amount": "1232"
			}
		},
		{
			"address": "comdex1qvtpawav3dfft7y7q08jd5ykkrtpw9jsm4ytmk",
			"reward": {
				"denom": "ucmdx",
				"amount": "874"
			}
		},
		{
			"address": "comdex1qvtluu0mp7n23rwgw5r4ja3ert7fue8pahn8ed",
			"reward": {
				"denom": "ucmdx",
				"amount": "6742"
			}
		},
		{
			"address": "comdex1qvdke4vd9qxjx0fxhxkvzh93a6k38zrf9temyf",
			"reward": {
				"denom": "ucmdx",
				"amount": "129802"
			}
		},
		{
			"address": "comdex1qv3cv46u5klmwaj2d04tya26en9qwg0tr3qa6p",
			"reward": {
				"denom": "ucmdx",
				"amount": "12364"
			}
		},
		{
			"address": "comdex1qv3at230lvgjfjqstyswfn3nkrayu5s333nky6",
			"reward": {
				"denom": "ucmdx",
				"amount": "1725"
			}
		},
		{
			"address": "comdex1qvn6ghe68l4g0k7s25rujr6yfpyrm6h3fgvvr0",
			"reward": {
				"denom": "ucmdx",
				"amount": "125"
			}
		},
		{
			"address": "comdex1qv50app8hhacfpdktggfrrl047cflyxzmsr79v",
			"reward": {
				"denom": "ucmdx",
				"amount": "6835"
			}
		},
		{
			"address": "comdex1qvkvcar33wnwkuv56yurnu4lsfkkdcz6snlr8z",
			"reward": {
				"denom": "ucmdx",
				"amount": "3206"
			}
		},
		{
			"address": "comdex1qve6yzh6f4r8sqflfmxq7ywxh5cflzexxuqghh",
			"reward": {
				"denom": "ucmdx",
				"amount": "158673"
			}
		},
		{
			"address": "comdex1qv6ee9vtu5j00dl4ky9tyqps34cftclwhqwkce",
			"reward": {
				"denom": "ucmdx",
				"amount": "86"
			}
		},
		{
			"address": "comdex1qvmpf9y7g8hqu68f5mq5w9fkdkjkv5jw7pdqmf",
			"reward": {
				"denom": "ucmdx",
				"amount": "7284"
			}
		},
		{
			"address": "comdex1qdqpcnn27cxndwr3z7h2m3h7g65zu7fnfcfnre",
			"reward": {
				"denom": "ucmdx",
				"amount": "10499"
			}
		},
		{
			"address": "comdex1qdzgugsx5zueucmj54ugr7rtllkksh3733ep2t",
			"reward": {
				"denom": "ucmdx",
				"amount": "1233"
			}
		},
		{
			"address": "comdex1qdrwa4d4pu8da9ydkmwr53dfypr30mkfewm3fj",
			"reward": {
				"denom": "ucmdx",
				"amount": "32591"
			}
		},
		{
			"address": "comdex1qd9v3rejvj0lggcqt8mna7v5w45qj7asw5qymf",
			"reward": {
				"denom": "ucmdx",
				"amount": "1747"
			}
		},
		{
			"address": "comdex1qdx4qw80x6mcdf3m9kmt9rrsw7dp4zjwdamj55",
			"reward": {
				"denom": "ucmdx",
				"amount": "142"
			}
		},
		{
			"address": "comdex1qd2apugvtly67s6yfzpet5cgpxaa0nqxpv9v6h",
			"reward": {
				"denom": "ucmdx",
				"amount": "4514"
			}
		},
		{
			"address": "comdex1qdw67wjf84qvul2g5dkpv09k3tnmdt4t7x8yeg",
			"reward": {
				"denom": "ucmdx",
				"amount": "36513"
			}
		},
		{
			"address": "comdex1qd5cvs4tel4a5zzgq2qgsmvyfzp7ytpxdare5g",
			"reward": {
				"denom": "ucmdx",
				"amount": "7780"
			}
		},
		{
			"address": "comdex1qdkvk964l7hmac4jqfpv4eqjwch9uy2xm5tzkm",
			"reward": {
				"denom": "ucmdx",
				"amount": "28"
			}
		},
		{
			"address": "comdex1qdh9gzkzeaxfdlg9cqlq805ksa8y43vt3gd098",
			"reward": {
				"denom": "ucmdx",
				"amount": "13993"
			}
		},
		{
			"address": "comdex1qdhsdudftpkcdgfpsydp0dtcc37w533m7v7848",
			"reward": {
				"denom": "ucmdx",
				"amount": "12320"
			}
		},
		{
			"address": "comdex1qdexvf9yae7wxj8gny8shzylhxrdcyl8uh9nxc",
			"reward": {
				"denom": "ucmdx",
				"amount": "397"
			}
		},
		{
			"address": "comdex1qd6rxzlgqzgz82cnjdc3jh3ju9yssay8sgwgr5",
			"reward": {
				"denom": "ucmdx",
				"amount": "16765"
			}
		},
		{
			"address": "comdex1qd6ynzaxcmuhgw8j42d2k865nky36fzmx4vanj",
			"reward": {
				"denom": "ucmdx",
				"amount": "1461"
			}
		},
		{
			"address": "comdex1qdu59su8skkrj9jzl8az696fd3v3f7fvjukarx",
			"reward": {
				"denom": "ucmdx",
				"amount": "143"
			}
		},
		{
			"address": "comdex1qdag0rqapegtr4w8udu69mgf4cp50cuf9v9dnm",
			"reward": {
				"denom": "ucmdx",
				"amount": "1732"
			}
		},
		{
			"address": "comdex1qda3nxh86ht7q3lmrtsh2jmkw5hn05w5dn6vw6",
			"reward": {
				"denom": "ucmdx",
				"amount": "1984"
			}
		},
		{
			"address": "comdex1qdl0puccs33ckna5d5gtcduk2hmse2xn4g0g8c",
			"reward": {
				"denom": "ucmdx",
				"amount": "207"
			}
		},
		{
			"address": "comdex1qdlug9864zj5mh76p7n70f2pcn56jfgtjfglvt",
			"reward": {
				"denom": "ucmdx",
				"amount": "1235"
			}
		},
		{
			"address": "comdex1qwgjjxcklma6kxwjn4zs9s49mce8mk6uj6fvxm",
			"reward": {
				"denom": "ucmdx",
				"amount": "1747"
			}
		},
		{
			"address": "comdex1qw264x7dvqcfuxj7zcsj2j9gglpg4wj8ssrv0z",
			"reward": {
				"denom": "ucmdx",
				"amount": "1999"
			}
		},
		{
			"address": "comdex1qwtfg8smvgft0vsaya7fg4d646mpnw9z3q97y0",
			"reward": {
				"denom": "ucmdx",
				"amount": "38861"
			}
		},
		{
			"address": "comdex1qw39tlpa5akjrvtzm08e375affrhnpz5x96qr4",
			"reward": {
				"denom": "ucmdx",
				"amount": "1724"
			}
		},
		{
			"address": "comdex1qwnxm22hrt6jm444c8d2guw3zrll0alrpwxyeq",
			"reward": {
				"denom": "ucmdx",
				"amount": "176"
			}
		},
		{
			"address": "comdex1qw590hrrv84uypn0amkp7zkc3xl76g4vz0p46t",
			"reward": {
				"denom": "ucmdx",
				"amount": "2026"
			}
		},
		{
			"address": "comdex1qwcc3nunfpkhln38sjaj3t5yw3szakgtct53yn",
			"reward": {
				"denom": "ucmdx",
				"amount": "645"
			}
		},
		{
			"address": "comdex1qwc70zrau90c5tj8pp0zzreklgvfem8vglnjs5",
			"reward": {
				"denom": "ucmdx",
				"amount": "5458"
			}
		},
		{
			"address": "comdex1qw6dlh5ucegm3hjp8pcmcma57xh5m8m455e2ap",
			"reward": {
				"denom": "ucmdx",
				"amount": "291"
			}
		},
		{
			"address": "comdex1qwmpmw5r6edwxry3f50wxr7zhxeqder4wg283n",
			"reward": {
				"denom": "ucmdx",
				"amount": "2950"
			}
		},
		{
			"address": "comdex1qwm9swxxtxfjzz49ctnx023cpcz53njd90mel6",
			"reward": {
				"denom": "ucmdx",
				"amount": "123396"
			}
		},
		{
			"address": "comdex1qwm86e67myvf282sua85hawn6nq8jeja8ejlc5",
			"reward": {
				"denom": "ucmdx",
				"amount": "36"
			}
		},
		{
			"address": "comdex1qwasamh38qgch2crlxad7l2agyax6lcqvumlwl",
			"reward": {
				"denom": "ucmdx",
				"amount": "19507"
			}
		},
		{
			"address": "comdex1qw7t5te03mvcdcvvuvnx8lg45hq2dh88fljqec",
			"reward": {
				"denom": "ucmdx",
				"amount": "137"
			}
		},
		{
			"address": "comdex1qw70xh24jcn86a4a8qgnxlvxw3wzggjlx3kqmf",
			"reward": {
				"denom": "ucmdx",
				"amount": "4011"
			}
		},
		{
			"address": "comdex1qw76444u658xf99g9vkr64n05kq0gl09q6g546",
			"reward": {
				"denom": "ucmdx",
				"amount": "21583"
			}
		},
		{
			"address": "comdex1qwl0h82l3edmvlx3afryag64ln2255t56q8enr",
			"reward": {
				"denom": "ucmdx",
				"amount": "192"
			}
		},
		{
			"address": "comdex1qwlkgfsxkr7v64a8w2m0j0h73znhw83uhgymd5",
			"reward": {
				"denom": "ucmdx",
				"amount": "173"
			}
		},
		{
			"address": "comdex1q0pja02396sjvyrpeypguhwf8rfg79tq2ghrcm",
			"reward": {
				"denom": "ucmdx",
				"amount": "3007"
			}
		},
		{
			"address": "comdex1q0r5ell5rfj7t9ndpa7975zz82ghfsh85wpeuw",
			"reward": {
				"denom": "ucmdx",
				"amount": "210"
			}
		},
		{
			"address": "comdex1q0fh6l29w7e6vj6x0n0hd5trp02nc8u9td0c70",
			"reward": {
				"denom": "ucmdx",
				"amount": "100892"
			}
		},
		{
			"address": "comdex1q0tfewa60enuu8utyjs9w7um58njhd7fjpeukx",
			"reward": {
				"denom": "ucmdx",
				"amount": "16811"
			}
		},
		{
			"address": "comdex1q0wycdla9unye0tkrxwtfh23r6atczhj0dqyrc",
			"reward": {
				"denom": "ucmdx",
				"amount": "176"
			}
		},
		{
			"address": "comdex1q0whsalla09e9hyctuw6kylkzx9a779sdsw9ur",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1q0wc3ps7ehwayy372slg2eqf6fzgqdrnk3u8yh",
			"reward": {
				"denom": "ucmdx",
				"amount": "2856"
			}
		},
		{
			"address": "comdex1q0k3zstemzx0k5x7v7jz8zwcu7qqyzysvd7jnr",
			"reward": {
				"denom": "ucmdx",
				"amount": "1641"
			}
		},
		{
			"address": "comdex1q0u7zqyxnsvrh7wt22mtk3hlj39rg2sdvescux",
			"reward": {
				"denom": "ucmdx",
				"amount": "6875"
			}
		},
		{
			"address": "comdex1qsqjxsnxq3gnav7n60kzlrrvacq7zrnc65myyn",
			"reward": {
				"denom": "ucmdx",
				"amount": "471"
			}
		},
		{
			"address": "comdex1qsqj8dr0tdn855rkm0mcyvyuw58cp4lkht5cyc",
			"reward": {
				"denom": "ucmdx",
				"amount": "1740"
			}
		},
		{
			"address": "comdex1qsp85ghzfq5af43veaa6vm8xglun0fx5ckvfh8",
			"reward": {
				"denom": "ucmdx",
				"amount": "7171"
			}
		},
		{
			"address": "comdex1qsr8tgynj69dmngrt34uaxtkuy0ylqhs9tmftv",
			"reward": {
				"denom": "ucmdx",
				"amount": "21247"
			}
		},
		{
			"address": "comdex1qsxhsqpgydnwy5kxz9akw2zkkuy5ppmar2wz9g",
			"reward": {
				"denom": "ucmdx",
				"amount": "10578"
			}
		},
		{
			"address": "comdex1qsfd5vkynrs2glpzmmj2xlnl7wk20k3z0sgf6p",
			"reward": {
				"denom": "ucmdx",
				"amount": "2"
			}
		},
		{
			"address": "comdex1qs2t36aa40gp0ttljq6pc3vxctnrrp5d4427k5",
			"reward": {
				"denom": "ucmdx",
				"amount": "203"
			}
		},
		{
			"address": "comdex1qswp0luvdvuhussx5nuzryv7wk3m5pjp6rx205",
			"reward": {
				"denom": "ucmdx",
				"amount": "16672"
			}
		},
		{
			"address": "comdex1qsjhtgcnm88asrkj86dsy3hev9mhtrn5yrshgq",
			"reward": {
				"denom": "ucmdx",
				"amount": "4116"
			}
		},
		{
			"address": "comdex1qs4cv7zmnrwqjx2mxmh4vwv7y9cu7esfnnk5tr",
			"reward": {
				"denom": "ucmdx",
				"amount": "61191"
			}
		},
		{
			"address": "comdex1qsemc4a3rklk5qxh3p4attkkav3pdlal2zd6ht",
			"reward": {
				"denom": "ucmdx",
				"amount": "200108"
			}
		},
		{
			"address": "comdex1qseu7wmcmk3tdy082z7j0sreh05gklr052qxds",
			"reward": {
				"denom": "ucmdx",
				"amount": "63681"
			}
		},
		{
			"address": "comdex1qs654gylw7mzp7nfqs0l7xq7xg847wry6ww47a",
			"reward": {
				"denom": "ucmdx",
				"amount": "6041"
			}
		},
		{
			"address": "comdex1qs647n8wdt68s0h56tvmwa2zn084d7z2urjkek",
			"reward": {
				"denom": "ucmdx",
				"amount": "179"
			}
		},
		{
			"address": "comdex1qsull9s00frk6l039g6jcklcxxnnehlt4fy97t",
			"reward": {
				"denom": "ucmdx",
				"amount": "724"
			}
		},
		{
			"address": "comdex1qs76t9gy78sk5hxnpceque5a9rc6gf7er3damu",
			"reward": {
				"denom": "ucmdx",
				"amount": "192"
			}
		},
		{
			"address": "comdex1qslxfc8ksjxm5shpug7nxt4mgdr270hf2qpkhv",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1q3z9j333a8x9uewemk92qcvlnuqna9j4605795",
			"reward": {
				"denom": "ucmdx",
				"amount": "1793"
			}
		},
		{
			"address": "comdex1q3yku7vztxdgwakhk07c9tt0cuy5r5d9tf6lvg",
			"reward": {
				"denom": "ucmdx",
				"amount": "315"
			}
		},
		{
			"address": "comdex1q3gfatr5g9u4h7k4485g372rtvg8xtjuw4704c",
			"reward": {
				"denom": "ucmdx",
				"amount": "1788"
			}
		},
		{
			"address": "comdex1q3vfmmu8w3vpar95x4yarwv43epx6pf3nachn0",
			"reward": {
				"denom": "ucmdx",
				"amount": "203"
			}
		},
		{
			"address": "comdex1q3sp4yxk3epf92sgs2l9tz84543kedpu5jvl3d",
			"reward": {
				"denom": "ucmdx",
				"amount": "89"
			}
		},
		{
			"address": "comdex1q3smvv69w0yphzvm9xyr4htuvemz2uyrqnmx2q",
			"reward": {
				"denom": "ucmdx",
				"amount": "363"
			}
		},
		{
			"address": "comdex1q3h26qgy7n0th9jr33adl07rmjr4ph4ymkagxp",
			"reward": {
				"denom": "ucmdx",
				"amount": "26277"
			}
		},
		{
			"address": "comdex1q3expyp88m3u73pvj5vknwcaa9s8620kajg6ql",
			"reward": {
				"denom": "ucmdx",
				"amount": "2551"
			}
		},
		{
			"address": "comdex1q3e7ecyjajwshj5atp8uckm69vhmudqkp0879h",
			"reward": {
				"denom": "ucmdx",
				"amount": "2211"
			}
		},
		{
			"address": "comdex1q3ue4jwx7e9zezwlnt0z39jmfwq3wkmzvw6vlh",
			"reward": {
				"denom": "ucmdx",
				"amount": "1607"
			}
		},
		{
			"address": "comdex1q3a63z295y2rgnjglhw2hangmj5n0gmlyenq7v",
			"reward": {
				"denom": "ucmdx",
				"amount": "768"
			}
		},
		{
			"address": "comdex1q3lzf7a0hlefqqmp5kjxdf66s8s8ma74fc40wv",
			"reward": {
				"denom": "ucmdx",
				"amount": "151448"
			}
		},
		{
			"address": "comdex1qjq3c5gx6az9ymmkda28yjvk86a04vussyewvr",
			"reward": {
				"denom": "ucmdx",
				"amount": "8830"
			}
		},
		{
			"address": "comdex1qjq57elgv4wkjpy5mf053l6n8v45txzjca999s",
			"reward": {
				"denom": "ucmdx",
				"amount": "35984"
			}
		},
		{
			"address": "comdex1qjpv7t2lnv4kxv5maa94vhfq4q6jkcspd23d7l",
			"reward": {
				"denom": "ucmdx",
				"amount": "9609"
			}
		},
		{
			"address": "comdex1qjzzu0pd9c47k6kvl8f3mfj0he67xwkfq3dvn9",
			"reward": {
				"denom": "ucmdx",
				"amount": "14"
			}
		},
		{
			"address": "comdex1qjzg059uxkesfzpwcup3yzje7sdk9at72v7k2p",
			"reward": {
				"denom": "ucmdx",
				"amount": "1762"
			}
		},
		{
			"address": "comdex1qjzu84a9sgsa0fsfk2uxatwuv90f2tgrpxrkcw",
			"reward": {
				"denom": "ucmdx",
				"amount": "15754"
			}
		},
		{
			"address": "comdex1qjrzy8grgc044pxxykhfkhvzvhr7c32uf0wn9t",
			"reward": {
				"denom": "ucmdx",
				"amount": "18549"
			}
		},
		{
			"address": "comdex1qj9qs3mppdjf67jhndp4trnr9wwqvyummqttax",
			"reward": {
				"denom": "ucmdx",
				"amount": "65808"
			}
		},
		{
			"address": "comdex1qjfsdrz7u9udnhkug06dcl9cumtldkkd2agnrm",
			"reward": {
				"denom": "ucmdx",
				"amount": "21961"
			}
		},
		{
			"address": "comdex1qjf6mu52gxed2el237ydrtprsky0aqkcsmeylf",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1qj2787n872jn8t2zhreac27lh0hlpq6zzygv3d",
			"reward": {
				"denom": "ucmdx",
				"amount": "1024"
			}
		},
		{
			"address": "comdex1qj0d8jzlqkueldzmvz9yxdzvzheacytlmax9yc",
			"reward": {
				"denom": "ucmdx",
				"amount": "1736"
			}
		},
		{
			"address": "comdex1qj0atwdrthht7ah8k6vvpnvz9jget6s5zv8mcv",
			"reward": {
				"denom": "ucmdx",
				"amount": "306"
			}
		},
		{
			"address": "comdex1qjs5ya42auyw4nfyc7y62ky5qjv9w7g35fy57j",
			"reward": {
				"denom": "ucmdx",
				"amount": "712"
			}
		},
		{
			"address": "comdex1qjs6kass0t24uqpcf96z792nyzjvhg6hm2gurr",
			"reward": {
				"denom": "ucmdx",
				"amount": "21998"
			}
		},
		{
			"address": "comdex1qj3rzpu4eknfcufzvtxvdmelc3hu2fa2faw0m9",
			"reward": {
				"denom": "ucmdx",
				"amount": "7320"
			}
		},
		{
			"address": "comdex1qj3yznwq8c6n55ma27wqwxjxzl5ykjwl47g05k",
			"reward": {
				"denom": "ucmdx",
				"amount": "166"
			}
		},
		{
			"address": "comdex1qjhysxkw2qpp6q77d3eeqn9dnmfw3dma080n3x",
			"reward": {
				"denom": "ucmdx",
				"amount": "144"
			}
		},
		{
			"address": "comdex1qjh7a82uhty0ashvqtf5nszcrcptsdphx8q0a7",
			"reward": {
				"denom": "ucmdx",
				"amount": "3985"
			}
		},
		{
			"address": "comdex1qj7nm3gxxypwqp5lvrjg27jc7tlhlys7r7shmh",
			"reward": {
				"denom": "ucmdx",
				"amount": "195759"
			}
		},
		{
			"address": "comdex1qj77duwrudg9tl6ekksnt3ehnpnwvf8qdxxeeu",
			"reward": {
				"denom": "ucmdx",
				"amount": "7714"
			}
		},
		{
			"address": "comdex1qjlp6j7tg8f7p827s8yz04hv5rgu7wcm3wklww",
			"reward": {
				"denom": "ucmdx",
				"amount": "39"
			}
		},
		{
			"address": "comdex1qnqxa8sny7f87devd6ys7etyet4lp0gr3s9kzu",
			"reward": {
				"denom": "ucmdx",
				"amount": "7096"
			}
		},
		{
			"address": "comdex1qnpx3y9jrx6z439r55aekmelkwe3fp42f3kavs",
			"reward": {
				"denom": "ucmdx",
				"amount": "649"
			}
		},
		{
			"address": "comdex1qnryfwn6hw02xuqhatw7d4tahxpwmns5angmd6",
			"reward": {
				"denom": "ucmdx",
				"amount": "15"
			}
		},
		{
			"address": "comdex1qnr3g3zjegemn88fh7nfn4rrw3k90rdd2hcyqc",
			"reward": {
				"denom": "ucmdx",
				"amount": "97"
			}
		},
		{
			"address": "comdex1qnylt6xhhk6mgx88gs0q4dfjpt80xq5sp285q9",
			"reward": {
				"denom": "ucmdx",
				"amount": "148"
			}
		},
		{
			"address": "comdex1qnxrqxm6cfkjmk4x8tcusfusdyann2e0z0373t",
			"reward": {
				"denom": "ucmdx",
				"amount": "4078"
			}
		},
		{
			"address": "comdex1qnfppzd4hjvyge0h740ftwvugjkakwaf6mtcck",
			"reward": {
				"denom": "ucmdx",
				"amount": "13394"
			}
		},
		{
			"address": "comdex1qnw8ncancruullzrh4aes2ksh2cfx7tcvamx26",
			"reward": {
				"denom": "ucmdx",
				"amount": "750"
			}
		},
		{
			"address": "comdex1qn082mqg5uzu0z2ddf6ess705dy6k30sh70gc3",
			"reward": {
				"denom": "ucmdx",
				"amount": "6165"
			}
		},
		{
			"address": "comdex1qnc3zk9qmstucrt68c7pp8fhurxgedyaphu8sw",
			"reward": {
				"denom": "ucmdx",
				"amount": "2472"
			}
		},
		{
			"address": "comdex1qnm7vjlp06set0r9dxmlxlssz4l7l9qmy4s0c0",
			"reward": {
				"denom": "ucmdx",
				"amount": "3600"
			}
		},
		{
			"address": "comdex1q5p35mgk04vxw3s7y4caf93yku2mfygf34lztc",
			"reward": {
				"denom": "ucmdx",
				"amount": "5288"
			}
		},
		{
			"address": "comdex1q5zgptmpjuxfl3n0uhg9g4slwggnnwwt5rqj73",
			"reward": {
				"denom": "ucmdx",
				"amount": "891"
			}
		},
		{
			"address": "comdex1q5raruhvrtzc2e4xxd2ev8gv7k95paywhzy9q8",
			"reward": {
				"denom": "ucmdx",
				"amount": "6225"
			}
		},
		{
			"address": "comdex1q5ylux9egnn4guvqlqzgl6s5e8y5cekd2hkxzk",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1q5x63dwfq76cgqcry05rpmyf9ylr9jcqkuhpj3",
			"reward": {
				"denom": "ucmdx",
				"amount": "75490"
			}
		},
		{
			"address": "comdex1q5fqy26g9snvq88jek7wgqtxy4nxrauwfnrq7x",
			"reward": {
				"denom": "ucmdx",
				"amount": "6556"
			}
		},
		{
			"address": "comdex1q52t7djt7gvnaeg3mdrx7stq6lat0fcygrlty8",
			"reward": {
				"denom": "ucmdx",
				"amount": "15489"
			}
		},
		{
			"address": "comdex1q5v4twzqtje33qgm49dtfgmet5u7rh337xxv2g",
			"reward": {
				"denom": "ucmdx",
				"amount": "271"
			}
		},
		{
			"address": "comdex1q5dm00zxvnymy96gc49agpxvgqug3mfzuvwpn3",
			"reward": {
				"denom": "ucmdx",
				"amount": "143"
			}
		},
		{
			"address": "comdex1q500462h3f0m5xdq7uvp94kjd3h8hg93p8stzh",
			"reward": {
				"denom": "ucmdx",
				"amount": "1901"
			}
		},
		{
			"address": "comdex1q50hkqt3lx3d9lpsj523n0qj3epr9hfsaeaus8",
			"reward": {
				"denom": "ucmdx",
				"amount": "6378"
			}
		},
		{
			"address": "comdex1q558d3wxxz0zjct3vu5mfgp60lwvgyr7sgaedq",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1q55wnqg0cf0u0pg5futf9zrl8zqagretmxm8mg",
			"reward": {
				"denom": "ucmdx",
				"amount": "6821"
			}
		},
		{
			"address": "comdex1q5k6gsdkhylj77cmjay0ergzukdcldt5z2aaj4",
			"reward": {
				"denom": "ucmdx",
				"amount": "143"
			}
		},
		{
			"address": "comdex1q5hqhr6q5qglmweh6adlul4zjwfhclpj200nst",
			"reward": {
				"denom": "ucmdx",
				"amount": "4106"
			}
		},
		{
			"address": "comdex1q5hhpcd2sxggn3ay8ea9330affk3mtxhflyjun",
			"reward": {
				"denom": "ucmdx",
				"amount": "123"
			}
		},
		{
			"address": "comdex1q5ma4wnrc6vgyqjagxe9lmq5xnh5r8k6r52kqy",
			"reward": {
				"denom": "ucmdx",
				"amount": "15460"
			}
		},
		{
			"address": "comdex1q5utz24aghrxuzmue0tde7q30vukq6nwyaww74",
			"reward": {
				"denom": "ucmdx",
				"amount": "6842"
			}
		},
		{
			"address": "comdex1q5acjnezmf2gl2uqvwamzf22pdt0dt008mdvdk",
			"reward": {
				"denom": "ucmdx",
				"amount": "145"
			}
		},
		{
			"address": "comdex1q4z4a6uzr682py5zluj03rudssj5n90myf68lc",
			"reward": {
				"denom": "ucmdx",
				"amount": "526"
			}
		},
		{
			"address": "comdex1q4ysu4yl060g9pxp57nrv6q4t6vjzunlfu7ttm",
			"reward": {
				"denom": "ucmdx",
				"amount": "144111"
			}
		},
		{
			"address": "comdex1q4x6s0dd9fym7l563tdn6kw2vd56t8ws68qe7r",
			"reward": {
				"denom": "ucmdx",
				"amount": "114"
			}
		},
		{
			"address": "comdex1q429fu7xwrjy9uk78r9f96rh2qaqrjtgeetswn",
			"reward": {
				"denom": "ucmdx",
				"amount": "892"
			}
		},
		{
			"address": "comdex1q4tw3ejl9xfajyh43pkcuahhuaaylkurjh234w",
			"reward": {
				"denom": "ucmdx",
				"amount": "10"
			}
		},
		{
			"address": "comdex1q4vfp4cs5sws92eff0fgtvegpjdftm9r3a6m2h",
			"reward": {
				"denom": "ucmdx",
				"amount": "6163"
			}
		},
		{
			"address": "comdex1q4vjkvtj8ux4m0lvxkpk5j9gf3s3we6rp3ggqs",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1q4wpsytkq99zar49p0adkjzcad5c3rtlhrdtju",
			"reward": {
				"denom": "ucmdx",
				"amount": "2332"
			}
		},
		{
			"address": "comdex1q4wxylcshtumrdp84sx3r9g79gtex5qcs4twd8",
			"reward": {
				"denom": "ucmdx",
				"amount": "144"
			}
		},
		{
			"address": "comdex1q404w4hlgu9thaa2vq608t29pd0y9lecnnnmza",
			"reward": {
				"denom": "ucmdx",
				"amount": "1975"
			}
		},
		{
			"address": "comdex1q4jarumugqc4vvmq2cy6fjddwfje95v46p8y7f",
			"reward": {
				"denom": "ucmdx",
				"amount": "754"
			}
		},
		{
			"address": "comdex1q45ql4ajwycnlh0ruutk3qt0z2fre6uzym82zz",
			"reward": {
				"denom": "ucmdx",
				"amount": "176"
			}
		},
		{
			"address": "comdex1q442zg75xjpxam9595tj6d265uf5asr6el5n3x",
			"reward": {
				"denom": "ucmdx",
				"amount": "108751"
			}
		},
		{
			"address": "comdex1q4h6pvn7s2r9e8h4fad5wepk8j9mw3g4n8dhf6",
			"reward": {
				"denom": "ucmdx",
				"amount": "143"
			}
		},
		{
			"address": "comdex1q4ug8tl7ghz5j0q23dhuszv98he392h2z8h4qf",
			"reward": {
				"denom": "ucmdx",
				"amount": "38124"
			}
		},
		{
			"address": "comdex1q47lnyzdy2dvs33wpz2q90xcuxnqr3edfvl2wj",
			"reward": {
				"denom": "ucmdx",
				"amount": "2181"
			}
		},
		{
			"address": "comdex1q4lqj27y44m4kf30s9q5jggp93cqpfxsk6udty",
			"reward": {
				"denom": "ucmdx",
				"amount": "2015"
			}
		},
		{
			"address": "comdex1qkp8g37rsdf9xahgeap2uuj0d3sq4hmrkxvraj",
			"reward": {
				"denom": "ucmdx",
				"amount": "882"
			}
		},
		{
			"address": "comdex1qkzqdx4l2r6gf0nmw3n0jmv2dust2kc4v8kpy0",
			"reward": {
				"denom": "ucmdx",
				"amount": "645716"
			}
		},
		{
			"address": "comdex1qkyufhpmtvfzklrdcsjtw4yjwj7vcdr3ljrqd3",
			"reward": {
				"denom": "ucmdx",
				"amount": "1971"
			}
		},
		{
			"address": "comdex1qkfymnkre3js86kt75ryl6k58ekmcswg9yat3p",
			"reward": {
				"denom": "ucmdx",
				"amount": "2279"
			}
		},
		{
			"address": "comdex1qkfmakm3h97uafajfwx24lecafl5e8s5j9y82p",
			"reward": {
				"denom": "ucmdx",
				"amount": "4113"
			}
		},
		{
			"address": "comdex1qktg2es04r3q75jpd03kryasrdn450du3k8t2t",
			"reward": {
				"denom": "ucmdx",
				"amount": "3495"
			}
		},
		{
			"address": "comdex1qktml8av5a9tn0dvc38m99y20v3y50jnvwnl4m",
			"reward": {
				"denom": "ucmdx",
				"amount": "353"
			}
		},
		{
			"address": "comdex1qk0rtrnk0xtf2g99u06z09sax28a67n0lhe9f8",
			"reward": {
				"denom": "ucmdx",
				"amount": "13554"
			}
		},
		{
			"address": "comdex1qksy9vf9p7d76qh9tmhjvenv8qu8dyr2yhv4qa",
			"reward": {
				"denom": "ucmdx",
				"amount": "2144"
			}
		},
		{
			"address": "comdex1qksdaex4ws9rdakalcm2yl59uxfcxm5d3va65c",
			"reward": {
				"denom": "ucmdx",
				"amount": "1107"
			}
		},
		{
			"address": "comdex1qk32s5ndae8n8ucle7dermqv6gakxs8mxaeamt",
			"reward": {
				"denom": "ucmdx",
				"amount": "1760"
			}
		},
		{
			"address": "comdex1qk35qaj0qen70pz70q6f4uwrqw6dpkdhftcp6l",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1qknz80zsnft88ppt5vchrgnkl7hnfp732wtfjv",
			"reward": {
				"denom": "ucmdx",
				"amount": "13362"
			}
		},
		{
			"address": "comdex1qk4g4pv07gymxchkucjfluepm5q59v8ahm3c6a",
			"reward": {
				"denom": "ucmdx",
				"amount": "268"
			}
		},
		{
			"address": "comdex1qkhfl2x79f5u8deuu4656acmghq836ljtus9rf",
			"reward": {
				"denom": "ucmdx",
				"amount": "135407"
			}
		},
		{
			"address": "comdex1qkmlemn5rwlmah6dd8yzahw4ak3mr499p20p3y",
			"reward": {
				"denom": "ucmdx",
				"amount": "16290"
			}
		},
		{
			"address": "comdex1qkugghsmeyxutcmqc8tlqwwfq4nskyqw8unm23",
			"reward": {
				"denom": "ucmdx",
				"amount": "257284"
			}
		},
		{
			"address": "comdex1qk7zrahe8snnuevhj63m9gzev7ucxnjtf3xvjd",
			"reward": {
				"denom": "ucmdx",
				"amount": "1227"
			}
		},
		{
			"address": "comdex1qklsjdgrjsjkh6qvzfs8uecyphdn5ceewzpx3p",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1qhg39z3ez6a7vw6crclgx2fvk0y04wga2gsfxq",
			"reward": {
				"denom": "ucmdx",
				"amount": "1438"
			}
		},
		{
			"address": "comdex1qh25fcfkwh0mxcdh5563kgyylq77fqep77p74s",
			"reward": {
				"denom": "ucmdx",
				"amount": "37106"
			}
		},
		{
			"address": "comdex1qhv9xtyely0qalwlueuer6vanfgz5w0ulndp45",
			"reward": {
				"denom": "ucmdx",
				"amount": "561"
			}
		},
		{
			"address": "comdex1qhv8vyl4gmmz69ggz72ktxa2pc34hvsc32fw3k",
			"reward": {
				"denom": "ucmdx",
				"amount": "4104"
			}
		},
		{
			"address": "comdex1qhw97e8v8hwzem0eh2xd0d44q44mgtf837tlw3",
			"reward": {
				"denom": "ucmdx",
				"amount": "3263"
			}
		},
		{
			"address": "comdex1qhwgpaw6cv5zlv3znd4gnferhp6dkkwp4y4403",
			"reward": {
				"denom": "ucmdx",
				"amount": "109322"
			}
		},
		{
			"address": "comdex1qh3v4xsadl8f5s23t4ytxm7fcdt7dczzyfcenr",
			"reward": {
				"denom": "ucmdx",
				"amount": "40901"
			}
		},
		{
			"address": "comdex1qhhenw5ydq05tes7ardrxhu37zpu6ux5thkdk6",
			"reward": {
				"denom": "ucmdx",
				"amount": "1634"
			}
		},
		{
			"address": "comdex1qhc5d04lfqphtuh6dtme6dwwr4fvkrvng6zthz",
			"reward": {
				"denom": "ucmdx",
				"amount": "2336"
			}
		},
		{
			"address": "comdex1qhazunwq5wes6fxjdkvdcqenqqf9kdys3ewd09",
			"reward": {
				"denom": "ucmdx",
				"amount": "2270"
			}
		},
		{
			"address": "comdex1qcq6s3prgn2ny98dsd2g9yu4d3vnc7j70pd7mw",
			"reward": {
				"denom": "ucmdx",
				"amount": "180"
			}
		},
		{
			"address": "comdex1qczet66z8x3lrsvgmnc9c724w8m6nv5c3rzjdu",
			"reward": {
				"denom": "ucmdx",
				"amount": "1055"
			}
		},
		{
			"address": "comdex1qcrjhyp0x7ft5k0duz8gxx4cmwqaxv7ld0gscx",
			"reward": {
				"denom": "ucmdx",
				"amount": "7134"
			}
		},
		{
			"address": "comdex1qcr5cqha6nj5kk7lkaf4ecvcvp73drxcn50vhe",
			"reward": {
				"denom": "ucmdx",
				"amount": "471"
			}
		},
		{
			"address": "comdex1qc99y4t293llhaukk6lq74cq8ul38e9da543se",
			"reward": {
				"denom": "ucmdx",
				"amount": "67396"
			}
		},
		{
			"address": "comdex1qcvz08dyk4r7wypv9ry3qjr7vl6kr5w4p6p39f",
			"reward": {
				"denom": "ucmdx",
				"amount": "3114"
			}
		},
		{
			"address": "comdex1qc38akh79mpw5885tjxynmw22gjgk3lvrctc9j",
			"reward": {
				"denom": "ucmdx",
				"amount": "6131"
			}
		},
		{
			"address": "comdex1qcnse5py9eq2l38yl463m6lhqnw2sqjnq6cmjy",
			"reward": {
				"denom": "ucmdx",
				"amount": "54957"
			}
		},
		{
			"address": "comdex1qccenhxav93jl7d4hdc0xha9s0y5pxnnycr5l6",
			"reward": {
				"denom": "ucmdx",
				"amount": "3378"
			}
		},
		{
			"address": "comdex1qc6l0fgxjpft8uwg2fq7wgp5zp6gfadde94kdd",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1qcual5kgmw3gqc9q22hlp0aluc3t7rnsxvmv3n",
			"reward": {
				"denom": "ucmdx",
				"amount": "174"
			}
		},
		{
			"address": "comdex1qcakekgtc86cj8w304vvl5suwx8rm72hp6qucr",
			"reward": {
				"denom": "ucmdx",
				"amount": "13338"
			}
		},
		{
			"address": "comdex1qc7z6an5xkkv7eah0jp2vkqyddwzxay30arq2e",
			"reward": {
				"denom": "ucmdx",
				"amount": "1281"
			}
		},
		{
			"address": "comdex1qcl8m4qje7qm33am779ynx2te9dp3y58l3sm8q",
			"reward": {
				"denom": "ucmdx",
				"amount": "3030"
			}
		},
		{
			"address": "comdex1qezxexfn73n0j5sxt0w20w39lxmtwhmzkuhrll",
			"reward": {
				"denom": "ucmdx",
				"amount": "10664"
			}
		},
		{
			"address": "comdex1qeygrcf0qjrzmr7ky3uzs6ekz9yz5ccx5yy4a4",
			"reward": {
				"denom": "ucmdx",
				"amount": "24832"
			}
		},
		{
			"address": "comdex1qetamgyh93dfwll5jg50jqwn2wdex6vqset80n",
			"reward": {
				"denom": "ucmdx",
				"amount": "294"
			}
		},
		{
			"address": "comdex1qevx2el05ar4grshwpu2mna6hznv27wn756c6n",
			"reward": {
				"denom": "ucmdx",
				"amount": "417"
			}
		},
		{
			"address": "comdex1qed2s9939gj6599lu0qx2hrysm0fmhf6ru3nyn",
			"reward": {
				"denom": "ucmdx",
				"amount": "12798"
			}
		},
		{
			"address": "comdex1qejdt7kts75njqyymh4zm4dudyzvkk8v6axxdt",
			"reward": {
				"denom": "ucmdx",
				"amount": "6308"
			}
		},
		{
			"address": "comdex1qenh2sdg87n0hcr50csmr0gewr6mp9pfhhql7q",
			"reward": {
				"denom": "ucmdx",
				"amount": "2830"
			}
		},
		{
			"address": "comdex1qe5ezh69ttf6559xp6d2hylak4ngw0n0y3kg75",
			"reward": {
				"denom": "ucmdx",
				"amount": "24831"
			}
		},
		{
			"address": "comdex1qe4gmpywlsqah4n9mfanwlazyfew82nl37a4nj",
			"reward": {
				"denom": "ucmdx",
				"amount": "371"
			}
		},
		{
			"address": "comdex1qe4se7ndfmnuc0txfg3w8fmem8505kygun4r8n",
			"reward": {
				"denom": "ucmdx",
				"amount": "1783"
			}
		},
		{
			"address": "comdex1qekwlmgl0uudgcxnkqrp5zfxn3sag9mgxdkrtg",
			"reward": {
				"denom": "ucmdx",
				"amount": "2289"
			}
		},
		{
			"address": "comdex1qemdyfwpaxasm8r75psaskc6s8jxmd3qz6sl8h",
			"reward": {
				"denom": "ucmdx",
				"amount": "100"
			}
		},
		{
			"address": "comdex1qeuhnsyyzaj0xz7xhv75lkfkpadnts0fl7mu3h",
			"reward": {
				"denom": "ucmdx",
				"amount": "3539"
			}
		},
		{
			"address": "comdex1qeaehf7jp7vf0hlzgf056lyvt2qdly7sf6jhf7",
			"reward": {
				"denom": "ucmdx",
				"amount": "49747"
			}
		},
		{
			"address": "comdex1qe70wsxl5vp7ddl4hfxd64zfmgv95hv5ulafvv",
			"reward": {
				"denom": "ucmdx",
				"amount": "11738"
			}
		},
		{
			"address": "comdex1q6rfdu5ky9kuu6agwmltw7ld5vqv2facc9gzga",
			"reward": {
				"denom": "ucmdx",
				"amount": "2181"
			}
		},
		{
			"address": "comdex1q6xw2p673xt9glmkr9nk3tz5fvx7qy6wncqnqv",
			"reward": {
				"denom": "ucmdx",
				"amount": "527"
			}
		},
		{
			"address": "comdex1q62spfecscjrv09mafk5syrsm44z2se8a94t96",
			"reward": {
				"denom": "ucmdx",
				"amount": "3477"
			}
		},
		{
			"address": "comdex1q6v9xfl0lz3d4qyvkydnduasgtgg3e9qwuv8j5",
			"reward": {
				"denom": "ucmdx",
				"amount": "16"
			}
		},
		{
			"address": "comdex1q6wutgt57f3cycaex8q0sze53tuc99e0fk8rdu",
			"reward": {
				"denom": "ucmdx",
				"amount": "39479"
			}
		},
		{
			"address": "comdex1q6syecdedalgfd6egyeey86k46kjz4rqzzywaf",
			"reward": {
				"denom": "ucmdx",
				"amount": "6421"
			}
		},
		{
			"address": "comdex1q6sc3hsaa3yr5h6c4aeu2qh3545c0v0t97uh5u",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1q634y0cc4hwf6cycnarrgw5t8v9rgfdfjdcp3a",
			"reward": {
				"denom": "ucmdx",
				"amount": "34240"
			}
		},
		{
			"address": "comdex1q6n2w6t6ek028tvraju2sluyxa47cku8kpaeng",
			"reward": {
				"denom": "ucmdx",
				"amount": "32"
			}
		},
		{
			"address": "comdex1q64c0ht8ey5j003x4wfkwzs6acvgj4w9k8xq23",
			"reward": {
				"denom": "ucmdx",
				"amount": "13616"
			}
		},
		{
			"address": "comdex1q6h20dn7vp7cla54a2vs4gylq3h23sl8fzxj48",
			"reward": {
				"denom": "ucmdx",
				"amount": "353"
			}
		},
		{
			"address": "comdex1q6hw3zyms2rn2ef92ef6s5sk0syf335p2kx6k8",
			"reward": {
				"denom": "ucmdx",
				"amount": "7115"
			}
		},
		{
			"address": "comdex1q6lyq4xc0pmngpm3wy53u0th8mcy5c5wxm706n",
			"reward": {
				"denom": "ucmdx",
				"amount": "391"
			}
		},
		{
			"address": "comdex1qmqutemah5tzymh9prcr7v7qcwufk3csecgdm7",
			"reward": {
				"denom": "ucmdx",
				"amount": "189"
			}
		},
		{
			"address": "comdex1qmpds0qvrkpj7jzvw5m42k3ptnx2lrsyr5rak4",
			"reward": {
				"denom": "ucmdx",
				"amount": "142525"
			}
		},
		{
			"address": "comdex1qmzcjn5m62895w5apqfxl4j0xgn0cprpyd7eq8",
			"reward": {
				"denom": "ucmdx",
				"amount": "2596"
			}
		},
		{
			"address": "comdex1qm9gkn8xrs87z57z4reevevjy4qztkdtucmkqh",
			"reward": {
				"denom": "ucmdx",
				"amount": "27"
			}
		},
		{
			"address": "comdex1qm9su76qg3nck0sfpvvckzwftf9t3l2qzqkfr5",
			"reward": {
				"denom": "ucmdx",
				"amount": "7232"
			}
		},
		{
			"address": "comdex1qm95d3ygfq7stzuwpnhd0z7r3rgtlt8mxkgp6w",
			"reward": {
				"denom": "ucmdx",
				"amount": "1331"
			}
		},
		{
			"address": "comdex1qm8ttsqm99uucnjqg7h2qzlrl6rxcay7kfrf99",
			"reward": {
				"denom": "ucmdx",
				"amount": "54028"
			}
		},
		{
			"address": "comdex1qmfe576aayfk4rlcp4278cyqgl68wfl7h9s0th",
			"reward": {
				"denom": "ucmdx",
				"amount": "39"
			}
		},
		{
			"address": "comdex1qmdjz5w9ndl90lku3yw87q288r4gvl0qghnsjt",
			"reward": {
				"denom": "ucmdx",
				"amount": "734"
			}
		},
		{
			"address": "comdex1qmwm85yvnz6qhsl2z5g5z6ulhyk7zxktsm2v5a",
			"reward": {
				"denom": "ucmdx",
				"amount": "66203"
			}
		},
		{
			"address": "comdex1qmsf0nuz283t8vujd2cyp2cnu52ncs4fgpwphj",
			"reward": {
				"denom": "ucmdx",
				"amount": "16667"
			}
		},
		{
			"address": "comdex1qmhrq0f5cmu0hclkhc40n80n2srefq8qsqxadn",
			"reward": {
				"denom": "ucmdx",
				"amount": "284"
			}
		},
		{
			"address": "comdex1qmh6n68tk9k3csvgnxt43wutzkr098gqhjp2a9",
			"reward": {
				"denom": "ucmdx",
				"amount": "1967"
			}
		},
		{
			"address": "comdex1qmh7x83kvy9qfkxzhtqcnkmlfl6709muh2vzjg",
			"reward": {
				"denom": "ucmdx",
				"amount": "176"
			}
		},
		{
			"address": "comdex1qmc05n6as5gjlrgt80y4rxp05fl8rtpcau7rpx",
			"reward": {
				"denom": "ucmdx",
				"amount": "97"
			}
		},
		{
			"address": "comdex1qmcsvf4nxctgxdl8fnyv2eldxv597nlxr2x4qt",
			"reward": {
				"denom": "ucmdx",
				"amount": "670"
			}
		},
		{
			"address": "comdex1qm6neyuuma44yw2sqdjudkxwhkakxd4k932792",
			"reward": {
				"denom": "ucmdx",
				"amount": "3702"
			}
		},
		{
			"address": "comdex1qmarfgc9rsuet8m98gv8qa7wuw29w99pygqca9",
			"reward": {
				"denom": "ucmdx",
				"amount": "1240"
			}
		},
		{
			"address": "comdex1qmask3sqaayqp09jcu4aawspzxuxcmqh2s60je",
			"reward": {
				"denom": "ucmdx",
				"amount": "12448"
			}
		},
		{
			"address": "comdex1qup2gufj9rlnqqsv8phrg3u8ddmz9jk3s9kars",
			"reward": {
				"denom": "ucmdx",
				"amount": "134"
			}
		},
		{
			"address": "comdex1quze8x8h29sf772hek66lm9qlhtpkn54efzuj2",
			"reward": {
				"denom": "ucmdx",
				"amount": "1254"
			}
		},
		{
			"address": "comdex1qurknvsheueeepku52naxu297z22rmclpwmajs",
			"reward": {
				"denom": "ucmdx",
				"amount": "28126"
			}
		},
		{
			"address": "comdex1quvfjqv3szjacksqzt0z5dz4muntxga5llhwvx",
			"reward": {
				"denom": "ucmdx",
				"amount": "8857"
			}
		},
		{
			"address": "comdex1quv00c7fu7mgx7h08dftx78ldj50uvzgp99gvc",
			"reward": {
				"denom": "ucmdx",
				"amount": "967"
			}
		},
		{
			"address": "comdex1quwvcyhrqqu2zcx4g4v3ltfm0dlkeqpac6d36u",
			"reward": {
				"denom": "ucmdx",
				"amount": "980"
			}
		},
		{
			"address": "comdex1qus6rx95zx88fj8j4v22zdxnh595grqvuhzm4v",
			"reward": {
				"denom": "ucmdx",
				"amount": "7370"
			}
		},
		{
			"address": "comdex1qunju5mghglplvj32ka0fvwc5hle8yue4du6ue",
			"reward": {
				"denom": "ucmdx",
				"amount": "5745"
			}
		},
		{
			"address": "comdex1qu5s4pmwdjc5533z7705rp8860x5k7t6uyz5r9",
			"reward": {
				"denom": "ucmdx",
				"amount": "166"
			}
		},
		{
			"address": "comdex1quu44sen50gafvcyl8k4k5apt0tm30a2dxv9y5",
			"reward": {
				"denom": "ucmdx",
				"amount": "12819"
			}
		},
		{
			"address": "comdex1quumcuv68zravu6lvddllrw7gy8rawj825cl46",
			"reward": {
				"denom": "ucmdx",
				"amount": "1933"
			}
		},
		{
			"address": "comdex1qu7hlk44y4jkm6407kt3q5ns9uewtksj5amjkp",
			"reward": {
				"denom": "ucmdx",
				"amount": "19869"
			}
		},
		{
			"address": "comdex1qu7mds7kxlhde7u3htkpmdnpmpqrl2ycgjsugg",
			"reward": {
				"denom": "ucmdx",
				"amount": "271"
			}
		},
		{
			"address": "comdex1qaqtsgh3t7a0nctuepfhkw2pecmdcf0qun247a",
			"reward": {
				"denom": "ucmdx",
				"amount": "46024"
			}
		},
		{
			"address": "comdex1qazqc3wr40y3z5gnmydrpxnz5zvxlqwpx7mar6",
			"reward": {
				"denom": "ucmdx",
				"amount": "14263"
			}
		},
		{
			"address": "comdex1qazw2n4ksnhsjt2e4u035mptqfghlvyxsl0rnm",
			"reward": {
				"denom": "ucmdx",
				"amount": "65317"
			}
		},
		{
			"address": "comdex1qarqnp88a82vehgkgqrq5tnf83xhqnsec2l4ge",
			"reward": {
				"denom": "ucmdx",
				"amount": "4501"
			}
		},
		{
			"address": "comdex1qar67xaeshq2sv2ykcrqgtveszftgpw452jkhv",
			"reward": {
				"denom": "ucmdx",
				"amount": "877"
			}
		},
		{
			"address": "comdex1qa9n2hntgwtgjrnl2tl5cwn6k5t7wuz68sk9t5",
			"reward": {
				"denom": "ucmdx",
				"amount": "7017"
			}
		},
		{
			"address": "comdex1qaxrz36ssshztctp3zf54jn53p7jurp3437vvg",
			"reward": {
				"denom": "ucmdx",
				"amount": "4138"
			}
		},
		{
			"address": "comdex1qaxlhdxvv6n5wqtmum0e6dar35ss0rkjvktwst",
			"reward": {
				"denom": "ucmdx",
				"amount": "111"
			}
		},
		{
			"address": "comdex1qag3dqncu6heep4wu8y5fmmasdyr039w69lrh7",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1qa2x7pf83j2hu65h6ldsn2mefh7jlk552j4aa0",
			"reward": {
				"denom": "ucmdx",
				"amount": "1795"
			}
		},
		{
			"address": "comdex1qawhpj28d9s404lkmr5qn7ndj08ahatpdpm2ru",
			"reward": {
				"denom": "ucmdx",
				"amount": "203"
			}
		},
		{
			"address": "comdex1qa3c8fk74xu5m2c7acpdj8uuycugmdzrvxahhe",
			"reward": {
				"denom": "ucmdx",
				"amount": "715"
			}
		},
		{
			"address": "comdex1qajttzprqmq7j6qpj8e8wf2l3nfnux9s2quxu4",
			"reward": {
				"denom": "ucmdx",
				"amount": "1993"
			}
		},
		{
			"address": "comdex1qajjkjrkzz2hjxedmt426w53ky4fdracp9cz52",
			"reward": {
				"denom": "ucmdx",
				"amount": "13302"
			}
		},
		{
			"address": "comdex1qa5srf72t7k8rn8n07qczfqnqsv2e20sasqhls",
			"reward": {
				"denom": "ucmdx",
				"amount": "1082"
			}
		},
		{
			"address": "comdex1qa538wmzzpg323v3eapt3pv4cur09jdw5f2n7p",
			"reward": {
				"denom": "ucmdx",
				"amount": "4974"
			}
		},
		{
			"address": "comdex1qakfg47kavvlm9ycnun3ld2v9c28w35cmz4e2t",
			"reward": {
				"denom": "ucmdx",
				"amount": "6693"
			}
		},
		{
			"address": "comdex1qak3fa94ltmuwd0sjyz0sq5te9tnv6xtnc08h5",
			"reward": {
				"denom": "ucmdx",
				"amount": "3484"
			}
		},
		{
			"address": "comdex1qae37x9ks2t9gakmr874uvud6g0whzh0l02ugc",
			"reward": {
				"denom": "ucmdx",
				"amount": "1220"
			}
		},
		{
			"address": "comdex1qa784nvlwj9nhkvxexh6xj7aa0aus8v4l0u4rq",
			"reward": {
				"denom": "ucmdx",
				"amount": "18329"
			}
		},
		{
			"address": "comdex1q7qf27nf2lkklg6fr7yal9qxq2maemytj9wpzp",
			"reward": {
				"denom": "ucmdx",
				"amount": "1284"
			}
		},
		{
			"address": "comdex1q7pdf8eqhna9x26xjngsrqeg00x4tsq2nj0alm",
			"reward": {
				"denom": "ucmdx",
				"amount": "8886"
			}
		},
		{
			"address": "comdex1q7rx6vrmc9epxze5ne3g39hs7evxj4ltpa4le7",
			"reward": {
				"denom": "ucmdx",
				"amount": "13238"
			}
		},
		{
			"address": "comdex1q7fvden2y88lwvp8x7sptyg26pevyvseynnlf0",
			"reward": {
				"denom": "ucmdx",
				"amount": "14"
			}
		},
		{
			"address": "comdex1q732vjzez0zk04wntf4cjr2wah2wucg220a6pl",
			"reward": {
				"denom": "ucmdx",
				"amount": "69"
			}
		},
		{
			"address": "comdex1q7hg3wg0xyv043s4fpjn7yypu5pna2m0y0t4gg",
			"reward": {
				"denom": "ucmdx",
				"amount": "1234"
			}
		},
		{
			"address": "comdex1q7htpg4phqk7z8lawakew3ne6kvgwmmutua8sy",
			"reward": {
				"denom": "ucmdx",
				"amount": "11239"
			}
		},
		{
			"address": "comdex1q7cflw2u8z4zx29c70kq6qympepq20gyhhxfla",
			"reward": {
				"denom": "ucmdx",
				"amount": "62568"
			}
		},
		{
			"address": "comdex1q7lxm726uxp9dnwjzguuclj7rj7fmtngreywv7",
			"reward": {
				"denom": "ucmdx",
				"amount": "166"
			}
		},
		{
			"address": "comdex1q7lfnfl3jl39grkctkf78zq4qa05d9efdv0cpn",
			"reward": {
				"denom": "ucmdx",
				"amount": "246"
			}
		},
		{
			"address": "comdex1qlznpx885ur2hnd5rr9wu3lw3ecnkdu0c0pxch",
			"reward": {
				"denom": "ucmdx",
				"amount": "16"
			}
		},
		{
			"address": "comdex1ql9w2ulcera7s2vhy2hkerk60zdcm3jqpkmey3",
			"reward": {
				"denom": "ucmdx",
				"amount": "65378"
			}
		},
		{
			"address": "comdex1qlfgdaf734lmw3lfnkky9pcnj88809kjqscupd",
			"reward": {
				"denom": "ucmdx",
				"amount": "16807"
			}
		},
		{
			"address": "comdex1qlvnt85wf7t9ee0yrscfr7fc2n64geek0d2gxz",
			"reward": {
				"denom": "ucmdx",
				"amount": "8288"
			}
		},
		{
			"address": "comdex1qln8y24qn0qgwcqyel9wfx3q6n7cl3ulslf068",
			"reward": {
				"denom": "ucmdx",
				"amount": "2002"
			}
		},
		{
			"address": "comdex1ql5d08k8tmhs6mc6pgu6k69tf5e0e6hx5prqjs",
			"reward": {
				"denom": "ucmdx",
				"amount": "1449"
			}
		},
		{
			"address": "comdex1qlkffk7qp8xxmq6c2xa2y86h8hc9xuukv2cqnt",
			"reward": {
				"denom": "ucmdx",
				"amount": "14110"
			}
		},
		{
			"address": "comdex1ql6wzmjcrpmms2ehv7r0nkekdahkmqwynp6u56",
			"reward": {
				"denom": "ucmdx",
				"amount": "1214"
			}
		},
		{
			"address": "comdex1qlmjr6fagn0usmsqfu73s6c3w7gzqlllyf9rwv",
			"reward": {
				"denom": "ucmdx",
				"amount": "14400"
			}
		},
		{
			"address": "comdex1qlupyxuffv643gny4d9kelpnusj82zgtp8eh4k",
			"reward": {
				"denom": "ucmdx",
				"amount": "413"
			}
		},
		{
			"address": "comdex1qlu0f4e39e5t5mcalvdzhfz8r90fsxv02f56gy",
			"reward": {
				"denom": "ucmdx",
				"amount": "28"
			}
		},
		{
			"address": "comdex1qlljvx0ggg0qmg3ffz4f3hd0tsd259n44z5fan",
			"reward": {
				"denom": "ucmdx",
				"amount": "4720"
			}
		},
		{
			"address": "comdex1pqp2f4p0uvl2jymy4hls69j3jgu492uepaaeem",
			"reward": {
				"denom": "ucmdx",
				"amount": "1788"
			}
		},
		{
			"address": "comdex1pqxta95c83t6r0kjgnn4a0ueh5pdpmyctel4zu",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1pqxes0lne939fjtyk5z93w3q3xea2s6crn389a",
			"reward": {
				"denom": "ucmdx",
				"amount": "1028"
			}
		},
		{
			"address": "comdex1pqtr7e7u9l2p0gv4wuy6nh895n356x9utq03v5",
			"reward": {
				"denom": "ucmdx",
				"amount": "173"
			}
		},
		{
			"address": "comdex1pqt2vefx5weulr8l6vjymz3zjz5kncrf266qq9",
			"reward": {
				"denom": "ucmdx",
				"amount": "169"
			}
		},
		{
			"address": "comdex1pqwzw5aqz8022ysgs2rsgtgr3kk3dfq73u24un",
			"reward": {
				"denom": "ucmdx",
				"amount": "174"
			}
		},
		{
			"address": "comdex1pqw5wkha87378mwjvakl9svaze06lr8x474ca8",
			"reward": {
				"denom": "ucmdx",
				"amount": "5915"
			}
		},
		{
			"address": "comdex1pqnarwhqk229ac79wcakh57gtzwzqwxrmpyjrr",
			"reward": {
				"denom": "ucmdx",
				"amount": "58936"
			}
		},
		{
			"address": "comdex1pqk0mgg8v44x80fm8e6w44dutglhj8sstm9xae",
			"reward": {
				"denom": "ucmdx",
				"amount": "169"
			}
		},
		{
			"address": "comdex1pqc6dvc8fsa3hj2dt5vcwqnq0u3pgq7ff4x5ca",
			"reward": {
				"denom": "ucmdx",
				"amount": "3850"
			}
		},
		{
			"address": "comdex1pqe9ujjuh44enqstdpwv6r90ag9y7ycs6zujgy",
			"reward": {
				"denom": "ucmdx",
				"amount": "145"
			}
		},
		{
			"address": "comdex1pq6j4c6p3uhf9ypmvw3hut0u92xw3vjf57y4wl",
			"reward": {
				"denom": "ucmdx",
				"amount": "174"
			}
		},
		{
			"address": "comdex1pquq3gngtc4v0f3vhx85wxhd4tr84nexfjk946",
			"reward": {
				"denom": "ucmdx",
				"amount": "15"
			}
		},
		{
			"address": "comdex1pquw8s67qykzx9swmhkgzrtaxke00ak2q6fyvn",
			"reward": {
				"denom": "ucmdx",
				"amount": "5634"
			}
		},
		{
			"address": "comdex1pquw2k6cdf6tul657wtda2p93uq8jac4cx4dwl",
			"reward": {
				"denom": "ucmdx",
				"amount": "152"
			}
		},
		{
			"address": "comdex1pqaqkr9echxjydytsqc8m0y2tx39cwqpxdcgz5",
			"reward": {
				"denom": "ucmdx",
				"amount": "14140"
			}
		},
		{
			"address": "comdex1pqa86qxnrjeupcqr42xq3tla3ghf0gxes4ttc8",
			"reward": {
				"denom": "ucmdx",
				"amount": "1429"
			}
		},
		{
			"address": "comdex1ppz54s8yvlmg8kr0009uprlrcatcy578qg0nss",
			"reward": {
				"denom": "ucmdx",
				"amount": "22849"
			}
		},
		{
			"address": "comdex1pprnuuqt25agmue0jg0za49w4w6dcrs4adje8a",
			"reward": {
				"denom": "ucmdx",
				"amount": "6976"
			}
		},
		{
			"address": "comdex1pptln64s9zllu9rzkfsnjl9zx6x3vnshxuu5p7",
			"reward": {
				"denom": "ucmdx",
				"amount": "5825"
			}
		},
		{
			"address": "comdex1ppdegtvc8lw2xxajw2sftj7wf4838hrw92essv",
			"reward": {
				"denom": "ucmdx",
				"amount": "26"
			}
		},
		{
			"address": "comdex1ppw6nzkhkpfl69heq3qd4cvwh3z0y0wz8wzzda",
			"reward": {
				"denom": "ucmdx",
				"amount": "185"
			}
		},
		{
			"address": "comdex1ppwmgd463a5ka8ea846gu0jmwyxcjpyr2s23ul",
			"reward": {
				"denom": "ucmdx",
				"amount": "76046"
			}
		},
		{
			"address": "comdex1pp0xtad9k74hncrdxqt4e779ttadtwc499nxrj",
			"reward": {
				"denom": "ucmdx",
				"amount": "20476"
			}
		},
		{
			"address": "comdex1ppska05vdnx5gyq46su0y96ukdnte82h8urzn2",
			"reward": {
				"denom": "ucmdx",
				"amount": "846"
			}
		},
		{
			"address": "comdex1ppk3ascqs5zrt9qqg3y579lp005uvnkhgz0w6l",
			"reward": {
				"denom": "ucmdx",
				"amount": "4157"
			}
		},
		{
			"address": "comdex1ppkenk3j9j7as062djdk7hx7ps346vsa8nsrnh",
			"reward": {
				"denom": "ucmdx",
				"amount": "10221"
			}
		},
		{
			"address": "comdex1ppuj9x33e6c4ee5axkdfdjavkpg2y6msmnaguh",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1pzzdvazgat8t9epvh2n5xn6wk4zcfc54fdqyqq",
			"reward": {
				"denom": "ucmdx",
				"amount": "83959"
			}
		},
		{
			"address": "comdex1pz9r9j7m2u72ale0yulrjmps9tpnkpnvjwjg8r",
			"reward": {
				"denom": "ucmdx",
				"amount": "22864"
			}
		},
		{
			"address": "comdex1pz9cawmf0jfa8d9emew5mxsla2nyutu0rse26m",
			"reward": {
				"denom": "ucmdx",
				"amount": "699"
			}
		},
		{
			"address": "comdex1pz8cqgv97kd9552w9nffr7zwxhcjkw00j94d04",
			"reward": {
				"denom": "ucmdx",
				"amount": "2889"
			}
		},
		{
			"address": "comdex1pzgmkzvw2nuav5j8yp6gd8ymjpfr55ym9wtmge",
			"reward": {
				"denom": "ucmdx",
				"amount": "6310"
			}
		},
		{
			"address": "comdex1pzdtnxnaexxqrv48fsa98fxqtsm5l4rtzm0hfk",
			"reward": {
				"denom": "ucmdx",
				"amount": "143"
			}
		},
		{
			"address": "comdex1pzwel4a5ft2256au8xqy89k04yxavgs5fg3e73",
			"reward": {
				"denom": "ucmdx",
				"amount": "21383"
			}
		},
		{
			"address": "comdex1pz0gys67zv6dl0g69dree3yees48plzu6mcake",
			"reward": {
				"denom": "ucmdx",
				"amount": "530"
			}
		},
		{
			"address": "comdex1pz0a4rzkszss4dzs7yfukuskqtjg739gqar3w2",
			"reward": {
				"denom": "ucmdx",
				"amount": "18835219"
			}
		},
		{
			"address": "comdex1pzsdh699967u7mcsh9vpxncpq05lpvzxxefxj8",
			"reward": {
				"denom": "ucmdx",
				"amount": "3119"
			}
		},
		{
			"address": "comdex1pz3atcgt6ujam2amyu4jdcaavtpkkykmx0247r",
			"reward": {
				"denom": "ucmdx",
				"amount": "406"
			}
		},
		{
			"address": "comdex1pzncy4h4mh685kkrn2ln2fc3wnsd0cl4lucxe3",
			"reward": {
				"denom": "ucmdx",
				"amount": "3461"
			}
		},
		{
			"address": "comdex1pzkvnc36hqk74qgk0e0403n3t647rm0wd20hr2",
			"reward": {
				"denom": "ucmdx",
				"amount": "3022"
			}
		},
		{
			"address": "comdex1pzknn7p9k5w2ryvlxgh886qc8yqpa42ucsjyhu",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1pzmzdd4enwl203mdxyt3fftwxvl7f7u7032c79",
			"reward": {
				"denom": "ucmdx",
				"amount": "6332"
			}
		},
		{
			"address": "comdex1pzmxp7qlqmnj72t4qy7tr825x9ssm508qjaggf",
			"reward": {
				"denom": "ucmdx",
				"amount": "3454"
			}
		},
		{
			"address": "comdex1pzu3pe95fenm3tyrwh70r55dsj2hyxsdgvecx9",
			"reward": {
				"denom": "ucmdx",
				"amount": "152"
			}
		},
		{
			"address": "comdex1pzunc9mys4lkr02nncdsg40we93qhj7n9j4sz5",
			"reward": {
				"denom": "ucmdx",
				"amount": "14"
			}
		},
		{
			"address": "comdex1pz78njwa8h0f63jwu8cx7m6xt7a5jakl527jdg",
			"reward": {
				"denom": "ucmdx",
				"amount": "2443"
			}
		},
		{
			"address": "comdex1prz29f6d3pqhs6l22wlulqxvc3pw3khev93myl",
			"reward": {
				"denom": "ucmdx",
				"amount": "1122"
			}
		},
		{
			"address": "comdex1prz77l7vz5ecgm6h5utu2qvcx72f2evamg3fdz",
			"reward": {
				"denom": "ucmdx",
				"amount": "10432"
			}
		},
		{
			"address": "comdex1prr2x6gy692cllg2zxtuh5pz70n2a7tjqg49n2",
			"reward": {
				"denom": "ucmdx",
				"amount": "124"
			}
		},
		{
			"address": "comdex1pr953ns5zej5wzmrvm9rxa380xmzaacryk80pk",
			"reward": {
				"denom": "ucmdx",
				"amount": "1189"
			}
		},
		{
			"address": "comdex1pr97avak3g7rvhf0k9w93fwsqlr8k8ekyjqljj",
			"reward": {
				"denom": "ucmdx",
				"amount": "528"
			}
		},
		{
			"address": "comdex1prg6sgs5dfga8r0t8vtj8t4tmee7v5gd7hrf7h",
			"reward": {
				"denom": "ucmdx",
				"amount": "0"
			}
		},
		{
			"address": "comdex1pr2le0sz9gsuaatms36dt9lv3ttumwlk6uxj2n",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1prtp3ruy507mjqw0t2p9whhewqlzy60wewld6l",
			"reward": {
				"denom": "ucmdx",
				"amount": "13546"
			}
		},
		{
			"address": "comdex1pr3n9rlpyj8zvaau9cl383d6anfly7l2ufg7f5",
			"reward": {
				"denom": "ucmdx",
				"amount": "167"
			}
		},
		{
			"address": "comdex1prkt7s379wupzh8nmuefa9hw5xyddqwk5sqth6",
			"reward": {
				"denom": "ucmdx",
				"amount": "1441"
			}
		},
		{
			"address": "comdex1prk4pyjhk0z8jwn7qz9ppsz98aanj0rrqtnzae",
			"reward": {
				"denom": "ucmdx",
				"amount": "1401"
			}
		},
		{
			"address": "comdex1prm9rdpf9eqx8kjq2qgsyrkg4x2m9pxvk8x6wr",
			"reward": {
				"denom": "ucmdx",
				"amount": "151"
			}
		},
		{
			"address": "comdex1pruhnn42hy2gy2qp4ndhknnugpjsh5ur3wmxwv",
			"reward": {
				"denom": "ucmdx",
				"amount": "1"
			}
		},
		{
			"address": "comdex1pruhh5tk27ya8twfee5m55tznk7tkdjv7l0t2c",
			"reward": {
				"denom": "ucmdx",
				"amount": "5294"
			}
		},
		{
			"address": "comdex1pr7pcnzts7cjmdte983psa2wkqemz3wwa4dwgq",
			"reward": {
				"denom": "ucmdx",
				"amount": "1778"
			}
		},
		{
			"address": "comdex1py9nkcxgn05qdgrhvfw2z7pezjqjj8ml4r0zzu",
			"reward": {
				"denom": "ucmdx",
				"amount": "302"
			}
		},
		{
			"address": "comdex1pyxg24gsfcnqzz3n03m7s50kyh078rmxdwkl7e",
			"reward": {
				"denom": "ucmdx",
				"amount": "163460"
			}
		},
		{
			"address": "comdex1py0nm22hy3gc6nxg6rw8k6luz4twr6fht0te78",
			"reward": {
				"denom": "ucmdx",
				"amount": "5062"
			}
		},
		{
			"address": "comdex1py3xvgkzr4m3nc98xnh4rvvjmngj9lmncdk8uf",
			"reward": {
				"denom": "ucmdx",
				"amount": "71729"
			}
		},
		{
			"address": "comdex1pyjqy5l25fl7953n7ptfarr0mx3n4a6c7fuf3j",
			"reward": {
				"denom": "ucmdx",
				"amount": "2631"
			}
		},
		{
			"address": "comdex1pynvf0cpaeysqkqrdd733ynuh96p43wyatuwa9",
			"reward": {
				"denom": "ucmdx",
				"amount": "55855"
			}
		},
		{
			"address": "comdex1pyh3q3u08urdjfq08qqaafl4kqfhyr6kx73y5m",
			"reward": {
				"denom": "ucmdx",
				"amount": "185"
			}
		},
		{
			"address": "comdex1pycrutgzg8e0w4yfkslgr8vfnpulem2zzef3y0",
			"reward": {
				"denom": "ucmdx",
				"amount": "709"
			}
		},
		{
			"address": "comdex1pyc4kpalk2dk7cqu63yecnuyg5cpcdzh8ghcqh",
			"reward": {
				"denom": "ucmdx",
				"amount": "48334"
			}
		},
		{
			"address": "comdex1pye9afwsyhupdk306vvalv7sh2mu76e20yyvwm",
			"reward": {
				"denom": "ucmdx",
				"amount": "35"
			}
		},
		{
			"address": "comdex1p9pypu9wk3q3cnvn70jq7q294drz40qasa3at8",
			"reward": {
				"denom": "ucmdx",
				"amount": "195"
			}
		},
		{
			"address": "comdex1p98hhcjx6wzfd4j87parzrdzefwr3jzhhhfxh7",
			"reward": {
				"denom": "ucmdx",
				"amount": "70033"
			}
		},
		{
			"address": "comdex1p9gywh35t4ljy8ez84ydmz3x8nf2v3lfrsthga",
			"reward": {
				"denom": "ucmdx",
				"amount": "674"
			}
		},
		{
			"address": "comdex1p9g22jhpsh47nccdgpdnfmm488y26l8zxnmxcd",
			"reward": {
				"denom": "ucmdx",
				"amount": "1428"
			}
		},
		{
			"address": "comdex1p90jjtth9pclh9258cxvszunysjezufx0hdre8",
			"reward": {
				"denom": "ucmdx",
				"amount": "4546"
			}
		},
		{
			"address": "comdex1p93yh228z48c649kspt52nwxp3pfu95a80rxqq",
			"reward": {
				"denom": "ucmdx",
				"amount": "4717"
			}
		},
		{
			"address": "comdex1p9jnw5vtmhg85hlvqvreluacuxc95yxv94pvd7",
			"reward": {
				"denom": "ucmdx",
				"amount": "17648"
			}
		},
		{
			"address": "comdex1p9n3455yz6s6w0avvdmhlgcrskexzje4m4kt0m",
			"reward": {
				"denom": "ucmdx",
				"amount": "1869"
			}
		},
		{
			"address": "comdex1p95ec7z25fuvteuqmplav6xz9dg8wpqfc6hs5c",
			"reward": {
				"denom": "ucmdx",
				"amount": "498"
			}
		},
		{
			"address": "comdex1p9cmpq4uwg48q9x9jnykqppl3j78an07lqra22",
			"reward": {
				"denom": "ucmdx",
				"amount": "1278"
			}
		},
		{
			"address": "comdex1p9erq3t2shufd9lyzqelykd2sfh5206aj80h4g",
			"reward": {
				"denom": "ucmdx",
				"amount": "1761"
			}
		},
		{
			"address": "comdex1p96xymq2k3kwa7r9shkusgjgwuy0sl8hl0ny55",
			"reward": {
				"denom": "ucmdx",
				"amount": "1266"
			}
		},
		{
			"address": "comdex1p9ltgksqnlhye6gazupsx4ecmshl585tqktq6r",
			"reward": {
				"denom": "ucmdx",
				"amount": "527"
			}
		},
		{
			"address": "comdex1pxqrk6n4lrzez4k08qct7hytlt4v2rl6eqyzyh",
			"reward": {
				"denom": "ucmdx",
				"amount": "168"
			}
		},
		{
			"address": "comdex1pxqhkgj6tr7ejfg2pyv9v5sekqncv8mmqywkh5",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1pxylcm978an30qyz82pcfatgskgqrcs4cu8zum",
			"reward": {
				"denom": "ucmdx",
				"amount": "1399"
			}
		},
		{
			"address": "comdex1px98h5f5gvp3ferer6w2d7wsrxmg4fzgrqk7xr",
			"reward": {
				"denom": "ucmdx",
				"amount": "2277"
			}
		},
		{
			"address": "comdex1px28v5f9mk0039gurky40uwwkqllun93eu3duz",
			"reward": {
				"denom": "ucmdx",
				"amount": "13862"
			}
		},
		{
			"address": "comdex1px0l5j8wmr8ar9h0p94z3fqgn6rna3kuycal9n",
			"reward": {
				"denom": "ucmdx",
				"amount": "1434"
			}
		},
		{
			"address": "comdex1px3zpahc6g0vjuve2ztcnrpx2mldnzty0yh3fn",
			"reward": {
				"denom": "ucmdx",
				"amount": "303"
			}
		},
		{
			"address": "comdex1pxnmvwh2eg2gkzatdfnxdl37wjzzr4xy9xqxch",
			"reward": {
				"denom": "ucmdx",
				"amount": "3351"
			}
		},
		{
			"address": "comdex1px5ws844l7nd60tmwes7upvrqmjhnn3036yzp9",
			"reward": {
				"denom": "ucmdx",
				"amount": "6209"
			}
		},
		{
			"address": "comdex1px5u5kzf8fdqtuwn73aw258t20h3r335dheqmf",
			"reward": {
				"denom": "ucmdx",
				"amount": "9617"
			}
		},
		{
			"address": "comdex1px5a056tnwhtr5dh9pt8qs8q4vmkdyqw5tcpj4",
			"reward": {
				"denom": "ucmdx",
				"amount": "1733"
			}
		},
		{
			"address": "comdex1pxh55qj7xfkayesj4kneqrrwlkf000f8qndale",
			"reward": {
				"denom": "ucmdx",
				"amount": "59508"
			}
		},
		{
			"address": "comdex1pxe95c9wlhh040c4uygmnsazx9t2rkz6493zqh",
			"reward": {
				"denom": "ucmdx",
				"amount": "494"
			}
		},
		{
			"address": "comdex1pxu39lh7xt53tmsxmfzyltev2ke995lchyyufy",
			"reward": {
				"denom": "ucmdx",
				"amount": "18325"
			}
		},
		{
			"address": "comdex1pxa93n7xrvvmzd9rv0zsxgfk0kjude6tk3r7rs",
			"reward": {
				"denom": "ucmdx",
				"amount": "18035"
			}
		},
		{
			"address": "comdex1p8y9wkftsjurpuwv59ppzntt4696ss2a4v0nwm",
			"reward": {
				"denom": "ucmdx",
				"amount": "2828"
			}
		},
		{
			"address": "comdex1p89n4mm2zmydeq3smlzydwcnrsqp9jzn54m6k2",
			"reward": {
				"denom": "ucmdx",
				"amount": "753"
			}
		},
		{
			"address": "comdex1p897yys9cj2r28w3hljqczsaf2rjde88zn4fp6",
			"reward": {
				"denom": "ucmdx",
				"amount": "125"
			}
		},
		{
			"address": "comdex1p824ept88e349z7tak2nwznjdyn6wpckgmqlwv",
			"reward": {
				"denom": "ucmdx",
				"amount": "349"
			}
		},
		{
			"address": "comdex1p80y4t4ct57sv70kmdyxalmrkhp6ssdtdznknv",
			"reward": {
				"denom": "ucmdx",
				"amount": "47854"
			}
		},
		{
			"address": "comdex1p8sydzqswrga8wwgxwqu5p7ru8yrmrqa95c4xj",
			"reward": {
				"denom": "ucmdx",
				"amount": "4051"
			}
		},
		{
			"address": "comdex1p8hf5sffxwhute06g2gzpwejtwdt97mcvrur7v",
			"reward": {
				"denom": "ucmdx",
				"amount": "17523"
			}
		},
		{
			"address": "comdex1p8hak8tyt9xr0ayh7grrvu3fmdzgdhr63f3hk8",
			"reward": {
				"denom": "ucmdx",
				"amount": "143"
			}
		},
		{
			"address": "comdex1p8cc6z6e4uh0yw7aldlnfhc6g8cdk8mcj62dz0",
			"reward": {
				"denom": "ucmdx",
				"amount": "1337"
			}
		},
		{
			"address": "comdex1p8epswmef4uwcr44guy3427rkaht08r5s9jrmt",
			"reward": {
				"denom": "ucmdx",
				"amount": "323"
			}
		},
		{
			"address": "comdex1p87ekhls82shd9f6mu7kh3fyqxjwx5cmp8r0pk",
			"reward": {
				"denom": "ucmdx",
				"amount": "1692"
			}
		},
		{
			"address": "comdex1p8lsptxgk6rnmhlkfg7y5yy8rra0n33srz2qq6",
			"reward": {
				"denom": "ucmdx",
				"amount": "8250"
			}
		},
		{
			"address": "comdex1p8ljzm5vdxu65h5vn28du94wt260tzezjg5v6j",
			"reward": {
				"denom": "ucmdx",
				"amount": "17793"
			}
		},
		{
			"address": "comdex1pgypskthfrl70g9n72yrrf9xstqtmvdaez46vk",
			"reward": {
				"denom": "ucmdx",
				"amount": "19"
			}
		},
		{
			"address": "comdex1pg8dr463ngkzlqxussm23yugtepm94va3x5w6y",
			"reward": {
				"denom": "ucmdx",
				"amount": "7223"
			}
		},
		{
			"address": "comdex1pggzlafefwe7e0zyh9y8cn3gg6sjv86k5f9ydx",
			"reward": {
				"denom": "ucmdx",
				"amount": "1768"
			}
		},
		{
			"address": "comdex1pggsssjdvf2kdnfg9gpzffz9vlfyssns4udmzq",
			"reward": {
				"denom": "ucmdx",
				"amount": "173"
			}
		},
		{
			"address": "comdex1pg23whjnnheuuq83qy9ehap9md28jzdfuwknc8",
			"reward": {
				"denom": "ucmdx",
				"amount": "81954"
			}
		},
		{
			"address": "comdex1pg25yypnkc5dhl2ec56na0ulvanpm0ppy0ckyv",
			"reward": {
				"denom": "ucmdx",
				"amount": "1760"
			}
		},
		{
			"address": "comdex1pgw6gqp8w7gu7vv337ky6h35egqu9efx4n3qpw",
			"reward": {
				"denom": "ucmdx",
				"amount": "2010"
			}
		},
		{
			"address": "comdex1pgw6ups2nsgr3y76q6sv7z9qkvldwnu4cn7c55",
			"reward": {
				"denom": "ucmdx",
				"amount": "3016"
			}
		},
		{
			"address": "comdex1pg0lfu3t045qj76k59kfj24xvdv8d4g92w8rq4",
			"reward": {
				"denom": "ucmdx",
				"amount": "1903"
			}
		},
		{
			"address": "comdex1pgjsu5996dk9gr8wy4z357k5lq4x7werv6002g",
			"reward": {
				"denom": "ucmdx",
				"amount": "5542"
			}
		},
		{
			"address": "comdex1pgnpjs3zg6sgm9pdhtrduqyyx4lq3l4f3z89y3",
			"reward": {
				"denom": "ucmdx",
				"amount": "165593"
			}
		},
		{
			"address": "comdex1pg4scf0eajyzz97fq6lcmjckmsravygn53zr2y",
			"reward": {
				"denom": "ucmdx",
				"amount": "2896"
			}
		},
		{
			"address": "comdex1pgkuwtythc4zdj3z6pjzrlau0gtpqscvdrg0sq",
			"reward": {
				"denom": "ucmdx",
				"amount": "718"
			}
		},
		{
			"address": "comdex1pghu87sprsw9kk6d78uxrmjru9rcs48d0s3lsq",
			"reward": {
				"denom": "ucmdx",
				"amount": "29257"
			}
		},
		{
			"address": "comdex1pgmp5f28r0p24g5xqjk73es3nwfse79q83v3yl",
			"reward": {
				"denom": "ucmdx",
				"amount": "7009"
			}
		},
		{
			"address": "comdex1pgaqzkvudes4nrm0hm3ktzu064ckwcnt583scy",
			"reward": {
				"denom": "ucmdx",
				"amount": "17528"
			}
		},
		{
			"address": "comdex1pganmyl5ty6xmnxwcyu8azg3n97x2n8l37kkqz",
			"reward": {
				"denom": "ucmdx",
				"amount": "5309"
			}
		},
		{
			"address": "comdex1pgakvl4737jtqwvk8jvut7kwrrl9trdk4l2u7m",
			"reward": {
				"denom": "ucmdx",
				"amount": "2642"
			}
		},
		{
			"address": "comdex1pgamxw340zsv2fgcas0mqglw8tnm44795tg59f",
			"reward": {
				"denom": "ucmdx",
				"amount": "1519"
			}
		},
		{
			"address": "comdex1pgau0y96k6y2j2493rmq0f290vwj8ekxtezurh",
			"reward": {
				"denom": "ucmdx",
				"amount": "3545"
			}
		},
		{
			"address": "comdex1pfqs0edvzjawzv9ynjssny78mrf65t68c4xkhy",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1pfpk6rsgc0x7g8q7y0yre86ej2enr8rhr9s4z9",
			"reward": {
				"denom": "ucmdx",
				"amount": "3393"
			}
		},
		{
			"address": "comdex1pfzsn5r30jun35k8xmw59tkylkqedfkr8erckv",
			"reward": {
				"denom": "ucmdx",
				"amount": "7631"
			}
		},
		{
			"address": "comdex1pfrgaw87sxlhldjz9p38murx2dtu3h4v7pgtjy",
			"reward": {
				"denom": "ucmdx",
				"amount": "2839"
			}
		},
		{
			"address": "comdex1pfye4l3ngvyzh5ax8lxug3cezvdqkq85qm93m5",
			"reward": {
				"denom": "ucmdx",
				"amount": "1944"
			}
		},
		{
			"address": "comdex1pf99ax4gzjz7fkxnyxglydc04xy7vmt0u2w3ay",
			"reward": {
				"denom": "ucmdx",
				"amount": "88"
			}
		},
		{
			"address": "comdex1pf9nd2hwqu0dwufn76wtn3rxj3ms7eyhkqefd9",
			"reward": {
				"denom": "ucmdx",
				"amount": "1248"
			}
		},
		{
			"address": "comdex1pf835gu54ale0g73jcnkpt6s86f8c35wctj2vr",
			"reward": {
				"denom": "ucmdx",
				"amount": "69830"
			}
		},
		{
			"address": "comdex1pfgg366ueh3caztdtx3r098520kt5pnkx520q5",
			"reward": {
				"denom": "ucmdx",
				"amount": "17906"
			}
		},
		{
			"address": "comdex1pfgdrap4z6w0de7f0wuuzwl5lm8td73jyurwxc",
			"reward": {
				"denom": "ucmdx",
				"amount": "140"
			}
		},
		{
			"address": "comdex1pfvdn5903gp45g2vadlejdwgkw2yjul55udjg2",
			"reward": {
				"denom": "ucmdx",
				"amount": "5877"
			}
		},
		{
			"address": "comdex1pfd8c7j9sm05gl2q6nhgv5guykdd9eszf9d2d8",
			"reward": {
				"denom": "ucmdx",
				"amount": "4125"
			}
		},
		{
			"address": "comdex1pfwn33xcjt54rkkphx66xqqp098f5lprxuwmjv",
			"reward": {
				"denom": "ucmdx",
				"amount": "20153"
			}
		},
		{
			"address": "comdex1pfwcevm899rzuwek99774am7q6glzpzsgurev0",
			"reward": {
				"denom": "ucmdx",
				"amount": "874"
			}
		},
		{
			"address": "comdex1pf0gjuduse5rsne403j9w5re0h06dd3fv4ml3s",
			"reward": {
				"denom": "ucmdx",
				"amount": "1931"
			}
		},
		{
			"address": "comdex1pfsm2neda65f8gzfpn3lecym62n8lxxumesyls",
			"reward": {
				"denom": "ucmdx",
				"amount": "180"
			}
		},
		{
			"address": "comdex1pf3t7czgnqvqyruaag9nynp0v5n5sggf45t7pr",
			"reward": {
				"denom": "ucmdx",
				"amount": "2872"
			}
		},
		{
			"address": "comdex1pf37xevfm0e4gnc2l3h0vxc5n0yy0vleckhkm3",
			"reward": {
				"denom": "ucmdx",
				"amount": "16"
			}
		},
		{
			"address": "comdex1pfjmpqzapgfm5utxrp6ctp0c3l4sw5revuwl82",
			"reward": {
				"denom": "ucmdx",
				"amount": "1443"
			}
		},
		{
			"address": "comdex1pf4aveqq2td8c7uwq76zwmjntwlutvkle2lac5",
			"reward": {
				"denom": "ucmdx",
				"amount": "61848"
			}
		},
		{
			"address": "comdex1pfhmse3detzwv8qyey6gc9m3rqzpnsvpqkhyjw",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1pfcgcg92jnmx4afg2g889vq4nl88hzu8szxvsc",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1pfewwatgpkpyrnrjrxn5q5wuyxm3fv9lwgvmxn",
			"reward": {
				"denom": "ucmdx",
				"amount": "42382"
			}
		},
		{
			"address": "comdex1pfmazly4w46wjj5x8espwve8lgc6u337ryjmsh",
			"reward": {
				"denom": "ucmdx",
				"amount": "10907"
			}
		},
		{
			"address": "comdex1p2py5f4znc7mqr06hrh737q5egl9updpj8hd2t",
			"reward": {
				"denom": "ucmdx",
				"amount": "6915"
			}
		},
		{
			"address": "comdex1p2peta3lrhuvez6fux4zh0q7vn064t78tz4w65",
			"reward": {
				"denom": "ucmdx",
				"amount": "2245"
			}
		},
		{
			"address": "comdex1p2xvlepwqf5cqasvpwt8h0je5af28ghu2neyk8",
			"reward": {
				"denom": "ucmdx",
				"amount": "1764"
			}
		},
		{
			"address": "comdex1p2g0ztqcattvs5wsauhckpdfz5pl8qe9z2d4vm",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1p2g4wzsv7cu6nlvqm8yvkaw8dx924k83yygltp",
			"reward": {
				"denom": "ucmdx",
				"amount": "1400"
			}
		},
		{
			"address": "comdex1p22f5a6vk7lclm0rwa0rs4q70afv8lf6h8pyav",
			"reward": {
				"denom": "ucmdx",
				"amount": "1838"
			}
		},
		{
			"address": "comdex1p2vlp30qlp0g005hjh2ryyy9tae0ya72zdulzc",
			"reward": {
				"denom": "ucmdx",
				"amount": "1931"
			}
		},
		{
			"address": "comdex1p23xncur3nhq24p4maw3hcv3ftv8dzguqrhupr",
			"reward": {
				"denom": "ucmdx",
				"amount": "96099"
			}
		},
		{
			"address": "comdex1p233jwaqxsacp8eyqzj4gtxf46jjeejhj2dh6l",
			"reward": {
				"denom": "ucmdx",
				"amount": "1308"
			}
		},
		{
			"address": "comdex1p2js62ftaqa6tg4gn7ds9xe5h07eyremqc4vtd",
			"reward": {
				"denom": "ucmdx",
				"amount": "606"
			}
		},
		{
			"address": "comdex1p2n7hz5lkcxft7x5tlnf6l5p33t3najvx2wzp7",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1p24gny9u8jqsvs7cqn5pkzu4xr6zeelmve34jr",
			"reward": {
				"denom": "ucmdx",
				"amount": "34974"
			}
		},
		{
			"address": "comdex1p2kgzxlq3afe3zwael2xnf8gmuymrp3k9ked5w",
			"reward": {
				"denom": "ucmdx",
				"amount": "1310"
			}
		},
		{
			"address": "comdex1p2k2w0enr4vsjchfexr8079y6sjy5u9p9ugsht",
			"reward": {
				"denom": "ucmdx",
				"amount": "4175"
			}
		},
		{
			"address": "comdex1p2kskgltakdqdz042uvaw4x2lz3f92kwq23vyr",
			"reward": {
				"denom": "ucmdx",
				"amount": "143"
			}
		},
		{
			"address": "comdex1p2hhm8qm5csqrk2pdh5vz7edfm0th8q86k7rdm",
			"reward": {
				"denom": "ucmdx",
				"amount": "149"
			}
		},
		{
			"address": "comdex1p26kvfeg69mghcnpvhgly3qqnxqz4ntheety3n",
			"reward": {
				"denom": "ucmdx",
				"amount": "1811"
			}
		},
		{
			"address": "comdex1p2uqpr9rhjh94kqf4wluqa9ze2rklaqnrsh5jv",
			"reward": {
				"denom": "ucmdx",
				"amount": "2932"
			}
		},
		{
			"address": "comdex1p27w2f0ju60f8ynuwqj7tappd0k3jyrnt4c2ph",
			"reward": {
				"denom": "ucmdx",
				"amount": "17910"
			}
		},
		{
			"address": "comdex1ptqk8rlpyc4yx5jqrqj8dyjjkg23235hjp6s6w",
			"reward": {
				"denom": "ucmdx",
				"amount": "16666"
			}
		},
		{
			"address": "comdex1ptpau0ksgezw9tmvzsexg6hxj0rjelrdswjfrj",
			"reward": {
				"denom": "ucmdx",
				"amount": "2"
			}
		},
		{
			"address": "comdex1ptyzewnns2kn37ewtmv6ppsvhdnmeapvfjwj3s",
			"reward": {
				"denom": "ucmdx",
				"amount": "19595"
			}
		},
		{
			"address": "comdex1pt93qp6l7k22qmtwyf2qc2vdx453jexad8makg",
			"reward": {
				"denom": "ucmdx",
				"amount": "1020"
			}
		},
		{
			"address": "comdex1ptxfhlh83vp9yxkgfjffvxrp226jc6xqcmfa3e",
			"reward": {
				"denom": "ucmdx",
				"amount": "6417"
			}
		},
		{
			"address": "comdex1pt8v4wxwac3nzmmvwpwa9s7gj3ehu2hwlxt4am",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1pt8atmzp05z7q32nr42zf7xu0t9gzdky25ww38",
			"reward": {
				"denom": "ucmdx",
				"amount": "203"
			}
		},
		{
			"address": "comdex1ptg3dwrukw3nzygk20zzj555l633vw9ysqvweh",
			"reward": {
				"denom": "ucmdx",
				"amount": "1034"
			}
		},
		{
			"address": "comdex1ptt8dz47jmwnypnmy5num8kh7rvuf2w3m3wmxc",
			"reward": {
				"denom": "ucmdx",
				"amount": "34556"
			}
		},
		{
			"address": "comdex1pttg3s34aqpehnksjk45hwte3qdmug3uw7l949",
			"reward": {
				"denom": "ucmdx",
				"amount": "8162"
			}
		},
		{
			"address": "comdex1ptd4ud0l7uezd8jhlj97cna5rl8l7njl2huwu2",
			"reward": {
				"denom": "ucmdx",
				"amount": "1213"
			}
		},
		{
			"address": "comdex1pt082dcc38wxf0m9lu6kfwhwr2vxf6vjf77wmd",
			"reward": {
				"denom": "ucmdx",
				"amount": "1723"
			}
		},
		{
			"address": "comdex1ptjxctv247qs33tfny5nsq9s49jacvk7ftvrec",
			"reward": {
				"denom": "ucmdx",
				"amount": "378"
			}
		},
		{
			"address": "comdex1ptn7q5l3qcxs3erwsc623p22ym5v7dx0827edq",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1pt5v3nfce6arjaquevtedd6j36v84jm85zfadx",
			"reward": {
				"denom": "ucmdx",
				"amount": "1279"
			}
		},
		{
			"address": "comdex1pt5kqgwcr6vm3xz22hdkrs894q54djpvasfedd",
			"reward": {
				"denom": "ucmdx",
				"amount": "15148"
			}
		},
		{
			"address": "comdex1pt57grkz0rl95ydjxnwr4zc2ftdlzj2cxqa2ej",
			"reward": {
				"denom": "ucmdx",
				"amount": "353"
			}
		},
		{
			"address": "comdex1ptkyz90pfemjy5dvnwyuexn72pv8qdkseq0va2",
			"reward": {
				"denom": "ucmdx",
				"amount": "3581"
			}
		},
		{
			"address": "comdex1ptupdv8yhw67klgn8zd5gfzl8w8esvxha3hrmx",
			"reward": {
				"denom": "ucmdx",
				"amount": "6151"
			}
		},
		{
			"address": "comdex1pt7c492jg0gr663ua6rauee4aa9n45w60klz8v",
			"reward": {
				"denom": "ucmdx",
				"amount": "396"
			}
		},
		{
			"address": "comdex1ptlrz32rvykae5adk946evevxzspcnr85ltfq7",
			"reward": {
				"denom": "ucmdx",
				"amount": "174"
			}
		},
		{
			"address": "comdex1pvzssder5a5md4wdhjyzdhtxmcrpvp9c6q763t",
			"reward": {
				"denom": "ucmdx",
				"amount": "1761"
			}
		},
		{
			"address": "comdex1pvrzeqg5kslvmjdrmy2ekg82739jmx9mvz2gcu",
			"reward": {
				"denom": "ucmdx",
				"amount": "143655"
			}
		},
		{
			"address": "comdex1pv93gycahd0evntus6qvj4zr0xm0gj6m56puek",
			"reward": {
				"denom": "ucmdx",
				"amount": "176"
			}
		},
		{
			"address": "comdex1pv8elrve2gzgpuh27dzdu78x367sjuuja662za",
			"reward": {
				"denom": "ucmdx",
				"amount": "698"
			}
		},
		{
			"address": "comdex1pv8a4jpkem62n6x8ml8tsf4gr2mk2lcv9pakzy",
			"reward": {
				"denom": "ucmdx",
				"amount": "4521"
			}
		},
		{
			"address": "comdex1pvf4ju252as6yzl6lasc45a348vhcnz5caqcf2",
			"reward": {
				"denom": "ucmdx",
				"amount": "6180"
			}
		},
		{
			"address": "comdex1pv29j6xfgspv0cl8wqh0qyea24dt5ygaexnm4m",
			"reward": {
				"denom": "ucmdx",
				"amount": "20583"
			}
		},
		{
			"address": "comdex1pv25l0v7e073mkhaz97vm74280zl5mhznt2t2e",
			"reward": {
				"denom": "ucmdx",
				"amount": "123"
			}
		},
		{
			"address": "comdex1pv2k843yglwzql3mckpaqepcjqzqwhtz5en9mu",
			"reward": {
				"denom": "ucmdx",
				"amount": "2358"
			}
		},
		{
			"address": "comdex1pvwussxtdv2m2cutv8atgpy40e7g947ugjy86r",
			"reward": {
				"denom": "ucmdx",
				"amount": "8971"
			}
		},
		{
			"address": "comdex1pv03yz7y7atd72xtzzpyfqy8w3eu5rwvmx9yec",
			"reward": {
				"denom": "ucmdx",
				"amount": "1785"
			}
		},
		{
			"address": "comdex1pvhgy9cwc789q6tc5l86ec8hmk6xhfdlvxtegd",
			"reward": {
				"denom": "ucmdx",
				"amount": "7090"
			}
		},
		{
			"address": "comdex1pvcn0jjv2mugdevlue9vj2ks7ryjgtcckhvaa9",
			"reward": {
				"denom": "ucmdx",
				"amount": "22025"
			}
		},
		{
			"address": "comdex1pvcm3nx4qtuljps2uu6c7teaqmzr9jvl2wnx3q",
			"reward": {
				"denom": "ucmdx",
				"amount": "8970"
			}
		},
		{
			"address": "comdex1pveqhnemyn88va477avrghqze0t37rxh3499n6",
			"reward": {
				"denom": "ucmdx",
				"amount": "5018"
			}
		},
		{
			"address": "comdex1pvldyu5ueuzxdqgkgh5svck0fpyj27k0hxvdjc",
			"reward": {
				"denom": "ucmdx",
				"amount": "9055"
			}
		},
		{
			"address": "comdex1pdrx0plfp6wuu6adfkus2tpny2kxrua7fsd92z",
			"reward": {
				"denom": "ucmdx",
				"amount": "2656"
			}
		},
		{
			"address": "comdex1pdf4pajs0glat4ql5zdcuv5yv7aqsp0q9jdrql",
			"reward": {
				"denom": "ucmdx",
				"amount": "166"
			}
		},
		{
			"address": "comdex1pd2udttgkcv6wwm6x2w32makuslsf9x4jkcwdz",
			"reward": {
				"denom": "ucmdx",
				"amount": "10209"
			}
		},
		{
			"address": "comdex1pdtg3wnep2sh0v8g26mdueq9lk86amv5s8ztt7",
			"reward": {
				"denom": "ucmdx",
				"amount": "755"
			}
		},
		{
			"address": "comdex1pd046vyzjglfu5896vl04ngr9gf5xzl8ldsls4",
			"reward": {
				"denom": "ucmdx",
				"amount": "12297"
			}
		},
		{
			"address": "comdex1pds0qjfqm3d59a055r7a5k2w0rpd5lhttv63gq",
			"reward": {
				"denom": "ucmdx",
				"amount": "19"
			}
		},
		{
			"address": "comdex1pdsepwmf42mw53tcq70fe2xjfshgkrgn97pkqy",
			"reward": {
				"denom": "ucmdx",
				"amount": "5305"
			}
		},
		{
			"address": "comdex1pd4rgc2pkc0wmrzcyef4dnxvkl9h7eclcuw2ky",
			"reward": {
				"denom": "ucmdx",
				"amount": "33"
			}
		},
		{
			"address": "comdex1pd6phr8wa9x59q8v6wtfp96qkhgu66qml6xrx4",
			"reward": {
				"denom": "ucmdx",
				"amount": "1441"
			}
		},
		{
			"address": "comdex1pd68chcg2yvvseudap643whw86wvl9majwlt47",
			"reward": {
				"denom": "ucmdx",
				"amount": "2103887"
			}
		},
		{
			"address": "comdex1pd6t5u3wy2ygwwddcjs0g44akepgzr9cvwm6wk",
			"reward": {
				"denom": "ucmdx",
				"amount": "173"
			}
		},
		{
			"address": "comdex1pwpr44q0jkr49x3ed9ennrlw053v7jywzjl5yx",
			"reward": {
				"denom": "ucmdx",
				"amount": "1012"
			}
		},
		{
			"address": "comdex1pwp8vcg6jxa6r5t5dcwwmxjnf8yux0uppmec7e",
			"reward": {
				"denom": "ucmdx",
				"amount": "1244"
			}
		},
		{
			"address": "comdex1pwza48p44cu9caawe47kvlzqcjzcuhkmflehdx",
			"reward": {
				"denom": "ucmdx",
				"amount": "7200"
			}
		},
		{
			"address": "comdex1pw8nh9c9g9792jt88ncas9znsgn5xrxyu9jgqk",
			"reward": {
				"denom": "ucmdx",
				"amount": "185"
			}
		},
		{
			"address": "comdex1pwvhq9pgnpa447vn8jqfuvnesdmqljvzvuvn73",
			"reward": {
				"denom": "ucmdx",
				"amount": "3901"
			}
		},
		{
			"address": "comdex1pwdnwrdh2uplcur4cw4r2846mvj3jl4axydkt9",
			"reward": {
				"denom": "ucmdx",
				"amount": "64566"
			}
		},
		{
			"address": "comdex1pwj2k29mq2u4sfm89hv5prhuajw320tfevyvew",
			"reward": {
				"denom": "ucmdx",
				"amount": "7493"
			}
		},
		{
			"address": "comdex1pwn8j7n98l4j2vt5cmjk57knspl8xdt60xdjpt",
			"reward": {
				"denom": "ucmdx",
				"amount": "961"
			}
		},
		{
			"address": "comdex1pw5tkkajgjg9xdtkh05v3ks6fcdlwv3ychy9p3",
			"reward": {
				"denom": "ucmdx",
				"amount": "33349"
			}
		},
		{
			"address": "comdex1pw4j58t52qe265q3j99h3fpgz5zkxej6nt3jt2",
			"reward": {
				"denom": "ucmdx",
				"amount": "71304"
			}
		},
		{
			"address": "comdex1pwhhezdspshfc30mgcctm9t7kd606rt66rrgp7",
			"reward": {
				"denom": "ucmdx",
				"amount": "205"
			}
		},
		{
			"address": "comdex1pwchad2ymm9tvpfmgr59sptkm3ate3p6eac2hs",
			"reward": {
				"denom": "ucmdx",
				"amount": "18127"
			}
		},
		{
			"address": "comdex1pwakle67923k9crvyr6f5h886g0rhm72gncnm4",
			"reward": {
				"denom": "ucmdx",
				"amount": "1617"
			}
		},
		{
			"address": "comdex1p0qgznl2fwtmx3jqnua793ppwnq5jwf2r9y920",
			"reward": {
				"denom": "ucmdx",
				"amount": "4099"
			}
		},
		{
			"address": "comdex1p0r9v25c8tjtpdzfcncyma6qygkvxknqkguxcq",
			"reward": {
				"denom": "ucmdx",
				"amount": "284"
			}
		},
		{
			"address": "comdex1p0xc8vz6m72velt8xera2quz3xs9ee6c7rddsy",
			"reward": {
				"denom": "ucmdx",
				"amount": "3238"
			}
		},
		{
			"address": "comdex1p02gahyrqn68weajhy7gaxyx4gzlvfn66p4cqq",
			"reward": {
				"denom": "ucmdx",
				"amount": "25314"
			}
		},
		{
			"address": "comdex1p0tsnu2hj4u86cy8f84ns9wykhg99mm34cee69",
			"reward": {
				"denom": "ucmdx",
				"amount": "7120"
			}
		},
		{
			"address": "comdex1p0s9rqfrmgkre8ueq72qqu65569hz6ttj0rkrg",
			"reward": {
				"denom": "ucmdx",
				"amount": "71725"
			}
		},
		{
			"address": "comdex1p0jn9hurhn9xm36gceltxn3hewn8xs42ylm99z",
			"reward": {
				"denom": "ucmdx",
				"amount": "2003"
			}
		},
		{
			"address": "comdex1p04amnegradn6gumvva6hg6yzzds38raxu5ru8",
			"reward": {
				"denom": "ucmdx",
				"amount": "22699"
			}
		},
		{
			"address": "comdex1p0c0d7faqjlhw5y07kxwppatahu22wdfd5v387",
			"reward": {
				"denom": "ucmdx",
				"amount": "16"
			}
		},
		{
			"address": "comdex1p0c4l80pjfrlfn7t50a3tsqcl8lc3h4m5mz7jh",
			"reward": {
				"denom": "ucmdx",
				"amount": "579"
			}
		},
		{
			"address": "comdex1p0e65wwvve7l38enukd83u7atavz9ctezehlx7",
			"reward": {
				"denom": "ucmdx",
				"amount": "180"
			}
		},
		{
			"address": "comdex1p0aufl58nh8vwp97z6wrdcep45rhszwzdd57r7",
			"reward": {
				"denom": "ucmdx",
				"amount": "9275"
			}
		},
		{
			"address": "comdex1p07vy4ncjj2e2yjmyn76z76kwe9893uv33ulra",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1p07d3q3ltyr3r0w5wv8xkz6853dc45s72cq2v2",
			"reward": {
				"denom": "ucmdx",
				"amount": "62953"
			}
		},
		{
			"address": "comdex1pszxqfv2qez5d8j0xh6svjdcpm7zltc5msg3e4",
			"reward": {
				"denom": "ucmdx",
				"amount": "186363"
			}
		},
		{
			"address": "comdex1pszd8q6434ljva8d2w4wdqr5envpkcpcdp3wkp",
			"reward": {
				"denom": "ucmdx",
				"amount": "8649"
			}
		},
		{
			"address": "comdex1psz0dvuka0du625h8hfwy0gdlu2etmknpazqux",
			"reward": {
				"denom": "ucmdx",
				"amount": "62"
			}
		},
		{
			"address": "comdex1psxgymvlca9w876c2dkw6yddxecmwuuaq62e5e",
			"reward": {
				"denom": "ucmdx",
				"amount": "37577"
			}
		},
		{
			"address": "comdex1psg256sg2dqnzrvvq3ge9g7seh2thyystl2htl",
			"reward": {
				"denom": "ucmdx",
				"amount": "2506"
			}
		},
		{
			"address": "comdex1ps238zq9wq6wr86cs56wflyv04tqlan6xg9u9m",
			"reward": {
				"denom": "ucmdx",
				"amount": "1799"
			}
		},
		{
			"address": "comdex1pst6axf05m2cqdevfd4ka5dmdx53prax04p7c8",
			"reward": {
				"denom": "ucmdx",
				"amount": "24"
			}
		},
		{
			"address": "comdex1psvusgarsvwpv9mqggddlfwwn2mendzcz82qxm",
			"reward": {
				"denom": "ucmdx",
				"amount": "5714"
			}
		},
		{
			"address": "comdex1psdhnw746sej28ruup7uyk8ym227zl6emlysce",
			"reward": {
				"denom": "ucmdx",
				"amount": "64359"
			}
		},
		{
			"address": "comdex1psdm7lrc29hlzj5axv5tlrt8huj8lmgh823txj",
			"reward": {
				"denom": "ucmdx",
				"amount": "1284"
			}
		},
		{
			"address": "comdex1pswe4zg8agsn8vjqtp5y7p69c2uu9znc5gsy52",
			"reward": {
				"denom": "ucmdx",
				"amount": "1951"
			}
		},
		{
			"address": "comdex1ps02ssvsat8mmwl8ly0hmyx57e0zxqfrpj0r3u",
			"reward": {
				"denom": "ucmdx",
				"amount": "36592"
			}
		},
		{
			"address": "comdex1pssf4s2m70ka74c29jrf922sm5mp0mh0r7wwdw",
			"reward": {
				"denom": "ucmdx",
				"amount": "24374"
			}
		},
		{
			"address": "comdex1psnx9kds82xjak2uuctwr8emqsfwpy5wlp4kuj",
			"reward": {
				"denom": "ucmdx",
				"amount": "746"
			}
		},
		{
			"address": "comdex1psnkjlqt3xhqgww6etluzk8l0lsey3cayw7gaa",
			"reward": {
				"denom": "ucmdx",
				"amount": "530"
			}
		},
		{
			"address": "comdex1ps5ys77fzcpph5670a29vwjc880p7z4gyswumu",
			"reward": {
				"denom": "ucmdx",
				"amount": "495"
			}
		},
		{
			"address": "comdex1ps4ppkfdt4vllr055vdaz98mz9mucfcgp3flc8",
			"reward": {
				"denom": "ucmdx",
				"amount": "1327"
			}
		},
		{
			"address": "comdex1ps4vckxcjgcyg7zqwkshjr62ca06kvpyvhe498",
			"reward": {
				"denom": "ucmdx",
				"amount": "9997"
			}
		},
		{
			"address": "comdex1ps6f67034fp22nqpvkur5pry039nfjputfv0sg",
			"reward": {
				"denom": "ucmdx",
				"amount": "1902"
			}
		},
		{
			"address": "comdex1ps6hkgfv6y0t7waumj9gcym5hvnypd4f3ete9n",
			"reward": {
				"denom": "ucmdx",
				"amount": "702"
			}
		},
		{
			"address": "comdex1psmcg6xp93leljuurf5z8vatxw7ay8dtz9sqjx",
			"reward": {
				"denom": "ucmdx",
				"amount": "21"
			}
		},
		{
			"address": "comdex1psaav9mslme6kcsspqjemjfn7m0acwyz2nqqra",
			"reward": {
				"denom": "ucmdx",
				"amount": "833"
			}
		},
		{
			"address": "comdex1pslzpwmgxp7tvhzexfpjz7v4zcsxdn2zwks5sj",
			"reward": {
				"denom": "ucmdx",
				"amount": "23569"
			}
		},
		{
			"address": "comdex1psl29hcvyd5svrn5v4pptcltc6pstxg49hfyew",
			"reward": {
				"denom": "ucmdx",
				"amount": "3292"
			}
		},
		{
			"address": "comdex1pslmcukgd2t8fw99xvpqrtxfykf3jz4mqqas2y",
			"reward": {
				"denom": "ucmdx",
				"amount": "15381"
			}
		},
		{
			"address": "comdex1p3qvyssf056epz799c64zcsaqeauku7u0hqpdr",
			"reward": {
				"denom": "ucmdx",
				"amount": "2649"
			}
		},
		{
			"address": "comdex1p3z4l0gcnr7yywxcuk5hccq64srpqajw8wwz76",
			"reward": {
				"denom": "ucmdx",
				"amount": "2877"
			}
		},
		{
			"address": "comdex1p3rjgws5fd8a9f6f3tje95js2e7x0ws00rujqu",
			"reward": {
				"denom": "ucmdx",
				"amount": "24610"
			}
		},
		{
			"address": "comdex1p3y6svk6p8wrul58al7q5xmmd30z2qaqxlvse2",
			"reward": {
				"denom": "ucmdx",
				"amount": "528"
			}
		},
		{
			"address": "comdex1p383tst9addrgmf320a6zckvsc3sfh5h4ulazz",
			"reward": {
				"denom": "ucmdx",
				"amount": "8842"
			}
		},
		{
			"address": "comdex1p3gz6qxtmjpz97zsa50s709znypsyhyn5rz2hl",
			"reward": {
				"denom": "ucmdx",
				"amount": "3614"
			}
		},
		{
			"address": "comdex1p3fa5a5ncjk8jrlna2ld9njy6jveu9j05mtauf",
			"reward": {
				"denom": "ucmdx",
				"amount": "6458"
			}
		},
		{
			"address": "comdex1p3226qgrcqmrwp0rh8whydzyx06mdj4tqz0jw3",
			"reward": {
				"denom": "ucmdx",
				"amount": "1564"
			}
		},
		{
			"address": "comdex1p32kfrpgy3g2lc7q2nl2rcdhwfuuzuak6ecpjn",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1p32msm7k6x3mrj5g0vk3zagjhs72nzrqyje43u",
			"reward": {
				"denom": "ucmdx",
				"amount": "12375"
			}
		},
		{
			"address": "comdex1p3vzy48erskut3z36wgkquthplvyq7h6lxeyel",
			"reward": {
				"denom": "ucmdx",
				"amount": "19618"
			}
		},
		{
			"address": "comdex1p3whfmm83r5eepfjapuhggzfvhaqm72r07xm8k",
			"reward": {
				"denom": "ucmdx",
				"amount": "895"
			}
		},
		{
			"address": "comdex1p3s0atpw6uhnzt63tpylmszhzxw8qnagf6muak",
			"reward": {
				"denom": "ucmdx",
				"amount": "2863"
			}
		},
		{
			"address": "comdex1p3elewj6tz2t5t30s09r6lqswu8hpk5sljzg0m",
			"reward": {
				"denom": "ucmdx",
				"amount": "6100"
			}
		},
		{
			"address": "comdex1p3mst4t3qyd5yyfgxdl5u58kf2xx0ke4dsf30a",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1p372ux8phmq0a055qwjd2rry8035m25xn49f57",
			"reward": {
				"denom": "ucmdx",
				"amount": "2671"
			}
		},
		{
			"address": "comdex1pjr9lfl5v3049pu62qhr6fz38ujhuxrhl885zd",
			"reward": {
				"denom": "ucmdx",
				"amount": "5669"
			}
		},
		{
			"address": "comdex1pjy32mrrtqz0tghpfrquycnennux8femvscj02",
			"reward": {
				"denom": "ucmdx",
				"amount": "14360"
			}
		},
		{
			"address": "comdex1pj9pcccv2ufdwqc25vq297j5qk7h45zjp92dtu",
			"reward": {
				"denom": "ucmdx",
				"amount": "123"
			}
		},
		{
			"address": "comdex1pj9rxy7nf2xglp5wzjl933ywn80848uq6jn5fp",
			"reward": {
				"denom": "ucmdx",
				"amount": "178"
			}
		},
		{
			"address": "comdex1pjxcp00m8ewgdpjwvhvj6gryyme8nyza4jwdju",
			"reward": {
				"denom": "ucmdx",
				"amount": "18198"
			}
		},
		{
			"address": "comdex1pjx7nsd7w9n4w7k8qexkttrgmxwusk8uf2u2dp",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1pj8s4q2u6np2gfh3sdcacckzzx68nqslr6kjec",
			"reward": {
				"denom": "ucmdx",
				"amount": "14346"
			}
		},
		{
			"address": "comdex1pjgvlykn6pm2rd5z5n3k0v8tzndz0z2s6feuqa",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1pjgkt8r887dssr4gfnydp0qc4kuyjwss83lp0g",
			"reward": {
				"denom": "ucmdx",
				"amount": "495"
			}
		},
		{
			"address": "comdex1pj2y9spq3s9akuvalpgp3sk9tyea8sjext4rjz",
			"reward": {
				"denom": "ucmdx",
				"amount": "17384"
			}
		},
		{
			"address": "comdex1pjv7c0vtjmyvqtfvq05qw3z3paepv7hlkwwfw5",
			"reward": {
				"denom": "ucmdx",
				"amount": "2104"
			}
		},
		{
			"address": "comdex1pjd2d2uft9kf6rwpv4xl67jf8qlhetwwe5uh40",
			"reward": {
				"denom": "ucmdx",
				"amount": "311"
			}
		},
		{
			"address": "comdex1pj06vncdufcn0k3sjj2kxl5n9zazct5ahnxlh8",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1pjufxjljq9dv679rwjkuhh2anllslav7jmy3wv",
			"reward": {
				"denom": "ucmdx",
				"amount": "1424"
			}
		},
		{
			"address": "comdex1pnq0qusxfhfrkm8th04mv8d7dmzg8762yrj7y7",
			"reward": {
				"denom": "ucmdx",
				"amount": "4979"
			}
		},
		{
			"address": "comdex1pnqn85ssx3den6pvyxk69z6jkkatrtd4hl95c5",
			"reward": {
				"denom": "ucmdx",
				"amount": "640"
			}
		},
		{
			"address": "comdex1pnzy54q6d2xxftryg5tw3fvpxx8tf6rty4ac5g",
			"reward": {
				"denom": "ucmdx",
				"amount": "3166"
			}
		},
		{
			"address": "comdex1pnx4s7qtrghv06rzp05cv9u8vjp6yztl2dsu67",
			"reward": {
				"denom": "ucmdx",
				"amount": "10089"
			}
		},
		{
			"address": "comdex1pngexrfz0pk99rfex9ck6vfdjryhr48rmn40f0",
			"reward": {
				"denom": "ucmdx",
				"amount": "498"
			}
		},
		{
			"address": "comdex1pnt7glcvahp0qtdltlggxqjek3ed587afqsvx7",
			"reward": {
				"denom": "ucmdx",
				"amount": "285"
			}
		},
		{
			"address": "comdex1pn075qddqdj34xwxd985047ru08lratnn8tz8q",
			"reward": {
				"denom": "ucmdx",
				"amount": "3391"
			}
		},
		{
			"address": "comdex1pnspa8qw4p7m7qkr4pd9tpuzjcy5uuxa0c9td3",
			"reward": {
				"denom": "ucmdx",
				"amount": "4466"
			}
		},
		{
			"address": "comdex1pnnkklufzhx9w4muj3dqpdxrqa2dhhml5j7eyg",
			"reward": {
				"denom": "ucmdx",
				"amount": "440"
			}
		},
		{
			"address": "comdex1pn5wktqrcw6leve0qsmzc3fh09avf3p3tpl2qk",
			"reward": {
				"denom": "ucmdx",
				"amount": "28"
			}
		},
		{
			"address": "comdex1pnk5lmfg7n4tgma646lx89k2q25ysjf67ktqft",
			"reward": {
				"denom": "ucmdx",
				"amount": "1812"
			}
		},
		{
			"address": "comdex1pncdulpncl2qazdexc27vl5nsw8fgd65amqpxy",
			"reward": {
				"denom": "ucmdx",
				"amount": "8206"
			}
		},
		{
			"address": "comdex1pneu3u2n24rea05nx7vp4kkhmpmyxnuhwqzdru",
			"reward": {
				"denom": "ucmdx",
				"amount": "30972"
			}
		},
		{
			"address": "comdex1pn6gcdkcsu44u3g5qw5586nvwl59hwkxy7jhcx",
			"reward": {
				"denom": "ucmdx",
				"amount": "26429"
			}
		},
		{
			"address": "comdex1pnmy7kqjxfz2kenffhcl0desyaccl6y564696k",
			"reward": {
				"denom": "ucmdx",
				"amount": "614"
			}
		},
		{
			"address": "comdex1pnm92lf6tvpsqxkz9y92naduwwkymfz3jr837x",
			"reward": {
				"denom": "ucmdx",
				"amount": "290"
			}
		},
		{
			"address": "comdex1pnml5wlwl94trvfhvq308jj7ptxvwz6r4v2s7z",
			"reward": {
				"denom": "ucmdx",
				"amount": "88"
			}
		},
		{
			"address": "comdex1p5ygt5ay4s4n5hfhcx7v7evre6cfvfrnslahtq",
			"reward": {
				"denom": "ucmdx",
				"amount": "1815"
			}
		},
		{
			"address": "comdex1p59lf9z0qzkj9v3s7w22tfsedslaxvqzmylxcr",
			"reward": {
				"denom": "ucmdx",
				"amount": "16"
			}
		},
		{
			"address": "comdex1p5g6ff5942yvpmkjtjglr2t9ldx3ts0g2m2zct",
			"reward": {
				"denom": "ucmdx",
				"amount": "1237"
			}
		},
		{
			"address": "comdex1p5fft5dv2fmn2pljkzfg5c3lg9x3qxhxunw9ld",
			"reward": {
				"denom": "ucmdx",
				"amount": "85"
			}
		},
		{
			"address": "comdex1p52fqvqe0pl59pwzndw4r7gkk486ay55t9mvwg",
			"reward": {
				"denom": "ucmdx",
				"amount": "73926"
			}
		},
		{
			"address": "comdex1p5vkazhh6yjxlrgc40n604rwdwpkv5r0taf6tp",
			"reward": {
				"denom": "ucmdx",
				"amount": "944"
			}
		},
		{
			"address": "comdex1p5w0r67kp6wlpc3e2cvnd4ql4yfhm7zzw4yvzl",
			"reward": {
				"denom": "ucmdx",
				"amount": "1425"
			}
		},
		{
			"address": "comdex1p50suef29mtxwzv4mst7rzqxd3ycjc7lmhvu4v",
			"reward": {
				"denom": "ucmdx",
				"amount": "896"
			}
		},
		{
			"address": "comdex1p5397g5djutq3kmc7tlx5j6ucwpelavvsathuv",
			"reward": {
				"denom": "ucmdx",
				"amount": "171"
			}
		},
		{
			"address": "comdex1p535fpjqehfs8nw30fku9fyzejk9a3wu946ytq",
			"reward": {
				"denom": "ucmdx",
				"amount": "43"
			}
		},
		{
			"address": "comdex1p5cnp7qjwv25jp06gwqte3af0x6mkusft9fmpq",
			"reward": {
				"denom": "ucmdx",
				"amount": "1431"
			}
		},
		{
			"address": "comdex1p5cnyt8as23jdjsxfkvnerafj2d53wlufn46tk",
			"reward": {
				"denom": "ucmdx",
				"amount": "29531"
			}
		},
		{
			"address": "comdex1p5als07xw4h4an92wgvpnajasll9p85r0cw9ls",
			"reward": {
				"denom": "ucmdx",
				"amount": "877"
			}
		},
		{
			"address": "comdex1p5leww2pmwddsre6z4qlug6g96pqmswausjlcm",
			"reward": {
				"denom": "ucmdx",
				"amount": "3950"
			}
		},
		{
			"address": "comdex1p4r5lv6pdchdnxw9ekqjwa08kzefrcrgy0nccr",
			"reward": {
				"denom": "ucmdx",
				"amount": "8914"
			}
		},
		{
			"address": "comdex1p493mra2kwlwh8pg525alxsu3efpwc5842fyc8",
			"reward": {
				"denom": "ucmdx",
				"amount": "7177"
			}
		},
		{
			"address": "comdex1p4xptl8d88uhlm8kr5060nvk7wg864zlvutgux",
			"reward": {
				"denom": "ucmdx",
				"amount": "6922"
			}
		},
		{
			"address": "comdex1p4vq3nvkk79n3gzp0f5wkyjjptwhsfw9flaeqy",
			"reward": {
				"denom": "ucmdx",
				"amount": "100916"
			}
		},
		{
			"address": "comdex1p4due0lp6f463e25dy2mhwfyrm6r7z2jg0u3pg",
			"reward": {
				"denom": "ucmdx",
				"amount": "1768"
			}
		},
		{
			"address": "comdex1p4wp9fy4kamsd7acz99w042xkgthyz40w7x2hl",
			"reward": {
				"denom": "ucmdx",
				"amount": "546"
			}
		},
		{
			"address": "comdex1p43ps3kx3hyqzsnuljjvudtv8qmmf5np4sc2wv",
			"reward": {
				"denom": "ucmdx",
				"amount": "2021"
			}
		},
		{
			"address": "comdex1p43srkxlpaqa6nn3udr43r5a34cuj5vfakn6ag",
			"reward": {
				"denom": "ucmdx",
				"amount": "1003690"
			}
		},
		{
			"address": "comdex1p4nyj2wcm8qctxfs55leqtj08xaj3tmu5xy8up",
			"reward": {
				"denom": "ucmdx",
				"amount": "5823"
			}
		},
		{
			"address": "comdex1p4kx0cfsj3qxaje65fghn058pkquplhtw8aynt",
			"reward": {
				"denom": "ucmdx",
				"amount": "11912"
			}
		},
		{
			"address": "comdex1p4kajj7nv2lvapcfx5tfsxmke85qrvsat0uq3w",
			"reward": {
				"denom": "ucmdx",
				"amount": "97831"
			}
		},
		{
			"address": "comdex1p4eh7hpd3ejt7a3pjz7mq4j0fjvf3jpvqsdhl8",
			"reward": {
				"denom": "ucmdx",
				"amount": "10775"
			}
		},
		{
			"address": "comdex1p4m0yvvn7m66a66dt0yup4harpdv7l2gh6rpf2",
			"reward": {
				"denom": "ucmdx",
				"amount": "316"
			}
		},
		{
			"address": "comdex1p4unv3rarrqq9valkj3puchqr77csrgutfwx4x",
			"reward": {
				"denom": "ucmdx",
				"amount": "3723"
			}
		},
		{
			"address": "comdex1p4a8a5m7345jkeuqgqfdytfka4dyzawhyzwwpp",
			"reward": {
				"denom": "ucmdx",
				"amount": "2015"
			}
		},
		{
			"address": "comdex1p479c8cgfz8al6hyvlxfev7j9dcgrv35ft6fcl",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1p4lunqxy3pjuzulwnk7ywqjvljv35qrdmvx6e3",
			"reward": {
				"denom": "ucmdx",
				"amount": "1765"
			}
		},
		{
			"address": "comdex1pkz5setcxsc35py3jtakpv4dr4rkkvxt23xdjv",
			"reward": {
				"denom": "ucmdx",
				"amount": "51347"
			}
		},
		{
			"address": "comdex1pkxrg43lqn24zssus7wyz8emvtg2uksxcg4ytl",
			"reward": {
				"denom": "ucmdx",
				"amount": "12441"
			}
		},
		{
			"address": "comdex1pkgj687y29u5zw3u27l39ccv5sefczdszanxpg",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1pk29uy4ha5djpedtzswzwe027769xzljtdtz0k",
			"reward": {
				"denom": "ucmdx",
				"amount": "1254"
			}
		},
		{
			"address": "comdex1pkwq44a7x7ta5lh3czwr837psv4l9ta3gnl96d",
			"reward": {
				"denom": "ucmdx",
				"amount": "165"
			}
		},
		{
			"address": "comdex1pk0d82scu4aud4s3yzschn8t3ghc67l3uv065m",
			"reward": {
				"denom": "ucmdx",
				"amount": "5943"
			}
		},
		{
			"address": "comdex1pk0n426dfe0d65rlpr90q4v6vhpcwl293x4zw5",
			"reward": {
				"denom": "ucmdx",
				"amount": "90648"
			}
		},
		{
			"address": "comdex1pk07jt3d2j5rqenlnpl4asdakp9us3nzmwuvs4",
			"reward": {
				"denom": "ucmdx",
				"amount": "31696"
			}
		},
		{
			"address": "comdex1pk3lemu5u2np669aalce6k6az35tyqyrsfmzzj",
			"reward": {
				"denom": "ucmdx",
				"amount": "3567"
			}
		},
		{
			"address": "comdex1pkj8kjhdz84grm7d3nx2jpynmau7h3l8whjycc",
			"reward": {
				"denom": "ucmdx",
				"amount": "410"
			}
		},
		{
			"address": "comdex1pkk0z57pcvkzrn7uxq9u7ztvxc5lfcayznufzy",
			"reward": {
				"denom": "ucmdx",
				"amount": "2831"
			}
		},
		{
			"address": "comdex1pkcv0s35yx9u0e4jr7tkqp0edrv7ghuq4s9tyj",
			"reward": {
				"denom": "ucmdx",
				"amount": "241948"
			}
		},
		{
			"address": "comdex1pkc05h2df5g0jvkmtsph3myryl2f30waspfu8z",
			"reward": {
				"denom": "ucmdx",
				"amount": "11454"
			}
		},
		{
			"address": "comdex1pkm998xw4vtluu0k4agya7e6zf4rwxsl5zu4z9",
			"reward": {
				"denom": "ucmdx",
				"amount": "124832"
			}
		},
		{
			"address": "comdex1pku7uwddkugqgzel8qrueh97v5ggvtt6eylfmy",
			"reward": {
				"denom": "ucmdx",
				"amount": "140153"
			}
		},
		{
			"address": "comdex1pk7s0p20eueduzcmnxagpr6vfj84cxfkvusl7d",
			"reward": {
				"denom": "ucmdx",
				"amount": "27370"
			}
		},
		{
			"address": "comdex1phqu8c7kq9vjdlp5jr0xtca7hefd5zr3p2ks6x",
			"reward": {
				"denom": "ucmdx",
				"amount": "13938"
			}
		},
		{
			"address": "comdex1phpdgyw0hupdu5fs66xxtwcf8yhrqjnrvqs4t4",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1phr97ne5l6fufgqvcrtekuurmx4df02j98elkd",
			"reward": {
				"denom": "ucmdx",
				"amount": "1513"
			}
		},
		{
			"address": "comdex1phrw03kmjl34pjkvv5zl5lrskhdp6374wvg27h",
			"reward": {
				"denom": "ucmdx",
				"amount": "13673"
			}
		},
		{
			"address": "comdex1ph9ejdjfuvluqr2543yw24z9f0z5nd4zk2ydz6",
			"reward": {
				"denom": "ucmdx",
				"amount": "146"
			}
		},
		{
			"address": "comdex1phxzzqftrxyx6f84ge8xyzpy5mgns5cqzg3y72",
			"reward": {
				"denom": "ucmdx",
				"amount": "818"
			}
		},
		{
			"address": "comdex1phxt5kdjkjpewsgu4sl3rshth2a6feuu6ndqgr",
			"reward": {
				"denom": "ucmdx",
				"amount": "1186"
			}
		},
		{
			"address": "comdex1phxkjqdm7xc4hseg6ke0waejplm5u2t8ah4603",
			"reward": {
				"denom": "ucmdx",
				"amount": "185"
			}
		},
		{
			"address": "comdex1ph890tyxqsycteuv8sg6t7056d94zd5vvlwws8",
			"reward": {
				"denom": "ucmdx",
				"amount": "2902"
			}
		},
		{
			"address": "comdex1ph2j5knpu3032r4axpqz9c9nphn6xpm38989gx",
			"reward": {
				"denom": "ucmdx",
				"amount": "30280"
			}
		},
		{
			"address": "comdex1phtedlu40slfvmdne9pkn8ql6c54hjrfhn503m",
			"reward": {
				"denom": "ucmdx",
				"amount": "7101"
			}
		},
		{
			"address": "comdex1phvddcjw2ng6ma6e46xk6k645wa52d6af42rqt",
			"reward": {
				"denom": "ucmdx",
				"amount": "155"
			}
		},
		{
			"address": "comdex1phv0ek7wjq63nfll7fy9673p9m8k8zkyhv8wnr",
			"reward": {
				"denom": "ucmdx",
				"amount": "616"
			}
		},
		{
			"address": "comdex1phw4wfqhcl4y9az2zhvfxphyhkg9u24cwx2eng",
			"reward": {
				"denom": "ucmdx",
				"amount": "6188"
			}
		},
		{
			"address": "comdex1ph0qvwfxqw0xrjwzsndy4ll3tvyv7g2cn8vrvu",
			"reward": {
				"denom": "ucmdx",
				"amount": "198"
			}
		},
		{
			"address": "comdex1ph3v5taxm86aq5v7yqpftmr3sgjss0fdz2q6kh",
			"reward": {
				"denom": "ucmdx",
				"amount": "574745"
			}
		},
		{
			"address": "comdex1phjay9jyx7wty20qpts2t78eqq22dvaq4vkx2d",
			"reward": {
				"denom": "ucmdx",
				"amount": "6494"
			}
		},
		{
			"address": "comdex1phnpmu3w7q82hfw8g2ykrus2s00yac8yfwqhza",
			"reward": {
				"denom": "ucmdx",
				"amount": "1439"
			}
		},
		{
			"address": "comdex1phk8v0qs56dhh5el62xee3dflczc98vwd3azwq",
			"reward": {
				"denom": "ucmdx",
				"amount": "445"
			}
		},
		{
			"address": "comdex1phcmh6rt32n05yfujxfphvv4y2a06nfznde355",
			"reward": {
				"denom": "ucmdx",
				"amount": "1390"
			}
		},
		{
			"address": "comdex1phukleuhwjgjynjzff0m2uu90r98t0fleu9g25",
			"reward": {
				"denom": "ucmdx",
				"amount": "179"
			}
		},
		{
			"address": "comdex1ph7gwddzsmcvlv5nn6fn9qpg3q7vewldyjngax",
			"reward": {
				"denom": "ucmdx",
				"amount": "1449"
			}
		},
		{
			"address": "comdex1pczddfv3fdzl9gchc62rwh3ua0shqtlgpv2pe2",
			"reward": {
				"denom": "ucmdx",
				"amount": "1426"
			}
		},
		{
			"address": "comdex1pc95zh8ceqst7cmwt8seqdeyfuh8rjaekslf8t",
			"reward": {
				"denom": "ucmdx",
				"amount": "1025"
			}
		},
		{
			"address": "comdex1pcx9wwzcdcdyjgpyeaufu5x44qwmf7es65359q",
			"reward": {
				"denom": "ucmdx",
				"amount": "4049"
			}
		},
		{
			"address": "comdex1pc8w4rj7ecq80z3s0qckfmy8n2n6u9s74qwqcy",
			"reward": {
				"denom": "ucmdx",
				"amount": "182"
			}
		},
		{
			"address": "comdex1pcwd9shqsuq7aj005kwtnlrcsfzq4p2uaytg3d",
			"reward": {
				"denom": "ucmdx",
				"amount": "28653"
			}
		},
		{
			"address": "comdex1pc0hyuwxyy63phpgh6ewau05ues7njhqdd80x2",
			"reward": {
				"denom": "ucmdx",
				"amount": "150"
			}
		},
		{
			"address": "comdex1pcsaq68d70ycvdxxv3tdl6uvamhkg9nvngvjmk",
			"reward": {
				"denom": "ucmdx",
				"amount": "13793"
			}
		},
		{
			"address": "comdex1pcs7h6egtmctr3ahahq93ppawwls7srz7aprpa",
			"reward": {
				"denom": "ucmdx",
				"amount": "417"
			}
		},
		{
			"address": "comdex1pcjf6gpl7vdw4t2aufwqtlpjwe05kgks68e7vh",
			"reward": {
				"denom": "ucmdx",
				"amount": "10708"
			}
		},
		{
			"address": "comdex1pcn785ct0da6xkp0er3r6tcl8s89s2l55mmwyx",
			"reward": {
				"denom": "ucmdx",
				"amount": "445"
			}
		},
		{
			"address": "comdex1pckq70nv36rjlsxpkqwf75yv963wg65gw827d0",
			"reward": {
				"denom": "ucmdx",
				"amount": "40107"
			}
		},
		{
			"address": "comdex1pccavde92gay3fq9c7e4ne828w0a2p22tud2va",
			"reward": {
				"denom": "ucmdx",
				"amount": "7555"
			}
		},
		{
			"address": "comdex1pceq7npe9ylnyrt02nzskt54szgxzwtjpzc8z6",
			"reward": {
				"denom": "ucmdx",
				"amount": "156480"
			}
		},
		{
			"address": "comdex1pc669ty4r0nt4navrnfevl2ss9ww8tyuxq850j",
			"reward": {
				"denom": "ucmdx",
				"amount": "5065"
			}
		},
		{
			"address": "comdex1pcm9emd04qn44yr997m3ywhd29t4te27cyv38a",
			"reward": {
				"denom": "ucmdx",
				"amount": "1234"
			}
		},
		{
			"address": "comdex1pcu8z2umew5mprurpawrekxrf0t7kke9ykjqdy",
			"reward": {
				"denom": "ucmdx",
				"amount": "7269"
			}
		},
		{
			"address": "comdex1pcaz4gzf598ep8jncgsyfgf9gv9fafurt076kg",
			"reward": {
				"denom": "ucmdx",
				"amount": "8143"
			}
		},
		{
			"address": "comdex1pcllsgzsfer97ujrhrep3892y3hmdzst9383g9",
			"reward": {
				"denom": "ucmdx",
				"amount": "42906"
			}
		},
		{
			"address": "comdex1peq96jqmp9anxft7plu9qp4cpz4yr2k29hg90p",
			"reward": {
				"denom": "ucmdx",
				"amount": "163915"
			}
		},
		{
			"address": "comdex1pep8j0pq7lkzhpa280kvk3nytsjv2v2wx8amqw",
			"reward": {
				"denom": "ucmdx",
				"amount": "6799"
			}
		},
		{
			"address": "comdex1pez2f5ynw9nxm3027a4rqlcxaesr4t7rgmmcz3",
			"reward": {
				"denom": "ucmdx",
				"amount": "19189"
			}
		},
		{
			"address": "comdex1persgmglpz0spklwqfef49p0qn3weryyf0tk3f",
			"reward": {
				"denom": "ucmdx",
				"amount": "106"
			}
		},
		{
			"address": "comdex1pe9w8h680jk05qnxqj38w3gfzj2xnatszevkpr",
			"reward": {
				"denom": "ucmdx",
				"amount": "7711"
			}
		},
		{
			"address": "comdex1pexgymrz9u2t4759kngw2vu0s304gn4ewenrq7",
			"reward": {
				"denom": "ucmdx",
				"amount": "11993"
			}
		},
		{
			"address": "comdex1peteje73eklqau66mr7h7rmewmt2vt994h5s2f",
			"reward": {
				"denom": "ucmdx",
				"amount": "14231"
			}
		},
		{
			"address": "comdex1pedckxe0dl4l53e5q8k0ljq3utjcn9fzw7vzvh",
			"reward": {
				"denom": "ucmdx",
				"amount": "57"
			}
		},
		{
			"address": "comdex1pew5c6mwplw35z9jpy6uyzv6ez94074cmzkzu6",
			"reward": {
				"denom": "ucmdx",
				"amount": "11300"
			}
		},
		{
			"address": "comdex1pe0ms48pr4lle5zy6zp7vn9dp0y9zykkruhc49",
			"reward": {
				"denom": "ucmdx",
				"amount": "1768"
			}
		},
		{
			"address": "comdex1pesn7v7ts6ns9strsya8xe9hkgnym2xcmg950w",
			"reward": {
				"denom": "ucmdx",
				"amount": "1954"
			}
		},
		{
			"address": "comdex1penglcm6kls4qxa6pac8t8ntjagugcezrppe00",
			"reward": {
				"denom": "ucmdx",
				"amount": "170"
			}
		},
		{
			"address": "comdex1pehvghr7mzgmlm96yasv08n42plp0h3wln5cte",
			"reward": {
				"denom": "ucmdx",
				"amount": "1517"
			}
		},
		{
			"address": "comdex1pehs6wzz67ss4etk25vhxg2yfvkx2jsjdw9v59",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1pec4e67nev2323el4g08sjl7pkjam465zl0qe5",
			"reward": {
				"denom": "ucmdx",
				"amount": "176"
			}
		},
		{
			"address": "comdex1pelmkne5gtdkqsfawj308mdkup4xf5563tssj6",
			"reward": {
				"denom": "ucmdx",
				"amount": "51292"
			}
		},
		{
			"address": "comdex1p6pp8438x6zk4eupapv564e6mze474sjmhljn9",
			"reward": {
				"denom": "ucmdx",
				"amount": "1530"
			}
		},
		{
			"address": "comdex1p6yvduezggrdxx2hfrdg2zp4xp348dcpvyew42",
			"reward": {
				"denom": "ucmdx",
				"amount": "7727"
			}
		},
		{
			"address": "comdex1p6y3x9zy0d05mnptm6sef0y9ly6xmgxfjtpjvw",
			"reward": {
				"denom": "ucmdx",
				"amount": "1384"
			}
		},
		{
			"address": "comdex1p6yceu9phpzsudtc95lm94d9pskcn7kpwulz65",
			"reward": {
				"denom": "ucmdx",
				"amount": "2252"
			}
		},
		{
			"address": "comdex1p6d3hknwt0wuys96u0pm50kw8ufvt3dlsqec2x",
			"reward": {
				"denom": "ucmdx",
				"amount": "6930"
			}
		},
		{
			"address": "comdex1p6j965hn595mq09h6fhj7rdasphta7yxnhnw32",
			"reward": {
				"denom": "ucmdx",
				"amount": "2549"
			}
		},
		{
			"address": "comdex1p6jhhlg7dtyapqy6apzjwsa2a55mq4j2wwa49y",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1p6kgx5t2da8amkmxxnvuu26fghmz2gzwf24mft",
			"reward": {
				"denom": "ucmdx",
				"amount": "201"
			}
		},
		{
			"address": "comdex1p6el23j6getrxw6al7wyrtxnkf8hmdl9h9u0t9",
			"reward": {
				"denom": "ucmdx",
				"amount": "1439"
			}
		},
		{
			"address": "comdex1p66wnp98x0sms62cnn5kyw3aat8leh594v5nv3",
			"reward": {
				"denom": "ucmdx",
				"amount": "324"
			}
		},
		{
			"address": "comdex1p6ucfx2rj3vpgdtq97df0mgp44qnsqar07hvu9",
			"reward": {
				"denom": "ucmdx",
				"amount": "2853"
			}
		},
		{
			"address": "comdex1p6ulzdwd038srpv40ldakr77pwxncnp7pdn7hm",
			"reward": {
				"denom": "ucmdx",
				"amount": "4144"
			}
		},
		{
			"address": "comdex1p67k5ppggmrhh74zkrh2kkv2uc490ylmhppvp4",
			"reward": {
				"denom": "ucmdx",
				"amount": "3505"
			}
		},
		{
			"address": "comdex1p6ly5kpq8j00fnhm9h3sltfxtr08x4qquyzqgm",
			"reward": {
				"denom": "ucmdx",
				"amount": "1267"
			}
		},
		{
			"address": "comdex1p6lgdq2jtln67tj6g893487lcz8ytmjr5f4p65",
			"reward": {
				"denom": "ucmdx",
				"amount": "5295"
			}
		},
		{
			"address": "comdex1pmpph2cargktq7c0c4ss3utwgtua4kfrx0grdd",
			"reward": {
				"denom": "ucmdx",
				"amount": "317"
			}
		},
		{
			"address": "comdex1pmr80rfvqyn426cug3wdqf92h5lh3fklr6pcqq",
			"reward": {
				"denom": "ucmdx",
				"amount": "16"
			}
		},
		{
			"address": "comdex1pm9320azsu26d094p4uj4rkcarraxwrxl443l4",
			"reward": {
				"denom": "ucmdx",
				"amount": "35428"
			}
		},
		{
			"address": "comdex1pm9cfs8dzzpwerj50vpaxjkem86np4lkjsype7",
			"reward": {
				"denom": "ucmdx",
				"amount": "6486"
			}
		},
		{
			"address": "comdex1pm8jhqcmjhh85u5ja0yqcahk42n5dsgc0fapwe",
			"reward": {
				"denom": "ucmdx",
				"amount": "26250"
			}
		},
		{
			"address": "comdex1pm2qczsm0ney0ev092mewmrg4wd5hxt75pxvwa",
			"reward": {
				"denom": "ucmdx",
				"amount": "20338"
			}
		},
		{
			"address": "comdex1pmwhwww4gfz4c8msg5mwfdzdsgfgjvzelhgzn4",
			"reward": {
				"denom": "ucmdx",
				"amount": "5798"
			}
		},
		{
			"address": "comdex1pm026e9dajvx82vhfxqj4dggj3vp7kph0g8zmw",
			"reward": {
				"denom": "ucmdx",
				"amount": "1420"
			}
		},
		{
			"address": "comdex1pm0k39pulhjh76xx9rhclhe7fezsy8kggquu0r",
			"reward": {
				"denom": "ucmdx",
				"amount": "1040"
			}
		},
		{
			"address": "comdex1pmnvxw6kn08ptlr4rq9xa59wp2p0wzf5fpgrc6",
			"reward": {
				"denom": "ucmdx",
				"amount": "2947"
			}
		},
		{
			"address": "comdex1pmnhst7a6rrrw05ck62jrg5xsnwqxs8s9n5ddd",
			"reward": {
				"denom": "ucmdx",
				"amount": "19"
			}
		},
		{
			"address": "comdex1pmkl4znqff9ylnmet0v38n7nqmtm3gqqwpqdy5",
			"reward": {
				"denom": "ucmdx",
				"amount": "716"
			}
		},
		{
			"address": "comdex1pmhg8m52u7umygtf4erfzn00rvkae46ys8rmm7",
			"reward": {
				"denom": "ucmdx",
				"amount": "8771"
			}
		},
		{
			"address": "comdex1pmccf8ke0pwnfhh3uhldfq9gm689c8u27edv7f",
			"reward": {
				"denom": "ucmdx",
				"amount": "91811"
			}
		},
		{
			"address": "comdex1pml0lnrt5r7ykkaxpm888vshx94274gqk7nxpg",
			"reward": {
				"denom": "ucmdx",
				"amount": "145"
			}
		},
		{
			"address": "comdex1pupf79nl6ka3jaxauram0fmdpv49uurygcj8jn",
			"reward": {
				"denom": "ucmdx",
				"amount": "22203"
			}
		},
		{
			"address": "comdex1purja9dkex682z3j2xxhz4d9hxrpmllthp4j5v",
			"reward": {
				"denom": "ucmdx",
				"amount": "14314"
			}
		},
		{
			"address": "comdex1puyrpr3p86rzclqyzeqk5f3mnm7zrz25svx766",
			"reward": {
				"denom": "ucmdx",
				"amount": "1954"
			}
		},
		{
			"address": "comdex1puy6v8g69p0dkge7ztznp4ma70nlvdzmswglfr",
			"reward": {
				"denom": "ucmdx",
				"amount": "19037"
			}
		},
		{
			"address": "comdex1pu94026qsp7dgjcn0h69jddfj4mt0g9z9a8534",
			"reward": {
				"denom": "ucmdx",
				"amount": "43"
			}
		},
		{
			"address": "comdex1pu8tk2kn2xx7skgq480jhsff2e9c04jdzyetyv",
			"reward": {
				"denom": "ucmdx",
				"amount": "7178"
			}
		},
		{
			"address": "comdex1pugajwawrsahvnm3md2tchkdkwl2zwktlxk364",
			"reward": {
				"denom": "ucmdx",
				"amount": "189"
			}
		},
		{
			"address": "comdex1pufy8qm4qsnxxapugc3l2vvmtpp0vk0s4n76w4",
			"reward": {
				"denom": "ucmdx",
				"amount": "353"
			}
		},
		{
			"address": "comdex1putjg6lzz4u3pv5wyeer7nwmkh6t6lev3u9ge0",
			"reward": {
				"denom": "ucmdx",
				"amount": "17848"
			}
		},
		{
			"address": "comdex1puvsznn7wlrxkf4ulvfeekg2j8nkmmezw3tly3",
			"reward": {
				"denom": "ucmdx",
				"amount": "26478"
			}
		},
		{
			"address": "comdex1pudl0y3tmja7928nwaxqtque7v9x6tdlz77qs2",
			"reward": {
				"denom": "ucmdx",
				"amount": "5312"
			}
		},
		{
			"address": "comdex1puw3maf5sz9j49n8v45aa93y8egcvyequ6l79w",
			"reward": {
				"denom": "ucmdx",
				"amount": "1431"
			}
		},
		{
			"address": "comdex1puwjs3s066jyv4hvgg9q5wyzv70suyjrkwvw0m",
			"reward": {
				"denom": "ucmdx",
				"amount": "6686"
			}
		},
		{
			"address": "comdex1pu0q2fq6s823u7kt3ljtgs8vt58kkrwg8dn0ql",
			"reward": {
				"denom": "ucmdx",
				"amount": "60012"
			}
		},
		{
			"address": "comdex1pu0spyadnykf548l3uw0phmng53t7ufrpsy0f0",
			"reward": {
				"denom": "ucmdx",
				"amount": "1430"
			}
		},
		{
			"address": "comdex1pu3hg0c6f93vachp3kvlycd7338uuydxedhgfp",
			"reward": {
				"denom": "ucmdx",
				"amount": "1884369"
			}
		},
		{
			"address": "comdex1puj6jlrjyvh76shdjmjd4e4q0vqjd44tkhf5cg",
			"reward": {
				"denom": "ucmdx",
				"amount": "4908"
			}
		},
		{
			"address": "comdex1pu5q8ul43t5qng3xx0c7znx3d0lthxq98c7zp5",
			"reward": {
				"denom": "ucmdx",
				"amount": "7866536"
			}
		},
		{
			"address": "comdex1pu5mn99hqt74grmrt5j54cnd8ratp89d2hlayt",
			"reward": {
				"denom": "ucmdx",
				"amount": "208"
			}
		},
		{
			"address": "comdex1puc2nzm2dnjjdfffcm822eqytvxhljww3syfu0",
			"reward": {
				"denom": "ucmdx",
				"amount": "2041"
			}
		},
		{
			"address": "comdex1pumtu9qmxw5d3wfvwfxyfj2syu3e6sw7tczhql",
			"reward": {
				"denom": "ucmdx",
				"amount": "61509"
			}
		},
		{
			"address": "comdex1pua2ne85l4hqjs5haqtu4sez2jgr67ej7ywhep",
			"reward": {
				"denom": "ucmdx",
				"amount": "723"
			}
		},
		{
			"address": "comdex1pulqzulckgnp0v2d2p4pma3rwm3wmtcxqmd7d8",
			"reward": {
				"denom": "ucmdx",
				"amount": "1768"
			}
		},
		{
			"address": "comdex1pape55v5vvt87ckraruyt6krrd8k4k9m24fc95",
			"reward": {
				"denom": "ucmdx",
				"amount": "43"
			}
		},
		{
			"address": "comdex1pa9pnsmavz7tysuje2f6fyrljaklzqa69vcnqk",
			"reward": {
				"denom": "ucmdx",
				"amount": "18936"
			}
		},
		{
			"address": "comdex1pa9rqy8e2wct65v0ku2pzf9egpdynl3cacgpyy",
			"reward": {
				"denom": "ucmdx",
				"amount": "3561"
			}
		},
		{
			"address": "comdex1pa9a9n9ed4y09vpr3w97pdxx4e4l6chexxezse",
			"reward": {
				"denom": "ucmdx",
				"amount": "151"
			}
		},
		{
			"address": "comdex1pawkhzdm0puy8y0w6aqfgnxhcmq7dy9ddne8es",
			"reward": {
				"denom": "ucmdx",
				"amount": "319"
			}
		},
		{
			"address": "comdex1pasnu3wly4gy0a5qwkrc5j99cv0zqj8g3y9ky3",
			"reward": {
				"denom": "ucmdx",
				"amount": "33287"
			}
		},
		{
			"address": "comdex1pa3jyyfuccped495vrkp8kueu2zq4ye55677gw",
			"reward": {
				"denom": "ucmdx",
				"amount": "1767"
			}
		},
		{
			"address": "comdex1pakxhdkmy8swne7wpg4prk60sdj78e8qnnw3ga",
			"reward": {
				"denom": "ucmdx",
				"amount": "2014"
			}
		},
		{
			"address": "comdex1pamd4ywj0jngtzdzu9usdjcf0f5ckcyfy4gl58",
			"reward": {
				"denom": "ucmdx",
				"amount": "49451"
			}
		},
		{
			"address": "comdex1paa2ur6c2z0t8w4qz0ty946gclpckrunmakzdl",
			"reward": {
				"denom": "ucmdx",
				"amount": "1251"
			}
		},
		{
			"address": "comdex1pa7ypzl0e5a7hxct28m7vdmfznazmmnevp0um8",
			"reward": {
				"denom": "ucmdx",
				"amount": "426"
			}
		},
		{
			"address": "comdex1p7pxhn6s25qjrjlqgvxa8q35mp6rl099s48e3g",
			"reward": {
				"denom": "ucmdx",
				"amount": "62046"
			}
		},
		{
			"address": "comdex1p7pv442wq3zchwena3kjx7nd0txdk2kx3jhq2g",
			"reward": {
				"denom": "ucmdx",
				"amount": "176"
			}
		},
		{
			"address": "comdex1p79svu4y6farp9zwz4026wrk0v4yp0fnh4axvp",
			"reward": {
				"denom": "ucmdx",
				"amount": "75"
			}
		},
		{
			"address": "comdex1p7gyng5y74pvrfdnlg2h7yr7vmfs9r6g5p2zj4",
			"reward": {
				"denom": "ucmdx",
				"amount": "2052"
			}
		},
		{
			"address": "comdex1p7g2n095unpdhcnw7et87z6r2usurukjc6a06k",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1p72xfy6scz0f9je8sj5mqcx262jcf3yj6hgdvm",
			"reward": {
				"denom": "ucmdx",
				"amount": "977"
			}
		},
		{
			"address": "comdex1p7td7t0htw596asmh5tveadh2fdz3wmhlhvhtg",
			"reward": {
				"denom": "ucmdx",
				"amount": "25887"
			}
		},
		{
			"address": "comdex1p7dmhx00k85qgyrntm9tpqektdqpwulh25njul",
			"reward": {
				"denom": "ucmdx",
				"amount": "255460"
			}
		},
		{
			"address": "comdex1p7wwlxyf73pc5fr4077hel9s049t72mm26hx2r",
			"reward": {
				"denom": "ucmdx",
				"amount": "14057"
			}
		},
		{
			"address": "comdex1p705srvvq7nsk6m0r3jhd4ud7u8mwenf7j2rmx",
			"reward": {
				"denom": "ucmdx",
				"amount": "353"
			}
		},
		{
			"address": "comdex1p73e2a947253ap2s4acrvzcs453lzg0e9evc85",
			"reward": {
				"denom": "ucmdx",
				"amount": "106"
			}
		},
		{
			"address": "comdex1p7j7dv4fxs80tydjjyyynu0z668r03k46tldxh",
			"reward": {
				"denom": "ucmdx",
				"amount": "26316"
			}
		},
		{
			"address": "comdex1p7km7027hnj2r5mjguk376r2zzwuv2q7syttye",
			"reward": {
				"denom": "ucmdx",
				"amount": "20252"
			}
		},
		{
			"address": "comdex1p7e66wk63j2respzxgy43ha5kyv9gj2dr6l8ez",
			"reward": {
				"denom": "ucmdx",
				"amount": "168"
			}
		},
		{
			"address": "comdex1p76vlex4m84wr8t8687d6t29watej9x7nd2zr9",
			"reward": {
				"denom": "ucmdx",
				"amount": "2048"
			}
		},
		{
			"address": "comdex1plpg5c303smnfgkfnwh4c7e40gxsu32negen7q",
			"reward": {
				"denom": "ucmdx",
				"amount": "646"
			}
		},
		{
			"address": "comdex1plz08zdlp36c8lmpydt5akdx0eulc6ejpqjl9n",
			"reward": {
				"denom": "ucmdx",
				"amount": "352"
			}
		},
		{
			"address": "comdex1plynttvhulryjv8cdej5r4ue4zxx82d6grjjph",
			"reward": {
				"denom": "ucmdx",
				"amount": "433"
			}
		},
		{
			"address": "comdex1plt8zgvnw6nrrh2876l94r5af2up77cjsahcqp",
			"reward": {
				"denom": "ucmdx",
				"amount": "3094"
			}
		},
		{
			"address": "comdex1pl5jtmhlkmvgjv9qhgeucw6u4qxtxp6k9pete7",
			"reward": {
				"denom": "ucmdx",
				"amount": "5979"
			}
		},
		{
			"address": "comdex1plhyktscnq22u96h8sg6tmltrv6vrl07jzst5k",
			"reward": {
				"denom": "ucmdx",
				"amount": "533"
			}
		},
		{
			"address": "comdex1plh2vje0yrmejvx3v0739cxu237586ng3sqtym",
			"reward": {
				"denom": "ucmdx",
				"amount": "6937"
			}
		},
		{
			"address": "comdex1plc8vwgxf2quey3mw5qe9fedmvdsjyzzp7klun",
			"reward": {
				"denom": "ucmdx",
				"amount": "1077"
			}
		},
		{
			"address": "comdex1ple90t3fahaaj4va7x86kx73780vzl7xr827wg",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1ple5cf4jftm8eafxa7wevpy9gae7dnnmwy7eh8",
			"reward": {
				"denom": "ucmdx",
				"amount": "69527"
			}
		},
		{
			"address": "comdex1pl60dp0rp6gjehe3fmtczxcyv4c9ufxla93zx6",
			"reward": {
				"denom": "ucmdx",
				"amount": "5751"
			}
		},
		{
			"address": "comdex1plmu2s7ecqv7pdfmcwlcdwmnwgpwgwjezuxw9f",
			"reward": {
				"denom": "ucmdx",
				"amount": "128967"
			}
		},
		{
			"address": "comdex1zqqj40n9vuve4hl84a0rhdzqkxzk42ttfj7m6r",
			"reward": {
				"denom": "ucmdx",
				"amount": "26006"
			}
		},
		{
			"address": "comdex1zqp246htc27mw56ea2ydvkg8ukx3q445glupyn",
			"reward": {
				"denom": "ucmdx",
				"amount": "7741"
			}
		},
		{
			"address": "comdex1zq9xc79svjvz8nua3weczg5csehnx8y429uuy5",
			"reward": {
				"denom": "ucmdx",
				"amount": "18926"
			}
		},
		{
			"address": "comdex1zq8fh2l8h77v8kvp2jt75adktnypcvvjaxajql",
			"reward": {
				"denom": "ucmdx",
				"amount": "2134"
			}
		},
		{
			"address": "comdex1zqfsqam6jd9wxnkefeq4xpaaujkvf72la2qt3n",
			"reward": {
				"denom": "ucmdx",
				"amount": "1768"
			}
		},
		{
			"address": "comdex1zq2jcz4rqauaayj83v0e6p7ulyclhtqz34wwls",
			"reward": {
				"denom": "ucmdx",
				"amount": "3527"
			}
		},
		{
			"address": "comdex1zqjkr6qgu7zmvj6gwy7gmzedtqsrsrune8l4es",
			"reward": {
				"denom": "ucmdx",
				"amount": "37706"
			}
		},
		{
			"address": "comdex1zq4q8qsrwztenme2a3ehhknxadr9n5qw3hmh2s",
			"reward": {
				"denom": "ucmdx",
				"amount": "157636"
			}
		},
		{
			"address": "comdex1zqm5tsatulpds9m4czxzq2rz7kdjav9cj6veh3",
			"reward": {
				"denom": "ucmdx",
				"amount": "1803"
			}
		},
		{
			"address": "comdex1zqlq0slcear0rggnvcgsxpv99mtlyck85fkvfh",
			"reward": {
				"denom": "ucmdx",
				"amount": "146000"
			}
		},
		{
			"address": "comdex1zpzr2hmgtmmeagpdrtzqg3tla8fhetdv68v36d",
			"reward": {
				"denom": "ucmdx",
				"amount": "397"
			}
		},
		{
			"address": "comdex1zpy0haerx08v0nvt9ekddlpufykjca6gv49l5q",
			"reward": {
				"denom": "ucmdx",
				"amount": "205"
			}
		},
		{
			"address": "comdex1zpgqfejgwqfrzc79xsl9pjy3w8fc6atqtr2e66",
			"reward": {
				"denom": "ucmdx",
				"amount": "34352"
			}
		},
		{
			"address": "comdex1zptq6h4dpjc9xqyxyyw8v8jxpqmfk5d4h2ayed",
			"reward": {
				"denom": "ucmdx",
				"amount": "2002"
			}
		},
		{
			"address": "comdex1zp0c3dzhks3clhwm2emte6e99pqr0s5n8ur3tn",
			"reward": {
				"denom": "ucmdx",
				"amount": "6809"
			}
		},
		{
			"address": "comdex1zps8rhlrfut6xdy0x0e9w8athp8nhh2r7yf3j8",
			"reward": {
				"denom": "ucmdx",
				"amount": "28243"
			}
		},
		{
			"address": "comdex1zpj0nsejtjjwamplemp83y4tpzj8fcln2pxg7n",
			"reward": {
				"denom": "ucmdx",
				"amount": "985"
			}
		},
		{
			"address": "comdex1zpjc5aesp4p9d7k4gyt4xava2mrnna5fam8dk9",
			"reward": {
				"denom": "ucmdx",
				"amount": "3315"
			}
		},
		{
			"address": "comdex1zpcapwrkds3g22qcpkmaqst6tzug05hpnv3g4c",
			"reward": {
				"denom": "ucmdx",
				"amount": "83"
			}
		},
		{
			"address": "comdex1zpe2h5ypc0z0yf4vzwx6z0wh3rgnqkp8vvjjvu",
			"reward": {
				"denom": "ucmdx",
				"amount": "600"
			}
		},
		{
			"address": "comdex1zpewufa0kyrfh37menttyraft3j3vzfysq9uty",
			"reward": {
				"denom": "ucmdx",
				"amount": "3562"
			}
		},
		{
			"address": "comdex1zz279r50n08rvhllrhwxk0fg5r8h6j8c4vr76a",
			"reward": {
				"denom": "ucmdx",
				"amount": "13309"
			}
		},
		{
			"address": "comdex1zz3echudaf6xrces6rksgsy6xqpcmjgc6xz636",
			"reward": {
				"denom": "ucmdx",
				"amount": "3523"
			}
		},
		{
			"address": "comdex1zrq9ag6mwjpkll3wc9garchy5sqr0856gtp50x",
			"reward": {
				"denom": "ucmdx",
				"amount": "1452"
			}
		},
		{
			"address": "comdex1zrzyw6l2spjpxsm23qkvfwsrntrhaqx5vgr8at",
			"reward": {
				"denom": "ucmdx",
				"amount": "20585"
			}
		},
		{
			"address": "comdex1zrzgz2epyh7jtsk6an84jp7x7xacng22s3d7ua",
			"reward": {
				"denom": "ucmdx",
				"amount": "1224"
			}
		},
		{
			"address": "comdex1zrrnwpjxwhsnxrj7lxn453dffj70q8yvlurect",
			"reward": {
				"denom": "ucmdx",
				"amount": "1066"
			}
		},
		{
			"address": "comdex1zr985yvxnw6uprkjxcm4z6shd9l450j2vem8p3",
			"reward": {
				"denom": "ucmdx",
				"amount": "199"
			}
		},
		{
			"address": "comdex1zrxzlzn0k4ppyqjuy8kylkykj42tclrvuk8h7q",
			"reward": {
				"denom": "ucmdx",
				"amount": "148"
			}
		},
		{
			"address": "comdex1zrxt0m74adk5lh37ys3pvfh5alve45zce9elp3",
			"reward": {
				"denom": "ucmdx",
				"amount": "7182"
			}
		},
		{
			"address": "comdex1zr8m533nh6lc637vxsnnkzw4rtawl3aq7mj5u6",
			"reward": {
				"denom": "ucmdx",
				"amount": "8946"
			}
		},
		{
			"address": "comdex1zrfjkxfvf7v7dk4fl09hzztl8mdp6r9v5u2awu",
			"reward": {
				"denom": "ucmdx",
				"amount": "6565"
			}
		},
		{
			"address": "comdex1zrdh2azcq8rmsgczk8qcyw5ea3sf790zqe6v7u",
			"reward": {
				"denom": "ucmdx",
				"amount": "7187"
			}
		},
		{
			"address": "comdex1zrwpux76j7r3ts2r95gtqrf2j22w9v75n2zdhg",
			"reward": {
				"denom": "ucmdx",
				"amount": "1566"
			}
		},
		{
			"address": "comdex1zrs53496tm53ucuw8c2g88l6hhtjfzrx3shp85",
			"reward": {
				"denom": "ucmdx",
				"amount": "442"
			}
		},
		{
			"address": "comdex1zrkaam6tmkxfzgcccndc7tw59q2xryv4xcggnl",
			"reward": {
				"denom": "ucmdx",
				"amount": "180"
			}
		},
		{
			"address": "comdex1zrequh7kyglhzgj0hf6fd6apfnea6d7sf2llyp",
			"reward": {
				"denom": "ucmdx",
				"amount": "337"
			}
		},
		{
			"address": "comdex1zrlvslywrvvaafxzthclfx0f0emd496w77264a",
			"reward": {
				"denom": "ucmdx",
				"amount": "690496"
			}
		},
		{
			"address": "comdex1zyqxtvzqfadl8l7nksvtfs6letxkzhl79xch06",
			"reward": {
				"denom": "ucmdx",
				"amount": "126397"
			}
		},
		{
			"address": "comdex1zyz2wpajw754zcdx4m9u5ar6hlxt8m7wlar2wm",
			"reward": {
				"denom": "ucmdx",
				"amount": "6572"
			}
		},
		{
			"address": "comdex1zyzlsrx3vwx84eru0w8qad34842y6rqdzvzyu5",
			"reward": {
				"denom": "ucmdx",
				"amount": "1745"
			}
		},
		{
			"address": "comdex1zyrqgg9ma4jwvy9cydq0mgyg6jjmk3cnlgaffc",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1zy9x0lxydx0lavkwuz97agc9c6x548cgtmx50c",
			"reward": {
				"denom": "ucmdx",
				"amount": "136"
			}
		},
		{
			"address": "comdex1zyx0kmlqg6payckm3s03qu5tph300unl5nruus",
			"reward": {
				"denom": "ucmdx",
				"amount": "2531"
			}
		},
		{
			"address": "comdex1zy8knphzhf9u8xeyjc6k3eps96q48hlmkn8nhx",
			"reward": {
				"denom": "ucmdx",
				"amount": "7445"
			}
		},
		{
			"address": "comdex1zyfhfmj9scarwcvnuw7sndlqmdl0mfelc5zp5s",
			"reward": {
				"denom": "ucmdx",
				"amount": "8912"
			}
		},
		{
			"address": "comdex1zyf6g2cn8k0jdnhqde6rf2vr0hnlqqsyapkpzk",
			"reward": {
				"denom": "ucmdx",
				"amount": "14485"
			}
		},
		{
			"address": "comdex1zyty63st4234xvypd5389tgqdvkzv6e6rv2swv",
			"reward": {
				"denom": "ucmdx",
				"amount": "711"
			}
		},
		{
			"address": "comdex1zyvx8em5h7pzr0azzvnj52vg09ny62wc3qexfh",
			"reward": {
				"denom": "ucmdx",
				"amount": "265"
			}
		},
		{
			"address": "comdex1zy0ypxgym856rd55qru40c9sxda4ly9eaeyd55",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1zy3nxgwtw9hr99tzfaxasgwklm9ev843dw2590",
			"reward": {
				"denom": "ucmdx",
				"amount": "144"
			}
		},
		{
			"address": "comdex1zy553za8nenzukmv65240323jhuvxzym3aftf4",
			"reward": {
				"denom": "ucmdx",
				"amount": "43463"
			}
		},
		{
			"address": "comdex1zy4ya2wwjfc0rtkq3uss4nl98zsny9yfndy36e",
			"reward": {
				"denom": "ucmdx",
				"amount": "1795"
			}
		},
		{
			"address": "comdex1zyku25927uuy203aphat7apmmy8qg38pn40yr5",
			"reward": {
				"denom": "ucmdx",
				"amount": "2848"
			}
		},
		{
			"address": "comdex1zyczu9sc6fm02y3n20wyqpskdj5pyyz5dm3xje",
			"reward": {
				"denom": "ucmdx",
				"amount": "2304"
			}
		},
		{
			"address": "comdex1zy7uuu6cd5fde3uunlh5l40jjf24ypd6sy9ej4",
			"reward": {
				"denom": "ucmdx",
				"amount": "287489"
			}
		},
		{
			"address": "comdex1zylt8ncjs29rkw3p3yz0gvzskpfpuqm5469n3j",
			"reward": {
				"denom": "ucmdx",
				"amount": "1761"
			}
		},
		{
			"address": "comdex1z9qv2zqnc4ksusxjxvswzx6j6rdqewgfcpj7yp",
			"reward": {
				"denom": "ucmdx",
				"amount": "19"
			}
		},
		{
			"address": "comdex1z9qkkvqnx4kytx3qpf8nm3r9ltjys6chgnkpyx",
			"reward": {
				"denom": "ucmdx",
				"amount": "6813"
			}
		},
		{
			"address": "comdex1z9p077a57nygyscrt54tvulpj8hukwafmnte0j",
			"reward": {
				"denom": "ucmdx",
				"amount": "76354"
			}
		},
		{
			"address": "comdex1z9pev3agkqpll4nsc33539rp5snm8ye9uyane4",
			"reward": {
				"denom": "ucmdx",
				"amount": "6917"
			}
		},
		{
			"address": "comdex1z9zyhgk59af064gaemclnf396gj9uep0cx6mez",
			"reward": {
				"denom": "ucmdx",
				"amount": "347753"
			}
		},
		{
			"address": "comdex1z9yj3xqwxefyzvkxvzzhduvsy8z5ryjh3447ah",
			"reward": {
				"denom": "ucmdx",
				"amount": "16"
			}
		},
		{
			"address": "comdex1z99arptjs2jlf0fnx6u4mr5nst07qrzz9ezy6h",
			"reward": {
				"denom": "ucmdx",
				"amount": "7183"
			}
		},
		{
			"address": "comdex1z9xpts34c996lk07j8neym6qc5smeyw99ggg28",
			"reward": {
				"denom": "ucmdx",
				"amount": "2473"
			}
		},
		{
			"address": "comdex1z98eg2ztdp2glyla62629nrlvczg8s7fgyen53",
			"reward": {
				"denom": "ucmdx",
				"amount": "1756"
			}
		},
		{
			"address": "comdex1z92n47uyeyxwtmfsw44nawy5plxfdxpnnxq65s",
			"reward": {
				"denom": "ucmdx",
				"amount": "3929"
			}
		},
		{
			"address": "comdex1z9w2xyjvz9punnkyw53wcpxerpd8mah5pfns8w",
			"reward": {
				"denom": "ucmdx",
				"amount": "4945"
			}
		},
		{
			"address": "comdex1z9wss2vhgj36pevpwn8fsynasms2urzxs768f3",
			"reward": {
				"denom": "ucmdx",
				"amount": "407"
			}
		},
		{
			"address": "comdex1z9nqkcmcld8mjmahvcwwx7thsk5n9l3gzngxye",
			"reward": {
				"denom": "ucmdx",
				"amount": "14198"
			}
		},
		{
			"address": "comdex1z94pggc3tpapmx8uetnlle7ysz0de6aykf3ycg",
			"reward": {
				"denom": "ucmdx",
				"amount": "21481"
			}
		},
		{
			"address": "comdex1z942qu7887lv5wt32gw03rjs7mmev6nx6n78h2",
			"reward": {
				"denom": "ucmdx",
				"amount": "24547"
			}
		},
		{
			"address": "comdex1z9h9jr6lev2mcwkzdnnqjh9q90zxaqrewkdeu0",
			"reward": {
				"denom": "ucmdx",
				"amount": "173"
			}
		},
		{
			"address": "comdex1z9h73y9fjpsc2zj3y7swushfgredeu99exvln4",
			"reward": {
				"denom": "ucmdx",
				"amount": "647"
			}
		},
		{
			"address": "comdex1z9m9cz2f0h0ntwym8hfvqfpqjr26a22jjw5ely",
			"reward": {
				"denom": "ucmdx",
				"amount": "89"
			}
		},
		{
			"address": "comdex1z9u3lv7hsxjnfss490w7d5t22v5s4tnpukspjj",
			"reward": {
				"denom": "ucmdx",
				"amount": "15621"
			}
		},
		{
			"address": "comdex1z97paefs6s9tm38fq5ekpcd9l4fkzeeqkmfwyk",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1z97acu3x6pr9zqnk80gvla5qrt7l28k295kkdc",
			"reward": {
				"denom": "ucmdx",
				"amount": "15100"
			}
		},
		{
			"address": "comdex1z9ltadgyhrqncq8e4m9f3klstxwuv035seqjtq",
			"reward": {
				"denom": "ucmdx",
				"amount": "12347"
			}
		},
		{
			"address": "comdex1z9l3jqpypdur7zv992sytghmxmv04fdnqlnc6m",
			"reward": {
				"denom": "ucmdx",
				"amount": "14253"
			}
		},
		{
			"address": "comdex1zxx8j75ngm8m38v9l5wreaavwnsuun7gkqrwe4",
			"reward": {
				"denom": "ucmdx",
				"amount": "1811"
			}
		},
		{
			"address": "comdex1zxxep8e6rc9tpxqunjt7xpsgk0glaqqmt6298u",
			"reward": {
				"denom": "ucmdx",
				"amount": "201"
			}
		},
		{
			"address": "comdex1zxfgrq5azczc4dtulpjcf6zaf3ma4y3r8uzdcz",
			"reward": {
				"denom": "ucmdx",
				"amount": "141"
			}
		},
		{
			"address": "comdex1zxde3a6g7lf6dgqjjquu63ef9egk7x7nu7y9yz",
			"reward": {
				"denom": "ucmdx",
				"amount": "5416"
			}
		},
		{
			"address": "comdex1zx0da3r7naa25z9g0wjcrwdx0tx4nzg6mwa3lf",
			"reward": {
				"denom": "ucmdx",
				"amount": "434"
			}
		},
		{
			"address": "comdex1zxszlqhqz9vypw4067zv60pygg0hcayewytu35",
			"reward": {
				"denom": "ucmdx",
				"amount": "6352"
			}
		},
		{
			"address": "comdex1zxnfvnzdv9e87gqmhz7cz940kmawjzpgllslgk",
			"reward": {
				"denom": "ucmdx",
				"amount": "4950"
			}
		},
		{
			"address": "comdex1zxnf5aa4af7ntsdf5jahwmdm99y3xvad5efdfh",
			"reward": {
				"denom": "ucmdx",
				"amount": "555"
			}
		},
		{
			"address": "comdex1zx5jmlyfpeckq75vl6pvfk9wvre6n0drrx9nuf",
			"reward": {
				"denom": "ucmdx",
				"amount": "52"
			}
		},
		{
			"address": "comdex1zxk9h3434pgwc69k6mvmawmkhc97zrucq32p3l",
			"reward": {
				"denom": "ucmdx",
				"amount": "9845"
			}
		},
		{
			"address": "comdex1zxkgxn4zmehh08gqf3t6pkj9gqp73acfs7mpz5",
			"reward": {
				"denom": "ucmdx",
				"amount": "1369"
			}
		},
		{
			"address": "comdex1zxck6wq3sdelgu9cluadfyt7c0swuazq23xqw5",
			"reward": {
				"denom": "ucmdx",
				"amount": "134"
			}
		},
		{
			"address": "comdex1zxe50kyr88dkmq9e3u0lyspjkdprarmvd0spnc",
			"reward": {
				"denom": "ucmdx",
				"amount": "1517"
			}
		},
		{
			"address": "comdex1zxm2lsvq8cspyhr98trk2xfe0hm7jjywasjyx4",
			"reward": {
				"denom": "ucmdx",
				"amount": "704"
			}
		},
		{
			"address": "comdex1z8rcv5qdpxtfa8v78fhh8lvwz7su0gquzvyhgk",
			"reward": {
				"denom": "ucmdx",
				"amount": "358"
			}
		},
		{
			"address": "comdex1z8yue4jj7xg2m59tqrtduynwvmcall4dcz58r4",
			"reward": {
				"denom": "ucmdx",
				"amount": "1435"
			}
		},
		{
			"address": "comdex1z8xtct40ljrk585xgktpjlf5x63rv2z4q6w8xl",
			"reward": {
				"denom": "ucmdx",
				"amount": "2720"
			}
		},
		{
			"address": "comdex1z8xemk98v03q7ql8zf5ks3syg432j7p246wj56",
			"reward": {
				"denom": "ucmdx",
				"amount": "123"
			}
		},
		{
			"address": "comdex1z8ttqqngajl4x5xw8fm9anqsmrf70d5x6yw9sw",
			"reward": {
				"denom": "ucmdx",
				"amount": "5139"
			}
		},
		{
			"address": "comdex1z8tdvvcnq4f0kujfhxf8dh4xnr573tp4rmsqwq",
			"reward": {
				"denom": "ucmdx",
				"amount": "173"
			}
		},
		{
			"address": "comdex1z8tkdrurtevzurndhuzheh5pxpndtnje9ma6tu",
			"reward": {
				"denom": "ucmdx",
				"amount": "251287"
			}
		},
		{
			"address": "comdex1z8vqvhu5qwjuyjyx62arlk2yaheafd4hchza00",
			"reward": {
				"denom": "ucmdx",
				"amount": "98"
			}
		},
		{
			"address": "comdex1z83d5zdh6jx6ks5cw6qygpwg32gf0fpjrwqq5s",
			"reward": {
				"denom": "ucmdx",
				"amount": "40"
			}
		},
		{
			"address": "comdex1z8jwfhfmthgaqkgt82zt3yjv0tm7vetxygpjyj",
			"reward": {
				"denom": "ucmdx",
				"amount": "23"
			}
		},
		{
			"address": "comdex1z8jm224m7pnd4hav905hx0djwhktegs4m0yfp5",
			"reward": {
				"denom": "ucmdx",
				"amount": "1695"
			}
		},
		{
			"address": "comdex1z8nkfutrqdzrqzdqu2fmsv9yfvh8uuf7xe2rg2",
			"reward": {
				"denom": "ucmdx",
				"amount": "100277"
			}
		},
		{
			"address": "comdex1z85v6kzkp9d2g4mkz7grq69lvpjrc6u9ggs7zq",
			"reward": {
				"denom": "ucmdx",
				"amount": "1459"
			}
		},
		{
			"address": "comdex1z84tvu9lxfkkgw3yjuvqwv2dq84u25963nnkss",
			"reward": {
				"denom": "ucmdx",
				"amount": "943"
			}
		},
		{
			"address": "comdex1z84shuuq0pu3lx8qsqxee53vthhr7fczpluyum",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1z8cj0z7pctrhyly5uxlra4kwhygqn2j9ydwjnj",
			"reward": {
				"denom": "ucmdx",
				"amount": "7085"
			}
		},
		{
			"address": "comdex1z8eecvkxtca7qng72k2g8ayzmkk99njl93gxt0",
			"reward": {
				"denom": "ucmdx",
				"amount": "28"
			}
		},
		{
			"address": "comdex1z8amf5zauum3vwrf3963yu645lzt2x4m3gqj0c",
			"reward": {
				"denom": "ucmdx",
				"amount": "17555"
			}
		},
		{
			"address": "comdex1z8a702hldm8xae5v62jhcz2du4pepyutnxc5wr",
			"reward": {
				"denom": "ucmdx",
				"amount": "3032"
			}
		},
		{
			"address": "comdex1z8lhhf37w0u02sgpzlg94zptwddkh2zrm0cp68",
			"reward": {
				"denom": "ucmdx",
				"amount": "21078"
			}
		},
		{
			"address": "comdex1z8l7umujphfzrnvhvnpsaa7gajffts8jttes0d",
			"reward": {
				"denom": "ucmdx",
				"amount": "6875"
			}
		},
		{
			"address": "comdex1zgqqgqhp2c75j6w2k06je2k32f420l4krs0rt2",
			"reward": {
				"denom": "ucmdx",
				"amount": "659"
			}
		},
		{
			"address": "comdex1zgp0gvr67fn5zrjup42sxm40ppu05sju3qr6gt",
			"reward": {
				"denom": "ucmdx",
				"amount": "1787"
			}
		},
		{
			"address": "comdex1zgyzssl32r9cwcql9wpn7t6sdkntuh7ax37mzt",
			"reward": {
				"denom": "ucmdx",
				"amount": "178"
			}
		},
		{
			"address": "comdex1zg94v0zhl4trlg2hqcjx2w6cgy6x77eglgmqj6",
			"reward": {
				"denom": "ucmdx",
				"amount": "183"
			}
		},
		{
			"address": "comdex1zgx636tcaxndz43wjf4qtraa0rkkwj5ydlj5fr",
			"reward": {
				"denom": "ucmdx",
				"amount": "53"
			}
		},
		{
			"address": "comdex1zgtpc7kur4uqdscax88mlnm3rvfl5f282s3nsa",
			"reward": {
				"denom": "ucmdx",
				"amount": "1831"
			}
		},
		{
			"address": "comdex1zgtvrhzf0stpxcsj6qt0kmq70fl3ewet4f9w54",
			"reward": {
				"denom": "ucmdx",
				"amount": "29555"
			}
		},
		{
			"address": "comdex1zgd6069zp596w50xym0rw3y55xf5srjnygjl3a",
			"reward": {
				"denom": "ucmdx",
				"amount": "49907"
			}
		},
		{
			"address": "comdex1zg0mar95nrl0kalgfxwergfqqvp08lrt62swnj",
			"reward": {
				"denom": "ucmdx",
				"amount": "7456"
			}
		},
		{
			"address": "comdex1zgsm8qwpxs48l7678fzfef0t0m38qu7tj03p2w",
			"reward": {
				"denom": "ucmdx",
				"amount": "881"
			}
		},
		{
			"address": "comdex1zg3qqff5qta9f6gqkkxhhmn5ceszzfnstzm5v5",
			"reward": {
				"denom": "ucmdx",
				"amount": "536"
			}
		},
		{
			"address": "comdex1zg42hhc2tkawgemndchyjwq6mtxzxqp96pvq0q",
			"reward": {
				"denom": "ucmdx",
				"amount": "63"
			}
		},
		{
			"address": "comdex1zg453l326xgrk5nm2prru4teqkym8zc9wd93q8",
			"reward": {
				"denom": "ucmdx",
				"amount": "4830"
			}
		},
		{
			"address": "comdex1zg4cwm5j3keugcmprr9j67433kn9m46tw2svfs",
			"reward": {
				"denom": "ucmdx",
				"amount": "611"
			}
		},
		{
			"address": "comdex1zgksg42ylta26axgpvz4zw8l8jydy68jm8eu7h",
			"reward": {
				"denom": "ucmdx",
				"amount": "14573"
			}
		},
		{
			"address": "comdex1zgaj2w950llsludydz5xmcp609uk02ysk9l893",
			"reward": {
				"denom": "ucmdx",
				"amount": "57506"
			}
		},
		{
			"address": "comdex1zgals4zkgun8fd68yhtncl0nxn9x3ss08j6d9q",
			"reward": {
				"denom": "ucmdx",
				"amount": "626946"
			}
		},
		{
			"address": "comdex1zglcvru504scyum4vscgz3chw906qxm8q020z4",
			"reward": {
				"denom": "ucmdx",
				"amount": "5097"
			}
		},
		{
			"address": "comdex1zfr9vn8r9yxyp755e33gn5gt7sa3mwse7xwrkr",
			"reward": {
				"denom": "ucmdx",
				"amount": "17327"
			}
		},
		{
			"address": "comdex1zfystmr33pk843ph4zcecyedfj8muvzsnuxdyz",
			"reward": {
				"denom": "ucmdx",
				"amount": "349"
			}
		},
		{
			"address": "comdex1zf9zpzfflhy9j7yp8987t92vyw5743nrgmus3p",
			"reward": {
				"denom": "ucmdx",
				"amount": "352"
			}
		},
		{
			"address": "comdex1zf90768qdhr08xnygm5d2w2zdpfcq8ane70jhp",
			"reward": {
				"denom": "ucmdx",
				"amount": "2907"
			}
		},
		{
			"address": "comdex1zf80hu9mczv4vzfjsxeql3pzmr5h9g2pnwejtp",
			"reward": {
				"denom": "ucmdx",
				"amount": "2041"
			}
		},
		{
			"address": "comdex1zfgdvmmzs3r8pwzyy806vf2jeynem9fln57mmz",
			"reward": {
				"denom": "ucmdx",
				"amount": "582"
			}
		},
		{
			"address": "comdex1zffyupsg4ckan05w5ecdylnzx09clt5535qy4c",
			"reward": {
				"denom": "ucmdx",
				"amount": "285"
			}
		},
		{
			"address": "comdex1zf2rqyuxk05pqmpz5z3nyz42hhewnrhn8w42f0",
			"reward": {
				"denom": "ucmdx",
				"amount": "193891"
			}
		},
		{
			"address": "comdex1zftpa0wva7ddp06r4hfpn27gv9mam6q2z0rnge",
			"reward": {
				"denom": "ucmdx",
				"amount": "13451"
			}
		},
		{
			"address": "comdex1zftvvwskk990fq6lggaz3r9m5j0sf3dptnr9j7",
			"reward": {
				"denom": "ucmdx",
				"amount": "19"
			}
		},
		{
			"address": "comdex1zfv5u4fm2sptsvq48n9h02ctk55xcl57tkw9mw",
			"reward": {
				"denom": "ucmdx",
				"amount": "7046"
			}
		},
		{
			"address": "comdex1zfsyymf2852d8y07lrzfm2cc7dm6uqahaazycl",
			"reward": {
				"denom": "ucmdx",
				"amount": "62083"
			}
		},
		{
			"address": "comdex1zfngg3euaseg5csr7fml4stz4rn8wwns7mkc2z",
			"reward": {
				"denom": "ucmdx",
				"amount": "408"
			}
		},
		{
			"address": "comdex1zfnkez04pq0djns3yh59hqkkkk0ndmqcl9p45w",
			"reward": {
				"denom": "ucmdx",
				"amount": "35318"
			}
		},
		{
			"address": "comdex1zf5z9hr57h20sdjjh2370qhqqp3sa7n4wwlgww",
			"reward": {
				"denom": "ucmdx",
				"amount": "14"
			}
		},
		{
			"address": "comdex1zfkclxhd4443vek7g86nwcq0vce2jf3vf5nr9j",
			"reward": {
				"denom": "ucmdx",
				"amount": "74790"
			}
		},
		{
			"address": "comdex1zfcj782j2fydqjfagdnl0v7p95vm52wgafja2w",
			"reward": {
				"denom": "ucmdx",
				"amount": "101"
			}
		},
		{
			"address": "comdex1z2p8yxg9282pvyehht725l48w8qkpc8ujadghs",
			"reward": {
				"denom": "ucmdx",
				"amount": "14110"
			}
		},
		{
			"address": "comdex1z2r8c2hfeajzn54q02gujvhpgk3ag66tng58jk",
			"reward": {
				"denom": "ucmdx",
				"amount": "3196"
			}
		},
		{
			"address": "comdex1z2y028kr6ynapczvxkvkcrlmu7w9y7t06whsj6",
			"reward": {
				"denom": "ucmdx",
				"amount": "1511"
			}
		},
		{
			"address": "comdex1z29zepzaq0th2a5c43hg7hkgp8lgx3ddlnrnpc",
			"reward": {
				"denom": "ucmdx",
				"amount": "173"
			}
		},
		{
			"address": "comdex1z28ytt72ef5dgtnf2j36lk7fn8m6q55tmqhmqk",
			"reward": {
				"denom": "ucmdx",
				"amount": "32944"
			}
		},
		{
			"address": "comdex1z2w333ttrqhkvpu4x9v059l7q6k9avcdmhm96v",
			"reward": {
				"denom": "ucmdx",
				"amount": "50748"
			}
		},
		{
			"address": "comdex1z2saeflyzxl6ss2k4f365uavf2uvyka97jf966",
			"reward": {
				"denom": "ucmdx",
				"amount": "7299"
			}
		},
		{
			"address": "comdex1z2j02xqwgsehdk3g2zlpryc0vraeku3e3ehvfe",
			"reward": {
				"denom": "ucmdx",
				"amount": "519"
			}
		},
		{
			"address": "comdex1z2j5qnujajl2l7fepega92kfqcfrertsgn72ps",
			"reward": {
				"denom": "ucmdx",
				"amount": "4104"
			}
		},
		{
			"address": "comdex1z2nupzll49pg5737klwm025dcz8e5sp3h4wdz6",
			"reward": {
				"denom": "ucmdx",
				"amount": "1440"
			}
		},
		{
			"address": "comdex1z25uywjrzeuvw0rmnhxney4fdcql8d4w8ff5j3",
			"reward": {
				"denom": "ucmdx",
				"amount": "72436"
			}
		},
		{
			"address": "comdex1z25uef9lxk54u77hy3ssv2erfhzh6jz82y9vxc",
			"reward": {
				"denom": "ucmdx",
				"amount": "12958"
			}
		},
		{
			"address": "comdex1z24npynn7cunwshk2fq0lu8ckslrja8h903rx3",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1z2h8qstwqt2jpljuw6ma0yz6qhce6fe4grfnjd",
			"reward": {
				"denom": "ucmdx",
				"amount": "1232"
			}
		},
		{
			"address": "comdex1z2cgkhumpypjmd9md495we0ne66z63gdc0r2ut",
			"reward": {
				"denom": "ucmdx",
				"amount": "1400"
			}
		},
		{
			"address": "comdex1z2cmdk7atwfefl4a3had7a2tsamxrwgucmhutx",
			"reward": {
				"denom": "ucmdx",
				"amount": "49831"
			}
		},
		{
			"address": "comdex1z26qsal3w2txaw58egw8nmtccwcd0tzufpqgv8",
			"reward": {
				"denom": "ucmdx",
				"amount": "4057"
			}
		},
		{
			"address": "comdex1z2uecn5khjexks0reldjjfm5a9ssj6zvze4lel",
			"reward": {
				"denom": "ucmdx",
				"amount": "12531"
			}
		},
		{
			"address": "comdex1z2ul6vm9lvj0lz4zrtte78kegt2v2hgx34qrfj",
			"reward": {
				"denom": "ucmdx",
				"amount": "10616"
			}
		},
		{
			"address": "comdex1z2lwvgzw4ut9z9azegzz7kh7fc4wwada9yc0np",
			"reward": {
				"denom": "ucmdx",
				"amount": "127"
			}
		},
		{
			"address": "comdex1ztpm0a57zlax2vf5e2u0elhdn6hupj7cqcel4x",
			"reward": {
				"denom": "ucmdx",
				"amount": "1772"
			}
		},
		{
			"address": "comdex1ztz8t88xgqvueap24ufpfun5mxc3jhmchwqwd6",
			"reward": {
				"denom": "ucmdx",
				"amount": "156265"
			}
		},
		{
			"address": "comdex1zt9dvaje2t6vzhd2d2qwm4d3sswq60clrrzad0",
			"reward": {
				"denom": "ucmdx",
				"amount": "5604"
			}
		},
		{
			"address": "comdex1ztx53a02e292sf7gedj56f29ynvmcas0lk4eq2",
			"reward": {
				"denom": "ucmdx",
				"amount": "5181"
			}
		},
		{
			"address": "comdex1ztfzxph6fyfsrnkjqhj9efw3lmdsn2l7rkacq0",
			"reward": {
				"denom": "ucmdx",
				"amount": "298"
			}
		},
		{
			"address": "comdex1ztfk3dglwhskqyktms2ts4n8dgq62s46m8cjuf",
			"reward": {
				"denom": "ucmdx",
				"amount": "351815"
			}
		},
		{
			"address": "comdex1zttl5ud0klt90zlcp6kxwczncz7tvhl94dyz7g",
			"reward": {
				"denom": "ucmdx",
				"amount": "12442"
			}
		},
		{
			"address": "comdex1ztd2g5ckqglzhrl3tf72yj43zjgs46yu9cufm0",
			"reward": {
				"denom": "ucmdx",
				"amount": "5621"
			}
		},
		{
			"address": "comdex1ztjgcrp8zzh4sy7uuxsjdj27pa8v696dlf25rq",
			"reward": {
				"denom": "ucmdx",
				"amount": "143"
			}
		},
		{
			"address": "comdex1ztj2gc3v2xsjsap4gt8xx83x3ff89ar4f84ntx",
			"reward": {
				"denom": "ucmdx",
				"amount": "2374"
			}
		},
		{
			"address": "comdex1ztk8k8eaz43l7m23zx6cr034sf76dnvg57zc0m",
			"reward": {
				"denom": "ucmdx",
				"amount": "44"
			}
		},
		{
			"address": "comdex1zthyparxg29vq28csu683curze3wajh3pkvn6e",
			"reward": {
				"denom": "ucmdx",
				"amount": "7204"
			}
		},
		{
			"address": "comdex1zt6dnk4fpaq94lfpx00hx26j8ue3dusemwslhr",
			"reward": {
				"denom": "ucmdx",
				"amount": "2830"
			}
		},
		{
			"address": "comdex1ztaj960x4vh44e5r4cpc0uku8d0due5dtwk4y3",
			"reward": {
				"denom": "ucmdx",
				"amount": "409"
			}
		},
		{
			"address": "comdex1zt7t6l84sc4c00p6479uuqyzzhkme7ywv8ldpv",
			"reward": {
				"denom": "ucmdx",
				"amount": "3986"
			}
		},
		{
			"address": "comdex1ztl2nu2hv6ulpgteqh9mfr00uv3qjfmjju4qwv",
			"reward": {
				"denom": "ucmdx",
				"amount": "8296"
			}
		},
		{
			"address": "comdex1ztlt6upprz0ss5jwdn8wh8k82c3jw97pnyun4z",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1zvy8thtvr6vztynkpm3yuw7tvd6y3xy8sdsa4t",
			"reward": {
				"denom": "ucmdx",
				"amount": "141"
			}
		},
		{
			"address": "comdex1zvyardw7kde5w8flkfp3ghmg7u4834gcwx5z4c",
			"reward": {
				"denom": "ucmdx",
				"amount": "34222"
			}
		},
		{
			"address": "comdex1zv8qvgwxs5r4cfqqhtw2wjv53u4vccfevr6p5j",
			"reward": {
				"denom": "ucmdx",
				"amount": "27645"
			}
		},
		{
			"address": "comdex1zvgz5ryz3rn0slt269cpqvhx8k63jpzf34e7vh",
			"reward": {
				"denom": "ucmdx",
				"amount": "178"
			}
		},
		{
			"address": "comdex1zvgr8mwtdtzel0p9tvumcm6dkr3n8pp62n44uf",
			"reward": {
				"denom": "ucmdx",
				"amount": "1424"
			}
		},
		{
			"address": "comdex1zvgyalpnjs0838tzmnk4panj200yjy6aqmwv20",
			"reward": {
				"denom": "ucmdx",
				"amount": "12368"
			}
		},
		{
			"address": "comdex1zvfwnkglqu3vqac4x0qjzra8aalsjhkammu7u9",
			"reward": {
				"denom": "ucmdx",
				"amount": "6600"
			}
		},
		{
			"address": "comdex1zvf4thvqgsdf770dwnpuptcaztv630lm6snhxk",
			"reward": {
				"denom": "ucmdx",
				"amount": "7114"
			}
		},
		{
			"address": "comdex1zvd3c0x7kaaauavguevnln9nwzsu7sewfg9q2c",
			"reward": {
				"denom": "ucmdx",
				"amount": "1994"
			}
		},
		{
			"address": "comdex1zvdnv9ytarls8arep9ny9cfcrcffpd68cz7y88",
			"reward": {
				"denom": "ucmdx",
				"amount": "361"
			}
		},
		{
			"address": "comdex1zv0t2ca3lmgt465dsnpkdtlfl6qr6rn27x42nm",
			"reward": {
				"denom": "ucmdx",
				"amount": "243"
			}
		},
		{
			"address": "comdex1zvj80m9ckmtp5cglwt6wn0dcs0y9ceddt63wpr",
			"reward": {
				"denom": "ucmdx",
				"amount": "64320"
			}
		},
		{
			"address": "comdex1zv4pxj3v3zyddxffcue80zqwgpj9y024rkxtgm",
			"reward": {
				"denom": "ucmdx",
				"amount": "1765"
			}
		},
		{
			"address": "comdex1zvmg5xynymhrhajqa0aerzk3hlw2qa9psre932",
			"reward": {
				"denom": "ucmdx",
				"amount": "12177"
			}
		},
		{
			"address": "comdex1zvm56xzrykqwpxejfvw2fcvsftdryakv3gf3vf",
			"reward": {
				"denom": "ucmdx",
				"amount": "18995"
			}
		},
		{
			"address": "comdex1zdqj4mtc7wldhvu7qcak88rxqyfp6vu38l7f70",
			"reward": {
				"denom": "ucmdx",
				"amount": "7199"
			}
		},
		{
			"address": "comdex1zdxzpcjafydkw30j8lwsjdd2dxnhdaa40kn3z5",
			"reward": {
				"denom": "ucmdx",
				"amount": "13903"
			}
		},
		{
			"address": "comdex1zdtmgwpzwtwepqywvdhzjg6h0tatddzempalhh",
			"reward": {
				"denom": "ucmdx",
				"amount": "54998"
			}
		},
		{
			"address": "comdex1zddaqz5mrewea0u0xt64lmwckakss2mytadmsg",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1zdwtk2gz573nvh4caga9ea23gue2mg2xt2lvmc",
			"reward": {
				"denom": "ucmdx",
				"amount": "3014"
			}
		},
		{
			"address": "comdex1zd0vemquu329xyl9esmqfgmv23pv675r2ptt8v",
			"reward": {
				"denom": "ucmdx",
				"amount": "1457"
			}
		},
		{
			"address": "comdex1zds38355lhrnumlqertcxwclqdy0y3u0dmg3w9",
			"reward": {
				"denom": "ucmdx",
				"amount": "331"
			}
		},
		{
			"address": "comdex1zdjt49vtgdylssjsxdyx9q2kujjzej5kur5gcd",
			"reward": {
				"denom": "ucmdx",
				"amount": "859"
			}
		},
		{
			"address": "comdex1zdnfqlakqr8sgg4z9ch72y3pd2pyp6usagsdwc",
			"reward": {
				"denom": "ucmdx",
				"amount": "1190"
			}
		},
		{
			"address": "comdex1zdnh7vj6w0t54sc7xacdcj6ffu4t0ssqwax8gs",
			"reward": {
				"denom": "ucmdx",
				"amount": "509"
			}
		},
		{
			"address": "comdex1zd5yfrakr5f23hzzlw4x99kwyfy47s0t3qmjkx",
			"reward": {
				"denom": "ucmdx",
				"amount": "1758"
			}
		},
		{
			"address": "comdex1zd5xv7fq0f64u9cv4j0xte7mrepex0mcee5ssa",
			"reward": {
				"denom": "ucmdx",
				"amount": "9632"
			}
		},
		{
			"address": "comdex1zdht8zshum4xk4v6x6rprk2k67ylms7e8qxvup",
			"reward": {
				"denom": "ucmdx",
				"amount": "1500"
			}
		},
		{
			"address": "comdex1zdhc5g9n9y4s0fzx98pxk6vrgw9tg2v4qqq80j",
			"reward": {
				"denom": "ucmdx",
				"amount": "102"
			}
		},
		{
			"address": "comdex1zdek2gepft4mm5m5m56z4sdg9e67ddkcqw4dyg",
			"reward": {
				"denom": "ucmdx",
				"amount": "1549"
			}
		},
		{
			"address": "comdex1zduwtppgh690zpvvzplr3y0l8e8yh55e6qzfn0",
			"reward": {
				"denom": "ucmdx",
				"amount": "1693"
			}
		},
		{
			"address": "comdex1zdap4hxc3accq5lfj04ey4wed3up02354g63gy",
			"reward": {
				"denom": "ucmdx",
				"amount": "4571"
			}
		},
		{
			"address": "comdex1zda94m6kzhk9cecyl8j8m6ud4f7wqc8c3zwzhy",
			"reward": {
				"denom": "ucmdx",
				"amount": "1938"
			}
		},
		{
			"address": "comdex1zdaavhnhx8sfdqgtw9ygfxfn304d88lhaq7mx7",
			"reward": {
				"denom": "ucmdx",
				"amount": "17409"
			}
		},
		{
			"address": "comdex1zd7zdtnhefkyc8pp9t4ekmene6rpgjkpcn5kpd",
			"reward": {
				"denom": "ucmdx",
				"amount": "35428"
			}
		},
		{
			"address": "comdex1zwqu43hnhm34pp9ldh5zz2ryh2g0lpe2vg7uzc",
			"reward": {
				"denom": "ucmdx",
				"amount": "78"
			}
		},
		{
			"address": "comdex1zwyys6qq3llh5esw7twdpwhc6tmm53cdqcea6l",
			"reward": {
				"denom": "ucmdx",
				"amount": "17288"
			}
		},
		{
			"address": "comdex1zwygeutz5wfug2g5ffn23eq0vtdwahc3553x7e",
			"reward": {
				"denom": "ucmdx",
				"amount": "524"
			}
		},
		{
			"address": "comdex1zw80u6ehx09q3l2ce72q2z6qtmrajdkvms9tj8",
			"reward": {
				"denom": "ucmdx",
				"amount": "71152"
			}
		},
		{
			"address": "comdex1zwftw570zheq5scnzmm830pjmsl3d32erua28m",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1zwsvkfyc2l0nfruemm5t7mv9lc0rlkwcdgfcrn",
			"reward": {
				"denom": "ucmdx",
				"amount": "101914"
			}
		},
		{
			"address": "comdex1zws09vrt3sk26x235fadqhlc0szn844nqpgx72",
			"reward": {
				"denom": "ucmdx",
				"amount": "4215"
			}
		},
		{
			"address": "comdex1zwnzthzckesnxh22q2g6w8l32h6cm4c0epqt40",
			"reward": {
				"denom": "ucmdx",
				"amount": "16382"
			}
		},
		{
			"address": "comdex1zw49fkp3jwa27g4nnsuz2qaqjzxe4lq34rmjth",
			"reward": {
				"denom": "ucmdx",
				"amount": "241"
			}
		},
		{
			"address": "comdex1zwk3vye3lcxe76vpdz8qyd2gf6n2vcx6ska3gu",
			"reward": {
				"denom": "ucmdx",
				"amount": "466"
			}
		},
		{
			"address": "comdex1zwk4enxs404xdwxutzwsqy3r6w9c272x50uf6f",
			"reward": {
				"denom": "ucmdx",
				"amount": "9305"
			}
		},
		{
			"address": "comdex1zwhajzfa43xn58fpsa039qeyqmnyzjcnl4uns9",
			"reward": {
				"denom": "ucmdx",
				"amount": "200"
			}
		},
		{
			"address": "comdex1zwmnvveh9469q7j7xrpdp5z2345t755vztp64x",
			"reward": {
				"denom": "ucmdx",
				"amount": "30390"
			}
		},
		{
			"address": "comdex1zwunv3ag8pyd69sl6fffqxhxpp3mh4p5jeg2xn",
			"reward": {
				"denom": "ucmdx",
				"amount": "1790"
			}
		},
		{
			"address": "comdex1zwu7tmfuc040g0sreutsf64huv3s6xr7ms2g20",
			"reward": {
				"denom": "ucmdx",
				"amount": "10749"
			}
		},
		{
			"address": "comdex1zwanpg0p3lyycqemrm4ju6r0npk824xdl7t0mk",
			"reward": {
				"denom": "ucmdx",
				"amount": "937"
			}
		},
		{
			"address": "comdex1zw7za06w00ha65jr9dcrrccr5m0y6r0aesc7st",
			"reward": {
				"denom": "ucmdx",
				"amount": "251"
			}
		},
		{
			"address": "comdex1zwlutvq3veqzjf8vewhe0k8t7jj7d3hdnfun5n",
			"reward": {
				"denom": "ucmdx",
				"amount": "1190"
			}
		},
		{
			"address": "comdex1z0qzwl72aj6pct2ct35lp0s0djg4p5snqqhxg0",
			"reward": {
				"denom": "ucmdx",
				"amount": "14"
			}
		},
		{
			"address": "comdex1z0zchsuqtxpetp54j7pl557fkf93ya3cmw00ny",
			"reward": {
				"denom": "ucmdx",
				"amount": "29476"
			}
		},
		{
			"address": "comdex1z0gsymslenehmhqj9tw6m7zs6p33ldurz96s74",
			"reward": {
				"denom": "ucmdx",
				"amount": "135"
			}
		},
		{
			"address": "comdex1z0ga39egyf0q0crake03a5te9yc54puph0qwkj",
			"reward": {
				"denom": "ucmdx",
				"amount": "9990"
			}
		},
		{
			"address": "comdex1z0293ptk7as6z47qds95lavwrerv7qajx0lx8m",
			"reward": {
				"denom": "ucmdx",
				"amount": "180"
			}
		},
		{
			"address": "comdex1z02dwh3p5u8sjhdeq33fg5clnp3vvg4qf2efuw",
			"reward": {
				"denom": "ucmdx",
				"amount": "960"
			}
		},
		{
			"address": "comdex1z0se2v8nw270sgueh02wwz2qrr9wx8x5300h9j",
			"reward": {
				"denom": "ucmdx",
				"amount": "107"
			}
		},
		{
			"address": "comdex1z0jlmnqwa6rszu8ksnh9trnajqpdwuw4fd0fx3",
			"reward": {
				"denom": "ucmdx",
				"amount": "23886"
			}
		},
		{
			"address": "comdex1z0nmvynk5fadw6lamdpewr9dnrjap49hau0rl6",
			"reward": {
				"denom": "ucmdx",
				"amount": "10221"
			}
		},
		{
			"address": "comdex1z0khlnpg8j5jruvzljsm4ynm4vmzr7vv07s6n2",
			"reward": {
				"denom": "ucmdx",
				"amount": "303"
			}
		},
		{
			"address": "comdex1z06q4htuzheuntqp85ls39vt96yhtsj4parfxg",
			"reward": {
				"denom": "ucmdx",
				"amount": "416"
			}
		},
		{
			"address": "comdex1z07ucv5gz23d2x4mk66nunjehk94mme4u39wnd",
			"reward": {
				"denom": "ucmdx",
				"amount": "7090"
			}
		},
		{
			"address": "comdex1z0lx2fhw3xrajgg0dae98flfafe5543nc0nppm",
			"reward": {
				"denom": "ucmdx",
				"amount": "2439"
			}
		},
		{
			"address": "comdex1zsp33an6uqwfyhfhzruzgmc2859w5eejs55025",
			"reward": {
				"denom": "ucmdx",
				"amount": "143"
			}
		},
		{
			"address": "comdex1zszznjnt57hy8y7y3cq7d6ldspxmjmtlhcmpvp",
			"reward": {
				"denom": "ucmdx",
				"amount": "620"
			}
		},
		{
			"address": "comdex1zsz6wh57axhgdxaz2em8f0q5hjv3v8863sme4z",
			"reward": {
				"denom": "ucmdx",
				"amount": "75"
			}
		},
		{
			"address": "comdex1zsygz8es37kyfnhsg3xp59d43re4tjtqntlcfh",
			"reward": {
				"denom": "ucmdx",
				"amount": "353"
			}
		},
		{
			"address": "comdex1zsyhd2z2m40xm098tqhvjfcgacml59erqpkrh8",
			"reward": {
				"denom": "ucmdx",
				"amount": "7343"
			}
		},
		{
			"address": "comdex1zs9dl6yv3dalwl88sx6t3p3mprrhq8cdsq4rs7",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1zsxvswh0sgpa8gkq6nrjpmj38gr2yaqt07w3nh",
			"reward": {
				"denom": "ucmdx",
				"amount": "805"
			}
		},
		{
			"address": "comdex1zsxcl4k76cwffjfrtqzmk38rdsx0ehger3upzg",
			"reward": {
				"denom": "ucmdx",
				"amount": "24321"
			}
		},
		{
			"address": "comdex1zs80vefq4c99mwuaxps9m5xl8cng2jvzc92g9p",
			"reward": {
				"denom": "ucmdx",
				"amount": "6954"
			}
		},
		{
			"address": "comdex1zs8k656xzml49uzlx8nex2l8ur2lfs97amwsd0",
			"reward": {
				"denom": "ucmdx",
				"amount": "8271"
			}
		},
		{
			"address": "comdex1zs2nylv6w5a9tmyhwhfjjnyrvtswmlkf6lxkay",
			"reward": {
				"denom": "ucmdx",
				"amount": "6111"
			}
		},
		{
			"address": "comdex1zsv6cvpd9y86hhxzatpgzu6ymxdeuc5nmtlgvn",
			"reward": {
				"denom": "ucmdx",
				"amount": "353"
			}
		},
		{
			"address": "comdex1zsdc50wkpg7u7l6k424eqhf76l8hrk8cesvt3p",
			"reward": {
				"denom": "ucmdx",
				"amount": "348"
			}
		},
		{
			"address": "comdex1zs0ltml24ys8kengpl8hvntfyghld2zcpqxflq",
			"reward": {
				"denom": "ucmdx",
				"amount": "886"
			}
		},
		{
			"address": "comdex1zs59qzf9g87jtpg4p0pmp5par72en9v3n93ctp",
			"reward": {
				"denom": "ucmdx",
				"amount": "28764"
			}
		},
		{
			"address": "comdex1zs400s4srmx7gwtwxsvevzx88tpdu3u9swlqpv",
			"reward": {
				"denom": "ucmdx",
				"amount": "48698"
			}
		},
		{
			"address": "comdex1zskq7uc72mazdxe6tdnm4ezzj46t35sre073af",
			"reward": {
				"denom": "ucmdx",
				"amount": "415"
			}
		},
		{
			"address": "comdex1zshz7gymvmljrtwplu83dzpm3e8cxrg94nlk6w",
			"reward": {
				"denom": "ucmdx",
				"amount": "6784"
			}
		},
		{
			"address": "comdex1zsm9gh9vephm6s73r2l4drddx7grsx4t972pdy",
			"reward": {
				"denom": "ucmdx",
				"amount": "5051"
			}
		},
		{
			"address": "comdex1zsuh63f4axsxzsmceen9ccd4309fee6akdaveq",
			"reward": {
				"denom": "ucmdx",
				"amount": "14322"
			}
		},
		{
			"address": "comdex1z3pzzw84d6xn00pw9dy3yapqypfde7vgtcmdh0",
			"reward": {
				"denom": "ucmdx",
				"amount": "6035"
			}
		},
		{
			"address": "comdex1z3pm4fkxyqzgau5qr2q4en425qgketkttwvde3",
			"reward": {
				"denom": "ucmdx",
				"amount": "44883"
			}
		},
		{
			"address": "comdex1z3z365lydry32xghj9m8gt62mvddl55gnj0kwa",
			"reward": {
				"denom": "ucmdx",
				"amount": "1761"
			}
		},
		{
			"address": "comdex1z3znw7lall94lym3vazl6xxd36a2wtqgrncemz",
			"reward": {
				"denom": "ucmdx",
				"amount": "125124"
			}
		},
		{
			"address": "comdex1z3z42cka5j8ph93xd60guw6pzkchaqwqw38wd5",
			"reward": {
				"denom": "ucmdx",
				"amount": "20001"
			}
		},
		{
			"address": "comdex1z3fa02a89qthsdmvzn53lsw5u8pa8uy07lsec5",
			"reward": {
				"denom": "ucmdx",
				"amount": "612"
			}
		},
		{
			"address": "comdex1z32qsutjlad4yufkanxl4mlvhe9wqpj4sfyelr",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1z307neur6kmcmlm7ekfssjgsmqx2raujhedfqa",
			"reward": {
				"denom": "ucmdx",
				"amount": "764"
			}
		},
		{
			"address": "comdex1z3nwatrw9gdtu90nmt3ap34nzz04zrdu552w0e",
			"reward": {
				"denom": "ucmdx",
				"amount": "80562"
			}
		},
		{
			"address": "comdex1z3560xux55hlxmah5lwm5x6r0pfh9wug0w4xzw",
			"reward": {
				"denom": "ucmdx",
				"amount": "275"
			}
		},
		{
			"address": "comdex1z3c533elehr90858stvvhnhzwywwyacddg5nja",
			"reward": {
				"denom": "ucmdx",
				"amount": "70"
			}
		},
		{
			"address": "comdex1z3mwy44tc3rcr7c4dvt78ycm3a9gclat8xgewd",
			"reward": {
				"denom": "ucmdx",
				"amount": "149"
			}
		},
		{
			"address": "comdex1z3uptrkrzt8muc7e009kaa7azn35h4wdge2pzw",
			"reward": {
				"denom": "ucmdx",
				"amount": "176"
			}
		},
		{
			"address": "comdex1z3lh0ekm4du093pwlwvreklda7d9e0z9gzapr3",
			"reward": {
				"denom": "ucmdx",
				"amount": "2019"
			}
		},
		{
			"address": "comdex1zjrhwl85gnc2yanj8zl6gf3wv6f03r97twelcd",
			"reward": {
				"denom": "ucmdx",
				"amount": "128"
			}
		},
		{
			"address": "comdex1zj9dvwdhsljqz5etmturcacul60cnxwzn46rhc",
			"reward": {
				"denom": "ucmdx",
				"amount": "7216"
			}
		},
		{
			"address": "comdex1zj9jskp28hptldrg2607egwt95vyk77w0zd98q",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1zjx7r6peme554r3krx20pcfwxu9mzhs2yf5yuf",
			"reward": {
				"denom": "ucmdx",
				"amount": "1"
			}
		},
		{
			"address": "comdex1zjgz9thlrk49h0enu0jzx2gf28h56pdqx8vxjx",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1zj2er283xluk9v8sp4nlm7y3n7jghzly4mrgry",
			"reward": {
				"denom": "ucmdx",
				"amount": "1977"
			}
		},
		{
			"address": "comdex1zjtvs7na8mck7l2cyhdt8cnz2jzusxtdtv79y3",
			"reward": {
				"denom": "ucmdx",
				"amount": "2047"
			}
		},
		{
			"address": "comdex1zjtmu95dnmh07mlvfq6x8e5unmxzxeh0az4a4u",
			"reward": {
				"denom": "ucmdx",
				"amount": "43274"
			}
		},
		{
			"address": "comdex1zjhj6z0mtv7pyluvm74swz7mx9mvduv3cvrf3k",
			"reward": {
				"denom": "ucmdx",
				"amount": "12718"
			}
		},
		{
			"address": "comdex1zjef84tw8rl5z09dpxnvx9v56najtq8hvlwyeg",
			"reward": {
				"denom": "ucmdx",
				"amount": "1786"
			}
		},
		{
			"address": "comdex1znr08a9mes4kkmz982wp5tw42wz0pjpm9vf9sl",
			"reward": {
				"denom": "ucmdx",
				"amount": "5422"
			}
		},
		{
			"address": "comdex1zn23yqpur5lyduk4cdg5zney43mjzascp9gkpx",
			"reward": {
				"denom": "ucmdx",
				"amount": "28213"
			}
		},
		{
			"address": "comdex1zn2aawjcyc7zacyealzswrqpgxdhta9yhvmngg",
			"reward": {
				"denom": "ucmdx",
				"amount": "7043"
			}
		},
		{
			"address": "comdex1zndwczwzym55vxv5pntahtq2ht9ctf4j3rhp2a",
			"reward": {
				"denom": "ucmdx",
				"amount": "173"
			}
		},
		{
			"address": "comdex1znweueafq992uxfecpc64cz23z8wh7ca548udd",
			"reward": {
				"denom": "ucmdx",
				"amount": "409"
			}
		},
		{
			"address": "comdex1zn0pvk8v0xwgljj4ejnl0zkdndr584lyuhrcqz",
			"reward": {
				"denom": "ucmdx",
				"amount": "698"
			}
		},
		{
			"address": "comdex1zn08yefads2sru8v3x4dup6gqyaa2py8d0j00g",
			"reward": {
				"denom": "ucmdx",
				"amount": "5921"
			}
		},
		{
			"address": "comdex1zn52qdes2eu2zvwu7aj22vs3xuawndva7ue2le",
			"reward": {
				"denom": "ucmdx",
				"amount": "5324"
			}
		},
		{
			"address": "comdex1znkkfljclph397z0wlgcqth49pmk99dkp2an2e",
			"reward": {
				"denom": "ucmdx",
				"amount": "12903"
			}
		},
		{
			"address": "comdex1znkkjxreyfwae9m2wpau8n4wsm405rwr0mdsmz",
			"reward": {
				"denom": "ucmdx",
				"amount": "176"
			}
		},
		{
			"address": "comdex1znmc3jd0lmp569x2rsgrjkcm5u6xd2aa0r8kmg",
			"reward": {
				"denom": "ucmdx",
				"amount": "8609"
			}
		},
		{
			"address": "comdex1znlvf7mqj3k6w9uauzycmeymlf5pasqpu57n45",
			"reward": {
				"denom": "ucmdx",
				"amount": "3440"
			}
		},
		{
			"address": "comdex1z5r3hd043lmt03uc60ypmqpew3kguquw04meu8",
			"reward": {
				"denom": "ucmdx",
				"amount": "1558"
			}
		},
		{
			"address": "comdex1z5yn4n6t33uss7cxpql6jw4j0vyj52g8m36svy",
			"reward": {
				"denom": "ucmdx",
				"amount": "88"
			}
		},
		{
			"address": "comdex1z5x7q5v4v7a2pec2ml26w69v5suw970ul0cygz",
			"reward": {
				"denom": "ucmdx",
				"amount": "5378"
			}
		},
		{
			"address": "comdex1z58lfqa2005hs34qupar265v3e6crf6xmcvp7v",
			"reward": {
				"denom": "ucmdx",
				"amount": "566"
			}
		},
		{
			"address": "comdex1z5frenwnytzsw5eqmn4jgd8axl052tsy0yqges",
			"reward": {
				"denom": "ucmdx",
				"amount": "168"
			}
		},
		{
			"address": "comdex1z5fm20363luv47pmlw3cn5fru4ly75ecfc27g5",
			"reward": {
				"denom": "ucmdx",
				"amount": "410"
			}
		},
		{
			"address": "comdex1z52zc9y2czf7jtffsyvt9taamy0329s8sq2jzx",
			"reward": {
				"denom": "ucmdx",
				"amount": "1237"
			}
		},
		{
			"address": "comdex1z527g9nj856p3surrgjwk9asdv75a9zj4mdh04",
			"reward": {
				"denom": "ucmdx",
				"amount": "19"
			}
		},
		{
			"address": "comdex1z5tq60s9zxu339j9pvdud4qeaued0xptfnsw4h",
			"reward": {
				"denom": "ucmdx",
				"amount": "36734"
			}
		},
		{
			"address": "comdex1z5v4nkzyrsvgd226w5ympqjd8rfcqc5htdcz7f",
			"reward": {
				"denom": "ucmdx",
				"amount": "11158"
			}
		},
		{
			"address": "comdex1z5dlp3eehqh20lgv8u25devedkjm567a9knfsh",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1z5wdx6anglnsj50e46ukcm4tf2095ppk82f4wn",
			"reward": {
				"denom": "ucmdx",
				"amount": "10437"
			}
		},
		{
			"address": "comdex1z5sq7suwta4xyztp6f0h42dtexajjdyslz8x29",
			"reward": {
				"denom": "ucmdx",
				"amount": "1078"
			}
		},
		{
			"address": "comdex1z53jx2lpkdwa5a7u89vgf7xv6gqxukaj5kg5xe",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1z553uh3kr90mhug4ygmpsvy6w5q3s4ad0ux3hg",
			"reward": {
				"denom": "ucmdx",
				"amount": "1804"
			}
		},
		{
			"address": "comdex1z54ldlgmy50xfzldutzsxuze8fz5kz5l667tfj",
			"reward": {
				"denom": "ucmdx",
				"amount": "143"
			}
		},
		{
			"address": "comdex1z5k8fxfk957z4llk6zgqgt3tnujmlr7dmljvhm",
			"reward": {
				"denom": "ucmdx",
				"amount": "143"
			}
		},
		{
			"address": "comdex1z5ccq32ugw2pdz3s5ten0w3n3x9qmezypk2p3r",
			"reward": {
				"denom": "ucmdx",
				"amount": "14192"
			}
		},
		{
			"address": "comdex1z5ujym7z2vfamctaz2pzyu9229g8qetfa6dujy",
			"reward": {
				"denom": "ucmdx",
				"amount": "6980"
			}
		},
		{
			"address": "comdex1z5u5m59nwvu6pmtasl7433sz5e2c40fncz3r0z",
			"reward": {
				"denom": "ucmdx",
				"amount": "216905"
			}
		},
		{
			"address": "comdex1z5aj6efhpkfh8jnl2j3m70w46s2x96rnm7n762",
			"reward": {
				"denom": "ucmdx",
				"amount": "2039"
			}
		},
		{
			"address": "comdex1z57hzn2c534v2d346avppr4h0a3rd0rjadcvdl",
			"reward": {
				"denom": "ucmdx",
				"amount": "2779"
			}
		},
		{
			"address": "comdex1z4quax3k5hlj4xt7wkdkqp5l4hpgn89y83t7zw",
			"reward": {
				"denom": "ucmdx",
				"amount": "920329"
			}
		},
		{
			"address": "comdex1z4pak2ssenus02rjz2xthpcrr8qpzkhn5rz30f",
			"reward": {
				"denom": "ucmdx",
				"amount": "809"
			}
		},
		{
			"address": "comdex1z4yu58vu90np5sh04j3p62vaqwl8y2t0nj22uq",
			"reward": {
				"denom": "ucmdx",
				"amount": "623"
			}
		},
		{
			"address": "comdex1z4xrdn8e080r04lz9qyf8lq5576gg3mdlnmp07",
			"reward": {
				"denom": "ucmdx",
				"amount": "7155"
			}
		},
		{
			"address": "comdex1z4f3c8z3wd399d8y242ggr9qxqz6dh84k9zvcy",
			"reward": {
				"denom": "ucmdx",
				"amount": "13345"
			}
		},
		{
			"address": "comdex1z4tuz7anknerl504cvs8tzyz6tagplaqmryx32",
			"reward": {
				"denom": "ucmdx",
				"amount": "6101"
			}
		},
		{
			"address": "comdex1z4v55tku262akr660zv0s7e6f4hghya5h7j2sc",
			"reward": {
				"denom": "ucmdx",
				"amount": "27645"
			}
		},
		{
			"address": "comdex1z4vumgph4ymwd9d5mhxt08uuyj5ptg5r2uhnu3",
			"reward": {
				"denom": "ucmdx",
				"amount": "3367"
			}
		},
		{
			"address": "comdex1z405fjpe26trpjtgvk5k4erh29wp7llalzrjs0",
			"reward": {
				"denom": "ucmdx",
				"amount": "1959"
			}
		},
		{
			"address": "comdex1z4347pzzm0fy5v6dwaks9ldyw55hl2km566lqf",
			"reward": {
				"denom": "ucmdx",
				"amount": "2065"
			}
		},
		{
			"address": "comdex1z43m99pkg67y5k3xfsenz0s0vqk5zmfhs5ldte",
			"reward": {
				"denom": "ucmdx",
				"amount": "323"
			}
		},
		{
			"address": "comdex1z4374xv43seaxh5r2gr03h990hfaj056knv90v",
			"reward": {
				"denom": "ucmdx",
				"amount": "88"
			}
		},
		{
			"address": "comdex1z4nc7uqh6ep0tjy7gtke6250gyeynjswcs20d3",
			"reward": {
				"denom": "ucmdx",
				"amount": "7069"
			}
		},
		{
			"address": "comdex1z4k2czlve6x4wmru3xnxjd76myuaef0xcfgpt6",
			"reward": {
				"denom": "ucmdx",
				"amount": "14584"
			}
		},
		{
			"address": "comdex1z4c5fn5yy6rrr0kg87wt723240m4hcf8u9vulx",
			"reward": {
				"denom": "ucmdx",
				"amount": "19118"
			}
		},
		{
			"address": "comdex1z4mjtlxrlkmk9njp6rnhg2tcvj95cfwz3fc6me",
			"reward": {
				"denom": "ucmdx",
				"amount": "1710"
			}
		},
		{
			"address": "comdex1zkqk99407ytmzn77zgqp8hpupa7dhy2emqu804",
			"reward": {
				"denom": "ucmdx",
				"amount": "96473"
			}
		},
		{
			"address": "comdex1zkxqtnu8upawfa57vmr23ukkxm37kgqljx9mcf",
			"reward": {
				"denom": "ucmdx",
				"amount": "71102"
			}
		},
		{
			"address": "comdex1zk8a808r70tqmuf4md9xe6f2v88e8yvqer9qrx",
			"reward": {
				"denom": "ucmdx",
				"amount": "323"
			}
		},
		{
			"address": "comdex1zkwn6fycdkham0c5n9ch28cqydur3dgh6gj6uv",
			"reward": {
				"denom": "ucmdx",
				"amount": "1932"
			}
		},
		{
			"address": "comdex1zksh98qtllqr8cjqsezvmtzwmax8emf0a4v8cf",
			"reward": {
				"denom": "ucmdx",
				"amount": "309093"
			}
		},
		{
			"address": "comdex1zk3ekyresv9xenax23u2eeuf4a3yfccc5vtxda",
			"reward": {
				"denom": "ucmdx",
				"amount": "2048"
			}
		},
		{
			"address": "comdex1zk36qayzye9wv22x9p6ctcllph6le2xtj5gtsq",
			"reward": {
				"denom": "ucmdx",
				"amount": "346"
			}
		},
		{
			"address": "comdex1zknrpqac3vusrv8yx4vz5r07mh6rnqfyv8rd7t",
			"reward": {
				"denom": "ucmdx",
				"amount": "900"
			}
		},
		{
			"address": "comdex1zknc9vjyzs0p3grrpkr66lwdlp5racpw9f36q0",
			"reward": {
				"denom": "ucmdx",
				"amount": "6866"
			}
		},
		{
			"address": "comdex1zknm2w2q88ef93putmttykktepc9knznvle7jn",
			"reward": {
				"denom": "ucmdx",
				"amount": "7048"
			}
		},
		{
			"address": "comdex1zkhdv5p2fruxv3xxc0ulvyj78hgpjxvlfswamc",
			"reward": {
				"denom": "ucmdx",
				"amount": "1616"
			}
		},
		{
			"address": "comdex1zkh5e4rzsuszpe920jsa2qgnpgusmre09ckx5v",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1zk69egdn06gu09rusu2p8ps6jtqxxlj9smh5gg",
			"reward": {
				"denom": "ucmdx",
				"amount": "1775"
			}
		},
		{
			"address": "comdex1zkaa6lz4x3e5ca8p0586ly5qj0tywqrn7tuy5v",
			"reward": {
				"denom": "ucmdx",
				"amount": "5650"
			}
		},
		{
			"address": "comdex1zhp0mxtv6dcffc9a6yj48rtkfv4qvvyxjvjfd7",
			"reward": {
				"denom": "ucmdx",
				"amount": "2801"
			}
		},
		{
			"address": "comdex1zhxp0w7rjrf6jnse3vz89ray9r38yxqehp6wsq",
			"reward": {
				"denom": "ucmdx",
				"amount": "167"
			}
		},
		{
			"address": "comdex1zh82g0v6mf4gulch22umjmd9xygr496a4t3mef",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1zhgd9u24663ga7a0svrta3uhuuy4c645hxummx",
			"reward": {
				"denom": "ucmdx",
				"amount": "2181"
			}
		},
		{
			"address": "comdex1zhfxp8ppplv0qcw6r0sflyzwk50ggrgfjs25eu",
			"reward": {
				"denom": "ucmdx",
				"amount": "247"
			}
		},
		{
			"address": "comdex1zh25yp662y40q4nw76t40yaq25qg8s42prnmwg",
			"reward": {
				"denom": "ucmdx",
				"amount": "140"
			}
		},
		{
			"address": "comdex1zhtvyz0al3xzwk34w0rwerp8qrlgxzw7ycy0v6",
			"reward": {
				"denom": "ucmdx",
				"amount": "3218"
			}
		},
		{
			"address": "comdex1zhtsgpks7ydp9e7g9l4s29dsv8cn4a6fc2wep9",
			"reward": {
				"denom": "ucmdx",
				"amount": "19174"
			}
		},
		{
			"address": "comdex1zhvm5fggn750dn82akysj8gnqcz68vkvz3vpu2",
			"reward": {
				"denom": "ucmdx",
				"amount": "2244"
			}
		},
		{
			"address": "comdex1zhwx6rfg65d4hntwz9a3h4d2aqdn47kke0lxnv",
			"reward": {
				"denom": "ucmdx",
				"amount": "7270"
			}
		},
		{
			"address": "comdex1zhwanlmwa806wpkjgrcjnrty9hzw4va97csz79",
			"reward": {
				"denom": "ucmdx",
				"amount": "896"
			}
		},
		{
			"address": "comdex1zh0482ywnsalnrratk2lmg9dklhv7alxsa0fqq",
			"reward": {
				"denom": "ucmdx",
				"amount": "1237"
			}
		},
		{
			"address": "comdex1zhsfu2mkm273yehdclzvvhz56fgk73vl5wprkr",
			"reward": {
				"denom": "ucmdx",
				"amount": "2384"
			}
		},
		{
			"address": "comdex1zh3atlxnu9q0nvsz6yldjl3n76j7wha7w6n5c6",
			"reward": {
				"denom": "ucmdx",
				"amount": "2048"
			}
		},
		{
			"address": "comdex1zh3l87t36qtgau5xgg5vur8zqt4dfszehkf45u",
			"reward": {
				"denom": "ucmdx",
				"amount": "203"
			}
		},
		{
			"address": "comdex1zhkynyc8esmlvcpqq4v5jl8m33ttmc5fmmjysl",
			"reward": {
				"denom": "ucmdx",
				"amount": "213094"
			}
		},
		{
			"address": "comdex1zhhkwhncfql8f63c35m72qarjzp2mg537huneu",
			"reward": {
				"denom": "ucmdx",
				"amount": "1006"
			}
		},
		{
			"address": "comdex1zh6zalur2m4qnutnephcmy4k2a574legpant56",
			"reward": {
				"denom": "ucmdx",
				"amount": "28191"
			}
		},
		{
			"address": "comdex1zhae93zee8mhvwuy98dlmh8prphnda6ucwy7js",
			"reward": {
				"denom": "ucmdx",
				"amount": "7083"
			}
		},
		{
			"address": "comdex1zcpn44y7k692dp65xdz5w32zux6s430p2c3mxm",
			"reward": {
				"denom": "ucmdx",
				"amount": "35628"
			}
		},
		{
			"address": "comdex1zczm4l40msc5j66764tf9pzv2rx7h98zq284sx",
			"reward": {
				"denom": "ucmdx",
				"amount": "2466"
			}
		},
		{
			"address": "comdex1zcvdr8xh2vvk5mzqepvf3lvu2rhzu9f43jtkz9",
			"reward": {
				"denom": "ucmdx",
				"amount": "44210"
			}
		},
		{
			"address": "comdex1zc0d9272kuzfllgnqsk3pkcmn6xrg4ljqr66qx",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1zcnas2ufhnph43lu04ylpmj5m07p9qyedhyyms",
			"reward": {
				"denom": "ucmdx",
				"amount": "959"
			}
		},
		{
			"address": "comdex1zchscpeu5m5x6nhztnfdjmalhhh6rdf5m0ccz5",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1zchkdfa34zhqvvr55vn3gxqyk6wfqld9rna7cd",
			"reward": {
				"denom": "ucmdx",
				"amount": "7583"
			}
		},
		{
			"address": "comdex1zch695nea92jqn94sdd0nqjhnvpgfffwk8png8",
			"reward": {
				"denom": "ucmdx",
				"amount": "19999"
			}
		},
		{
			"address": "comdex1zccyu6mwrdpulcrsfkevw8thxr8pvpw4fsrgfq",
			"reward": {
				"denom": "ucmdx",
				"amount": "3223"
			}
		},
		{
			"address": "comdex1zc68l9s3sjsayac3pkr55al45y33mf2ufg8zqr",
			"reward": {
				"denom": "ucmdx",
				"amount": "19"
			}
		},
		{
			"address": "comdex1zclu338794as6f3fekrs0rkfrf3q39u9xt5ndm",
			"reward": {
				"denom": "ucmdx",
				"amount": "2460"
			}
		},
		{
			"address": "comdex1zeq3cgssfjxzardcm7vnuf2n7j4mhmvzgr9vlm",
			"reward": {
				"denom": "ucmdx",
				"amount": "7630"
			}
		},
		{
			"address": "comdex1zepz6pe5ff4muwkl7nqkgdshwjusypae9m4nk7",
			"reward": {
				"denom": "ucmdx",
				"amount": "32444"
			}
		},
		{
			"address": "comdex1zef2heg434vaqm0ss634v0hq6enkdel8xph45u",
			"reward": {
				"denom": "ucmdx",
				"amount": "1"
			}
		},
		{
			"address": "comdex1zefuuym6kypfscr0fr30xymr80avjc6uqrr6t2",
			"reward": {
				"denom": "ucmdx",
				"amount": "701"
			}
		},
		{
			"address": "comdex1ze2qd4dn0q96nmuq3d8l26r8cfafhvghp2qah4",
			"reward": {
				"denom": "ucmdx",
				"amount": "1765"
			}
		},
		{
			"address": "comdex1ze2ye5u5k3qdlexvt2e0nn0508p040944qwels",
			"reward": {
				"denom": "ucmdx",
				"amount": "1762"
			}
		},
		{
			"address": "comdex1ze2e8f4gptrttluvayhd8fjccsw8gm8vq7jsyk",
			"reward": {
				"denom": "ucmdx",
				"amount": "4287"
			}
		},
		{
			"address": "comdex1zevpqm9ktrrdz8v2vt8smnp69wmsqfvn85pkzh",
			"reward": {
				"denom": "ucmdx",
				"amount": "2849"
			}
		},
		{
			"address": "comdex1zejzxwund2xq00aw27wt6taany0d9ulrrae457",
			"reward": {
				"denom": "ucmdx",
				"amount": "718"
			}
		},
		{
			"address": "comdex1zej7u6g6jqlhyeeyq5cz48javnkq0awl0vgps9",
			"reward": {
				"denom": "ucmdx",
				"amount": "2005"
			}
		},
		{
			"address": "comdex1zen2nw9q44zrgdsxe9acnt52kyj4gznvn2hslz",
			"reward": {
				"denom": "ucmdx",
				"amount": "5453"
			}
		},
		{
			"address": "comdex1zekppcu735ltlpr4la8gmasn7pt7ajyltlumjg",
			"reward": {
				"denom": "ucmdx",
				"amount": "10698"
			}
		},
		{
			"address": "comdex1zehndhtyuq0cg0jhgjy6cn64psuk45sq66862s",
			"reward": {
				"denom": "ucmdx",
				"amount": "74"
			}
		},
		{
			"address": "comdex1zec5glry6tj3zaqjxx827gzp8am5txr5p5gl8f",
			"reward": {
				"denom": "ucmdx",
				"amount": "43598"
			}
		},
		{
			"address": "comdex1zemqgkmldkea8mws0lz55syrxtzznsjcuqn2p5",
			"reward": {
				"denom": "ucmdx",
				"amount": "56949"
			}
		},
		{
			"address": "comdex1zemwck40euwxthkev2m20px3wz2jhudjcexjy7",
			"reward": {
				"denom": "ucmdx",
				"amount": "1438"
			}
		},
		{
			"address": "comdex1zeupf9znx68t84l6huaj2tpj7cx32dep23mevg",
			"reward": {
				"denom": "ucmdx",
				"amount": "10293"
			}
		},
		{
			"address": "comdex1zelysc3s9vch8ka90lu3x9gfcwagy8kjm20fry",
			"reward": {
				"denom": "ucmdx",
				"amount": "1733"
			}
		},
		{
			"address": "comdex1zely6ek6zeegl4fl52358cv4hk6jef7ewt3ah7",
			"reward": {
				"denom": "ucmdx",
				"amount": "1374"
			}
		},
		{
			"address": "comdex1z6pqtr53tkde6xxypt7mhw8gse8ld4duwhkcav",
			"reward": {
				"denom": "ucmdx",
				"amount": "2461"
			}
		},
		{
			"address": "comdex1z6pv8774nht09pnpq9sklns752p9seppw2gtgl",
			"reward": {
				"denom": "ucmdx",
				"amount": "1791"
			}
		},
		{
			"address": "comdex1z6re8g3k2adtntwpw6qpk6kddsxn7pcetqz6xs",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1z6y6zg67j0ce85j5fd7kadf3399vvzgh9jdggz",
			"reward": {
				"denom": "ucmdx",
				"amount": "717"
			}
		},
		{
			"address": "comdex1z692deslv3tgqcf7yjzu6l0x5p23jzw5xrwgrm",
			"reward": {
				"denom": "ucmdx",
				"amount": "359"
			}
		},
		{
			"address": "comdex1z68tn9f3z4fzkq8k57wxxpznj6m0ldfcms2j7v",
			"reward": {
				"denom": "ucmdx",
				"amount": "878"
			}
		},
		{
			"address": "comdex1z6gcxstn4qse9gdkpzaxlrznwjgw22jh8sp2a9",
			"reward": {
				"denom": "ucmdx",
				"amount": "1859"
			}
		},
		{
			"address": "comdex1z6guce5rt3wtvefpz29ewtnzvqterdqj9vcfkg",
			"reward": {
				"denom": "ucmdx",
				"amount": "5641"
			}
		},
		{
			"address": "comdex1z6ntarcvd07uu6yj4ahnv0nhr84dce00psszea",
			"reward": {
				"denom": "ucmdx",
				"amount": "89264"
			}
		},
		{
			"address": "comdex1z64pnp6agc4qatq0z4xzs3mq9qn4ye55n09spd",
			"reward": {
				"denom": "ucmdx",
				"amount": "6156"
			}
		},
		{
			"address": "comdex1z64u0uvvvvafcz5ywevkt8y9udxfh37tc8z9l9",
			"reward": {
				"denom": "ucmdx",
				"amount": "5273"
			}
		},
		{
			"address": "comdex1z6kqy230t8w038re8lfne2gftz0pg8wksgfpef",
			"reward": {
				"denom": "ucmdx",
				"amount": "8"
			}
		},
		{
			"address": "comdex1z6k6cpdcthpgfk2px4lflu8ml0jkkpr4lml83h",
			"reward": {
				"denom": "ucmdx",
				"amount": "420"
			}
		},
		{
			"address": "comdex1z6huujnznhp075jy28edexgzmwq6ya03dn432e",
			"reward": {
				"denom": "ucmdx",
				"amount": "20985"
			}
		},
		{
			"address": "comdex1z6mm2l43mvdfczn0glel7f3q0pc7zvmyrtp8a9",
			"reward": {
				"denom": "ucmdx",
				"amount": "15156"
			}
		},
		{
			"address": "comdex1z6a9p43svtkzf5efmy78tae7qqz8p8cccngg6f",
			"reward": {
				"denom": "ucmdx",
				"amount": "3788"
			}
		},
		{
			"address": "comdex1z6a9d2ta4dtedemzt6vzlhphmyv6fytx7gefdw",
			"reward": {
				"denom": "ucmdx",
				"amount": "6312"
			}
		},
		{
			"address": "comdex1zmze5wt5qjh63u4qskdfycpplt4erenzp324hh",
			"reward": {
				"denom": "ucmdx",
				"amount": "261"
			}
		},
		{
			"address": "comdex1zmrwdgkect3vlt0plruyzzed6zsvpwknmjta77",
			"reward": {
				"denom": "ucmdx",
				"amount": "71707"
			}
		},
		{
			"address": "comdex1zmyx0tegzetxfvlutf09fagmugz2r044yu4d6f",
			"reward": {
				"denom": "ucmdx",
				"amount": "17577"
			}
		},
		{
			"address": "comdex1zmxqd36hhd7q7q84fv9sd7aajjrlj422hn0a4p",
			"reward": {
				"denom": "ucmdx",
				"amount": "182"
			}
		},
		{
			"address": "comdex1zmxqn8ffwhq3vp5h2alklgeujttlq78lzqw3j6",
			"reward": {
				"denom": "ucmdx",
				"amount": "2010"
			}
		},
		{
			"address": "comdex1zm8epp3rul4lsu47szwmjcuxxa5xekgfymm73s",
			"reward": {
				"denom": "ucmdx",
				"amount": "71286"
			}
		},
		{
			"address": "comdex1zmwvt8fluwu6hj9vw85edh3zr5vap2auv37sp6",
			"reward": {
				"denom": "ucmdx",
				"amount": "194"
			}
		},
		{
			"address": "comdex1zmjees6ganrja79kw8e7f5scsfdrk59hygvkxa",
			"reward": {
				"denom": "ucmdx",
				"amount": "18034"
			}
		},
		{
			"address": "comdex1zm45ymwqcfstrkakvyrsc4w4lfh4dkqgvkhjf2",
			"reward": {
				"denom": "ucmdx",
				"amount": "14095"
			}
		},
		{
			"address": "comdex1zm47g2npu467mkyhv9r8unfjsuufh6jr7gy3xx",
			"reward": {
				"denom": "ucmdx",
				"amount": "2980"
			}
		},
		{
			"address": "comdex1zmknnzeez869un39defkj0gcy5nqxaufyz79xv",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1zmkm6606vsnvtlul8yxgfl99jntdjefrlfe46z",
			"reward": {
				"denom": "ucmdx",
				"amount": "6016"
			}
		},
		{
			"address": "comdex1zmek4j5jcmm7fuf7sxf6kl2zjxrrh52lqxa5sm",
			"reward": {
				"denom": "ucmdx",
				"amount": "28324"
			}
		},
		{
			"address": "comdex1zmawmqmd8pkvm5f2433efkfraq5rpm8ykr4qxw",
			"reward": {
				"denom": "ucmdx",
				"amount": "36187"
			}
		},
		{
			"address": "comdex1zm79pz6s5tvf4dsxv2vg45hnee8hynsv0r2avd",
			"reward": {
				"denom": "ucmdx",
				"amount": "616"
			}
		},
		{
			"address": "comdex1zm7csef9pkarhfvt8h5mk9r3gw730s0qtwkdnp",
			"reward": {
				"denom": "ucmdx",
				"amount": "141"
			}
		},
		{
			"address": "comdex1zuqs7dcncty8uj3c0pplxd8339mvazvw6qtzv9",
			"reward": {
				"denom": "ucmdx",
				"amount": "6602"
			}
		},
		{
			"address": "comdex1zu9v2vp34djqa0cjle8cwylvu5euqhj7k5taut",
			"reward": {
				"denom": "ucmdx",
				"amount": "323"
			}
		},
		{
			"address": "comdex1zugnekqxxkrp0mdfva6sguzhwf3r99g2uk6700",
			"reward": {
				"denom": "ucmdx",
				"amount": "8499"
			}
		},
		{
			"address": "comdex1zut6u8z4dhzaxzu252ufe7fxudcx97cpnxeuza",
			"reward": {
				"denom": "ucmdx",
				"amount": "6060"
			}
		},
		{
			"address": "comdex1zuw5cy0d33xslueltr74htdsztwl2ccjak49he",
			"reward": {
				"denom": "ucmdx",
				"amount": "144"
			}
		},
		{
			"address": "comdex1zu0pvlf9wf9y2lpfyjmfkqjc6hpld9e6n399pc",
			"reward": {
				"denom": "ucmdx",
				"amount": "7222"
			}
		},
		{
			"address": "comdex1zujhphf37s3g6kd7axsuf7w85m23y2r254sf3x",
			"reward": {
				"denom": "ucmdx",
				"amount": "8860"
			}
		},
		{
			"address": "comdex1zuhn7530kdckmtkx7yh4l42tp9kmqa66lk39kp",
			"reward": {
				"denom": "ucmdx",
				"amount": "1594"
			}
		},
		{
			"address": "comdex1zuejdlg5x90c38fz6fsrgkglrcp4ahzsepu3q2",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1zu6ptf0y8s2tsqz80mx0p3z2ejd9c3scnqrudw",
			"reward": {
				"denom": "ucmdx",
				"amount": "853"
			}
		},
		{
			"address": "comdex1zumcq9erp4mapl6m8rksghmq9c7g2wmc958wwd",
			"reward": {
				"denom": "ucmdx",
				"amount": "1216"
			}
		},
		{
			"address": "comdex1zualynuzrv27q3y7cv2fff0r9u70h6zt77m93w",
			"reward": {
				"denom": "ucmdx",
				"amount": "7574"
			}
		},
		{
			"address": "comdex1zulemzv7e9yyzkzn3tcg7rehe47dde4jjcs0z8",
			"reward": {
				"denom": "ucmdx",
				"amount": "12902"
			}
		},
		{
			"address": "comdex1zaze2d9sr326cwfmm794nceg728sg5kqvrc922",
			"reward": {
				"denom": "ucmdx",
				"amount": "30165"
			}
		},
		{
			"address": "comdex1zay8pjl256wgc3ufcvmmcsnf6kf8jz4qstadng",
			"reward": {
				"denom": "ucmdx",
				"amount": "1524"
			}
		},
		{
			"address": "comdex1za8jtdnfd74dy4zn9qq0y6n02tfdhc6pescxkl",
			"reward": {
				"denom": "ucmdx",
				"amount": "3219"
			}
		},
		{
			"address": "comdex1zaglk5q6n70dk8mcr9x5r0lqd5lrkqvk6ycms6",
			"reward": {
				"denom": "ucmdx",
				"amount": "843"
			}
		},
		{
			"address": "comdex1za2tfq4spx4hyza38ynawutn264cjcarq55q4u",
			"reward": {
				"denom": "ucmdx",
				"amount": "24559"
			}
		},
		{
			"address": "comdex1zaw4kk7fcza0g8jllu4aag9yg882v8artepdpn",
			"reward": {
				"denom": "ucmdx",
				"amount": "28"
			}
		},
		{
			"address": "comdex1za5z4ne6hmqp9evrqdrh227pn5np05pj3zvha7",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1zacx48t2msr5fuvqc93hutjp87n0lg9n0x2fww",
			"reward": {
				"denom": "ucmdx",
				"amount": "316"
			}
		},
		{
			"address": "comdex1zaertzj3cuqs0yzw4arglrukph3vyvmrxdeqjs",
			"reward": {
				"denom": "ucmdx",
				"amount": "594"
			}
		},
		{
			"address": "comdex1za67sw3y2fcdtyypkht3435x9zpm2dh0xj2c3m",
			"reward": {
				"denom": "ucmdx",
				"amount": "33821"
			}
		},
		{
			"address": "comdex1zau5l60vgm7vcs66smgeatherr6mj3995l2d8f",
			"reward": {
				"denom": "ucmdx",
				"amount": "1"
			}
		},
		{
			"address": "comdex1za7236kpvf0k8n7j5n353rm3ae444c6zguztqj",
			"reward": {
				"denom": "ucmdx",
				"amount": "1788"
			}
		},
		{
			"address": "comdex1z7qvrdv3dym46rdz52gyzeygqemjwc688w5nw9",
			"reward": {
				"denom": "ucmdx",
				"amount": "17870"
			}
		},
		{
			"address": "comdex1z7rjzn4vzh6kqha6p0a9jfmf7jn93wh2596xnt",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1z79wmvqy6kmlwpgr4rdkqgpvadw6ygnac7x2tl",
			"reward": {
				"denom": "ucmdx",
				"amount": "9112"
			}
		},
		{
			"address": "comdex1z7293l0m85qgj8s8xd34eaxmz9j7av6auegk87",
			"reward": {
				"denom": "ucmdx",
				"amount": "537"
			}
		},
		{
			"address": "comdex1z7dlknnjk3ztcz6gs8eykc3rleqz2eu7cxzz0d",
			"reward": {
				"denom": "ucmdx",
				"amount": "26781"
			}
		},
		{
			"address": "comdex1z7wvg2hxej4ak9zlagyhlu36ghdu7pf8cxks8p",
			"reward": {
				"denom": "ucmdx",
				"amount": "35050"
			}
		},
		{
			"address": "comdex1z7s43rfwmnl55ag4zx6jqkrzz409wcq9dwqfxu",
			"reward": {
				"denom": "ucmdx",
				"amount": "130"
			}
		},
		{
			"address": "comdex1z730h7vla30v9hg3z8tmud393jsqyedmht7rny",
			"reward": {
				"denom": "ucmdx",
				"amount": "89"
			}
		},
		{
			"address": "comdex1z752hx4mptrqxc2z3mqu602yglhys9w9fvr9e2",
			"reward": {
				"denom": "ucmdx",
				"amount": "60174"
			}
		},
		{
			"address": "comdex1z75hq00728pagjzmezhml3y49hwpys5k8gfrgr",
			"reward": {
				"denom": "ucmdx",
				"amount": "1416410"
			}
		},
		{
			"address": "comdex1z7ekxv6uhlld88vwrp7kf24p9n84vcwd0v5wa6",
			"reward": {
				"denom": "ucmdx",
				"amount": "3543"
			}
		},
		{
			"address": "comdex1z7ut2nntntr669qqvlqadtj7sw28q565lz6ce6",
			"reward": {
				"denom": "ucmdx",
				"amount": "176"
			}
		},
		{
			"address": "comdex1z77a73t80uxh6tgnesfylqj7n5nt42yttxkcfu",
			"reward": {
				"denom": "ucmdx",
				"amount": "395"
			}
		},
		{
			"address": "comdex1zlqmq2jlmq66e4zfll3h20jrmynvrf2nxtr2mu",
			"reward": {
				"denom": "ucmdx",
				"amount": "123014"
			}
		},
		{
			"address": "comdex1zlpegyjjudzpe7gwl5vjpmenea58ap4v72vu95",
			"reward": {
				"denom": "ucmdx",
				"amount": "13871"
			}
		},
		{
			"address": "comdex1zlxqae2zah2v68nntkckzlu0vy08vy5vxruh9t",
			"reward": {
				"denom": "ucmdx",
				"amount": "17593"
			}
		},
		{
			"address": "comdex1zlfne473lemyeafv2jlms3x0xr6j3g0jjdpy4g",
			"reward": {
				"denom": "ucmdx",
				"amount": "3991"
			}
		},
		{
			"address": "comdex1zlfe706vtcp2v9r4hnm8vck4qs6s3u88zp4f99",
			"reward": {
				"denom": "ucmdx",
				"amount": "861"
			}
		},
		{
			"address": "comdex1zltat33twv0596lqpztem34ruj7m29x9h0pc5k",
			"reward": {
				"denom": "ucmdx",
				"amount": "2870"
			}
		},
		{
			"address": "comdex1zlw5yk99wsw0e4sucyq8guma6crdm937gl0lre",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1zl30ec3sx36jtd86g5tqhgx0jlh5gnwec96yqy",
			"reward": {
				"denom": "ucmdx",
				"amount": "3630"
			}
		},
		{
			"address": "comdex1zl54ztksmefhz8unhwxus9q5mspej5nmzpcnm3",
			"reward": {
				"denom": "ucmdx",
				"amount": "1241"
			}
		},
		{
			"address": "comdex1zlcn84t3yda0drypkwalctxeqvzc872vz0x8m3",
			"reward": {
				"denom": "ucmdx",
				"amount": "1647"
			}
		},
		{
			"address": "comdex1zlcevrzusz6ccawra8c2ff7vx7uyy4kwu2l0pg",
			"reward": {
				"denom": "ucmdx",
				"amount": "2652"
			}
		},
		{
			"address": "comdex1zlehs444wranjxs6arnw9gy8yt9jmf8am0c6yp",
			"reward": {
				"denom": "ucmdx",
				"amount": "298"
			}
		},
		{
			"address": "comdex1zlm0u44mnf90wr9450cfqm3yfzxwv5zpymqjse",
			"reward": {
				"denom": "ucmdx",
				"amount": "1800"
			}
		},
		{
			"address": "comdex1zlu953zzkl445hpnr9q4jqt3x4lhxuuqp5pssl",
			"reward": {
				"denom": "ucmdx",
				"amount": "3331"
			}
		},
		{
			"address": "comdex1zlufvmffz3cf4pp0lulv6tcgd5scqz9vwvamah",
			"reward": {
				"denom": "ucmdx",
				"amount": "38790"
			}
		},
		{
			"address": "comdex1zla0t3xe6m2uvlvda3acp44cu24j5arxsegcfg",
			"reward": {
				"denom": "ucmdx",
				"amount": "102"
			}
		},
		{
			"address": "comdex1zlas4xpx8lyu4dl39han2aq9h5t9qvtdyr0hqp",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1zlly4fmwfdznlhuw5swqdk32ujt8c6cw5j27se",
			"reward": {
				"denom": "ucmdx",
				"amount": "1747"
			}
		},
		{
			"address": "comdex1rqzq7key9nr4lmqqsywv8r873jltt658edlahp",
			"reward": {
				"denom": "ucmdx",
				"amount": "1190"
			}
		},
		{
			"address": "comdex1rqzssm6xuy52evmthj8zn3nw4294ya26v5x2ue",
			"reward": {
				"denom": "ucmdx",
				"amount": "180"
			}
		},
		{
			"address": "comdex1rqrxh56flv7rk6yypuw3f98jdvdf86jdua2022",
			"reward": {
				"denom": "ucmdx",
				"amount": "170996"
			}
		},
		{
			"address": "comdex1rqrswydpthppjrpcqjgzlffcqrql0gy6dp5mcg",
			"reward": {
				"denom": "ucmdx",
				"amount": "5333"
			}
		},
		{
			"address": "comdex1rqrck8d5u34l8p6zvanh9zn4kkd4v6sh64vddu",
			"reward": {
				"denom": "ucmdx",
				"amount": "1750"
			}
		},
		{
			"address": "comdex1rqyyttavvav6ghjn37uklutgz6a9z6ax3qqq0n",
			"reward": {
				"denom": "ucmdx",
				"amount": "60174"
			}
		},
		{
			"address": "comdex1rqx7h77w3gjvwxekg6uecnfsxjdz98hp0dfnqr",
			"reward": {
				"denom": "ucmdx",
				"amount": "6391"
			}
		},
		{
			"address": "comdex1rqgyv853uctcmc8ex5g5s8h3n58dwcrewcazm3",
			"reward": {
				"denom": "ucmdx",
				"amount": "1794"
			}
		},
		{
			"address": "comdex1rqg87xgeymk28r5seadgjl9fk0a7k3u9qgklv4",
			"reward": {
				"denom": "ucmdx",
				"amount": "1982"
			}
		},
		{
			"address": "comdex1rq2rk77v2pz04z23pucufpcnfsvk2z4kv7e7h8",
			"reward": {
				"denom": "ucmdx",
				"amount": "13430"
			}
		},
		{
			"address": "comdex1rqtppjpw9lvxxdr4lwsha5p8c9amaszypd2kyc",
			"reward": {
				"denom": "ucmdx",
				"amount": "47468"
			}
		},
		{
			"address": "comdex1rqd0wl48u6nu5622r7h5t4ca58g3j3u2jq74r8",
			"reward": {
				"denom": "ucmdx",
				"amount": "685"
			}
		},
		{
			"address": "comdex1rqwxunk8ujl3ngsxdwt05q4fplx6gprke7vj5f",
			"reward": {
				"denom": "ucmdx",
				"amount": "14472"
			}
		},
		{
			"address": "comdex1rqsnvdk38cjur62yqaffuq6ddpm6ztuuzjt8xs",
			"reward": {
				"denom": "ucmdx",
				"amount": "738"
			}
		},
		{
			"address": "comdex1rq3eqd9pmxsc44nukrw9a2jv6zfzhv5rfllamn",
			"reward": {
				"denom": "ucmdx",
				"amount": "1768"
			}
		},
		{
			"address": "comdex1rqjrpkntamtyjf0yvswctmzdkqftf5z9zkq2vl",
			"reward": {
				"denom": "ucmdx",
				"amount": "8383"
			}
		},
		{
			"address": "comdex1rqjhlm4vny43qhupcfvdvl0t3pau6hjfw0fzjv",
			"reward": {
				"denom": "ucmdx",
				"amount": "4807"
			}
		},
		{
			"address": "comdex1rq5y58saqv57m2jeqkz27lls8aucdfq0n9tsyz",
			"reward": {
				"denom": "ucmdx",
				"amount": "8982"
			}
		},
		{
			"address": "comdex1rq5hr6vyexs9c3nyjdjggda7kelwg9jt8n73hn",
			"reward": {
				"denom": "ucmdx",
				"amount": "1018"
			}
		},
		{
			"address": "comdex1rq4ffl3ywqmm52dre2hkggchhxgn49hu3d8fy4",
			"reward": {
				"denom": "ucmdx",
				"amount": "1487"
			}
		},
		{
			"address": "comdex1rqcfqmmff8qgltq2ry3d0eu8t7vq0dqv556474",
			"reward": {
				"denom": "ucmdx",
				"amount": "689"
			}
		},
		{
			"address": "comdex1rquh3r2gfkyfr46xtrg7h55sxawacjtsfme78e",
			"reward": {
				"denom": "ucmdx",
				"amount": "35374"
			}
		},
		{
			"address": "comdex1rqar83ammh0j8xkdmy676ryzn20770x66mqsnr",
			"reward": {
				"denom": "ucmdx",
				"amount": "5651"
			}
		},
		{
			"address": "comdex1rprt5hgr85m50kup8pztcq2ezev5udl6kxyf9h",
			"reward": {
				"denom": "ucmdx",
				"amount": "143"
			}
		},
		{
			"address": "comdex1rpr3l7jcjhunpdxu8ffaxy4qy58r5j7jhpn0ed",
			"reward": {
				"denom": "ucmdx",
				"amount": "80988"
			}
		},
		{
			"address": "comdex1rp9s45fegm46czm5vj2s044v7z35hhjt4v54e7",
			"reward": {
				"denom": "ucmdx",
				"amount": "1"
			}
		},
		{
			"address": "comdex1rpxkxxcm0xudn6y9gftqf872ffzav7sxkaalfn",
			"reward": {
				"denom": "ucmdx",
				"amount": "201"
			}
		},
		{
			"address": "comdex1rp80q82lz54nfytatg0w2dpwyzeuluuyyyk85v",
			"reward": {
				"denom": "ucmdx",
				"amount": "9301"
			}
		},
		{
			"address": "comdex1rp8jfnknrhteqz7regf4y87nrcztpweqtppj2n",
			"reward": {
				"denom": "ucmdx",
				"amount": "1998"
			}
		},
		{
			"address": "comdex1rpfnamlga4wgrr2q33xjegtgv6l3043dstzw30",
			"reward": {
				"denom": "ucmdx",
				"amount": "97942"
			}
		},
		{
			"address": "comdex1rpvcjeq77wes6lkrv9se28983ju9jxpgfemzpj",
			"reward": {
				"denom": "ucmdx",
				"amount": "140572"
			}
		},
		{
			"address": "comdex1rp3cfrq55t9qc2p8p3z5mzuy8khtafrwq70dth",
			"reward": {
				"denom": "ucmdx",
				"amount": "1724"
			}
		},
		{
			"address": "comdex1rpjm0ulyz679sd93zjc2j2xu5dsaqd3u5nm0p3",
			"reward": {
				"denom": "ucmdx",
				"amount": "19"
			}
		},
		{
			"address": "comdex1rp4fk7728d647rvafpgsdacsg9uallcet7fmlk",
			"reward": {
				"denom": "ucmdx",
				"amount": "2470"
			}
		},
		{
			"address": "comdex1rp43sezpnlpvds3d46l4w55rhcvh8mdvxh7ers",
			"reward": {
				"denom": "ucmdx",
				"amount": "2079"
			}
		},
		{
			"address": "comdex1rph6jlvr653vt70n5wqm2nm87sc3xzs0q79rk6",
			"reward": {
				"denom": "ucmdx",
				"amount": "174"
			}
		},
		{
			"address": "comdex1rpc00f77qs693vrpux362sca3dtw0y87ufjknq",
			"reward": {
				"denom": "ucmdx",
				"amount": "91508"
			}
		},
		{
			"address": "comdex1rpmzafq0lsnyzan4cc8s75v60lrhnyedja9m3t",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1rzqlze8gedwvw0tctxl2k9v4fx7lg6tcuxpw0j",
			"reward": {
				"denom": "ucmdx",
				"amount": "128588"
			}
		},
		{
			"address": "comdex1rzpx9c0asu3zey73h04zf4w98wmngcjk9zf7te",
			"reward": {
				"denom": "ucmdx",
				"amount": "198179"
			}
		},
		{
			"address": "comdex1rzpmtevqxcr93tw9ulqf6445cvfftqcl7nrf6a",
			"reward": {
				"denom": "ucmdx",
				"amount": "77"
			}
		},
		{
			"address": "comdex1rzrpcseyl09c2elkwqypt6egyach3z9zvxtveq",
			"reward": {
				"denom": "ucmdx",
				"amount": "2000"
			}
		},
		{
			"address": "comdex1rzr5jtfyl4fs7v8mt0aqylnhkc296s8q2qwzel",
			"reward": {
				"denom": "ucmdx",
				"amount": "61"
			}
		},
		{
			"address": "comdex1rz9d9khl7vrd7nyckjrkpzs20yvk5gc709yfak",
			"reward": {
				"denom": "ucmdx",
				"amount": "1772"
			}
		},
		{
			"address": "comdex1rzgdc23zh7d8sm3g5vxx3thtpssu99ecf64vaq",
			"reward": {
				"denom": "ucmdx",
				"amount": "35808"
			}
		},
		{
			"address": "comdex1rzf0252vyzcdhadk3xuwqmsfkhdg7tx4nzngz0",
			"reward": {
				"denom": "ucmdx",
				"amount": "2790"
			}
		},
		{
			"address": "comdex1rz3qkrhagrmn5ucc7fjunapacw8s9pjswrfcrw",
			"reward": {
				"denom": "ucmdx",
				"amount": "25"
			}
		},
		{
			"address": "comdex1rz3a9caf0pwd095rznpzfyxdad87w25c8580k3",
			"reward": {
				"denom": "ucmdx",
				"amount": "35236"
			}
		},
		{
			"address": "comdex1rznpp8rv5lkzag5ky96s06whfmv4axcv8fmehk",
			"reward": {
				"denom": "ucmdx",
				"amount": "1125"
			}
		},
		{
			"address": "comdex1rznvumap5tksk8hmaxdnjys2ec4p25jjqvdny2",
			"reward": {
				"denom": "ucmdx",
				"amount": "6046"
			}
		},
		{
			"address": "comdex1rz5kp9hf88xch8prw8l8wcdegvuwr4h2lruewe",
			"reward": {
				"denom": "ucmdx",
				"amount": "8662"
			}
		},
		{
			"address": "comdex1rz563gtr7j55x3hpev7mnn43dlcr6zaw8nrz6h",
			"reward": {
				"denom": "ucmdx",
				"amount": "15116"
			}
		},
		{
			"address": "comdex1rz6krrj8tlkphphyxxz34k9gyd079ey48rxjgn",
			"reward": {
				"denom": "ucmdx",
				"amount": "880"
			}
		},
		{
			"address": "comdex1rz6ackm5p23evrw90el8j5376m625hdm6s904v",
			"reward": {
				"denom": "ucmdx",
				"amount": "367"
			}
		},
		{
			"address": "comdex1rzm25ul3vmw5djnwzhgnthuqx0zvqyccave5gd",
			"reward": {
				"denom": "ucmdx",
				"amount": "4108"
			}
		},
		{
			"address": "comdex1rz7sqm9gndhlfkzywqfenz065z30vv9gxwkk6f",
			"reward": {
				"denom": "ucmdx",
				"amount": "71"
			}
		},
		{
			"address": "comdex1rz73s7ys2cx7753pne8j6dhtkkfw4eupvawvwt",
			"reward": {
				"denom": "ucmdx",
				"amount": "176"
			}
		},
		{
			"address": "comdex1rrqv99tnkrlhzeav8x86fyz00dg4nwset033p2",
			"reward": {
				"denom": "ucmdx",
				"amount": "13"
			}
		},
		{
			"address": "comdex1rrpmradq3mwpuemq9ukmmnmkd0pe9pyprydlrk",
			"reward": {
				"denom": "ucmdx",
				"amount": "8883"
			}
		},
		{
			"address": "comdex1rrrmtxt8gxaxt5ffrvnusm7hrq0ul2z2hmj5tr",
			"reward": {
				"denom": "ucmdx",
				"amount": "12254"
			}
		},
		{
			"address": "comdex1rry6wc9tp39rgj0sp9s0nu88jrjf3yfmqek4v2",
			"reward": {
				"denom": "ucmdx",
				"amount": "145"
			}
		},
		{
			"address": "comdex1rrdqa4ad24cs43wlmure23js8qyc5ssl7fzmcm",
			"reward": {
				"denom": "ucmdx",
				"amount": "8825"
			}
		},
		{
			"address": "comdex1rrw6fhj48a5s58rwf0w2kucc5zvymln5dyc05t",
			"reward": {
				"denom": "ucmdx",
				"amount": "3521"
			}
		},
		{
			"address": "comdex1rr00ymah9cqgxapfljphknv8uzne67ewz2g5c2",
			"reward": {
				"denom": "ucmdx",
				"amount": "623"
			}
		},
		{
			"address": "comdex1rrshfyn2uwnqc5mcgcdllrhfm4nu0c7n24aujz",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1rrseaw3xehdyslm43nrqwk44kx69phtyvctmxz",
			"reward": {
				"denom": "ucmdx",
				"amount": "2987"
			}
		},
		{
			"address": "comdex1rrngdjs75v5z8muj6s4w9q7uwmrqrzhr0q47ag",
			"reward": {
				"denom": "ucmdx",
				"amount": "1056"
			}
		},
		{
			"address": "comdex1rrn636u3vc3qvuu9lyg8eef7yqemn2t6ry75gu",
			"reward": {
				"denom": "ucmdx",
				"amount": "144"
			}
		},
		{
			"address": "comdex1rrk0hgawmg5lzlkf948f395ksqmhr867f4kgem",
			"reward": {
				"denom": "ucmdx",
				"amount": "884"
			}
		},
		{
			"address": "comdex1rrhfst9rrxun9z05mnxyvkvz5vg9gd4ht0fach",
			"reward": {
				"denom": "ucmdx",
				"amount": "268"
			}
		},
		{
			"address": "comdex1rrhmq379fk5hh750y5aw85xpp6z8emfwczdxjn",
			"reward": {
				"denom": "ucmdx",
				"amount": "1508"
			}
		},
		{
			"address": "comdex1rra9yp9zvej6msj0l5zd5gk8yygsgwyry8yq4x",
			"reward": {
				"denom": "ucmdx",
				"amount": "359"
			}
		},
		{
			"address": "comdex1ryphxjvgrre5ym2wfw7l34s87y9p7qucxhtt9q",
			"reward": {
				"denom": "ucmdx",
				"amount": "7434"
			}
		},
		{
			"address": "comdex1ryrjl3ezykj935ant3enkn8wz8jzzdkh48ayry",
			"reward": {
				"denom": "ucmdx",
				"amount": "880"
			}
		},
		{
			"address": "comdex1ryre0akzpm8dptgunewpnau93hrxh6sr0wml4y",
			"reward": {
				"denom": "ucmdx",
				"amount": "1426"
			}
		},
		{
			"address": "comdex1ry9rfhpqueyz9mlk8l0l9d4qudn79cdn0umj4l",
			"reward": {
				"denom": "ucmdx",
				"amount": "10343"
			}
		},
		{
			"address": "comdex1ryxlj2dv7ccj9zpg03qcm5zezkunve3nmu0rn3",
			"reward": {
				"denom": "ucmdx",
				"amount": "749859"
			}
		},
		{
			"address": "comdex1rygyppvm9g33grcvmrmhj9jf9ztcvcd79q3n6k",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1ry249tywevy4chyrh5qwrg8t0hvcwyp0z6ec0f",
			"reward": {
				"denom": "ucmdx",
				"amount": "7881"
			}
		},
		{
			"address": "comdex1rytcu2xljq4cch37n6wk4cvkq8ky9at7m72409",
			"reward": {
				"denom": "ucmdx",
				"amount": "3353"
			}
		},
		{
			"address": "comdex1ryv9atk6akqud9ha9q509g5jg9l9sl2v72scht",
			"reward": {
				"denom": "ucmdx",
				"amount": "5838"
			}
		},
		{
			"address": "comdex1ryw2apvkp0fnyd5cq59592uywhq0732yke9m4v",
			"reward": {
				"denom": "ucmdx",
				"amount": "149"
			}
		},
		{
			"address": "comdex1rysg83j95sqjm3z22qhtgy7arey3uqcf0c06zt",
			"reward": {
				"denom": "ucmdx",
				"amount": "5321"
			}
		},
		{
			"address": "comdex1ry3e3fvjn6e6276ym5kt8f0xr0padqtj5ee9th",
			"reward": {
				"denom": "ucmdx",
				"amount": "21580"
			}
		},
		{
			"address": "comdex1ryke5qm0zhey25clxctz6qcln0yr0ytcvgpeu7",
			"reward": {
				"denom": "ucmdx",
				"amount": "2060"
			}
		},
		{
			"address": "comdex1rye252zfralen9uep5hc8eekea4ur46lhm2zuq",
			"reward": {
				"denom": "ucmdx",
				"amount": "139"
			}
		},
		{
			"address": "comdex1ry6k5lns384d7xvpyuml5d762uv8p2g7rcsmys",
			"reward": {
				"denom": "ucmdx",
				"amount": "5685"
			}
		},
		{
			"address": "comdex1ry6l63fmffstskm77jxqmq3guu4cmvvfdet69h",
			"reward": {
				"denom": "ucmdx",
				"amount": "544"
			}
		},
		{
			"address": "comdex1ryupx0avcx5uu7u4smz5whsxy00he7z70tah3e",
			"reward": {
				"denom": "ucmdx",
				"amount": "141"
			}
		},
		{
			"address": "comdex1ryu2zpqjyf6zd6cxhlcszrg4ykt5lk2a8khfzl",
			"reward": {
				"denom": "ucmdx",
				"amount": "19030"
			}
		},
		{
			"address": "comdex1ryu5y7sqjwet76wjxdugkaturzek3vqpjap4fd",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1ry7pysmssuhat9lwujcuvq5x33pcfzylvlm8qf",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1ry7xdkejjm0fnafg60des2gcqw6663rpmulqaj",
			"reward": {
				"denom": "ucmdx",
				"amount": "8040"
			}
		},
		{
			"address": "comdex1rylxu4nc5yvkd64kk8tmtcplkmrsv8d8svjfzu",
			"reward": {
				"denom": "ucmdx",
				"amount": "13614"
			}
		},
		{
			"address": "comdex1r9yk3fqj8744shyk2sgaz48ypp6mt6clhlfpcv",
			"reward": {
				"denom": "ucmdx",
				"amount": "115839"
			}
		},
		{
			"address": "comdex1r9xzrjhkmy3wxhs4t9dyskan9y0dvqp0eh6lhe",
			"reward": {
				"denom": "ucmdx",
				"amount": "3936"
			}
		},
		{
			"address": "comdex1r9fpcvz3p6n8mpavgpjvsd7ttdecc96wrdy7wz",
			"reward": {
				"denom": "ucmdx",
				"amount": "17132"
			}
		},
		{
			"address": "comdex1r9frh3xv4ksmd58qltvwr0dtyeu47jf0n4fkkx",
			"reward": {
				"denom": "ucmdx",
				"amount": "3030"
			}
		},
		{
			"address": "comdex1r92xpez906std82ne6y22vege4dcdypruqgxdv",
			"reward": {
				"denom": "ucmdx",
				"amount": "303"
			}
		},
		{
			"address": "comdex1r9wge25wgn3httw4aetkgnwge9ejjplqjwuwz4",
			"reward": {
				"denom": "ucmdx",
				"amount": "12063"
			}
		},
		{
			"address": "comdex1r90lvn9w2t6tamlggcxwk9zvwu9gtt3e7npm88",
			"reward": {
				"denom": "ucmdx",
				"amount": "176"
			}
		},
		{
			"address": "comdex1r9szg4zt2xjfd68tw5cf30ucpdejy8ass6rk7d",
			"reward": {
				"denom": "ucmdx",
				"amount": "22384"
			}
		},
		{
			"address": "comdex1r9syd7940sekh2hpxjhmy24tjqehsxh029lmw0",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1r9np5sweyx34wl7pus6ysz5fmylyec962d20wf",
			"reward": {
				"denom": "ucmdx",
				"amount": "50790"
			}
		},
		{
			"address": "comdex1r95u9jdp97g8rk9lp7rdmlph8gmd5nhnvh7805",
			"reward": {
				"denom": "ucmdx",
				"amount": "56845"
			}
		},
		{
			"address": "comdex1r945h3s0erc5czds5lwlzxvm9c0cls9pz3wpv5",
			"reward": {
				"denom": "ucmdx",
				"amount": "27446"
			}
		},
		{
			"address": "comdex1r9k74ntguz4lyvz7wzmhuay0acx2z0pqjw02lf",
			"reward": {
				"denom": "ucmdx",
				"amount": "227"
			}
		},
		{
			"address": "comdex1r9h47qy8zwcz5axcgahufzg3yud22cp9fytzgm",
			"reward": {
				"denom": "ucmdx",
				"amount": "1774"
			}
		},
		{
			"address": "comdex1r9c00n6qh4n0kz0ux5kea4tcmfjf2nr43v9lgn",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1r9e6vddj4tewr9vh05ejapscqs24etm7zjklqd",
			"reward": {
				"denom": "ucmdx",
				"amount": "1426"
			}
		},
		{
			"address": "comdex1r9umzg5v8nvczue3l3l65j333p64xkh0qn5v0q",
			"reward": {
				"denom": "ucmdx",
				"amount": "1771"
			}
		},
		{
			"address": "comdex1r9a6xjlknunsjrw29rzk3fqjjprw0ujw26y47g",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1rx9g9rtg6g7hzzc6elht7w5dcskc9mspyvqsq7",
			"reward": {
				"denom": "ucmdx",
				"amount": "3780"
			}
		},
		{
			"address": "comdex1rxx5xrkuxhh4ma4mg9f3kd7uja75cqaazz6g59",
			"reward": {
				"denom": "ucmdx",
				"amount": "10639"
			}
		},
		{
			"address": "comdex1rxg5qu49kn7azcmntla0pc2rg69a36pqheqp53",
			"reward": {
				"denom": "ucmdx",
				"amount": "2859"
			}
		},
		{
			"address": "comdex1rx2zdndnv5wurkgvysm8l7nxcgyk7umd3e7480",
			"reward": {
				"denom": "ucmdx",
				"amount": "7193"
			}
		},
		{
			"address": "comdex1rxvpyzfwp5xa4r6l8umy087hhxfzesk36l860z",
			"reward": {
				"denom": "ucmdx",
				"amount": "142"
			}
		},
		{
			"address": "comdex1rxv2yvsgr5cergymypcs3hys5dmnnpxclg58g7",
			"reward": {
				"denom": "ucmdx",
				"amount": "14290"
			}
		},
		{
			"address": "comdex1rxdpz02sxzm30vpcptsxg8rfe66zhvdvn338he",
			"reward": {
				"denom": "ucmdx",
				"amount": "127557"
			}
		},
		{
			"address": "comdex1rxd7erv4d70fqkl82m52le0uw0mlj80k0y8jdf",
			"reward": {
				"denom": "ucmdx",
				"amount": "1786"
			}
		},
		{
			"address": "comdex1rxwnxskx7w9jjecfzjve35juw3qxughsvusnn3",
			"reward": {
				"denom": "ucmdx",
				"amount": "2852"
			}
		},
		{
			"address": "comdex1rxw5k3055swtcqfq3vnw6gd9w0guxp5p3dwvel",
			"reward": {
				"denom": "ucmdx",
				"amount": "166"
			}
		},
		{
			"address": "comdex1rxnnhqcmwyhl6u7eq709tkdgj9ylhlqvldqugf",
			"reward": {
				"denom": "ucmdx",
				"amount": "5665"
			}
		},
		{
			"address": "comdex1rxn72sgckrzq0xrz53cewglg2833d5jc8wuut8",
			"reward": {
				"denom": "ucmdx",
				"amount": "36268"
			}
		},
		{
			"address": "comdex1rx59hzyrq3r23d9ys4yjwff2q0y9gpxcpjtrwm",
			"reward": {
				"denom": "ucmdx",
				"amount": "12371"
			}
		},
		{
			"address": "comdex1rx4xgtlkrwzsm0n00dqztkrqyn8cpq0m8ncsxl",
			"reward": {
				"denom": "ucmdx",
				"amount": "2065"
			}
		},
		{
			"address": "comdex1rxkw7gt7wm4vxs6gxzpnf8ln8jgtvr20us7uzy",
			"reward": {
				"denom": "ucmdx",
				"amount": "26328"
			}
		},
		{
			"address": "comdex1rxkmasezf2n63edgvxlp7q2udwatqf33md05t0",
			"reward": {
				"denom": "ucmdx",
				"amount": "895"
			}
		},
		{
			"address": "comdex1rxhswdq9852jldxvfxs53x3aql9t927ksepkq4",
			"reward": {
				"denom": "ucmdx",
				"amount": "353"
			}
		},
		{
			"address": "comdex1rxcy8s7zaj6zmd395952g8v23f9aa8mw62g24c",
			"reward": {
				"denom": "ucmdx",
				"amount": "144"
			}
		},
		{
			"address": "comdex1rxueh47zeuh9ga3fzmrztxjxnv0g4mykwh9xl7",
			"reward": {
				"denom": "ucmdx",
				"amount": "14598"
			}
		},
		{
			"address": "comdex1rxaf5qqupxmr7sdqtrudq29th6x7mfyagp8lmc",
			"reward": {
				"denom": "ucmdx",
				"amount": "1236"
			}
		},
		{
			"address": "comdex1rx7llzdw95seekuzsnrhf8nnt6vwvaxcrwtq6e",
			"reward": {
				"denom": "ucmdx",
				"amount": "2795"
			}
		},
		{
			"address": "comdex1r8f94fykkrpxqtdsdej2sl7a7k5kmp7pqc7zfk",
			"reward": {
				"denom": "ucmdx",
				"amount": "571"
			}
		},
		{
			"address": "comdex1r8vx25y6a88qdaek59rvtda76qmn96z2atd28m",
			"reward": {
				"denom": "ucmdx",
				"amount": "158"
			}
		},
		{
			"address": "comdex1r8vmqyrlhda7nlt76t9yx5f5vx2hz0rncmntn2",
			"reward": {
				"denom": "ucmdx",
				"amount": "5713"
			}
		},
		{
			"address": "comdex1r8046fdz9nfg8ayanlg0ujxuu8h3y4uvceary8",
			"reward": {
				"denom": "ucmdx",
				"amount": "271"
			}
		},
		{
			"address": "comdex1r80mgfsjgwtyf9rad7xezun8gdum6wp7pkk9gl",
			"reward": {
				"denom": "ucmdx",
				"amount": "343"
			}
		},
		{
			"address": "comdex1r83x0xk3hqz53f7uk7p00ls8jfggaev6xszzth",
			"reward": {
				"denom": "ucmdx",
				"amount": "7479"
			}
		},
		{
			"address": "comdex1r836dhm0v2jfwemnruu33yrnngux4n60ppeznc",
			"reward": {
				"denom": "ucmdx",
				"amount": "28254"
			}
		},
		{
			"address": "comdex1r8jzsvvy9es4t32pe9ydejragudzfepp3d2x2n",
			"reward": {
				"denom": "ucmdx",
				"amount": "1419"
			}
		},
		{
			"address": "comdex1r8n42mgrv0rkhjpnt8h4wyt8zykyr9wf9qdzkv",
			"reward": {
				"denom": "ucmdx",
				"amount": "1780"
			}
		},
		{
			"address": "comdex1r8kwusdhfj7d6r20jh0uunpd2wyll0du35yxme",
			"reward": {
				"denom": "ucmdx",
				"amount": "17698"
			}
		},
		{
			"address": "comdex1r8hhuayvfwep8n0y290f3x6apxaf459udjhx9d",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1r8mgteg6t79e22f09akt4yk37ryvrgjpp7yqra",
			"reward": {
				"denom": "ucmdx",
				"amount": "749"
			}
		},
		{
			"address": "comdex1r8anr3nwpst5n0kg9ch0g9395m536uv87s952s",
			"reward": {
				"denom": "ucmdx",
				"amount": "39528"
			}
		},
		{
			"address": "comdex1rgpw69sh3enp878luu5kfg8f606melgeaqefpd",
			"reward": {
				"denom": "ucmdx",
				"amount": "678"
			}
		},
		{
			"address": "comdex1rgrvl0yxwz58fmg5qmnwcnf6feg9qf5xf2ynez",
			"reward": {
				"denom": "ucmdx",
				"amount": "553"
			}
		},
		{
			"address": "comdex1rgx95evwls47zct6w2ej8jh6efp8qa6h650u6m",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1rgxk34akfv34ekdeda59ax5a08pnu20rn4ephr",
			"reward": {
				"denom": "ucmdx",
				"amount": "20538"
			}
		},
		{
			"address": "comdex1rgf04mkxpyepmg7xxucxu5y6fntg82mtrellrp",
			"reward": {
				"denom": "ucmdx",
				"amount": "32082"
			}
		},
		{
			"address": "comdex1rg206wlgdykvjkg3d9k6j4j7exhw9uxu2px278",
			"reward": {
				"denom": "ucmdx",
				"amount": "1688"
			}
		},
		{
			"address": "comdex1rgt2r9psgrh4x0rmkqau2llrgyc66khlctkl4d",
			"reward": {
				"denom": "ucmdx",
				"amount": "303"
			}
		},
		{
			"address": "comdex1rgs259qgcwcny5euyjjc2spwkv6ttvmvh8vryh",
			"reward": {
				"denom": "ucmdx",
				"amount": "7687"
			}
		},
		{
			"address": "comdex1rgstwmdzuk4wetzsn3jywxzhmkmsu9aldnn8v8",
			"reward": {
				"denom": "ucmdx",
				"amount": "1214"
			}
		},
		{
			"address": "comdex1rg3fqyyz7sfc6gqawyxty05lkpky8q6ukuhpf8",
			"reward": {
				"denom": "ucmdx",
				"amount": "41"
			}
		},
		{
			"address": "comdex1rgnd3wftjyzm4j0jgkxlcvfhs4u6qafwvh4ck4",
			"reward": {
				"denom": "ucmdx",
				"amount": "13590"
			}
		},
		{
			"address": "comdex1rg57r7lv6ffl2srw80upm6xyssx05txy9zz8pa",
			"reward": {
				"denom": "ucmdx",
				"amount": "904"
			}
		},
		{
			"address": "comdex1rg4c6lnaufp3wnlqwewrecw5dhk84qx4tspnnc",
			"reward": {
				"denom": "ucmdx",
				"amount": "16323"
			}
		},
		{
			"address": "comdex1rgksnjtvr2q0spxce88dlykya5uz6cgjvqqnw8",
			"reward": {
				"denom": "ucmdx",
				"amount": "2552"
			}
		},
		{
			"address": "comdex1rgch3ypxx2fnjkdlcarj62mas3rpdzylh6jqju",
			"reward": {
				"denom": "ucmdx",
				"amount": "7083"
			}
		},
		{
			"address": "comdex1rgekc6qqxsjed0st5tlv2qjj5jj5c4srp6kser",
			"reward": {
				"denom": "ucmdx",
				"amount": "481"
			}
		},
		{
			"address": "comdex1rgad40wyhk0xwsry6ydxfrlqqrwj5m8afrfxrc",
			"reward": {
				"denom": "ucmdx",
				"amount": "1002"
			}
		},
		{
			"address": "comdex1rg7xdsvszrh4nqpd7vs9hs2x9cs2xgq82h2tvf",
			"reward": {
				"denom": "ucmdx",
				"amount": "204"
			}
		},
		{
			"address": "comdex1rfpljk3k6n69alr8ksx5vdsfe7qm8d3lv2nd6a",
			"reward": {
				"denom": "ucmdx",
				"amount": "88"
			}
		},
		{
			"address": "comdex1rfyweuf82wdrcu0dg6u22cxrsgnfhuuxz5lara",
			"reward": {
				"denom": "ucmdx",
				"amount": "288"
			}
		},
		{
			"address": "comdex1rfgjlsy5f8r7vy40tnyd6x578wqfvvudgxsy0c",
			"reward": {
				"denom": "ucmdx",
				"amount": "1190"
			}
		},
		{
			"address": "comdex1rfgc2vs0ntupdmdkrnc8fj4sq9e0m9u6q0jwsn",
			"reward": {
				"denom": "ucmdx",
				"amount": "101"
			}
		},
		{
			"address": "comdex1rf2he7f862qenl9rr7td9d849mlwhx8hy85nf6",
			"reward": {
				"denom": "ucmdx",
				"amount": "989"
			}
		},
		{
			"address": "comdex1rftx407fg2agnsta2ute6ckxnd82csalqd3xpq",
			"reward": {
				"denom": "ucmdx",
				"amount": "166"
			}
		},
		{
			"address": "comdex1rfvqatz6gmdj0czx72pa947ckfges6wx25mp35",
			"reward": {
				"denom": "ucmdx",
				"amount": "13039"
			}
		},
		{
			"address": "comdex1rfv7petmhsg6q4y3c4z0gx6gcuvq0cqsmf6cjw",
			"reward": {
				"denom": "ucmdx",
				"amount": "1019"
			}
		},
		{
			"address": "comdex1rfd6sts3agkprv6f08846wct5p35uhtdqhjcre",
			"reward": {
				"denom": "ucmdx",
				"amount": "1160"
			}
		},
		{
			"address": "comdex1rf3yzg6yd2jpzl4z22xmudjgh778tvnhh4knzx",
			"reward": {
				"denom": "ucmdx",
				"amount": "63490"
			}
		},
		{
			"address": "comdex1rf3vnqlzjlzvvpq33fr9u2pzqjx0sghxaa56ra",
			"reward": {
				"denom": "ucmdx",
				"amount": "1768"
			}
		},
		{
			"address": "comdex1rfj90xyl5c53zna5jvvuzmdwvnpajss89rk8yq",
			"reward": {
				"denom": "ucmdx",
				"amount": "7910"
			}
		},
		{
			"address": "comdex1rf5n4qh5ulmyhw02uuwlkptg90ucrmpzkc8zwl",
			"reward": {
				"denom": "ucmdx",
				"amount": "142"
			}
		},
		{
			"address": "comdex1rf5hprmts0h50ww7mktkly9pw9sdkyt66nk7xp",
			"reward": {
				"denom": "ucmdx",
				"amount": "1465"
			}
		},
		{
			"address": "comdex1rf4yydw4lmwr665mz3j3ekkhn6nh3r5t3zn3up",
			"reward": {
				"denom": "ucmdx",
				"amount": "168"
			}
		},
		{
			"address": "comdex1rfc6me4474jmgfc6054q4sl95ywpkf7ver53xe",
			"reward": {
				"denom": "ucmdx",
				"amount": "6771"
			}
		},
		{
			"address": "comdex1rf65vhkjamj0qx7nsu544hjzzyuzzp7gmf7vr8",
			"reward": {
				"denom": "ucmdx",
				"amount": "13905"
			}
		},
		{
			"address": "comdex1rfac3wa0dgq76s2rgc5kxmjksyfujtha76htfc",
			"reward": {
				"denom": "ucmdx",
				"amount": "34"
			}
		},
		{
			"address": "comdex1rf7dsrvv8rtxlkj90jq4vttsj0z49mrf05c8c8",
			"reward": {
				"denom": "ucmdx",
				"amount": "1575"
			}
		},
		{
			"address": "comdex1r2q6xwa6gp0yl8vjnhgzv99cz4dymuykau3jnf",
			"reward": {
				"denom": "ucmdx",
				"amount": "70288"
			}
		},
		{
			"address": "comdex1r2zwhvf6z4tw8da0w9pekk96knwud4gw5l4uyr",
			"reward": {
				"denom": "ucmdx",
				"amount": "7065"
			}
		},
		{
			"address": "comdex1r2rh0l432fkda2r8yawfady5hj7lwek09f9895",
			"reward": {
				"denom": "ucmdx",
				"amount": "29234"
			}
		},
		{
			"address": "comdex1r2rmymxq34hrqp4rf5v737jpj8ghspcltxcz2n",
			"reward": {
				"denom": "ucmdx",
				"amount": "92731"
			}
		},
		{
			"address": "comdex1r2gphsxrxr8kxul58taqgmzn7n4yqdugms06a8",
			"reward": {
				"denom": "ucmdx",
				"amount": "178"
			}
		},
		{
			"address": "comdex1r2fspw6vt773x2a5f8rk0xuc73q7lw9200xh4m",
			"reward": {
				"denom": "ucmdx",
				"amount": "369"
			}
		},
		{
			"address": "comdex1r2wjemwugdlauhm7qrnclwqgdhtk89ldrajhl5",
			"reward": {
				"denom": "ucmdx",
				"amount": "3189"
			}
		},
		{
			"address": "comdex1r2nujddqfxvw3nuvdfgsnskxgrdhpfz7nj37ln",
			"reward": {
				"denom": "ucmdx",
				"amount": "3527"
			}
		},
		{
			"address": "comdex1r25g7gvxt23nd96dkxpszefqz98e562u9fvvea",
			"reward": {
				"denom": "ucmdx",
				"amount": "7107"
			}
		},
		{
			"address": "comdex1r25mhpts8j2wjfq27cw6rdfmhp7s7kzdhh4r09",
			"reward": {
				"denom": "ucmdx",
				"amount": "8621"
			}
		},
		{
			"address": "comdex1r277pfc6tce8v3s4jl3k4jfntt44w56aa24r0j",
			"reward": {
				"denom": "ucmdx",
				"amount": "1454"
			}
		},
		{
			"address": "comdex1rtq6xrrnc977qaj7e6wnmgvsx90hpypgulcm2h",
			"reward": {
				"denom": "ucmdx",
				"amount": "815"
			}
		},
		{
			"address": "comdex1rt9p2kf33hgpt50z52rnu2vn67g4hev65gr46t",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1rt8rugxtck9p7fv2h3twepwca439k2g8npd8j6",
			"reward": {
				"denom": "ucmdx",
				"amount": "5889"
			}
		},
		{
			"address": "comdex1rt8e2rx3fzpp0cz3x9q3kyl955qnq0yj9vh7zz",
			"reward": {
				"denom": "ucmdx",
				"amount": "8888"
			}
		},
		{
			"address": "comdex1rttj8zuy97eurgckrwkffl54v0u7a7sqw5uvna",
			"reward": {
				"denom": "ucmdx",
				"amount": "181"
			}
		},
		{
			"address": "comdex1rtthrjf5aatwhlm57avpedesywk2svwnavmn4c",
			"reward": {
				"denom": "ucmdx",
				"amount": "88"
			}
		},
		{
			"address": "comdex1rtdr28me7lhzr2hvehaxp7zktqx5vpguc4czfh",
			"reward": {
				"denom": "ucmdx",
				"amount": "15054"
			}
		},
		{
			"address": "comdex1rt0mun97l0whrd66f6t78qw386h3l0tg6nhuzt",
			"reward": {
				"denom": "ucmdx",
				"amount": "313"
			}
		},
		{
			"address": "comdex1rtsce8m32ah3fnp9lu5ut9mjkczd7wl6j9el2w",
			"reward": {
				"denom": "ucmdx",
				"amount": "17599"
			}
		},
		{
			"address": "comdex1rt3qgjs5g5mu79mfxwm6cx4nj9ufg0e6lyegxh",
			"reward": {
				"denom": "ucmdx",
				"amount": "17671"
			}
		},
		{
			"address": "comdex1rtnj905gy6nq46rw8nye7kcc2h94e4kp0a3tva",
			"reward": {
				"denom": "ucmdx",
				"amount": "12597"
			}
		},
		{
			"address": "comdex1rtn4ydkuqpurcwvycc5pquwzmfp0rwvh83sw9u",
			"reward": {
				"denom": "ucmdx",
				"amount": "7107"
			}
		},
		{
			"address": "comdex1rt5vpu5uyxvxh4r8a3ugqht8vk5rkayl5gfw34",
			"reward": {
				"denom": "ucmdx",
				"amount": "5632"
			}
		},
		{
			"address": "comdex1rta9aa5dzsr3kvcf2x8x8zeycd0yw4stfpqf9t",
			"reward": {
				"denom": "ucmdx",
				"amount": "1829"
			}
		},
		{
			"address": "comdex1rv9226wsqtcs4rpzthkch23s6hve6uxwxv7swh",
			"reward": {
				"denom": "ucmdx",
				"amount": "6792"
			}
		},
		{
			"address": "comdex1rv9czkxhnwxc3788qml4cx7hpzeqxzc7yy36yy",
			"reward": {
				"denom": "ucmdx",
				"amount": "19"
			}
		},
		{
			"address": "comdex1rvx0xz6pj69m495ks39rfv4dlp35u2hd97f026",
			"reward": {
				"denom": "ucmdx",
				"amount": "14394"
			}
		},
		{
			"address": "comdex1rv2mp5ravyurqcn7dlu0hg24vyum9uc000vemq",
			"reward": {
				"denom": "ucmdx",
				"amount": "6263"
			}
		},
		{
			"address": "comdex1rvwskwvy9xsj98kuve84jlah5ln49v2fnhqmxr",
			"reward": {
				"denom": "ucmdx",
				"amount": "18890"
			}
		},
		{
			"address": "comdex1rv0jvzu7vjv6hk6nslrv7mzmdk89w7anqemufr",
			"reward": {
				"denom": "ucmdx",
				"amount": "2683"
			}
		},
		{
			"address": "comdex1rvj2fahykm3hmmzqeleeylt2nhj9xah8mzcmfe",
			"reward": {
				"denom": "ucmdx",
				"amount": "17792"
			}
		},
		{
			"address": "comdex1rvjlwh77dphcj026u9dj0z257c9s28qdlvjwlv",
			"reward": {
				"denom": "ucmdx",
				"amount": "25"
			}
		},
		{
			"address": "comdex1rvnmsrr4hdnf44wzp944xsej2afup6m40tc86k",
			"reward": {
				"denom": "ucmdx",
				"amount": "7035"
			}
		},
		{
			"address": "comdex1rv52shjza8lv7pv4avr24nqpqmq4z90y9qsgsc",
			"reward": {
				"denom": "ucmdx",
				"amount": "6953"
			}
		},
		{
			"address": "comdex1rv45jtcrv660rhku2r40kvantdz4xfc8swkzlj",
			"reward": {
				"denom": "ucmdx",
				"amount": "16"
			}
		},
		{
			"address": "comdex1rvkqktatv8dv9lhm73y45qvczqzlnp3ymlf88v",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1rvcwvlp2vf95l3pxjkk80np4n7r98yqkszjezc",
			"reward": {
				"denom": "ucmdx",
				"amount": "166"
			}
		},
		{
			"address": "comdex1rvm2yp6lr58pdc9zdunnqdzxgdvvwtkm80a5xz",
			"reward": {
				"denom": "ucmdx",
				"amount": "1726"
			}
		},
		{
			"address": "comdex1rdqf2gqpfncd06u6hjexjswpsss32d66xya90f",
			"reward": {
				"denom": "ucmdx",
				"amount": "16642"
			}
		},
		{
			"address": "comdex1rdq536awt627r95mxdte3y0vf7fup832mrdj0w",
			"reward": {
				"denom": "ucmdx",
				"amount": "732"
			}
		},
		{
			"address": "comdex1rd96vv2uqtehkvj3ju42y5cd68ywukwumnf4j4",
			"reward": {
				"denom": "ucmdx",
				"amount": "184"
			}
		},
		{
			"address": "comdex1rdxzema2k3j5mqn940hwum4m04j466cllh9hrh",
			"reward": {
				"denom": "ucmdx",
				"amount": "1452"
			}
		},
		{
			"address": "comdex1rdghq6rw5xfg9qe6vmn377vt97j59escr4yeyy",
			"reward": {
				"denom": "ucmdx",
				"amount": "61998"
			}
		},
		{
			"address": "comdex1rd264dy5wh2kxk86558lvyxkdh9r2rrzntr8ql",
			"reward": {
				"denom": "ucmdx",
				"amount": "407"
			}
		},
		{
			"address": "comdex1rdvtj2eem8gyualzhw8295slk443prfngm2mvt",
			"reward": {
				"denom": "ucmdx",
				"amount": "41561"
			}
		},
		{
			"address": "comdex1rddtrewxpnqz6e6s0gsevaaah83hkplrtq9mjz",
			"reward": {
				"denom": "ucmdx",
				"amount": "285"
			}
		},
		{
			"address": "comdex1rdwdk8x735yf7d284zr2shjvvlxxshm2gsw43h",
			"reward": {
				"denom": "ucmdx",
				"amount": "215642"
			}
		},
		{
			"address": "comdex1rd335jfxa0ytxp4adkc37cnmt49snel2rj79gk",
			"reward": {
				"denom": "ucmdx",
				"amount": "1765"
			}
		},
		{
			"address": "comdex1rd3ep5eqqusarda6pu3u0dwrdfwlx6lmueqdxy",
			"reward": {
				"denom": "ucmdx",
				"amount": "37584"
			}
		},
		{
			"address": "comdex1rdjp4g6tw4vs8nm7awy25cr893ydqzgt0mllgu",
			"reward": {
				"denom": "ucmdx",
				"amount": "568"
			}
		},
		{
			"address": "comdex1rd4rg8v3v4m7rlhytgcg0esyw7de5th7269aa0",
			"reward": {
				"denom": "ucmdx",
				"amount": "5309"
			}
		},
		{
			"address": "comdex1rd42684uvsw2732nnp60ufmw33awrjdlkedfnv",
			"reward": {
				"denom": "ucmdx",
				"amount": "173"
			}
		},
		{
			"address": "comdex1rdk8qx34q55pt3r8m9lsx268zufes8ew5tnglp",
			"reward": {
				"denom": "ucmdx",
				"amount": "14167"
			}
		},
		{
			"address": "comdex1rdhwv67pffuyxawqmkafym6vx4uq4tvzfj347y",
			"reward": {
				"denom": "ucmdx",
				"amount": "1802052"
			}
		},
		{
			"address": "comdex1rde3fuhlz9c07h7uxvsld2vczc2rwm64ta5ufn",
			"reward": {
				"denom": "ucmdx",
				"amount": "1357"
			}
		},
		{
			"address": "comdex1rd7p3u3y5skxk08j5p4ugffua8vfr7xuspcjg8",
			"reward": {
				"denom": "ucmdx",
				"amount": "87"
			}
		},
		{
			"address": "comdex1rd78ktc3j8mhnvr8kt4mmyyex89uegywl842a7",
			"reward": {
				"denom": "ucmdx",
				"amount": "199"
			}
		},
		{
			"address": "comdex1rwrp4765cxyhawa54emrte8j8rpgq098qk2drl",
			"reward": {
				"denom": "ucmdx",
				"amount": "4654"
			}
		},
		{
			"address": "comdex1rwytctrqqaq08mdd4c7v57sqah44guxvu47lva",
			"reward": {
				"denom": "ucmdx",
				"amount": "7180"
			}
		},
		{
			"address": "comdex1rwy3ksejz79nvsx3ajqft8twxgtwaadrar54ur",
			"reward": {
				"denom": "ucmdx",
				"amount": "84"
			}
		},
		{
			"address": "comdex1rw994l234kcd37q6j9c95x56rpqnh8fp7xh74w",
			"reward": {
				"denom": "ucmdx",
				"amount": "613"
			}
		},
		{
			"address": "comdex1rwxakwwswkmn60uhyglvzrx49ydgww653y2at5",
			"reward": {
				"denom": "ucmdx",
				"amount": "12407"
			}
		},
		{
			"address": "comdex1rw8zczwpsehy0jwl2dhvwz48kz3w624vs4zqrh",
			"reward": {
				"denom": "ucmdx",
				"amount": "394"
			}
		},
		{
			"address": "comdex1rwgy4khg7m3lduyzuuvxa62h46w99kl26ejcvn",
			"reward": {
				"denom": "ucmdx",
				"amount": "110324"
			}
		},
		{
			"address": "comdex1rwfljkxymf0sj7909hknsel7uaf6kvkraarf3j",
			"reward": {
				"denom": "ucmdx",
				"amount": "4268"
			}
		},
		{
			"address": "comdex1rw22mlq6ak2ypl8xeqtxfcz7lt2ctm55m2q7dg",
			"reward": {
				"denom": "ucmdx",
				"amount": "2280"
			}
		},
		{
			"address": "comdex1rwdjjmzq9e8wx9255676lak7uzqzh4yjtutk5e",
			"reward": {
				"denom": "ucmdx",
				"amount": "44"
			}
		},
		{
			"address": "comdex1rwwpegehvlrcwl068mt0ftc55h5ez63xxrj3xj",
			"reward": {
				"denom": "ucmdx",
				"amount": "182"
			}
		},
		{
			"address": "comdex1rwn7uhk9mnvkcgjqek7p597ff4jnlgju5yzjzx",
			"reward": {
				"denom": "ucmdx",
				"amount": "1775"
			}
		},
		{
			"address": "comdex1rwnlgzg9gl7s22cedwc773ejav932phss2vgky",
			"reward": {
				"denom": "ucmdx",
				"amount": "247"
			}
		},
		{
			"address": "comdex1rw5wgsd9vgnn0czyrhlwk65ghz0p3ynglg7plm",
			"reward": {
				"denom": "ucmdx",
				"amount": "1757"
			}
		},
		{
			"address": "comdex1rw5eyc2cdjjqv06zgu97pl8lfr2sjr5lc9fscf",
			"reward": {
				"denom": "ucmdx",
				"amount": "533"
			}
		},
		{
			"address": "comdex1rwhqvkm43y0v4af66efa0gk2k27ftcwtqv8sew",
			"reward": {
				"denom": "ucmdx",
				"amount": "1670"
			}
		},
		{
			"address": "comdex1rweyjna4mu37s8aygsaxcz9tk4aa3rwt0p6slr",
			"reward": {
				"denom": "ucmdx",
				"amount": "889"
			}
		},
		{
			"address": "comdex1rw70cyc0jchm6825qtegfl05n40fnqukejy6w4",
			"reward": {
				"denom": "ucmdx",
				"amount": "91472"
			}
		},
		{
			"address": "comdex1rw76mkr59cecdkqr5r4h56xvwuxjtf54q8syur",
			"reward": {
				"denom": "ucmdx",
				"amount": "9938"
			}
		},
		{
			"address": "comdex1rwldprss28md49qkdh5w05e86f4q8ah4afgcw2",
			"reward": {
				"denom": "ucmdx",
				"amount": "15081"
			}
		},
		{
			"address": "comdex1r0qpy0ckahfzs45q0vzqks2a7spx34zmq9r7jz",
			"reward": {
				"denom": "ucmdx",
				"amount": "993"
			}
		},
		{
			"address": "comdex1r08c5w29x4rwjqjfa0fkx73drwtpkyq4gfqemv",
			"reward": {
				"denom": "ucmdx",
				"amount": "1580"
			}
		},
		{
			"address": "comdex1r0gt0ngrdpkujucq7xmfslej6gc0nkcwxefrzu",
			"reward": {
				"denom": "ucmdx",
				"amount": "28649"
			}
		},
		{
			"address": "comdex1r0ghv5k4g0xcved3v3f7u8fug5l77rk5uq86qa",
			"reward": {
				"denom": "ucmdx",
				"amount": "1380"
			}
		},
		{
			"address": "comdex1r000nhpy9fvtm3km404tngn8y8qyv4fp8yqh9d",
			"reward": {
				"denom": "ucmdx",
				"amount": "1434"
			}
		},
		{
			"address": "comdex1r003xgy3c8srztwdwmeye67w8neap25san85n9",
			"reward": {
				"denom": "ucmdx",
				"amount": "20928"
			}
		},
		{
			"address": "comdex1r0s8fxs8uugk5zzlxd4308g3mjv2vtncayn3q6",
			"reward": {
				"denom": "ucmdx",
				"amount": "3485"
			}
		},
		{
			"address": "comdex1r04kt4fqr26rq9z6zhglkqqrp98qlyxrq6rzg8",
			"reward": {
				"denom": "ucmdx",
				"amount": "2346"
			}
		},
		{
			"address": "comdex1r0hv3p0zx26hn0fs2l2h8rmjcpfu0lna53qyk7",
			"reward": {
				"denom": "ucmdx",
				"amount": "12760"
			}
		},
		{
			"address": "comdex1r0mjrsjfprau9rwwxe04hkzfud696n3x84nhmm",
			"reward": {
				"denom": "ucmdx",
				"amount": "2010"
			}
		},
		{
			"address": "comdex1r0mlu77nlkqlll55uvqfpsnyyt9ndcrlvpnyx4",
			"reward": {
				"denom": "ucmdx",
				"amount": "6440"
			}
		},
		{
			"address": "comdex1r0u7nmj9z3559wzwuymt06ye48m682vwvg8qws",
			"reward": {
				"denom": "ucmdx",
				"amount": "16"
			}
		},
		{
			"address": "comdex1r0l2y46802hvlznrnuds9y04t4mak5gfdqp2yh",
			"reward": {
				"denom": "ucmdx",
				"amount": "172"
			}
		},
		{
			"address": "comdex1rszu34t0zhc6x8g3d579xlxk35cl6mxms2pnya",
			"reward": {
				"denom": "ucmdx",
				"amount": "1415"
			}
		},
		{
			"address": "comdex1rsyp9aypxa8tmgtc7rucqyssl6muf887a2pmw8",
			"reward": {
				"denom": "ucmdx",
				"amount": "555"
			}
		},
		{
			"address": "comdex1rs8ag8xqsyzlymt5599rl4a8kmqtt3ncg0gy3e",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1rswpf59rmu7ky6nsffyz6ruuy5u030fd4ca740",
			"reward": {
				"denom": "ucmdx",
				"amount": "6297"
			}
		},
		{
			"address": "comdex1rsw7rpfe4wr00gmq2c9c652urtn5flyqmy5krx",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1rsnx075yzdw32tjenj7qydj5annh630p9wldcl",
			"reward": {
				"denom": "ucmdx",
				"amount": "1990"
			}
		},
		{
			"address": "comdex1rs4qsusunnehxzcgmg4x6s37kpdncrhdz9mswf",
			"reward": {
				"denom": "ucmdx",
				"amount": "10436"
			}
		},
		{
			"address": "comdex1rs4asrnd94yd850txd5x3j46ltccd6et034gd2",
			"reward": {
				"denom": "ucmdx",
				"amount": "227"
			}
		},
		{
			"address": "comdex1rsks3gfv2rcf3atj8gs09yyhrftyn7j5rgdhrm",
			"reward": {
				"denom": "ucmdx",
				"amount": "140"
			}
		},
		{
			"address": "comdex1rs6yum0lgxu2yde5gy0cgua6vavfua0yj7hczr",
			"reward": {
				"denom": "ucmdx",
				"amount": "2321"
			}
		},
		{
			"address": "comdex1rs60hpfysn4ff5c2a4gf5wnnzpd27g0sg8r0zu",
			"reward": {
				"denom": "ucmdx",
				"amount": "1253"
			}
		},
		{
			"address": "comdex1rsmxq269hzgxqdah6mf5xdqvqm0rrhsnr808dl",
			"reward": {
				"denom": "ucmdx",
				"amount": "6924"
			}
		},
		{
			"address": "comdex1rslg7tjk2x6qhg6xaels6lflcx4k6vlj8cf5n8",
			"reward": {
				"denom": "ucmdx",
				"amount": "144"
			}
		},
		{
			"address": "comdex1r3pvn9yawgfych2vmd953fenu5p5xu3mkvsvxz",
			"reward": {
				"denom": "ucmdx",
				"amount": "151"
			}
		},
		{
			"address": "comdex1r3r377258lyfga4uxcgvxtlvmae6yu7rw7yjzy",
			"reward": {
				"denom": "ucmdx",
				"amount": "1427"
			}
		},
		{
			"address": "comdex1r3xmfmxah0pm0acvsu0pahg8mx4ayjkzmueks3",
			"reward": {
				"denom": "ucmdx",
				"amount": "581"
			}
		},
		{
			"address": "comdex1r32utmdljapuz7a94qtahd89j0m4m9lzdvwl5a",
			"reward": {
				"denom": "ucmdx",
				"amount": "196"
			}
		},
		{
			"address": "comdex1r3vdmrjpd29zfh58k5n9p7ccl9x6a5wmtj8uu6",
			"reward": {
				"denom": "ucmdx",
				"amount": "2937"
			}
		},
		{
			"address": "comdex1r3sglnjc34q2l0nf2rl2dh6slum564mmq3e8xp",
			"reward": {
				"denom": "ucmdx",
				"amount": "6908"
			}
		},
		{
			"address": "comdex1r3sfq6qfxmnna9ldfl5mf4cxkdfj0fac3u3t3r",
			"reward": {
				"denom": "ucmdx",
				"amount": "104867"
			}
		},
		{
			"address": "comdex1r3s4d5cmh449akj40stvyjhckvmjpw7yyez9zc",
			"reward": {
				"denom": "ucmdx",
				"amount": "19811"
			}
		},
		{
			"address": "comdex1r3hd68phwrwldulgmj5me9g2klsxh9rt7cyp97",
			"reward": {
				"denom": "ucmdx",
				"amount": "1190"
			}
		},
		{
			"address": "comdex1r3my6p39lmmu2lekc0pdae6mj46def2586nlph",
			"reward": {
				"denom": "ucmdx",
				"amount": "8611"
			}
		},
		{
			"address": "comdex1r3780rwha6l00p438xuz0jj9jlasrnezzeu9av",
			"reward": {
				"denom": "ucmdx",
				"amount": "18939"
			}
		},
		{
			"address": "comdex1r3lhpuag4ks7kxcrd3y56ktu7z0mvhtldyck8n",
			"reward": {
				"denom": "ucmdx",
				"amount": "2502"
			}
		},
		{
			"address": "comdex1rjphl3h9rlmmmucyndrj04h4mtexplvkdjuag7",
			"reward": {
				"denom": "ucmdx",
				"amount": "2841"
			}
		},
		{
			"address": "comdex1rjzsjxsfel6wvgynrw54y36mt4hxrpnl7npkzn",
			"reward": {
				"denom": "ucmdx",
				"amount": "88"
			}
		},
		{
			"address": "comdex1rjxd6qmdxnruf27tsr8z7xfw9fcgs80l2h27k7",
			"reward": {
				"denom": "ucmdx",
				"amount": "2365"
			}
		},
		{
			"address": "comdex1rjgff0l2l8rd3gza08w9m50qttgwe3h2l5eexf",
			"reward": {
				"denom": "ucmdx",
				"amount": "167"
			}
		},
		{
			"address": "comdex1rjg0vg9hlu8kygnrdehzt3z0t4ea5ssul4rsda",
			"reward": {
				"denom": "ucmdx",
				"amount": "6226"
			}
		},
		{
			"address": "comdex1rjf699kz55upwy8ggkgnl8r9rv3epavzfp3esv",
			"reward": {
				"denom": "ucmdx",
				"amount": "1753"
			}
		},
		{
			"address": "comdex1rjt58y27yex5cvmw52csmamxk06n2zkupdwhws",
			"reward": {
				"denom": "ucmdx",
				"amount": "12420579"
			}
		},
		{
			"address": "comdex1rjwkcgz0rgz4ytqez8dc076nhf9dkxn9cgtstq",
			"reward": {
				"denom": "ucmdx",
				"amount": "1419"
			}
		},
		{
			"address": "comdex1rj0ynahrwgtnd5x4ffflzd2kzg6jshu62kanz3",
			"reward": {
				"denom": "ucmdx",
				"amount": "6803"
			}
		},
		{
			"address": "comdex1rj3xkgz5h4ahmk7ftdy3fjnnqd65lajwlwz093",
			"reward": {
				"denom": "ucmdx",
				"amount": "1562"
			}
		},
		{
			"address": "comdex1rj565ettms8l8ffud0w0ldan3d4vjvpuwpn4ys",
			"reward": {
				"denom": "ucmdx",
				"amount": "409"
			}
		},
		{
			"address": "comdex1rj4k76hr8u3ra60crjpggt483x0tl60uedhz72",
			"reward": {
				"denom": "ucmdx",
				"amount": "298763"
			}
		},
		{
			"address": "comdex1rjct35zmggzpzu3e3wspq7yrypyyh4y90sjz7c",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1rjc0xq3ygy8qvjwjlf3jkhmjdzgy24tcaf64yc",
			"reward": {
				"denom": "ucmdx",
				"amount": "271"
			}
		},
		{
			"address": "comdex1rjcjnrf6pqcjh7nt8kmfhf8wa2k534l0x4u37p",
			"reward": {
				"denom": "ucmdx",
				"amount": "2800"
			}
		},
		{
			"address": "comdex1rjm5q68kxqzkxlqgalavk7rk2mle4vdk9u8zv9",
			"reward": {
				"denom": "ucmdx",
				"amount": "71970"
			}
		},
		{
			"address": "comdex1rjue4gtl04ucphfqdum8ndlaxvzk09877x0c25",
			"reward": {
				"denom": "ucmdx",
				"amount": "317452"
			}
		},
		{
			"address": "comdex1rjacag2azwmyump8q2psnzptfpjzt6weg6kr04",
			"reward": {
				"denom": "ucmdx",
				"amount": "72"
			}
		},
		{
			"address": "comdex1rjlr8gdl6kxhqlgrndaufu2n56je0fsq5devpg",
			"reward": {
				"denom": "ucmdx",
				"amount": "720876"
			}
		},
		{
			"address": "comdex1rjl2f7nwxv0arg5yvxrpgreq3jfgapy06hfh08",
			"reward": {
				"denom": "ucmdx",
				"amount": "1436"
			}
		},
		{
			"address": "comdex1rnq9tss3uwk9htafk0lwf5chkfuhnyar7sh2lr",
			"reward": {
				"denom": "ucmdx",
				"amount": "66353"
			}
		},
		{
			"address": "comdex1rnryv26aq043e4a7mg62mv6ndkpea2f96wlce7",
			"reward": {
				"denom": "ucmdx",
				"amount": "14846"
			}
		},
		{
			"address": "comdex1rnx7zjznm4urxuq88pt722r8frhe0cvurxvd66",
			"reward": {
				"denom": "ucmdx",
				"amount": "124"
			}
		},
		{
			"address": "comdex1rn28lfuhmcm8nwvhk9gve2hgjzymwxghrq48hn",
			"reward": {
				"denom": "ucmdx",
				"amount": "28630"
			}
		},
		{
			"address": "comdex1rn2hgqdfe2985ce4he4gh68s2rc4p57v9d9lzp",
			"reward": {
				"denom": "ucmdx",
				"amount": "7423"
			}
		},
		{
			"address": "comdex1rntv887vmzu5f7wt38lwdgsqu3stpmxynnuzx9",
			"reward": {
				"denom": "ucmdx",
				"amount": "176"
			}
		},
		{
			"address": "comdex1rnd99lt4t7lf85y7mf653w35wp7jj4sl3zfhyw",
			"reward": {
				"denom": "ucmdx",
				"amount": "3870"
			}
		},
		{
			"address": "comdex1rnnuk7rmjwj3ep97s5d0w5ek5che53ek8qnnm3",
			"reward": {
				"denom": "ucmdx",
				"amount": "5111"
			}
		},
		{
			"address": "comdex1rn54u9kxkverhcvkqg67z6dt5a303ak3wnrla6",
			"reward": {
				"denom": "ucmdx",
				"amount": "271"
			}
		},
		{
			"address": "comdex1rnchnq9w7n234es4t7g5jgamddfek7tqjwr34x",
			"reward": {
				"denom": "ucmdx",
				"amount": "2867"
			}
		},
		{
			"address": "comdex1rnmfqlfhp7apsdaaeum57rauz2s66zjac67sx0",
			"reward": {
				"denom": "ucmdx",
				"amount": "5940"
			}
		},
		{
			"address": "comdex1rnmka79q0p5m29ujtrlz3pfy2kqn847j0v7pfm",
			"reward": {
				"denom": "ucmdx",
				"amount": "291"
			}
		},
		{
			"address": "comdex1rnmect2sakfhx2xc9as085ddfpep2jz5wjdsyr",
			"reward": {
				"denom": "ucmdx",
				"amount": "27"
			}
		},
		{
			"address": "comdex1r5qkmvn9hnv0pugejr73639w07d2mugh55uzl2",
			"reward": {
				"denom": "ucmdx",
				"amount": "86847"
			}
		},
		{
			"address": "comdex1r5pkkpcg0436puweny46kmkytkehnrhyxnc6yh",
			"reward": {
				"denom": "ucmdx",
				"amount": "14402"
			}
		},
		{
			"address": "comdex1r5g2ngyjzq29xk3ccw230fh376v5l0250hzyk9",
			"reward": {
				"denom": "ucmdx",
				"amount": "1777"
			}
		},
		{
			"address": "comdex1r5gjlh6dwgp32zls3svc9kfdse6v2eelcnnnwy",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1r5fpwsprsw70l04cz33lh3rrud54xj23cfudu4",
			"reward": {
				"denom": "ucmdx",
				"amount": "6189"
			}
		},
		{
			"address": "comdex1r5f9ruu0tc2tpdt2y6qfegkhtxdv4k2un6ql9e",
			"reward": {
				"denom": "ucmdx",
				"amount": "15209"
			}
		},
		{
			"address": "comdex1r5trk7akrgukgcvvf852kwzgzydmla3r2kwq9g",
			"reward": {
				"denom": "ucmdx",
				"amount": "8781"
			}
		},
		{
			"address": "comdex1r5tjkj5x4kc98n5syy9spputcvgsfgq0h73fkk",
			"reward": {
				"denom": "ucmdx",
				"amount": "6565"
			}
		},
		{
			"address": "comdex1r5sv9uqx7zfd7epg806z82z3lk30k7fand8fjk",
			"reward": {
				"denom": "ucmdx",
				"amount": "57110"
			}
		},
		{
			"address": "comdex1r5s36chuefsknf8y0kmcma5r64kzfz4h3cpxaz",
			"reward": {
				"denom": "ucmdx",
				"amount": "6494"
			}
		},
		{
			"address": "comdex1r5suv47u6ludwrlu4w390huwccx7u8drmad28r",
			"reward": {
				"denom": "ucmdx",
				"amount": "147"
			}
		},
		{
			"address": "comdex1r53fsf392ty2ag6e3edj9722muqhnptdmltxvj",
			"reward": {
				"denom": "ucmdx",
				"amount": "6899"
			}
		},
		{
			"address": "comdex1r54xcpm6ujh3tkr3m0x90fppc3fqrtgay0pw63",
			"reward": {
				"denom": "ucmdx",
				"amount": "3360"
			}
		},
		{
			"address": "comdex1r5hqmrywk33ffchwqgca8yndw3q6ryaa3gjvda",
			"reward": {
				"denom": "ucmdx",
				"amount": "12716"
			}
		},
		{
			"address": "comdex1r5cwl8zk7pryh73md8r5sfvvxuthlmurfdm7rl",
			"reward": {
				"denom": "ucmdx",
				"amount": "17196"
			}
		},
		{
			"address": "comdex1r5muk8te0z2url8kntgjst7wygnfmr8xrrdlrz",
			"reward": {
				"denom": "ucmdx",
				"amount": "20481"
			}
		},
		{
			"address": "comdex1r5agl7dka6ernsl6rmyungvhrh76pmadvg5kmg",
			"reward": {
				"denom": "ucmdx",
				"amount": "8633"
			}
		},
		{
			"address": "comdex1r5a6qswz5pwakxdlc3gh4mze7pn4rwa5e5yta4",
			"reward": {
				"denom": "ucmdx",
				"amount": "8152"
			}
		},
		{
			"address": "comdex1r57zznvh43zyjhe3ru2mjqvsnzygl7rr7lvqsx",
			"reward": {
				"denom": "ucmdx",
				"amount": "947"
			}
		},
		{
			"address": "comdex1r57687yledhc4qenf33cpcjg6x3wf2qh7axalt",
			"reward": {
				"denom": "ucmdx",
				"amount": "2010"
			}
		},
		{
			"address": "comdex1r4qaw2swemz0s8t28snq5z6hsxjtnwnj8h0zk9",
			"reward": {
				"denom": "ucmdx",
				"amount": "10885"
			}
		},
		{
			"address": "comdex1r4rkqyxu0fej8g884874rv9wuwp86gaa6afvwq",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1r4r77hw4ft6jn7zny7tz29st4tkqdm9hwjpm4k",
			"reward": {
				"denom": "ucmdx",
				"amount": "1220"
			}
		},
		{
			"address": "comdex1r4rlplpges8qp9znjkuu680z6lgjqr6muumfjw",
			"reward": {
				"denom": "ucmdx",
				"amount": "1115"
			}
		},
		{
			"address": "comdex1r4ylr8knuljvhlwpnztzu8pwtehmkppz7lgh02",
			"reward": {
				"denom": "ucmdx",
				"amount": "10648"
			}
		},
		{
			"address": "comdex1r4f6f35qc8fk6k0ydstd48pn00zhc4y0pc94cx",
			"reward": {
				"denom": "ucmdx",
				"amount": "8962"
			}
		},
		{
			"address": "comdex1r42cyntdmq0hqwl8kqe2rrk0rx53w2nrh2xzf8",
			"reward": {
				"denom": "ucmdx",
				"amount": "251"
			}
		},
		{
			"address": "comdex1r4thrg4c9crc2e2lppx07dhu04x26xdezvlp32",
			"reward": {
				"denom": "ucmdx",
				"amount": "146197"
			}
		},
		{
			"address": "comdex1r4d2ux7z3vyhw7968luml6kavvtaz0z3rkzwwd",
			"reward": {
				"denom": "ucmdx",
				"amount": "142"
			}
		},
		{
			"address": "comdex1r4jh8jspuaq6urq8eqkts7daun5entam7cauma",
			"reward": {
				"denom": "ucmdx",
				"amount": "359"
			}
		},
		{
			"address": "comdex1r4jecy60py35ljaef0esqmvsmsfu52235ltcf3",
			"reward": {
				"denom": "ucmdx",
				"amount": "4214"
			}
		},
		{
			"address": "comdex1r4nqllj5u4l5pacqfnjt2c0nhv8et5lt9la3p9",
			"reward": {
				"denom": "ucmdx",
				"amount": "166"
			}
		},
		{
			"address": "comdex1r4n5la7hefgeewu55awevjy4sdtvzurqa7ceql",
			"reward": {
				"denom": "ucmdx",
				"amount": "74"
			}
		},
		{
			"address": "comdex1r44wvdka4tf5ygj6unrhfwcls4us0jkp5mwmll",
			"reward": {
				"denom": "ucmdx",
				"amount": "21339"
			}
		},
		{
			"address": "comdex1r4kgdflhtshyl4lt8v3dsnh3detwpqreyujkf2",
			"reward": {
				"denom": "ucmdx",
				"amount": "51001"
			}
		},
		{
			"address": "comdex1r4hsgw2xr70s8lngttylcat670e4rrsxrjxkz6",
			"reward": {
				"denom": "ucmdx",
				"amount": "5931"
			}
		},
		{
			"address": "comdex1r4lrqxmyzv6lm7fqknrd9khna6w8e3xxp7flft",
			"reward": {
				"denom": "ucmdx",
				"amount": "1753"
			}
		},
		{
			"address": "comdex1r4lf6phm5950yq8f4g4yqjudvh6cst7es7gs8x",
			"reward": {
				"denom": "ucmdx",
				"amount": "10535"
			}
		},
		{
			"address": "comdex1rkpsf84h6wqlvwdpjpv557q858qkzrpfypdy5s",
			"reward": {
				"denom": "ucmdx",
				"amount": "431"
			}
		},
		{
			"address": "comdex1rkfwwydm4zwm436m4m86dqpkfps3qxpqvuzrrn",
			"reward": {
				"denom": "ucmdx",
				"amount": "13593"
			}
		},
		{
			"address": "comdex1rkdcav650kt9hejngqv8jajhwmug4ra8vx7a2k",
			"reward": {
				"denom": "ucmdx",
				"amount": "179"
			}
		},
		{
			"address": "comdex1rk0tvj7l5zdq8fdtpn2srw9fd40l8hxuyqyvf4",
			"reward": {
				"denom": "ucmdx",
				"amount": "6801"
			}
		},
		{
			"address": "comdex1rksdun4z9k54p94hepn3gden99728xwa974cjc",
			"reward": {
				"denom": "ucmdx",
				"amount": "9733"
			}
		},
		{
			"address": "comdex1rks34negumft9fa2u3q3d65w23py74qf95x3c6",
			"reward": {
				"denom": "ucmdx",
				"amount": "1955"
			}
		},
		{
			"address": "comdex1rknrapqtccrj8qflfaywajg2yj2dd5409c6qwz",
			"reward": {
				"denom": "ucmdx",
				"amount": "8969"
			}
		},
		{
			"address": "comdex1rkk8tfk8xkhln3v64776tq6e3wjs9jh9kpxzql",
			"reward": {
				"denom": "ucmdx",
				"amount": "1431"
			}
		},
		{
			"address": "comdex1rk6qdahk8s7e5pn7354ddfzu90qfcd8jkjx54n",
			"reward": {
				"denom": "ucmdx",
				"amount": "216"
			}
		},
		{
			"address": "comdex1rkas75auffvr7vjt57fl3637ryleamlrlam76m",
			"reward": {
				"denom": "ucmdx",
				"amount": "247"
			}
		},
		{
			"address": "comdex1rka7s5vs8uwagde3l95g0jfehkculh08namccu",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1rkl8s9x4fm2353ht5cxuf8lh6wqrhpydq8uhcc",
			"reward": {
				"denom": "ucmdx",
				"amount": "2251"
			}
		},
		{
			"address": "comdex1rhpa7rrd03dflt9q40qcluyhc6zwhuwpf86v8p",
			"reward": {
				"denom": "ucmdx",
				"amount": "9190"
			}
		},
		{
			"address": "comdex1rhzqlga6p5pv9lfc464usx7el5q375vz6f6g6s",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1rhy88a6zcc5rmhfnzvrpsm3cnc0u7a7m70d98k",
			"reward": {
				"denom": "ucmdx",
				"amount": "6433"
			}
		},
		{
			"address": "comdex1rhyakd4ctadxq2cmavgmpredq8t09xax9hzyd5",
			"reward": {
				"denom": "ucmdx",
				"amount": "140"
			}
		},
		{
			"address": "comdex1rh9kyfkq2t0e0z99e0asxfys39raxcq9mskzef",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1rhd7lywrqh6hg9qgmv3xqz3n6np42mj0squu3c",
			"reward": {
				"denom": "ucmdx",
				"amount": "849751"
			}
		},
		{
			"address": "comdex1rh0qta8wzcgnqg8eh7h5lqqnp367dje3qd2uuj",
			"reward": {
				"denom": "ucmdx",
				"amount": "167"
			}
		},
		{
			"address": "comdex1rhj9wp5u3frarxkwer67vufxwvnpfglcf64wz3",
			"reward": {
				"denom": "ucmdx",
				"amount": "15419"
			}
		},
		{
			"address": "comdex1rhnycf27euj0k6e9w8g7wtyy6uxvdgugat3f86",
			"reward": {
				"denom": "ucmdx",
				"amount": "178"
			}
		},
		{
			"address": "comdex1rhnuphhhpyzm3dwzwrnwwj79lf50yyt40a92kd",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1rh4dluvm59qe8dumxn6xrrx4qaxcswq6446m4k",
			"reward": {
				"denom": "ucmdx",
				"amount": "7213"
			}
		},
		{
			"address": "comdex1rhhk93luzchn54vvpr3uwyetymufq592a0uakv",
			"reward": {
				"denom": "ucmdx",
				"amount": "72"
			}
		},
		{
			"address": "comdex1rhe3eqns0ysehxpnj9xjtxc4wzrlu4l4746wqk",
			"reward": {
				"denom": "ucmdx",
				"amount": "1740"
			}
		},
		{
			"address": "comdex1rhe5kw3zk49h34tf24vg3pa6vef8uasu45nfr0",
			"reward": {
				"denom": "ucmdx",
				"amount": "199"
			}
		},
		{
			"address": "comdex1rh6ua96hvnh2muw6vlrz25nv2r0z5leacq3n7n",
			"reward": {
				"denom": "ucmdx",
				"amount": "546"
			}
		},
		{
			"address": "comdex1rh77t8g7hlajsxdeke2uveypv6v35dpqcgzhyz",
			"reward": {
				"denom": "ucmdx",
				"amount": "20520"
			}
		},
		{
			"address": "comdex1rhl6xexsnernd4zt3l7cuuwnp6zn33hr4dg2yl",
			"reward": {
				"denom": "ucmdx",
				"amount": "2121"
			}
		},
		{
			"address": "comdex1rczjkh9mjavf26x6u0uxnwd62rumlq3t5k8msw",
			"reward": {
				"denom": "ucmdx",
				"amount": "705"
			}
		},
		{
			"address": "comdex1rcy27tz8edk5udhtrqwymws5vymr5ne4ct3wxw",
			"reward": {
				"denom": "ucmdx",
				"amount": "27994"
			}
		},
		{
			"address": "comdex1rc8lmtpcq065gkmphz7m0248wwtxwm05tu5xaj",
			"reward": {
				"denom": "ucmdx",
				"amount": "2742"
			}
		},
		{
			"address": "comdex1rc2647d9nftgwwh9n5vaw26lnjdrewc4fqdtqc",
			"reward": {
				"denom": "ucmdx",
				"amount": "29178"
			}
		},
		{
			"address": "comdex1rc0m7vua7flfd3t3ar48p8w6glcz8k3dtelgy5",
			"reward": {
				"denom": "ucmdx",
				"amount": "89"
			}
		},
		{
			"address": "comdex1rcsw6lxwpxzeuv96dutv0j4eny6x83n095jwjp",
			"reward": {
				"denom": "ucmdx",
				"amount": "152"
			}
		},
		{
			"address": "comdex1rc3k3xpngt5lcrltkw4m7zt0s07nnjttqdgq49",
			"reward": {
				"denom": "ucmdx",
				"amount": "381"
			}
		},
		{
			"address": "comdex1rcjnhd5z0r6yly9r53e2tvwu8vpx7m5nrcqsk8",
			"reward": {
				"denom": "ucmdx",
				"amount": "6786"
			}
		},
		{
			"address": "comdex1rc5qffknrxu3dq7uwfcmqccqdmwzgrjkmnzwze",
			"reward": {
				"denom": "ucmdx",
				"amount": "3369"
			}
		},
		{
			"address": "comdex1rc5yr5w0h7kjrw5hyk7s340kn0fhr4fy02w6gt",
			"reward": {
				"denom": "ucmdx",
				"amount": "3720"
			}
		},
		{
			"address": "comdex1rc49a0ssmqktq456zkljw95qg0ghz3g68ucclp",
			"reward": {
				"denom": "ucmdx",
				"amount": "2859"
			}
		},
		{
			"address": "comdex1rck3y44nvdn4sefkxdkwsctv0eq0w6l2yevhrd",
			"reward": {
				"denom": "ucmdx",
				"amount": "157821"
			}
		},
		{
			"address": "comdex1rchnjkr78mz7p08a9lqdvy53l9sayfqsulyyul",
			"reward": {
				"denom": "ucmdx",
				"amount": "2596"
			}
		},
		{
			"address": "comdex1rc6zk2k0las34k4zn5msp4d48yuukg0j8c62hg",
			"reward": {
				"denom": "ucmdx",
				"amount": "174"
			}
		},
		{
			"address": "comdex1rc6d2fuh7e00wqmls7czngljwwwl3gd93lrcsp",
			"reward": {
				"denom": "ucmdx",
				"amount": "69"
			}
		},
		{
			"address": "comdex1rcuft35qpjzpezpg6ytcrf0nmvk3l96qpzrecq",
			"reward": {
				"denom": "ucmdx",
				"amount": "3794"
			}
		},
		{
			"address": "comdex1rclcwvu6fhzgxvde92vr5jpf4gdgtmd9jf484n",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1rerzet35eymqvgx2yc7yveq020w77tjc27jc7z",
			"reward": {
				"denom": "ucmdx",
				"amount": "2181"
			}
		},
		{
			"address": "comdex1rerw0v6u8d9pn9a3ruf5l680nw8m6365hjcz26",
			"reward": {
				"denom": "ucmdx",
				"amount": "8819"
			}
		},
		{
			"address": "comdex1rer3qx83un3x96zkhc2gy6285r2ysjshej659h",
			"reward": {
				"denom": "ucmdx",
				"amount": "86827"
			}
		},
		{
			"address": "comdex1re9qnpwdyylzfu6z7spt0xergk3gp6v4vaa473",
			"reward": {
				"denom": "ucmdx",
				"amount": "2254"
			}
		},
		{
			"address": "comdex1re8a4xf2xter8cuqffyexj7dgpt9342zm7vd2a",
			"reward": {
				"denom": "ucmdx",
				"amount": "37460"
			}
		},
		{
			"address": "comdex1retq3mz5fdj2w09adkzxqu48lvp44a7gqaajtd",
			"reward": {
				"denom": "ucmdx",
				"amount": "180"
			}
		},
		{
			"address": "comdex1retz0scvhmpar90jaewg5ln0pxemqtjnd7zu5d",
			"reward": {
				"denom": "ucmdx",
				"amount": "646"
			}
		},
		{
			"address": "comdex1rew6jgajftk24ztavd9mapr47m24e0c8sqsdqr",
			"reward": {
				"denom": "ucmdx",
				"amount": "2501"
			}
		},
		{
			"address": "comdex1re0guzzpw6uuz47a07hqx8ul0sdydyysufum8g",
			"reward": {
				"denom": "ucmdx",
				"amount": "3531"
			}
		},
		{
			"address": "comdex1re35gutd7ewvdckt6c8jualhfcyudlu24zv9n2",
			"reward": {
				"denom": "ucmdx",
				"amount": "713"
			}
		},
		{
			"address": "comdex1rejm7sec650wwgvqcdm7r3hnyprd0a57ykhtu6",
			"reward": {
				"denom": "ucmdx",
				"amount": "1887"
			}
		},
		{
			"address": "comdex1renjgkzppgk36dlrudnay2zlw77jgsczyxyj8p",
			"reward": {
				"denom": "ucmdx",
				"amount": "28323"
			}
		},
		{
			"address": "comdex1re4fyk6cpkpepukwhk6h34ynw5pye226pahfdp",
			"reward": {
				"denom": "ucmdx",
				"amount": "17924"
			}
		},
		{
			"address": "comdex1re4adwnyqjaxwvjaxu9rxlrw8elxtdqd42tzcr",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1reerg23sp9u79sz2emra5kgz00delnaudq2vxz",
			"reward": {
				"denom": "ucmdx",
				"amount": "176"
			}
		},
		{
			"address": "comdex1re6p223rwfefd2ptd547mv2aqwd0c4vump6cjr",
			"reward": {
				"denom": "ucmdx",
				"amount": "8714"
			}
		},
		{
			"address": "comdex1r6q46m952lema5hyl04rjpxmqrp6u2rkdkwxzr",
			"reward": {
				"denom": "ucmdx",
				"amount": "180"
			}
		},
		{
			"address": "comdex1r6z8mdgsmuuh7pv5u3utqrkglvmakeryf54yky",
			"reward": {
				"denom": "ucmdx",
				"amount": "171"
			}
		},
		{
			"address": "comdex1r6zdqatd5lz47p49gw0q0httlrzy56j6d76npd",
			"reward": {
				"denom": "ucmdx",
				"amount": "125"
			}
		},
		{
			"address": "comdex1r6yghsn9kptt0gl8sr0yevva2rluwtr7g9syz0",
			"reward": {
				"denom": "ucmdx",
				"amount": "4970"
			}
		},
		{
			"address": "comdex1r6ynldq2z3zgn8afr3cxqzdfwestcxsvprgzxv",
			"reward": {
				"denom": "ucmdx",
				"amount": "1142"
			}
		},
		{
			"address": "comdex1r69mmsgzvwc23fdtvm4l66qgknnkxh9zddvhhf",
			"reward": {
				"denom": "ucmdx",
				"amount": "1761"
			}
		},
		{
			"address": "comdex1r69unxe63naxxczsd3yus7rf0j4fasnzydvctt",
			"reward": {
				"denom": "ucmdx",
				"amount": "202"
			}
		},
		{
			"address": "comdex1r68ht4rszyttxt29p2fvlc9kx97056udvajxyc",
			"reward": {
				"denom": "ucmdx",
				"amount": "370474"
			}
		},
		{
			"address": "comdex1r627q55sy9gzmx54nggzyesmuu36x8694reu92",
			"reward": {
				"denom": "ucmdx",
				"amount": "1776"
			}
		},
		{
			"address": "comdex1r6tvjsw6zk4xhw05xjs36mptw65rdzer7u8z3k",
			"reward": {
				"denom": "ucmdx",
				"amount": "1027"
			}
		},
		{
			"address": "comdex1r6tjzt03hnafx4c6h499ph5h7p8wpwwl7vd5yc",
			"reward": {
				"denom": "ucmdx",
				"amount": "86"
			}
		},
		{
			"address": "comdex1r6vt4re2d6vu75hx2appcpn2u3anzjqj5ddqy9",
			"reward": {
				"denom": "ucmdx",
				"amount": "144"
			}
		},
		{
			"address": "comdex1r6ndkkkmq66sgfctjn79aak9v4n0yyqk3ucrlw",
			"reward": {
				"denom": "ucmdx",
				"amount": "14449"
			}
		},
		{
			"address": "comdex1r6n0ukjek8r9p6mm6dzpgckfw5mk0l7qvlu2ut",
			"reward": {
				"denom": "ucmdx",
				"amount": "527"
			}
		},
		{
			"address": "comdex1r6nj99d8k5tcjh3dakxyas6vxkwpnnqua9m6a3",
			"reward": {
				"denom": "ucmdx",
				"amount": "12922"
			}
		},
		{
			"address": "comdex1r6kftt936662dnezdunv6v8gmemw2ljw6l32uz",
			"reward": {
				"denom": "ucmdx",
				"amount": "1021"
			}
		},
		{
			"address": "comdex1r6e8ylmzaedku6zwsa786u5h2pmhhw0zcj9lza",
			"reward": {
				"denom": "ucmdx",
				"amount": "1021"
			}
		},
		{
			"address": "comdex1r66kj28pxgv9aqwafznz3f9cg8dfuqvk8055vz",
			"reward": {
				"denom": "ucmdx",
				"amount": "178"
			}
		},
		{
			"address": "comdex1r6m0nuclskpv90pdzygtjcvzp0s5ghm9x9t50y",
			"reward": {
				"denom": "ucmdx",
				"amount": "555"
			}
		},
		{
			"address": "comdex1rmqqxuv762wv39fhw7pgesssf0lx2lh3n78anm",
			"reward": {
				"denom": "ucmdx",
				"amount": "17904"
			}
		},
		{
			"address": "comdex1rmqwr69ua0yzwtjvzxwtuwycqm9vkqlpa75am7",
			"reward": {
				"denom": "ucmdx",
				"amount": "12737"
			}
		},
		{
			"address": "comdex1rmqmdt32ff2xs5cx9lfrcekan6jqdrvunzggux",
			"reward": {
				"denom": "ucmdx",
				"amount": "7274"
			}
		},
		{
			"address": "comdex1rmpmkqmzzjx4t0s7unj7tehvj7p5g9dfdsr0da",
			"reward": {
				"denom": "ucmdx",
				"amount": "88"
			}
		},
		{
			"address": "comdex1rmzptymwq859ts34zu92z55a5rpadm6w7d50v7",
			"reward": {
				"denom": "ucmdx",
				"amount": "19"
			}
		},
		{
			"address": "comdex1rmrdpty2n8mq67uqc8cm3uj5anaqulkn6hqs8r",
			"reward": {
				"denom": "ucmdx",
				"amount": "17089"
			}
		},
		{
			"address": "comdex1rmrkctju50d9vmsvsvn9dz2xjczqg689q88mga",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1rm9s2a3e4lze6qt5277su6hwzg5vtrlreeq4rz",
			"reward": {
				"denom": "ucmdx",
				"amount": "1765"
			}
		},
		{
			"address": "comdex1rmxl4fps24pe8s9uv8an3nqpng3ggyf8r6gn3f",
			"reward": {
				"denom": "ucmdx",
				"amount": "59215"
			}
		},
		{
			"address": "comdex1rm8v4efqwp3zawf67cuuln88mxgdzyzyf52znz",
			"reward": {
				"denom": "ucmdx",
				"amount": "888"
			}
		},
		{
			"address": "comdex1rm8s6cnafgahfafxrnekt43wr0fx5xh5u3nqvw",
			"reward": {
				"denom": "ucmdx",
				"amount": "717"
			}
		},
		{
			"address": "comdex1rmfmxffwag6qeqmeqa39jyd8c39ujczuj2elfj",
			"reward": {
				"denom": "ucmdx",
				"amount": "6296"
			}
		},
		{
			"address": "comdex1rmvvhwnyh744c52k87tfgkepfr06frpkhwsp4t",
			"reward": {
				"denom": "ucmdx",
				"amount": "28"
			}
		},
		{
			"address": "comdex1rmvn9m369snsmxflfv4ktrs6n40rsvs4280dyg",
			"reward": {
				"denom": "ucmdx",
				"amount": "395122"
			}
		},
		{
			"address": "comdex1rmdz40wfkyf6740yfk0z6sccxmw5n3pcevfyf4",
			"reward": {
				"denom": "ucmdx",
				"amount": "74223"
			}
		},
		{
			"address": "comdex1rm3zh79z7hr62cunfdwcpd5w8gwd8xxszu0dgg",
			"reward": {
				"denom": "ucmdx",
				"amount": "1229"
			}
		},
		{
			"address": "comdex1rmnf4t0jy2xxwl5hj4mxqhqkdy56ujj8av2hev",
			"reward": {
				"denom": "ucmdx",
				"amount": "534"
			}
		},
		{
			"address": "comdex1rm5qxlr352yhz3v4pj5urehtlgs7frnl95lkvx",
			"reward": {
				"denom": "ucmdx",
				"amount": "1772"
			}
		},
		{
			"address": "comdex1rmkys62qut4sryg3x4mszw9pcv9c2wn3mr2vy8",
			"reward": {
				"denom": "ucmdx",
				"amount": "143"
			}
		},
		{
			"address": "comdex1rmkyu39fgt9cppydwycc5p4kgmc0flrvc7kvhs",
			"reward": {
				"denom": "ucmdx",
				"amount": "149"
			}
		},
		{
			"address": "comdex1rmk972wd09fjx5r2pcy6assujskykq988mtfyl",
			"reward": {
				"denom": "ucmdx",
				"amount": "1473"
			}
		},
		{
			"address": "comdex1rmemz0vyxjzqgy4hljymvzq3e3uf4630j757uu",
			"reward": {
				"denom": "ucmdx",
				"amount": "17775"
			}
		},
		{
			"address": "comdex1rmmym5wfl8xy82dhe3ghh2k5qzks638v6w0zn9",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1rmajk4pdv6y40kq0evck05gksuru6mtnw30as3",
			"reward": {
				"denom": "ucmdx",
				"amount": "14096"
			}
		},
		{
			"address": "comdex1rm7299qn4pm2zh5m8lrh5rjms3uu230enu9as8",
			"reward": {
				"denom": "ucmdx",
				"amount": "176"
			}
		},
		{
			"address": "comdex1rml8k35002dpcq4kenwk4h4jfhukyqytfu42ca",
			"reward": {
				"denom": "ucmdx",
				"amount": "10046"
			}
		},
		{
			"address": "comdex1ruphl50638rpc93vpr6nzaaw76ludcsl9clc8a",
			"reward": {
				"denom": "ucmdx",
				"amount": "4452"
			}
		},
		{
			"address": "comdex1ru9346jkcy762neudmnvuswpnl8jmktejefsup",
			"reward": {
				"denom": "ucmdx",
				"amount": "1409"
			}
		},
		{
			"address": "comdex1ru9jxwmxlh9zem6y4xhscklqx9fvv9zm7sslky",
			"reward": {
				"denom": "ucmdx",
				"amount": "1900"
			}
		},
		{
			"address": "comdex1rux5j9tfsl92zxvqv04m4ntwlxezsl6a6ncgxg",
			"reward": {
				"denom": "ucmdx",
				"amount": "4225"
			}
		},
		{
			"address": "comdex1ru8n0ryexq33uerw28wk9s9vrjwjtyazrwzgnp",
			"reward": {
				"denom": "ucmdx",
				"amount": "1230"
			}
		},
		{
			"address": "comdex1rudrkgsmyrm8uxez7m4xsfklq0nmzmfplaclax",
			"reward": {
				"denom": "ucmdx",
				"amount": "7828"
			}
		},
		{
			"address": "comdex1ruwraamthss0nsvgye7pwnl6jarxwk0vq4tnvy",
			"reward": {
				"denom": "ucmdx",
				"amount": "1461"
			}
		},
		{
			"address": "comdex1ruw9jyp5yax2ztkt05cx20x6zre3kp6cpcka9t",
			"reward": {
				"denom": "ucmdx",
				"amount": "619"
			}
		},
		{
			"address": "comdex1ru0nx4cm8r3g9dyw2kad7d3cz3uym2yuxp06dd",
			"reward": {
				"denom": "ucmdx",
				"amount": "1410"
			}
		},
		{
			"address": "comdex1rusqjyx3k2j680rnm2kxae8zvka5s4t9lp0egw",
			"reward": {
				"denom": "ucmdx",
				"amount": "1671"
			}
		},
		{
			"address": "comdex1ruslszl7z5q2lrnmkknvznpw3hrguduttdjwff",
			"reward": {
				"denom": "ucmdx",
				"amount": "224"
			}
		},
		{
			"address": "comdex1ru3pnm29lpu8hjlcsw3hp0eywc4nd4fp5f99xc",
			"reward": {
				"denom": "ucmdx",
				"amount": "1848"
			}
		},
		{
			"address": "comdex1ru3tdzdnfww7du3pkfxfxdxesvwclu7rtqmqt4",
			"reward": {
				"denom": "ucmdx",
				"amount": "522"
			}
		},
		{
			"address": "comdex1run8se9letzfxxxpp2mktwdukrm9nm27kn2rek",
			"reward": {
				"denom": "ucmdx",
				"amount": "6203"
			}
		},
		{
			"address": "comdex1rukgach9v7znq7dajy09vwxy9qf9ql7ysf6dyl",
			"reward": {
				"denom": "ucmdx",
				"amount": "71720"
			}
		},
		{
			"address": "comdex1rue3krma7x5tfedfnmdm9d032wxmnd5t602uaf",
			"reward": {
				"denom": "ucmdx",
				"amount": "51438"
			}
		},
		{
			"address": "comdex1ru6246qn5t2jfdvz82m0m4axn9w4duu53g7lvz",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1ru62e377a5yvw2vus6tx6u4sclqrs3yye4q7pk",
			"reward": {
				"denom": "ucmdx",
				"amount": "1040"
			}
		},
		{
			"address": "comdex1ru6j24uz2hasdww5ur6v8fn9k9knrhev30tycd",
			"reward": {
				"denom": "ucmdx",
				"amount": "2014"
			}
		},
		{
			"address": "comdex1ruuqqmf734xx7qf56xdksu3ulrgfx492mpquys",
			"reward": {
				"denom": "ucmdx",
				"amount": "6186"
			}
		},
		{
			"address": "comdex1ruup2az5gwavvrqcn82kg87ke8j6s3077304l8",
			"reward": {
				"denom": "ucmdx",
				"amount": "136"
			}
		},
		{
			"address": "comdex1ruuh8t8fy0hs98pnas7ercpmxmn7ecyg3sq2xh",
			"reward": {
				"denom": "ucmdx",
				"amount": "1902"
			}
		},
		{
			"address": "comdex1ru7fupnqmfwhqhgh3qgq49whrl9n0z9ufpsm50",
			"reward": {
				"denom": "ucmdx",
				"amount": "149"
			}
		},
		{
			"address": "comdex1rulmmg0ee9ffm9q3rxq59eakeke6nts64whaxc",
			"reward": {
				"denom": "ucmdx",
				"amount": "75"
			}
		},
		{
			"address": "comdex1rarlv7e5c7px2f7lfljq58uwzwmwcxzul7cwt8",
			"reward": {
				"denom": "ucmdx",
				"amount": "34"
			}
		},
		{
			"address": "comdex1rayxf0asjq3yrdg0kxz96prp7l4gk4ncgnphqj",
			"reward": {
				"denom": "ucmdx",
				"amount": "19016"
			}
		},
		{
			"address": "comdex1ra9563z5tn2lmqhydt2atrgzftk2d7umx3cckj",
			"reward": {
				"denom": "ucmdx",
				"amount": "8791"
			}
		},
		{
			"address": "comdex1rax350s6dgmjh4373hx2vkmq0ekmzegguwtpyg",
			"reward": {
				"denom": "ucmdx",
				"amount": "750"
			}
		},
		{
			"address": "comdex1raxlzjunlwq335v3jgktmf0sfu80ez8vm7z4z2",
			"reward": {
				"denom": "ucmdx",
				"amount": "473"
			}
		},
		{
			"address": "comdex1rag8rjr6d66qpdzxqzls6ks89slpjx6yaye44d",
			"reward": {
				"denom": "ucmdx",
				"amount": "222"
			}
		},
		{
			"address": "comdex1rafqlhkvz38d9jhzkxumr78nfg8gyce4a8k4y9",
			"reward": {
				"denom": "ucmdx",
				"amount": "106435"
			}
		},
		{
			"address": "comdex1ravrrr3un9hk4zfhu3cfayc9zqn3z96cujeq8m",
			"reward": {
				"denom": "ucmdx",
				"amount": "1763"
			}
		},
		{
			"address": "comdex1ra3aelqq4u5czfr9kcpww0gfl84nmqjxl63l7q",
			"reward": {
				"denom": "ucmdx",
				"amount": "5078"
			}
		},
		{
			"address": "comdex1ra5u3yj5aaw5pas88p8ekgl2g3kczx62vhaau3",
			"reward": {
				"denom": "ucmdx",
				"amount": "4044"
			}
		},
		{
			"address": "comdex1racq9362dw7urvfung32xrtsujdygj0augt0vq",
			"reward": {
				"denom": "ucmdx",
				"amount": "154454"
			}
		},
		{
			"address": "comdex1ra6mrwfuhka8a2xed0w2969grw30wtgradr0hu",
			"reward": {
				"denom": "ucmdx",
				"amount": "155"
			}
		},
		{
			"address": "comdex1rau40jn64r0yrd9xh9jjx6zl5gqh7h7u3cr5ss",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1ra7kcta5k5nx97xdy7y2mmyfejyvw0lluycmpl",
			"reward": {
				"denom": "ucmdx",
				"amount": "1726"
			}
		},
		{
			"address": "comdex1ralz0egzm87anu44wjqjg6723k3wc654cphpg0",
			"reward": {
				"denom": "ucmdx",
				"amount": "993"
			}
		},
		{
			"address": "comdex1r7q0m4dqxywkk7xsxe0j8ul56cqxvapxvw3rpk",
			"reward": {
				"denom": "ucmdx",
				"amount": "5342"
			}
		},
		{
			"address": "comdex1r7pnv302qhmn0st3kx78jk8e7hrq2870gxvxp6",
			"reward": {
				"denom": "ucmdx",
				"amount": "1948"
			}
		},
		{
			"address": "comdex1r7zdu8e8j2skj3ff330clfu2qwecp3vd6fvfmc",
			"reward": {
				"denom": "ucmdx",
				"amount": "1919"
			}
		},
		{
			"address": "comdex1r7z7j90hqanx7yv67nafru3d6au5e37mvjwn8z",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1r7r0gvavvtmmnplk8h263hj5xk44x2kaqgj7g3",
			"reward": {
				"denom": "ucmdx",
				"amount": "12380"
			}
		},
		{
			"address": "comdex1r78h0c8zmd3lwnhamwwdf3ahe2czvz0ud58aqn",
			"reward": {
				"denom": "ucmdx",
				"amount": "1184"
			}
		},
		{
			"address": "comdex1r7gt4ep5ccaqckrz45fuuwy74a8zvtfm384w32",
			"reward": {
				"denom": "ucmdx",
				"amount": "1190"
			}
		},
		{
			"address": "comdex1r7fteg3ugtwh00vglgrc85r7rq09s8vdfy9v8q",
			"reward": {
				"denom": "ucmdx",
				"amount": "150"
			}
		},
		{
			"address": "comdex1r7v49skmt25nuz2ny0xp75et5q593axk574qfr",
			"reward": {
				"denom": "ucmdx",
				"amount": "71"
			}
		},
		{
			"address": "comdex1r7d6mkjpwracqdv90rm2kdgewmrr8xwuzpj66p",
			"reward": {
				"denom": "ucmdx",
				"amount": "174"
			}
		},
		{
			"address": "comdex1r70gaxncn2nwhqd70ee3eczdteppy8f93avsc0",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1r7nreqpf9hdxdyw3jmdd080lxkzdy0ugyr282c",
			"reward": {
				"denom": "ucmdx",
				"amount": "71"
			}
		},
		{
			"address": "comdex1r75p3ss6cdhvvna2j9g4qkk80tst2zj4n0dgde",
			"reward": {
				"denom": "ucmdx",
				"amount": "90387"
			}
		},
		{
			"address": "comdex1r74s3zzye0ukwt04eveump6kuftm6c7f5vkzvp",
			"reward": {
				"denom": "ucmdx",
				"amount": "32761"
			}
		},
		{
			"address": "comdex1r7krypf66v0yxddhazd28gv05x0zfhl0g4nn0k",
			"reward": {
				"denom": "ucmdx",
				"amount": "559296"
			}
		},
		{
			"address": "comdex1r7k96aufpzznw0mrq9xg3sx7khjhc2azskem6c",
			"reward": {
				"denom": "ucmdx",
				"amount": "2698"
			}
		},
		{
			"address": "comdex1r7h3ak485snphcjkk8pt4yc06tagqgjzpzfhyr",
			"reward": {
				"denom": "ucmdx",
				"amount": "16"
			}
		},
		{
			"address": "comdex1r7cejjd3z6sgmg878447grjgye60qrjkvwgncw",
			"reward": {
				"denom": "ucmdx",
				"amount": "30586"
			}
		},
		{
			"address": "comdex1r76cjply5304v5wzlfsj4fx7c50lx2v7cuz0my",
			"reward": {
				"denom": "ucmdx",
				"amount": "39"
			}
		},
		{
			"address": "comdex1r7mul5yxz4e44dlxm0s5r894g3y2tkkq8e2cxe",
			"reward": {
				"denom": "ucmdx",
				"amount": "12863"
			}
		},
		{
			"address": "comdex1r7udhmde2u2r7t9x4e2k8waj7nkdhxdms7uvh8",
			"reward": {
				"denom": "ucmdx",
				"amount": "144"
			}
		},
		{
			"address": "comdex1r7lah5y3qyq453hqfgmxl2y0s5wesh09nqj8ms",
			"reward": {
				"denom": "ucmdx",
				"amount": "1053"
			}
		},
		{
			"address": "comdex1rl8rxfee8ugumza9yzgssarsj6qkq0reewvz6l",
			"reward": {
				"denom": "ucmdx",
				"amount": "24288"
			}
		},
		{
			"address": "comdex1rlfw3suplyawg24k85dgmnxja43ernpkmt54n0",
			"reward": {
				"denom": "ucmdx",
				"amount": "1435"
			}
		},
		{
			"address": "comdex1rlfezq7fedawqp6w7p0vmyq42a9kcpma869jem",
			"reward": {
				"denom": "ucmdx",
				"amount": "14187"
			}
		},
		{
			"address": "comdex1rl2qh2ecr2rkjes76q6y6h3zkmynmw0zx8s3gx",
			"reward": {
				"denom": "ucmdx",
				"amount": "346"
			}
		},
		{
			"address": "comdex1rl23we4zmztwv9zvg8vvpnmv02cp8e7y4a92af",
			"reward": {
				"denom": "ucmdx",
				"amount": "1255"
			}
		},
		{
			"address": "comdex1rl03vh459c56v6rw72kfdac64desxnk783m75s",
			"reward": {
				"denom": "ucmdx",
				"amount": "1419"
			}
		},
		{
			"address": "comdex1rlsdrt7m2aat7cvjvks00n5srq5khllqcjndl2",
			"reward": {
				"denom": "ucmdx",
				"amount": "1437"
			}
		},
		{
			"address": "comdex1rls6gupt0vt42s3c0ffr9dagh2vlwprlrw5whv",
			"reward": {
				"denom": "ucmdx",
				"amount": "871512"
			}
		},
		{
			"address": "comdex1rl39c8hsd55cd7c999e9mk36lua4u6np93vuq7",
			"reward": {
				"denom": "ucmdx",
				"amount": "287"
			}
		},
		{
			"address": "comdex1rl3dhfqretkzwktn4fnazyssdglrul4ah9c2nd",
			"reward": {
				"denom": "ucmdx",
				"amount": "1504"
			}
		},
		{
			"address": "comdex1rlj6j087x9q989w6ljt4jkcx5mv4l8mtew9w24",
			"reward": {
				"denom": "ucmdx",
				"amount": "14405"
			}
		},
		{
			"address": "comdex1rl5wj8ldwn2n95ancxkhhqfsmmtw8ew3n638sw",
			"reward": {
				"denom": "ucmdx",
				"amount": "168"
			}
		},
		{
			"address": "comdex1rl4tuq777hxt0fec46x8dy68utzs0qr5932jp8",
			"reward": {
				"denom": "ucmdx",
				"amount": "25656"
			}
		},
		{
			"address": "comdex1rl4exzwq4dvgcwyljtc7qqdwxnyvy9j3g9fazw",
			"reward": {
				"denom": "ucmdx",
				"amount": "2000"
			}
		},
		{
			"address": "comdex1rlkdml7dpuedpmkxhval5qlqzp9hxfpukjvgfw",
			"reward": {
				"denom": "ucmdx",
				"amount": "1026"
			}
		},
		{
			"address": "comdex1rlk4p5hjsa0ve8ncjftqgdgjr2vrfm8690xhtk",
			"reward": {
				"denom": "ucmdx",
				"amount": "8625"
			}
		},
		{
			"address": "comdex1rl65kece87y2plffkl756rdny0r3mp23jjc577",
			"reward": {
				"denom": "ucmdx",
				"amount": "696"
			}
		},
		{
			"address": "comdex1rlmyshum4r996hgewmj57j33pll4788ujuj4t4",
			"reward": {
				"denom": "ucmdx",
				"amount": "151"
			}
		},
		{
			"address": "comdex1rlm8hhq7xcs50ln06skyl2er2dnkdw2d9tgf55",
			"reward": {
				"denom": "ucmdx",
				"amount": "889"
			}
		},
		{
			"address": "comdex1rlmuqh8f96yzdetdnctj4r6prktpt77nwjvg4p",
			"reward": {
				"denom": "ucmdx",
				"amount": "70"
			}
		},
		{
			"address": "comdex1rlu09uxyqz4cf4ce8j7qgjaruyp357yu3gtjz3",
			"reward": {
				"denom": "ucmdx",
				"amount": "10064"
			}
		},
		{
			"address": "comdex1rlavrcmh7c54kusfydu56h39gksxa2v3vlkeqp",
			"reward": {
				"denom": "ucmdx",
				"amount": "895"
			}
		},
		{
			"address": "comdex1rla36l65dkzcvetv4scaf9tsqh02522hracgvm",
			"reward": {
				"denom": "ucmdx",
				"amount": "2444"
			}
		},
		{
			"address": "comdex1rll0jfn2u6waz2vujq4a5rc2rtdct6l8qh6w2q",
			"reward": {
				"denom": "ucmdx",
				"amount": "1938"
			}
		},
		{
			"address": "comdex1rlla0qxsy72vpp6870t8kvergzeph6lgva920v",
			"reward": {
				"denom": "ucmdx",
				"amount": "312"
			}
		},
		{
			"address": "comdex1yqqdy6ez96atd5yqtcvpd5azdnf0n069j7z55f",
			"reward": {
				"denom": "ucmdx",
				"amount": "98"
			}
		},
		{
			"address": "comdex1yqp32nyaks7qna8uz4h944vq4nj58ce7n8efy8",
			"reward": {
				"denom": "ucmdx",
				"amount": "73"
			}
		},
		{
			"address": "comdex1yqx7eqt4hyayd9pp5rpju6uzlv2xwqada7xs3x",
			"reward": {
				"denom": "ucmdx",
				"amount": "9755"
			}
		},
		{
			"address": "comdex1yq8g7javg0n4uczlcw64nmxwf4kc4cyes77gts",
			"reward": {
				"denom": "ucmdx",
				"amount": "22956"
			}
		},
		{
			"address": "comdex1yq8tl54yydesau7fkc24ywfkc05ldnekmjklwf",
			"reward": {
				"denom": "ucmdx",
				"amount": "339"
			}
		},
		{
			"address": "comdex1yq2e5mva2sxawyajzdpdh2ahg7mykwqycvfr6f",
			"reward": {
				"denom": "ucmdx",
				"amount": "1066"
			}
		},
		{
			"address": "comdex1yq2aen5x6u0hu3pedt94lhxunsznhyslxvk8st",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1yqth74nsx8pgmhccssqrwe7luvjmz8wqvlnyx4",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1yqvxyhzqd4e6mg6hhd0x08jr5yg6mdzje39j2a",
			"reward": {
				"denom": "ucmdx",
				"amount": "8734"
			}
		},
		{
			"address": "comdex1yq3xs9u634ew8gjfxw8r6sg8vvxkcmshgdk82l",
			"reward": {
				"denom": "ucmdx",
				"amount": "180"
			}
		},
		{
			"address": "comdex1yqj6yfxle223jljkmhwzcfqy5gkmsprvjms9w3",
			"reward": {
				"denom": "ucmdx",
				"amount": "5856"
			}
		},
		{
			"address": "comdex1yq5n2f8qeratvj8jhwxz26mxsyf7r83tu2ga6x",
			"reward": {
				"denom": "ucmdx",
				"amount": "976"
			}
		},
		{
			"address": "comdex1yq4jxk3k5uenj0nlvwgrm9jqe004sv50lhqflz",
			"reward": {
				"denom": "ucmdx",
				"amount": "3517"
			}
		},
		{
			"address": "comdex1yq4u4u3mjdtfsfw8au92wzs5v0ra7j20tnlpt3",
			"reward": {
				"denom": "ucmdx",
				"amount": "579"
			}
		},
		{
			"address": "comdex1yq6lvcn803zv8pq5qcrsrw6wxpg7t9hevup8th",
			"reward": {
				"denom": "ucmdx",
				"amount": "4167"
			}
		},
		{
			"address": "comdex1yqudky3w28w4gsytctjjr0ex8rhrhqv2ssdpr7",
			"reward": {
				"denom": "ucmdx",
				"amount": "990"
			}
		},
		{
			"address": "comdex1yqlelm3tunsscg3nq3v3lgkw4j9d52e8we247g",
			"reward": {
				"denom": "ucmdx",
				"amount": "2614"
			}
		},
		{
			"address": "comdex1ypq22n8y8aw7jtxux9kx860mnpp8phkvd08wjx",
			"reward": {
				"denom": "ucmdx",
				"amount": "11733"
			}
		},
		{
			"address": "comdex1yp923424vva3eeprfvmjzztry36dcc30mf0wzl",
			"reward": {
				"denom": "ucmdx",
				"amount": "892"
			}
		},
		{
			"address": "comdex1yp95ns7exf4l9jgh4rm58lmk3s6j80zyffzwza",
			"reward": {
				"denom": "ucmdx",
				"amount": "4941"
			}
		},
		{
			"address": "comdex1yp8xntz69060tdxw3df8ndgs03yqh0cxwset40",
			"reward": {
				"denom": "ucmdx",
				"amount": "6732"
			}
		},
		{
			"address": "comdex1ypt5gvkd6gy2lupnwzhems5h6mh42ztzvw8qny",
			"reward": {
				"denom": "ucmdx",
				"amount": "90"
			}
		},
		{
			"address": "comdex1ypd73l2zgdymx9h039lrjhkyccvkx8spmhp2h9",
			"reward": {
				"denom": "ucmdx",
				"amount": "608"
			}
		},
		{
			"address": "comdex1yp06lcmpgf4p0zqfu902hxs5gq6panez82jshr",
			"reward": {
				"denom": "ucmdx",
				"amount": "77"
			}
		},
		{
			"address": "comdex1ypjpuwmclph6m587aqld5zu0t0yuxe523whrjq",
			"reward": {
				"denom": "ucmdx",
				"amount": "1732"
			}
		},
		{
			"address": "comdex1ypj8tej66ccncvxzf9mlp29j0crcx3ljcsfnrt",
			"reward": {
				"denom": "ucmdx",
				"amount": "48029"
			}
		},
		{
			"address": "comdex1ypj5f9xprf2kadtgu6mhyuzu6pzrk77vzv4r8m",
			"reward": {
				"denom": "ucmdx",
				"amount": "89"
			}
		},
		{
			"address": "comdex1ypjcqq02nfvxgckdmafh9z5deg9lr89wneuwyp",
			"reward": {
				"denom": "ucmdx",
				"amount": "781"
			}
		},
		{
			"address": "comdex1ypn2m0qsvqqd3aacptt90y8tj442ktkmdt73hp",
			"reward": {
				"denom": "ucmdx",
				"amount": "151"
			}
		},
		{
			"address": "comdex1yphk4jddl06zc93t36klpukzu30y5tsudeq763",
			"reward": {
				"denom": "ucmdx",
				"amount": "90"
			}
		},
		{
			"address": "comdex1yp6xqp5avdhgah9nk6vaam9mw7yflyv4zm4zrz",
			"reward": {
				"denom": "ucmdx",
				"amount": "1505"
			}
		},
		{
			"address": "comdex1ypakuerp2fz5resn74cngarknzlgj87sp7r4vp",
			"reward": {
				"denom": "ucmdx",
				"amount": "32629"
			}
		},
		{
			"address": "comdex1ypl0y7devmgjw86gg8jn36pa3ez8rjm5278upu",
			"reward": {
				"denom": "ucmdx",
				"amount": "15761"
			}
		},
		{
			"address": "comdex1yzqcgl0tarldzhvq9s54lhrukafzasa3kal4ry",
			"reward": {
				"denom": "ucmdx",
				"amount": "74583"
			}
		},
		{
			"address": "comdex1yzrct893pnwpxqetkgchpmtvyfz3dp02wnncw4",
			"reward": {
				"denom": "ucmdx",
				"amount": "3658"
			}
		},
		{
			"address": "comdex1yzxs365xllkcv4rnvgz9g67t7ac8h4vrgq76pn",
			"reward": {
				"denom": "ucmdx",
				"amount": "1190"
			}
		},
		{
			"address": "comdex1yzgud8pe4fvgc9jy0l92pf7jg74ljf0u28gzgc",
			"reward": {
				"denom": "ucmdx",
				"amount": "271"
			}
		},
		{
			"address": "comdex1yzf45ckty87kf63d5627rqwdcff7jv34ey4czl",
			"reward": {
				"denom": "ucmdx",
				"amount": "169"
			}
		},
		{
			"address": "comdex1yz239ef9axpuprqsfnw4puuyflr2yqzzg0zjkj",
			"reward": {
				"denom": "ucmdx",
				"amount": "8933"
			}
		},
		{
			"address": "comdex1yzt8gy6pttcszpdfemweh59cnhdxpgwgzd00yd",
			"reward": {
				"denom": "ucmdx",
				"amount": "184"
			}
		},
		{
			"address": "comdex1yzvfg3f0mmfmsj3vjrg8neevq2gqcpzy3xhunk",
			"reward": {
				"denom": "ucmdx",
				"amount": "4551"
			}
		},
		{
			"address": "comdex1yz0ky304um9hzezz7xfu9nwm8jju46h06sznuf",
			"reward": {
				"denom": "ucmdx",
				"amount": "195"
			}
		},
		{
			"address": "comdex1yzs6c7d4ewpw5yskrrrwusccdtgq5d4atj0gz4",
			"reward": {
				"denom": "ucmdx",
				"amount": "1410"
			}
		},
		{
			"address": "comdex1yz5k7edv0m284m673wynr80tsf67kdjjhxzfqn",
			"reward": {
				"denom": "ucmdx",
				"amount": "112"
			}
		},
		{
			"address": "comdex1yz4klhg8ecfz2hugwtc3glvkvjh8vs06vpaeca",
			"reward": {
				"denom": "ucmdx",
				"amount": "14241"
			}
		},
		{
			"address": "comdex1yzc8mx67ngdfvxch3gjevd6t3gvwl0rnw0l0c8",
			"reward": {
				"denom": "ucmdx",
				"amount": "3536"
			}
		},
		{
			"address": "comdex1yzc2ucyztjyg7aveje0yjx805zwgrcy89j8xdf",
			"reward": {
				"denom": "ucmdx",
				"amount": "4383"
			}
		},
		{
			"address": "comdex1yzcl4gs7pdlpreyzzd59jtqnrd4u795dvu85gq",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1yzuk73u9rx2ef8xfnx288ljyxm02yfwwdz8rjg",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1yzumu9aaxcrqssq6lgz05ak73znx6yum6t9hfv",
			"reward": {
				"denom": "ucmdx",
				"amount": "1020"
			}
		},
		{
			"address": "comdex1yza3c69h46efwmdwa5wh0kafemexpv3zsedz79",
			"reward": {
				"denom": "ucmdx",
				"amount": "523"
			}
		},
		{
			"address": "comdex1yza508rvq3kgvx6zyrlwua9a7me40ke3q9hwcu",
			"reward": {
				"denom": "ucmdx",
				"amount": "1957"
			}
		},
		{
			"address": "comdex1yz7xedge0xg7cej8xzfq9jzjrev3d2d85jd4y0",
			"reward": {
				"denom": "ucmdx",
				"amount": "12321"
			}
		},
		{
			"address": "comdex1yz7cs0gfp0wzjcveulvcd6gfuqxd8kmlawkk6r",
			"reward": {
				"denom": "ucmdx",
				"amount": "194"
			}
		},
		{
			"address": "comdex1yrq4zzqjlq5t2fxnhgazmazy4sufe2j9846fpy",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1yrpy9s2ftwe2p86wzsleu0hdxvk8u4g7054k2f",
			"reward": {
				"denom": "ucmdx",
				"amount": "3652"
			}
		},
		{
			"address": "comdex1yrr3rmg7e0xxztjcvxs8c6kx4f5vw4m8423ynk",
			"reward": {
				"denom": "ucmdx",
				"amount": "1230"
			}
		},
		{
			"address": "comdex1yr9ydgz2yzwwyyf78qx9fralq6t8qytgh9ue99",
			"reward": {
				"denom": "ucmdx",
				"amount": "15597"
			}
		},
		{
			"address": "comdex1yr8fgts6d76g0u847zkng2e9l9nk4stwcxyjpe",
			"reward": {
				"denom": "ucmdx",
				"amount": "57359"
			}
		},
		{
			"address": "comdex1yrd730g3znkpwdgm65lvydmpta04usj7ue0dfk",
			"reward": {
				"denom": "ucmdx",
				"amount": "1266"
			}
		},
		{
			"address": "comdex1yrwc2e5q20asr73ndmdvxu33xnae25tkauc8yk",
			"reward": {
				"denom": "ucmdx",
				"amount": "16"
			}
		},
		{
			"address": "comdex1yr0pk8jtxpg2cu08qcgjwt36lm0dmr4hly7kl4",
			"reward": {
				"denom": "ucmdx",
				"amount": "72401"
			}
		},
		{
			"address": "comdex1yrs6s0tnnqr0n9926q6zhhdzhnmd4h0c2a9fha",
			"reward": {
				"denom": "ucmdx",
				"amount": "2005"
			}
		},
		{
			"address": "comdex1yr3ql5g49ffqz4pme5949wv8kw5zetnxtfpfu0",
			"reward": {
				"denom": "ucmdx",
				"amount": "1126"
			}
		},
		{
			"address": "comdex1yrjq3ya5njp6n0sn4kt3urgwhlduxe5jvvvu72",
			"reward": {
				"denom": "ucmdx",
				"amount": "1759"
			}
		},
		{
			"address": "comdex1yrjlqyf9r3kaxd3zea3ma02vk6meqgcjshsg0u",
			"reward": {
				"denom": "ucmdx",
				"amount": "6031"
			}
		},
		{
			"address": "comdex1yrnxx2j674wgc4jz4z8rvxllqrn4m6xaf59da2",
			"reward": {
				"denom": "ucmdx",
				"amount": "148"
			}
		},
		{
			"address": "comdex1yr56rjg87txr44cepgkj8kwwnjukmkqyamj9qf",
			"reward": {
				"denom": "ucmdx",
				"amount": "284"
			}
		},
		{
			"address": "comdex1yr4q5thgngxhan2prwpqnly7lde5ggz4tcxq3l",
			"reward": {
				"denom": "ucmdx",
				"amount": "1190"
			}
		},
		{
			"address": "comdex1yreqngeudeftdgpaasfkz6yp75kje7rlqvmxpr",
			"reward": {
				"denom": "ucmdx",
				"amount": "966"
			}
		},
		{
			"address": "comdex1yreychdns4z978d9fq8zmyl5ftlfm2and6a59m",
			"reward": {
				"denom": "ucmdx",
				"amount": "215"
			}
		},
		{
			"address": "comdex1yrmrz5xq8whqa7rl4pt4x8gs2ptv76nle88vnm",
			"reward": {
				"denom": "ucmdx",
				"amount": "13952"
			}
		},
		{
			"address": "comdex1yrmn29rdft22kds8yysqc776yd43y25un5nr9q",
			"reward": {
				"denom": "ucmdx",
				"amount": "179"
			}
		},
		{
			"address": "comdex1yrawdeqtp37njrlzwuexkp3qe5w4n40xteyx39",
			"reward": {
				"denom": "ucmdx",
				"amount": "112"
			}
		},
		{
			"address": "comdex1yyqupa3zf4cehydh597cwj6v7a9rm50hkxrwzh",
			"reward": {
				"denom": "ucmdx",
				"amount": "700"
			}
		},
		{
			"address": "comdex1yyzn63eadu9kryag052d2dcsuut8gmt5ac3q4n",
			"reward": {
				"denom": "ucmdx",
				"amount": "1752"
			}
		},
		{
			"address": "comdex1yy9vw8pskk6586y57tanfpjmx546ea40nzv666",
			"reward": {
				"denom": "ucmdx",
				"amount": "4659"
			}
		},
		{
			"address": "comdex1yy80gc85ypfe9uxr8mykfu73v0z6jcywxtuq82",
			"reward": {
				"denom": "ucmdx",
				"amount": "447"
			}
		},
		{
			"address": "comdex1yygx69d3x7jy3qjgwwe8zr0j2xhdarcp9zl40x",
			"reward": {
				"denom": "ucmdx",
				"amount": "83"
			}
		},
		{
			"address": "comdex1yy22ldc8h495hmlk258tkeslry6xlgpukgcrn0",
			"reward": {
				"denom": "ucmdx",
				"amount": "8333"
			}
		},
		{
			"address": "comdex1yyv5lj9s0yt9ar6d7mmr9jrtnz4sz95yutjn2t",
			"reward": {
				"denom": "ucmdx",
				"amount": "7191"
			}
		},
		{
			"address": "comdex1yyjcck3rw7qfhqf7ayvqdvfnvj3n34jxwjmanu",
			"reward": {
				"denom": "ucmdx",
				"amount": "4380"
			}
		},
		{
			"address": "comdex1yy5zyeuqfvld25uf3lvyyqet99rahmrvjxveld",
			"reward": {
				"denom": "ucmdx",
				"amount": "4711"
			}
		},
		{
			"address": "comdex1yyex40yfa4t3lvgcud9hvazepc8wmjpar932zx",
			"reward": {
				"denom": "ucmdx",
				"amount": "1235"
			}
		},
		{
			"address": "comdex1yy7fpny3pkfzj9zdsq8s0p3n9ldjscrmc65ujs",
			"reward": {
				"denom": "ucmdx",
				"amount": "1678"
			}
		},
		{
			"address": "comdex1yy7lh4cqjljtqtdnuc430r2w4d58jmu9w2khtl",
			"reward": {
				"denom": "ucmdx",
				"amount": "3592"
			}
		},
		{
			"address": "comdex1y9qny2z56e7d8ydpj6cfrchclwklxqhqjlvq0f",
			"reward": {
				"denom": "ucmdx",
				"amount": "1056"
			}
		},
		{
			"address": "comdex1y9qkmk47rca4gzfqxuyh7m28lxhprcu7w7zt3p",
			"reward": {
				"denom": "ucmdx",
				"amount": "2016"
			}
		},
		{
			"address": "comdex1y9rusf025aa3a7tusjhdyycu3nwewksrc855n0",
			"reward": {
				"denom": "ucmdx",
				"amount": "6905"
			}
		},
		{
			"address": "comdex1y9xk0csuget547turu3uqrv5zt6nctpjd6unp8",
			"reward": {
				"denom": "ucmdx",
				"amount": "28440"
			}
		},
		{
			"address": "comdex1y98wpm473qvpg9gtyfs5cu0tr7hymg4qhyk4yn",
			"reward": {
				"denom": "ucmdx",
				"amount": "1463"
			}
		},
		{
			"address": "comdex1y9tf0gkjv4llh4v62p6z4v87gd585rgnp07pgd",
			"reward": {
				"denom": "ucmdx",
				"amount": "479"
			}
		},
		{
			"address": "comdex1y9da46mj76n0ly4hvtglfnfwuezspksyj9r4j4",
			"reward": {
				"denom": "ucmdx",
				"amount": "17604"
			}
		},
		{
			"address": "comdex1y9373yy8fgplazsg2pjnjq3lqg4048jnrjrcrw",
			"reward": {
				"denom": "ucmdx",
				"amount": "1418"
			}
		},
		{
			"address": "comdex1y9na0yxs4lgse025vvdx9vuz0a5h4gtxem0957",
			"reward": {
				"denom": "ucmdx",
				"amount": "18568"
			}
		},
		{
			"address": "comdex1y946vr3nxmxp4ueuunt4hazmvf5srdk9a3uwcv",
			"reward": {
				"denom": "ucmdx",
				"amount": "122940"
			}
		},
		{
			"address": "comdex1y94uk0hdf0sp60amfpyt95pnvx95wtaq4nd4gk",
			"reward": {
				"denom": "ucmdx",
				"amount": "42415"
			}
		},
		{
			"address": "comdex1y9e576jqgrgxsze0qtf45ylt8cn3glc3mr9fzq",
			"reward": {
				"denom": "ucmdx",
				"amount": "178"
			}
		},
		{
			"address": "comdex1y9mhz5qxwjhlcw9cn4nwkzkht6suhv93zhq4f2",
			"reward": {
				"denom": "ucmdx",
				"amount": "166"
			}
		},
		{
			"address": "comdex1y9ucazx4mxze9lwxx4g4wpdmzjxx6d9d2smcu0",
			"reward": {
				"denom": "ucmdx",
				"amount": "151"
			}
		},
		{
			"address": "comdex1y9ls725ta2pzqzdp35n6uf29kxwn65nkg8ux90",
			"reward": {
				"denom": "ucmdx",
				"amount": "2574"
			}
		},
		{
			"address": "comdex1yxp82qlay8768867clrg64kje9dyv4y6ycgwm0",
			"reward": {
				"denom": "ucmdx",
				"amount": "2215"
			}
		},
		{
			"address": "comdex1yxye2anen8t30vlhpj36hfdzr7wk4ml04gty70",
			"reward": {
				"denom": "ucmdx",
				"amount": "9332"
			}
		},
		{
			"address": "comdex1yxg3va3h27sfhxv0fxccuxp2fhk5wlqsfam76q",
			"reward": {
				"denom": "ucmdx",
				"amount": "11098"
			}
		},
		{
			"address": "comdex1yxgefa7h4hnta2cfglm2e7r5cnlwpkxf5chpv4",
			"reward": {
				"denom": "ucmdx",
				"amount": "72522"
			}
		},
		{
			"address": "comdex1yx2h780kl992l3vmvm5ll6la2wwr5r7ljut7gp",
			"reward": {
				"denom": "ucmdx",
				"amount": "16410"
			}
		},
		{
			"address": "comdex1yxtfgaee36t0v94uaquqyw35lpladqxe7taelw",
			"reward": {
				"denom": "ucmdx",
				"amount": "9573"
			}
		},
		{
			"address": "comdex1yxd3a88mfvqnygaxp8t79zkepxlmzw7t2q62ep",
			"reward": {
				"denom": "ucmdx",
				"amount": "14459"
			}
		},
		{
			"address": "comdex1yx0jt6zkz3u6xq4ah8am7k43qqawg7acd2qt65",
			"reward": {
				"denom": "ucmdx",
				"amount": "882"
			}
		},
		{
			"address": "comdex1yxjracs7gkgcrckz9m8wanhc57y0w0wvgk43eq",
			"reward": {
				"denom": "ucmdx",
				"amount": "80"
			}
		},
		{
			"address": "comdex1yx423vp5tk7unduhjsuffpnanuwe0u4hvh480x",
			"reward": {
				"denom": "ucmdx",
				"amount": "427"
			}
		},
		{
			"address": "comdex1yxktgcpq2c282huagmwv7qmedceeul8dthttna",
			"reward": {
				"denom": "ucmdx",
				"amount": "32565"
			}
		},
		{
			"address": "comdex1yxhwzyhercfga0xla236jkucqedddkc7lnpngj",
			"reward": {
				"denom": "ucmdx",
				"amount": "2853"
			}
		},
		{
			"address": "comdex1yxhjnc3u2avecdqwj63n4xeymg958ausgglcpw",
			"reward": {
				"denom": "ucmdx",
				"amount": "334122"
			}
		},
		{
			"address": "comdex1yxchl7wqeleevxz9k09xcmc2pgtc0amcatxm2r",
			"reward": {
				"denom": "ucmdx",
				"amount": "122385"
			}
		},
		{
			"address": "comdex1yx7xge0tsfsf6jg8xnuldnqkqxk8fdn65uy3yr",
			"reward": {
				"denom": "ucmdx",
				"amount": "722"
			}
		},
		{
			"address": "comdex1yx72f2eksw5c0adfrcc9xkyv38czn4gwt7yy4a",
			"reward": {
				"denom": "ucmdx",
				"amount": "6042"
			}
		},
		{
			"address": "comdex1y899dtyacpmmya099jrskefn3r8gmnt8t9gheg",
			"reward": {
				"denom": "ucmdx",
				"amount": "7456"
			}
		},
		{
			"address": "comdex1y895wujvctese26gylm74e4wej9q3r404dc6yx",
			"reward": {
				"denom": "ucmdx",
				"amount": "1983"
			}
		},
		{
			"address": "comdex1y8x9pqh7uklanah5u00rs8mpzn280vg2ds55hv",
			"reward": {
				"denom": "ucmdx",
				"amount": "14243"
			}
		},
		{
			"address": "comdex1y8f69czxpah3p9qdsggz62ykd0m5pmj6p9wr3k",
			"reward": {
				"denom": "ucmdx",
				"amount": "67666"
			}
		},
		{
			"address": "comdex1y8tg6rkzgwav0ef99rwa487y07vut8vt56fggt",
			"reward": {
				"denom": "ucmdx",
				"amount": "3141"
			}
		},
		{
			"address": "comdex1y8vfkas22fxqf5rl2crpjh7zs0567ks3mpkpgy",
			"reward": {
				"denom": "ucmdx",
				"amount": "49166"
			}
		},
		{
			"address": "comdex1y8dxp5jawvtqrn8hd7cj7d6uy5kxc68g2hgvj4",
			"reward": {
				"denom": "ucmdx",
				"amount": "1764"
			}
		},
		{
			"address": "comdex1y8df4y0ydgwgp22txm97pfq4yje895fwxzr7rr",
			"reward": {
				"denom": "ucmdx",
				"amount": "14"
			}
		},
		{
			"address": "comdex1y8d0pg05ldx6q7kr2u3c3w3hk0eaqhetm9fvte",
			"reward": {
				"denom": "ucmdx",
				"amount": "2915"
			}
		},
		{
			"address": "comdex1y8wfswmqnan7ysdfgaqh603wmyjr574m75rptm",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1y8s9maf6r9y936x3gccns3q93qgh5tu0tzw0u7",
			"reward": {
				"denom": "ucmdx",
				"amount": "1412"
			}
		},
		{
			"address": "comdex1y8n2dnzva2ksv4ymxg05ty5a29szwxuw43hdd4",
			"reward": {
				"denom": "ucmdx",
				"amount": "8215"
			}
		},
		{
			"address": "comdex1y8nuctv8z6zsjf5673ql6auhy6x3xngwhtsrnd",
			"reward": {
				"denom": "ucmdx",
				"amount": "1733"
			}
		},
		{
			"address": "comdex1y8k9zz0wq8qdf44uktnze03klen3lt7hfv996s",
			"reward": {
				"denom": "ucmdx",
				"amount": "386"
			}
		},
		{
			"address": "comdex1y8c6ggewsz45euutqax36qww5xqww884n9l0h8",
			"reward": {
				"denom": "ucmdx",
				"amount": "1734"
			}
		},
		{
			"address": "comdex1y8eyz3vf48ydt9ycjlwdrtsfffs79zj3ld95cd",
			"reward": {
				"denom": "ucmdx",
				"amount": "163359"
			}
		},
		{
			"address": "comdex1y8e9msx7xtskng3l4dayh6jjzmgcs90lpkcspr",
			"reward": {
				"denom": "ucmdx",
				"amount": "3378"
			}
		},
		{
			"address": "comdex1y8euullau2ms69ycdkjpsk58ddkv9ekuf9fe8s",
			"reward": {
				"denom": "ucmdx",
				"amount": "29"
			}
		},
		{
			"address": "comdex1y86xrqwegs0hwc3n7hekf7ecyxs3qp8ntq8dcs",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1y86f4qj6gs8fzcrdgh7ye335hkp75czs4dgqvl",
			"reward": {
				"denom": "ucmdx",
				"amount": "419977"
			}
		},
		{
			"address": "comdex1y8upemamwka3m79h35jq5vz89xzjcpt7t5qeph",
			"reward": {
				"denom": "ucmdx",
				"amount": "4284"
			}
		},
		{
			"address": "comdex1y8atckguzrp8yk67j7zvv8l2hy3vyjmzskptqt",
			"reward": {
				"denom": "ucmdx",
				"amount": "743"
			}
		},
		{
			"address": "comdex1y8akmskm44l69ag2lf9t5ahtns8nxadmaek9a3",
			"reward": {
				"denom": "ucmdx",
				"amount": "5103"
			}
		},
		{
			"address": "comdex1y8lzp8danh26us8jef74jgjjpshxnjxhcknvln",
			"reward": {
				"denom": "ucmdx",
				"amount": "42143"
			}
		},
		{
			"address": "comdex1ygqvpc3qrtfjjw8wakz4vc0l27ql26wd4gx7qf",
			"reward": {
				"denom": "ucmdx",
				"amount": "7461"
			}
		},
		{
			"address": "comdex1ygqssma02fnh24net3r99vgt0clrwlqg9lw5ew",
			"reward": {
				"denom": "ucmdx",
				"amount": "667"
			}
		},
		{
			"address": "comdex1ygqcu6c6jsc443apwnad4ck0rrgpcv7hj6a0zh",
			"reward": {
				"denom": "ucmdx",
				"amount": "373"
			}
		},
		{
			"address": "comdex1ygppqsvww2g4h2ev6ze4cutq83t03ugnjz4246",
			"reward": {
				"denom": "ucmdx",
				"amount": "71"
			}
		},
		{
			"address": "comdex1ygphw2ke3clscpm0jsxk3998cy2dd4hts5khct",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1yg9zevglatk27qk8hj98cgc25r7etya72ymmgd",
			"reward": {
				"denom": "ucmdx",
				"amount": "13282"
			}
		},
		{
			"address": "comdex1yg9023cqj3a3v67c0xshcsgpcd9a59pp6u0w0d",
			"reward": {
				"denom": "ucmdx",
				"amount": "84"
			}
		},
		{
			"address": "comdex1yg93sten48lv682wkqg3jlhzee5wxwft0rfhak",
			"reward": {
				"denom": "ucmdx",
				"amount": "1430"
			}
		},
		{
			"address": "comdex1ygxyccsq7nxl5rwg5xe0zmqwup4wdn6tmthn86",
			"reward": {
				"denom": "ucmdx",
				"amount": "268"
			}
		},
		{
			"address": "comdex1yggfnr9kccfms5yh3nterydjfhm7mu5jmj5xmd",
			"reward": {
				"denom": "ucmdx",
				"amount": "5657"
			}
		},
		{
			"address": "comdex1ygf45trm9aw39maze2lxd044lluz332lfccnvr",
			"reward": {
				"denom": "ucmdx",
				"amount": "5542"
			}
		},
		{
			"address": "comdex1ygt5q28smu68ltywdkamwa57ps7w8vvxajfh7w",
			"reward": {
				"denom": "ucmdx",
				"amount": "3518"
			}
		},
		{
			"address": "comdex1ygtmw5sf774h6t0ml0m8jkudrsukwjwlycwc00",
			"reward": {
				"denom": "ucmdx",
				"amount": "2679"
			}
		},
		{
			"address": "comdex1ygvjf8pfguws8crgl59k9n8w0zh7ajev0jr2m6",
			"reward": {
				"denom": "ucmdx",
				"amount": "2652"
			}
		},
		{
			"address": "comdex1ygvm6wdxv4zgpkhurqaddlvtwherryvm8jyqel",
			"reward": {
				"denom": "ucmdx",
				"amount": "374020"
			}
		},
		{
			"address": "comdex1ygvaedvgxdvxfj9qada0pwvdjm4t75xqghacul",
			"reward": {
				"denom": "ucmdx",
				"amount": "34723"
			}
		},
		{
			"address": "comdex1yg0fa6aa8npeyweedzct2z0qsj0q7xypzkqwnp",
			"reward": {
				"denom": "ucmdx",
				"amount": "8980"
			}
		},
		{
			"address": "comdex1ygsgtmn057xn7xl30hz0a3y5p7p7ltyhfuju79",
			"reward": {
				"denom": "ucmdx",
				"amount": "14155"
			}
		},
		{
			"address": "comdex1yg382agy3f49qg6nhh3vdxfnvxnr3kaftm3h70",
			"reward": {
				"denom": "ucmdx",
				"amount": "6171"
			}
		},
		{
			"address": "comdex1yg3lqm8degawqwd5f8pf9mt7nf76elv3d78y3r",
			"reward": {
				"denom": "ucmdx",
				"amount": "4966"
			}
		},
		{
			"address": "comdex1ygkxhl355eg2v0p5hc5x6x58ycs5le5vjs6hma",
			"reward": {
				"denom": "ucmdx",
				"amount": "3735"
			}
		},
		{
			"address": "comdex1ygkthe8l5p6sg68hn4myjfddmmvxselhtqggsz",
			"reward": {
				"denom": "ucmdx",
				"amount": "244"
			}
		},
		{
			"address": "comdex1ygc3r3sltxc3da6g3u3eew9ae0mtrt34a34wmn",
			"reward": {
				"denom": "ucmdx",
				"amount": "169"
			}
		},
		{
			"address": "comdex1ygm527cwaxh2y629pdrc7ky5pr3cdk8zqjjyzn",
			"reward": {
				"denom": "ucmdx",
				"amount": "12567"
			}
		},
		{
			"address": "comdex1ygarrkevfqehg4hlgwtw3mdzygm3rzpc6vmnk3",
			"reward": {
				"denom": "ucmdx",
				"amount": "33553"
			}
		},
		{
			"address": "comdex1ygafkw446ur0lzye6kupyg0rg6waqfs7036sk2",
			"reward": {
				"denom": "ucmdx",
				"amount": "697"
			}
		},
		{
			"address": "comdex1ygadp7chakalsju5thj4aw26ehwkuz3teakct8",
			"reward": {
				"denom": "ucmdx",
				"amount": "14097"
			}
		},
		{
			"address": "comdex1yfqmdwqdzdynwqu4akd4ev344njw42v3g94f7d",
			"reward": {
				"denom": "ucmdx",
				"amount": "6507"
			}
		},
		{
			"address": "comdex1yfqu3r6d26n8jckgtrauw9gtyt35s90xk6j4jd",
			"reward": {
				"denom": "ucmdx",
				"amount": "16824"
			}
		},
		{
			"address": "comdex1yfy82r2tkxqy58jsfuezffx2e5wwstkszjvtn2",
			"reward": {
				"denom": "ucmdx",
				"amount": "23426"
			}
		},
		{
			"address": "comdex1yfyef54dtrreenpyl7v779muqxslhyw6datlvf",
			"reward": {
				"denom": "ucmdx",
				"amount": "19183"
			}
		},
		{
			"address": "comdex1yfxrjz33javqe06lnwwyjw2mrf40ylpkdm7j92",
			"reward": {
				"denom": "ucmdx",
				"amount": "1792"
			}
		},
		{
			"address": "comdex1yfxm2lkwp4476nq362hew60mf5yw7x60wplze0",
			"reward": {
				"denom": "ucmdx",
				"amount": "125"
			}
		},
		{
			"address": "comdex1yftwzn4lfugsc3sh93c5gc2txaj4me8v7m7ujh",
			"reward": {
				"denom": "ucmdx",
				"amount": "530"
			}
		},
		{
			"address": "comdex1yf0ng6dmgsq63pxw7k0kt2pd3370gq99duj53j",
			"reward": {
				"denom": "ucmdx",
				"amount": "529"
			}
		},
		{
			"address": "comdex1yfne7gg52umva4e72qa95w6wd09ads8ezl5tya",
			"reward": {
				"denom": "ucmdx",
				"amount": "6336"
			}
		},
		{
			"address": "comdex1yf6p95r50l7a8fjuqalg5vxx425xdu7c84hjn2",
			"reward": {
				"denom": "ucmdx",
				"amount": "5721"
			}
		},
		{
			"address": "comdex1y2qzkezsdrs93nzewfsehj40k6mgj5n2p2mqed",
			"reward": {
				"denom": "ucmdx",
				"amount": "14889"
			}
		},
		{
			"address": "comdex1y2ptynyf78vjl7agjj7e7x9t2k3gwgv550q8lz",
			"reward": {
				"denom": "ucmdx",
				"amount": "173"
			}
		},
		{
			"address": "comdex1y2rfrgh374gd40gj0yusjasw2paaysah7hdzzw",
			"reward": {
				"denom": "ucmdx",
				"amount": "8804"
			}
		},
		{
			"address": "comdex1y2r0wnw4veqchnzzhpsd7n4sw5xeykrq6lsame",
			"reward": {
				"denom": "ucmdx",
				"amount": "181"
			}
		},
		{
			"address": "comdex1y2rnzjck3lpuy0zljn0t7g6l4kjeesyd5ukytj",
			"reward": {
				"denom": "ucmdx",
				"amount": "57"
			}
		},
		{
			"address": "comdex1y2yffdmjl8jjde2yfkwj30e0c7m7efgrqcxukh",
			"reward": {
				"denom": "ucmdx",
				"amount": "3456"
			}
		},
		{
			"address": "comdex1y293l23yyk8vzn9cgl4tvd8yqujpc69crn3egd",
			"reward": {
				"denom": "ucmdx",
				"amount": "1816"
			}
		},
		{
			"address": "comdex1y2tkzjdjvcdj3v8l0jmfjukqhumgsqd0va5gdj",
			"reward": {
				"denom": "ucmdx",
				"amount": "1782"
			}
		},
		{
			"address": "comdex1y2ducl7a7mal56vjh9a0g3d3dcnx93k8qm5azs",
			"reward": {
				"denom": "ucmdx",
				"amount": "44123"
			}
		},
		{
			"address": "comdex1y2jfyamd4h3cpc70vussezmvkegwgx666y3jl9",
			"reward": {
				"denom": "ucmdx",
				"amount": "6236"
			}
		},
		{
			"address": "comdex1y2jaag9mjls82u0cvefer3gxdf6akus54pqscd",
			"reward": {
				"denom": "ucmdx",
				"amount": "3846"
			}
		},
		{
			"address": "comdex1y24rcvqehf0fh7sppg2amgw4s4s3dd2d7h8023",
			"reward": {
				"denom": "ucmdx",
				"amount": "3076"
			}
		},
		{
			"address": "comdex1y24yd9dmvgr4va2utj2836ctrt4znf905q4m0g",
			"reward": {
				"denom": "ucmdx",
				"amount": "1464"
			}
		},
		{
			"address": "comdex1y24d6p64lmt98ttl3z9al96xezxyet2zta4g7y",
			"reward": {
				"denom": "ucmdx",
				"amount": "17104"
			}
		},
		{
			"address": "comdex1y2k3ykd0n5w88hcswl2zut74ltpnngq524an7q",
			"reward": {
				"denom": "ucmdx",
				"amount": "126"
			}
		},
		{
			"address": "comdex1y2knaz4gga4rrce38w5evnztsueq937pcvf9dd",
			"reward": {
				"denom": "ucmdx",
				"amount": "35523"
			}
		},
		{
			"address": "comdex1y2k4nj4s3e96hqzu6xm36424axq35e7f7ahzew",
			"reward": {
				"denom": "ucmdx",
				"amount": "15530"
			}
		},
		{
			"address": "comdex1y2kux97fe8a3zjkxd9j3ymmgq2eu0ptj529en0",
			"reward": {
				"denom": "ucmdx",
				"amount": "12312"
			}
		},
		{
			"address": "comdex1y2lan2anmfwsz9p00l5j2rd7agclezyq4gmzj0",
			"reward": {
				"denom": "ucmdx",
				"amount": "74"
			}
		},
		{
			"address": "comdex1ytr0nujljr44t7kw2vhe566ecjz8mtn859wuln",
			"reward": {
				"denom": "ucmdx",
				"amount": "16598"
			}
		},
		{
			"address": "comdex1yt8zcghfq8pmhxsc95fqqua44y4rautw6z0gux",
			"reward": {
				"denom": "ucmdx",
				"amount": "13082"
			}
		},
		{
			"address": "comdex1yt8yzj806dvfl6ya6jjr0anxg4rs245n6ps83d",
			"reward": {
				"denom": "ucmdx",
				"amount": "10553"
			}
		},
		{
			"address": "comdex1yt8jd7695ynmrpsfjuccl7vsvxhdj6yhllu6s8",
			"reward": {
				"denom": "ucmdx",
				"amount": "57430"
			}
		},
		{
			"address": "comdex1ytg3crmdxpm7amn2s6hs699f0sz7cws4lppyhh",
			"reward": {
				"denom": "ucmdx",
				"amount": "62094"
			}
		},
		{
			"address": "comdex1ytghjrywcnr3ntdvhvfu7jkk3gqmyx4e6k2zln",
			"reward": {
				"denom": "ucmdx",
				"amount": "315"
			}
		},
		{
			"address": "comdex1ytf6yu225p7yjy4cnfxmanj5f6xlu6wsq3f534",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1ytvfhfgmss5ffh6sd0u3x4fntxd34a5cv552us",
			"reward": {
				"denom": "ucmdx",
				"amount": "2653"
			}
		},
		{
			"address": "comdex1ytwjccum7632twayxyel47f8t4srhvvu2y8tmc",
			"reward": {
				"denom": "ucmdx",
				"amount": "176"
			}
		},
		{
			"address": "comdex1yt3rua6q546h86hyqtphjs0jg5ccm7s7esynlw",
			"reward": {
				"denom": "ucmdx",
				"amount": "1067332"
			}
		},
		{
			"address": "comdex1ytjal6t0qagtuqfs36v0lv8wlkmq6qtpvjqq7c",
			"reward": {
				"denom": "ucmdx",
				"amount": "136717"
			}
		},
		{
			"address": "comdex1yt5gyfega3h6c58fsfvwkj7xm679h9lzy6khwv",
			"reward": {
				"denom": "ucmdx",
				"amount": "6991"
			}
		},
		{
			"address": "comdex1yth5zyerana5uhzasgzd2v7hvtym628nvjewvf",
			"reward": {
				"denom": "ucmdx",
				"amount": "1301"
			}
		},
		{
			"address": "comdex1ytem9skzmq7n7tlfcqjw6wqfgd587cn3u76haw",
			"reward": {
				"denom": "ucmdx",
				"amount": "358"
			}
		},
		{
			"address": "comdex1yt6pkw70sk3rfzt6fft2vdhh0nvvpwmzqm3gc3",
			"reward": {
				"denom": "ucmdx",
				"amount": "142"
			}
		},
		{
			"address": "comdex1ytmk4d84agsrp4zhg0wtsgpqfux05aed5gha35",
			"reward": {
				"denom": "ucmdx",
				"amount": "9065"
			}
		},
		{
			"address": "comdex1ytagn5zvpltzg4y5k45qm7amyfscpqjewcs2s8",
			"reward": {
				"denom": "ucmdx",
				"amount": "14372"
			}
		},
		{
			"address": "comdex1ytlcfhy73wrteqnsg67ew3j3uphq5yl2gmz47w",
			"reward": {
				"denom": "ucmdx",
				"amount": "3047"
			}
		},
		{
			"address": "comdex1yvqxf8n0tywq5hm9vqw9ypkk40x832t34qgudc",
			"reward": {
				"denom": "ucmdx",
				"amount": "1734"
			}
		},
		{
			"address": "comdex1yvr5w8cmcsqec03pa9vpvdnm20w3kundyned7j",
			"reward": {
				"denom": "ucmdx",
				"amount": "86"
			}
		},
		{
			"address": "comdex1yv9dtu7ke2s6txjy3qqm8qr0yn0h3jc9dgyr4a",
			"reward": {
				"denom": "ucmdx",
				"amount": "697"
			}
		},
		{
			"address": "comdex1yv8y7let7pkam2wzntf7xagqry306y3k4g2p5t",
			"reward": {
				"denom": "ucmdx",
				"amount": "1442"
			}
		},
		{
			"address": "comdex1yvg48vsfj295xyl66equg0v9rulg09ghkq34lf",
			"reward": {
				"denom": "ucmdx",
				"amount": "173"
			}
		},
		{
			"address": "comdex1yvvqdumhffxx2vtnhqarjp28kd93s7m8nut2gm",
			"reward": {
				"denom": "ucmdx",
				"amount": "180"
			}
		},
		{
			"address": "comdex1yvd3kz09q8wf8ta4tmtmqlc20erqa58k7052p8",
			"reward": {
				"denom": "ucmdx",
				"amount": "11851"
			}
		},
		{
			"address": "comdex1yv0dk6hapz4cjz2kyty0rgwp6ju609h89q9egm",
			"reward": {
				"denom": "ucmdx",
				"amount": "2026"
			}
		},
		{
			"address": "comdex1yvsp4fvj6wwh9wwdqrj4l86nf2sz4czg2fhssr",
			"reward": {
				"denom": "ucmdx",
				"amount": "833"
			}
		},
		{
			"address": "comdex1yvsdumzty3k4tasp9krzr454lgtkwgzemgwate",
			"reward": {
				"denom": "ucmdx",
				"amount": "141624"
			}
		},
		{
			"address": "comdex1yv5eyg5nfmh5cnvaj0xe6dptr07tf7lmeg8yq2",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1yvklrnhjtuqrfuvl2umqaf24u8gu6al43v0kmx",
			"reward": {
				"denom": "ucmdx",
				"amount": "523"
			}
		},
		{
			"address": "comdex1yve4z5g0keuqgu6al53n6dxtpgfqunmmtdn0pn",
			"reward": {
				"denom": "ucmdx",
				"amount": "124891"
			}
		},
		{
			"address": "comdex1yvelgm0sgypdyu8qxe28c5fnyfhwhyvrj3fpsw",
			"reward": {
				"denom": "ucmdx",
				"amount": "823"
			}
		},
		{
			"address": "comdex1yv6grsrwrrq6p49f54h9l8r3uszav4gquzqqk5",
			"reward": {
				"denom": "ucmdx",
				"amount": "70855"
			}
		},
		{
			"address": "comdex1yv6nx4xz0cnp3da5yaudm5325upk38zf8j7ctg",
			"reward": {
				"denom": "ucmdx",
				"amount": "7493"
			}
		},
		{
			"address": "comdex1yvmp37yd6ppxpn97ns44en3ae6vw0as4c8jnrf",
			"reward": {
				"denom": "ucmdx",
				"amount": "1236"
			}
		},
		{
			"address": "comdex1yvmkxpts6rl9eaptpy2c0lxrmm8wxxn4xfskdh",
			"reward": {
				"denom": "ucmdx",
				"amount": "5279"
			}
		},
		{
			"address": "comdex1ydrdqwjp9ejxclsd4c4zavlc3tuh94gpplpleu",
			"reward": {
				"denom": "ucmdx",
				"amount": "1766"
			}
		},
		{
			"address": "comdex1ydrskl8w27n3dyxzgl233uazpgxarerd75ft2a",
			"reward": {
				"denom": "ucmdx",
				"amount": "26797"
			}
		},
		{
			"address": "comdex1yd88m4nu97z7lkastdkhqdql7wxxzfyatytnff",
			"reward": {
				"denom": "ucmdx",
				"amount": "1750"
			}
		},
		{
			"address": "comdex1ydgfy54c638dnv5dqcl453azum29ezsnevrck0",
			"reward": {
				"denom": "ucmdx",
				"amount": "52"
			}
		},
		{
			"address": "comdex1ydgl3ttt9akrakl6kdpxynp4r0trqgl9s9aqjl",
			"reward": {
				"denom": "ucmdx",
				"amount": "791"
			}
		},
		{
			"address": "comdex1ydfz3mr0m452ntk4t8cw2myefqm3sf3ngtdryq",
			"reward": {
				"denom": "ucmdx",
				"amount": "20704"
			}
		},
		{
			"address": "comdex1ydf7vallans4czw000g242vxqdjan020vxpgyf",
			"reward": {
				"denom": "ucmdx",
				"amount": "71"
			}
		},
		{
			"address": "comdex1yd2ary9extd5t3l3k29ktlq9gff9h5l8gxcmga",
			"reward": {
				"denom": "ucmdx",
				"amount": "6501"
			}
		},
		{
			"address": "comdex1ydtssfe6lmn34hhz03j94h003d35ynqrcj5l4l",
			"reward": {
				"denom": "ucmdx",
				"amount": "3456"
			}
		},
		{
			"address": "comdex1ydthxdntqc48ajqy87ecd343l0sfm3e625xa04",
			"reward": {
				"denom": "ucmdx",
				"amount": "2728"
			}
		},
		{
			"address": "comdex1ydtaj62cm9z0nz2pskh2mld0t9d444p6r9h5cz",
			"reward": {
				"denom": "ucmdx",
				"amount": "219644"
			}
		},
		{
			"address": "comdex1ydvzhetzz0k9y0hmhh8syw566gy48nkhnskmh7",
			"reward": {
				"denom": "ucmdx",
				"amount": "64446"
			}
		},
		{
			"address": "comdex1ydsmt3qjq9aev9l7us4uh0g4fydwhu5c9dx3vq",
			"reward": {
				"denom": "ucmdx",
				"amount": "4126"
			}
		},
		{
			"address": "comdex1ydn2luca7w0xapuen04axaw3fsea437syhlzf8",
			"reward": {
				"denom": "ucmdx",
				"amount": "1314"
			}
		},
		{
			"address": "comdex1ydnaxvjppm45jfuck6zkcsdskleggsezcnd5v6",
			"reward": {
				"denom": "ucmdx",
				"amount": "6192"
			}
		},
		{
			"address": "comdex1yd50s0fsn0eyp8kx6wvgrfmjjgr6tejq3znrj5",
			"reward": {
				"denom": "ucmdx",
				"amount": "21605"
			}
		},
		{
			"address": "comdex1ydk94ht9d6s0uexgzmdnx4h6xecz0qn05fe9cd",
			"reward": {
				"denom": "ucmdx",
				"amount": "150"
			}
		},
		{
			"address": "comdex1ydcpt926p7wlurlsvll5y2u7w7vu5hj9dvvjmj",
			"reward": {
				"denom": "ucmdx",
				"amount": "27278"
			}
		},
		{
			"address": "comdex1yda82h060gxnqaahaz9x0uvram8xm665fyq88m",
			"reward": {
				"denom": "ucmdx",
				"amount": "1163"
			}
		},
		{
			"address": "comdex1ydav5f82ztcjdpxuc94v5ff6z9pgcp6ffzw24m",
			"reward": {
				"denom": "ucmdx",
				"amount": "836"
			}
		},
		{
			"address": "comdex1ywqvj9dhxzm3hjmtnwznqexhw4neg2jn2m9p9n",
			"reward": {
				"denom": "ucmdx",
				"amount": "3568"
			}
		},
		{
			"address": "comdex1ywyk2knhmsql652dysecaxtev35sly96fec283",
			"reward": {
				"denom": "ucmdx",
				"amount": "8982"
			}
		},
		{
			"address": "comdex1yw9tqtt6ugtm9lae77fx6lwa0ecfxg9n2ednce",
			"reward": {
				"denom": "ucmdx",
				"amount": "16692"
			}
		},
		{
			"address": "comdex1ywxltzxvyxel8nualanz9vvj4q3yp4cjjnaewt",
			"reward": {
				"denom": "ucmdx",
				"amount": "350"
			}
		},
		{
			"address": "comdex1ywga7phr2vdqq7jpd5mc0sa2y2musy3mwdtsgz",
			"reward": {
				"denom": "ucmdx",
				"amount": "39"
			}
		},
		{
			"address": "comdex1ywv2tl5hvzlvd0snv3ykhkv4p899pdtvgm95a5",
			"reward": {
				"denom": "ucmdx",
				"amount": "723"
			}
		},
		{
			"address": "comdex1ywwr6elt0h2v7stjzqf3t7cceqd6njk590tz5v",
			"reward": {
				"denom": "ucmdx",
				"amount": "181"
			}
		},
		{
			"address": "comdex1ywjmt4cvwjrc40vl5zkl73tme9j6yup6exc4wr",
			"reward": {
				"denom": "ucmdx",
				"amount": "13612"
			}
		},
		{
			"address": "comdex1yw49hzz0muyvukj53g0ueqg3zwkkuuxkv4wykk",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1ywkgk4txq9selrldfljdwkjkgrl5l3epl8s2z4",
			"reward": {
				"denom": "ucmdx",
				"amount": "58058"
			}
		},
		{
			"address": "comdex1yw7p526sdmkc4lyqd3vksy3qtw9fs57qluvcsy",
			"reward": {
				"denom": "ucmdx",
				"amount": "1774"
			}
		},
		{
			"address": "comdex1yw7plgs2vaaz5wpay4gcu0t4lkgyk9au3a32gq",
			"reward": {
				"denom": "ucmdx",
				"amount": "10448"
			}
		},
		{
			"address": "comdex1yw7yw5hjqkmgf9rfncsl75dpf8sj5326j30l9d",
			"reward": {
				"denom": "ucmdx",
				"amount": "16"
			}
		},
		{
			"address": "comdex1y0rt5fney27qyxrzcvyawa7m8nw2y7eedrjlqe",
			"reward": {
				"denom": "ucmdx",
				"amount": "530"
			}
		},
		{
			"address": "comdex1y0ypty2sf860crwnqmvyprd46la6xx6e8rqfua",
			"reward": {
				"denom": "ucmdx",
				"amount": "1696"
			}
		},
		{
			"address": "comdex1y0xw2wh78syt7rzr9lfx2dnnd5ya79ttt92xpd",
			"reward": {
				"denom": "ucmdx",
				"amount": "1768"
			}
		},
		{
			"address": "comdex1y0f3unk3zr6w8k33nz0yfnuphdewy7mskym286",
			"reward": {
				"denom": "ucmdx",
				"amount": "297121"
			}
		},
		{
			"address": "comdex1y0wn2xmh3tu74qzrr8v86rnrgawejzatjrtrck",
			"reward": {
				"denom": "ucmdx",
				"amount": "706"
			}
		},
		{
			"address": "comdex1y03r3dgyy29gcfdhdph7revhstjgknfwlytw8a",
			"reward": {
				"denom": "ucmdx",
				"amount": "1402"
			}
		},
		{
			"address": "comdex1y0k86wx0jpfqs9qcq6alalfwvyackwfqw5px8f",
			"reward": {
				"denom": "ucmdx",
				"amount": "604"
			}
		},
		{
			"address": "comdex1y0knneeefwzhejlh6a09g8h94ru0dtucxwwrh5",
			"reward": {
				"denom": "ucmdx",
				"amount": "1666"
			}
		},
		{
			"address": "comdex1y0cclpqqc386lqc7ekp3ae4jw4hmqe2qaw0v8f",
			"reward": {
				"denom": "ucmdx",
				"amount": "167"
			}
		},
		{
			"address": "comdex1y0madmkzkhqcytafmtetgmfh5u0ymxu22gamgp",
			"reward": {
				"denom": "ucmdx",
				"amount": "34967"
			}
		},
		{
			"address": "comdex1y0avy9vnrgflmydzpe3y0ndl8fz8djzl0rnmdm",
			"reward": {
				"denom": "ucmdx",
				"amount": "21591"
			}
		},
		{
			"address": "comdex1y0a7xtd0e2jpsg6xmwyszztq2zyruplfq6z7t3",
			"reward": {
				"denom": "ucmdx",
				"amount": "180"
			}
		},
		{
			"address": "comdex1y07sn3nhlt43u6xdvcnk20xd98pfn79k9hd87p",
			"reward": {
				"denom": "ucmdx",
				"amount": "3094"
			}
		},
		{
			"address": "comdex1ysq9gst3vx5yt4wmkeedl6jhqq0zkxtthqrep5",
			"reward": {
				"denom": "ucmdx",
				"amount": "7579"
			}
		},
		{
			"address": "comdex1yspgv274rrewz43hpvwk45qwafle4jj58qpxnk",
			"reward": {
				"denom": "ucmdx",
				"amount": "28879"
			}
		},
		{
			"address": "comdex1ys966422848fj6zz6ya3nvyfwqqqdps4eqkult",
			"reward": {
				"denom": "ucmdx",
				"amount": "1264"
			}
		},
		{
			"address": "comdex1ys9ml0ke0r9d5e0t9nsneltx0fy25xqsvdg2ua",
			"reward": {
				"denom": "ucmdx",
				"amount": "16291"
			}
		},
		{
			"address": "comdex1ysx3fw33mqrvedq0zhfduvgvecml2ctsy8mqrn",
			"reward": {
				"denom": "ucmdx",
				"amount": "2850"
			}
		},
		{
			"address": "comdex1ysx42gtfu2wfzu8fhggvzfxwp7ehznx64z883s",
			"reward": {
				"denom": "ucmdx",
				"amount": "3109"
			}
		},
		{
			"address": "comdex1ysg7gh0jd5rf9lzfj9nqevztpld7nm3e0kju0p",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1ysfd9up9g7pw0kdx97xcsh778y33cq0xt8revk",
			"reward": {
				"denom": "ucmdx",
				"amount": "303"
			}
		},
		{
			"address": "comdex1ysflmrj4x5lafn6fj0ufrtcx3k723wy3kyvx5g",
			"reward": {
				"denom": "ucmdx",
				"amount": "1754"
			}
		},
		{
			"address": "comdex1ysv30ul5a2fm5ssfx8hlsn0nrqdvw2mvz6c3er",
			"reward": {
				"denom": "ucmdx",
				"amount": "2164"
			}
		},
		{
			"address": "comdex1ysva577vu0wwaare4ne5ms8pyz97cxmyaxyyzs",
			"reward": {
				"denom": "ucmdx",
				"amount": "390"
			}
		},
		{
			"address": "comdex1ysdahlm4yz78asyhpp7ywkjymwqxs284xavpu8",
			"reward": {
				"denom": "ucmdx",
				"amount": "26571"
			}
		},
		{
			"address": "comdex1yswwg94jj6pesrsrsjm0435gduvle7qlrjaclf",
			"reward": {
				"denom": "ucmdx",
				"amount": "19"
			}
		},
		{
			"address": "comdex1ys0vwlesj93sjmglvtxr0sjwrrng74kpns74rq",
			"reward": {
				"denom": "ucmdx",
				"amount": "5599"
			}
		},
		{
			"address": "comdex1ys3uh7v2sk4kn0czs3q5mhlt4cqyzzdh5gesdn",
			"reward": {
				"denom": "ucmdx",
				"amount": "1438"
			}
		},
		{
			"address": "comdex1ys5qgqw8q5jttqmythqcxwz4sexc5m5fqnmlgc",
			"reward": {
				"denom": "ucmdx",
				"amount": "1984"
			}
		},
		{
			"address": "comdex1ysmaql0ng57jslg628x985tm4a4z7cw2pwsqwy",
			"reward": {
				"denom": "ucmdx",
				"amount": "289"
			}
		},
		{
			"address": "comdex1yslf2mgsa4cn9c35vvvzsngwz4thtd4s9ka6ca",
			"reward": {
				"denom": "ucmdx",
				"amount": "84"
			}
		},
		{
			"address": "comdex1y3zvvetqyqr28sztw9ge89gm584f63vlwr5395",
			"reward": {
				"denom": "ucmdx",
				"amount": "115"
			}
		},
		{
			"address": "comdex1y3rc0sgpvjuy94f4sumfw7c0t3qkyyymuq47s2",
			"reward": {
				"denom": "ucmdx",
				"amount": "1717"
			}
		},
		{
			"address": "comdex1y3yezqnqas0nm4m607z45c56ngr245czapyxry",
			"reward": {
				"denom": "ucmdx",
				"amount": "2778"
			}
		},
		{
			"address": "comdex1y3xn3qhllq5a07p385cpr5thl026zek4j5zn2q",
			"reward": {
				"denom": "ucmdx",
				"amount": "14345"
			}
		},
		{
			"address": "comdex1y38y5vat305vllk4se8vqhy53lylxqqlvac6ph",
			"reward": {
				"denom": "ucmdx",
				"amount": "757"
			}
		},
		{
			"address": "comdex1y3fqr58xdsayxqqj2wzqgqzl6gqjm9xj96x3sv",
			"reward": {
				"denom": "ucmdx",
				"amount": "2305"
			}
		},
		{
			"address": "comdex1y3f267ejv7skw9rndwtutq9satcguuzw24w3lc",
			"reward": {
				"denom": "ucmdx",
				"amount": "1389"
			}
		},
		{
			"address": "comdex1y32n53d2uynw6fxrrhvat3z4dwgpm5tml3355f",
			"reward": {
				"denom": "ucmdx",
				"amount": "403"
			}
		},
		{
			"address": "comdex1y3d4skv5f08fh68sxzqxezm62k94sz0wkfd3k2",
			"reward": {
				"denom": "ucmdx",
				"amount": "156"
			}
		},
		{
			"address": "comdex1y3wmyt70ws9zuk7dj88x29dskhh36weygq8vjp",
			"reward": {
				"denom": "ucmdx",
				"amount": "28936"
			}
		},
		{
			"address": "comdex1y3spp8pm8006fd58x8uv6rzlgsca2uq8cxr7g0",
			"reward": {
				"denom": "ucmdx",
				"amount": "3276"
			}
		},
		{
			"address": "comdex1y33ykptmccjasaasw9eql8mjr52ucqyjudf367",
			"reward": {
				"denom": "ucmdx",
				"amount": "229"
			}
		},
		{
			"address": "comdex1y3596u7nkyjxfzg5qh69vdnpzrxfzgzqsngzzu",
			"reward": {
				"denom": "ucmdx",
				"amount": "16021"
			}
		},
		{
			"address": "comdex1y3cv85veatdv6wupp6u25jqtk4nq6n0684zydn",
			"reward": {
				"denom": "ucmdx",
				"amount": "102"
			}
		},
		{
			"address": "comdex1y3eh67nthcvnc3epn39ltv7qf2478me7nch9wr",
			"reward": {
				"denom": "ucmdx",
				"amount": "61"
			}
		},
		{
			"address": "comdex1y3aq3jmurtswgkkpp74nj8ute9yc82uujlrukm",
			"reward": {
				"denom": "ucmdx",
				"amount": "523"
			}
		},
		{
			"address": "comdex1yjq0n2ewufluenyyvj2y9sead9jfstpxu5acf8",
			"reward": {
				"denom": "ucmdx",
				"amount": "419"
			}
		},
		{
			"address": "comdex1yjq59fx4a4ahns9uje8jqvasyaafzxfa9huf6s",
			"reward": {
				"denom": "ucmdx",
				"amount": "87"
			}
		},
		{
			"address": "comdex1yjpuqyt44u4wshvn00qdkvkflyzu3zqqv6f7qm",
			"reward": {
				"denom": "ucmdx",
				"amount": "3960"
			}
		},
		{
			"address": "comdex1yj8gd4ql3dxjcgucuy09lg955yg3mxpm0xkplk",
			"reward": {
				"denom": "ucmdx",
				"amount": "2232"
			}
		},
		{
			"address": "comdex1yj849zqcl9vgcd4hxdxtl7p6qahck6nlanaweq",
			"reward": {
				"denom": "ucmdx",
				"amount": "142"
			}
		},
		{
			"address": "comdex1yj87fzm6xwwpd5y92frwn0gg52aeusn6fka6er",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1yjfq497pa2y6lc7qxpmx80v5sv7kgdal55kpah",
			"reward": {
				"denom": "ucmdx",
				"amount": "5168"
			}
		},
		{
			"address": "comdex1yjfxtgr2jpe8f5m80gu4lna5n4p39dacpveyy4",
			"reward": {
				"denom": "ucmdx",
				"amount": "17809"
			}
		},
		{
			"address": "comdex1yjv2wc0m0w5a6ze5c2az72jp8093y8ttwr6cde",
			"reward": {
				"denom": "ucmdx",
				"amount": "29"
			}
		},
		{
			"address": "comdex1yjwglsq96qx5l7a7pehvf49p9suegfdyu0hxvk",
			"reward": {
				"denom": "ucmdx",
				"amount": "89"
			}
		},
		{
			"address": "comdex1yjs2njh5yv0x7etxn53fmxvyl757dj9wf984ze",
			"reward": {
				"denom": "ucmdx",
				"amount": "345246"
			}
		},
		{
			"address": "comdex1yjs0a9fuzxlr4hnd7r7q2z03duur2lww7jsp3l",
			"reward": {
				"denom": "ucmdx",
				"amount": "12216"
			}
		},
		{
			"address": "comdex1yjh92ya4f4fmfu4q8tmkkqu6e7l3nlt4udycfv",
			"reward": {
				"denom": "ucmdx",
				"amount": "7120"
			}
		},
		{
			"address": "comdex1yjhtjtw73czsf3580frmvg92ftjgjv2zrct9e2",
			"reward": {
				"denom": "ucmdx",
				"amount": "400"
			}
		},
		{
			"address": "comdex1yjcddkk7xfz2hg7srs2q23u2hne02xww8khq8d",
			"reward": {
				"denom": "ucmdx",
				"amount": "181726"
			}
		},
		{
			"address": "comdex1yjcaq0ym2krfscx5cqz09pj0503y3ppnpgm8z5",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1yj6qg3ac5tym0lerh0jjnwkdzlywgzydc24qlr",
			"reward": {
				"denom": "ucmdx",
				"amount": "4284"
			}
		},
		{
			"address": "comdex1yj6ldz9tez9tatv6wmxfnrg2q2nl6m0zfphlz0",
			"reward": {
				"denom": "ucmdx",
				"amount": "345"
			}
		},
		{
			"address": "comdex1yjmeqn8ejucxnh8rk9308smjj9e3sqrcjm0w54",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1yjupvz04v2pry6qan358rj6rv0c6s0trxsqxva",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1yj7ddxu2j4yct94zw2wscez7czm4gcl8m0kr5g",
			"reward": {
				"denom": "ucmdx",
				"amount": "172"
			}
		},
		{
			"address": "comdex1ynq7mv7an353ug66h8tcgvj75ax582rsrhnyyu",
			"reward": {
				"denom": "ucmdx",
				"amount": "716"
			}
		},
		{
			"address": "comdex1ynp5697gjtesj2qd4qnmrdyfl89sfy7p0lumj8",
			"reward": {
				"denom": "ucmdx",
				"amount": "3694"
			}
		},
		{
			"address": "comdex1ynyfsrlk325x8wddrsv0caus5u59ufghxldkhr",
			"reward": {
				"denom": "ucmdx",
				"amount": "1190"
			}
		},
		{
			"address": "comdex1yn9yrgh5f3an23y8g0ktt5c6uc285lg8c2s3nw",
			"reward": {
				"denom": "ucmdx",
				"amount": "2915885"
			}
		},
		{
			"address": "comdex1yn9teqyxd66fmwtsyn3l4tldu3zykuat547yvp",
			"reward": {
				"denom": "ucmdx",
				"amount": "352"
			}
		},
		{
			"address": "comdex1ynxnxp6t5fehpwtumx7adxymd4pljum65ke5pj",
			"reward": {
				"denom": "ucmdx",
				"amount": "1412"
			}
		},
		{
			"address": "comdex1ynx5j74nyjl4h8defje77vndnurl9up2s0f4xm",
			"reward": {
				"denom": "ucmdx",
				"amount": "12399"
			}
		},
		{
			"address": "comdex1yng27xvqv0z7uszrzvcfkggdpnxgartqd0tfrx",
			"reward": {
				"denom": "ucmdx",
				"amount": "4111"
			}
		},
		{
			"address": "comdex1yn22396antzf5n6z89jdxskjc4tz4j329cjpcg",
			"reward": {
				"denom": "ucmdx",
				"amount": "10264"
			}
		},
		{
			"address": "comdex1ynv040cnmetr4x43pmczx26cww62kpjqtkd8n5",
			"reward": {
				"denom": "ucmdx",
				"amount": "434"
			}
		},
		{
			"address": "comdex1ynd4859f58lz0999hsl9ggaaex2g47d7rxgxkp",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1yndaeyl63lqyqr8hks0dtqyy27swkmyv48n034",
			"reward": {
				"denom": "ucmdx",
				"amount": "2034"
			}
		},
		{
			"address": "comdex1ynsg9ejgz9rc5rnm7aej8p7q25xxd95v4v5czy",
			"reward": {
				"denom": "ucmdx",
				"amount": "18079"
			}
		},
		{
			"address": "comdex1ynskedg6pkntpqv302d0l0p54vemxvffyw6sgp",
			"reward": {
				"denom": "ucmdx",
				"amount": "1445"
			}
		},
		{
			"address": "comdex1ync9qf3femhh5cu5qjwcd5mvyqxn9lys2wgpkr",
			"reward": {
				"denom": "ucmdx",
				"amount": "2273"
			}
		},
		{
			"address": "comdex1yncfajtquwg3u96v3w3a6pzz5y4qrnyeh0cyry",
			"reward": {
				"denom": "ucmdx",
				"amount": "25051"
			}
		},
		{
			"address": "comdex1yn6fewlyd5n2x69e3eyvx6vg2j5rydtyf6zyrc",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1yna059jt7geugufhv0t9xdfayfj7m33cpny03c",
			"reward": {
				"denom": "ucmdx",
				"amount": "166"
			}
		},
		{
			"address": "comdex1y5qy0wr4vd6gtuwy3cafl5v6c8ff6y2russmp4",
			"reward": {
				"denom": "ucmdx",
				"amount": "57085"
			}
		},
		{
			"address": "comdex1y5p0938uhm0lxfzc8taqu9p4n5g9fcpqax96w8",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1y5zf4stv0lfecvf0d5lqdh9w4gu7z4stld3768",
			"reward": {
				"denom": "ucmdx",
				"amount": "3297"
			}
		},
		{
			"address": "comdex1y5zsxx43u5uqgnq0enljzjuut3ad9xlh56gs0g",
			"reward": {
				"denom": "ucmdx",
				"amount": "3922"
			}
		},
		{
			"address": "comdex1y5y7flv7sqaxamty4290hma4lzs4vhpgvqk8ep",
			"reward": {
				"denom": "ucmdx",
				"amount": "2042"
			}
		},
		{
			"address": "comdex1y59yvs7cpev6rft9x29uuqkjg3ekd8tfggts97",
			"reward": {
				"denom": "ucmdx",
				"amount": "6553"
			}
		},
		{
			"address": "comdex1y59nphv4zn0mlahwpyw3wjjk04vlumrvu96ypx",
			"reward": {
				"denom": "ucmdx",
				"amount": "1137"
			}
		},
		{
			"address": "comdex1y596vmpf00utexaw8s3tetdjrg97mqeczq27lx",
			"reward": {
				"denom": "ucmdx",
				"amount": "9678"
			}
		},
		{
			"address": "comdex1y588p0whdnmgna2wyhdwzhd6678ptsxd23dryg",
			"reward": {
				"denom": "ucmdx",
				"amount": "317"
			}
		},
		{
			"address": "comdex1y52qj2lhgtp0serwpcl358pfdgj8s58hk6n0py",
			"reward": {
				"denom": "ucmdx",
				"amount": "29768"
			}
		},
		{
			"address": "comdex1y5trdczjtupl96xqzr48tg6k4wp7f8enjx705r",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1y5ts0p3yas9lscru9m2xuhx97k4ylqcpcxz78z",
			"reward": {
				"denom": "ucmdx",
				"amount": "143"
			}
		},
		{
			"address": "comdex1y5vm9jkwklt9djjw2fcd5yzc3at28kcqv773hg",
			"reward": {
				"denom": "ucmdx",
				"amount": "6039"
			}
		},
		{
			"address": "comdex1y53fyd2nc3ctzkhuaphxy9lw7c3aefc3f3fshz",
			"reward": {
				"denom": "ucmdx",
				"amount": "3998"
			}
		},
		{
			"address": "comdex1y53t6lr5a33g4ku3x6apksqw45k8w6h3aqjtr9",
			"reward": {
				"denom": "ucmdx",
				"amount": "1762"
			}
		},
		{
			"address": "comdex1y54rjjlmyselkufxrrdcua3w3ukhzyhg2g9z9a",
			"reward": {
				"denom": "ucmdx",
				"amount": "26785"
			}
		},
		{
			"address": "comdex1y5hp8f08pjpwhwxl7hw5xpudw7jrx5fk2lqkz4",
			"reward": {
				"denom": "ucmdx",
				"amount": "16"
			}
		},
		{
			"address": "comdex1y5apl6xuljxs4mre8e890d55lngzrlsnneruy8",
			"reward": {
				"denom": "ucmdx",
				"amount": "10100"
			}
		},
		{
			"address": "comdex1y4qk37dtjqk3j73agm29kwxmwfaz3lat5nwpty",
			"reward": {
				"denom": "ucmdx",
				"amount": "625"
			}
		},
		{
			"address": "comdex1y4rzzrgl66eyhzt6gse2k7ej3zgwmngemngr9x",
			"reward": {
				"denom": "ucmdx",
				"amount": "4116"
			}
		},
		{
			"address": "comdex1y4rs2wsnply909j0zpjgzlludtg4ythay8guhj",
			"reward": {
				"denom": "ucmdx",
				"amount": "1524"
			}
		},
		{
			"address": "comdex1y4yrvx7qje6mh5lfypyrensrfhyv3t47g9dvuj",
			"reward": {
				"denom": "ucmdx",
				"amount": "15585"
			}
		},
		{
			"address": "comdex1y49xws42jp5gwfmh3slg3zyapn0jqsflkdczpv",
			"reward": {
				"denom": "ucmdx",
				"amount": "2041"
			}
		},
		{
			"address": "comdex1y4x9zm5usesmjen432pvd4jtuc47mrtal36tgm",
			"reward": {
				"denom": "ucmdx",
				"amount": "2790"
			}
		},
		{
			"address": "comdex1y48nh8jy6wj9m2z8mpakam7p6g9rz5ztfy2zuc",
			"reward": {
				"denom": "ucmdx",
				"amount": "14427"
			}
		},
		{
			"address": "comdex1y48ujlcsl038gshqzdrxx9ck32en47drced3u5",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1y4gm4vflfzfuv5s6h35m7tpz67erpwqgr3lwdh",
			"reward": {
				"denom": "ucmdx",
				"amount": "16821"
			}
		},
		{
			"address": "comdex1y4t7lfgwvalncnvfge0tcc7l08tfdmrln3v28t",
			"reward": {
				"denom": "ucmdx",
				"amount": "182"
			}
		},
		{
			"address": "comdex1y4wepr0q57al8zr2h9kwudcm5t4lw9ckaxkvsn",
			"reward": {
				"denom": "ucmdx",
				"amount": "4048"
			}
		},
		{
			"address": "comdex1y40meywe59k8kmlsap2dz9ws3dnjx957kfp0x2",
			"reward": {
				"denom": "ucmdx",
				"amount": "528"
			}
		},
		{
			"address": "comdex1y4sn4wl8nwdh7ntp52xct4phhyzwez8xgg6r76",
			"reward": {
				"denom": "ucmdx",
				"amount": "68728"
			}
		},
		{
			"address": "comdex1y43x2pdgj69vff52cagz7mcz092azxjgk07x9r",
			"reward": {
				"denom": "ucmdx",
				"amount": "7070"
			}
		},
		{
			"address": "comdex1y43ut5pl33aq78zvufqd999w9p22gw4ghcmrrj",
			"reward": {
				"denom": "ucmdx",
				"amount": "1252"
			}
		},
		{
			"address": "comdex1y4j2pzf24p29pr9zegxvq7pp8vw398wme580mv",
			"reward": {
				"denom": "ucmdx",
				"amount": "17657"
			}
		},
		{
			"address": "comdex1y4j3p8p94rtf5xsa3pez9htqhck63yxrft2rlt",
			"reward": {
				"denom": "ucmdx",
				"amount": "26431"
			}
		},
		{
			"address": "comdex1y4jjwr9d4w0dg8wve6fn9wfgrgr3armjt90nca",
			"reward": {
				"denom": "ucmdx",
				"amount": "14266"
			}
		},
		{
			"address": "comdex1y4nvllu3sucusqepys6elsl63xkhsw4lht3xzt",
			"reward": {
				"denom": "ucmdx",
				"amount": "54408"
			}
		},
		{
			"address": "comdex1y46y7rwcu5kcvcaqwrckayq9etg003760hjfan",
			"reward": {
				"denom": "ucmdx",
				"amount": "539"
			}
		},
		{
			"address": "comdex1y4mxd9rankzysaugn52ac7afljr22myef56fel",
			"reward": {
				"denom": "ucmdx",
				"amount": "62354"
			}
		},
		{
			"address": "comdex1y4ue39e28x90qvpykz9h0pr8qct24rx0xn2rzy",
			"reward": {
				"denom": "ucmdx",
				"amount": "5375"
			}
		},
		{
			"address": "comdex1y4acdjhfrweytvx90kynkd0rd4sw9tu3lx4kuq",
			"reward": {
				"denom": "ucmdx",
				"amount": "144"
			}
		},
		{
			"address": "comdex1y47jfdvnt4p3jt5y0wdavy3fmsx9zaqnnwn6np",
			"reward": {
				"denom": "ucmdx",
				"amount": "1"
			}
		},
		{
			"address": "comdex1y4l7ztx4xpmrrat79tngfvsr5jjays37cl9dyp",
			"reward": {
				"denom": "ucmdx",
				"amount": "1338"
			}
		},
		{
			"address": "comdex1ykqfsmrw0uk4fcnycpp60a0cdg0xlsmcgsmped",
			"reward": {
				"denom": "ucmdx",
				"amount": "17619"
			}
		},
		{
			"address": "comdex1ykqd7czfjfvwan777jq0uly0ll6gvwv7yxrsrx",
			"reward": {
				"denom": "ucmdx",
				"amount": "197636"
			}
		},
		{
			"address": "comdex1ykzg545k2kapf7fwmgmpvewlnfp88jhvflts6a",
			"reward": {
				"denom": "ucmdx",
				"amount": "7041"
			}
		},
		{
			"address": "comdex1ykyk43nxapm48n60ft0sr8ndrg8deu9prz2aey",
			"reward": {
				"denom": "ucmdx",
				"amount": "1974"
			}
		},
		{
			"address": "comdex1ykxtlchw93fcu0q8xkxss7glg2z2mqflmsv6ex",
			"reward": {
				"denom": "ucmdx",
				"amount": "61397"
			}
		},
		{
			"address": "comdex1yk8dsn62an3p8c67a047qv9lvqa7crpafek5vm",
			"reward": {
				"denom": "ucmdx",
				"amount": "481"
			}
		},
		{
			"address": "comdex1yk8anx9kam5ep6ajkwhvd9dy2wqrfkvuufeeh7",
			"reward": {
				"denom": "ucmdx",
				"amount": "15"
			}
		},
		{
			"address": "comdex1yk877d8cglr4c3rywg7ys7qcp4u48kwd98dsdd",
			"reward": {
				"denom": "ucmdx",
				"amount": "18165"
			}
		},
		{
			"address": "comdex1ykfqk7836ejwdpxazpq3suy2d26km8wfe6ttwy",
			"reward": {
				"denom": "ucmdx",
				"amount": "6432"
			}
		},
		{
			"address": "comdex1yktng538a2np8266xr74sd3k598plh7gwv2thh",
			"reward": {
				"denom": "ucmdx",
				"amount": "131644"
			}
		},
		{
			"address": "comdex1yktlxmklpf3rrlnm6jmcy5uxk3f4lh6wemlk0z",
			"reward": {
				"denom": "ucmdx",
				"amount": "3282"
			}
		},
		{
			"address": "comdex1ykvyhjaf4eltac89kypud7rc226ckutkx7vwtp",
			"reward": {
				"denom": "ucmdx",
				"amount": "1595"
			}
		},
		{
			"address": "comdex1ykdyn32m2st760tt5udu7ljskky0y7cxr073u6",
			"reward": {
				"denom": "ucmdx",
				"amount": "19"
			}
		},
		{
			"address": "comdex1yk3nm9n8kd94rr5rsczuu86nuxp9qt2aapj264",
			"reward": {
				"denom": "ucmdx",
				"amount": "1432"
			}
		},
		{
			"address": "comdex1yk3md3meme3ntlaa23kskdhjmz5taw9pknlsy4",
			"reward": {
				"denom": "ucmdx",
				"amount": "30995"
			}
		},
		{
			"address": "comdex1yknlqplw7sf7adsg9uxjprurx7mmqjyfqvugna",
			"reward": {
				"denom": "ucmdx",
				"amount": "25065"
			}
		},
		{
			"address": "comdex1ykkvmkpun5d6ce2t4s04pdkva5mzhmpxxn9whu",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1ykh92n7j2h930m8gf7v80m87mzxcljamcd7ldc",
			"reward": {
				"denom": "ucmdx",
				"amount": "71477"
			}
		},
		{
			"address": "comdex1ykhd5n0zj722d9dm0s0209s7a3y2g0vvzacw0w",
			"reward": {
				"denom": "ucmdx",
				"amount": "53287"
			}
		},
		{
			"address": "comdex1ykczm52yegslsy0wwpqj7mnrhm8r0dwqdd8wmj",
			"reward": {
				"denom": "ucmdx",
				"amount": "16667"
			}
		},
		{
			"address": "comdex1yke2q9whj5gcdc8u046ec6jxk8r0l9w0f5upgf",
			"reward": {
				"denom": "ucmdx",
				"amount": "2018"
			}
		},
		{
			"address": "comdex1ykeesr75yf4yga60chpt54a4ynz8rn2ftchkd7",
			"reward": {
				"denom": "ucmdx",
				"amount": "3001"
			}
		},
		{
			"address": "comdex1yk67s35ndjytl2vwxr56se8vq809nx48jl77vr",
			"reward": {
				"denom": "ucmdx",
				"amount": "275"
			}
		},
		{
			"address": "comdex1ykazj74lxsvsclswgss7p0a07ksvauqsskcd6k",
			"reward": {
				"denom": "ucmdx",
				"amount": "4373"
			}
		},
		{
			"address": "comdex1yhzvn7jq8j248cx5asnjzxwr79c20flsq0xqxx",
			"reward": {
				"denom": "ucmdx",
				"amount": "2159"
			}
		},
		{
			"address": "comdex1yhymjgzq7lnmkh974lmv6cy8ckl5gxry8ak9j8",
			"reward": {
				"denom": "ucmdx",
				"amount": "1762"
			}
		},
		{
			"address": "comdex1yh9p85xhnxmhtxz6tzde5rxquekcyhftthfrtj",
			"reward": {
				"denom": "ucmdx",
				"amount": "759"
			}
		},
		{
			"address": "comdex1yh8v0uyejum2rvw0glcla7hfvpf40gjsmrjd3g",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1yh8df7p3wlltqxfv4k7q49sgx84sa90h0xcvpr",
			"reward": {
				"denom": "ucmdx",
				"amount": "10605"
			}
		},
		{
			"address": "comdex1yhvr59ywfndwyncf7naa37a6np68a4xcfvce7w",
			"reward": {
				"denom": "ucmdx",
				"amount": "864"
			}
		},
		{
			"address": "comdex1yhvesf6xk5nefcry0wyxgxtjadxrcrgkq3a7ug",
			"reward": {
				"denom": "ucmdx",
				"amount": "1475"
			}
		},
		{
			"address": "comdex1yhwefhmxapfc92d0s4emxhq03fd070l4kkhek0",
			"reward": {
				"denom": "ucmdx",
				"amount": "597"
			}
		},
		{
			"address": "comdex1yh4wlek6mqywj7v4kuxfyw6lzd0dy0wxu9ay8f",
			"reward": {
				"denom": "ucmdx",
				"amount": "2876"
			}
		},
		{
			"address": "comdex1yhcwpk3waq6yj4clzzvgjyr9h6cdmc5vhhudrk",
			"reward": {
				"denom": "ucmdx",
				"amount": "17469"
			}
		},
		{
			"address": "comdex1yhc65cz47854k7hmj0yndmkpyuu3r7yefcc2j2",
			"reward": {
				"denom": "ucmdx",
				"amount": "189"
			}
		},
		{
			"address": "comdex1yhmqs9l6qrqz2638hnxvxglyjjv65cmz0648h5",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1yhu8hqhpl4w7nk8v5kyldqdhqkmlyuxvxgsare",
			"reward": {
				"denom": "ucmdx",
				"amount": "15"
			}
		},
		{
			"address": "comdex1yh7vedp8k0ulf4m9jswdyk8a6ywalzegg6759m",
			"reward": {
				"denom": "ucmdx",
				"amount": "525"
			}
		},
		{
			"address": "comdex1yh7h6axju3u6xsq8nu3nsv962juhgl2pkf5jkk",
			"reward": {
				"denom": "ucmdx",
				"amount": "612"
			}
		},
		{
			"address": "comdex1ycrqyzx4726nra6c73dhwth2svp5kujypt8py7",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1yc80f9fv44n2pkslcstc4dzegdq45rh03a9utk",
			"reward": {
				"denom": "ucmdx",
				"amount": "6225"
			}
		},
		{
			"address": "comdex1yc85rpju8rq0tn32hhaq8pwukzzrcxl9p3clya",
			"reward": {
				"denom": "ucmdx",
				"amount": "2023"
			}
		},
		{
			"address": "comdex1ycfh3lwhatggcxnu3kmylrhy4w2sr06fqj77cz",
			"reward": {
				"denom": "ucmdx",
				"amount": "92449"
			}
		},
		{
			"address": "comdex1yc23cxk2d3v7xf3ew0s6z2gg2tqa0z0haenhhn",
			"reward": {
				"denom": "ucmdx",
				"amount": "169"
			}
		},
		{
			"address": "comdex1ycvmjmx0ysqakan8lejlvd7zd4vsp3glr8d4uz",
			"reward": {
				"denom": "ucmdx",
				"amount": "1899"
			}
		},
		{
			"address": "comdex1ycwp584dujn2vngmurqwqqrey9hwxmyc5wfcmc",
			"reward": {
				"denom": "ucmdx",
				"amount": "1066"
			}
		},
		{
			"address": "comdex1ycwthrxmtrcew8cf30q6yj0gxgkw5msz4pruyw",
			"reward": {
				"denom": "ucmdx",
				"amount": "32232"
			}
		},
		{
			"address": "comdex1ycwc3mpf5wkdkje3tjthylpux7ffc2zqkcy2vv",
			"reward": {
				"denom": "ucmdx",
				"amount": "75559"
			}
		},
		{
			"address": "comdex1ycsm86qpuk3zv6820lflttynd90vpdcgj9j025",
			"reward": {
				"denom": "ucmdx",
				"amount": "336"
			}
		},
		{
			"address": "comdex1ycsmamp6hhealghtqt8rkxnmlw94aa7xlnyf6w",
			"reward": {
				"denom": "ucmdx",
				"amount": "26494"
			}
		},
		{
			"address": "comdex1yc4tulhndy4de55p5kkwza7sg5gpsz7yqfzu0p",
			"reward": {
				"denom": "ucmdx",
				"amount": "3017"
			}
		},
		{
			"address": "comdex1ycke0hm0q8zywzvl32chap88tfr6gcekpuxh6w",
			"reward": {
				"denom": "ucmdx",
				"amount": "508"
			}
		},
		{
			"address": "comdex1ycmshxs3fk05jztderp6gvap8zzsmpzqqj6nwd",
			"reward": {
				"denom": "ucmdx",
				"amount": "77320"
			}
		},
		{
			"address": "comdex1ycl5c4dmd8zmhxfx3udm08zk86nay5r9duec2f",
			"reward": {
				"denom": "ucmdx",
				"amount": "10207"
			}
		},
		{
			"address": "comdex1yezx3f3ad85xlpysmf7p5yqtnq6pmcezggpra6",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1yer9rccnxvcxslsl0zu97tyd6hvpyf22jvcw69",
			"reward": {
				"denom": "ucmdx",
				"amount": "18963"
			}
		},
		{
			"address": "comdex1yera5r8ancdhedeyqcumjwgpxe53au9568a2e5",
			"reward": {
				"denom": "ucmdx",
				"amount": "17106"
			}
		},
		{
			"address": "comdex1yeyafzl0nczpdwx2q2wkwkez7za7uwt2wjqz93",
			"reward": {
				"denom": "ucmdx",
				"amount": "3326578"
			}
		},
		{
			"address": "comdex1yex3lcnw6lry48gd7m6d9ys0q70dchne0ep7gj",
			"reward": {
				"denom": "ucmdx",
				"amount": "1491"
			}
		},
		{
			"address": "comdex1ye8ttaw42e4ykfxzgv848e4lkhlqnf8j9fzyf5",
			"reward": {
				"denom": "ucmdx",
				"amount": "14106"
			}
		},
		{
			"address": "comdex1ye2n5x2y3efytekcq0zs86tf0hjjzavxwcxlwc",
			"reward": {
				"denom": "ucmdx",
				"amount": "90"
			}
		},
		{
			"address": "comdex1yet6h7n48dxc3xz20l8s7wynslhcgkee0vpwct",
			"reward": {
				"denom": "ucmdx",
				"amount": "143"
			}
		},
		{
			"address": "comdex1yevnk9vftraykaastq4j6s3y9seusxl3lf72na",
			"reward": {
				"denom": "ucmdx",
				"amount": "8722"
			}
		},
		{
			"address": "comdex1ye0kh9z36dn0esfps070h6acg8s27kmqhc099v",
			"reward": {
				"denom": "ucmdx",
				"amount": "142"
			}
		},
		{
			"address": "comdex1ye0hjttuxx29zqk764sjqrhppvmcsaltej3cu5",
			"reward": {
				"denom": "ucmdx",
				"amount": "39"
			}
		},
		{
			"address": "comdex1yen5f0ej9njg0d9pa8nn2hwjpqqm36zjq627yy",
			"reward": {
				"denom": "ucmdx",
				"amount": "12414"
			}
		},
		{
			"address": "comdex1ye5dea0xs2ha2mvzxkggdk4junrqcswvgfdxls",
			"reward": {
				"denom": "ucmdx",
				"amount": "8416"
			}
		},
		{
			"address": "comdex1yekf2vgg84cj9nmtyz3efz9z4wwu62xq0sppue",
			"reward": {
				"denom": "ucmdx",
				"amount": "2059"
			}
		},
		{
			"address": "comdex1yek25khr0ktsagm47rcyryzea9f2y3qq4p98ax",
			"reward": {
				"denom": "ucmdx",
				"amount": "4385"
			}
		},
		{
			"address": "comdex1yeuxnj3yldscyxvdgzr8zckke0ztwxr2mxgu73",
			"reward": {
				"denom": "ucmdx",
				"amount": "3696"
			}
		},
		{
			"address": "comdex1yeu6chjmdys6w8r4qhcl33e2ppnenj4p67uzt3",
			"reward": {
				"denom": "ucmdx",
				"amount": "6527"
			}
		},
		{
			"address": "comdex1yearf9axnmne76y40ylsqu2cw0kmc6rfn440cv",
			"reward": {
				"denom": "ucmdx",
				"amount": "2954"
			}
		},
		{
			"address": "comdex1yeacdju8fsuy3u9dmr23w8d92ehz9fwz2kuv5x",
			"reward": {
				"denom": "ucmdx",
				"amount": "15140"
			}
		},
		{
			"address": "comdex1yelthfx0snxvn6gte4q980y7hkkwy8swghy4mm",
			"reward": {
				"denom": "ucmdx",
				"amount": "151"
			}
		},
		{
			"address": "comdex1y6p7jgsr63gm62lxj730j6vdd3ushgypk0cz2u",
			"reward": {
				"denom": "ucmdx",
				"amount": "48517"
			}
		},
		{
			"address": "comdex1y6y8kv9ahfxq3jqlupj2xw9gf3ce2myan0dh3k",
			"reward": {
				"denom": "ucmdx",
				"amount": "649"
			}
		},
		{
			"address": "comdex1y69hzkwecnu8ax9ta2qm4xv8h6n9cgdhfgc9sg",
			"reward": {
				"denom": "ucmdx",
				"amount": "1951"
			}
		},
		{
			"address": "comdex1y6xegcjplad68u5m94cs0gemgp9dgxl285rrcj",
			"reward": {
				"denom": "ucmdx",
				"amount": "2030"
			}
		},
		{
			"address": "comdex1y68svh2swz0ljkg3sp2nn8lksnwx72scp5rsa3",
			"reward": {
				"denom": "ucmdx",
				"amount": "11770"
			}
		},
		{
			"address": "comdex1y6260p0m7nykkxpmmlgw44n5vlagh77h97nyyk",
			"reward": {
				"denom": "ucmdx",
				"amount": "403"
			}
		},
		{
			"address": "comdex1y6vtgyeekjrxk5zm3zuud6kxfu50emmmk2n95t",
			"reward": {
				"denom": "ucmdx",
				"amount": "143"
			}
		},
		{
			"address": "comdex1y6vmmrz93xyhhv5qct80cd3n9mkvxzk5ud4ev7",
			"reward": {
				"denom": "ucmdx",
				"amount": "1283"
			}
		},
		{
			"address": "comdex1y6vauvc9kls0kjc2a2gs2q6h89fzryrpq848fs",
			"reward": {
				"denom": "ucmdx",
				"amount": "1752"
			}
		},
		{
			"address": "comdex1y6s4w36azugqkmgxs4jcy07quz9ruegsa4043y",
			"reward": {
				"denom": "ucmdx",
				"amount": "166"
			}
		},
		{
			"address": "comdex1y6n3hrnuh9dcfluzylsc5vey97za30r5cw4vvt",
			"reward": {
				"denom": "ucmdx",
				"amount": "38827"
			}
		},
		{
			"address": "comdex1y64x3l3pqjxjrvq2cyj48dsh66rgptaczdwqj0",
			"reward": {
				"denom": "ucmdx",
				"amount": "1190"
			}
		},
		{
			"address": "comdex1y6hm3h83ntcskqfxgvv330j7psurzz0vy5t3g0",
			"reward": {
				"denom": "ucmdx",
				"amount": "9005"
			}
		},
		{
			"address": "comdex1y6c86rjd6eqeq24hwpjs70l4wuhyuvswl4nqjj",
			"reward": {
				"denom": "ucmdx",
				"amount": "182"
			}
		},
		{
			"address": "comdex1y6603qh08ttze6vaqr3ucmclq4xxj88x4unxfz",
			"reward": {
				"denom": "ucmdx",
				"amount": "1492"
			}
		},
		{
			"address": "comdex1y67w03eqljlxmh6nny2r9zwcc9dzgv4al6mn79",
			"reward": {
				"denom": "ucmdx",
				"amount": "54313"
			}
		},
		{
			"address": "comdex1ymp6mtwdefclcmgk3whp6n6f5tcgrpruvsayx3",
			"reward": {
				"denom": "ucmdx",
				"amount": "1790"
			}
		},
		{
			"address": "comdex1ymp739xnnqgxsv8gp205zvvkrjfjvf27lf3d2e",
			"reward": {
				"denom": "ucmdx",
				"amount": "5450"
			}
		},
		{
			"address": "comdex1ymzks3udfcn5rg8n5wf9n47yw8jrz80fxhfdsk",
			"reward": {
				"denom": "ucmdx",
				"amount": "70669"
			}
		},
		{
			"address": "comdex1ym9guccydzgqhqrutyphe6g72f0as0nj46kkj8",
			"reward": {
				"denom": "ucmdx",
				"amount": "143363"
			}
		},
		{
			"address": "comdex1ym8mq092q22lf8nw0jlvug7kf5kx0h7x998xcw",
			"reward": {
				"denom": "ucmdx",
				"amount": "6871"
			}
		},
		{
			"address": "comdex1ymg877c9t0wrlxn3tudl0svdkg9w9600y0z4gh",
			"reward": {
				"denom": "ucmdx",
				"amount": "2041"
			}
		},
		{
			"address": "comdex1ymfug75xsv0ljh5ywkwxtan3l96ylf27f36nqg",
			"reward": {
				"denom": "ucmdx",
				"amount": "12623"
			}
		},
		{
			"address": "comdex1ymtsrltulp2ne7q3355xmm9xm96qmf3e5vqmr6",
			"reward": {
				"denom": "ucmdx",
				"amount": "890"
			}
		},
		{
			"address": "comdex1ymv3hncpk54vxtvkwz9gvq58dcfq9fk9gcq3mm",
			"reward": {
				"denom": "ucmdx",
				"amount": "9477"
			}
		},
		{
			"address": "comdex1ym3z4tvjkx54jaz2spgt0a59akzwtt8hhr38r3",
			"reward": {
				"denom": "ucmdx",
				"amount": "3508"
			}
		},
		{
			"address": "comdex1ym3wguk83xzrlj58kwnwv4wvqrclan34c5jc5f",
			"reward": {
				"denom": "ucmdx",
				"amount": "123503"
			}
		},
		{
			"address": "comdex1ymnpjqj7lcx50n8j307mwxv4xzzr2trxgfr7uq",
			"reward": {
				"denom": "ucmdx",
				"amount": "59451"
			}
		},
		{
			"address": "comdex1ym445tjd0tjqmgcmq6pj6ewf8usrsuyj226cdm",
			"reward": {
				"denom": "ucmdx",
				"amount": "712"
			}
		},
		{
			"address": "comdex1ym47l7zvp286f949ylpsqkzz6e4vv6357zzpfp",
			"reward": {
				"denom": "ucmdx",
				"amount": "1"
			}
		},
		{
			"address": "comdex1ymc4d62tq5ujglvnm23vn84qf5007h8f0gtmf7",
			"reward": {
				"denom": "ucmdx",
				"amount": "284"
			}
		},
		{
			"address": "comdex1ymepwnyrl9yjsy4fgwj59xdnh2gutn0nk4mfet",
			"reward": {
				"denom": "ucmdx",
				"amount": "5943"
			}
		},
		{
			"address": "comdex1ymeugaq0evp02v3mk5rgauuwu5c86ay3430k47",
			"reward": {
				"denom": "ucmdx",
				"amount": "1523"
			}
		},
		{
			"address": "comdex1ymm8jzym5a5adarzwgsqkuge9422gs27ahltxt",
			"reward": {
				"denom": "ucmdx",
				"amount": "1460"
			}
		},
		{
			"address": "comdex1ymau7uztmww04awddql2a5wym5m9x57e7096p6",
			"reward": {
				"denom": "ucmdx",
				"amount": "7166"
			}
		},
		{
			"address": "comdex1yuqyqyyvemh8zmler9la3fntvh3dqw762gje6l",
			"reward": {
				"denom": "ucmdx",
				"amount": "68759"
			}
		},
		{
			"address": "comdex1yuzy0tnj7w33csrh8dflt682dlefgt09uv5e5u",
			"reward": {
				"denom": "ucmdx",
				"amount": "13449"
			}
		},
		{
			"address": "comdex1yuz6ak80k9859ecq4h2tgr98e8jyues047yqw8",
			"reward": {
				"denom": "ucmdx",
				"amount": "3470"
			}
		},
		{
			"address": "comdex1yurfkr4z9z4qs7tejflqgr674k2ed8pu6my9qe",
			"reward": {
				"denom": "ucmdx",
				"amount": "126"
			}
		},
		{
			"address": "comdex1yuy5lz22lx8lvt6gn9fxrczj5g2s5z6amk2cpq",
			"reward": {
				"denom": "ucmdx",
				"amount": "12020"
			}
		},
		{
			"address": "comdex1yuycef80txefx4fu4gc7kwzvql0ckuyht7n3x5",
			"reward": {
				"denom": "ucmdx",
				"amount": "14319"
			}
		},
		{
			"address": "comdex1yu9njx2t8w8nxk97qtz0t8jx3nq8pgg78xr9wn",
			"reward": {
				"denom": "ucmdx",
				"amount": "18074"
			}
		},
		{
			"address": "comdex1yu8xykg6h7vdntj58t8rh89qfylq3clzt2zy8u",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1yu8td2qa49u0qac0l345efh0r4u68uhl05yrre",
			"reward": {
				"denom": "ucmdx",
				"amount": "1411"
			}
		},
		{
			"address": "comdex1yu8t4udttksnrq6w6qqyuvjesu92qqkc7hd7qq",
			"reward": {
				"denom": "ucmdx",
				"amount": "148"
			}
		},
		{
			"address": "comdex1yug5gfjyxuqa4kk937232qvku0f5urpujjpe20",
			"reward": {
				"denom": "ucmdx",
				"amount": "9494"
			}
		},
		{
			"address": "comdex1yugcyr0tfh5rtvjwgsps2hh4rdujepfnaqs2u8",
			"reward": {
				"denom": "ucmdx",
				"amount": "2416"
			}
		},
		{
			"address": "comdex1yutt3n3ja0d5m86tlcg8zszae70hggtmyp4hw4",
			"reward": {
				"denom": "ucmdx",
				"amount": "36"
			}
		},
		{
			"address": "comdex1yuvrsgknm2pfzastdn9enykdgykepzt2u95vz3",
			"reward": {
				"denom": "ucmdx",
				"amount": "317"
			}
		},
		{
			"address": "comdex1yuvezz5ldrz44w9g3g8xaa4rjxvx6d0xy398rp",
			"reward": {
				"denom": "ucmdx",
				"amount": "258"
			}
		},
		{
			"address": "comdex1yusrw0tzsam6wk5r8003g33usswfpyzn0zz5rh",
			"reward": {
				"denom": "ucmdx",
				"amount": "6520"
			}
		},
		{
			"address": "comdex1yuswg8jck57cak9f2xn8gtsfxzn06s0wa2yzkr",
			"reward": {
				"denom": "ucmdx",
				"amount": "1262"
			}
		},
		{
			"address": "comdex1yu39qgjc7gcczguggknxlpqh3cy2lks2xyce7f",
			"reward": {
				"denom": "ucmdx",
				"amount": "410"
			}
		},
		{
			"address": "comdex1yu3u7xguyn79xp4h08wc9za6k5zes3dzl0wt4m",
			"reward": {
				"denom": "ucmdx",
				"amount": "129477"
			}
		},
		{
			"address": "comdex1yujaep7vnr6j3au3zsr7gnvk7gj0z6f5esqhp7",
			"reward": {
				"denom": "ucmdx",
				"amount": "1401"
			}
		},
		{
			"address": "comdex1yu50qwvcktzklclfj25ukndrxwngjet9s34mrn",
			"reward": {
				"denom": "ucmdx",
				"amount": "7085"
			}
		},
		{
			"address": "comdex1yu4gsgzhylmrsjvy6jvys9les7z66n8wlxcdtg",
			"reward": {
				"denom": "ucmdx",
				"amount": "537"
			}
		},
		{
			"address": "comdex1yuke6answsw9rtqapa3y776z80uupma6c44fdc",
			"reward": {
				"denom": "ucmdx",
				"amount": "1041"
			}
		},
		{
			"address": "comdex1yua2z68wxggv3nsq8zqevpnr4ar63qq6cy9hw2",
			"reward": {
				"denom": "ucmdx",
				"amount": "1697"
			}
		},
		{
			"address": "comdex1yu7qk5z54vd5p62dzfvtlhssdfasxhfmqeld2y",
			"reward": {
				"denom": "ucmdx",
				"amount": "17724"
			}
		},
		{
			"address": "comdex1yu7pg0hnc878z2qe0hc95dm2mprru2pwyk0r6a",
			"reward": {
				"denom": "ucmdx",
				"amount": "10106"
			}
		},
		{
			"address": "comdex1yadcdtkwk9j0a6gskm6h2pvnlx9qjfukn7dejq",
			"reward": {
				"denom": "ucmdx",
				"amount": "27359"
			}
		},
		{
			"address": "comdex1yak8j4ssp6vq65ue57cfd2w6464w70kh0fyslt",
			"reward": {
				"denom": "ucmdx",
				"amount": "38598"
			}
		},
		{
			"address": "comdex1yam90ez4dsmw3a08c4g49v5lwrhlyfvnrv6fcq",
			"reward": {
				"denom": "ucmdx",
				"amount": "196"
			}
		},
		{
			"address": "comdex1yaahw59q2kly8znvw4w87gypl4e4j2jhkuugg2",
			"reward": {
				"denom": "ucmdx",
				"amount": "28694"
			}
		},
		{
			"address": "comdex1ya7yg6wupzfgynduj0ch5vcvc56asuln8hrsx0",
			"reward": {
				"denom": "ucmdx",
				"amount": "210"
			}
		},
		{
			"address": "comdex1y7zrmlz96dnwrdxelqccwp0ffvzz9v6e39h5k5",
			"reward": {
				"denom": "ucmdx",
				"amount": "1784"
			}
		},
		{
			"address": "comdex1y78al3x5xu5qc0ycxz2m8wplx4x84tzln0cl49",
			"reward": {
				"denom": "ucmdx",
				"amount": "5894"
			}
		},
		{
			"address": "comdex1y78l0cnun58rqlkj60w2wvf7a6lsc3extwa7jz",
			"reward": {
				"denom": "ucmdx",
				"amount": "486"
			}
		},
		{
			"address": "comdex1y7tdejjz7jf2yd6jdgueaeh7vc4g36pjl53xt3",
			"reward": {
				"denom": "ucmdx",
				"amount": "6475"
			}
		},
		{
			"address": "comdex1y7tc6g66ewgkux44u5htnyje59p4mnj2at6sls",
			"reward": {
				"denom": "ucmdx",
				"amount": "304978"
			}
		},
		{
			"address": "comdex1y7vcfhj9g7t3aeg50hlt42m4jhzy6uxjgvrzq8",
			"reward": {
				"denom": "ucmdx",
				"amount": "166"
			}
		},
		{
			"address": "comdex1y7s6xcqs5dul55rhfygrlq029jea0pfsj7fga2",
			"reward": {
				"denom": "ucmdx",
				"amount": "14529"
			}
		},
		{
			"address": "comdex1y7jpq65f86n5fgrnfevasvmg9ewj5tc7ac849g",
			"reward": {
				"denom": "ucmdx",
				"amount": "3369"
			}
		},
		{
			"address": "comdex1y7kxurnq74w2huk6usfq6zcmdcz752km9jajjd",
			"reward": {
				"denom": "ucmdx",
				"amount": "1758"
			}
		},
		{
			"address": "comdex1y7k3h87awsxqvps7e6d7f2nvk8mvsxwp2h2muw",
			"reward": {
				"denom": "ucmdx",
				"amount": "34164"
			}
		},
		{
			"address": "comdex1y7hkq5zcggfm9wyzw3me307hkc4ld3tx3j4qn6",
			"reward": {
				"denom": "ucmdx",
				"amount": "1248"
			}
		},
		{
			"address": "comdex1y7h673t8uvtu40jyqhr5zv0qffqgv7tchzx9pt",
			"reward": {
				"denom": "ucmdx",
				"amount": "1424"
			}
		},
		{
			"address": "comdex1y7cv7psfvled3xnrtsraedeke2p39lcguwvxxt",
			"reward": {
				"denom": "ucmdx",
				"amount": "5696"
			}
		},
		{
			"address": "comdex1y7ez5s7fxnajaw5t8u4fu9ut8lpn757mtme5c9",
			"reward": {
				"denom": "ucmdx",
				"amount": "8"
			}
		},
		{
			"address": "comdex1y7mh5e8aw0vfqsm9j6x2c9e776l52uegpk4lff",
			"reward": {
				"denom": "ucmdx",
				"amount": "715"
			}
		},
		{
			"address": "comdex1y7l6th236vu8hwyd76dyzg4gl07qvtlcy7kwm6",
			"reward": {
				"denom": "ucmdx",
				"amount": "4361"
			}
		},
		{
			"address": "comdex1ylzxzuac7cxp75rxhtv5t672gsewh4fy3k289h",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1ylxg2fn0h6x7vua0n3kz3w7q8aqrl4z93sn0hp",
			"reward": {
				"denom": "ucmdx",
				"amount": "179"
			}
		},
		{
			"address": "comdex1ylx03vqp2zkpjayudpx2k0psgdx4uc3d38mlx5",
			"reward": {
				"denom": "ucmdx",
				"amount": "5300"
			}
		},
		{
			"address": "comdex1ylf0n7tzlfrd8u8798w9m7075p462p9qpe0s2n",
			"reward": {
				"denom": "ucmdx",
				"amount": "1427"
			}
		},
		{
			"address": "comdex1ylfkjdcmfryy2023jp4hpnrsxhcwp9j0dkxzvm",
			"reward": {
				"denom": "ucmdx",
				"amount": "725"
			}
		},
		{
			"address": "comdex1yl2p7f8glqmwq2gpz4tlu6x02vxjksvz20k62v",
			"reward": {
				"denom": "ucmdx",
				"amount": "3348"
			}
		},
		{
			"address": "comdex1ylvt7wuj05m5wml0822n2p9tgczxn8g330t4m6",
			"reward": {
				"denom": "ucmdx",
				"amount": "12169"
			}
		},
		{
			"address": "comdex1yldcae6uy5q55k55v0524cfw0eaj8cta576uku",
			"reward": {
				"denom": "ucmdx",
				"amount": "3254"
			}
		},
		{
			"address": "comdex1yldaeduhnx9uru9de2rw8ln068tmz5en7a0w4h",
			"reward": {
				"denom": "ucmdx",
				"amount": "4241"
			}
		},
		{
			"address": "comdex1ylwgt5kkwqse5kc4tg53z68urjz3jsd0khr32l",
			"reward": {
				"denom": "ucmdx",
				"amount": "1009"
			}
		},
		{
			"address": "comdex1yljf0upv9zcdv62wdshjey6lx8357rjkjylaru",
			"reward": {
				"denom": "ucmdx",
				"amount": "985"
			}
		},
		{
			"address": "comdex1yl475lm0ehtsndah3wcdp7xqvwslpk4n62c9ya",
			"reward": {
				"denom": "ucmdx",
				"amount": "73"
			}
		},
		{
			"address": "comdex1ylh9dw289624ghtkrsgc73759z8rgqnvpcpa77",
			"reward": {
				"denom": "ucmdx",
				"amount": "32"
			}
		},
		{
			"address": "comdex1ylhd8q24au3pum0daq4el4y20scxdz29ufah4t",
			"reward": {
				"denom": "ucmdx",
				"amount": "9704"
			}
		},
		{
			"address": "comdex1ylh6rtney8hljckmjtcup05upuxctpv3fjnau8",
			"reward": {
				"denom": "ucmdx",
				"amount": "22026"
			}
		},
		{
			"address": "comdex1ylc88uqww6j7xrycj09d9mtas9ykjy94wuq4as",
			"reward": {
				"denom": "ucmdx",
				"amount": "142"
			}
		},
		{
			"address": "comdex1yll0ylmltrltcw09qn7jt42hqjke99yucfcucg",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex19qqza66v3f3srgyn6nfyjqxgzrd04fpnqduv76",
			"reward": {
				"denom": "ucmdx",
				"amount": "978"
			}
		},
		{
			"address": "comdex19qypu6hpwpyqmd24xufq8zkua3fmrhapvla06w",
			"reward": {
				"denom": "ucmdx",
				"amount": "326"
			}
		},
		{
			"address": "comdex19qy6gu2ata49fqp6l8nl27p07z24nk6eccnl45",
			"reward": {
				"denom": "ucmdx",
				"amount": "23343"
			}
		},
		{
			"address": "comdex19qfqj8g5ch5g4eydl5dtq73legq8mja4yjk9h5",
			"reward": {
				"denom": "ucmdx",
				"amount": "9492"
			}
		},
		{
			"address": "comdex19qfw6xkejzeswz6h8gfhf50waz0f2kcgc24q40",
			"reward": {
				"denom": "ucmdx",
				"amount": "10257"
			}
		},
		{
			"address": "comdex19q2nyhnrdhw7h5s05kafv4mrgwrerexzdn39h4",
			"reward": {
				"denom": "ucmdx",
				"amount": "16"
			}
		},
		{
			"address": "comdex19q2h8237906pk5w3e2l7lauej3j4uzuh4lzw28",
			"reward": {
				"denom": "ucmdx",
				"amount": "1519"
			}
		},
		{
			"address": "comdex19q2mc4ajvaah2dv7zssusywggheyuhhse6agcz",
			"reward": {
				"denom": "ucmdx",
				"amount": "103"
			}
		},
		{
			"address": "comdex19qv4r542te30wnq24x4zrgdauu0ux4p27zauaz",
			"reward": {
				"denom": "ucmdx",
				"amount": "6515"
			}
		},
		{
			"address": "comdex19qdxkzphe0zw6c3eqnf3w938zfe3k0jsk67lyd",
			"reward": {
				"denom": "ucmdx",
				"amount": "5867"
			}
		},
		{
			"address": "comdex19qdtt7mpfurn9rhyy0mwgpkpgnvct7lmky2s40",
			"reward": {
				"denom": "ucmdx",
				"amount": "88927"
			}
		},
		{
			"address": "comdex19qs5uceyzemm2tmw55w53ydm6963hdscds2ph4",
			"reward": {
				"denom": "ucmdx",
				"amount": "1420"
			}
		},
		{
			"address": "comdex19qskdhawwk3s2l043mh9lyrwx9v96tksqt3ygq",
			"reward": {
				"denom": "ucmdx",
				"amount": "51070"
			}
		},
		{
			"address": "comdex19qs798zha4r6xc38dnpctc0hryrj8g8rdmy6wm",
			"reward": {
				"denom": "ucmdx",
				"amount": "8784"
			}
		},
		{
			"address": "comdex19q32hxdexs77qt7jrzz4977ryaxqypm9jcfpdq",
			"reward": {
				"denom": "ucmdx",
				"amount": "924"
			}
		},
		{
			"address": "comdex19q3slfkwmr939ctjd7ejagnfxtrw46yeujnryq",
			"reward": {
				"denom": "ucmdx",
				"amount": "656"
			}
		},
		{
			"address": "comdex19qng6fmn9rk6q36e50t9s8j3rrwhadauvny7y5",
			"reward": {
				"denom": "ucmdx",
				"amount": "4958"
			}
		},
		{
			"address": "comdex19qhjd5hfh3865sdphuethct6jf4gamfeddemak",
			"reward": {
				"denom": "ucmdx",
				"amount": "192"
			}
		},
		{
			"address": "comdex19qc25l2nwztml9nvhqpg7y55420uvg59rnphnm",
			"reward": {
				"denom": "ucmdx",
				"amount": "4030"
			}
		},
		{
			"address": "comdex19qcsgsj9m7znk54x33a9m0p79dtym93s79ke2n",
			"reward": {
				"denom": "ucmdx",
				"amount": "615"
			}
		},
		{
			"address": "comdex19qe6s5m0a3cc726pyzwfwzrfx5tkwvd8whjyxr",
			"reward": {
				"denom": "ucmdx",
				"amount": "3049"
			}
		},
		{
			"address": "comdex19q6nwnhr05sfmh5j0wlfp5xc9gf4rnjax3hyrm",
			"reward": {
				"denom": "ucmdx",
				"amount": "251"
			}
		},
		{
			"address": "comdex19q6aqzkv3kwpc5wcxkgdu32wnsgkvtwqu7q4lz",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex19qu8nwj9kfytuwj60p2ad6cm3mzfdgkk4ad5rv",
			"reward": {
				"denom": "ucmdx",
				"amount": "6658"
			}
		},
		{
			"address": "comdex19ppk89xzyfel7y9ym43srgaemsjmm8j7aajec4",
			"reward": {
				"denom": "ucmdx",
				"amount": "1440"
			}
		},
		{
			"address": "comdex19pz88924pu0tjp980pmuv5hvu08p64mqqcluus",
			"reward": {
				"denom": "ucmdx",
				"amount": "726"
			}
		},
		{
			"address": "comdex19prs7rglnl0mques2e462wd9ylev86eclc7km9",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex19pr49x2he3kfa0tvgdcjugrgs6h8zt4gjr9vkl",
			"reward": {
				"denom": "ucmdx",
				"amount": "17857"
			}
		},
		{
			"address": "comdex19pys8f5krknrxa35dk25rw7fyu66suv5nr6kul",
			"reward": {
				"denom": "ucmdx",
				"amount": "12586"
			}
		},
		{
			"address": "comdex19p9tw2xcelkkcrllfeh7ed5g9pcdqvtu4h959c",
			"reward": {
				"denom": "ucmdx",
				"amount": "774"
			}
		},
		{
			"address": "comdex19p9wxmsz5gcec205hg5fjnr2rwagf2g044my78",
			"reward": {
				"denom": "ucmdx",
				"amount": "1389"
			}
		},
		{
			"address": "comdex19p8wunnpxq3srpzy64l6er56nyr9qalwe6e8el",
			"reward": {
				"denom": "ucmdx",
				"amount": "3011"
			}
		},
		{
			"address": "comdex19pw0r50uxwahtd2z9dvazvldjg3kfjpynt8ht3",
			"reward": {
				"denom": "ucmdx",
				"amount": "1751"
			}
		},
		{
			"address": "comdex19p0q0a0pz8s6649wtq34pjzwj2v3n908dnvu5y",
			"reward": {
				"denom": "ucmdx",
				"amount": "991"
			}
		},
		{
			"address": "comdex19p3uzdflshsag86g3z7fqamm5e4hz67zqugz4f",
			"reward": {
				"denom": "ucmdx",
				"amount": "7460"
			}
		},
		{
			"address": "comdex19phpk6eg5ehaengmnqnew04temdq645sngmfz0",
			"reward": {
				"denom": "ucmdx",
				"amount": "5455"
			}
		},
		{
			"address": "comdex19putahtcfue5pgcmnldc5jhes06tc6xffdfzsk",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex19pa99uc974p5055sqezy2uk3prstwvrpn3mwwd",
			"reward": {
				"denom": "ucmdx",
				"amount": "2612"
			}
		},
		{
			"address": "comdex19pa0wr268zanxcpym9djyvmws307a8s0cmrrnm",
			"reward": {
				"denom": "ucmdx",
				"amount": "2891"
			}
		},
		{
			"address": "comdex19z95u4839qydhcu7uptg0lhen0q4y2j58m8j3t",
			"reward": {
				"denom": "ucmdx",
				"amount": "3377"
			}
		},
		{
			"address": "comdex19z9maw73f489tk8rt63x2txlxg56lxhfkrjjpl",
			"reward": {
				"denom": "ucmdx",
				"amount": "3338"
			}
		},
		{
			"address": "comdex19zf93r5uf63x9anlpfzagn0xu35vtfjkwl6th6",
			"reward": {
				"denom": "ucmdx",
				"amount": "10537"
			}
		},
		{
			"address": "comdex19zfswdtfhal3h639r2kw0pwar405g4g5dn67fc",
			"reward": {
				"denom": "ucmdx",
				"amount": "6983"
			}
		},
		{
			"address": "comdex19ztrw2k2gnzushs45fra0zteztm4dlu3u5x30n",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex19zvajszwzuyh0dcqtd92jkva7n9qlr0kq7vt9n",
			"reward": {
				"denom": "ucmdx",
				"amount": "417"
			}
		},
		{
			"address": "comdex19zjcxqw9347dx869nnu6nj4g08q22uaywx8784",
			"reward": {
				"denom": "ucmdx",
				"amount": "79207"
			}
		},
		{
			"address": "comdex19zj7tk7guj48sejnua86k9zam9kmkvszp570nc",
			"reward": {
				"denom": "ucmdx",
				"amount": "42"
			}
		},
		{
			"address": "comdex19z5xkvnxmcxxv9nl94hs0ulhaxazlu3hh5f3lx",
			"reward": {
				"denom": "ucmdx",
				"amount": "20273"
			}
		},
		{
			"address": "comdex19z5jztky4dgr952k659f39gyrdeka834e2uuax",
			"reward": {
				"denom": "ucmdx",
				"amount": "1749"
			}
		},
		{
			"address": "comdex19zhltzh3rsd94p7szge4sp7zts82hlth6r2h8n",
			"reward": {
				"denom": "ucmdx",
				"amount": "4251"
			}
		},
		{
			"address": "comdex19zcad9yuyqpp84qr29ujnqph8mrkrgj03enrag",
			"reward": {
				"denom": "ucmdx",
				"amount": "1803"
			}
		},
		{
			"address": "comdex19za9rcrsrtx5skk7ap47ajftlvc47mrn5przvv",
			"reward": {
				"denom": "ucmdx",
				"amount": "4093"
			}
		},
		{
			"address": "comdex19zl79prfjrv9x532g750pjzhfpjgqvyf5srrnk",
			"reward": {
				"denom": "ucmdx",
				"amount": "4032"
			}
		},
		{
			"address": "comdex19rq5w8svtw7jq2elcjfzlw8lghpjlxapn603fq",
			"reward": {
				"denom": "ucmdx",
				"amount": "12440"
			}
		},
		{
			"address": "comdex19rzhw7yka25vf3mdpcdkerzc93n8ds4chtenft",
			"reward": {
				"denom": "ucmdx",
				"amount": "6560"
			}
		},
		{
			"address": "comdex19rzc4nyq44356wruv0wk820hk0j0ftdff7gde2",
			"reward": {
				"denom": "ucmdx",
				"amount": "1190"
			}
		},
		{
			"address": "comdex19ryx4h40yrr602pg32ekg6704l3nzqhjwsgz4y",
			"reward": {
				"denom": "ucmdx",
				"amount": "146"
			}
		},
		{
			"address": "comdex19rxzn30zvfvaxgmzlk4qntypkgs9uffetefnqg",
			"reward": {
				"denom": "ucmdx",
				"amount": "12175"
			}
		},
		{
			"address": "comdex19rgysq3h2f4e34yt220x5cu9yqyjkv6svzdqk9",
			"reward": {
				"denom": "ucmdx",
				"amount": "103"
			}
		},
		{
			"address": "comdex19rfqwv3equna5eyypmhrw7yvc0qz0y9mkgyt3c",
			"reward": {
				"denom": "ucmdx",
				"amount": "69293"
			}
		},
		{
			"address": "comdex19r25ky6nud0mucs3m8zvrprseal89k26fms3j2",
			"reward": {
				"denom": "ucmdx",
				"amount": "215"
			}
		},
		{
			"address": "comdex19rv6wh22f55kng4x2u0sl0rgqws5xew9fnjl0c",
			"reward": {
				"denom": "ucmdx",
				"amount": "595"
			}
		},
		{
			"address": "comdex19rsp6jjxssp4v0r360tc3ahkraclvyq3v2k6c8",
			"reward": {
				"denom": "ucmdx",
				"amount": "57"
			}
		},
		{
			"address": "comdex19rswkpgq7jfv0uysvem4sqylxvtue4kykneqzm",
			"reward": {
				"denom": "ucmdx",
				"amount": "2177"
			}
		},
		{
			"address": "comdex19rnz79mcxl33qg9rzv0tqlcefmq9m6yvltvqx7",
			"reward": {
				"denom": "ucmdx",
				"amount": "480"
			}
		},
		{
			"address": "comdex19rn52vj06h7mykzj9nue30ye3hcqe24p0wt4ue",
			"reward": {
				"denom": "ucmdx",
				"amount": "89701"
			}
		},
		{
			"address": "comdex19rhhn2f9dfugltk2r7zsf76e39925dat8u0vms",
			"reward": {
				"denom": "ucmdx",
				"amount": "16667"
			}
		},
		{
			"address": "comdex19rhhktgdp9g4mauk3q5p9rdu79xemf9qjn8fhw",
			"reward": {
				"denom": "ucmdx",
				"amount": "184"
			}
		},
		{
			"address": "comdex19rch2exfcg72vtt8p7a4t9j5lglwwspevweky6",
			"reward": {
				"denom": "ucmdx",
				"amount": "1"
			}
		},
		{
			"address": "comdex19rcusumrfwe25jm39t92sc0q0nljmwyu8xh55w",
			"reward": {
				"denom": "ucmdx",
				"amount": "8784"
			}
		},
		{
			"address": "comdex19r645ycpjw9jzat8jaz8jehw90q7lr8wlt4frg",
			"reward": {
				"denom": "ucmdx",
				"amount": "2738"
			}
		},
		{
			"address": "comdex19rm0sfgzhtregmvl42lt6el9fq7pvuaxal628j",
			"reward": {
				"denom": "ucmdx",
				"amount": "1931"
			}
		},
		{
			"address": "comdex19ru2g73f2hvjk343rmpp02g34zrznc5vah0lf9",
			"reward": {
				"denom": "ucmdx",
				"amount": "323"
			}
		},
		{
			"address": "comdex19ypepz7f5rgmf2rcng2tt75yny03umvegzd9aq",
			"reward": {
				"denom": "ucmdx",
				"amount": "2277"
			}
		},
		{
			"address": "comdex19yrdsqjjsg78vmqaqd84q2wf84vj2n8vx2w05d",
			"reward": {
				"denom": "ucmdx",
				"amount": "1992"
			}
		},
		{
			"address": "comdex19yrexd5s8j9eg8lvxtah2y04np8a6lcv45vulc",
			"reward": {
				"denom": "ucmdx",
				"amount": "1495"
			}
		},
		{
			"address": "comdex19y9j8exv64vrvzg77nupjz4w3eqk7qyd8kgec3",
			"reward": {
				"denom": "ucmdx",
				"amount": "389"
			}
		},
		{
			"address": "comdex19yg70ssswww6wwlyt2z95wq4r7jj65j86edq3v",
			"reward": {
				"denom": "ucmdx",
				"amount": "99754"
			}
		},
		{
			"address": "comdex19yf3w9qhw0cfgaq658twnx9zzt48zdqpganwnr",
			"reward": {
				"denom": "ucmdx",
				"amount": "36848"
			}
		},
		{
			"address": "comdex19y2rd84wkgz549rse3v75m24433c8lljhlh0h7",
			"reward": {
				"denom": "ucmdx",
				"amount": "143"
			}
		},
		{
			"address": "comdex19yvvl0cz068qu4sg7v2pxln0cqm8nt77pesvze",
			"reward": {
				"denom": "ucmdx",
				"amount": "122"
			}
		},
		{
			"address": "comdex19yjs2wlfmspz06e9u3wca8hg4kzmznt0aeja3q",
			"reward": {
				"denom": "ucmdx",
				"amount": "171"
			}
		},
		{
			"address": "comdex19y4qt2l4apntf2pgr86j4usc2ude6yrec3n4at",
			"reward": {
				"denom": "ucmdx",
				"amount": "150"
			}
		},
		{
			"address": "comdex19ycu6zt5dlm3yw3207f36pfjs0qg9wv9888amy",
			"reward": {
				"denom": "ucmdx",
				"amount": "365"
			}
		},
		{
			"address": "comdex19ymugy3xr7eql5fq03r60e8vt9s58pdm5q204d",
			"reward": {
				"denom": "ucmdx",
				"amount": "6775"
			}
		},
		{
			"address": "comdex19yu3y9gwj5avfkx9c0mxg2p8yr3hm55q8ve4tx",
			"reward": {
				"denom": "ucmdx",
				"amount": "40862"
			}
		},
		{
			"address": "comdex199x2qrqzchy9rkv4nc2l3xwld9azz29lfnpl08",
			"reward": {
				"denom": "ucmdx",
				"amount": "9678"
			}
		},
		{
			"address": "comdex199gxe69xgw8qxpnld7ya8eh92ajdaqzwu690et",
			"reward": {
				"denom": "ucmdx",
				"amount": "43619"
			}
		},
		{
			"address": "comdex199tlvkk93jyleervlvf6em3udxcy8gvxh27vfh",
			"reward": {
				"denom": "ucmdx",
				"amount": "6540"
			}
		},
		{
			"address": "comdex199ljhagf55twqa2weywrqacj0wsaax9esdvgnw",
			"reward": {
				"denom": "ucmdx",
				"amount": "419"
			}
		},
		{
			"address": "comdex19xpdvmuj4u5p8mczdyf85tsdmpjs0ezkkcmdez",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex19xzeq8lglatf8eglx3q4sxg4zhqgm8k5ydug4h",
			"reward": {
				"denom": "ucmdx",
				"amount": "61869"
			}
		},
		{
			"address": "comdex19xwgxltsnj0h5xvce0sp0wvkfmlcagjh9k20az",
			"reward": {
				"denom": "ucmdx",
				"amount": "185"
			}
		},
		{
			"address": "comdex19xjv76dg85p7gy5v65qw9ap8c7tjf8qhquq5s2",
			"reward": {
				"denom": "ucmdx",
				"amount": "1262"
			}
		},
		{
			"address": "comdex19xkqaw6lwk5uuq8wte829wa7mmnrhdaelfpnak",
			"reward": {
				"denom": "ucmdx",
				"amount": "125754"
			}
		},
		{
			"address": "comdex19xc9vce0mh4pw02l6tdc93d4k44we7hg8enpc3",
			"reward": {
				"denom": "ucmdx",
				"amount": "7965"
			}
		},
		{
			"address": "comdex19xus8kr64rptmtrgc7557tkmmya0qq6uk94c03",
			"reward": {
				"denom": "ucmdx",
				"amount": "2012"
			}
		},
		{
			"address": "comdex19x7ud2veye04aqxhfv0sql08lh0pe4h9w0a78w",
			"reward": {
				"denom": "ucmdx",
				"amount": "90"
			}
		},
		{
			"address": "comdex198r47wg90fdc3l4qg506l3vdlqe75ul3n3c0cm",
			"reward": {
				"denom": "ucmdx",
				"amount": "94"
			}
		},
		{
			"address": "comdex1982f5umpj48enr5rgf6tzc92u7m829eeqtp2yk",
			"reward": {
				"denom": "ucmdx",
				"amount": "680"
			}
		},
		{
			"address": "comdex1982us7eyttvkselj7hrnfqgguu7la92z9e225d",
			"reward": {
				"denom": "ucmdx",
				"amount": "2822"
			}
		},
		{
			"address": "comdex1980hcusmw4w8j0x09uuct0j5tff8k27s3nx0sz",
			"reward": {
				"denom": "ucmdx",
				"amount": "16"
			}
		},
		{
			"address": "comdex198nwjee6042qusptzg8xmg2s6z7etxsheldrjy",
			"reward": {
				"denom": "ucmdx",
				"amount": "32199"
			}
		},
		{
			"address": "comdex1984xgt0e3q4txtjmkdexqq38yls4r2kqfvj8nk",
			"reward": {
				"denom": "ucmdx",
				"amount": "130112"
			}
		},
		{
			"address": "comdex1984xcee3ejd65w9zmjr2stawy7a2nqzqwv768s",
			"reward": {
				"denom": "ucmdx",
				"amount": "12986"
			}
		},
		{
			"address": "comdex198ckuxk6epfktgujqf22qtunppmwh9f75z5em4",
			"reward": {
				"denom": "ucmdx",
				"amount": "4238"
			}
		},
		{
			"address": "comdex1986lw40c5c2x4957frmw6rt00jxv0mwau8dw26",
			"reward": {
				"denom": "ucmdx",
				"amount": "181"
			}
		},
		{
			"address": "comdex198m2cfnlm3dgul9yjzfkvjn00nuvkt77m5c2pn",
			"reward": {
				"denom": "ucmdx",
				"amount": "2596"
			}
		},
		{
			"address": "comdex198ladsstypgajnjnkgf0g8w89wdmjvcq8qj802",
			"reward": {
				"denom": "ucmdx",
				"amount": "166"
			}
		},
		{
			"address": "comdex19gqrpgympu9875gn3j72yr6v2amas0tjuxln5x",
			"reward": {
				"denom": "ucmdx",
				"amount": "129004"
			}
		},
		{
			"address": "comdex19gqgjg3qwvwa5l5czxfkd5dadxejrdc7zycp7l",
			"reward": {
				"denom": "ucmdx",
				"amount": "28806"
			}
		},
		{
			"address": "comdex19gpqukys0n38zegmu4adwq8yacw9kktx6ckeva",
			"reward": {
				"denom": "ucmdx",
				"amount": "526"
			}
		},
		{
			"address": "comdex19gphkpkf29fq6v55rxp7a3wmq97wyhqka7srxr",
			"reward": {
				"denom": "ucmdx",
				"amount": "1383"
			}
		},
		{
			"address": "comdex19gy80qz4nzuu2np2vfyjhu0wszxg6addm633jf",
			"reward": {
				"denom": "ucmdx",
				"amount": "7911"
			}
		},
		{
			"address": "comdex19g9pnfpa4jeusz087lt52g8vewcphhkchnffa6",
			"reward": {
				"denom": "ucmdx",
				"amount": "14327"
			}
		},
		{
			"address": "comdex19ggwu0vsfw0kcgx9tkxnev09ctv59auc9wgyv9",
			"reward": {
				"denom": "ucmdx",
				"amount": "16237"
			}
		},
		{
			"address": "comdex19gfwy3h3ywz27t668qptqu2rp050d3c9tq9fny",
			"reward": {
				"denom": "ucmdx",
				"amount": "199"
			}
		},
		{
			"address": "comdex19gw0twaee29007kxdu2zfzvn4ut4cvwty9uenv",
			"reward": {
				"denom": "ucmdx",
				"amount": "1190"
			}
		},
		{
			"address": "comdex19gsmm6fsmc2xhk2xzs5yp7am9a4f8u5a9fknda",
			"reward": {
				"denom": "ucmdx",
				"amount": "1885"
			}
		},
		{
			"address": "comdex19g3yrtt2jz52e0h9u6cxrg7j22yh5uff9f6jcm",
			"reward": {
				"denom": "ucmdx",
				"amount": "1933"
			}
		},
		{
			"address": "comdex19g48nzz2xtwn2p8cskarqsypvx5jw8ncvnhujv",
			"reward": {
				"denom": "ucmdx",
				"amount": "1475"
			}
		},
		{
			"address": "comdex19g4tfg8axrt47e9uhwd0fqje9j7fakfh9ft2ks",
			"reward": {
				"denom": "ucmdx",
				"amount": "1736"
			}
		},
		{
			"address": "comdex19gk75r4pd0e6hzz60nqgg3nps5xh980ntpgkj4",
			"reward": {
				"denom": "ucmdx",
				"amount": "92915"
			}
		},
		{
			"address": "comdex19ghln2cm5jsne4mcc88lxhukdpqf9y66szjnd0",
			"reward": {
				"denom": "ucmdx",
				"amount": "1572"
			}
		},
		{
			"address": "comdex19gcza5dmm9dsrm8qquts5vhfhm3y2c7y5qzsta",
			"reward": {
				"denom": "ucmdx",
				"amount": "708"
			}
		},
		{
			"address": "comdex19fq8wms5gve6599xm2wh573gspkccyq80xtepn",
			"reward": {
				"denom": "ucmdx",
				"amount": "1430"
			}
		},
		{
			"address": "comdex19frvek60qkhax7uekw5wm5na66l03jay67mfx9",
			"reward": {
				"denom": "ucmdx",
				"amount": "13208"
			}
		},
		{
			"address": "comdex19fgljxshl04guenjt9a58sfxauc7p34txf2ftk",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex19f2edvj89efvw58f43vypzxkd25flcd39qs6jq",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex19fs5dy8yhgg5dvzc4863vq6llszegc5apl0t55",
			"reward": {
				"denom": "ucmdx",
				"amount": "28736"
			}
		},
		{
			"address": "comdex19f3ydhvj9e5d72m808fmyz3jxyxf5z502a2qrl",
			"reward": {
				"denom": "ucmdx",
				"amount": "35"
			}
		},
		{
			"address": "comdex19fnw5kv6x36t938expps89agd4dcq4wynq9x85",
			"reward": {
				"denom": "ucmdx",
				"amount": "90"
			}
		},
		{
			"address": "comdex19f4kqfnpldpcx04v6j286wz8en7yfpzsnznmjh",
			"reward": {
				"denom": "ucmdx",
				"amount": "5942"
			}
		},
		{
			"address": "comdex19fcgsm32kyljpc3dlxw7czyk7xkq26h5fu7j7u",
			"reward": {
				"denom": "ucmdx",
				"amount": "2965"
			}
		},
		{
			"address": "comdex19fck52a2zlvcr2cy80x0npxm7jqg6megyjpx7w",
			"reward": {
				"denom": "ucmdx",
				"amount": "365531"
			}
		},
		{
			"address": "comdex19fu95j620avmtz0wxh3jkcr2tmznduqwlfcp6d",
			"reward": {
				"denom": "ucmdx",
				"amount": "527127"
			}
		},
		{
			"address": "comdex19farfc38qn8tj3jd8nc6nsgm7de3ssuzu8xhy0",
			"reward": {
				"denom": "ucmdx",
				"amount": "687"
			}
		},
		{
			"address": "comdex19fad64jvq420yeu27zmmxergakalrek48qfrcp",
			"reward": {
				"denom": "ucmdx",
				"amount": "14563"
			}
		},
		{
			"address": "comdex19f7f7ylewuf6wzenqn5p4l4nsrrwcta6j7zgvh",
			"reward": {
				"denom": "ucmdx",
				"amount": "147"
			}
		},
		{
			"address": "comdex19f7nq7cs5lxgartse9u3jutaz98w20lc747q23",
			"reward": {
				"denom": "ucmdx",
				"amount": "2912"
			}
		},
		{
			"address": "comdex192yswn23xs4m342ytyp58lmr6mnyjmpvyhxvrf",
			"reward": {
				"denom": "ucmdx",
				"amount": "1592"
			}
		},
		{
			"address": "comdex192yc3zp7gl987py3cvlsvh28x5dc5pmhyqzuvf",
			"reward": {
				"denom": "ucmdx",
				"amount": "5715"
			}
		},
		{
			"address": "comdex1928edtzxcy9qk6hg7w4d9deqx8ghgcz4pgrfda",
			"reward": {
				"denom": "ucmdx",
				"amount": "1552"
			}
		},
		{
			"address": "comdex192guner3ahpkfegtptmw2dfxscdqvg3f9wxkh4",
			"reward": {
				"denom": "ucmdx",
				"amount": "123004"
			}
		},
		{
			"address": "comdex1922hzxp2xk348mr9jmqdl5efmkh495qp85mjz7",
			"reward": {
				"denom": "ucmdx",
				"amount": "10259"
			}
		},
		{
			"address": "comdex192vk7dyp2cac79txtuupe3xvnnxkgxasvmmfuu",
			"reward": {
				"denom": "ucmdx",
				"amount": "1792"
			}
		},
		{
			"address": "comdex192wpc837s94vvg5uyagqsyc29yf0r6w3tsrunr",
			"reward": {
				"denom": "ucmdx",
				"amount": "14065"
			}
		},
		{
			"address": "comdex192wxe6u96s98m649hajmtfesganwnx3gzjg4vk",
			"reward": {
				"denom": "ucmdx",
				"amount": "181"
			}
		},
		{
			"address": "comdex192s2y8jzjmwcpsmrvqmthtxvl6ycn8yf23arln",
			"reward": {
				"denom": "ucmdx",
				"amount": "14305"
			}
		},
		{
			"address": "comdex192slevrmncrchrngq90cfuag58lsgmdkxd7eta",
			"reward": {
				"denom": "ucmdx",
				"amount": "5951"
			}
		},
		{
			"address": "comdex1923utg7e27z6cztq7ky89hx5n5k0awkhqr36mq",
			"reward": {
				"denom": "ucmdx",
				"amount": "266"
			}
		},
		{
			"address": "comdex192kfv9yqk8eucfvret74lc93cxxmxkz3ku8c8x",
			"reward": {
				"denom": "ucmdx",
				"amount": "1349"
			}
		},
		{
			"address": "comdex192uvwncyl6gyhyxww3uaqq3jtcsvv2hamgyulu",
			"reward": {
				"denom": "ucmdx",
				"amount": "166"
			}
		},
		{
			"address": "comdex192uuw9a8vgfytqm0gughqklxmqkj68fvsf027h",
			"reward": {
				"denom": "ucmdx",
				"amount": "194"
			}
		},
		{
			"address": "comdex1927khwp24rudv4yyux9dhp05ajuvqhkd2xexps",
			"reward": {
				"denom": "ucmdx",
				"amount": "3981"
			}
		},
		{
			"address": "comdex19tq4cl0euaa92s6hfu2fl3n9y88khenta3wxjr",
			"reward": {
				"denom": "ucmdx",
				"amount": "178"
			}
		},
		{
			"address": "comdex19tx2rvtmd904p8txh27cah0j3duy2s0gx939tt",
			"reward": {
				"denom": "ucmdx",
				"amount": "8742"
			}
		},
		{
			"address": "comdex19tg3mpsydqrqm8s9h43vqfh2wct0a4pu9slj68",
			"reward": {
				"denom": "ucmdx",
				"amount": "143"
			}
		},
		{
			"address": "comdex19t2j534y998y0863mwyczzxz8nf6h0ag4tu3yq",
			"reward": {
				"denom": "ucmdx",
				"amount": "614"
			}
		},
		{
			"address": "comdex19tvd38rhdzhnuxv34a5jgayp7uej5e6j9rs82d",
			"reward": {
				"denom": "ucmdx",
				"amount": "2838"
			}
		},
		{
			"address": "comdex19t0djs44lqj0uh3exye5lv307h20h4cl7hfd6g",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex19ts5hyt3whayule98t6lrvpfxp0jg4hu89xkqc",
			"reward": {
				"denom": "ucmdx",
				"amount": "2238"
			}
		},
		{
			"address": "comdex19t5tqtykmd8xfmuc2h0nqy27arxeyuap70j6rp",
			"reward": {
				"denom": "ucmdx",
				"amount": "16"
			}
		},
		{
			"address": "comdex19tkec6ptf90qu69cnsau03c6yc55rpwvkrt3uu",
			"reward": {
				"denom": "ucmdx",
				"amount": "66262"
			}
		},
		{
			"address": "comdex19tcvllfk8ps49ph7zv647ewlgg4dwvc2kagvgt",
			"reward": {
				"denom": "ucmdx",
				"amount": "2380"
			}
		},
		{
			"address": "comdex19te6prqalzfldkvj8p49wgrldex9rpmxp229td",
			"reward": {
				"denom": "ucmdx",
				"amount": "486"
			}
		},
		{
			"address": "comdex19t6kcrckngw3mz7xg6d0vr7tqt84wxv9atl9hj",
			"reward": {
				"denom": "ucmdx",
				"amount": "5861"
			}
		},
		{
			"address": "comdex19tmhae3hy2gy6vfsq8k4ur90lmdgajtt0swha2",
			"reward": {
				"denom": "ucmdx",
				"amount": "17718"
			}
		},
		{
			"address": "comdex19tu6j665jkglj678qje7py5qzwxemztjtneugf",
			"reward": {
				"denom": "ucmdx",
				"amount": "6089"
			}
		},
		{
			"address": "comdex19t7k5wxtc9pvmxku34ltj030czt906e5v9dpew",
			"reward": {
				"denom": "ucmdx",
				"amount": "6796"
			}
		},
		{
			"address": "comdex19vqdwvedl9gkvcte0khm96ur7gyzx75pkd7r8p",
			"reward": {
				"denom": "ucmdx",
				"amount": "1471"
			}
		},
		{
			"address": "comdex19vqjj6a4wcsnd82v50y9payde726f5kexcrf8s",
			"reward": {
				"denom": "ucmdx",
				"amount": "97658"
			}
		},
		{
			"address": "comdex19vz6m9cxysjlftmg0y9xglntxf232e74szdy6c",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex19vrzd4y6dex46rzy0t5kavcjv3v622pxvu2hun",
			"reward": {
				"denom": "ucmdx",
				"amount": "13359"
			}
		},
		{
			"address": "comdex19vr047wr08gt8kme76r8m02nkf9hdq57jnax9u",
			"reward": {
				"denom": "ucmdx",
				"amount": "976"
			}
		},
		{
			"address": "comdex19v9pkens9uaplrdh6kl7rp9j482uhv3r0pfa0u",
			"reward": {
				"denom": "ucmdx",
				"amount": "4284"
			}
		},
		{
			"address": "comdex19v2fjt0sjfmqjeraa7j2deyjwkpexjs2lzepd8",
			"reward": {
				"denom": "ucmdx",
				"amount": "24874"
			}
		},
		{
			"address": "comdex19vdqe9elnvgkwx3q3vm76g9hv8mu4xujagyuvt",
			"reward": {
				"denom": "ucmdx",
				"amount": "1877"
			}
		},
		{
			"address": "comdex19v3gfr7y9dvxr2z65wah3fwucp4utf4tux5v77",
			"reward": {
				"denom": "ucmdx",
				"amount": "1375"
			}
		},
		{
			"address": "comdex19v4gyxhf3z9q97wmpq7ztwzqvpdzxavqqasuaf",
			"reward": {
				"denom": "ucmdx",
				"amount": "31827"
			}
		},
		{
			"address": "comdex19vm8jct0ql0ux86cvzpqje6tp6wwkm4eep064s",
			"reward": {
				"denom": "ucmdx",
				"amount": "2057"
			}
		},
		{
			"address": "comdex19vm0uc3mdaehlmlfmvp9ray9wt2tr3htwzeesf",
			"reward": {
				"denom": "ucmdx",
				"amount": "173"
			}
		},
		{
			"address": "comdex19v7p0233gn7n3qmar6rd2afdq053pphg0twj2h",
			"reward": {
				"denom": "ucmdx",
				"amount": "145140"
			}
		},
		{
			"address": "comdex19dq0ra270296mde89qx8jzemahtt2da87mrdmp",
			"reward": {
				"denom": "ucmdx",
				"amount": "185"
			}
		},
		{
			"address": "comdex19dpm2mseuymddclzshjvjx3a63mrhx9puy0xyw",
			"reward": {
				"denom": "ucmdx",
				"amount": "86"
			}
		},
		{
			"address": "comdex19druaq2ff3fjrvujnaxwf63krsrwhghz50tj5f",
			"reward": {
				"denom": "ucmdx",
				"amount": "26152"
			}
		},
		{
			"address": "comdex19d97guw9swjwruwwg30xgfzu0e0l3kgm9xtn2q",
			"reward": {
				"denom": "ucmdx",
				"amount": "2043"
			}
		},
		{
			"address": "comdex19dx6xr5pqg83f45d43xfsuczvnlk37w8acrfnh",
			"reward": {
				"denom": "ucmdx",
				"amount": "502"
			}
		},
		{
			"address": "comdex19d88a65deer4jt7etqgd2rsnv6l75ydh8ha8yf",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex19dde6w7akqt3uwluh2ypvyezztd9w2pulfz2h6",
			"reward": {
				"denom": "ucmdx",
				"amount": "35986"
			}
		},
		{
			"address": "comdex19d03tk62npt7vdemr0aszt7t53k6vvu8d2zmct",
			"reward": {
				"denom": "ucmdx",
				"amount": "10136"
			}
		},
		{
			"address": "comdex19djm47wsk5h08t7je3mnrxhu6j9749zdf2d38l",
			"reward": {
				"denom": "ucmdx",
				"amount": "2647"
			}
		},
		{
			"address": "comdex19dn5aalascydkpca725msfgkjmxh8hpj3t2pnf",
			"reward": {
				"denom": "ucmdx",
				"amount": "74022"
			}
		},
		{
			"address": "comdex19d5uz5tyssz77e5m2vqenfxs2s85975a2lf3kg",
			"reward": {
				"denom": "ucmdx",
				"amount": "1096"
			}
		},
		{
			"address": "comdex19dhrq6fuhfmc5rny83y4ldsj3ux2qm8ms5wwng",
			"reward": {
				"denom": "ucmdx",
				"amount": "28"
			}
		},
		{
			"address": "comdex19dh608wxj82zqljzxl65zcd69a48kygp4hk5zn",
			"reward": {
				"denom": "ucmdx",
				"amount": "698"
			}
		},
		{
			"address": "comdex19dupdkqnlacthy8mezxvzsy263247rdchd5zfg",
			"reward": {
				"denom": "ucmdx",
				"amount": "176"
			}
		},
		{
			"address": "comdex19dlpu8wk43eazrjpyk9068uey4s2vrmtwvlw7t",
			"reward": {
				"denom": "ucmdx",
				"amount": "34"
			}
		},
		{
			"address": "comdex19wzrpwe9e6pksmptxx8ldn3qtf2vhf5gz7swzg",
			"reward": {
				"denom": "ucmdx",
				"amount": "1362"
			}
		},
		{
			"address": "comdex19wz7f9rh356pfk5ffe3ytvx54jdgkra8h883ur",
			"reward": {
				"denom": "ucmdx",
				"amount": "12905"
			}
		},
		{
			"address": "comdex19w9j3dkl2fa7v94tlf9t7vud9zf44mx42wc0qf",
			"reward": {
				"denom": "ucmdx",
				"amount": "18966"
			}
		},
		{
			"address": "comdex19w8c98qalk0hrf229zaetjew3h5xak4mvtyxev",
			"reward": {
				"denom": "ucmdx",
				"amount": "10734"
			}
		},
		{
			"address": "comdex19wv7srdrfn9qtt69jf3ep8ks280qgvk0x3pfu0",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex19wsl7r5a3xt22f9kj6x506s28gmy3lc525kzgs",
			"reward": {
				"denom": "ucmdx",
				"amount": "2247"
			}
		},
		{
			"address": "comdex19wn6kfldgl4k7ul6gxthnrgng6jq28d4u922g0",
			"reward": {
				"denom": "ucmdx",
				"amount": "16672"
			}
		},
		{
			"address": "comdex19whhfes94tn8dl84f8jy8c5a0fyfkx7jc5xpht",
			"reward": {
				"denom": "ucmdx",
				"amount": "697"
			}
		},
		{
			"address": "comdex19wheec06zqlpk87vtxh3dd9c72v7uf9sgv38q9",
			"reward": {
				"denom": "ucmdx",
				"amount": "70305"
			}
		},
		{
			"address": "comdex19w6d9yrlly2tw0sf0f3lz7gay0m0jlfhl78atx",
			"reward": {
				"denom": "ucmdx",
				"amount": "25"
			}
		},
		{
			"address": "comdex19w7yqx73qt62j4hdg9k0emxsluv05h3jva555k",
			"reward": {
				"denom": "ucmdx",
				"amount": "5704"
			}
		},
		{
			"address": "comdex19w7e74fj4hl3eppl46y53y48llf7jd2m2c5qjx",
			"reward": {
				"denom": "ucmdx",
				"amount": "10359"
			}
		},
		{
			"address": "comdex19wlqhpvz4pqpck94zpmmhlhsqgmufx82cs58x7",
			"reward": {
				"denom": "ucmdx",
				"amount": "6660"
			}
		},
		{
			"address": "comdex190yklca66sqhfljdpc4mmqkmjmg95y3agkcntr",
			"reward": {
				"denom": "ucmdx",
				"amount": "8869"
			}
		},
		{
			"address": "comdex1909zsug47a9j5pwnkwk363gcl0uh6clyc60eg5",
			"reward": {
				"denom": "ucmdx",
				"amount": "1276"
			}
		},
		{
			"address": "comdex190f36xptgrpdzdezf285c306jnvgtzvrml5c6y",
			"reward": {
				"denom": "ucmdx",
				"amount": "16866"
			}
		},
		{
			"address": "comdex1900zeh5v405w7pu4mvw7tzcjkqqzhvqrhstnzu",
			"reward": {
				"denom": "ucmdx",
				"amount": "1778"
			}
		},
		{
			"address": "comdex19002ghhvt5ywlmrd5gsjm7zdtedjvsesxjyvej",
			"reward": {
				"denom": "ucmdx",
				"amount": "4497"
			}
		},
		{
			"address": "comdex190sj2y34wu46927z4lrs3pkey8j8k68gml542r",
			"reward": {
				"denom": "ucmdx",
				"amount": "45976"
			}
		},
		{
			"address": "comdex1906zvfm4hcd9kx4cqsgc2dp2dpwscr0p06rg7r",
			"reward": {
				"denom": "ucmdx",
				"amount": "70079"
			}
		},
		{
			"address": "comdex190a4q8x8tk5uhp3cnqsl2uwuhcjra8s49ra6wa",
			"reward": {
				"denom": "ucmdx",
				"amount": "1452"
			}
		},
		{
			"address": "comdex19073hgge95fw5hvh4ejg97stvy3hhdp4xfsn6z",
			"reward": {
				"denom": "ucmdx",
				"amount": "7410"
			}
		},
		{
			"address": "comdex19sq6rqusnyeyg873w534m9sdtmyg8u22wr7cj5",
			"reward": {
				"denom": "ucmdx",
				"amount": "171"
			}
		},
		{
			"address": "comdex19sph5wrfs2kh7xgh27fksadv93dctjpzepnn5t",
			"reward": {
				"denom": "ucmdx",
				"amount": "26130"
			}
		},
		{
			"address": "comdex19s9xt8vm9js0jk7r3pe5kc4q7sxz9wp57jzq85",
			"reward": {
				"denom": "ucmdx",
				"amount": "20992"
			}
		},
		{
			"address": "comdex19s9870729lngap2kzls3dcyeu7a4m2tchsv74j",
			"reward": {
				"denom": "ucmdx",
				"amount": "319"
			}
		},
		{
			"address": "comdex19sxwwn09hjweu4fmslldm85txwcs8esenrhm2d",
			"reward": {
				"denom": "ucmdx",
				"amount": "2888"
			}
		},
		{
			"address": "comdex19s839fk4cyk784nsrduuzchywcz9t9kgsufzvy",
			"reward": {
				"denom": "ucmdx",
				"amount": "980"
			}
		},
		{
			"address": "comdex19s06gae5mjepzyg2zpr3s4en6qxy8xdqa6kv9m",
			"reward": {
				"denom": "ucmdx",
				"amount": "3885"
			}
		},
		{
			"address": "comdex19sjw0mrctwkhjga2956kar6mhrq7trv6kjsfsk",
			"reward": {
				"denom": "ucmdx",
				"amount": "15066"
			}
		},
		{
			"address": "comdex19sns754s5kxwgtfvcldzq3wu77xn2m2pm0vn5h",
			"reward": {
				"denom": "ucmdx",
				"amount": "1398"
			}
		},
		{
			"address": "comdex19s5wmpavsqc0vsfahj8lnll05c5jtnelkjg2te",
			"reward": {
				"denom": "ucmdx",
				"amount": "1511"
			}
		},
		{
			"address": "comdex19shw0nvlhwzhu4jzen273pfudaujkus0hf00qd",
			"reward": {
				"denom": "ucmdx",
				"amount": "177490"
			}
		},
		{
			"address": "comdex19shwmrhqdznlar4y7mgth2ugaysqtpyf3zy227",
			"reward": {
				"denom": "ucmdx",
				"amount": "131"
			}
		},
		{
			"address": "comdex19s6zq4xwwhevglch5kl7zundvhj2lw4g4c6sr8",
			"reward": {
				"denom": "ucmdx",
				"amount": "14824"
			}
		},
		{
			"address": "comdex19sackq272vs2kh02f78xxnkq79sj77a7u9l70m",
			"reward": {
				"denom": "ucmdx",
				"amount": "6057"
			}
		},
		{
			"address": "comdex19s7nrzxwlwwt4pdk4dnk9wltnr67w6sd7plxjd",
			"reward": {
				"denom": "ucmdx",
				"amount": "37498"
			}
		},
		{
			"address": "comdex193rzkz9kkq54n57jfuj26vpgcv4tskq6yq4tt2",
			"reward": {
				"denom": "ucmdx",
				"amount": "179"
			}
		},
		{
			"address": "comdex193r2hq2u2usqrf9cjeg6n97kx4m76cxnmyfdwf",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex193y2auy23epr53qt9xj80zaujvk5wk7lflen4x",
			"reward": {
				"denom": "ucmdx",
				"amount": "150"
			}
		},
		{
			"address": "comdex193ywksg84mxxsnr2p93546z7s03n084w4xm3gz",
			"reward": {
				"denom": "ucmdx",
				"amount": "169"
			}
		},
		{
			"address": "comdex1938qljjfgdzm2pwef4v9894k6vg6ymege9cz38",
			"reward": {
				"denom": "ucmdx",
				"amount": "3564"
			}
		},
		{
			"address": "comdex193fzh5pz8w49kq2vztmtv4uwamud79lrlkemrp",
			"reward": {
				"denom": "ucmdx",
				"amount": "2510"
			}
		},
		{
			"address": "comdex193tu8jn0ruzvguchuj0sfxxce5rq6g8g0nvzm9",
			"reward": {
				"denom": "ucmdx",
				"amount": "27176"
			}
		},
		{
			"address": "comdex193vm38q84uz2475r6g04mhl085wcm64rej50yr",
			"reward": {
				"denom": "ucmdx",
				"amount": "2113367"
			}
		},
		{
			"address": "comdex193d2qsuk25e9mx0quan2msl2y9fqtt0g84jrk2",
			"reward": {
				"denom": "ucmdx",
				"amount": "354"
			}
		},
		{
			"address": "comdex193d6vvajjxzjfvrlqd6dq8nq2hgxt09zagzf4v",
			"reward": {
				"denom": "ucmdx",
				"amount": "23942"
			}
		},
		{
			"address": "comdex1933nc6q9gxyxyrp06ctardejecpapgnugf99tn",
			"reward": {
				"denom": "ucmdx",
				"amount": "165"
			}
		},
		{
			"address": "comdex193eduvlfn7g3rhjpk0xl7zunh4ed2f6ec8gzvk",
			"reward": {
				"denom": "ucmdx",
				"amount": "176"
			}
		},
		{
			"address": "comdex193mwhjpa62zymyjnkzdne50ru85c2nmy0pquyj",
			"reward": {
				"denom": "ucmdx",
				"amount": "151"
			}
		},
		{
			"address": "comdex19jqqqnsk8uush6tnjjyercqatyj2f7kglvn6g2",
			"reward": {
				"denom": "ucmdx",
				"amount": "9569"
			}
		},
		{
			"address": "comdex19jrw4dceet0q8ls6pc9krkr0p0lgkm3w3el2lz",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex19jyn93h8f2pafve0wxtpyktpdpjgrxvknw5n4n",
			"reward": {
				"denom": "ucmdx",
				"amount": "35266"
			}
		},
		{
			"address": "comdex19j97lsr8d6arcjhmf993nsdc8xwzt0tmq607pc",
			"reward": {
				"denom": "ucmdx",
				"amount": "3100"
			}
		},
		{
			"address": "comdex19jxaymxxdkr2uscqcvk23zf50e8y0tyydgd99d",
			"reward": {
				"denom": "ucmdx",
				"amount": "395"
			}
		},
		{
			"address": "comdex19j8nt48al76j92np9p84awnkdwukg95gspknak",
			"reward": {
				"denom": "ucmdx",
				"amount": "940"
			}
		},
		{
			"address": "comdex19jgzvmgjas6pzgq2gk0hug3w83qnrnvglsamyq",
			"reward": {
				"denom": "ucmdx",
				"amount": "5373"
			}
		},
		{
			"address": "comdex19jge2paz4czad6wkaqwewkydnjc83j4p02lcz2",
			"reward": {
				"denom": "ucmdx",
				"amount": "5657"
			}
		},
		{
			"address": "comdex19jflm0as89elx7sfhx9mtxm7vt3qlyuftylhel",
			"reward": {
				"denom": "ucmdx",
				"amount": "130"
			}
		},
		{
			"address": "comdex19jtn9xx65x5q5z5d5k2egkk2h25a9pqffljca6",
			"reward": {
				"denom": "ucmdx",
				"amount": "43583"
			}
		},
		{
			"address": "comdex19j0kx8g09nah8dkrw9f9svt4f0ntgyz94yssj0",
			"reward": {
				"denom": "ucmdx",
				"amount": "67904"
			}
		},
		{
			"address": "comdex19jsvr5luakkc57klq8nc8gx6fetsn2m83nxgjy",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex19jselk552hjyyhc9r3cmfvgfe7nfewm7ksq42j",
			"reward": {
				"denom": "ucmdx",
				"amount": "10143"
			}
		},
		{
			"address": "comdex19j48p95kpgme9qvw8p0eucp95z2557a3hg857u",
			"reward": {
				"denom": "ucmdx",
				"amount": "5372"
			}
		},
		{
			"address": "comdex19j44kytuewzcdc5euzanxj3xx7jfcafum77ev9",
			"reward": {
				"denom": "ucmdx",
				"amount": "23993"
			}
		},
		{
			"address": "comdex19jkt4vgzvmm9p0uhmhrutjg0mhac62fhwg932t",
			"reward": {
				"denom": "ucmdx",
				"amount": "28"
			}
		},
		{
			"address": "comdex19jk5xruuhkemyguj33gytcxjeuc49mlwehx972",
			"reward": {
				"denom": "ucmdx",
				"amount": "400"
			}
		},
		{
			"address": "comdex19jk5nyef4fdpc0v35q5jwpeje68ywdlvcw8rgp",
			"reward": {
				"denom": "ucmdx",
				"amount": "2365"
			}
		},
		{
			"address": "comdex19jcvkwl3uywz4mzxngf8488uy6n9psvx8jfksq",
			"reward": {
				"denom": "ucmdx",
				"amount": "5"
			}
		},
		{
			"address": "comdex19jmk5t4pu22zjzx6vc8upsdm5cugq530xaccrw",
			"reward": {
				"denom": "ucmdx",
				"amount": "173"
			}
		},
		{
			"address": "comdex19j7tfzwp7xdcw6j6ds7v2raldkc87jkahj20t6",
			"reward": {
				"denom": "ucmdx",
				"amount": "14586"
			}
		},
		{
			"address": "comdex19ny3l9lm2m4me7qhgvjrs9ya4fwlae44nc7vxs",
			"reward": {
				"denom": "ucmdx",
				"amount": "3604"
			}
		},
		{
			"address": "comdex19nfp657kqq8nqy9vu2mwewxjrahrc5qprld4je",
			"reward": {
				"denom": "ucmdx",
				"amount": "725"
			}
		},
		{
			"address": "comdex19nfe2l95pnmre8snp3eddvzun29s5rtjkj6az0",
			"reward": {
				"denom": "ucmdx",
				"amount": "1787"
			}
		},
		{
			"address": "comdex19nvkp97wnmtphc9jcfdqw8plfdh4ktxrgz7v0n",
			"reward": {
				"denom": "ucmdx",
				"amount": "1450"
			}
		},
		{
			"address": "comdex19nsyf4u803h9h50zclcph3ywy2tdnhpw6vvqjy",
			"reward": {
				"denom": "ucmdx",
				"amount": "151"
			}
		},
		{
			"address": "comdex19n5f66m9utezge4khzct2m5kvvd8wj4lweqfrg",
			"reward": {
				"denom": "ucmdx",
				"amount": "10061"
			}
		},
		{
			"address": "comdex19n5767xecsvffr46cx4jym67373gae83k0qv9m",
			"reward": {
				"denom": "ucmdx",
				"amount": "178"
			}
		},
		{
			"address": "comdex19n4rcrxwzxt9xfkdllc833pd7w0p39jjahsdh6",
			"reward": {
				"denom": "ucmdx",
				"amount": "11635"
			}
		},
		{
			"address": "comdex19nmt2j957apgy8r3nc5npjq2ex342r6urz6rwx",
			"reward": {
				"denom": "ucmdx",
				"amount": "4309"
			}
		},
		{
			"address": "comdex19nm3mhkxnaq2hceq9a37d9tmu6z0jz75g94zmd",
			"reward": {
				"denom": "ucmdx",
				"amount": "22484"
			}
		},
		{
			"address": "comdex19nuyr8mfkxkevcnr9gveznzvqpypgl5qj38kyz",
			"reward": {
				"denom": "ucmdx",
				"amount": "6293"
			}
		},
		{
			"address": "comdex19nuu5l7sgn582g27ec64aek4psau45fvmj22me",
			"reward": {
				"denom": "ucmdx",
				"amount": "742415"
			}
		},
		{
			"address": "comdex19n7n9vtnanucvx6akh0ly3f038nmet0ljqc3m7",
			"reward": {
				"denom": "ucmdx",
				"amount": "1021"
			}
		},
		{
			"address": "comdex195qkp24qhesuk6vnen7uxn6kvsqpulundgmzgz",
			"reward": {
				"denom": "ucmdx",
				"amount": "598"
			}
		},
		{
			"address": "comdex195pdvs7lnwlel73uq79l6nn0hkdhtc7gnvprm7",
			"reward": {
				"denom": "ucmdx",
				"amount": "895"
			}
		},
		{
			"address": "comdex1959zwl34lxxaadrkaeqgvrhfmljjt5y0r4x58u",
			"reward": {
				"denom": "ucmdx",
				"amount": "94179"
			}
		},
		{
			"address": "comdex195f0u84l3t2t3ds3crll66ldmj4pwm4kr42kul",
			"reward": {
				"denom": "ucmdx",
				"amount": "76022"
			}
		},
		{
			"address": "comdex195fhfz36cukv2zqdxuc2qrk5n7vxde82pxw3el",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1952kp8ndu63k4kd7qjaxzuyxe9yrcfmmquvqhp",
			"reward": {
				"denom": "ucmdx",
				"amount": "945"
			}
		},
		{
			"address": "comdex195t9yv9dvdtgwfh3rkmyuyx76gz9uv3l9n3a8c",
			"reward": {
				"denom": "ucmdx",
				"amount": "3354"
			}
		},
		{
			"address": "comdex195tavlpd8sy253z2mr748zv6c0qk3kl7ejrert",
			"reward": {
				"denom": "ucmdx",
				"amount": "394"
			}
		},
		{
			"address": "comdex195w6wdeyksrpkvw6nc7depcdp6gc6epsa0nks8",
			"reward": {
				"denom": "ucmdx",
				"amount": "74575"
			}
		},
		{
			"address": "comdex195j82u6zsl0gveg7aa9n6e00p3fu04fm7whjpq",
			"reward": {
				"denom": "ucmdx",
				"amount": "2264"
			}
		},
		{
			"address": "comdex195nf2aplw8emkag36sjgdpnfva9s9qkgcs5rzx",
			"reward": {
				"denom": "ucmdx",
				"amount": "166"
			}
		},
		{
			"address": "comdex195k393ez550cg7q2p68dee6xv8kqmsz4cl7mdj",
			"reward": {
				"denom": "ucmdx",
				"amount": "1483"
			}
		},
		{
			"address": "comdex195emnutzcz3vsapf0sepudkh4v82tn0xmwlfny",
			"reward": {
				"denom": "ucmdx",
				"amount": "530"
			}
		},
		{
			"address": "comdex195ea098ndztvdt286l02y0jcg6des7pysqzhmq",
			"reward": {
				"denom": "ucmdx",
				"amount": "52"
			}
		},
		{
			"address": "comdex1956zxfmd03w90vs7nnpwlw52t64vm3k6hzshaz",
			"reward": {
				"denom": "ucmdx",
				"amount": "6896"
			}
		},
		{
			"address": "comdex195u3ns7y4p88tm3namjeyzc07uqhkszxzuskuq",
			"reward": {
				"denom": "ucmdx",
				"amount": "8733"
			}
		},
		{
			"address": "comdex1957qaateklhv4sygqxawccrjedkv6l08jup4x2",
			"reward": {
				"denom": "ucmdx",
				"amount": "4279"
			}
		},
		{
			"address": "comdex195lfnvs7jeevj6prn5a9ghr563s6tmzl65fr8n",
			"reward": {
				"denom": "ucmdx",
				"amount": "179"
			}
		},
		{
			"address": "comdex194qp2uhhl9vqgvteeuewpatar94n0mdq3xfnes",
			"reward": {
				"denom": "ucmdx",
				"amount": "2855"
			}
		},
		{
			"address": "comdex194qu98jmfytgg6y5t59xhr4quuza6ll8y37ve9",
			"reward": {
				"denom": "ucmdx",
				"amount": "29549"
			}
		},
		{
			"address": "comdex194px257lyndy5sldjkqqr2evxljd9ymeqc0g0s",
			"reward": {
				"denom": "ucmdx",
				"amount": "1582"
			}
		},
		{
			"address": "comdex194xsjdzgn348d60lk5df44ykp6jk7wke4pg2y5",
			"reward": {
				"denom": "ucmdx",
				"amount": "16"
			}
		},
		{
			"address": "comdex19485fye7pm387x2fr2pgr49ruegquskc2lpw3q",
			"reward": {
				"denom": "ucmdx",
				"amount": "1213"
			}
		},
		{
			"address": "comdex194g63wd5y8ua8j7v475qca2w6s2wqq7paa6myv",
			"reward": {
				"denom": "ucmdx",
				"amount": "88"
			}
		},
		{
			"address": "comdex1942wfa3fwvela6gt80dzx4qykkl2ck63r0hz5r",
			"reward": {
				"denom": "ucmdx",
				"amount": "3555"
			}
		},
		{
			"address": "comdex194dw378xw2frmwkjct6kq5dqkldd3wwm8z8l6a",
			"reward": {
				"denom": "ucmdx",
				"amount": "40298"
			}
		},
		{
			"address": "comdex194dawsp29zcp4s9r6hdppdak0cy5kf3xvhfm9r",
			"reward": {
				"denom": "ucmdx",
				"amount": "8368"
			}
		},
		{
			"address": "comdex1940meq08et7gmxmljl6d3gxhr5hswd49l0p79n",
			"reward": {
				"denom": "ucmdx",
				"amount": "1117"
			}
		},
		{
			"address": "comdex1943pxeyv9982p2zekjq344zu9l9q86ew2tkvw2",
			"reward": {
				"denom": "ucmdx",
				"amount": "53194"
			}
		},
		{
			"address": "comdex194jxs9jvar80wtavkan5eppwve2jh4vmuvxgn0",
			"reward": {
				"denom": "ucmdx",
				"amount": "7176"
			}
		},
		{
			"address": "comdex194nez388h7mf9ncs4m8kaxyg23l5pgt95xfzlu",
			"reward": {
				"denom": "ucmdx",
				"amount": "298744"
			}
		},
		{
			"address": "comdex1945qxzyv5nswwfsa37h0v7spe9qa57ks20q3t8",
			"reward": {
				"denom": "ucmdx",
				"amount": "1098"
			}
		},
		{
			"address": "comdex1945wnglscxy8mzld3dt42u9j36xefl333qs9ym",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex194epjn65g69wwn6m6tff4p7zp5v0demq7pkwmq",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex194evedxwtff3lzz3p9vw6qwn6w2kd2p9apmuk5",
			"reward": {
				"denom": "ucmdx",
				"amount": "1958"
			}
		},
		{
			"address": "comdex1946l67mvw5xje4s8ud8vv6ky4n6lpmcslwhy35",
			"reward": {
				"denom": "ucmdx",
				"amount": "3541"
			}
		},
		{
			"address": "comdex194uclzlyl8et39ggg6gt2ndmaef5ktdadac3re",
			"reward": {
				"denom": "ucmdx",
				"amount": "5402"
			}
		},
		{
			"address": "comdex194agug687wrep7qfkl34vj5mgnhnv0ueyftjtk",
			"reward": {
				"denom": "ucmdx",
				"amount": "715"
			}
		},
		{
			"address": "comdex194lyh0fmcgdv5vafnqkxs0efjz7wg732ty4ss2",
			"reward": {
				"denom": "ucmdx",
				"amount": "1234"
			}
		},
		{
			"address": "comdex19kzhk63sq68ft2q2sqwxqfryw74fv0w2tsznrj",
			"reward": {
				"denom": "ucmdx",
				"amount": "84"
			}
		},
		{
			"address": "comdex19kreg0vke00n3qqa4s6ge427y8f7mwcp8exg7t",
			"reward": {
				"denom": "ucmdx",
				"amount": "285"
			}
		},
		{
			"address": "comdex19kyfagck4sh48fppchf2ufnhl8a8lrqy06ygdl",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex19ky0t3kmcky03nanj28v9d3urkx6qykwvn8fxv",
			"reward": {
				"denom": "ucmdx",
				"amount": "197204"
			}
		},
		{
			"address": "comdex19kfny9yae4wrgulaa6ru8fq9ssjcf9qs38z059",
			"reward": {
				"denom": "ucmdx",
				"amount": "171"
			}
		},
		{
			"address": "comdex19kva6tm9wpeeac8v3mx0kvdf45puejv4zvcttf",
			"reward": {
				"denom": "ucmdx",
				"amount": "728"
			}
		},
		{
			"address": "comdex19kd9jqhwa3ta9y8nhlyekt92rqwnpddln5zapu",
			"reward": {
				"denom": "ucmdx",
				"amount": "5616"
			}
		},
		{
			"address": "comdex19knqp0y5xvd4e9l6e3uf4v4zu2h9a3yggx6s0c",
			"reward": {
				"denom": "ucmdx",
				"amount": "14110"
			}
		},
		{
			"address": "comdex19k4qz76wljuxeh2t0hlxpkr5ll4malde8p64pc",
			"reward": {
				"denom": "ucmdx",
				"amount": "6124"
			}
		},
		{
			"address": "comdex19k4q96kd8xwkd3ykg7y772namlrswhgc7p4vw5",
			"reward": {
				"denom": "ucmdx",
				"amount": "113048"
			}
		},
		{
			"address": "comdex19kk2gcxyhp9qe8nrldctdjqygevqa3hn9nf68h",
			"reward": {
				"denom": "ucmdx",
				"amount": "18733"
			}
		},
		{
			"address": "comdex19kh5z7ujuuf4vlrphecw50vgywemdd7x2yz2eg",
			"reward": {
				"denom": "ucmdx",
				"amount": "22985"
			}
		},
		{
			"address": "comdex19kcazvgqzt6u9qh9fq7j6u2nym8ymcjwwat8fe",
			"reward": {
				"denom": "ucmdx",
				"amount": "8655"
			}
		},
		{
			"address": "comdex19k6uu6qh8jq8tp8ele030zd3n890kjmvf37f4t",
			"reward": {
				"denom": "ucmdx",
				"amount": "1470"
			}
		},
		{
			"address": "comdex19klagsh8f6kkje6s4u0kmex3dl9yg2vkw2vytg",
			"reward": {
				"denom": "ucmdx",
				"amount": "1757"
			}
		},
		{
			"address": "comdex19hq2m2rj04wdll8ngmwkx2lleaxmayty8n42sm",
			"reward": {
				"denom": "ucmdx",
				"amount": "2714"
			}
		},
		{
			"address": "comdex19hy5d3ha8v8krsgu34angnephkfzvtg9efl42w",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex19hyeze8tnvwx2tcd49g340zz8y0453fncsuln7",
			"reward": {
				"denom": "ucmdx",
				"amount": "179"
			}
		},
		{
			"address": "comdex19hxgdrufmct5v95meqfx8tjepryhf6p6v4px99",
			"reward": {
				"denom": "ucmdx",
				"amount": "1252"
			}
		},
		{
			"address": "comdex19hsyt5qrpyyjka9deny0hqyggwcwqte69l0u63",
			"reward": {
				"denom": "ucmdx",
				"amount": "4315"
			}
		},
		{
			"address": "comdex19h327payr2jr4qc2jvr55jfm9p76a07um9dqjj",
			"reward": {
				"denom": "ucmdx",
				"amount": "879"
			}
		},
		{
			"address": "comdex19hjq22c3u7a7d89yz7g8msmndxskzgfnyyjv8v",
			"reward": {
				"denom": "ucmdx",
				"amount": "18019"
			}
		},
		{
			"address": "comdex19hn7myfdsw7lff0ruqkfhaa35lr4h3re9xtddz",
			"reward": {
				"denom": "ucmdx",
				"amount": "289"
			}
		},
		{
			"address": "comdex19h5tvn8k07l70vdt3mp7tr505ksnmtm54rpx3j",
			"reward": {
				"denom": "ucmdx",
				"amount": "433"
			}
		},
		{
			"address": "comdex19h5u2ut2dnvzn866zr5288y8jupdfyfdn4k8q7",
			"reward": {
				"denom": "ucmdx",
				"amount": "271"
			}
		},
		{
			"address": "comdex19hkywmdum24d9sev8nlxhyg2j9mpdzfhfmmm6e",
			"reward": {
				"denom": "ucmdx",
				"amount": "6834"
			}
		},
		{
			"address": "comdex19hkmrn8sjwvplgcz5eefrc5mva5wgaxp58wxy9",
			"reward": {
				"denom": "ucmdx",
				"amount": "176"
			}
		},
		{
			"address": "comdex19hcvcq0xzhh2pgutp6574z3pvd4t5crdudmfy8",
			"reward": {
				"denom": "ucmdx",
				"amount": "7403"
			}
		},
		{
			"address": "comdex19hcegnj7va9du2p37ektj3k9c2lcdcqt78hqs3",
			"reward": {
				"denom": "ucmdx",
				"amount": "1403"
			}
		},
		{
			"address": "comdex19he6yj9j3fkl62jgkjdjpeq6tmap9nlr49lj5m",
			"reward": {
				"denom": "ucmdx",
				"amount": "3524"
			}
		},
		{
			"address": "comdex19h6yep42ps759mcr8kj90t4l4kwz8fhw6e3m8m",
			"reward": {
				"denom": "ucmdx",
				"amount": "1957"
			}
		},
		{
			"address": "comdex19havvttw7je3rmhyg8fgxdsywpkq8hek6yrvam",
			"reward": {
				"denom": "ucmdx",
				"amount": "323"
			}
		},
		{
			"address": "comdex19h74sazj0du3jqhlc6z8xlqx9hwmwmwnasrtz0",
			"reward": {
				"denom": "ucmdx",
				"amount": "2018"
			}
		},
		{
			"address": "comdex19cqwtnner2jlhnwhsdwwgwnkgvy97gfmhnhk0h",
			"reward": {
				"denom": "ucmdx",
				"amount": "4721"
			}
		},
		{
			"address": "comdex19cphf4mapvq9athkahry6a7l6haykcz7h5chyp",
			"reward": {
				"denom": "ucmdx",
				"amount": "94"
			}
		},
		{
			"address": "comdex19c9qe0zdtn2x6pc7skx5e5r86vjmex3wrggdq6",
			"reward": {
				"denom": "ucmdx",
				"amount": "14171"
			}
		},
		{
			"address": "comdex19cxlfvlpvztv62fr850gwdw6ttq32elvrz66q5",
			"reward": {
				"denom": "ucmdx",
				"amount": "205133"
			}
		},
		{
			"address": "comdex19c8mx8khjguz79gvavaere0c3n9c9j27e2frn6",
			"reward": {
				"denom": "ucmdx",
				"amount": "125052"
			}
		},
		{
			"address": "comdex19cgshdteaedfskxwycpja6cf4jlruyqy4hevma",
			"reward": {
				"denom": "ucmdx",
				"amount": "1488"
			}
		},
		{
			"address": "comdex19cddafns88ldqh3hp0fpsvxv76zulfxpcsrd83",
			"reward": {
				"denom": "ucmdx",
				"amount": "2002"
			}
		},
		{
			"address": "comdex19cwpg980fcs2e2s3l4xe9frvnl7a3vkeasw8w4",
			"reward": {
				"denom": "ucmdx",
				"amount": "978"
			}
		},
		{
			"address": "comdex19cwz3lusnqaqtyhxjk3gx9rul5qzxv9e5aafj5",
			"reward": {
				"denom": "ucmdx",
				"amount": "5098"
			}
		},
		{
			"address": "comdex19cw7sv2u6wmmnax52cjusl5fjkqw9p3azucnqs",
			"reward": {
				"denom": "ucmdx",
				"amount": "14394"
			}
		},
		{
			"address": "comdex19cnfek09q4uyqer92vcelyusvest6flhm878yr",
			"reward": {
				"denom": "ucmdx",
				"amount": "718"
			}
		},
		{
			"address": "comdex19c5lk03sx8d5n9wmkj8l6547gj35xwpqzt999h",
			"reward": {
				"denom": "ucmdx",
				"amount": "881"
			}
		},
		{
			"address": "comdex19ch732jmd57qcdyywx0q9zs87rkzudmj6qz8g5",
			"reward": {
				"denom": "ucmdx",
				"amount": "166"
			}
		},
		{
			"address": "comdex19ccsgjwp7f5umk0wuw5wqa4uld8mag7s5dq56l",
			"reward": {
				"denom": "ucmdx",
				"amount": "176"
			}
		},
		{
			"address": "comdex19cuhagpv8h3pt5espht77lathkhvzfvlh867ue",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex19caz5cgq4cwqdmfealqrq4uw5j56q8maltx7l6",
			"reward": {
				"denom": "ucmdx",
				"amount": "14379"
			}
		},
		{
			"address": "comdex19c706j05qh2ush03p9wdyudxx08cll6k4f3u2d",
			"reward": {
				"denom": "ucmdx",
				"amount": "1795"
			}
		},
		{
			"address": "comdex19c7ed7wuqmqd7sp790xagarfy8z7qvvup4tq5t",
			"reward": {
				"denom": "ucmdx",
				"amount": "1931"
			}
		},
		{
			"address": "comdex19c7av8chqutlhqncwch8ge7yynj0qqdqg6zrzq",
			"reward": {
				"denom": "ucmdx",
				"amount": "14"
			}
		},
		{
			"address": "comdex19epu6wwfulyxzcr0vu7xsh8efzv99hyf9zppgj",
			"reward": {
				"denom": "ucmdx",
				"amount": "1418"
			}
		},
		{
			"address": "comdex19e8qthuzl5tv98p5zzv8qjnmfgrz59xh4kz96h",
			"reward": {
				"denom": "ucmdx",
				"amount": "100"
			}
		},
		{
			"address": "comdex19egxezvj67ehjt060k3zxsfw7zwkrwehf0dlgl",
			"reward": {
				"denom": "ucmdx",
				"amount": "30135"
			}
		},
		{
			"address": "comdex19egvq66s6ktygvu3yd7xzzp4ffqedvw3e5rhcn",
			"reward": {
				"denom": "ucmdx",
				"amount": "324"
			}
		},
		{
			"address": "comdex19e2670v3akpgnj4wetq2g4c954zkpkt8r240mk",
			"reward": {
				"denom": "ucmdx",
				"amount": "52642"
			}
		},
		{
			"address": "comdex19etsf3gmt7z5cff06ts2v08cxn9kq25xwehssj",
			"reward": {
				"denom": "ucmdx",
				"amount": "22221"
			}
		},
		{
			"address": "comdex19edwnyvhle7v77pzjfdmzt38at6npm397ruu95",
			"reward": {
				"denom": "ucmdx",
				"amount": "12404"
			}
		},
		{
			"address": "comdex19ewc99penj77073slqslqvpwuyqhuxkuzw2zqm",
			"reward": {
				"denom": "ucmdx",
				"amount": "6946"
			}
		},
		{
			"address": "comdex19eslc5x5pawk6nts97zsa8mptal4t5skgsr39w",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex19enhux802sxs7rnnp5hkade39r4el7gn5lk8gm",
			"reward": {
				"denom": "ucmdx",
				"amount": "102151"
			}
		},
		{
			"address": "comdex19e59mx70cn56yeu8d90vxrgpkv7cztc0vsp6nt",
			"reward": {
				"denom": "ucmdx",
				"amount": "71121"
			}
		},
		{
			"address": "comdex19e4l3hnkzp0wzz2zmngpkqtdwt49zlgp82q4kt",
			"reward": {
				"denom": "ucmdx",
				"amount": "353"
			}
		},
		{
			"address": "comdex19ekdg5lnkk36a7gdcyv6238467zrvp0mag5yew",
			"reward": {
				"denom": "ucmdx",
				"amount": "13"
			}
		},
		{
			"address": "comdex19ehzs08e2e3myvcrnne435trt4svvxna6ztjk4",
			"reward": {
				"denom": "ucmdx",
				"amount": "167"
			}
		},
		{
			"address": "comdex19ecraycqchyg2utz4lgry3htcs5mpqz339mttn",
			"reward": {
				"denom": "ucmdx",
				"amount": "417"
			}
		},
		{
			"address": "comdex19ecxrrsr3qrjzgncdlvqclkxz707ta5cw0v0nu",
			"reward": {
				"denom": "ucmdx",
				"amount": "2480"
			}
		},
		{
			"address": "comdex19em3jrcrj4z909k9pp7tvkzae2cyymcx2h8mpp",
			"reward": {
				"denom": "ucmdx",
				"amount": "7228"
			}
		},
		{
			"address": "comdex19eu22zmzy5nz7cry2m20cqcky9sfm3nd0j0qku",
			"reward": {
				"denom": "ucmdx",
				"amount": "441"
			}
		},
		{
			"address": "comdex19earn0aynmtysgcheew5tcc5al66jkfehrlzvq",
			"reward": {
				"denom": "ucmdx",
				"amount": "17177"
			}
		},
		{
			"address": "comdex19eayv4ahr44lynh7ytzrl0gma4m794c7gwarka",
			"reward": {
				"denom": "ucmdx",
				"amount": "762"
			}
		},
		{
			"address": "comdex196p5wwcazp36s3nps6cw50am6tu7l0ucqvwnr9",
			"reward": {
				"denom": "ucmdx",
				"amount": "55"
			}
		},
		{
			"address": "comdex196xye6wdmy3xd7pjp0np6p3yj4d5adaq3djnhx",
			"reward": {
				"denom": "ucmdx",
				"amount": "44358"
			}
		},
		{
			"address": "comdex19684tn079l8rxhg0gprgcv52tk04a8x32u5xwt",
			"reward": {
				"denom": "ucmdx",
				"amount": "17451"
			}
		},
		{
			"address": "comdex1968mh9rmp0p29k6yzqxn6se4fca23882pmadx3",
			"reward": {
				"denom": "ucmdx",
				"amount": "1386"
			}
		},
		{
			"address": "comdex196gqwgmkyx8c4rgctvcv2zkm7kwfur3wspc848",
			"reward": {
				"denom": "ucmdx",
				"amount": "14143"
			}
		},
		{
			"address": "comdex196vjhefllcdvmpeyjezvctmve58ahenl0l5r66",
			"reward": {
				"denom": "ucmdx",
				"amount": "22120"
			}
		},
		{
			"address": "comdex196h0j26h5gpaz9n9hy7d30x9z7zfy9572g7m7s",
			"reward": {
				"denom": "ucmdx",
				"amount": "411"
			}
		},
		{
			"address": "comdex1966lqpaz0am0ynw22zfh02apquhhjgc6jcay66",
			"reward": {
				"denom": "ucmdx",
				"amount": "1664"
			}
		},
		{
			"address": "comdex196ud7fvtxtupp694sdmfv79d4jzssxun6xeu77",
			"reward": {
				"denom": "ucmdx",
				"amount": "54"
			}
		},
		{
			"address": "comdex19676a2qy59ttvcsuayqeyutfu5qx4jmnuhzzf3",
			"reward": {
				"denom": "ucmdx",
				"amount": "11468"
			}
		},
		{
			"address": "comdex19mr7h9feje33949f64yzks5s99xyaz5a8kzzj4",
			"reward": {
				"denom": "ucmdx",
				"amount": "5066"
			}
		},
		{
			"address": "comdex19mgx4f80wr4kpk0cun02z2kjz2tx3sucjl4pqs",
			"reward": {
				"denom": "ucmdx",
				"amount": "1224"
			}
		},
		{
			"address": "comdex19mf50s7zfv543fer6z8dzhuecj7nfdxfaq0e4a",
			"reward": {
				"denom": "ucmdx",
				"amount": "90620"
			}
		},
		{
			"address": "comdex19mfm9aa5shp59agt675nwrl9lf8n3lds68l37l",
			"reward": {
				"denom": "ucmdx",
				"amount": "145"
			}
		},
		{
			"address": "comdex19ms9rgvu38vr7zenw7npyynfvzv0mdwvdr797m",
			"reward": {
				"denom": "ucmdx",
				"amount": "1518"
			}
		},
		{
			"address": "comdex19mng8f95pwek6ey8ufjvp7yayhy9w94h9udr2d",
			"reward": {
				"denom": "ucmdx",
				"amount": "5282"
			}
		},
		{
			"address": "comdex19mk9lhs76a5h4s84qsfdjggt33xh7xz9gxvlwv",
			"reward": {
				"denom": "ucmdx",
				"amount": "179"
			}
		},
		{
			"address": "comdex19mc5kzqfqfn8h2fay90lpaadfchtmexvhqe0qw",
			"reward": {
				"denom": "ucmdx",
				"amount": "300"
			}
		},
		{
			"address": "comdex19mm44yzrfnugszn3mfe9yv34yw2y6zr78qpe3e",
			"reward": {
				"denom": "ucmdx",
				"amount": "3525"
			}
		},
		{
			"address": "comdex19mlujahncgnugxu05qaun6fwvpw2kcpwgww5r6",
			"reward": {
				"denom": "ucmdx",
				"amount": "1708"
			}
		},
		{
			"address": "comdex19upjme8fv5dvrajld8s9ka4h6062s3a3g74gzg",
			"reward": {
				"denom": "ucmdx",
				"amount": "35"
			}
		},
		{
			"address": "comdex19upnxgegz6k6t6e9fz65zpyum52s828ewl37p0",
			"reward": {
				"denom": "ucmdx",
				"amount": "57467"
			}
		},
		{
			"address": "comdex19uzwxj056a5jzqnyn8jshyh4s2g0ja77me7vtg",
			"reward": {
				"denom": "ucmdx",
				"amount": "3865"
			}
		},
		{
			"address": "comdex19uraxqsucmfvzvdhjhu2ycypqfh2z8k0d2mrxs",
			"reward": {
				"denom": "ucmdx",
				"amount": "58075"
			}
		},
		{
			"address": "comdex19u9xyt8m56e2x6hlpm8r33dps6m2sn09jtjnk7",
			"reward": {
				"denom": "ucmdx",
				"amount": "6826"
			}
		},
		{
			"address": "comdex19uxydn67lvxy0d2y3j3fqgj7d73dkkqt9pnrss",
			"reward": {
				"denom": "ucmdx",
				"amount": "31599"
			}
		},
		{
			"address": "comdex19ufpuv5jvesaxzqp8pdrdth4q3425zy93yuec9",
			"reward": {
				"denom": "ucmdx",
				"amount": "721"
			}
		},
		{
			"address": "comdex19utman0zhrc9qjty7tyju6wc3vw6qy5l85vt2p",
			"reward": {
				"denom": "ucmdx",
				"amount": "1068"
			}
		},
		{
			"address": "comdex19u0e2cd2au6lqgs67yrljyx98fqu92rvqh3g9f",
			"reward": {
				"denom": "ucmdx",
				"amount": "1430"
			}
		},
		{
			"address": "comdex19uj7gj8vgag5ds6csgalypy7x2vus907ct8jzk",
			"reward": {
				"denom": "ucmdx",
				"amount": "16"
			}
		},
		{
			"address": "comdex19u5vncf5zrl2na038crqek75h0mn0gzdfjn5e8",
			"reward": {
				"denom": "ucmdx",
				"amount": "1744"
			}
		},
		{
			"address": "comdex19u554mel2ur2g9z03xwh8mn762ugtcsyryxyx5",
			"reward": {
				"denom": "ucmdx",
				"amount": "513"
			}
		},
		{
			"address": "comdex19u42agwyk60dluv0ky36m6xmvgqj546qry7yjn",
			"reward": {
				"denom": "ucmdx",
				"amount": "2055"
			}
		},
		{
			"address": "comdex19uhj7kknqrzn763zwtmumpy20n08pezmcsc45g",
			"reward": {
				"denom": "ucmdx",
				"amount": "25612"
			}
		},
		{
			"address": "comdex19u7jlrhhga3eg7tmyk5pduzwsmpanf6e4k38xm",
			"reward": {
				"denom": "ucmdx",
				"amount": "1336"
			}
		},
		{
			"address": "comdex19arxx6afjjqttz3nutdyn09wr3qx8xujnm207m",
			"reward": {
				"denom": "ucmdx",
				"amount": "1978"
			}
		},
		{
			"address": "comdex19ayrkzcfpgw8ht7cugmcmm4ca96zwp5adyej06",
			"reward": {
				"denom": "ucmdx",
				"amount": "7582"
			}
		},
		{
			"address": "comdex19a87dkp33euhwhcuzef20h0gj5qtjr32nc46v9",
			"reward": {
				"denom": "ucmdx",
				"amount": "833"
			}
		},
		{
			"address": "comdex19a2ryujvk9xmzpdhz9g8mx9237usmgks0zc9zn",
			"reward": {
				"denom": "ucmdx",
				"amount": "1411"
			}
		},
		{
			"address": "comdex19a266ukkyjl72f0uve274f5kl34tdul7hm47yk",
			"reward": {
				"denom": "ucmdx",
				"amount": "345"
			}
		},
		{
			"address": "comdex19atggawer8mzqm282dgjat052j7nq76requze4",
			"reward": {
				"denom": "ucmdx",
				"amount": "14172"
			}
		},
		{
			"address": "comdex19awsl9264ljm7nxs0dy8hzg8guehrdu2f0hpx4",
			"reward": {
				"denom": "ucmdx",
				"amount": "33377"
			}
		},
		{
			"address": "comdex19ashlhtzzgmkq8l3dlaku4ukgn3mh0200k8qnd",
			"reward": {
				"denom": "ucmdx",
				"amount": "540"
			}
		},
		{
			"address": "comdex19ajaqtf7s6dzyr8c3e0zljgvvr3vn3ra6j6nz2",
			"reward": {
				"denom": "ucmdx",
				"amount": "329"
			}
		},
		{
			"address": "comdex19ajlxk4m03ncvrjzf6vu223v4fkwwflf44kqmd",
			"reward": {
				"denom": "ucmdx",
				"amount": "169709"
			}
		},
		{
			"address": "comdex19a5rmtgzaapj4zjhtr48ktnu8zhr7dpd8f50z8",
			"reward": {
				"denom": "ucmdx",
				"amount": "14340"
			}
		},
		{
			"address": "comdex19a4uev2tpx8kl37cgn6qdh03s5avf8cjldyf2y",
			"reward": {
				"denom": "ucmdx",
				"amount": "3552"
			}
		},
		{
			"address": "comdex19aufqnga9980f0hlpd2ulplu9erl35v3luj4xy",
			"reward": {
				"denom": "ucmdx",
				"amount": "4948"
			}
		},
		{
			"address": "comdex197q8kc2k6cc65cuan5yr0acwjtnf3gpc0s3fxt",
			"reward": {
				"denom": "ucmdx",
				"amount": "8636"
			}
		},
		{
			"address": "comdex197pp4slzf7pfvhvz7rkt3cuqqsc65fpxhz62pt",
			"reward": {
				"denom": "ucmdx",
				"amount": "4382"
			}
		},
		{
			"address": "comdex1979jfkvguhnwhmz2a9yaqqsavtukw7m87u6h04",
			"reward": {
				"denom": "ucmdx",
				"amount": "3537"
			}
		},
		{
			"address": "comdex197xzuu0mn0grzlrljjkpjqzlwcxznz3g4dah9c",
			"reward": {
				"denom": "ucmdx",
				"amount": "379"
			}
		},
		{
			"address": "comdex197x2emaletn274tczl9r4h3cclxkeap7lfn4h3",
			"reward": {
				"denom": "ucmdx",
				"amount": "143"
			}
		},
		{
			"address": "comdex197x0z8m097k76f9wsxnghjx0ngnu4szy7qns7d",
			"reward": {
				"denom": "ucmdx",
				"amount": "1503"
			}
		},
		{
			"address": "comdex19783t2qhzr64t420ekf03vwf5xtkc6082uwprv",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1972sugdnnxm4txq7hy540x94speal7kauvax8z",
			"reward": {
				"denom": "ucmdx",
				"amount": "4"
			}
		},
		{
			"address": "comdex1972aw784ukkx62w6vn4n7ek0y4qkl79elw0myw",
			"reward": {
				"denom": "ucmdx",
				"amount": "6187"
			}
		},
		{
			"address": "comdex197v60d6uaz4h2rw5qm06ylc38vg8ttdg3yj3h9",
			"reward": {
				"denom": "ucmdx",
				"amount": "144"
			}
		},
		{
			"address": "comdex197dq6ehzafs7vqwa73yhewyktg3xt4ucj34e7a",
			"reward": {
				"denom": "ucmdx",
				"amount": "59008"
			}
		},
		{
			"address": "comdex1970q655k7ss0fcuuv3kxnpxfl8mj2jfr3jnv76",
			"reward": {
				"denom": "ucmdx",
				"amount": "301"
			}
		},
		{
			"address": "comdex1973qhhswmsf0y99h88t40v7ecgzqhs6rc6fvq6",
			"reward": {
				"denom": "ucmdx",
				"amount": "7108"
			}
		},
		{
			"address": "comdex197nvjpgu79qsw27zzk0gyt9qh8cmq8wd852wt9",
			"reward": {
				"denom": "ucmdx",
				"amount": "17391"
			}
		},
		{
			"address": "comdex1975r8a82mvgjc6m8qm4k9ggw40ezjujptref3h",
			"reward": {
				"denom": "ucmdx",
				"amount": "300"
			}
		},
		{
			"address": "comdex19740m96hg4hn5u3vjm5rsfzt62kwz7f0j3a68w",
			"reward": {
				"denom": "ucmdx",
				"amount": "39814"
			}
		},
		{
			"address": "comdex197hljncjx27jrsy4yqfs7ka754uklt8yd0d0up",
			"reward": {
				"denom": "ucmdx",
				"amount": "184"
			}
		},
		{
			"address": "comdex197azlmrp8wvz4446jfnxu605snegznl46emneh",
			"reward": {
				"denom": "ucmdx",
				"amount": "456"
			}
		},
		{
			"address": "comdex197l8v50am7tp49qr3qf49f05f3ey6hmjk4sdgh",
			"reward": {
				"denom": "ucmdx",
				"amount": "12833"
			}
		},
		{
			"address": "comdex197lkxfc9rqf88zfj5zchtsekt9nv4qv06ade25",
			"reward": {
				"denom": "ucmdx",
				"amount": "197696"
			}
		},
		{
			"address": "comdex19lqhe9rr2zgjheqttg9kg33herf2d89aaz5as3",
			"reward": {
				"denom": "ucmdx",
				"amount": "13531"
			}
		},
		{
			"address": "comdex19l9uzenafg0r8atqune7g8e8mf4x6p2shxrzqa",
			"reward": {
				"denom": "ucmdx",
				"amount": "70"
			}
		},
		{
			"address": "comdex19lg0ty6qhmtk3fa87ng7lykus9g2s6540esnhx",
			"reward": {
				"denom": "ucmdx",
				"amount": "58310"
			}
		},
		{
			"address": "comdex19lg3uc00p9g9g5n96qlsxndgkt0glf5fanumst",
			"reward": {
				"denom": "ucmdx",
				"amount": "5679"
			}
		},
		{
			"address": "comdex19lf9spad2ddfshvf3p09erm390st6rstsvyjth",
			"reward": {
				"denom": "ucmdx",
				"amount": "28577"
			}
		},
		{
			"address": "comdex19lf0pxgrqaeuscgxm6gku9efsxa0fdfvhgjr5u",
			"reward": {
				"denom": "ucmdx",
				"amount": "24402"
			}
		},
		{
			"address": "comdex19ld2f77e525qy0ru69k6mtqhdja365v7dsny32",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex19l0gugs8n0ku3mva4n306xe85xpx5myfj8fx8m",
			"reward": {
				"denom": "ucmdx",
				"amount": "2017"
			}
		},
		{
			"address": "comdex19l399vr5ngtk4854xsqclrrs2gsnfpkuqpe2kt",
			"reward": {
				"denom": "ucmdx",
				"amount": "6785"
			}
		},
		{
			"address": "comdex19lhycf64lrnje76nhkexccmk3m9rqhzp44y4v7",
			"reward": {
				"denom": "ucmdx",
				"amount": "60"
			}
		},
		{
			"address": "comdex19lh5jkpm4y8y5070v5s7a3ftrpq56y8qrhhkl4",
			"reward": {
				"denom": "ucmdx",
				"amount": "5321"
			}
		},
		{
			"address": "comdex19lezakyp0kwsr97kmpvz56qpjvmet26fsyjwa3",
			"reward": {
				"denom": "ucmdx",
				"amount": "182517"
			}
		},
		{
			"address": "comdex19len5ya8xln30j5wmwk4q40ap3tv67a7ttat7c",
			"reward": {
				"denom": "ucmdx",
				"amount": "19635"
			}
		},
		{
			"address": "comdex19l6s67qqsdsw3jqty83gfw3p7fkagr8j28z25r",
			"reward": {
				"denom": "ucmdx",
				"amount": "5839"
			}
		},
		{
			"address": "comdex19lupd7xmlw6wc3r2qk8jsdvzrkluu05eey039x",
			"reward": {
				"denom": "ucmdx",
				"amount": "16"
			}
		},
		{
			"address": "comdex19luzvyfuudfm3guaukf9mfmnl6pgve7r9fvjht",
			"reward": {
				"denom": "ucmdx",
				"amount": "35"
			}
		},
		{
			"address": "comdex1xqqfpgqftkcnmrerwz8pgag2shfy9cwwgfvgyp",
			"reward": {
				"denom": "ucmdx",
				"amount": "17726"
			}
		},
		{
			"address": "comdex1xqpf37andd90quhjh8nqe8dsm8tdywqpym2dg9",
			"reward": {
				"denom": "ucmdx",
				"amount": "207"
			}
		},
		{
			"address": "comdex1xqzvt80c9zgnjkv2wd5kv4g6d3feutqd587zzc",
			"reward": {
				"denom": "ucmdx",
				"amount": "1944"
			}
		},
		{
			"address": "comdex1xqzc497frc7372xy25zvnqth4hyexvaxh7f0fm",
			"reward": {
				"denom": "ucmdx",
				"amount": "390"
			}
		},
		{
			"address": "comdex1xqfd0lr8w6kqutx4re73z74kt29r2hldsvt5gx",
			"reward": {
				"denom": "ucmdx",
				"amount": "151"
			}
		},
		{
			"address": "comdex1xqfwfvfret2cvm7mrpayda0xsx0v75zy360hvr",
			"reward": {
				"denom": "ucmdx",
				"amount": "5400"
			}
		},
		{
			"address": "comdex1xq2r0t4548t3qzks67m8svzj9f0q4d0flpe9ae",
			"reward": {
				"denom": "ucmdx",
				"amount": "1242"
			}
		},
		{
			"address": "comdex1xqtqcv02gdj4mpz5c33mrlkk09ss6sst4yedn9",
			"reward": {
				"denom": "ucmdx",
				"amount": "4769"
			}
		},
		{
			"address": "comdex1xqtlv8sap0vd3anrvpcf2wn6xv9xp02u2jnfgt",
			"reward": {
				"denom": "ucmdx",
				"amount": "1515"
			}
		},
		{
			"address": "comdex1xqvauc9myhlr87g80usn3vdsc5rdv5g00zdg4q",
			"reward": {
				"denom": "ucmdx",
				"amount": "946"
			}
		},
		{
			"address": "comdex1xqdpmpc2ayhgtydv0artuw0quvftp0uzz33kw7",
			"reward": {
				"denom": "ucmdx",
				"amount": "3355"
			}
		},
		{
			"address": "comdex1xqdn3ax3q82ujd920s4xu8w0y00jtlanc7lsqr",
			"reward": {
				"denom": "ucmdx",
				"amount": "7082"
			}
		},
		{
			"address": "comdex1xqw28cnhsxxpvqw0cl9lhr2z83hk2hjv52ks85",
			"reward": {
				"denom": "ucmdx",
				"amount": "10181"
			}
		},
		{
			"address": "comdex1xqwmypx2ty3puscujv3lklmsgyzkr5yds03e9t",
			"reward": {
				"denom": "ucmdx",
				"amount": "4313"
			}
		},
		{
			"address": "comdex1xq3khjktepxwyakdy67yrszk5tw9tpu82v92a9",
			"reward": {
				"denom": "ucmdx",
				"amount": "3523"
			}
		},
		{
			"address": "comdex1xqjup8fwt75afeyygrl3e5hskcsza54l0x6kmy",
			"reward": {
				"denom": "ucmdx",
				"amount": "156"
			}
		},
		{
			"address": "comdex1xq5dx8hy24rkkj6jh5sfkgxp5w8mdxa9udf9gf",
			"reward": {
				"denom": "ucmdx",
				"amount": "144"
			}
		},
		{
			"address": "comdex1xq44le24qhtc8v630cz0nquc9dan2sc5njevcp",
			"reward": {
				"denom": "ucmdx",
				"amount": "10037"
			}
		},
		{
			"address": "comdex1xqklf332j92wrmkpqyy6m00fkrq2w09xqg5x2l",
			"reward": {
				"denom": "ucmdx",
				"amount": "13047"
			}
		},
		{
			"address": "comdex1xqhdznzvq48g957h2p3hf5vj0e3lnddtdjqd32",
			"reward": {
				"denom": "ucmdx",
				"amount": "1240"
			}
		},
		{
			"address": "comdex1xqcxrwjfh0vay46gxj2a5wz6d8dxn3sugrsrd6",
			"reward": {
				"denom": "ucmdx",
				"amount": "12908"
			}
		},
		{
			"address": "comdex1xqmz5mcvltz380gu2g7mv88d0unpx9hn2090rs",
			"reward": {
				"denom": "ucmdx",
				"amount": "1732"
			}
		},
		{
			"address": "comdex1xqm2037e2879ewudegdw4y4nksthjrnjfaaeyz",
			"reward": {
				"denom": "ucmdx",
				"amount": "31634"
			}
		},
		{
			"address": "comdex1xqa3quymmjrke94s9uyjxpylqjjyuucuvl4fwt",
			"reward": {
				"denom": "ucmdx",
				"amount": "193570"
			}
		},
		{
			"address": "comdex1xpp8t3lehymp2gpnwmq7l7qpeuxckzwv3wqy3t",
			"reward": {
				"denom": "ucmdx",
				"amount": "1232"
			}
		},
		{
			"address": "comdex1xppjm3k3r2k8q844xh2tta6kx4c3v58n0a2hff",
			"reward": {
				"denom": "ucmdx",
				"amount": "8454"
			}
		},
		{
			"address": "comdex1xpruy6aaxw2nuw2shfg0uukcahst4h9gekj5ec",
			"reward": {
				"denom": "ucmdx",
				"amount": "6992"
			}
		},
		{
			"address": "comdex1xpyqxhyg742e4p9323676mqqzwe730e9el5ym0",
			"reward": {
				"denom": "ucmdx",
				"amount": "4728"
			}
		},
		{
			"address": "comdex1xp87a3ezugt0vshkjsrjwkfalgl3eyxk442pwk",
			"reward": {
				"denom": "ucmdx",
				"amount": "1415"
			}
		},
		{
			"address": "comdex1xpt4ynpf5ggpskj70l2mypdkgxz3mjrgzzcvpm",
			"reward": {
				"denom": "ucmdx",
				"amount": "2939"
			}
		},
		{
			"address": "comdex1xpw85kkrhua8f2h22utpactalqvcr27gmq737u",
			"reward": {
				"denom": "ucmdx",
				"amount": "2532"
			}
		},
		{
			"address": "comdex1xpwg3ks46q48u38828ammfa5serpj9qrhr9ptd",
			"reward": {
				"denom": "ucmdx",
				"amount": "1933"
			}
		},
		{
			"address": "comdex1xp0wea925y9dg25j3cvmsr4emjm6j239674ep9",
			"reward": {
				"denom": "ucmdx",
				"amount": "142"
			}
		},
		{
			"address": "comdex1xps00eq8fw5p3wjjmjfuar50kymc2qm8av89ks",
			"reward": {
				"denom": "ucmdx",
				"amount": "167"
			}
		},
		{
			"address": "comdex1xp3gq2e6cjpfqdefahq9da0l77z9lr0j0lk5xk",
			"reward": {
				"denom": "ucmdx",
				"amount": "17285"
			}
		},
		{
			"address": "comdex1xp3sp6q42v92xe4unuszc7h5ujt0772ayua47m",
			"reward": {
				"denom": "ucmdx",
				"amount": "114"
			}
		},
		{
			"address": "comdex1xpjenqlc2x8e67c9zcjp2cqskxwdzm992l5tna",
			"reward": {
				"denom": "ucmdx",
				"amount": "6352"
			}
		},
		{
			"address": "comdex1xph75n6d70cg0ppgykrzjj2s5y366fn92xm6ze",
			"reward": {
				"denom": "ucmdx",
				"amount": "1244"
			}
		},
		{
			"address": "comdex1xpe0fr4elxw00a0evs3l7kgqrzlv52ahsgesmr",
			"reward": {
				"denom": "ucmdx",
				"amount": "43946"
			}
		},
		{
			"address": "comdex1xp6j5kxgv5zl3qw5cw72s48cq3rdh6xpzczq76",
			"reward": {
				"denom": "ucmdx",
				"amount": "2874"
			}
		},
		{
			"address": "comdex1xpmpusnfk9qdrfpupnmjjyfatkqqz47sc3drj8",
			"reward": {
				"denom": "ucmdx",
				"amount": "148"
			}
		},
		{
			"address": "comdex1xpumk23t6jr6fxqr2g7fnt2m5hppryx93d2q3p",
			"reward": {
				"denom": "ucmdx",
				"amount": "1747"
			}
		},
		{
			"address": "comdex1xpa4hxmx0dmccqd78hqc5mnwtkhxqg3pva0t8c",
			"reward": {
				"denom": "ucmdx",
				"amount": "1234"
			}
		},
		{
			"address": "comdex1xp7psm2knh9459f6vyr4l7n57v56v5z0njxlqy",
			"reward": {
				"denom": "ucmdx",
				"amount": "2025"
			}
		},
		{
			"address": "comdex1xp78ulqa42yx8kma9lvzpn76s6y4a6cctu0juh",
			"reward": {
				"denom": "ucmdx",
				"amount": "3690"
			}
		},
		{
			"address": "comdex1xzqhvaczuwj0munmu3eyrhuh5pwragmvx9v45g",
			"reward": {
				"denom": "ucmdx",
				"amount": "1738"
			}
		},
		{
			"address": "comdex1xzpsdlhsxz3su4g2egherrknlllnvsq42krmxe",
			"reward": {
				"denom": "ucmdx",
				"amount": "5381"
			}
		},
		{
			"address": "comdex1xzz0ukyv0rkwfvxvqp3wd5rkym6p8etkeuzfpn",
			"reward": {
				"denom": "ucmdx",
				"amount": "369"
			}
		},
		{
			"address": "comdex1xzvtw63unx8y5j4a3qa9uv24nqlyj8rzsws38q",
			"reward": {
				"denom": "ucmdx",
				"amount": "1787"
			}
		},
		{
			"address": "comdex1xzv3ga06kzxu7rd2uxkq5k952x5xgff7vxk5tj",
			"reward": {
				"denom": "ucmdx",
				"amount": "2840"
			}
		},
		{
			"address": "comdex1xzwaddt87764u2fr3hf085s7gndhh6fsftjkm0",
			"reward": {
				"denom": "ucmdx",
				"amount": "12538"
			}
		},
		{
			"address": "comdex1xzwlsw84j8me2w2l9agcjhdxkej8ct6q7dmwx8",
			"reward": {
				"denom": "ucmdx",
				"amount": "126"
			}
		},
		{
			"address": "comdex1xzswnkuru0ztndknr25m92p03wj3hmg7z8nav2",
			"reward": {
				"denom": "ucmdx",
				"amount": "61689"
			}
		},
		{
			"address": "comdex1xzsu93xry7wze2xfggukprhvpcq8vsrdav5slm",
			"reward": {
				"denom": "ucmdx",
				"amount": "28488"
			}
		},
		{
			"address": "comdex1xz3pkz2lmatfc6cn6mnz72nfqfm5mpfl7jl734",
			"reward": {
				"denom": "ucmdx",
				"amount": "3522"
			}
		},
		{
			"address": "comdex1xz3w6rxum2sa6j7058as8mgdqe2pza8f7ljq4f",
			"reward": {
				"denom": "ucmdx",
				"amount": "566"
			}
		},
		{
			"address": "comdex1xz436dmuxuzpdt4a08s0a7z4tmc7806cnkznw3",
			"reward": {
				"denom": "ucmdx",
				"amount": "1128222"
			}
		},
		{
			"address": "comdex1xz4jjhxu5yq8q3zewx82qhzz2hyh8a2fa7s6p8",
			"reward": {
				"denom": "ucmdx",
				"amount": "448"
			}
		},
		{
			"address": "comdex1xz45ng9wztng9fwvzlsakhlcrdzlfdtevrrd6l",
			"reward": {
				"denom": "ucmdx",
				"amount": "3020"
			}
		},
		{
			"address": "comdex1xzktgumtn7kn2g665w8ukuwtwxtlw8mxrf657x",
			"reward": {
				"denom": "ucmdx",
				"amount": "60300"
			}
		},
		{
			"address": "comdex1xzk30g4k6egf8rep0axcjm93psw82y6vr78k6g",
			"reward": {
				"denom": "ucmdx",
				"amount": "1759"
			}
		},
		{
			"address": "comdex1xzhk0fyjumvvmkk4x8cmcy0fn86x5wmm9gzd2h",
			"reward": {
				"denom": "ucmdx",
				"amount": "71706"
			}
		},
		{
			"address": "comdex1xze0l22e234jwcwrnnavdhhlwre9df0m9cav72",
			"reward": {
				"denom": "ucmdx",
				"amount": "526"
			}
		},
		{
			"address": "comdex1xze63acalxpkj03lscc5ns66uhh9pedmk9p2cr",
			"reward": {
				"denom": "ucmdx",
				"amount": "8390"
			}
		},
		{
			"address": "comdex1xzu49233fstwwqswlsgjrrykq9dry29nhfl83u",
			"reward": {
				"denom": "ucmdx",
				"amount": "72645"
			}
		},
		{
			"address": "comdex1xzadnnz66q7pq424sstgdnlx95z894v5d748hq",
			"reward": {
				"denom": "ucmdx",
				"amount": "7190"
			}
		},
		{
			"address": "comdex1xz7c092kkmmfu2m6snt5j20hjevnr50uscchza",
			"reward": {
				"denom": "ucmdx",
				"amount": "697"
			}
		},
		{
			"address": "comdex1xr999hmrd3s6qaevcwqddc6wxly0pw9yzrqjpr",
			"reward": {
				"denom": "ucmdx",
				"amount": "38118"
			}
		},
		{
			"address": "comdex1xr9s7c4adydwezgafarrktfdjsm5j694qr8m0v",
			"reward": {
				"denom": "ucmdx",
				"amount": "883"
			}
		},
		{
			"address": "comdex1xr9ck0fcyupfetuuyqjm2hked6am42c5qu6zaz",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1xrd8qpm3p9qdz2f5zdns9z8hfagps5nwrnjwgq",
			"reward": {
				"denom": "ucmdx",
				"amount": "3959"
			}
		},
		{
			"address": "comdex1xrwpqtmkl7sdveutr3vna5yvauawvl9y2q8ytt",
			"reward": {
				"denom": "ucmdx",
				"amount": "8857"
			}
		},
		{
			"address": "comdex1xr0mfhue0fyqlzdx98ya55zvzhcwqwzmgpq5cg",
			"reward": {
				"denom": "ucmdx",
				"amount": "1618"
			}
		},
		{
			"address": "comdex1xrsjs7e569vcc6x3mdfhuwqc230wvs07305d67",
			"reward": {
				"denom": "ucmdx",
				"amount": "3743"
			}
		},
		{
			"address": "comdex1xrsj335h3euvpxnt0a320lns59t02ma2qg3u6w",
			"reward": {
				"denom": "ucmdx",
				"amount": "30938"
			}
		},
		{
			"address": "comdex1xr3pa5np3azmp32km3dfv624fw7km28hewtwd2",
			"reward": {
				"denom": "ucmdx",
				"amount": "32502"
			}
		},
		{
			"address": "comdex1xrj9hvrheluj5svvda7typzksw7uxut74g0j9a",
			"reward": {
				"denom": "ucmdx",
				"amount": "33172"
			}
		},
		{
			"address": "comdex1xr4r9nwrzaztww2y3kl3hmg8gehkgwzplt6fss",
			"reward": {
				"denom": "ucmdx",
				"amount": "23123"
			}
		},
		{
			"address": "comdex1xrkevy5vcftljgcsy60mwc5jd0jmj4g3kexhfn",
			"reward": {
				"denom": "ucmdx",
				"amount": "1020"
			}
		},
		{
			"address": "comdex1xruy9wenm8luk0388na7hvd3uysa549nawumdm",
			"reward": {
				"denom": "ucmdx",
				"amount": "7314"
			}
		},
		{
			"address": "comdex1xrad6az427fnde98d9u97wtf67reg8w88k6za7",
			"reward": {
				"denom": "ucmdx",
				"amount": "22208"
			}
		},
		{
			"address": "comdex1xrleyy97dkq6fyqwtxk6893wgfqmpn4f5a4yra",
			"reward": {
				"denom": "ucmdx",
				"amount": "14323"
			}
		},
		{
			"address": "comdex1xypexyphk0zspgrmdn2uc2f5elqdtj48qs72zt",
			"reward": {
				"denom": "ucmdx",
				"amount": "89"
			}
		},
		{
			"address": "comdex1xyzrmxm74xyyt0c85u264gkzfucuv368zjwdll",
			"reward": {
				"denom": "ucmdx",
				"amount": "66"
			}
		},
		{
			"address": "comdex1xy9d4zzl0ewj45mzha9ynu8syl8cs6x82d8cah",
			"reward": {
				"denom": "ucmdx",
				"amount": "401"
			}
		},
		{
			"address": "comdex1xyfnvx0dlv28jpc4freujfj2smw4cv5yr2s8j6",
			"reward": {
				"denom": "ucmdx",
				"amount": "109"
			}
		},
		{
			"address": "comdex1xy2r8u3a4ecs79lf472c39zzz0xng5axm5rdyp",
			"reward": {
				"denom": "ucmdx",
				"amount": "14023"
			}
		},
		{
			"address": "comdex1xysg279gkkeqs5yykqjtm3vyfmuz3tuhq52ucg",
			"reward": {
				"denom": "ucmdx",
				"amount": "6451"
			}
		},
		{
			"address": "comdex1xynszxq7pfc7c09n6s0c7hmefn875gu26lhuwg",
			"reward": {
				"denom": "ucmdx",
				"amount": "435"
			}
		},
		{
			"address": "comdex1xync4kuknalyx6xdh7h8p5wgdvrpdd6z8vduht",
			"reward": {
				"denom": "ucmdx",
				"amount": "192"
			}
		},
		{
			"address": "comdex1xylcv86u84gh7k9dssrvg4nak7u5f4h682rl02",
			"reward": {
				"denom": "ucmdx",
				"amount": "7537"
			}
		},
		{
			"address": "comdex1x9qmln79uh8z2v6xe4tr44kmvs0jtpz9ff8y7f",
			"reward": {
				"denom": "ucmdx",
				"amount": "14"
			}
		},
		{
			"address": "comdex1x9zpcuwvmacyzgr8888rfd7qnym7dm774cp05z",
			"reward": {
				"denom": "ucmdx",
				"amount": "14"
			}
		},
		{
			"address": "comdex1x9znuckerlnnfqtzpaenudntcln4qd7gv7yx3y",
			"reward": {
				"denom": "ucmdx",
				"amount": "1267"
			}
		},
		{
			"address": "comdex1x9x3yfq23ws4ks6xs6nkmechqflrj208ycn2yc",
			"reward": {
				"denom": "ucmdx",
				"amount": "40143"
			}
		},
		{
			"address": "comdex1x92x4ruwdh272y7p73ms297w9mrn0fwe27yht8",
			"reward": {
				"denom": "ucmdx",
				"amount": "169"
			}
		},
		{
			"address": "comdex1x9twy2gupqkcfn3eynaraym3h4p87hws7wtnen",
			"reward": {
				"denom": "ucmdx",
				"amount": "4074"
			}
		},
		{
			"address": "comdex1x93u4q54rvzz86vz5f5pjudejmlzzs3jqwej7n",
			"reward": {
				"denom": "ucmdx",
				"amount": "12182"
			}
		},
		{
			"address": "comdex1x945kqx3qyfnt4wln9m0dc4eedg8zyg66ee0lu",
			"reward": {
				"denom": "ucmdx",
				"amount": "215"
			}
		},
		{
			"address": "comdex1x9eftjw9nhnpaz44f655yxtkf4vgxqskzf9p6h",
			"reward": {
				"denom": "ucmdx",
				"amount": "268"
			}
		},
		{
			"address": "comdex1xxq4yyjd28kmx5d88wjva4rntlg98hvv978yst",
			"reward": {
				"denom": "ucmdx",
				"amount": "1784"
			}
		},
		{
			"address": "comdex1xxpwcxzz5h0g8vd4jxs429pfm8srt7dr7dlzf4",
			"reward": {
				"denom": "ucmdx",
				"amount": "16961"
			}
		},
		{
			"address": "comdex1xxzq7xkzkead285g6wv3j8xvtvs6wvyxfpsqkf",
			"reward": {
				"denom": "ucmdx",
				"amount": "1896"
			}
		},
		{
			"address": "comdex1xx2xqk5elvf0ry56jppgwg8twma02ygx2uaw32",
			"reward": {
				"denom": "ucmdx",
				"amount": "22360"
			}
		},
		{
			"address": "comdex1xx2xzlg7hlhf0229x2900mtrv5alufyn2f9vna",
			"reward": {
				"denom": "ucmdx",
				"amount": "1760"
			}
		},
		{
			"address": "comdex1xxwejyxdr73lnech5rrkgzuat8xduxyqpykfdc",
			"reward": {
				"denom": "ucmdx",
				"amount": "8937"
			}
		},
		{
			"address": "comdex1xx3srpv0x6ak7lxegvkjh3wug5drfdva4m422h",
			"reward": {
				"denom": "ucmdx",
				"amount": "3519"
			}
		},
		{
			"address": "comdex1xx33y2exg8s4vjusr4v2gye3upv4ethd6jky79",
			"reward": {
				"denom": "ucmdx",
				"amount": "1238"
			}
		},
		{
			"address": "comdex1xx335567pmatl7a2xm3g02xsj0fmxl2danazs8",
			"reward": {
				"denom": "ucmdx",
				"amount": "1762"
			}
		},
		{
			"address": "comdex1xxntzlrgekrk7nt9xdh5fk52c2kds2r0x9zq5s",
			"reward": {
				"denom": "ucmdx",
				"amount": "1787"
			}
		},
		{
			"address": "comdex1xxn3atej6vmxsx2pty2pd7kdvqc3fegnz73jd9",
			"reward": {
				"denom": "ucmdx",
				"amount": "969"
			}
		},
		{
			"address": "comdex1xx59zpfy65yyl47jfz8elr052c3n9ptfd25390",
			"reward": {
				"denom": "ucmdx",
				"amount": "8558"
			}
		},
		{
			"address": "comdex1xx4wkmhdxvt9xfe38rz23gp8gldw5egragy50h",
			"reward": {
				"denom": "ucmdx",
				"amount": "1775"
			}
		},
		{
			"address": "comdex1xxczyjnvxwp2708rmj4n5q2yccxemuxp0pu7ta",
			"reward": {
				"denom": "ucmdx",
				"amount": "5131"
			}
		},
		{
			"address": "comdex1xxen6w2fswjjqultamjgmrzmpcpz8dt3r9wf84",
			"reward": {
				"denom": "ucmdx",
				"amount": "60615"
			}
		},
		{
			"address": "comdex1xxeupzkflw8mu6045wxef2ft2hc0lleqpdavx2",
			"reward": {
				"denom": "ucmdx",
				"amount": "176"
			}
		},
		{
			"address": "comdex1xx6zyqf95948gxkcjcjdzy2xufl5p7mu7vpwqy",
			"reward": {
				"denom": "ucmdx",
				"amount": "57863"
			}
		},
		{
			"address": "comdex1xxlaanz4kltawzn4k8axydgmquns95eavh2uf7",
			"reward": {
				"denom": "ucmdx",
				"amount": "12310"
			}
		},
		{
			"address": "comdex1x8pt3hdf47fshryalrkvfcu8hzw0tly6mtaa0p",
			"reward": {
				"denom": "ucmdx",
				"amount": "5230"
			}
		},
		{
			"address": "comdex1x8rwal4g26mv3pkfs3fwcrvypfpmcvujhq88jp",
			"reward": {
				"denom": "ucmdx",
				"amount": "14544"
			}
		},
		{
			"address": "comdex1x8ywhue4y8x9hplstuxcklgj5aa36w6rc85z89",
			"reward": {
				"denom": "ucmdx",
				"amount": "1438"
			}
		},
		{
			"address": "comdex1x8958xm5el05p4thmg3vajv883z2w434h9r2z2",
			"reward": {
				"denom": "ucmdx",
				"amount": "1672"
			}
		},
		{
			"address": "comdex1x897tmle2eka82n0tly50m4un683dwdclus4js",
			"reward": {
				"denom": "ucmdx",
				"amount": "2715"
			}
		},
		{
			"address": "comdex1x8dl3mfhgr0sn66m0jlchhvyqljhv94yddjfmw",
			"reward": {
				"denom": "ucmdx",
				"amount": "5665"
			}
		},
		{
			"address": "comdex1x805l2asm9lslpca8zq7wx8a8cwysujkl0t3cq",
			"reward": {
				"denom": "ucmdx",
				"amount": "13458"
			}
		},
		{
			"address": "comdex1x8sf3s4j7h63l75dc9jd92d6z47lkxh3kgq8lw",
			"reward": {
				"denom": "ucmdx",
				"amount": "14097"
			}
		},
		{
			"address": "comdex1x8swu7ttzxjnvtszmh0mavdcr95xaymjqwcplu",
			"reward": {
				"denom": "ucmdx",
				"amount": "142"
			}
		},
		{
			"address": "comdex1x8s4ep3w6l3cgmeu45kxyrl0cg32qr2pyx4ugc",
			"reward": {
				"denom": "ucmdx",
				"amount": "1452"
			}
		},
		{
			"address": "comdex1x8jzm2nzq0xq6fv6rqvu4dwyl4wpy0dm946jtg",
			"reward": {
				"denom": "ucmdx",
				"amount": "182"
			}
		},
		{
			"address": "comdex1x8h069z68xtwg9ssf0gjvu80t4xwehx3gzksnp",
			"reward": {
				"denom": "ucmdx",
				"amount": "81"
			}
		},
		{
			"address": "comdex1x8cxgecujwqmskg2zkycnm4mrwnr57mfyy88jc",
			"reward": {
				"denom": "ucmdx",
				"amount": "35176"
			}
		},
		{
			"address": "comdex1x8lph9g94u84rkzfxkqwcja6dt8432mxdsn0t9",
			"reward": {
				"denom": "ucmdx",
				"amount": "179"
			}
		},
		{
			"address": "comdex1x8ltc5dr5h7pj9y44v9k9lyg8mcmeejh4pkhx2",
			"reward": {
				"denom": "ucmdx",
				"amount": "2489"
			}
		},
		{
			"address": "comdex1xgqqttzu3jh42q2h92pjtns7mjvxdlut5qd63u",
			"reward": {
				"denom": "ucmdx",
				"amount": "181"
			}
		},
		{
			"address": "comdex1xgptg4cxmstxj8yrpdut0fy6c43n4kqn7pfln5",
			"reward": {
				"denom": "ucmdx",
				"amount": "70"
			}
		},
		{
			"address": "comdex1xgzyc38k7tjckxd7f3556093v77u4jky6mysls",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1xgrus0wdywy7wkdqgaxcq3fdtt50zmhw542ldp",
			"reward": {
				"denom": "ucmdx",
				"amount": "2878"
			}
		},
		{
			"address": "comdex1xg9fnquzgfzg486h04qmntfrm64xufh40pk895",
			"reward": {
				"denom": "ucmdx",
				"amount": "28557"
			}
		},
		{
			"address": "comdex1xgxx0k47kzp349yw7hd33l35juzct66z8mzfe7",
			"reward": {
				"denom": "ucmdx",
				"amount": "1755"
			}
		},
		{
			"address": "comdex1xgxdq7cekxqvjq8vxzy9n6cm54plcn5kfqrpl9",
			"reward": {
				"denom": "ucmdx",
				"amount": "25"
			}
		},
		{
			"address": "comdex1xgx59wazg8t78danhppx7fl7tf6va4l63mhe67",
			"reward": {
				"denom": "ucmdx",
				"amount": "29"
			}
		},
		{
			"address": "comdex1xgxl5er2uvfkj0z5h46wdxkwa3pctuz5zaqqpd",
			"reward": {
				"denom": "ucmdx",
				"amount": "1486"
			}
		},
		{
			"address": "comdex1xgfww3pw0t0jdzdds8kmwgsknlmmhzpqcap0rd",
			"reward": {
				"denom": "ucmdx",
				"amount": "5245"
			}
		},
		{
			"address": "comdex1xg28h5wral9ccdda0u9z2zrh5w2jkuzjahn2k5",
			"reward": {
				"denom": "ucmdx",
				"amount": "16"
			}
		},
		{
			"address": "comdex1xg2a0yj5t2d9twrsp8x7lw3edwytfn3xle7g4j",
			"reward": {
				"denom": "ucmdx",
				"amount": "28147"
			}
		},
		{
			"address": "comdex1xgdx92trajgudzaj9y8e5gs8r307knmxw6edmc",
			"reward": {
				"denom": "ucmdx",
				"amount": "3343"
			}
		},
		{
			"address": "comdex1xg56degqnzufwd5nw0hc389nfpz886qszwttw6",
			"reward": {
				"denom": "ucmdx",
				"amount": "1103"
			}
		},
		{
			"address": "comdex1xg4q2p3mq72p9dtmgk8l65pq4w530j76e2gmjd",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1xgcwvcunndc9dtwr8ztuj6vnptjurm9m0ypkq6",
			"reward": {
				"denom": "ucmdx",
				"amount": "2530"
			}
		},
		{
			"address": "comdex1xguy6dv79gfsd98a5x75djy0z4us5dhhlskptr",
			"reward": {
				"denom": "ucmdx",
				"amount": "6839"
			}
		},
		{
			"address": "comdex1xglptwghr3ryrqyypyttkhchsfz5lt398d9lze",
			"reward": {
				"denom": "ucmdx",
				"amount": "13050"
			}
		},
		{
			"address": "comdex1xgl757g35nunclwsc4a42uup2eexth3qspa799",
			"reward": {
				"denom": "ucmdx",
				"amount": "20464"
			}
		},
		{
			"address": "comdex1xfq3kkj5hgyyfrqlvh9cj95z8nxx77r33kmcgt",
			"reward": {
				"denom": "ucmdx",
				"amount": "16"
			}
		},
		{
			"address": "comdex1xfq6utavfrznucp5a6cvqny8d9s0skepa432g5",
			"reward": {
				"denom": "ucmdx",
				"amount": "3587"
			}
		},
		{
			"address": "comdex1xfr3z2pvwf6mchdn5zx5uhqmekyudnjkqs8r9v",
			"reward": {
				"denom": "ucmdx",
				"amount": "275"
			}
		},
		{
			"address": "comdex1xf9xf6kyps79523zstg9usqxcutxantsj5hvfx",
			"reward": {
				"denom": "ucmdx",
				"amount": "58642"
			}
		},
		{
			"address": "comdex1xf2zpppdc6y38mn7785mkasgrmztxuajd09gqv",
			"reward": {
				"denom": "ucmdx",
				"amount": "151"
			}
		},
		{
			"address": "comdex1xfv3n4rhf3lfs69mpdcdjgjerqar6rt476udtc",
			"reward": {
				"denom": "ucmdx",
				"amount": "12266"
			}
		},
		{
			"address": "comdex1xf3x70azd3g8j8keh7mj77u6pzu5ehz0z6g5y5",
			"reward": {
				"denom": "ucmdx",
				"amount": "2035"
			}
		},
		{
			"address": "comdex1xf3fd4ew5l6vr96hs7f3qsx7hgfjplaegdazf4",
			"reward": {
				"denom": "ucmdx",
				"amount": "1222"
			}
		},
		{
			"address": "comdex1xfj2fzgq2xaag7xuy4y6yccdau6tkfv5rhyme6",
			"reward": {
				"denom": "ucmdx",
				"amount": "1737"
			}
		},
		{
			"address": "comdex1xfkg88s40w78jm3whdnfesvn2726tzymhs2uv4",
			"reward": {
				"denom": "ucmdx",
				"amount": "3567"
			}
		},
		{
			"address": "comdex1xfkltah9f43zups3tqaaud6dyefp2l3nlk6yrk",
			"reward": {
				"denom": "ucmdx",
				"amount": "4061"
			}
		},
		{
			"address": "comdex1xfczwn9nczhspyk98t7z5ydag5cs87u2dh0q8j",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1xfcr0fhww226ynd2yj6x0daxqgl55pda0h5k76",
			"reward": {
				"denom": "ucmdx",
				"amount": "16"
			}
		},
		{
			"address": "comdex1xfadke8nej6tzlhcukn3d3wkv6t76cdwp7nmgt",
			"reward": {
				"denom": "ucmdx",
				"amount": "3560"
			}
		},
		{
			"address": "comdex1xflcjsjhxx9wep27yncuurvxw8z23c08xf72g4",
			"reward": {
				"denom": "ucmdx",
				"amount": "3369"
			}
		},
		{
			"address": "comdex1xflluqvqlgmczj4wufalyw6dvy0n5neu3un8k7",
			"reward": {
				"denom": "ucmdx",
				"amount": "312"
			}
		},
		{
			"address": "comdex1x2prhkkcv69y5ujxufj0vl2uwkz3z9fw7v846z",
			"reward": {
				"denom": "ucmdx",
				"amount": "871"
			}
		},
		{
			"address": "comdex1x29vcwp6xretuwph2axafw7jemf6rdwd2zpeda",
			"reward": {
				"denom": "ucmdx",
				"amount": "1427"
			}
		},
		{
			"address": "comdex1x2xqqjsye4ea03pxd2urft5x5t7p2neqyuyawp",
			"reward": {
				"denom": "ucmdx",
				"amount": "148"
			}
		},
		{
			"address": "comdex1x28lg2kuy9gle08wsdpnhaag2907y850xad2w0",
			"reward": {
				"denom": "ucmdx",
				"amount": "16959"
			}
		},
		{
			"address": "comdex1x2g08slc8af8svz24pmlxdslhff76030uhqm7r",
			"reward": {
				"denom": "ucmdx",
				"amount": "3587"
			}
		},
		{
			"address": "comdex1x2fvgg7pjm6ew48x0h3sa760t63fj53zd57kjt",
			"reward": {
				"denom": "ucmdx",
				"amount": "24993"
			}
		},
		{
			"address": "comdex1x2teu6m75xpnt4aplpr98ywsecrlax32x6a55w",
			"reward": {
				"denom": "ucmdx",
				"amount": "53787"
			}
		},
		{
			"address": "comdex1x2vv3lqu5wrjnfnchke5mw0anv3v44cqy7t63m",
			"reward": {
				"denom": "ucmdx",
				"amount": "173"
			}
		},
		{
			"address": "comdex1x204fwguzz9c4qq2ee7g4yv26t9yw4sumt8fgg",
			"reward": {
				"denom": "ucmdx",
				"amount": "6424"
			}
		},
		{
			"address": "comdex1x24e6ust4dsdtewuh4vcxrqddgjfahz37gdy5f",
			"reward": {
				"denom": "ucmdx",
				"amount": "147488"
			}
		},
		{
			"address": "comdex1x2kqh6uavp79depu7eqmj3dv74mxpt9ngl2nt6",
			"reward": {
				"denom": "ucmdx",
				"amount": "509"
			}
		},
		{
			"address": "comdex1x2k9n836uastz2df7cjs6x6txtzv5jmy9r5m9w",
			"reward": {
				"denom": "ucmdx",
				"amount": "19"
			}
		},
		{
			"address": "comdex1x2mae67s52jwakxa7czyqkem3xyval0hnsk976",
			"reward": {
				"denom": "ucmdx",
				"amount": "269"
			}
		},
		{
			"address": "comdex1xtqp2zd2n2qwlxq0sf57krxd402zlu46vjn8nu",
			"reward": {
				"denom": "ucmdx",
				"amount": "215"
			}
		},
		{
			"address": "comdex1xtr8jq84cruu8scfuft04rpu7y58jd5cnk9nw9",
			"reward": {
				"denom": "ucmdx",
				"amount": "3704"
			}
		},
		{
			"address": "comdex1xt9pmsczktz7qww3udxhez50n9x2ekt0xeshzd",
			"reward": {
				"denom": "ucmdx",
				"amount": "1563"
			}
		},
		{
			"address": "comdex1xt8gjx303p0egj0cqzccaycumhesuvda7ll0gn",
			"reward": {
				"denom": "ucmdx",
				"amount": "1409"
			}
		},
		{
			"address": "comdex1xtgcv83v7ftpdtexzfqfyv73mfyenpa2jgh7lm",
			"reward": {
				"denom": "ucmdx",
				"amount": "571"
			}
		},
		{
			"address": "comdex1xtf8t3t7ths43yk6fwu57zyk6m5uj8gz4gddpr",
			"reward": {
				"denom": "ucmdx",
				"amount": "704"
			}
		},
		{
			"address": "comdex1xtty6sts47u628xpvw9lle4dw4evqh479vxacn",
			"reward": {
				"denom": "ucmdx",
				"amount": "8806"
			}
		},
		{
			"address": "comdex1xtvuq8a04qdwv9rt3rxs5d5uukxrj5unvp9thu",
			"reward": {
				"denom": "ucmdx",
				"amount": "2832"
			}
		},
		{
			"address": "comdex1xtwhy7madlrxwmupfw5a6v592tm5949w0f0j5x",
			"reward": {
				"denom": "ucmdx",
				"amount": "538"
			}
		},
		{
			"address": "comdex1xt05smu3wv35n4qydt5xrdwcqnucqvlwudx77y",
			"reward": {
				"denom": "ucmdx",
				"amount": "374054"
			}
		},
		{
			"address": "comdex1xtjnzwmf6gadsmw555tyl8dz05pcp9frnpfyp8",
			"reward": {
				"denom": "ucmdx",
				"amount": "104"
			}
		},
		{
			"address": "comdex1xt54085hc7wvgxlngesdpqhf4pguwulacmpnzc",
			"reward": {
				"denom": "ucmdx",
				"amount": "5963"
			}
		},
		{
			"address": "comdex1xtcxur9uenat3jp2cs2ux7d86knszsfz5trqa3",
			"reward": {
				"denom": "ucmdx",
				"amount": "1948"
			}
		},
		{
			"address": "comdex1xteuese3sw2xqemdax2kmey8gfqdga27y2syvp",
			"reward": {
				"denom": "ucmdx",
				"amount": "1142"
			}
		},
		{
			"address": "comdex1xt6z52s9n5zgjj2zursjyt949420shjmp28z6y",
			"reward": {
				"denom": "ucmdx",
				"amount": "170"
			}
		},
		{
			"address": "comdex1xtulwv9cze967j56ehgejdwduw9tzfy0dwnqu5",
			"reward": {
				"denom": "ucmdx",
				"amount": "4383"
			}
		},
		{
			"address": "comdex1xtazpmwpv3rtltcvlk0x5q43sx7y45qadt68qp",
			"reward": {
				"denom": "ucmdx",
				"amount": "2988"
			}
		},
		{
			"address": "comdex1xt77j9lkg63zxr4wk70paaf3gghwsu4r0gpjng",
			"reward": {
				"denom": "ucmdx",
				"amount": "140637"
			}
		},
		{
			"address": "comdex1xvqzm75rq8m63casdysufeu3aftt63pns76fkm",
			"reward": {
				"denom": "ucmdx",
				"amount": "5488"
			}
		},
		{
			"address": "comdex1xvq3a5ufwej95ne4e4dkcwdkw5zeudzyjsmtxk",
			"reward": {
				"denom": "ucmdx",
				"amount": "1751"
			}
		},
		{
			"address": "comdex1xvr6cjp9xec4ar87ecljp7mhwcqeuzr3eeqndg",
			"reward": {
				"denom": "ucmdx",
				"amount": "3541"
			}
		},
		{
			"address": "comdex1xv996lc2c8n2uvj5c90jgvd96ym6glh0v6frev",
			"reward": {
				"denom": "ucmdx",
				"amount": "204"
			}
		},
		{
			"address": "comdex1xvg0k075t6plunqhlyfw7ydm6p70gyhny45k95",
			"reward": {
				"denom": "ucmdx",
				"amount": "112"
			}
		},
		{
			"address": "comdex1xvg4zewj2kz7rnvxcgdckws6rhgt07h30ptf4x",
			"reward": {
				"denom": "ucmdx",
				"amount": "192"
			}
		},
		{
			"address": "comdex1xvvzrflekfkv649q94uhrmrlcnqp8pfqnf7yln",
			"reward": {
				"denom": "ucmdx",
				"amount": "2246"
			}
		},
		{
			"address": "comdex1xvdszgmhxpthv85p70n64tr3mtv80jfyxe5rs8",
			"reward": {
				"denom": "ucmdx",
				"amount": "6925"
			}
		},
		{
			"address": "comdex1xvd45767z3j27200n0czkayx5rmvnlud3hg5lh",
			"reward": {
				"denom": "ucmdx",
				"amount": "19941"
			}
		},
		{
			"address": "comdex1xv3htcp85yk78qpw58dm2vshvx0p89lf78hxnl",
			"reward": {
				"denom": "ucmdx",
				"amount": "12269"
			}
		},
		{
			"address": "comdex1xv3m3jshg0m8w8a233t7zr7t09zwaqqluk5x83",
			"reward": {
				"denom": "ucmdx",
				"amount": "297"
			}
		},
		{
			"address": "comdex1xvjdmslxgycpvr5ra4chaeyhmzx6jnuey85nvz",
			"reward": {
				"denom": "ucmdx",
				"amount": "3050"
			}
		},
		{
			"address": "comdex1xvkptnt95w9qtnzkta4uwhn62kh3h5zc3250mu",
			"reward": {
				"denom": "ucmdx",
				"amount": "210"
			}
		},
		{
			"address": "comdex1xvhm9f36l9yyv84sy2ke5gjjgmaeef52qtmj09",
			"reward": {
				"denom": "ucmdx",
				"amount": "27534"
			}
		},
		{
			"address": "comdex1xvm3crugf4rg84wm7hefz4v7dyc5caasu49ttk",
			"reward": {
				"denom": "ucmdx",
				"amount": "5509"
			}
		},
		{
			"address": "comdex1xva2jjjm4c3a27r2r28n2xa423xwe9rm34ltal",
			"reward": {
				"denom": "ucmdx",
				"amount": "11321"
			}
		},
		{
			"address": "comdex1xv7d5uq47cdf0pyf8gcerhg3gwp72jng2jph30",
			"reward": {
				"denom": "ucmdx",
				"amount": "16604"
			}
		},
		{
			"address": "comdex1xvlsa4kwj4d098vuzku4pcmumvh3femtnm4n0l",
			"reward": {
				"denom": "ucmdx",
				"amount": "1181"
			}
		},
		{
			"address": "comdex1xdzaygkwstu9tm2c5nwzct3qv89ne555m7u0jd",
			"reward": {
				"denom": "ucmdx",
				"amount": "2104"
			}
		},
		{
			"address": "comdex1xdgd0m4sg6yjsx0aeh0zg6hmwwyxlzehgq8rs2",
			"reward": {
				"denom": "ucmdx",
				"amount": "19604"
			}
		},
		{
			"address": "comdex1xd0wmw24c6v345tpkvcnhngkpvm23lu0knzdcc",
			"reward": {
				"denom": "ucmdx",
				"amount": "35036"
			}
		},
		{
			"address": "comdex1xd3k0dgzv5s2e7vetw99374qa05606pkhny7ed",
			"reward": {
				"denom": "ucmdx",
				"amount": "31742"
			}
		},
		{
			"address": "comdex1xd40eaa9cd2rhz7nq9u4m5hcrgnljaclq957gr",
			"reward": {
				"denom": "ucmdx",
				"amount": "873"
			}
		},
		{
			"address": "comdex1xdka5dczhqclc927egfcurc7thvasyus8s4tfv",
			"reward": {
				"denom": "ucmdx",
				"amount": "108138"
			}
		},
		{
			"address": "comdex1xdmy6usy3epu8ddl2ch2vk0269l2jz5xwwe9zc",
			"reward": {
				"denom": "ucmdx",
				"amount": "7527"
			}
		},
		{
			"address": "comdex1xduulk9tqgx93xj4npncdgpr525e36vtayq7g5",
			"reward": {
				"denom": "ucmdx",
				"amount": "1758"
			}
		},
		{
			"address": "comdex1xdavvafqz4pppx7netwts0m2dl4phpv5ewmuv9",
			"reward": {
				"denom": "ucmdx",
				"amount": "272"
			}
		},
		{
			"address": "comdex1xd7msw457czev06cawsm5h8hma8w5wrqufz8m3",
			"reward": {
				"denom": "ucmdx",
				"amount": "1506"
			}
		},
		{
			"address": "comdex1xwqsywzaz03zdjq2smysj3hp64fvq26lggnga5",
			"reward": {
				"denom": "ucmdx",
				"amount": "4560"
			}
		},
		{
			"address": "comdex1xwpgxsf39qv5zt2vva4fgsvrgmpxzh5a37j2ej",
			"reward": {
				"denom": "ucmdx",
				"amount": "512"
			}
		},
		{
			"address": "comdex1xwrp0zdqkxxwp0gqrzxhnq853cpcal9uh6f4qf",
			"reward": {
				"denom": "ucmdx",
				"amount": "1766"
			}
		},
		{
			"address": "comdex1xwwx4jz7uhdyjeal3gkuademcknwm03exqgmhy",
			"reward": {
				"denom": "ucmdx",
				"amount": "782"
			}
		},
		{
			"address": "comdex1xwszhh663ymaup824ufcr0d33v6jqznzx3jrcp",
			"reward": {
				"denom": "ucmdx",
				"amount": "5583"
			}
		},
		{
			"address": "comdex1xw5eypty5tves32wzrj635meesrv7gytzy00up",
			"reward": {
				"denom": "ucmdx",
				"amount": "204"
			}
		},
		{
			"address": "comdex1xw48qrrka0lkw9kmu4wq7vk730z4ga3p9qjeqx",
			"reward": {
				"denom": "ucmdx",
				"amount": "6733"
			}
		},
		{
			"address": "comdex1xwew6rey29tw44y59ayt522eszxahxt4fmzzs7",
			"reward": {
				"denom": "ucmdx",
				"amount": "23"
			}
		},
		{
			"address": "comdex1xwm26sgjh66hjsq8c9py6wucl5re6mdzu04wv6",
			"reward": {
				"denom": "ucmdx",
				"amount": "204"
			}
		},
		{
			"address": "comdex1xwu7jw8q9v50mrcfu3dfc84pdpx2qfeghzsr67",
			"reward": {
				"denom": "ucmdx",
				"amount": "69768"
			}
		},
		{
			"address": "comdex1x0zdfhz2hj9kt74cfp05n3gnn3v74t8tawgg94",
			"reward": {
				"denom": "ucmdx",
				"amount": "71990"
			}
		},
		{
			"address": "comdex1x0xuj55t8me68u390tu7e3mwyqrv30526chhy3",
			"reward": {
				"denom": "ucmdx",
				"amount": "151"
			}
		},
		{
			"address": "comdex1x0gy9tde8fffxexm5vsaa389aw9t5hq9fnjlmm",
			"reward": {
				"denom": "ucmdx",
				"amount": "78793"
			}
		},
		{
			"address": "comdex1x0tuap8wxe403x8e869nwr2kl534afh0c5sgs5",
			"reward": {
				"denom": "ucmdx",
				"amount": "87"
			}
		},
		{
			"address": "comdex1x0w0qqpysz4q5v6tv0fpay4nnrm3m2ju09xvxq",
			"reward": {
				"denom": "ucmdx",
				"amount": "6206"
			}
		},
		{
			"address": "comdex1x0w3cj289hxcgszncscpgrvmdauuv9lspwedzw",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1x00q68w364h7eza3grus94wvlkl0n8j72x9due",
			"reward": {
				"denom": "ucmdx",
				"amount": "2034"
			}
		},
		{
			"address": "comdex1x0slcaeanw0s6ttz7wt3arpavmaq3hcdqj9hmq",
			"reward": {
				"denom": "ucmdx",
				"amount": "1014"
			}
		},
		{
			"address": "comdex1x03xrh32rguw9f54w8503v25c87fnn07qurnuf",
			"reward": {
				"denom": "ucmdx",
				"amount": "47423"
			}
		},
		{
			"address": "comdex1x05e4czqcstg3f2c6yl9dsv5zwd8qt5zzy4u2u",
			"reward": {
				"denom": "ucmdx",
				"amount": "205"
			}
		},
		{
			"address": "comdex1x0hkutjturcpr8zn6wg7lampyt6s5gzjaxw0ud",
			"reward": {
				"denom": "ucmdx",
				"amount": "145"
			}
		},
		{
			"address": "comdex1x0u8207la8shnrldse4wnsfqmjudefjd39fma6",
			"reward": {
				"denom": "ucmdx",
				"amount": "95510"
			}
		},
		{
			"address": "comdex1xspr9c0e032zr8j4er79fw3vm85exsg2tgylg2",
			"reward": {
				"denom": "ucmdx",
				"amount": "583575"
			}
		},
		{
			"address": "comdex1xszxevp439dhy23gv40e764dm3a85p9yshkjqa",
			"reward": {
				"denom": "ucmdx",
				"amount": "524"
			}
		},
		{
			"address": "comdex1xszhf6fxfkhy9f760q4mw55pjuyumdn4ccq8sq",
			"reward": {
				"denom": "ucmdx",
				"amount": "4492"
			}
		},
		{
			"address": "comdex1xszafx3ejny3q98nm7qfnmkeczvd9qhz4g5fq0",
			"reward": {
				"denom": "ucmdx",
				"amount": "140"
			}
		},
		{
			"address": "comdex1xsrjd04hf43j2t34775xpx2x92ksv4z9atpmx0",
			"reward": {
				"denom": "ucmdx",
				"amount": "27513"
			}
		},
		{
			"address": "comdex1xsyyqtq9qryf705aqx58v7n8ndaupjayyx0zn2",
			"reward": {
				"denom": "ucmdx",
				"amount": "336"
			}
		},
		{
			"address": "comdex1xsy9pxr0g7txf3k4qj3zmdj7ms46av8tqvayfh",
			"reward": {
				"denom": "ucmdx",
				"amount": "26582"
			}
		},
		{
			"address": "comdex1xs9e279kk308uwelcl9m6mrv2xgttfdf0cvy22",
			"reward": {
				"denom": "ucmdx",
				"amount": "3409"
			}
		},
		{
			"address": "comdex1xs8dspc82ql7eat0vjk63dw8fsrs4dryyqt7h5",
			"reward": {
				"denom": "ucmdx",
				"amount": "13926"
			}
		},
		{
			"address": "comdex1xs8cu3pmc58fzlxqqz8pkghugy2hxh7rczaltg",
			"reward": {
				"denom": "ucmdx",
				"amount": "2018"
			}
		},
		{
			"address": "comdex1xsf5a29s89kft2gurs4wcjs07npy4lanhlz99m",
			"reward": {
				"denom": "ucmdx",
				"amount": "93301"
			}
		},
		{
			"address": "comdex1xs2yzlmswfmxenr850ctlng0aejd7lrhkkhgp8",
			"reward": {
				"denom": "ucmdx",
				"amount": "28440"
			}
		},
		{
			"address": "comdex1xs2whcrhrdumlrwfy7cfqztj66e3rn9dv43p2p",
			"reward": {
				"denom": "ucmdx",
				"amount": "51031"
			}
		},
		{
			"address": "comdex1xst4h9xfjuvunjxmsve7dasx6ujanqm2p3652p",
			"reward": {
				"denom": "ucmdx",
				"amount": "18076"
			}
		},
		{
			"address": "comdex1xssqr76amz0feeveux05fg9a9pxqjx6rj9fr8p",
			"reward": {
				"denom": "ucmdx",
				"amount": "481"
			}
		},
		{
			"address": "comdex1xs398n8swp0e8hnhxnhg27gz4w2u9qa9mdf06r",
			"reward": {
				"denom": "ucmdx",
				"amount": "12671"
			}
		},
		{
			"address": "comdex1xs3tfpx5esanwc4taz5n30md4633232auykmry",
			"reward": {
				"denom": "ucmdx",
				"amount": "28361"
			}
		},
		{
			"address": "comdex1xsnlq9mzrj40p2cya8pvupjqgquaatg7ufm4dp",
			"reward": {
				"denom": "ucmdx",
				"amount": "4977"
			}
		},
		{
			"address": "comdex1xsct98zczg6g8dss00rwwgp26k2t0f86x3rl9s",
			"reward": {
				"denom": "ucmdx",
				"amount": "46945"
			}
		},
		{
			"address": "comdex1xsesu2vkwfts0x3t06jkwrextyd9sdlzj4u7p8",
			"reward": {
				"denom": "ucmdx",
				"amount": "3377"
			}
		},
		{
			"address": "comdex1xseker45q4zt0zeqeu9cx5wydnreqje72s5q0c",
			"reward": {
				"denom": "ucmdx",
				"amount": "2801"
			}
		},
		{
			"address": "comdex1xs6qp5ajszhvvay7urxwk2rw8r9h4wrv46rnf3",
			"reward": {
				"denom": "ucmdx",
				"amount": "31678"
			}
		},
		{
			"address": "comdex1xsmzhhhpakf6c9qvazswq2p4mqkpuj25x7d5ul",
			"reward": {
				"denom": "ucmdx",
				"amount": "3634"
			}
		},
		{
			"address": "comdex1x3ppnxkzldlxf3l8n6gsvgaq6fvnrhum6xqd7q",
			"reward": {
				"denom": "ucmdx",
				"amount": "2877"
			}
		},
		{
			"address": "comdex1x3renzyuc03yn2zqer6hkqjrkm4ys4m4m7y4ss",
			"reward": {
				"denom": "ucmdx",
				"amount": "311"
			}
		},
		{
			"address": "comdex1x39mzput6ks5e47eapd62vzmllev7t9c0g7v2s",
			"reward": {
				"denom": "ucmdx",
				"amount": "878"
			}
		},
		{
			"address": "comdex1x3xfyyactagxda6k2wchukd2v76e44mnwau2vz",
			"reward": {
				"denom": "ucmdx",
				"amount": "1486"
			}
		},
		{
			"address": "comdex1x3gqmux5ht76lnpth3cgls79dfe7ew60fc9hc2",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1x32ehvfms20v8xh2akpcm9yd6kakfmxvspq4dr",
			"reward": {
				"denom": "ucmdx",
				"amount": "5838"
			}
		},
		{
			"address": "comdex1x3v3fnvnyx4v25jnsw6uuaq6y0xdtmc2nxkgr3",
			"reward": {
				"denom": "ucmdx",
				"amount": "2854"
			}
		},
		{
			"address": "comdex1x3vkrcnjw6y90rfyw84hwc67tz4y7z5rzw6urg",
			"reward": {
				"denom": "ucmdx",
				"amount": "508"
			}
		},
		{
			"address": "comdex1x3dvu2xr4cugpnxtj68amkf4eytdzvuhxd4g8e",
			"reward": {
				"denom": "ucmdx",
				"amount": "146105"
			}
		},
		{
			"address": "comdex1x3dcas6ld2c803sexm0nfppn25r4uc3gymnuqf",
			"reward": {
				"denom": "ucmdx",
				"amount": "205"
			}
		},
		{
			"address": "comdex1x33h20pmc8x0q0ge70sfy0x38eaycwcv6hxh8f",
			"reward": {
				"denom": "ucmdx",
				"amount": "786"
			}
		},
		{
			"address": "comdex1x35dqplgdtyavsuwusqeqclj2aujfvkv9j3j60",
			"reward": {
				"denom": "ucmdx",
				"amount": "1666"
			}
		},
		{
			"address": "comdex1x3krertljwsgfxvjkr5jq0ec3hy3s7npsx7xy6",
			"reward": {
				"denom": "ucmdx",
				"amount": "7066"
			}
		},
		{
			"address": "comdex1x3ksqdrqt4kq9v5f3rd8gtr0qtp7r3d28y5dr4",
			"reward": {
				"denom": "ucmdx",
				"amount": "19"
			}
		},
		{
			"address": "comdex1x3khp368gf9g464a8pd9f4x7hcmq9qwu7aqruw",
			"reward": {
				"denom": "ucmdx",
				"amount": "15200"
			}
		},
		{
			"address": "comdex1x3uempnf3qryzzq7a5l2vvhpad7l9lp44ex636",
			"reward": {
				"denom": "ucmdx",
				"amount": "6214"
			}
		},
		{
			"address": "comdex1x3apaxajum0qkka9fqxe9hhcydxdsxr2egqg84",
			"reward": {
				"denom": "ucmdx",
				"amount": "3952"
			}
		},
		{
			"address": "comdex1x3ltf7x3jsyu9ett0wc5c4ehe6ymxtlvnw9aak",
			"reward": {
				"denom": "ucmdx",
				"amount": "15322"
			}
		},
		{
			"address": "comdex1xjzk0nxv2yahgekwj6uh2je3698e85c6cgljng",
			"reward": {
				"denom": "ucmdx",
				"amount": "88440"
			}
		},
		{
			"address": "comdex1xjzkaclgkglwhr8lc5yp3falz7ylh680fy559g",
			"reward": {
				"denom": "ucmdx",
				"amount": "131"
			}
		},
		{
			"address": "comdex1xj9emmyy6jl35apku7ngcl8q8ugmynk33vtzwy",
			"reward": {
				"denom": "ucmdx",
				"amount": "747"
			}
		},
		{
			"address": "comdex1xjxhkzd034nt0u7mdcexg3u3pt4slc2jwxu8yt",
			"reward": {
				"denom": "ucmdx",
				"amount": "300"
			}
		},
		{
			"address": "comdex1xj8mu7ez23se7hpfn3nvy9mjc7zpszhhh8s4q0",
			"reward": {
				"denom": "ucmdx",
				"amount": "1949"
			}
		},
		{
			"address": "comdex1xjfgaw24t5yl094l72r9wdz03tzhpdajpty6xj",
			"reward": {
				"denom": "ucmdx",
				"amount": "1417"
			}
		},
		{
			"address": "comdex1xj2ngw820uy3f4yrucua2d9r6f4t79dgzrx07w",
			"reward": {
				"denom": "ucmdx",
				"amount": "138352"
			}
		},
		{
			"address": "comdex1xjvhpmflhue38mc8j75gqj2h0fds9gux0dz4dh",
			"reward": {
				"denom": "ucmdx",
				"amount": "17376"
			}
		},
		{
			"address": "comdex1xj0frc2kkwd8lxcnjlz0erjza6xcg4ak6m4kq3",
			"reward": {
				"denom": "ucmdx",
				"amount": "2258"
			}
		},
		{
			"address": "comdex1xj0nqn5k5dc7wjqlafxcu59dtpf40g4t0n9kxc",
			"reward": {
				"denom": "ucmdx",
				"amount": "145"
			}
		},
		{
			"address": "comdex1xj05n93arxu70en9xj67pevxpdmuc3mxnv7jav",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1xjsgwvp0mcn8jq52wk3xdgf5cwdy889gcqd9a9",
			"reward": {
				"denom": "ucmdx",
				"amount": "14770"
			}
		},
		{
			"address": "comdex1xjn96jm0jjuhclp3x9r5kuncazaegk62gxrfgl",
			"reward": {
				"denom": "ucmdx",
				"amount": "5748"
			}
		},
		{
			"address": "comdex1xjn4kvxqvy4lwn9jxc78umx6yw0zelqsldaq53",
			"reward": {
				"denom": "ucmdx",
				"amount": "1766"
			}
		},
		{
			"address": "comdex1xj4sm7vdrzpxefx6m6ae2mxvvf00zhth74hn6w",
			"reward": {
				"denom": "ucmdx",
				"amount": "2925"
			}
		},
		{
			"address": "comdex1xjkgvt2y3ss4d9p4zewmsh0ec3lr4n9pp9tczj",
			"reward": {
				"denom": "ucmdx",
				"amount": "184"
			}
		},
		{
			"address": "comdex1xjc0qz42t06cg3ddzpc0zx6fvz8apsfgpnq79g",
			"reward": {
				"denom": "ucmdx",
				"amount": "1327"
			}
		},
		{
			"address": "comdex1xj6tdtk594mjrruucl4g39yvfhsahg8vzkcvjp",
			"reward": {
				"denom": "ucmdx",
				"amount": "22822"
			}
		},
		{
			"address": "comdex1xj6nu35emjgkf5mwsellvye02zmr7rnhzeaq9l",
			"reward": {
				"denom": "ucmdx",
				"amount": "7932"
			}
		},
		{
			"address": "comdex1xjmwl5qpv0wjl9j2favtv6wc34zn9ccarn95k9",
			"reward": {
				"denom": "ucmdx",
				"amount": "5734"
			}
		},
		{
			"address": "comdex1xjax2m6mj62j7j8p3vpltq5ywr2gzhcqlwedqq",
			"reward": {
				"denom": "ucmdx",
				"amount": "6883"
			}
		},
		{
			"address": "comdex1xj78trfgnqwacvhaheqqrusaklum6fana7ku9h",
			"reward": {
				"denom": "ucmdx",
				"amount": "19"
			}
		},
		{
			"address": "comdex1xnpg74jrldt9g8kdz5uwd9ycqxtqep5nujex8n",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1xnzchh09z6ad8edk3z6a7qp6450lr9zw592ekh",
			"reward": {
				"denom": "ucmdx",
				"amount": "8104"
			}
		},
		{
			"address": "comdex1xnruylz68pef69e00pp3zsucgmyfkykusmqc7r",
			"reward": {
				"denom": "ucmdx",
				"amount": "68423"
			}
		},
		{
			"address": "comdex1xnyxtnzwgklrp9xmvhntnjwgjaujedwqk5dm69",
			"reward": {
				"denom": "ucmdx",
				"amount": "712"
			}
		},
		{
			"address": "comdex1xn9sjdhgjjpppj5ugzqzs4pv0j2dx6y5e2zpsc",
			"reward": {
				"denom": "ucmdx",
				"amount": "14782"
			}
		},
		{
			"address": "comdex1xn8xpf7z3s9adl348fdpp96krzpcvcvftyedeg",
			"reward": {
				"denom": "ucmdx",
				"amount": "6908"
			}
		},
		{
			"address": "comdex1xn2fzx2fw7qlgalycv0huschmcxz46gc48w038",
			"reward": {
				"denom": "ucmdx",
				"amount": "1422"
			}
		},
		{
			"address": "comdex1xndg3futnh8xvyc8kyp85r6rrlhl65v2kwlzde",
			"reward": {
				"denom": "ucmdx",
				"amount": "217"
			}
		},
		{
			"address": "comdex1xnwehxs04hse7qmvctw4ar9h59flfzhqgsrxjx",
			"reward": {
				"denom": "ucmdx",
				"amount": "1561"
			}
		},
		{
			"address": "comdex1xnwuzhdrwvyuqf37gtmcnt5gapr9nr3alzf5jx",
			"reward": {
				"denom": "ucmdx",
				"amount": "6893"
			}
		},
		{
			"address": "comdex1xn008hn0xdatua03mqv3u3664m4mu3cju4htsy",
			"reward": {
				"denom": "ucmdx",
				"amount": "17505"
			}
		},
		{
			"address": "comdex1xn65504lrt7z83awkq72xc860xejz5zfddp4k0",
			"reward": {
				"denom": "ucmdx",
				"amount": "168"
			}
		},
		{
			"address": "comdex1xnmgww3clpue7ly68v97u6rfh0sswtfz83xysz",
			"reward": {
				"denom": "ucmdx",
				"amount": "1686"
			}
		},
		{
			"address": "comdex1xnav3tux3t03k5u4hyylg5y4nf9s6wk9zxgest",
			"reward": {
				"denom": "ucmdx",
				"amount": "173"
			}
		},
		{
			"address": "comdex1x5pk9da5evgwcjje6npyg222etpq9erse8xjat",
			"reward": {
				"denom": "ucmdx",
				"amount": "4557"
			}
		},
		{
			"address": "comdex1x5p7nla6hjhjvpwkw7j80hzlan2vu062y54kvw",
			"reward": {
				"denom": "ucmdx",
				"amount": "15"
			}
		},
		{
			"address": "comdex1x5yh3fkm7h7l5k0zzc6pal2wx8g4nknv7sem9h",
			"reward": {
				"denom": "ucmdx",
				"amount": "31572"
			}
		},
		{
			"address": "comdex1x5xwn3gl2amhyzuzhesk7zzz7rraadz6xp0tk8",
			"reward": {
				"denom": "ucmdx",
				"amount": "183"
			}
		},
		{
			"address": "comdex1x5x0lnsv5ftfdce7zhdwurhlt327f38e6p7vvu",
			"reward": {
				"denom": "ucmdx",
				"amount": "17908"
			}
		},
		{
			"address": "comdex1x5drfhpdhg4az8peus0ua9hqhsyxtsve0aj7ye",
			"reward": {
				"denom": "ucmdx",
				"amount": "3423"
			}
		},
		{
			"address": "comdex1x538g7pvy23q5h64fqlv7hlrh0rtq8r8m3kw0l",
			"reward": {
				"denom": "ucmdx",
				"amount": "1792"
			}
		},
		{
			"address": "comdex1x5jvr8qje7psymr7afjpl6an07qwl72fwnjkjk",
			"reward": {
				"denom": "ucmdx",
				"amount": "87605"
			}
		},
		{
			"address": "comdex1x55xayfkjdmt4jpl4llzg70am7d330xtvegzrm",
			"reward": {
				"denom": "ucmdx",
				"amount": "0"
			}
		},
		{
			"address": "comdex1x5c8fduzk75wju9744htk6nw67rphwkln652nw",
			"reward": {
				"denom": "ucmdx",
				"amount": "44632"
			}
		},
		{
			"address": "comdex1x5c0e98plmkzelz0dgq94vtl555sfdrjvqy3as",
			"reward": {
				"denom": "ucmdx",
				"amount": "1433"
			}
		},
		{
			"address": "comdex1x5c66t4awtqwn4knw25mnqrc0y0z29vtj0qhdz",
			"reward": {
				"denom": "ucmdx",
				"amount": "286"
			}
		},
		{
			"address": "comdex1x5awyzfgwjvncsvjwxfeu6gk9257jjskq7yjfd",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1x57lpzlhuk7ylc99dpw265asawexawfuwxqc3a",
			"reward": {
				"denom": "ucmdx",
				"amount": "69961"
			}
		},
		{
			"address": "comdex1x4pd34fk6xyfgl9qnhr9zxmh59v46lu0ctunlx",
			"reward": {
				"denom": "ucmdx",
				"amount": "991"
			}
		},
		{
			"address": "comdex1x4pah2r6t35um69m7u9dhdn8e47vwgqa3cakv7",
			"reward": {
				"denom": "ucmdx",
				"amount": "533"
			}
		},
		{
			"address": "comdex1x498jwwxf0wtfm8t73k79q8t8jckf8eldlrdq8",
			"reward": {
				"denom": "ucmdx",
				"amount": "54011"
			}
		},
		{
			"address": "comdex1x49ua0gvedmsf8ufhtjul6cwsgkl2t2nt6uhw9",
			"reward": {
				"denom": "ucmdx",
				"amount": "17996"
			}
		},
		{
			"address": "comdex1x4gj676l5lq9j2uehl0mw8pq3e0jl3cmh46fku",
			"reward": {
				"denom": "ucmdx",
				"amount": "1542"
			}
		},
		{
			"address": "comdex1x4vhnxkfrqvgq6ehftrtp5l0286gawvmmh2y2a",
			"reward": {
				"denom": "ucmdx",
				"amount": "6711"
			}
		},
		{
			"address": "comdex1x4dhg0wsnpdl54ggvzh9ledcev4kmpe40tlr2d",
			"reward": {
				"denom": "ucmdx",
				"amount": "1013"
			}
		},
		{
			"address": "comdex1x4wv58h2rz7u6w8nzq9zpjkcyg9rdx7s9v39pc",
			"reward": {
				"denom": "ucmdx",
				"amount": "44070"
			}
		},
		{
			"address": "comdex1x4w0akjh7wr67t9wkf9nx3x06kd808ymknze8k",
			"reward": {
				"denom": "ucmdx",
				"amount": "6992"
			}
		},
		{
			"address": "comdex1x4std3j6mnergxnpepxftug6csd0df8wqfdx43",
			"reward": {
				"denom": "ucmdx",
				"amount": "1247"
			}
		},
		{
			"address": "comdex1x4n8mstcd3pfs9mslhnkxalcquv0p8camy4utu",
			"reward": {
				"denom": "ucmdx",
				"amount": "80628"
			}
		},
		{
			"address": "comdex1x4kxq36ayde766jw0lag3vjjtu38knjf4z98e2",
			"reward": {
				"denom": "ucmdx",
				"amount": "528151"
			}
		},
		{
			"address": "comdex1x4hrvetcm7mzhk7xqt46tk8aq27cpclewwflrq",
			"reward": {
				"denom": "ucmdx",
				"amount": "362"
			}
		},
		{
			"address": "comdex1x4c5u2gmr0l7f7e2trpct6rmpn7d864qkpe77x",
			"reward": {
				"denom": "ucmdx",
				"amount": "1686"
			}
		},
		{
			"address": "comdex1x4eey5gll3a9j3cjalypfkuxx9h3u5ehc6vk02",
			"reward": {
				"denom": "ucmdx",
				"amount": "92228"
			}
		},
		{
			"address": "comdex1x46ff45ctkak6rjqka9uqvwqt4vhdr45v4fvx5",
			"reward": {
				"denom": "ucmdx",
				"amount": "2811"
			}
		},
		{
			"address": "comdex1x47rcwukde03fah62jhu237yqauwmsjpkfl7pj",
			"reward": {
				"denom": "ucmdx",
				"amount": "3732"
			}
		},
		{
			"address": "comdex1xkzz8dgdhusv5zuhtvg7ucupkv65jz9xhefr99",
			"reward": {
				"denom": "ucmdx",
				"amount": "615"
			}
		},
		{
			"address": "comdex1xkz84rwd7umd2mmsl97he56f3zzl3hkyg9c7gj",
			"reward": {
				"denom": "ucmdx",
				"amount": "173"
			}
		},
		{
			"address": "comdex1xk9nguk28krahklva6a25umjvcqm2jqh2udz6r",
			"reward": {
				"denom": "ucmdx",
				"amount": "2215"
			}
		},
		{
			"address": "comdex1xk3e07sevn9703vcz4tmgpz8akr9pdygs48e2p",
			"reward": {
				"denom": "ucmdx",
				"amount": "173"
			}
		},
		{
			"address": "comdex1xk5y0qls7675lccl4sqlc6f8x6lye66ykqg6qy",
			"reward": {
				"denom": "ucmdx",
				"amount": "88916"
			}
		},
		{
			"address": "comdex1xk56mj6qj9a20ku0cpe3azggjzkk2mcyjp4cne",
			"reward": {
				"denom": "ucmdx",
				"amount": "15"
			}
		},
		{
			"address": "comdex1xkevqtytapkwtpv6cysc3eqzg72gf228mm7ls6",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1xk793azjrrkxjg0j8ytkutwrpuczccujyrt5f4",
			"reward": {
				"denom": "ucmdx",
				"amount": "16630"
			}
		},
		{
			"address": "comdex1xklqdrulta9d49ehpvxwx57v0a3snal3v3wgq7",
			"reward": {
				"denom": "ucmdx",
				"amount": "35"
			}
		},
		{
			"address": "comdex1xhpmld45kettu77v4v4cjhe36368cw2ftsgfpy",
			"reward": {
				"denom": "ucmdx",
				"amount": "321"
			}
		},
		{
			"address": "comdex1xhzmsv4cjfffjm85vlg6cuzwmaz2nukvw07fep",
			"reward": {
				"denom": "ucmdx",
				"amount": "170264"
			}
		},
		{
			"address": "comdex1xhrm5s6ahj8x5cdwmxzrut8aa6v9hdygwueds6",
			"reward": {
				"denom": "ucmdx",
				"amount": "3030"
			}
		},
		{
			"address": "comdex1xhr7fy2ayv777njm22fd2ht92ma49yej3x6cwn",
			"reward": {
				"denom": "ucmdx",
				"amount": "2850"
			}
		},
		{
			"address": "comdex1xh8ydf9rqxtyg7txh56wrcysadk3htfl0ux2c9",
			"reward": {
				"denom": "ucmdx",
				"amount": "1423"
			}
		},
		{
			"address": "comdex1xhtgcayrej5aa22h5pdjl38kfjwc3q33dfrrss",
			"reward": {
				"denom": "ucmdx",
				"amount": "13389"
			}
		},
		{
			"address": "comdex1xhttwcewy0fdefgykd3hv49xjyllwu0v7gk3r6",
			"reward": {
				"denom": "ucmdx",
				"amount": "7103"
			}
		},
		{
			"address": "comdex1xhddf4ata79dx3dglurcnun4lqqqv2vufdtw55",
			"reward": {
				"denom": "ucmdx",
				"amount": "19349"
			}
		},
		{
			"address": "comdex1xhd5ytc36kacfwja70rdwnvwmuz6p9kvwjccsc",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1xhwwdnnrt9h2gm5nsq2jw3q94c35cv66ryp793",
			"reward": {
				"denom": "ucmdx",
				"amount": "934998744"
			}
		},
		{
			"address": "comdex1xhskf0pe345rk54trmvrs0ax7d68fz5ygg6d7q",
			"reward": {
				"denom": "ucmdx",
				"amount": "12900"
			}
		},
		{
			"address": "comdex1xh3zlehfkyk3fh9yg08hr3n0jc7kqfmcypna9g",
			"reward": {
				"denom": "ucmdx",
				"amount": "8277"
			}
		},
		{
			"address": "comdex1xhhquwctvwm40qn3nmmqq067pc2gw22ezdcl7t",
			"reward": {
				"denom": "ucmdx",
				"amount": "3036"
			}
		},
		{
			"address": "comdex1xhh87eldwxjrmqywrgpxz3yr7kkatdv8wtara9",
			"reward": {
				"denom": "ucmdx",
				"amount": "74"
			}
		},
		{
			"address": "comdex1xhcr47wulnmyat6udu406sffz0vhp4093rmvd2",
			"reward": {
				"denom": "ucmdx",
				"amount": "202"
			}
		},
		{
			"address": "comdex1xhc67efcl6n63rdsd5serwrpl08c8f4rjp0gn0",
			"reward": {
				"denom": "ucmdx",
				"amount": "8983"
			}
		},
		{
			"address": "comdex1xhec3fg9wqkmy0de0hftlzs3avwe3k3la9fg95",
			"reward": {
				"denom": "ucmdx",
				"amount": "5487"
			}
		},
		{
			"address": "comdex1xh6cvttdl6c7ym5jsumwmsz4x68kxa2vkwwuwt",
			"reward": {
				"denom": "ucmdx",
				"amount": "128666"
			}
		},
		{
			"address": "comdex1xhm48vf0tnp927w9afc3ckwhmymrhyzn0uexm5",
			"reward": {
				"denom": "ucmdx",
				"amount": "201"
			}
		},
		{
			"address": "comdex1xhudm0vj9sxksgd2nsunshxs7v9xajxym69d09",
			"reward": {
				"denom": "ucmdx",
				"amount": "1305"
			}
		},
		{
			"address": "comdex1xh7478l056aefvs8n75xxc4dm6axw73uqthafy",
			"reward": {
				"denom": "ucmdx",
				"amount": "61000"
			}
		},
		{
			"address": "comdex1xcqdzkl5wshey5s6s65stfk3fak8admfcqewsa",
			"reward": {
				"denom": "ucmdx",
				"amount": "41721"
			}
		},
		{
			"address": "comdex1xcquwffc0efn34pazvlsrasvnu5ctunk94x4mp",
			"reward": {
				"denom": "ucmdx",
				"amount": "3181"
			}
		},
		{
			"address": "comdex1xcpy6x238n94tws37nq2xfy0tuue8s3sjx7qw8",
			"reward": {
				"denom": "ucmdx",
				"amount": "17422"
			}
		},
		{
			"address": "comdex1xczgqukp36umln6j69rmz39uaxce4j8l40srf3",
			"reward": {
				"denom": "ucmdx",
				"amount": "13870"
			}
		},
		{
			"address": "comdex1xcymgnpuvgm7tp92qggcju50yj39rmn9vtavpk",
			"reward": {
				"denom": "ucmdx",
				"amount": "19044"
			}
		},
		{
			"address": "comdex1xc90xdyj56sna0q0acyg8fqx3kt348hm2y3tln",
			"reward": {
				"denom": "ucmdx",
				"amount": "16"
			}
		},
		{
			"address": "comdex1xcxg4ek7qvxxlrv7qc6jeulf72xydpfv3guemw",
			"reward": {
				"denom": "ucmdx",
				"amount": "10118"
			}
		},
		{
			"address": "comdex1xcxu29mdf06u734d6qkryzrp798d2fm02dlus3",
			"reward": {
				"denom": "ucmdx",
				"amount": "1190"
			}
		},
		{
			"address": "comdex1xc8tjcmhfdn7ld2qgcc6s92r3wsp2x6hrmjtt0",
			"reward": {
				"denom": "ucmdx",
				"amount": "23897"
			}
		},
		{
			"address": "comdex1xc8nvj9wrkrvcye0kku48tv4tzl0nx4qj4pwp4",
			"reward": {
				"denom": "ucmdx",
				"amount": "176"
			}
		},
		{
			"address": "comdex1xcgmya6ryvcjsdxk0ywy23660h8qsh59gm2xkr",
			"reward": {
				"denom": "ucmdx",
				"amount": "84185"
			}
		},
		{
			"address": "comdex1xc2t364d97d0ejsrm2h5ruv8yv7kfuycgafhfa",
			"reward": {
				"denom": "ucmdx",
				"amount": "101373"
			}
		},
		{
			"address": "comdex1xctkzs0rt9s6saghqjcy27twgakcjrm55nh6xk",
			"reward": {
				"denom": "ucmdx",
				"amount": "3537"
			}
		},
		{
			"address": "comdex1xcd2pt3yvrsea8z2hmamzw0uehvc56j5y3y5e3",
			"reward": {
				"denom": "ucmdx",
				"amount": "1024"
			}
		},
		{
			"address": "comdex1xc3vtmesgqasaawxa9wpwj7y69y2f750pd3zg5",
			"reward": {
				"denom": "ucmdx",
				"amount": "1778"
			}
		},
		{
			"address": "comdex1xcnymjd3zqlr3kjqvnw505x07fq3v0hnp62gmx",
			"reward": {
				"denom": "ucmdx",
				"amount": "618"
			}
		},
		{
			"address": "comdex1xc4ter7gj3t97780wp8a23ncnv74mhal78s4vq",
			"reward": {
				"denom": "ucmdx",
				"amount": "10126"
			}
		},
		{
			"address": "comdex1xc4e9qped5x77k92h378jt3sq7frnp2ka6dywg",
			"reward": {
				"denom": "ucmdx",
				"amount": "2957048"
			}
		},
		{
			"address": "comdex1xccnldlwmgpq4hc855f3nzvshfl5nw744sh3dk",
			"reward": {
				"denom": "ucmdx",
				"amount": "7044"
			}
		},
		{
			"address": "comdex1xcmjeagxd0clehkux8xwmpy8u9mnpshx69az6d",
			"reward": {
				"denom": "ucmdx",
				"amount": "173"
			}
		},
		{
			"address": "comdex1xcuk5lls8pf08k5glj578afau223tv30x7h4zx",
			"reward": {
				"denom": "ucmdx",
				"amount": "1744"
			}
		},
		{
			"address": "comdex1xc7zeetcgfsqut8s0putllatus9xnu05n02gmf",
			"reward": {
				"denom": "ucmdx",
				"amount": "7176"
			}
		},
		{
			"address": "comdex1xc74m3qhwalhvpges9pvhteesnhx924504akyh",
			"reward": {
				"denom": "ucmdx",
				"amount": "347"
			}
		},
		{
			"address": "comdex1xclsxf5wg2mmw670y7axes4607gsncctqxqu0e",
			"reward": {
				"denom": "ucmdx",
				"amount": "113587"
			}
		},
		{
			"address": "comdex1xe99zszd0zqzxc75kvky0hxvvzgkj3lyjeh4k3",
			"reward": {
				"denom": "ucmdx",
				"amount": "1483"
			}
		},
		{
			"address": "comdex1xe9ncg88upmpnz2hqfwscfke8vu9dlw0rem3hp",
			"reward": {
				"denom": "ucmdx",
				"amount": "2012"
			}
		},
		{
			"address": "comdex1xex5n5q54qh5n7vslz9tls6k543hah3h0ytj54",
			"reward": {
				"denom": "ucmdx",
				"amount": "718"
			}
		},
		{
			"address": "comdex1xe886g8ahj70tasyuslv2a4vvm2kq33wx7st22",
			"reward": {
				"denom": "ucmdx",
				"amount": "285"
			}
		},
		{
			"address": "comdex1xe82zalz55v4cj8u0tx7c8zqsyx80hccd96wlc",
			"reward": {
				"denom": "ucmdx",
				"amount": "10306"
			}
		},
		{
			"address": "comdex1xe8tdssrv8vqajtfxk5v9ekxmv7k38clstpjgw",
			"reward": {
				"denom": "ucmdx",
				"amount": "3840"
			}
		},
		{
			"address": "comdex1xe2q09s2vrectm337rujg42vv2hp4349zhfgqg",
			"reward": {
				"denom": "ucmdx",
				"amount": "115243"
			}
		},
		{
			"address": "comdex1xe26rtgu0qgkmhce00g5d79xfcjff2faspkdjy",
			"reward": {
				"denom": "ucmdx",
				"amount": "6044"
			}
		},
		{
			"address": "comdex1xevyrnyg2dc8qcwmfhs4pa5zr23fnpch44yvrp",
			"reward": {
				"denom": "ucmdx",
				"amount": "142"
			}
		},
		{
			"address": "comdex1xed5800ya5wl7pq4r90dapa45lpv5xq9euw2qr",
			"reward": {
				"denom": "ucmdx",
				"amount": "58053"
			}
		},
		{
			"address": "comdex1xe45au792a3v8newdlhndy6zf9yq667rr9j324",
			"reward": {
				"denom": "ucmdx",
				"amount": "2492"
			}
		},
		{
			"address": "comdex1xecaqps7memq5l5z546slk93zhxnx5dq5xu8dt",
			"reward": {
				"denom": "ucmdx",
				"amount": "19034"
			}
		},
		{
			"address": "comdex1xeeuzlzfmhdfyg6a9qnssnl4ys5ay3tc8t8a7u",
			"reward": {
				"denom": "ucmdx",
				"amount": "899"
			}
		},
		{
			"address": "comdex1xeadxg0ee2muwh3082wp8mth4cwqsmjky25yua",
			"reward": {
				"denom": "ucmdx",
				"amount": "3385"
			}
		},
		{
			"address": "comdex1xe7dzscg7csqt990yg6j2368txm2dpptqshc52",
			"reward": {
				"denom": "ucmdx",
				"amount": "1212"
			}
		},
		{
			"address": "comdex1xel5yvwzxerhewmlhckn2rdczk5uqqdmv8r86p",
			"reward": {
				"denom": "ucmdx",
				"amount": "1878"
			}
		},
		{
			"address": "comdex1x6rf8pyzvmxpc8ydxzdnxyvzrpytsh88t3aeg5",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1x69jls25lldyytfhzaaved078v72rc3j5tj5r7",
			"reward": {
				"denom": "ucmdx",
				"amount": "71637"
			}
		},
		{
			"address": "comdex1x68h9mgxsnzq9ewfyp43e7shg3xmkl4vlacc90",
			"reward": {
				"denom": "ucmdx",
				"amount": "887"
			}
		},
		{
			"address": "comdex1x6dc0rwuvr7fxdwsnyarw6gafddycvcjqd2w2d",
			"reward": {
				"denom": "ucmdx",
				"amount": "14509"
			}
		},
		{
			"address": "comdex1x60q98yxy8shyd9a4jj0e75esfqchjll2um4lm",
			"reward": {
				"denom": "ucmdx",
				"amount": "727"
			}
		},
		{
			"address": "comdex1x63tv5csgv6mpsckq90p6h3tlufuxv35svg0tm",
			"reward": {
				"denom": "ucmdx",
				"amount": "7007"
			}
		},
		{
			"address": "comdex1x65ng38pytnkkjnqpz7y052cpfplzuazgl2jlx",
			"reward": {
				"denom": "ucmdx",
				"amount": "5390"
			}
		},
		{
			"address": "comdex1x64yeuk900fe0quex4fqrpg2ttn8q7sz99ehr4",
			"reward": {
				"denom": "ucmdx",
				"amount": "353"
			}
		},
		{
			"address": "comdex1x6ksjsj05ypgg8yfpw8cc4ct809sx52krcn78l",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1x6eyg69r0ugj2rd2tumsfk4tjpavwv68n0ctw2",
			"reward": {
				"denom": "ucmdx",
				"amount": "20222"
			}
		},
		{
			"address": "comdex1x6enzu323dn8yfjcjgxyqe75tct9kks7na6dca",
			"reward": {
				"denom": "ucmdx",
				"amount": "25964"
			}
		},
		{
			"address": "comdex1x668vnse277adcujc7vlkz5qpdghzdnt3xw33u",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1x6m588ylgx4nedyw5x7v06v3hafqer2egcdlmq",
			"reward": {
				"denom": "ucmdx",
				"amount": "6179"
			}
		},
		{
			"address": "comdex1x6mlaeg7fg095dmmhc9azmwkuczvfk6r0g40d3",
			"reward": {
				"denom": "ucmdx",
				"amount": "1960"
			}
		},
		{
			"address": "comdex1xmq55xuqmqk6k02m6f6ar2mrlpwslcqjx6dcg8",
			"reward": {
				"denom": "ucmdx",
				"amount": "1741"
			}
		},
		{
			"address": "comdex1xmzcs6uxy0jjaqkdac3xsvm35qkjjv37nd3mn2",
			"reward": {
				"denom": "ucmdx",
				"amount": "698512"
			}
		},
		{
			"address": "comdex1xm9ymlxpzzu6pe6kuzm60gpnj8cmc4ac6cxycx",
			"reward": {
				"denom": "ucmdx",
				"amount": "13166"
			}
		},
		{
			"address": "comdex1xm920fhqn4teg5pz9th30v8dzvmnv4488lh83u",
			"reward": {
				"denom": "ucmdx",
				"amount": "201"
			}
		},
		{
			"address": "comdex1xmx3s5qz8cggtga35ly74rvvq4p5wjjfc27h23",
			"reward": {
				"denom": "ucmdx",
				"amount": "1406"
			}
		},
		{
			"address": "comdex1xmx7ryxev6gwr86307y7305jlpsg4zwey6498k",
			"reward": {
				"denom": "ucmdx",
				"amount": "1440"
			}
		},
		{
			"address": "comdex1xm8fcwfaettsvufd80zfx5kahf67t3w04wt3vy",
			"reward": {
				"denom": "ucmdx",
				"amount": "7491"
			}
		},
		{
			"address": "comdex1xmgnwdpgslxstx0fv3mq23kqgwrrd35epq609q",
			"reward": {
				"denom": "ucmdx",
				"amount": "8855"
			}
		},
		{
			"address": "comdex1xm26944klp4ynssyn3dt5rztkd9lxfv0rwp2dt",
			"reward": {
				"denom": "ucmdx",
				"amount": "2324"
			}
		},
		{
			"address": "comdex1xmdq3khx9v9wgtr6ttynavh576tpam96rg234v",
			"reward": {
				"denom": "ucmdx",
				"amount": "367"
			}
		},
		{
			"address": "comdex1xmj9m30ndemg47dngxr042qjltgeedu4yn3mf6",
			"reward": {
				"denom": "ucmdx",
				"amount": "1486"
			}
		},
		{
			"address": "comdex1xm5mk3mhxfjrurx0ex768altry07sl4x059wf4",
			"reward": {
				"denom": "ucmdx",
				"amount": "1879"
			}
		},
		{
			"address": "comdex1xmklaad960c2shgv5c7kn4jsw3vyyc0kt7v6tu",
			"reward": {
				"denom": "ucmdx",
				"amount": "5065"
			}
		},
		{
			"address": "comdex1xmhfwc8m5mrzwfh7ea8hdxjpfft02z8ygftc2w",
			"reward": {
				"denom": "ucmdx",
				"amount": "179"
			}
		},
		{
			"address": "comdex1xmckehfvjpwfqm7ngamq6vlj7tum66u9gcee99",
			"reward": {
				"denom": "ucmdx",
				"amount": "8159"
			}
		},
		{
			"address": "comdex1xmclv8fkryrals9e8p9l6pht0sgja2nx9vxjj4",
			"reward": {
				"denom": "ucmdx",
				"amount": "709"
			}
		},
		{
			"address": "comdex1xmm3r3ers67ttnq4s8da4y2nu7sa2yw3jc4wfy",
			"reward": {
				"denom": "ucmdx",
				"amount": "712"
			}
		},
		{
			"address": "comdex1xmu76j2cy9q4fafxggy7846xvm2zkfw44vqycf",
			"reward": {
				"denom": "ucmdx",
				"amount": "11434"
			}
		},
		{
			"address": "comdex1xma8nyuuaejgs54krrpfrfyck66spcx870p7qm",
			"reward": {
				"denom": "ucmdx",
				"amount": "896761"
			}
		},
		{
			"address": "comdex1xm7lau0vas7elyxuyx9e6q4pg0uh3fwyugpaks",
			"reward": {
				"denom": "ucmdx",
				"amount": "1026"
			}
		},
		{
			"address": "comdex1xmlnhjnm4sdygxqwt7ewjxf67lnde0cky4g4xa",
			"reward": {
				"denom": "ucmdx",
				"amount": "32592"
			}
		},
		{
			"address": "comdex1xuyalg6qj9tpmx03ngy4jw983rjlua0077l9y0",
			"reward": {
				"denom": "ucmdx",
				"amount": "17049"
			}
		},
		{
			"address": "comdex1xu8nrezef7fnzqhzjgc8ers2mjjd436he6vav5",
			"reward": {
				"denom": "ucmdx",
				"amount": "67818"
			}
		},
		{
			"address": "comdex1xutnj32vxp8gz8x7txld3sdp808fquv823p2af",
			"reward": {
				"denom": "ucmdx",
				"amount": "613"
			}
		},
		{
			"address": "comdex1xuwmmz2r6kgzhlwdy726hve4wlqth2pcmva3yz",
			"reward": {
				"denom": "ucmdx",
				"amount": "181"
			}
		},
		{
			"address": "comdex1xus0kvfkav7h85ahmnxuxz7ujyyfrmknywxvkf",
			"reward": {
				"denom": "ucmdx",
				"amount": "5672"
			}
		},
		{
			"address": "comdex1xusmwr85xvl7vezpmv8gsm97svyfzczf76l6g4",
			"reward": {
				"denom": "ucmdx",
				"amount": "435"
			}
		},
		{
			"address": "comdex1xu3qn6sk4emmlma3qaj8gu5zwnszjysrdd4a89",
			"reward": {
				"denom": "ucmdx",
				"amount": "345"
			}
		},
		{
			"address": "comdex1xun3qq5veudusxys2rygxkduqld3yq9cltp53f",
			"reward": {
				"denom": "ucmdx",
				"amount": "13841"
			}
		},
		{
			"address": "comdex1xu5p7xjq7sc0t9kh5zrc9cqp8ve3p5eer79jgl",
			"reward": {
				"denom": "ucmdx",
				"amount": "174"
			}
		},
		{
			"address": "comdex1xu5fsc3jgcfwmr3a7uefcfs4r0u42q4cvsm33m",
			"reward": {
				"denom": "ucmdx",
				"amount": "30354"
			}
		},
		{
			"address": "comdex1xukzyulrkutrwefg035j99jqywgcnzejaahjw6",
			"reward": {
				"denom": "ucmdx",
				"amount": "189"
			}
		},
		{
			"address": "comdex1xuk8v9wghrhqktp2q4wvpnjj5l5jt2gy3djsc5",
			"reward": {
				"denom": "ucmdx",
				"amount": "7163"
			}
		},
		{
			"address": "comdex1xueq0l35yxegqk5g8kmq44mxqf3ld5l2dr62zm",
			"reward": {
				"denom": "ucmdx",
				"amount": "1654"
			}
		},
		{
			"address": "comdex1xuu9z949fxmy6r9pnwt6zkr55zvl9z5aqpfact",
			"reward": {
				"denom": "ucmdx",
				"amount": "888"
			}
		},
		{
			"address": "comdex1xazphqvfcdswkku476xee9uzjtqdjx0tp6g3pl",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1xart8ces8etr8gnan499860rf6mec3mzssnkdk",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1xar680e82wp3q35z7kaug3jcwrg4jxjp7lpuvu",
			"reward": {
				"denom": "ucmdx",
				"amount": "525"
			}
		},
		{
			"address": "comdex1xa9acgrvtrgqwupd83pmr9mtmawupeffwak3u5",
			"reward": {
				"denom": "ucmdx",
				"amount": "2116"
			}
		},
		{
			"address": "comdex1xa97zljdcazpjzgvy0d7qay7jfcfxaz5tx5tkk",
			"reward": {
				"denom": "ucmdx",
				"amount": "6174"
			}
		},
		{
			"address": "comdex1xa8h94ylkmfl33jr3d4zz8nppkz46jlx0ywfax",
			"reward": {
				"denom": "ucmdx",
				"amount": "84"
			}
		},
		{
			"address": "comdex1xagsu2e8retyrvrdjyf6608dn6tt5pewc3qka8",
			"reward": {
				"denom": "ucmdx",
				"amount": "144"
			}
		},
		{
			"address": "comdex1xata6hmxt5ql6ycquygeq6rjk7unc9k7ajfkuf",
			"reward": {
				"denom": "ucmdx",
				"amount": "5461"
			}
		},
		{
			"address": "comdex1xadfefqhfue3s54jrw2yl2fkchsktk4ukwllpt",
			"reward": {
				"denom": "ucmdx",
				"amount": "60749"
			}
		},
		{
			"address": "comdex1xadwk0ruwcq6cvd0qrk24tflrkrg3ndwqakh6d",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1xadkhcws8qnhnjet2vzqw8as8z3ul4h3w7w3eg",
			"reward": {
				"denom": "ucmdx",
				"amount": "179"
			}
		},
		{
			"address": "comdex1xadh02mh5230gkjzp9kckuqv2v36vurv04wewn",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1xadmlv23gm2mfr9hhylldwkffvfqwr8lxsq43w",
			"reward": {
				"denom": "ucmdx",
				"amount": "7537"
			}
		},
		{
			"address": "comdex1xa0g3j2562kksq5m4zqvufma2syvwqtf7czzhx",
			"reward": {
				"denom": "ucmdx",
				"amount": "2881"
			}
		},
		{
			"address": "comdex1xanyvw2yd9c9kkc45v42saqfhy8dk3ye0z24at",
			"reward": {
				"denom": "ucmdx",
				"amount": "179"
			}
		},
		{
			"address": "comdex1xacneta7w2yn9hvnagkuzzrz0j5pmfwf3zct4j",
			"reward": {
				"denom": "ucmdx",
				"amount": "7000"
			}
		},
		{
			"address": "comdex1xa6kpu9aektq0wefypw62sr7j85psjnl7tl8g6",
			"reward": {
				"denom": "ucmdx",
				"amount": "60155"
			}
		},
		{
			"address": "comdex1xaaxw8xzyxh6cflqd4g2xqmypsa4x8hzsn68ua",
			"reward": {
				"denom": "ucmdx",
				"amount": "1238"
			}
		},
		{
			"address": "comdex1x7qy5pl5r5843f5s8tsv9p8mc3h25spszp7ldf",
			"reward": {
				"denom": "ucmdx",
				"amount": "143"
			}
		},
		{
			"address": "comdex1x7qw2vnwx3ggzur7nxdfq0gs88r64fw4cs0lnz",
			"reward": {
				"denom": "ucmdx",
				"amount": "359"
			}
		},
		{
			"address": "comdex1x7prg55fwqnk7aj8e0lxfadrsde0kl85tc4rd9",
			"reward": {
				"denom": "ucmdx",
				"amount": "139"
			}
		},
		{
			"address": "comdex1x79hwup38e0pd66e6v070ra6kn5hz5y3jz8j4a",
			"reward": {
				"denom": "ucmdx",
				"amount": "6150"
			}
		},
		{
			"address": "comdex1x72wm0ktheec37unqj6da7mgzfddpt4xs956k0",
			"reward": {
				"denom": "ucmdx",
				"amount": "3637"
			}
		},
		{
			"address": "comdex1x7v8k46aqq2u53azhv3cnvtller6ljrj47tmgj",
			"reward": {
				"denom": "ucmdx",
				"amount": "37865"
			}
		},
		{
			"address": "comdex1x70ddy5et3xe8qk04xvfg43ehjh5vysg9avscm",
			"reward": {
				"denom": "ucmdx",
				"amount": "1994"
			}
		},
		{
			"address": "comdex1x7jxnyqaxy94mprg3t3anhh9s6vshw3uvakqu9",
			"reward": {
				"denom": "ucmdx",
				"amount": "1248"
			}
		},
		{
			"address": "comdex1x75tlkeshmlu2uzkvanhlqxaqqcz7dg8kd3sly",
			"reward": {
				"denom": "ucmdx",
				"amount": "61646"
			}
		},
		{
			"address": "comdex1x755asu299qwzp9z3f2r6886ceut99hjykxe7h",
			"reward": {
				"denom": "ucmdx",
				"amount": "200"
			}
		},
		{
			"address": "comdex1x7ktn95p0jwkh8ha56jn30yl33yy3ythtu5jew",
			"reward": {
				"denom": "ucmdx",
				"amount": "820"
			}
		},
		{
			"address": "comdex1x7hxgtx3n9qya3f64uqvkyhkvjumhlt2yu4p73",
			"reward": {
				"denom": "ucmdx",
				"amount": "204"
			}
		},
		{
			"address": "comdex1x76vl0y66k6pjj5jz9l7djyd7n7926s92v0yzq",
			"reward": {
				"denom": "ucmdx",
				"amount": "1789"
			}
		},
		{
			"address": "comdex1x7lff02j06jf9pdcca6fcmlyza3t33pwxlj6m4",
			"reward": {
				"denom": "ucmdx",
				"amount": "503"
			}
		},
		{
			"address": "comdex1xlrgm9eldzqmgf2mskr58gsanzsyssyjtnxlh4",
			"reward": {
				"denom": "ucmdx",
				"amount": "54"
			}
		},
		{
			"address": "comdex1xlykrvmztj07nkvw6aulkwahc2vyfqjrc37gyp",
			"reward": {
				"denom": "ucmdx",
				"amount": "143"
			}
		},
		{
			"address": "comdex1xl890retehlmgtkukn6y7l54wr84x7lfv9xtxd",
			"reward": {
				"denom": "ucmdx",
				"amount": "1765"
			}
		},
		{
			"address": "comdex1xl8tj209ncr39lggzfyxvvmzjrf37dfrxm97nz",
			"reward": {
				"denom": "ucmdx",
				"amount": "503"
			}
		},
		{
			"address": "comdex1xl8kl4ccalrrrt9h3sf2y44ft5ye3u9jw5r76l",
			"reward": {
				"denom": "ucmdx",
				"amount": "2346"
			}
		},
		{
			"address": "comdex1xl8etjf7k7dc0lh84ttm7uxv48ea7z7rashq2a",
			"reward": {
				"denom": "ucmdx",
				"amount": "1746"
			}
		},
		{
			"address": "comdex1xlfgkjdgrmwdpfmrd6u023tljxt6lzu0dyzx6z",
			"reward": {
				"denom": "ucmdx",
				"amount": "356"
			}
		},
		{
			"address": "comdex1xlsghhn4227an0d9wcwg9xp78wsq4kclf4rvpn",
			"reward": {
				"denom": "ucmdx",
				"amount": "75770"
			}
		},
		{
			"address": "comdex1xlsnsuvcw5a8hn6r3am48tmn8x6emm4h8dpmj9",
			"reward": {
				"denom": "ucmdx",
				"amount": "284"
			}
		},
		{
			"address": "comdex1xl3a0qhpk0kyudm7kkd2d0d2mr8z4q237s3ppv",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1xln4gfnrk7k545tnxt5qtgre9x2l3lhrz99q89",
			"reward": {
				"denom": "ucmdx",
				"amount": "1244"
			}
		},
		{
			"address": "comdex1xl5j5egqm63nnxluxtz6v7k7vxwp57lckdst2d",
			"reward": {
				"denom": "ucmdx",
				"amount": "849"
			}
		},
		{
			"address": "comdex1xl57czg4x0y0r85qejv0qjtgx4y5j39m4083v2",
			"reward": {
				"denom": "ucmdx",
				"amount": "279"
			}
		},
		{
			"address": "comdex1xlkz20pvfc7qhyc4rzrzscgtzrmu5z9jsyu9cd",
			"reward": {
				"denom": "ucmdx",
				"amount": "34897"
			}
		},
		{
			"address": "comdex1xlk88cqn9lljpeu47yu7vw5z8cf895zc4x07q6",
			"reward": {
				"denom": "ucmdx",
				"amount": "2288"
			}
		},
		{
			"address": "comdex1xlcua9lqtz2t86srptf5y7sdsn870524plgcs3",
			"reward": {
				"denom": "ucmdx",
				"amount": "5968"
			}
		},
		{
			"address": "comdex1xlm0gsj6y72ewkjyz00zwjzemjhyrz30rq9v44",
			"reward": {
				"denom": "ucmdx",
				"amount": "1317"
			}
		},
		{
			"address": "comdex1xlm5xh34hmjh8rrvw7krexupxsm86g2fddv7an",
			"reward": {
				"denom": "ucmdx",
				"amount": "7104"
			}
		},
		{
			"address": "comdex1xlmeapqyzrnnwcg65a8gmc6mh4lj48ack7qwkf",
			"reward": {
				"denom": "ucmdx",
				"amount": "4320"
			}
		},
		{
			"address": "comdex1xlaswhyluhp6tsq0h9wvglwnfamlnugwq6klqp",
			"reward": {
				"denom": "ucmdx",
				"amount": "179"
			}
		},
		{
			"address": "comdex1xl76c72v0qyz6qhcehquypw6pytk32jew2aszz",
			"reward": {
				"denom": "ucmdx",
				"amount": "55639"
			}
		},
		{
			"address": "comdex18qp6u0mvs7ac89gkn6r5j5jqu90zp74y53xsmf",
			"reward": {
				"denom": "ucmdx",
				"amount": "6758"
			}
		},
		{
			"address": "comdex18qz8ecysetqqhyv4qs9tlw5ylpgnx9ymwh06dw",
			"reward": {
				"denom": "ucmdx",
				"amount": "12451"
			}
		},
		{
			"address": "comdex18qyha72aytq4wdyus9ep0h34nqy338vm0kjsvv",
			"reward": {
				"denom": "ucmdx",
				"amount": "37064"
			}
		},
		{
			"address": "comdex18qyef5frjalrterst2xct67g7ss6n5u4unm8p7",
			"reward": {
				"denom": "ucmdx",
				"amount": "2361"
			}
		},
		{
			"address": "comdex18qshayud7c233sc670nugc5w4wr7szex29tjl7",
			"reward": {
				"denom": "ucmdx",
				"amount": "5877"
			}
		},
		{
			"address": "comdex18q32edxrasywvp4snlg9xk2wlmcc4n7hyf5jcu",
			"reward": {
				"denom": "ucmdx",
				"amount": "2818"
			}
		},
		{
			"address": "comdex18qj9t9e7maqsl7qt0a642232kyysgrxda64v5a",
			"reward": {
				"denom": "ucmdx",
				"amount": "6722"
			}
		},
		{
			"address": "comdex18qnmwu5ykwulswg3ldc0ddu5rksmpdwn82leg9",
			"reward": {
				"denom": "ucmdx",
				"amount": "61512"
			}
		},
		{
			"address": "comdex18qccy3r0p3z88c67r89gvu87anr9tzytqenqaq",
			"reward": {
				"denom": "ucmdx",
				"amount": "11810"
			}
		},
		{
			"address": "comdex18qm8ch3c85wt2r5hqe29ddey608nmehldnc0vy",
			"reward": {
				"denom": "ucmdx",
				"amount": "144"
			}
		},
		{
			"address": "comdex18qmtw0wmhrm0uan6xwx6d7dvn5373wavjaetxe",
			"reward": {
				"denom": "ucmdx",
				"amount": "3983"
			}
		},
		{
			"address": "comdex18qm7wjnstvyvvqqw6sarml47qxgxy8dfyzncln",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex18quet884vgkat3yu00yz05ca9wknh2d0vguvde",
			"reward": {
				"denom": "ucmdx",
				"amount": "3693"
			}
		},
		{
			"address": "comdex18pz8ftfzptnnd4xthvg4pmy5xn9vkuweaqxrly",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex18py44dux0qwm0mexydkq6ye66ycgrqt39juhrr",
			"reward": {
				"denom": "ucmdx",
				"amount": "50"
			}
		},
		{
			"address": "comdex18pgdl4qhu8s6px4f2fdz0m0lfjl0remdnsyvuf",
			"reward": {
				"denom": "ucmdx",
				"amount": "5861"
			}
		},
		{
			"address": "comdex18pdzj89fck3vaz3vwfx702kcnkqg7tz457rpxm",
			"reward": {
				"denom": "ucmdx",
				"amount": "151"
			}
		},
		{
			"address": "comdex18p339mq3yq2pnhwtpq08p74ypgmunxyr6dd2py",
			"reward": {
				"denom": "ucmdx",
				"amount": "5100"
			}
		},
		{
			"address": "comdex18p3lfnxdfuzjd9v6kmgvrjk0mun7rz8a5487yn",
			"reward": {
				"denom": "ucmdx",
				"amount": "1764"
			}
		},
		{
			"address": "comdex18pjf2rcwgygzt26d0g76z3fk3znh693pfrvf70",
			"reward": {
				"denom": "ucmdx",
				"amount": "3400"
			}
		},
		{
			"address": "comdex18p5wlzvvhn24k0w3h8zsasudehch43ezx8t43d",
			"reward": {
				"denom": "ucmdx",
				"amount": "897"
			}
		},
		{
			"address": "comdex18pud5xrke2wj88jydft7g6c9hylu96kg7sp9ap",
			"reward": {
				"denom": "ucmdx",
				"amount": "31232"
			}
		},
		{
			"address": "comdex18zph6wtfave8xc6q84r5fxzfqzsh6h2uejs3sy",
			"reward": {
				"denom": "ucmdx",
				"amount": "783"
			}
		},
		{
			"address": "comdex18zzs354zcvmp4kgg9ut0lza3wrtgtajqywk58x",
			"reward": {
				"denom": "ucmdx",
				"amount": "8868"
			}
		},
		{
			"address": "comdex18zz7vfr8vx3dghqxps4r57gu0jc9fl9vrfx9f9",
			"reward": {
				"denom": "ucmdx",
				"amount": "625"
			}
		},
		{
			"address": "comdex18zrz4v4jwk0nwjyl4c5q24mqps4xl0jyg08ttn",
			"reward": {
				"denom": "ucmdx",
				"amount": "1115"
			}
		},
		{
			"address": "comdex18z97tkm0uxggd4p9gsn3vezywqwajj2s5qraf7",
			"reward": {
				"denom": "ucmdx",
				"amount": "30978"
			}
		},
		{
			"address": "comdex18zxgp53e0x8k5z5wsanyzfutmxq62ukve0ql22",
			"reward": {
				"denom": "ucmdx",
				"amount": "4270"
			}
		},
		{
			"address": "comdex18z8ms73kalfgg6le8celn7wleecfx2e9utcpgj",
			"reward": {
				"denom": "ucmdx",
				"amount": "165"
			}
		},
		{
			"address": "comdex18z040fkvhxxvhgt5es7pzweksne0pn75qcwgc2",
			"reward": {
				"denom": "ucmdx",
				"amount": "28911"
			}
		},
		{
			"address": "comdex18zjqg4j695dcram66ge98zu2pvr94htnzgewcc",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex18zneqmyee9kqvtm8j24g5cauawym9k9xlwh8f9",
			"reward": {
				"denom": "ucmdx",
				"amount": "4096"
			}
		},
		{
			"address": "comdex18z6lyevyg36gmtwy5647pg724rjhf5amngnuy5",
			"reward": {
				"denom": "ucmdx",
				"amount": "1990"
			}
		},
		{
			"address": "comdex18rqhjv3yjc26vvk0z4vcarac77cz5x7wwhkwzz",
			"reward": {
				"denom": "ucmdx",
				"amount": "151"
			}
		},
		{
			"address": "comdex18r9wldyld430lrz9w8hcqkstnh8arjxuazdlgs",
			"reward": {
				"denom": "ucmdx",
				"amount": "19"
			}
		},
		{
			"address": "comdex18rxqjfzu5jtkvv8alpxt6dmshhpz89a5g4xd2t",
			"reward": {
				"denom": "ucmdx",
				"amount": "27362"
			}
		},
		{
			"address": "comdex18rxcw8ehuhn4dspyr9ky32t90ulpsphham5qnu",
			"reward": {
				"denom": "ucmdx",
				"amount": "1"
			}
		},
		{
			"address": "comdex18r8q3yn39ez2m9kre5g4307yeqala9pxnmx4ag",
			"reward": {
				"denom": "ucmdx",
				"amount": "181"
			}
		},
		{
			"address": "comdex18rgzexqd6c5j3jc5c907829e6an3wmsuk7j4an",
			"reward": {
				"denom": "ucmdx",
				"amount": "1985"
			}
		},
		{
			"address": "comdex18rgcw23nzx34u4wt54lv4a95ua0kxga0vjeqh3",
			"reward": {
				"denom": "ucmdx",
				"amount": "615"
			}
		},
		{
			"address": "comdex18rtv3t4525vsu3epyeqz6v57qaeafksnfd5c2m",
			"reward": {
				"denom": "ucmdx",
				"amount": "1768"
			}
		},
		{
			"address": "comdex18rdzwdpr0nsnds0ydqtcr4nuxl7kpk0tv9z085",
			"reward": {
				"denom": "ucmdx",
				"amount": "15178"
			}
		},
		{
			"address": "comdex18rwj6hygds8axwej2l959t4ejetf4hhzhhacfe",
			"reward": {
				"denom": "ucmdx",
				"amount": "697"
			}
		},
		{
			"address": "comdex18r3ua9wscepk848chfkn4ydthk5dq4ar9h2jen",
			"reward": {
				"denom": "ucmdx",
				"amount": "6298"
			}
		},
		{
			"address": "comdex18r44vuj0rca8jedfw836rzu7fge5ez53yw3w7t",
			"reward": {
				"denom": "ucmdx",
				"amount": "284"
			}
		},
		{
			"address": "comdex18r44epe6ju5pa5vd4rdf8yssqtrpjhzd3lteyz",
			"reward": {
				"denom": "ucmdx",
				"amount": "14440"
			}
		},
		{
			"address": "comdex18rk93f0m6actx32g5yca5ufgt3x8z9vqd0swd5",
			"reward": {
				"denom": "ucmdx",
				"amount": "343"
			}
		},
		{
			"address": "comdex18rh48fmz5e6vax6y0vpt4gwhu64vl9cne7fqyj",
			"reward": {
				"denom": "ucmdx",
				"amount": "1761"
			}
		},
		{
			"address": "comdex18r6vrf896fuu5976wacdj0k0pz5slh55vtzcsv",
			"reward": {
				"denom": "ucmdx",
				"amount": "526"
			}
		},
		{
			"address": "comdex18rmxxe4lz82kjzteuutl6q82enq33nx2kzqzyy",
			"reward": {
				"denom": "ucmdx",
				"amount": "1428"
			}
		},
		{
			"address": "comdex18ru2sffd4u2al42s3ns6kefhdkcput7v9f66zv",
			"reward": {
				"denom": "ucmdx",
				"amount": "168"
			}
		},
		{
			"address": "comdex18raa5zs3vhdsl45grs3ensgxnu4xk9lxrgq3xu",
			"reward": {
				"denom": "ucmdx",
				"amount": "23036"
			}
		},
		{
			"address": "comdex18rlnm5tc2zj6whmpu788lzk8fnusc3ppvcnlkc",
			"reward": {
				"denom": "ucmdx",
				"amount": "151"
			}
		},
		{
			"address": "comdex18yr6q5a5uzfwj4p24lajgyrlnv2duj263hqm6z",
			"reward": {
				"denom": "ucmdx",
				"amount": "9507"
			}
		},
		{
			"address": "comdex18yf77t533m77vagn6x2hdvqqc9ppg3z0ltgpap",
			"reward": {
				"denom": "ucmdx",
				"amount": "283"
			}
		},
		{
			"address": "comdex18y2vj9fcm2t43km2v5xys0judjs9ulvtfnfvuk",
			"reward": {
				"denom": "ucmdx",
				"amount": "61905"
			}
		},
		{
			"address": "comdex18y2kvfjhfc6e4tlxrwzjyxngdczs4dmvkxpnm8",
			"reward": {
				"denom": "ucmdx",
				"amount": "192"
			}
		},
		{
			"address": "comdex18y26gt56994yanj9092v4l4exf57zfvjx6aycj",
			"reward": {
				"denom": "ucmdx",
				"amount": "74903"
			}
		},
		{
			"address": "comdex18ytvgdpnacx7afy338w9rrtemmaquc0zm6lunf",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex18ytm3qa465p43xm8ss0lsnw0t45tmcf5sg33f0",
			"reward": {
				"denom": "ucmdx",
				"amount": "1755"
			}
		},
		{
			"address": "comdex18ydakgmca043fgawmv6cm676wa99957fjgjt28",
			"reward": {
				"denom": "ucmdx",
				"amount": "33285"
			}
		},
		{
			"address": "comdex18ywr8zc95akmaa7a5ylv5vr9wldzv6z8k54ztw",
			"reward": {
				"denom": "ucmdx",
				"amount": "1764"
			}
		},
		{
			"address": "comdex18ywwgc84luq0u8evna62hasdeyagzenaym245j",
			"reward": {
				"denom": "ucmdx",
				"amount": "353"
			}
		},
		{
			"address": "comdex18y3yehhv2kkn55za38vv46wu5wqvzvc53nhht0",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex18y30qlzs450wpf5sxjas2k2j6d05tcrjvf56sx",
			"reward": {
				"denom": "ucmdx",
				"amount": "0"
			}
		},
		{
			"address": "comdex18y3u88g6npwun8m3yuzc5cmudqgz5v5dkuje4e",
			"reward": {
				"denom": "ucmdx",
				"amount": "3512"
			}
		},
		{
			"address": "comdex18yj8tx8vl2l5l79xef5y7qqyu2ac75y0jmtzyv",
			"reward": {
				"denom": "ucmdx",
				"amount": "185"
			}
		},
		{
			"address": "comdex18yn5p2d3g8njm9lzh0ucy8ms2ahlvkj8mumz90",
			"reward": {
				"denom": "ucmdx",
				"amount": "3691"
			}
		},
		{
			"address": "comdex18ye3wgd28lu2qgq6r6jpsal5luy85wpmf9f2nd",
			"reward": {
				"denom": "ucmdx",
				"amount": "201"
			}
		},
		{
			"address": "comdex18y6vry4z982gz77ypxdl0cpj800g5js3zq8phz",
			"reward": {
				"denom": "ucmdx",
				"amount": "18543"
			}
		},
		{
			"address": "comdex18ym4dy9s6cczh3vgpzt3s2qstjynps2r9mqzuy",
			"reward": {
				"denom": "ucmdx",
				"amount": "17105"
			}
		},
		{
			"address": "comdex18yacz4kmy6chypkrhrfgl0an7e0xqv22dls80q",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex18ylsjechmwxc9cwh90xquwn8lw2exakekkms48",
			"reward": {
				"denom": "ucmdx",
				"amount": "1893"
			}
		},
		{
			"address": "comdex189qnw3zvwd0eht7m00pxtzhdr6pkh5uxwghlay",
			"reward": {
				"denom": "ucmdx",
				"amount": "1899"
			}
		},
		{
			"address": "comdex189rxgratl9r4kfeg9992sau5xvn6n66ewj3aq4",
			"reward": {
				"denom": "ucmdx",
				"amount": "185"
			}
		},
		{
			"address": "comdex189ym2ws9u8y6fk77dxe635kw9uyz8fr2nmmvsk",
			"reward": {
				"denom": "ucmdx",
				"amount": "19119"
			}
		},
		{
			"address": "comdex1899lnpdu6ntjpqrz9yv4j4xa07pc5h99y9dz76",
			"reward": {
				"denom": "ucmdx",
				"amount": "1021"
			}
		},
		{
			"address": "comdex189g9war7xmrl68g4rxjgafjl0e99jnh9zek07e",
			"reward": {
				"denom": "ucmdx",
				"amount": "2338"
			}
		},
		{
			"address": "comdex189g5e2hy7anxhqhaz9uy24z5ylvvqe0sx0ld0z",
			"reward": {
				"denom": "ucmdx",
				"amount": "15646"
			}
		},
		{
			"address": "comdex1892afnptlwyvzvngsx8x3m9w4g6f2q57x92hp7",
			"reward": {
				"denom": "ucmdx",
				"amount": "3699"
			}
		},
		{
			"address": "comdex189tenua4f2vtqlpsltaq8zese9q43ea40kz5yx",
			"reward": {
				"denom": "ucmdx",
				"amount": "11198"
			}
		},
		{
			"address": "comdex189vmrcqwaus996sma6mltx0qjhkxp5rn6dh4yg",
			"reward": {
				"denom": "ucmdx",
				"amount": "71129"
			}
		},
		{
			"address": "comdex189d6js7ghtyxkxx68rl00a4u77kgr2mkqvqauk",
			"reward": {
				"denom": "ucmdx",
				"amount": "71"
			}
		},
		{
			"address": "comdex1890d208737fjzkzc3z8v7pr9tzdc4544keuh2m",
			"reward": {
				"denom": "ucmdx",
				"amount": "34"
			}
		},
		{
			"address": "comdex1893vm2ecvpat3rg67qsyljmc6qe4h9xuf7smmv",
			"reward": {
				"denom": "ucmdx",
				"amount": "3063"
			}
		},
		{
			"address": "comdex189n8jcs42t56rzs7cur9wtz9pfz3sxachmt6ej",
			"reward": {
				"denom": "ucmdx",
				"amount": "74674"
			}
		},
		{
			"address": "comdex1894kfs7qzxca7x2et7xxfr0jhlr3flqg3ljq8h",
			"reward": {
				"denom": "ucmdx",
				"amount": "173"
			}
		},
		{
			"address": "comdex189hc0mdnukgp378wpj0srfdu7t6umjwq28txfy",
			"reward": {
				"denom": "ucmdx",
				"amount": "2015"
			}
		},
		{
			"address": "comdex18xyvurc7tz4fghuzkwafn3984a0r2sy9gywsze",
			"reward": {
				"denom": "ucmdx",
				"amount": "901"
			}
		},
		{
			"address": "comdex18xya0432fq3gythgugy76p0q3ugcyrgw0nayfm",
			"reward": {
				"denom": "ucmdx",
				"amount": "167"
			}
		},
		{
			"address": "comdex18xxj9y8l42ad8vxpu0pq3gpwq8nj98zpm28eex",
			"reward": {
				"denom": "ucmdx",
				"amount": "14342"
			}
		},
		{
			"address": "comdex18x8wrlsulemma9knd2mm74azyp3dkvun7l3fdt",
			"reward": {
				"denom": "ucmdx",
				"amount": "16"
			}
		},
		{
			"address": "comdex18xt39w2lw7q0u39z26hj0z6ft3k3cg3vna8lw2",
			"reward": {
				"denom": "ucmdx",
				"amount": "6175"
			}
		},
		{
			"address": "comdex18xd5n3x2e7rwaalggsuyejxdaw9g6hm8tsk8fj",
			"reward": {
				"denom": "ucmdx",
				"amount": "693"
			}
		},
		{
			"address": "comdex18xwazjxuu9h6v2h8ke30qz2z30d540c6ujecnf",
			"reward": {
				"denom": "ucmdx",
				"amount": "26"
			}
		},
		{
			"address": "comdex18xsrat8vknykw7dq9nzgcvjyq43ykecl7lsfde",
			"reward": {
				"denom": "ucmdx",
				"amount": "4488"
			}
		},
		{
			"address": "comdex18xsyk2j0l0nfj6xpmu9qs08xg85s4uj3vq8729",
			"reward": {
				"denom": "ucmdx",
				"amount": "5841"
			}
		},
		{
			"address": "comdex18xh8yml8kwsdu2k0hasdhvfz9rwc527fddpk6k",
			"reward": {
				"denom": "ucmdx",
				"amount": "179"
			}
		},
		{
			"address": "comdex18xhs7tmx833m36y7gvq9y4pqtymuu9xa0xytgz",
			"reward": {
				"denom": "ucmdx",
				"amount": "10833"
			}
		},
		{
			"address": "comdex18xevzusfne0kz50kha9zyf7smqx0r7vh5tmhpw",
			"reward": {
				"denom": "ucmdx",
				"amount": "13498"
			}
		},
		{
			"address": "comdex18xmsnlj209m6q35hfxt5dms9ugada9guc3p6m2",
			"reward": {
				"denom": "ucmdx",
				"amount": "6729"
			}
		},
		{
			"address": "comdex18xuqhpf52uy796ekfm6g0x7junfljwqp8qtsj3",
			"reward": {
				"denom": "ucmdx",
				"amount": "1687"
			}
		},
		{
			"address": "comdex18xu58cv8r3eqqa0rjxyvey5wpalmh4mwe9namv",
			"reward": {
				"denom": "ucmdx",
				"amount": "1413"
			}
		},
		{
			"address": "comdex18x7vnyynms56f40e7kxuaczkjdktxwy5qdcvzc",
			"reward": {
				"denom": "ucmdx",
				"amount": "6294"
			}
		},
		{
			"address": "comdex18xlu7xwusmmt2mgs0xenzk0nkwvagenwk07p2c",
			"reward": {
				"denom": "ucmdx",
				"amount": "3525"
			}
		},
		{
			"address": "comdex188pjlw6q8xtfhuzfgzstux8vg379na68xce3t6",
			"reward": {
				"denom": "ucmdx",
				"amount": "353"
			}
		},
		{
			"address": "comdex18897y62t5agv6rnu3ulcuy2q3ksys0nuqrmpg3",
			"reward": {
				"denom": "ucmdx",
				"amount": "434"
			}
		},
		{
			"address": "comdex18888m4mej8tfwedl4t4evcen5ap5f3d3rmdd9u",
			"reward": {
				"denom": "ucmdx",
				"amount": "142756"
			}
		},
		{
			"address": "comdex188g3v5cp5ea0tc0ales08vcmlc2kktygsvsjn5",
			"reward": {
				"denom": "ucmdx",
				"amount": "167"
			}
		},
		{
			"address": "comdex188fae2mtfak0xdrn0a66hyst2zlr94wds2dte5",
			"reward": {
				"denom": "ucmdx",
				"amount": "2046"
			}
		},
		{
			"address": "comdex1882dvsg9eurck0r8z40z78289f69h7hagt5avp",
			"reward": {
				"denom": "ucmdx",
				"amount": "166"
			}
		},
		{
			"address": "comdex188tqqsu8yv2vk9tmdre9yt66qd783h62qfc4ea",
			"reward": {
				"denom": "ucmdx",
				"amount": "6305"
			}
		},
		{
			"address": "comdex1880qjzhjp5v2dx5qgqn6l4nexcvzwss0lnmcfy",
			"reward": {
				"denom": "ucmdx",
				"amount": "12329"
			}
		},
		{
			"address": "comdex1880z9278v00hnksusunyxvgde9wvgkxkjz0p3d",
			"reward": {
				"denom": "ucmdx",
				"amount": "2812"
			}
		},
		{
			"address": "comdex1880z6724rrl8dnlmt82v59uw7ek0dvayx9sqtr",
			"reward": {
				"denom": "ucmdx",
				"amount": "176"
			}
		},
		{
			"address": "comdex188nh2zudndq8fca8a67mf0gghc55s88rl0lgnr",
			"reward": {
				"denom": "ucmdx",
				"amount": "4428"
			}
		},
		{
			"address": "comdex1884v32e4wdl48d7gskt97pq89xmcyedw64348s",
			"reward": {
				"denom": "ucmdx",
				"amount": "71"
			}
		},
		{
			"address": "comdex188cetgt7xulkt8glp264cety9jh600q7m55g4s",
			"reward": {
				"denom": "ucmdx",
				"amount": "26282"
			}
		},
		{
			"address": "comdex1886zc6l3x8uhrjs8a9uplnvqm5f8usa35xgzgy",
			"reward": {
				"denom": "ucmdx",
				"amount": "1524"
			}
		},
		{
			"address": "comdex188upmgk5dk06mpy6fuzk28njkk7uf9dh4nsruw",
			"reward": {
				"denom": "ucmdx",
				"amount": "1514"
			}
		},
		{
			"address": "comdex1887q2pnk07mu8nh7fcryrnkv5rgmlq3vrd27l4",
			"reward": {
				"denom": "ucmdx",
				"amount": "1603"
			}
		},
		{
			"address": "comdex18gqesp8xys8std82yjg8wr6067eqh0cvgqufqz",
			"reward": {
				"denom": "ucmdx",
				"amount": "3582"
			}
		},
		{
			"address": "comdex18grgrw89a3e88ttan78a04qsxujwjld8lrd5ep",
			"reward": {
				"denom": "ucmdx",
				"amount": "149436"
			}
		},
		{
			"address": "comdex18gxg34cp8r4gmsm5vgkaafg4fw4gzfurcfvmhf",
			"reward": {
				"denom": "ucmdx",
				"amount": "358"
			}
		},
		{
			"address": "comdex18g86k2tjareglenfflnpnt5w7gansz3clc0kv2",
			"reward": {
				"denom": "ucmdx",
				"amount": "560"
			}
		},
		{
			"address": "comdex18g2nswmcsyu46qva3y648q6d9m3wwzg7q37n49",
			"reward": {
				"denom": "ucmdx",
				"amount": "197"
			}
		},
		{
			"address": "comdex18gsq230k9lppqpgkaejarywxsxc06jq6v5vl50",
			"reward": {
				"denom": "ucmdx",
				"amount": "2894"
			}
		},
		{
			"address": "comdex18g45wgm2zyws40jd3w5xyzq66v450gt87lg8k7",
			"reward": {
				"denom": "ucmdx",
				"amount": "1769"
			}
		},
		{
			"address": "comdex18gcjmx5s7wacn07y906tjhpr6u54tqhexlc6jd",
			"reward": {
				"denom": "ucmdx",
				"amount": "3205"
			}
		},
		{
			"address": "comdex18gm5h80wfrxps5pk3l7vjqd8segu5ctphnsq8c",
			"reward": {
				"denom": "ucmdx",
				"amount": "52821"
			}
		},
		{
			"address": "comdex18gaj9c7fy8j07x6yjf0gpjhd37dz8dusg98gxw",
			"reward": {
				"denom": "ucmdx",
				"amount": "711"
			}
		},
		{
			"address": "comdex18glwdnhg8lqal60r2l3846kjpgqs06yn66jy4g",
			"reward": {
				"denom": "ucmdx",
				"amount": "289"
			}
		},
		{
			"address": "comdex18glslf6qzdphu7f79yd3j2a8ag9qyqcklwa288",
			"reward": {
				"denom": "ucmdx",
				"amount": "5206"
			}
		},
		{
			"address": "comdex18fx5me3yxp8pu6kd2e42gkjqthyv867pqssqc4",
			"reward": {
				"denom": "ucmdx",
				"amount": "824"
			}
		},
		{
			"address": "comdex18f8ty5jlct72vsnt4rwxzzk2auqt4rtwrs3jtx",
			"reward": {
				"denom": "ucmdx",
				"amount": "497"
			}
		},
		{
			"address": "comdex18fgflhpucd4ac7hun49z49mv4p6a68z3spwph8",
			"reward": {
				"denom": "ucmdx",
				"amount": "179"
			}
		},
		{
			"address": "comdex18fghxtdfeypc9j0r52lnw2kh78uhdt5dqj6p2w",
			"reward": {
				"denom": "ucmdx",
				"amount": "1419"
			}
		},
		{
			"address": "comdex18fg6lnzr2pa8fdgx59we6mjx8j9xazam5cl4z5",
			"reward": {
				"denom": "ucmdx",
				"amount": "1741"
			}
		},
		{
			"address": "comdex18f22qngc7s0xmzd2mal8gctafznqx6xshygqz2",
			"reward": {
				"denom": "ucmdx",
				"amount": "4078"
			}
		},
		{
			"address": "comdex18ft8chhqk5ggvd6fd4tujcevjvke7g56fs3gmh",
			"reward": {
				"denom": "ucmdx",
				"amount": "2043"
			}
		},
		{
			"address": "comdex18ftjvvzdva674ljp93ulvf5v72dhryn5gkvw2l",
			"reward": {
				"denom": "ucmdx",
				"amount": "6403"
			}
		},
		{
			"address": "comdex18ftcgpfmflcqpmy79vlqnenxhmz6na5uw7xpgy",
			"reward": {
				"denom": "ucmdx",
				"amount": "43604"
			}
		},
		{
			"address": "comdex18fvuta8pqgm6g9zjqax5lv8qxh42c436rk37kg",
			"reward": {
				"denom": "ucmdx",
				"amount": "1250"
			}
		},
		{
			"address": "comdex18fva927qr4ddzcm54sn2h87qvdky3cy6w7kx2f",
			"reward": {
				"denom": "ucmdx",
				"amount": "14201"
			}
		},
		{
			"address": "comdex18fd9nn5dce6z3qlt86k62wqzek8tzlhtpqm3an",
			"reward": {
				"denom": "ucmdx",
				"amount": "39860"
			}
		},
		{
			"address": "comdex18fsdl9jzhs78p3rmjyf5vyejcjmeffsvn9frkt",
			"reward": {
				"denom": "ucmdx",
				"amount": "142249"
			}
		},
		{
			"address": "comdex18fs3vnjzxst4jfjfzgxreke3mp20fmds3xz800",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex18fnfpvcj3pfr8cj2hfz4y0qlj4679f5zft7x5h",
			"reward": {
				"denom": "ucmdx",
				"amount": "936"
			}
		},
		{
			"address": "comdex18f5nduc77d6qv7kqjpcu48dkglyaeyxshq9rwn",
			"reward": {
				"denom": "ucmdx",
				"amount": "1615"
			}
		},
		{
			"address": "comdex18fkxdfs4va8amxdhvlthv0hevjc7dg0kvtwqyn",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex18fccqhr3agtd5q8r0r24fe8x6x4lks9c980pku",
			"reward": {
				"denom": "ucmdx",
				"amount": "26212"
			}
		},
		{
			"address": "comdex18fc79ke4r88f5fvxwsxeqemzs4sez68t9ffd4z",
			"reward": {
				"denom": "ucmdx",
				"amount": "13572"
			}
		},
		{
			"address": "comdex18f6tavz9v72mjee4xq6zuj32rgryepxnz87fk2",
			"reward": {
				"denom": "ucmdx",
				"amount": "795"
			}
		},
		{
			"address": "comdex18fmul4thsu4tpftsl6gt7ja6fg8u323gj4f9lf",
			"reward": {
				"denom": "ucmdx",
				"amount": "169"
			}
		},
		{
			"address": "comdex18f7lp0tprynwhkf9nj5h99fdacpq3df0cn65jv",
			"reward": {
				"denom": "ucmdx",
				"amount": "1720"
			}
		},
		{
			"address": "comdex18flrc9fjt0yp73e56mmq7hzprfakawdx843umc",
			"reward": {
				"denom": "ucmdx",
				"amount": "678"
			}
		},
		{
			"address": "comdex18fl9zr2luz5sff0kkd5h27sn5mk9g3gjkwuz25",
			"reward": {
				"denom": "ucmdx",
				"amount": "1440"
			}
		},
		{
			"address": "comdex182xlpqjrrcw9nxd8g2chzzwygx3zuk8m8wzep5",
			"reward": {
				"denom": "ucmdx",
				"amount": "3467"
			}
		},
		{
			"address": "comdex18280cjsrmzrvcmpcjds78zfm3vjgdwzws4ra27",
			"reward": {
				"denom": "ucmdx",
				"amount": "29"
			}
		},
		{
			"address": "comdex18209hfrqcwvu9gufukf73v0y0kmwlzzet0w36q",
			"reward": {
				"denom": "ucmdx",
				"amount": "10591"
			}
		},
		{
			"address": "comdex1820vq4v2vy9f70k8f7p2tff49lz504z9tdcjx4",
			"reward": {
				"denom": "ucmdx",
				"amount": "200"
			}
		},
		{
			"address": "comdex1820v59dwdzjfe7fqkc4fq2qrsenz9qn0gtfly0",
			"reward": {
				"denom": "ucmdx",
				"amount": "412"
			}
		},
		{
			"address": "comdex1823v4rrxyfwlyu3laxcqn8f46putq4zvddp5gh",
			"reward": {
				"denom": "ucmdx",
				"amount": "4011"
			}
		},
		{
			"address": "comdex1824ghchgmde22agcfkdgaqw5td8t2qr7h45xk6",
			"reward": {
				"denom": "ucmdx",
				"amount": "6647"
			}
		},
		{
			"address": "comdex182mzj3c7hg96fsaratn0wt2khgtgxgvl4d7p3f",
			"reward": {
				"denom": "ucmdx",
				"amount": "3842"
			}
		},
		{
			"address": "comdex182mvfd3m8fr4u74j3skvwnd3mrz56lallzz4my",
			"reward": {
				"denom": "ucmdx",
				"amount": "17055"
			}
		},
		{
			"address": "comdex182mng5t8tq4j7he3jc7ahyfdap2ldjw4j5k8hy",
			"reward": {
				"denom": "ucmdx",
				"amount": "1788"
			}
		},
		{
			"address": "comdex1827xsel9rtrenhevr5n80qsdah3vqjev5600dh",
			"reward": {
				"denom": "ucmdx",
				"amount": "181"
			}
		},
		{
			"address": "comdex1827h7hmwjf9ccacy79fhqz46lfxewujvdn80up",
			"reward": {
				"denom": "ucmdx",
				"amount": "806"
			}
		},
		{
			"address": "comdex18tq9gk04sfrdmms4e2kd5wa4lh5r62k28crvy2",
			"reward": {
				"denom": "ucmdx",
				"amount": "31255"
			}
		},
		{
			"address": "comdex18tzp963n2vyw79rvgm0p8k0msq6tcqxrjclztr",
			"reward": {
				"denom": "ucmdx",
				"amount": "178"
			}
		},
		{
			"address": "comdex18tzss72jrt60yfd7etwpy3mggmntl2slvyjhu8",
			"reward": {
				"denom": "ucmdx",
				"amount": "6297"
			}
		},
		{
			"address": "comdex18tzene5kjaj2z6ry7tj67h9a80mjmrd60l0j3p",
			"reward": {
				"denom": "ucmdx",
				"amount": "56485"
			}
		},
		{
			"address": "comdex18trqyqnv6rzaf0auzyc38cdyj5xvcv3dr7dtvm",
			"reward": {
				"denom": "ucmdx",
				"amount": "4069"
			}
		},
		{
			"address": "comdex18t9veg2vd6ep0wmxu2z7v5367s8zk9uhzend5t",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex18t97ph2esca72wmrgnhmljk4zux60wax370m2x",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex18tgy7tmpng6qspxgwl4alvc0myq6ln7ugt8wqv",
			"reward": {
				"denom": "ucmdx",
				"amount": "197"
			}
		},
		{
			"address": "comdex18tg684zh36se7aszz906mt0pmh0e8fd24eqhtx",
			"reward": {
				"denom": "ucmdx",
				"amount": "84660"
			}
		},
		{
			"address": "comdex18tfm8787jd4nducxgewhg2t33c2tstsuqqptl0",
			"reward": {
				"denom": "ucmdx",
				"amount": "52"
			}
		},
		{
			"address": "comdex18tvlchkpvnmettc7wzz42r86hn3harhj3a24ud",
			"reward": {
				"denom": "ucmdx",
				"amount": "697"
			}
		},
		{
			"address": "comdex18tn367aundw2kuhkg88tyul2cddx60v53l8vuv",
			"reward": {
				"denom": "ucmdx",
				"amount": "411"
			}
		},
		{
			"address": "comdex18tncp9gdgd8p9ktvcddrp00j6mnpqfwnshh5ep",
			"reward": {
				"denom": "ucmdx",
				"amount": "2977"
			}
		},
		{
			"address": "comdex18t46lsd89dp2uvae35fxpzenj43hr4l2gwj9qn",
			"reward": {
				"denom": "ucmdx",
				"amount": "16393"
			}
		},
		{
			"address": "comdex18thwcddjk76wkyytu0vym9z9sl3mkzhjchyw28",
			"reward": {
				"denom": "ucmdx",
				"amount": "9014"
			}
		},
		{
			"address": "comdex18th67pfq56ecq4m3985a5upergwqu5kssalkc2",
			"reward": {
				"denom": "ucmdx",
				"amount": "13197"
			}
		},
		{
			"address": "comdex18tcjak6ltt7fmg2lfkls93w6eamj3zjnjn07dd",
			"reward": {
				"denom": "ucmdx",
				"amount": "2827"
			}
		},
		{
			"address": "comdex18tey9dudnazg7r7wws6r6wm683n692g7s7zs8r",
			"reward": {
				"denom": "ucmdx",
				"amount": "1747"
			}
		},
		{
			"address": "comdex18t6l5m0vs5995h9mr9ngkh8vvq7aq3ah9lj7m7",
			"reward": {
				"denom": "ucmdx",
				"amount": "526"
			}
		},
		{
			"address": "comdex18tmnn4wtavenpxgxjk3qdpnnv0llpf9rl2a080",
			"reward": {
				"denom": "ucmdx",
				"amount": "132"
			}
		},
		{
			"address": "comdex18tmky7gmz76jxxgywhmt2w89tm8qu4rzluyeyl",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex18tmkvruqsz3xjwzsh94yjr2jrnl8v7sc2vd4kt",
			"reward": {
				"denom": "ucmdx",
				"amount": "8654"
			}
		},
		{
			"address": "comdex18tmheu6ykvl3vmhwr63480wknhpxe003xmnq65",
			"reward": {
				"denom": "ucmdx",
				"amount": "1404"
			}
		},
		{
			"address": "comdex18t7qytg9fdgklzac9lm4ym0cjxhjcxsdvalxdj",
			"reward": {
				"denom": "ucmdx",
				"amount": "2845"
			}
		},
		{
			"address": "comdex18t7juu62sxnlwj672l2cyve8cvqknunecprvre",
			"reward": {
				"denom": "ucmdx",
				"amount": "9302"
			}
		},
		{
			"address": "comdex18t7utq0wvtaa80t3sl3309h2staa74m924y885",
			"reward": {
				"denom": "ucmdx",
				"amount": "16"
			}
		},
		{
			"address": "comdex18vypcsuef67cp0yxteerpfufeykwea4qmk3u68",
			"reward": {
				"denom": "ucmdx",
				"amount": "60469"
			}
		},
		{
			"address": "comdex18vygp7wghlwqqcctsldl9d7qd38nn5d7vv52v8",
			"reward": {
				"denom": "ucmdx",
				"amount": "7139"
			}
		},
		{
			"address": "comdex18v8mgsx4m3p220pp3ywa09rmhle22sfd79t05j",
			"reward": {
				"denom": "ucmdx",
				"amount": "2534"
			}
		},
		{
			"address": "comdex18vgjru99uty22nszm9n8xex8hux2ydnae8p976",
			"reward": {
				"denom": "ucmdx",
				"amount": "17904"
			}
		},
		{
			"address": "comdex18vf0ngq6ecd87asmk96r3avpwg3gtl4vewu07w",
			"reward": {
				"denom": "ucmdx",
				"amount": "6312"
			}
		},
		{
			"address": "comdex18vflvmz0x6pceuv7yp49ez7jtjsjgfrcfxdvyy",
			"reward": {
				"denom": "ucmdx",
				"amount": "604"
			}
		},
		{
			"address": "comdex18v2epwkmh42pek7d8mndxk9hcaxu7n4m5pl8yv",
			"reward": {
				"denom": "ucmdx",
				"amount": "2988"
			}
		},
		{
			"address": "comdex18vv6c7ma8twmpwtuyrlx7se20ahna9lnv9f9xk",
			"reward": {
				"denom": "ucmdx",
				"amount": "6980"
			}
		},
		{
			"address": "comdex18vdtjm6f8xvz8w4xzwycha2duq4ugy7ryk3msy",
			"reward": {
				"denom": "ucmdx",
				"amount": "34"
			}
		},
		{
			"address": "comdex18v39t7z7nvdkdd3fjxqjrf32dqr6atr0e9ke2n",
			"reward": {
				"denom": "ucmdx",
				"amount": "14757"
			}
		},
		{
			"address": "comdex18v5sjcp8vz2zd7xpul8x98vyvwg5yc0ln9ur42",
			"reward": {
				"denom": "ucmdx",
				"amount": "7150"
			}
		},
		{
			"address": "comdex18v4ql94unanrcpgkc08e3d0509qd0ypf7f53x5",
			"reward": {
				"denom": "ucmdx",
				"amount": "3354"
			}
		},
		{
			"address": "comdex18vheqfa3fjvkwf2ucyaykhex6tp2mkdq2w4eud",
			"reward": {
				"denom": "ucmdx",
				"amount": "150"
			}
		},
		{
			"address": "comdex18vc047wgzgttm668564alc2t2j3nwhel7ts3gh",
			"reward": {
				"denom": "ucmdx",
				"amount": "2754243"
			}
		},
		{
			"address": "comdex18vc7asglvk6tlzk4nvkr59zv4duqeawj6wqdx0",
			"reward": {
				"denom": "ucmdx",
				"amount": "6222"
			}
		},
		{
			"address": "comdex18vmyj3348y0q6cgh74p2jm5cw3rqd9x0h7f005",
			"reward": {
				"denom": "ucmdx",
				"amount": "13346"
			}
		},
		{
			"address": "comdex18drp7uw2lvx4ju40ddcjjek6t4tnc0fey0zt56",
			"reward": {
				"denom": "ucmdx",
				"amount": "139"
			}
		},
		{
			"address": "comdex18d9msewku04ed760mvw9yr7ad7n8vag6uc4d5n",
			"reward": {
				"denom": "ucmdx",
				"amount": "26700"
			}
		},
		{
			"address": "comdex18dxzdrdu6t4us6dmc9q66qxecy8v3pnrmlplw5",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex18dxwd9qyhwkzw6h8ep9lld7q3ecvmuctpmp644",
			"reward": {
				"denom": "ucmdx",
				"amount": "28152"
			}
		},
		{
			"address": "comdex18dgpn2pg3hsmhl9eu5wjcr5zep4t3lvv6a9x76",
			"reward": {
				"denom": "ucmdx",
				"amount": "246"
			}
		},
		{
			"address": "comdex18dg4tkhkgq3tzdt0z8penwss2qw2zk9ashwrrm",
			"reward": {
				"denom": "ucmdx",
				"amount": "71359"
			}
		},
		{
			"address": "comdex18d2x8tnul8f5nurntag5nqrvsryarcjrz2jh33",
			"reward": {
				"denom": "ucmdx",
				"amount": "991"
			}
		},
		{
			"address": "comdex18d2gzh82jjvyva2e8xf2avlafg4g8zxj9cmvje",
			"reward": {
				"denom": "ucmdx",
				"amount": "301"
			}
		},
		{
			"address": "comdex18d2ge2rnttunzas75ygkjl5s8tvetevsd0v6wh",
			"reward": {
				"denom": "ucmdx",
				"amount": "353"
			}
		},
		{
			"address": "comdex18dvpf5yeettvxfwmem4hnsmmtqtr57vdznhzn2",
			"reward": {
				"denom": "ucmdx",
				"amount": "30039"
			}
		},
		{
			"address": "comdex18dsuaulan4tkqm9kug2s5klx84997rm7rwvcdx",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex18dj9q7nfxyrgy5gymq0nf5chpxvu2f7902fyvy",
			"reward": {
				"denom": "ucmdx",
				"amount": "116212"
			}
		},
		{
			"address": "comdex18d4mfwawt4lkd5y4uk9gt3caxqf2svcyw0ejmc",
			"reward": {
				"denom": "ucmdx",
				"amount": "36629"
			}
		},
		{
			"address": "comdex18dk8ew6ee40a3m0065se5et9smx5lykj77pg0r",
			"reward": {
				"denom": "ucmdx",
				"amount": "12522"
			}
		},
		{
			"address": "comdex18dhtsvjsasz9qfjh7qfnv40zsge5eecugphzwj",
			"reward": {
				"denom": "ucmdx",
				"amount": "133"
			}
		},
		{
			"address": "comdex18dcvv9xd8t7kyxucjg8g3s3d502qumqmnzm7qm",
			"reward": {
				"denom": "ucmdx",
				"amount": "14034"
			}
		},
		{
			"address": "comdex18d66t5p3x2tf5mjml58uchp22n3axuaxccd9ng",
			"reward": {
				"denom": "ucmdx",
				"amount": "17719"
			}
		},
		{
			"address": "comdex18duveg07j5k7x2jnvdd2309ccnfxaqtkmssnt4",
			"reward": {
				"denom": "ucmdx",
				"amount": "3142"
			}
		},
		{
			"address": "comdex18da054k6aavla030fcx3k3tvjm6z2fa4thn97v",
			"reward": {
				"denom": "ucmdx",
				"amount": "1887"
			}
		},
		{
			"address": "comdex18wzgwr458dhpqca2q8pccaz3p5eagevvk5x0t2",
			"reward": {
				"denom": "ucmdx",
				"amount": "17891"
			}
		},
		{
			"address": "comdex18wz2jqtzle8cr8756np6txl9gnwv5q42e7hpat",
			"reward": {
				"denom": "ucmdx",
				"amount": "4268"
			}
		},
		{
			"address": "comdex18wycp23y69ky0rla4p75muzhejn0cd76z5xhmq",
			"reward": {
				"denom": "ucmdx",
				"amount": "35452"
			}
		},
		{
			"address": "comdex18w9zzcnmqtypc689gkewxg964vvxxa6l3e35ae",
			"reward": {
				"denom": "ucmdx",
				"amount": "1248"
			}
		},
		{
			"address": "comdex18wg80p26wzjwdggclux0mzeetdapfhjcseqxnx",
			"reward": {
				"denom": "ucmdx",
				"amount": "804833"
			}
		},
		{
			"address": "comdex18wdhlhw9qs6hm9dwuadw72qwr259yj3nkj96h2",
			"reward": {
				"denom": "ucmdx",
				"amount": "25216"
			}
		},
		{
			"address": "comdex18ww0rqmfcduucqtz2s93x5u0dfuvs2wthfhcn0",
			"reward": {
				"denom": "ucmdx",
				"amount": "17313"
			}
		},
		{
			"address": "comdex18w05mmczm9e3z4cuww5ke07g6gsl8v5exj5kcf",
			"reward": {
				"denom": "ucmdx",
				"amount": "1876"
			}
		},
		{
			"address": "comdex18ws3fpnclluwjp6phr0hpyaukgu9neskv8y7ve",
			"reward": {
				"denom": "ucmdx",
				"amount": "113"
			}
		},
		{
			"address": "comdex18w3r86plh703fvpktjawepskq4gy0pgclgq47r",
			"reward": {
				"denom": "ucmdx",
				"amount": "766"
			}
		},
		{
			"address": "comdex18wncrwwjamua2y8wq89pctmg0jz5g6k0fpazxy",
			"reward": {
				"denom": "ucmdx",
				"amount": "1059"
			}
		},
		{
			"address": "comdex18w6pkfsvktjnccysylcx48lyxnfzp33q3ejfwp",
			"reward": {
				"denom": "ucmdx",
				"amount": "2081"
			}
		},
		{
			"address": "comdex18w70c6a8w7f8q3vlv0uktgkxnsr5rxrxxc0qac",
			"reward": {
				"denom": "ucmdx",
				"amount": "18288"
			}
		},
		{
			"address": "comdex180qhgwdjtqaav5f70du33nu6seasu2tz2m43v3",
			"reward": {
				"denom": "ucmdx",
				"amount": "1925"
			}
		},
		{
			"address": "comdex180pdyeng2am7mmqwm9ruvftsddmd2vyj7uv5f2",
			"reward": {
				"denom": "ucmdx",
				"amount": "13820"
			}
		},
		{
			"address": "comdex180g0wrwqu4e27mmsg0rtalmuhd9983gttndu5r",
			"reward": {
				"denom": "ucmdx",
				"amount": "382966"
			}
		},
		{
			"address": "comdex180g67lek5haptetllkazu0xr0fcya2uczu6vve",
			"reward": {
				"denom": "ucmdx",
				"amount": "19"
			}
		},
		{
			"address": "comdex180fu70aw4834d4m57je54aar3lx2mc80uksqwp",
			"reward": {
				"denom": "ucmdx",
				"amount": "181"
			}
		},
		{
			"address": "comdex1802qk6y02af8jd5ucdqzrqrrt9m08xfqnnul9e",
			"reward": {
				"denom": "ucmdx",
				"amount": "122179"
			}
		},
		{
			"address": "comdex180j05q24wpedqesguwdkta42k0gxd7j6mptuj3",
			"reward": {
				"denom": "ucmdx",
				"amount": "1228"
			}
		},
		{
			"address": "comdex180j36lh2axv6763xg0gd334yr6au727t3puanr",
			"reward": {
				"denom": "ucmdx",
				"amount": "12621"
			}
		},
		{
			"address": "comdex180c39urf2eekhssxxq2pn44smx7zw2r2ly8j9a",
			"reward": {
				"denom": "ucmdx",
				"amount": "61138"
			}
		},
		{
			"address": "comdex180cj9v695r4qkk3e6k7gcj5zurx674dyyg4gmr",
			"reward": {
				"denom": "ucmdx",
				"amount": "5726"
			}
		},
		{
			"address": "comdex180e96f5kacd9utuvvfr2fma0yvd7ypyypha67h",
			"reward": {
				"denom": "ucmdx",
				"amount": "6156"
			}
		},
		{
			"address": "comdex180eunxeqfdcz88rz5vap6zhsez76v3hele8zyk",
			"reward": {
				"denom": "ucmdx",
				"amount": "14295"
			}
		},
		{
			"address": "comdex18068g3r6chexz9dl9tkqp90yjlzm7es6km3c2v",
			"reward": {
				"denom": "ucmdx",
				"amount": "384"
			}
		},
		{
			"address": "comdex180u2pakfaczrurlw89wgmqr9qv7mt6q438s8ew",
			"reward": {
				"denom": "ucmdx",
				"amount": "2018"
			}
		},
		{
			"address": "comdex180a68ng3mrux3sq82lukx0kk608dznw0p557px",
			"reward": {
				"denom": "ucmdx",
				"amount": "1349"
			}
		},
		{
			"address": "comdex1807s9nqyzeshp7ukcawt0j0kvl4fxpqx30lz9h",
			"reward": {
				"denom": "ucmdx",
				"amount": "142"
			}
		},
		{
			"address": "comdex1807ejwpzd7yeqevdrp32yx5tyjcq35xg7nwyqn",
			"reward": {
				"denom": "ucmdx",
				"amount": "2797"
			}
		},
		{
			"address": "comdex180l5y8h6ygc96tqe3yeejltxcm3rst4zw0wtqn",
			"reward": {
				"denom": "ucmdx",
				"amount": "11361"
			}
		},
		{
			"address": "comdex180lm374735m4hxs6lvzah7hxwl28pc5xgr8l0g",
			"reward": {
				"denom": "ucmdx",
				"amount": "618"
			}
		},
		{
			"address": "comdex18sfflzjtercqwv3w4dtg9w6dlgvx0ps93e8388",
			"reward": {
				"denom": "ucmdx",
				"amount": "194"
			}
		},
		{
			"address": "comdex18stdn9kl56mg6epnrms8vunedqjcsrwrdhx64n",
			"reward": {
				"denom": "ucmdx",
				"amount": "2010"
			}
		},
		{
			"address": "comdex18svlnyrlt0kfh7t2x8extv7wrw2gc7xsryxwcf",
			"reward": {
				"denom": "ucmdx",
				"amount": "6939"
			}
		},
		{
			"address": "comdex18sd0v7u65rxd84mhngn0rmt0jl4jw05qhugm0z",
			"reward": {
				"denom": "ucmdx",
				"amount": "1425"
			}
		},
		{
			"address": "comdex18sdews8g2caqfcgq9wx483kvr7t6vcsegrgs29",
			"reward": {
				"denom": "ucmdx",
				"amount": "1167"
			}
		},
		{
			"address": "comdex18s34fescdunkdveayskhuxlsyqz5tv3jmduyp6",
			"reward": {
				"denom": "ucmdx",
				"amount": "12750"
			}
		},
		{
			"address": "comdex18s3kl3ze5rzz4s3unf9xexdprru03ywvk6xk2m",
			"reward": {
				"denom": "ucmdx",
				"amount": "2645"
			}
		},
		{
			"address": "comdex18snggy0cxl55f828eaw2z4nndcwrkexk0z5vcm",
			"reward": {
				"denom": "ucmdx",
				"amount": "125392"
			}
		},
		{
			"address": "comdex18sn4spm78jsyensvn5mlk7txj3jjwdfjxh3pka",
			"reward": {
				"denom": "ucmdx",
				"amount": "21252"
			}
		},
		{
			"address": "comdex18sh66ah6df0ptwugmwj550a4k6zvfpsmghjzjy",
			"reward": {
				"denom": "ucmdx",
				"amount": "29"
			}
		},
		{
			"address": "comdex18s6wg6emqtck0q688xurd0t3fr0wzf7dx4nzsr",
			"reward": {
				"denom": "ucmdx",
				"amount": "11291"
			}
		},
		{
			"address": "comdex18smgsfadpg320pclyn5l6h2n8fx0w0t7yddynm",
			"reward": {
				"denom": "ucmdx",
				"amount": "142"
			}
		},
		{
			"address": "comdex18suweqrppd9hg3ykdxuzn9gw8y9df2x8e9h6rg",
			"reward": {
				"denom": "ucmdx",
				"amount": "14405"
			}
		},
		{
			"address": "comdex18sax55hceu3nn66x9gsgp6hzm54ewz8ujsctka",
			"reward": {
				"denom": "ucmdx",
				"amount": "14328"
			}
		},
		{
			"address": "comdex18sajp6w6yyl7x2t0prd39wp5gspq3vth7vmh6s",
			"reward": {
				"denom": "ucmdx",
				"amount": "7139"
			}
		},
		{
			"address": "comdex183qypvjfzn7xhnfk2ya3c0xfnp94j8asmpas8x",
			"reward": {
				"denom": "ucmdx",
				"amount": "285"
			}
		},
		{
			"address": "comdex183q3z5ldllmczxuq620k8pauydwe6z74tqx0cq",
			"reward": {
				"denom": "ucmdx",
				"amount": "4103"
			}
		},
		{
			"address": "comdex183xw8l4vs7jxx872hucdp56dfgansunzqgl2an",
			"reward": {
				"denom": "ucmdx",
				"amount": "115199"
			}
		},
		{
			"address": "comdex183xh7f4umya7prz7v90z7lzq6g48wgeuhgyy4a",
			"reward": {
				"denom": "ucmdx",
				"amount": "261"
			}
		},
		{
			"address": "comdex183x7vjta83c3dm0lff08r73x4uzlrawyqhs9kf",
			"reward": {
				"denom": "ucmdx",
				"amount": "24883"
			}
		},
		{
			"address": "comdex183vejcy69mun7q6fa38gnpwk3mmgvrllp7dhr9",
			"reward": {
				"denom": "ucmdx",
				"amount": "506"
			}
		},
		{
			"address": "comdex183v76f0fz62jg5tk2sg0ysd5dmer4mt8q4vthj",
			"reward": {
				"denom": "ucmdx",
				"amount": "2179"
			}
		},
		{
			"address": "comdex18359mjlfkp5s9ev8tw798m8xcc8lrq0kp3symz",
			"reward": {
				"denom": "ucmdx",
				"amount": "132"
			}
		},
		{
			"address": "comdex183hrcuzzt3az7p4je6l4vdnr2qp7q82lfyaxx3",
			"reward": {
				"denom": "ucmdx",
				"amount": "9598"
			}
		},
		{
			"address": "comdex183h79x8lmh6wg8avt5hfncvgv98vps64x9ptan",
			"reward": {
				"denom": "ucmdx",
				"amount": "710"
			}
		},
		{
			"address": "comdex183czmqgy5rujz4per7ekjygtyx8pjvkl3jl47j",
			"reward": {
				"denom": "ucmdx",
				"amount": "9565"
			}
		},
		{
			"address": "comdex183elyd28thv5l5uv5w9u6k0kdruuzczdy49aml",
			"reward": {
				"denom": "ucmdx",
				"amount": "12597"
			}
		},
		{
			"address": "comdex1837ug8qw34e6rqv4837nkmgxv33aq6nk05wfrg",
			"reward": {
				"denom": "ucmdx",
				"amount": "15048"
			}
		},
		{
			"address": "comdex183lrvkzpfd5f5negk28preenq6yz3lg8hmk24k",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex18jrqe4954awdpv3yam8q2nu3k56dkq592j80qv",
			"reward": {
				"denom": "ucmdx",
				"amount": "302"
			}
		},
		{
			"address": "comdex18jy3dvdwv3usn7m4urglldj4tdsf6ldn94d5jq",
			"reward": {
				"denom": "ucmdx",
				"amount": "5247"
			}
		},
		{
			"address": "comdex18j2eg363cz8s8q8mp8hzsx44dlvukp73cm2p9t",
			"reward": {
				"denom": "ucmdx",
				"amount": "7295"
			}
		},
		{
			"address": "comdex18jtr958kvgc0q4306y572lx0ngltrdzy8wn348",
			"reward": {
				"denom": "ucmdx",
				"amount": "184"
			}
		},
		{
			"address": "comdex18jtj2tpzcd6vd0kkdcfdsymgzlmclw9exnl8l0",
			"reward": {
				"denom": "ucmdx",
				"amount": "3550"
			}
		},
		{
			"address": "comdex18jtcdz6n73g3sw0u9ek4s6d2wkhtpjdkf73uzj",
			"reward": {
				"denom": "ucmdx",
				"amount": "148"
			}
		},
		{
			"address": "comdex18jv2n4d2rvysv69e4q7d9k4p7wj3hvtkkpmwv2",
			"reward": {
				"denom": "ucmdx",
				"amount": "15225"
			}
		},
		{
			"address": "comdex18jdvz0gg4c2x2wt35eefvdums2mtk7hg02cya7",
			"reward": {
				"denom": "ucmdx",
				"amount": "967"
			}
		},
		{
			"address": "comdex18jwmcye8dz4s7we4ge5y6llekur3kq2uap5ylh",
			"reward": {
				"denom": "ucmdx",
				"amount": "181"
			}
		},
		{
			"address": "comdex18j08s9we7z2yvgdshdce4vqhqyacte2rfylf45",
			"reward": {
				"denom": "ucmdx",
				"amount": "4714"
			}
		},
		{
			"address": "comdex18j0adn0l606srstepmsje7k7v6cavt5urx776p",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex18jnqnds9xg085c3wr730ck4fvr6c0gu8yf3edj",
			"reward": {
				"denom": "ucmdx",
				"amount": "264622"
			}
		},
		{
			"address": "comdex18jem3e68gyyjekqyqy4aewznx2pheqmqeh4kme",
			"reward": {
				"denom": "ucmdx",
				"amount": "3818"
			}
		},
		{
			"address": "comdex18ju02hfx2fksxy8l4u4g8tp636e4ypaers2hh6",
			"reward": {
				"denom": "ucmdx",
				"amount": "44"
			}
		},
		{
			"address": "comdex18ja0e0pvt9hhulth73drshp66dpdxxakj3a5uz",
			"reward": {
				"denom": "ucmdx",
				"amount": "59156"
			}
		},
		{
			"address": "comdex18jlc56v4gxux4wcah5njcnkxv6lde2hn8leqt7",
			"reward": {
				"denom": "ucmdx",
				"amount": "93317"
			}
		},
		{
			"address": "comdex18nqkyt7hcwdq2qwgs8e3yen2lgnc3u849w43r7",
			"reward": {
				"denom": "ucmdx",
				"amount": "10237"
			}
		},
		{
			"address": "comdex18npp0s2tffdx29qc0wrw8d48xrgcx2apf40e4z",
			"reward": {
				"denom": "ucmdx",
				"amount": "15068"
			}
		},
		{
			"address": "comdex18nyx6kgwqlplsmqtg643h72mf0hg60p0l5xef7",
			"reward": {
				"denom": "ucmdx",
				"amount": "1732"
			}
		},
		{
			"address": "comdex18n8h9a7d07esmcq9pnng9dgxxm0fg5uzzpjkxz",
			"reward": {
				"denom": "ucmdx",
				"amount": "1501"
			}
		},
		{
			"address": "comdex18n8a7vmzs7ywrtdezsdt7cy9g5en25s5292v5s",
			"reward": {
				"denom": "ucmdx",
				"amount": "176"
			}
		},
		{
			"address": "comdex18ntfznfg62dyhl5wsh6zvc7acq2c57hyyhplkj",
			"reward": {
				"denom": "ucmdx",
				"amount": "1190"
			}
		},
		{
			"address": "comdex18nd3c4gcrje6hrwqew7nyqjak39aps0yj6qvs4",
			"reward": {
				"denom": "ucmdx",
				"amount": "8040"
			}
		},
		{
			"address": "comdex18nw2r20637zeepq2jle9jltp55fdvxhjtxc6wj",
			"reward": {
				"denom": "ucmdx",
				"amount": "144"
			}
		},
		{
			"address": "comdex18n0y689ap6cpn9ch9zury47dfyas7leg6h3tqm",
			"reward": {
				"denom": "ucmdx",
				"amount": "354"
			}
		},
		{
			"address": "comdex18n0j6vw9gp3r40wrewk48mwyrrm4tjl3s55qm6",
			"reward": {
				"denom": "ucmdx",
				"amount": "18599"
			}
		},
		{
			"address": "comdex18nszxha39n6x42ththv2m3t66rzvefsnvezgtd",
			"reward": {
				"denom": "ucmdx",
				"amount": "142"
			}
		},
		{
			"address": "comdex18n3vlcthxv5s2s4zgjdrjr5a6lw9e7tmq3rygw",
			"reward": {
				"denom": "ucmdx",
				"amount": "4073"
			}
		},
		{
			"address": "comdex18n33ywygvrrzqspewsh7f5c4dkp00le05xhmvz",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex18nnj6l7rn5mjc68s9e32luqx6ek9dz5ty7a5m7",
			"reward": {
				"denom": "ucmdx",
				"amount": "70"
			}
		},
		{
			"address": "comdex18n5f3weawsldmddxr8c6ffd97yuz9xp437ajqp",
			"reward": {
				"denom": "ucmdx",
				"amount": "4812"
			}
		},
		{
			"address": "comdex18nk707axx4tkwc4ce6cgnx8mlnyfsr8ttcz7vw",
			"reward": {
				"denom": "ucmdx",
				"amount": "531"
			}
		},
		{
			"address": "comdex18nhxaevvwfmanjqfqg363xu98rgfxezlqflx3r",
			"reward": {
				"denom": "ucmdx",
				"amount": "6149"
			}
		},
		{
			"address": "comdex18n69c3f8f4nmjyxg3n7rznq7sld6jnka5l56c5",
			"reward": {
				"denom": "ucmdx",
				"amount": "1237"
			}
		},
		{
			"address": "comdex18nm59wfyjghqcjx0x0zu5yz6ruf6qshpxpvw94",
			"reward": {
				"denom": "ucmdx",
				"amount": "949"
			}
		},
		{
			"address": "comdex18nuq93paw6njulne8pk5rmdpwtv2xtqkjpe4tu",
			"reward": {
				"denom": "ucmdx",
				"amount": "1384"
			}
		},
		{
			"address": "comdex18nawnka82e7zrkuunyddeyd99x0frenmxm3gyt",
			"reward": {
				"denom": "ucmdx",
				"amount": "65"
			}
		},
		{
			"address": "comdex18nl5lvh97h52gtwzzs4dp3mhrpcsys9arnzmvu",
			"reward": {
				"denom": "ucmdx",
				"amount": "61553"
			}
		},
		{
			"address": "comdex185p6x2xvclccwr9jvc9qgtysjlfhl787dt82fu",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex185z0c3uentr0wpvgv2f7q894646zt5gwvw30t2",
			"reward": {
				"denom": "ucmdx",
				"amount": "13956"
			}
		},
		{
			"address": "comdex185zcsm5x5r78jn2wt2au7m2ygyt45gh5jpk6tn",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex185rq54sg76ld0laf2v7zjey4keakrqaqauvc87",
			"reward": {
				"denom": "ucmdx",
				"amount": "17766"
			}
		},
		{
			"address": "comdex185yv4ka7ymn3f6meq7k3g42nxrgfeykj90gth7",
			"reward": {
				"denom": "ucmdx",
				"amount": "13151"
			}
		},
		{
			"address": "comdex185yd6en9k2zh06gtgwpdmcr8ddqk8tg4fhwam8",
			"reward": {
				"denom": "ucmdx",
				"amount": "3524"
			}
		},
		{
			"address": "comdex1859dye4uzl773fcvmrve6tnraxasnp2skclmcx",
			"reward": {
				"denom": "ucmdx",
				"amount": "576"
			}
		},
		{
			"address": "comdex185xr9cuf5f2u2tuhuq4fjw2guf0pxrdrwxw24n",
			"reward": {
				"denom": "ucmdx",
				"amount": "757"
			}
		},
		{
			"address": "comdex18589sa66fvhleskyz6fv7u6j0e6rte9xk9vk8c",
			"reward": {
				"denom": "ucmdx",
				"amount": "14"
			}
		},
		{
			"address": "comdex1852gkzprzcqxftk8etrqen9rgm0j8s48gtf038",
			"reward": {
				"denom": "ucmdx",
				"amount": "1643"
			}
		},
		{
			"address": "comdex185vfvhy7xj0ee8u3qmj4e9c9dsn300zpvj8eug",
			"reward": {
				"denom": "ucmdx",
				"amount": "72539"
			}
		},
		{
			"address": "comdex185ddspcfygexhcejeqy437q9zhac2zftdvtd0w",
			"reward": {
				"denom": "ucmdx",
				"amount": "166"
			}
		},
		{
			"address": "comdex1850rdjq4sf488sad4qsj64sj2s8h53z07czm03",
			"reward": {
				"denom": "ucmdx",
				"amount": "2268"
			}
		},
		{
			"address": "comdex1850cmyf023m3pqdvhknnfcsg03pj7r8r0mmet2",
			"reward": {
				"denom": "ucmdx",
				"amount": "142"
			}
		},
		{
			"address": "comdex1853phfjzqczmdvs9meltz5l3hkxksru37dzcl9",
			"reward": {
				"denom": "ucmdx",
				"amount": "347491"
			}
		},
		{
			"address": "comdex18539wnkwuyl0ltwpgrd4dt8zq85fjgmk8eklmd",
			"reward": {
				"denom": "ucmdx",
				"amount": "16"
			}
		},
		{
			"address": "comdex185jzean9qr88j9wd6lykjm0m9ka73a8wpwgl93",
			"reward": {
				"denom": "ucmdx",
				"amount": "4291"
			}
		},
		{
			"address": "comdex185440fvp384ssjq7a0q50k4jtt69yfgjwr8f2e",
			"reward": {
				"denom": "ucmdx",
				"amount": "4260"
			}
		},
		{
			"address": "comdex185k6aaxudayk8n5w0s6hkctnjj5r5rngf6ur3f",
			"reward": {
				"denom": "ucmdx",
				"amount": "9289"
			}
		},
		{
			"address": "comdex185cd5utlg02pw2jvz7fdfqgn0ud8gu6hntx8eg",
			"reward": {
				"denom": "ucmdx",
				"amount": "1690"
			}
		},
		{
			"address": "comdex185eyku8zusq4mukcl4724yuv9stcu2lehaveuv",
			"reward": {
				"denom": "ucmdx",
				"amount": "3151"
			}
		},
		{
			"address": "comdex1856rtz9tq5hz40elah4fp6c6qzwsk6r0fsdam5",
			"reward": {
				"denom": "ucmdx",
				"amount": "179"
			}
		},
		{
			"address": "comdex1856tlqceeevgcg698q82mhndzrvtlgqzqrcqlk",
			"reward": {
				"denom": "ucmdx",
				"amount": "6856"
			}
		},
		{
			"address": "comdex1856mterrqagdj3w38u2zrg93849v0ey34e23sp",
			"reward": {
				"denom": "ucmdx",
				"amount": "17961"
			}
		},
		{
			"address": "comdex185uhy9m6sqmvhweqy6jqm8hkxqtkd67xwu6004",
			"reward": {
				"denom": "ucmdx",
				"amount": "2641"
			}
		},
		{
			"address": "comdex1857v89ezmes6lsjwner5fe3u3aajk2cfaxfve5",
			"reward": {
				"denom": "ucmdx",
				"amount": "1061"
			}
		},
		{
			"address": "comdex184p2gkyhl0mc8x4nvvrrzqzzdz0st2gya2ng9a",
			"reward": {
				"denom": "ucmdx",
				"amount": "10347"
			}
		},
		{
			"address": "comdex184zvd9xdydvlndve78zj5wqpe478lnukfc6xdd",
			"reward": {
				"denom": "ucmdx",
				"amount": "24577"
			}
		},
		{
			"address": "comdex184zm6tug27h386y3fajgn39lf9keslstf0evta",
			"reward": {
				"denom": "ucmdx",
				"amount": "5307"
			}
		},
		{
			"address": "comdex184r5gzewfxxcm2n0kpas4plvl24lpm2r9tem5t",
			"reward": {
				"denom": "ucmdx",
				"amount": "8639"
			}
		},
		{
			"address": "comdex1849aeae0zp80s9e27cvpanxvqhj9j5tcl8jdhx",
			"reward": {
				"denom": "ucmdx",
				"amount": "25156"
			}
		},
		{
			"address": "comdex184tkxyprazepegkrp26frpp2h4j3syc3wmrvap",
			"reward": {
				"denom": "ucmdx",
				"amount": "384"
			}
		},
		{
			"address": "comdex184wap298zlfsq4adqn2t3tzgdzmqyxqlhdd6rp",
			"reward": {
				"denom": "ucmdx",
				"amount": "52"
			}
		},
		{
			"address": "comdex18402uxqshsujk42clftk5lfs9j4lhaskl76frn",
			"reward": {
				"denom": "ucmdx",
				"amount": "651"
			}
		},
		{
			"address": "comdex1840j8ayh2x7dj75fmhn7fwthy8jtkg80deus6k",
			"reward": {
				"denom": "ucmdx",
				"amount": "148"
			}
		},
		{
			"address": "comdex184sxruylzpe362rqqpxkk5xj9au6remlykemdn",
			"reward": {
				"denom": "ucmdx",
				"amount": "4332"
			}
		},
		{
			"address": "comdex1843r0v49jcday5ak62n8wv47f5avmxgq793au5",
			"reward": {
				"denom": "ucmdx",
				"amount": "2706"
			}
		},
		{
			"address": "comdex184ns0f5eqrss5t0vysk8ez6gxdfua0eklxvfmd",
			"reward": {
				"denom": "ucmdx",
				"amount": "16434"
			}
		},
		{
			"address": "comdex1846asz0qfrl3ktwyf9u2wkpzmgcs9klgxnt3cn",
			"reward": {
				"denom": "ucmdx",
				"amount": "182"
			}
		},
		{
			"address": "comdex184m5yur2w2fevtx850hpecqfausrc62ev2k60s",
			"reward": {
				"denom": "ucmdx",
				"amount": "32353"
			}
		},
		{
			"address": "comdex184uxwpj7vfk3xfdnxl7hzlsga0dtgcpf6kl3nm",
			"reward": {
				"denom": "ucmdx",
				"amount": "190498"
			}
		},
		{
			"address": "comdex184ue0drpfg05c7rg7uzcan2pwwgj4h9uynmapw",
			"reward": {
				"denom": "ucmdx",
				"amount": "895063"
			}
		},
		{
			"address": "comdex18kyhv2hrplak0emy3nfcqd6ac8lucvq2wt6m09",
			"reward": {
				"denom": "ucmdx",
				"amount": "8493"
			}
		},
		{
			"address": "comdex18k9y9cky4khcxu7svh8avyywa4xcgafc3caw8k",
			"reward": {
				"denom": "ucmdx",
				"amount": "25575"
			}
		},
		{
			"address": "comdex18kx06wn4s2a23sr2n798hh80ceq62u0ax93zgk",
			"reward": {
				"denom": "ucmdx",
				"amount": "991"
			}
		},
		{
			"address": "comdex18kx07q9knsfnyvud7zxm74wch435fkqgevc4sl",
			"reward": {
				"denom": "ucmdx",
				"amount": "638"
			}
		},
		{
			"address": "comdex18kxlcsn0jmwmnjkfp3rrevzt3k7raaefl5tkkk",
			"reward": {
				"denom": "ucmdx",
				"amount": "41"
			}
		},
		{
			"address": "comdex18kglw4m60363w8gjm6glcpnjgqaw8vdyavtpc7",
			"reward": {
				"denom": "ucmdx",
				"amount": "1907"
			}
		},
		{
			"address": "comdex18kfsu2qw34crm4fsjeze6n0nv8hhv7ll0408dr",
			"reward": {
				"denom": "ucmdx",
				"amount": "4934"
			}
		},
		{
			"address": "comdex18ktd4nayqq2m687pv5m4s834dcnfc3tvewk32x",
			"reward": {
				"denom": "ucmdx",
				"amount": "51"
			}
		},
		{
			"address": "comdex18kwympkj6mysy4sd3hqy0ewe6jz5cmp59jg9rk",
			"reward": {
				"denom": "ucmdx",
				"amount": "143"
			}
		},
		{
			"address": "comdex18kw6wv03lfk0usx30da9auzhx7ym5ue6222gqf",
			"reward": {
				"denom": "ucmdx",
				"amount": "1190"
			}
		},
		{
			"address": "comdex18k0wh2n2lqxp85h9rrp2qufuyt3etspdv354v4",
			"reward": {
				"denom": "ucmdx",
				"amount": "116"
			}
		},
		{
			"address": "comdex18knqcqyfpcxgz6wupk3karcpukqdh937w07nwd",
			"reward": {
				"denom": "ucmdx",
				"amount": "6960"
			}
		},
		{
			"address": "comdex18ke7g8dtvhnypgkhhrec8xg7mvaeulld5u7gxl",
			"reward": {
				"denom": "ucmdx",
				"amount": "691"
			}
		},
		{
			"address": "comdex18hpm55dz3juhcweagtnpel9qmy5mgl8nyp987s",
			"reward": {
				"denom": "ucmdx",
				"amount": "825"
			}
		},
		{
			"address": "comdex18hxmwtzeyjr7qzxtxq4cksajl4aspfjfu258w9",
			"reward": {
				"denom": "ucmdx",
				"amount": "9071"
			}
		},
		{
			"address": "comdex18hf79r67jyueez70htz7t3cjpxuq6m5g9f09vm",
			"reward": {
				"denom": "ucmdx",
				"amount": "7154"
			}
		},
		{
			"address": "comdex18hs3u4wm66k8f7xnak6xh7zghcxewdvkvdxks0",
			"reward": {
				"denom": "ucmdx",
				"amount": "771"
			}
		},
		{
			"address": "comdex18h3napq0pcy7sahav089lj863qx2mn9ahtnqc8",
			"reward": {
				"denom": "ucmdx",
				"amount": "258"
			}
		},
		{
			"address": "comdex18hkpa8yz2kt2fucnh53te9898k7pgww2998pp6",
			"reward": {
				"denom": "ucmdx",
				"amount": "14199"
			}
		},
		{
			"address": "comdex18hhw82dsrv4xy0r9xh3ld2z4hkn3lnakrtj9gq",
			"reward": {
				"denom": "ucmdx",
				"amount": "138"
			}
		},
		{
			"address": "comdex18hek3kesv03zgq9mxpxrfvf0pz57ms85mfqql3",
			"reward": {
				"denom": "ucmdx",
				"amount": "371"
			}
		},
		{
			"address": "comdex18h6znhp68m64ryh6wfal8yyc9782w4kl7z4fe2",
			"reward": {
				"denom": "ucmdx",
				"amount": "354"
			}
		},
		{
			"address": "comdex18h7y7c4wvaxn4axjnx0aq0zzzu2r0jefhmnjwl",
			"reward": {
				"denom": "ucmdx",
				"amount": "347"
			}
		},
		{
			"address": "comdex18cq0vtf3gz5ea3fwdfs8jnkafnxx34nm2ftu4z",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex18crh2zlxusj0pv6fg6yhawxqavmefzf2dsch7w",
			"reward": {
				"denom": "ucmdx",
				"amount": "17738"
			}
		},
		{
			"address": "comdex18cyayu3m69uhpdstq6pq80wz0qepgcmtfrvqvy",
			"reward": {
				"denom": "ucmdx",
				"amount": "19235"
			}
		},
		{
			"address": "comdex18cxu889yp5r87cp9arejl0stnd9p4pytvk0jze",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex18c83ermq7tr42azemt3jm735amfpqrwljf3t9e",
			"reward": {
				"denom": "ucmdx",
				"amount": "9797"
			}
		},
		{
			"address": "comdex18cgmudj0w6yrck0pzkk4pnf5z62tuv0waptrv4",
			"reward": {
				"denom": "ucmdx",
				"amount": "991"
			}
		},
		{
			"address": "comdex18cdsntyaezx3hk4sv0rfyk5t0plzgddurlct4v",
			"reward": {
				"denom": "ucmdx",
				"amount": "49013"
			}
		},
		{
			"address": "comdex18cwj2yn2560daducx5flv22k5y8xr8p7x7k7eg",
			"reward": {
				"denom": "ucmdx",
				"amount": "173"
			}
		},
		{
			"address": "comdex18c0grhdx96lw2u5t9qchl390n5weu9znzah45g",
			"reward": {
				"denom": "ucmdx",
				"amount": "20252"
			}
		},
		{
			"address": "comdex18c0u2n8620s584239xaajvgg8677dhnhjca560",
			"reward": {
				"denom": "ucmdx",
				"amount": "2841"
			}
		},
		{
			"address": "comdex18cj9kguj7upxxzef6xu7tde8s02ln3c3sg762u",
			"reward": {
				"denom": "ucmdx",
				"amount": "23374"
			}
		},
		{
			"address": "comdex18cnux0fjk2h32ts9u3dfr6dxyk4mj85gj6u3p9",
			"reward": {
				"denom": "ucmdx",
				"amount": "2984"
			}
		},
		{
			"address": "comdex18c43qqudknvvm8e2yt48l4x49vmlg82h9y4n7e",
			"reward": {
				"denom": "ucmdx",
				"amount": "4444"
			}
		},
		{
			"address": "comdex18cksurlfdsthuhf8vgm6ps0f8k2hqr497kclrc",
			"reward": {
				"denom": "ucmdx",
				"amount": "1247"
			}
		},
		{
			"address": "comdex18ckjkml8aed9glnmt4lgstv45x084s98sv9jnx",
			"reward": {
				"denom": "ucmdx",
				"amount": "6134"
			}
		},
		{
			"address": "comdex18ch2tskkcm2ejkt9mp2z29624dzq37005yp8wz",
			"reward": {
				"denom": "ucmdx",
				"amount": "165"
			}
		},
		{
			"address": "comdex18cmfry7khzu00mjrsgrzgqvdry0840f029quws",
			"reward": {
				"denom": "ucmdx",
				"amount": "13159"
			}
		},
		{
			"address": "comdex18cmt27xqn42dejcreh2hruudtqk9hedtsy3xkg",
			"reward": {
				"denom": "ucmdx",
				"amount": "214751"
			}
		},
		{
			"address": "comdex18cauz4ydtvkvnhsmgvnyhfaf499y6qqnrgx4xy",
			"reward": {
				"denom": "ucmdx",
				"amount": "8540"
			}
		},
		{
			"address": "comdex18eq9txpm4amym7kxx6w49s0szwgm3fxra4nadq",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex18e9g3kfkzc9n9aet5c9eqwvt6ay5raarg05t49",
			"reward": {
				"denom": "ucmdx",
				"amount": "12140"
			}
		},
		{
			"address": "comdex18e8vtvqnawvpjcuzaxf8dm5dfs6s5mpnuew4d3",
			"reward": {
				"denom": "ucmdx",
				"amount": "4003"
			}
		},
		{
			"address": "comdex18efn4s7e9sv5m67hqj0j43veryldrgdc06s6d9",
			"reward": {
				"denom": "ucmdx",
				"amount": "179"
			}
		},
		{
			"address": "comdex18ece6q4ds40stu24w6zh307vy5zpjfssllneyj",
			"reward": {
				"denom": "ucmdx",
				"amount": "884"
			}
		},
		{
			"address": "comdex18e6wtx8c3zach2a7wx3jf4n3yxhgrddhpt8yze",
			"reward": {
				"denom": "ucmdx",
				"amount": "52"
			}
		},
		{
			"address": "comdex18eaxywacv0g25v3799zw44s8h8pg8yycwycglg",
			"reward": {
				"denom": "ucmdx",
				"amount": "6472"
			}
		},
		{
			"address": "comdex18ea3vr0m4p7cnrvu03wqutu4ggc603ytntcysx",
			"reward": {
				"denom": "ucmdx",
				"amount": "1452"
			}
		},
		{
			"address": "comdex186y4suz4z87v6s6ecd6552alnef4fdmuchzp2c",
			"reward": {
				"denom": "ucmdx",
				"amount": "46875"
			}
		},
		{
			"address": "comdex1869e9sjmv2arftja9fthp8z5p72ntlyedwhs2d",
			"reward": {
				"denom": "ucmdx",
				"amount": "1475"
			}
		},
		{
			"address": "comdex186g3e3ewqxmfd7s4nv3z6glucmrp2aurvn8gj3",
			"reward": {
				"denom": "ucmdx",
				"amount": "271"
			}
		},
		{
			"address": "comdex186gj77s4gqvykqldd9egt76tmjn8k7p856nv9j",
			"reward": {
				"denom": "ucmdx",
				"amount": "66"
			}
		},
		{
			"address": "comdex186tr9qhp9md9c6yg8c0xtt3w5kr54xkfm6xdrs",
			"reward": {
				"denom": "ucmdx",
				"amount": "1786"
			}
		},
		{
			"address": "comdex186vrvcdwnv20jkh98x48a7yuevkymjrnz9cfn5",
			"reward": {
				"denom": "ucmdx",
				"amount": "2476"
			}
		},
		{
			"address": "comdex186de56xwv35alsh9llxcwt8n50d2lkf8af65pw",
			"reward": {
				"denom": "ucmdx",
				"amount": "1025"
			}
		},
		{
			"address": "comdex186s3ltqtq8mcdgdgyscat3ny9jgpscq92yk0w7",
			"reward": {
				"denom": "ucmdx",
				"amount": "20704"
			}
		},
		{
			"address": "comdex186sj2kq43zzeyn85pngf8rqdjqt647ygc4f53t",
			"reward": {
				"denom": "ucmdx",
				"amount": "17713"
			}
		},
		{
			"address": "comdex186nc7f5790hy2klqa6e6vy8ts4x52vekyyytkr",
			"reward": {
				"denom": "ucmdx",
				"amount": "26454"
			}
		},
		{
			"address": "comdex1864060m7x8sxegarghngvdxxum6yyg8aplh52f",
			"reward": {
				"denom": "ucmdx",
				"amount": "715"
			}
		},
		{
			"address": "comdex186ezrgr9yuqv3vcgx4c2qv0tkw2y0uzx6jp0mt",
			"reward": {
				"denom": "ucmdx",
				"amount": "134610"
			}
		},
		{
			"address": "comdex186mmwn5qnf2xpmshffz6rrs8pa7gzp3qw755sc",
			"reward": {
				"denom": "ucmdx",
				"amount": "4880"
			}
		},
		{
			"address": "comdex186umvct5qghk3vrlg8hed4ftsecmce07c05k9t",
			"reward": {
				"denom": "ucmdx",
				"amount": "176"
			}
		},
		{
			"address": "comdex186avevsqxkkk5p8h8x2fpefua58arakv0fdrr7",
			"reward": {
				"denom": "ucmdx",
				"amount": "173"
			}
		},
		{
			"address": "comdex186aemtlgrf5prcw4k5zjhxsadjzljgtfe6604l",
			"reward": {
				"denom": "ucmdx",
				"amount": "7558"
			}
		},
		{
			"address": "comdex18mqfafetckqqq75658flkhvvd3ygvxrd9swrgs",
			"reward": {
				"denom": "ucmdx",
				"amount": "882"
			}
		},
		{
			"address": "comdex18mp0vs9970pnannlf2g46zyv9d6gffglz6u3v0",
			"reward": {
				"denom": "ucmdx",
				"amount": "2227"
			}
		},
		{
			"address": "comdex18mzjqryx9aw3svux73jduxzyqpnzyx8p9lsau9",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex18myuwxmresr7xzuta5yqzj4pksr9wdvaergu2l",
			"reward": {
				"denom": "ucmdx",
				"amount": "13762"
			}
		},
		{
			"address": "comdex18msw2yscztxy57e6yf5v2yk07qj0uupsprga0t",
			"reward": {
				"denom": "ucmdx",
				"amount": "1699"
			}
		},
		{
			"address": "comdex18mkuwykd46xwypl4262gcpv7fllpgza558xumc",
			"reward": {
				"denom": "ucmdx",
				"amount": "7966"
			}
		},
		{
			"address": "comdex18mcmc9m8gqp5g4620c6t39ap3zwapszasw0ct7",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex18mme533j0gahs99kspc2g6gyklxmdlpenkgs6c",
			"reward": {
				"denom": "ucmdx",
				"amount": "1169"
			}
		},
		{
			"address": "comdex18uqe743vkrlauyhvfyyz6nk38qv0qss6ad4u0w",
			"reward": {
				"denom": "ucmdx",
				"amount": "3894"
			}
		},
		{
			"address": "comdex18uyrnkyg0pp3ts5g9w5k0fugmyhgy34ayfl0xs",
			"reward": {
				"denom": "ucmdx",
				"amount": "284"
			}
		},
		{
			"address": "comdex18u9622u6tw2z6evksfrq9mj4xdsu0rf07qxcy8",
			"reward": {
				"denom": "ucmdx",
				"amount": "1753"
			}
		},
		{
			"address": "comdex18ug8yac5h8j7vqwtssqpz8gp8cq6spwg77ncxp",
			"reward": {
				"denom": "ucmdx",
				"amount": "4294"
			}
		},
		{
			"address": "comdex18ufwenfk744c400p44nqrlfgs793yqau620fpx",
			"reward": {
				"denom": "ucmdx",
				"amount": "706"
			}
		},
		{
			"address": "comdex18utawmxge83fyg3agnjflvc88pdwa8r2pr6p7c",
			"reward": {
				"denom": "ucmdx",
				"amount": "20365"
			}
		},
		{
			"address": "comdex18us5yfw0eyndhekj38pktwg20t6wn9jmn93he8",
			"reward": {
				"denom": "ucmdx",
				"amount": "15"
			}
		},
		{
			"address": "comdex18u3kespvwnf6ggx69etlr8mg80hw33sdyd3mm6",
			"reward": {
				"denom": "ucmdx",
				"amount": "1130264"
			}
		},
		{
			"address": "comdex18u3hdpr8p993gnhj72r8gj9hhekktpm2l3zzha",
			"reward": {
				"denom": "ucmdx",
				"amount": "13292"
			}
		},
		{
			"address": "comdex18u5jr32x3yna6r488zpx0p6dhndtwd8c707duu",
			"reward": {
				"denom": "ucmdx",
				"amount": "6922"
			}
		},
		{
			"address": "comdex18umucrt44w3xtdj53q96cv64zkk63l4ngmdqkc",
			"reward": {
				"denom": "ucmdx",
				"amount": "7861"
			}
		},
		{
			"address": "comdex18uaadv5qehjxeqdhdhk5d69a5dvmjvs4cfkc6m",
			"reward": {
				"denom": "ucmdx",
				"amount": "393"
			}
		},
		{
			"address": "comdex18u7sm5kfzz4spgkjzjsmqlfce9v6hgtmwa9j5n",
			"reward": {
				"denom": "ucmdx",
				"amount": "17815"
			}
		},
		{
			"address": "comdex18uln99d672yeeftyuvc6eq6gg85fdw4kkplwch",
			"reward": {
				"denom": "ucmdx",
				"amount": "0"
			}
		},
		{
			"address": "comdex18aqw29mqhxsysqm76axnycy535j7n3065sxrcy",
			"reward": {
				"denom": "ucmdx",
				"amount": "19715"
			}
		},
		{
			"address": "comdex18aqjyx9md8hh0f8vvz0gu0jzufmtcxtc9hrndu",
			"reward": {
				"denom": "ucmdx",
				"amount": "1403"
			}
		},
		{
			"address": "comdex18aprs8ywf6tushjxmx5jgg2gkk465nnk8kn25w",
			"reward": {
				"denom": "ucmdx",
				"amount": "59"
			}
		},
		{
			"address": "comdex18ayfuqplf2u3q7lscw60w428j9wlcmdqhexegl",
			"reward": {
				"denom": "ucmdx",
				"amount": "571"
			}
		},
		{
			"address": "comdex18agclqfeqhtn5yjyaw586ugnk7fn87vv7qs2q7",
			"reward": {
				"denom": "ucmdx",
				"amount": "323"
			}
		},
		{
			"address": "comdex18afk9m67j0rsztg7waqazfvw77rpt6dsy7pgks",
			"reward": {
				"denom": "ucmdx",
				"amount": "137"
			}
		},
		{
			"address": "comdex18adpqq0ppa6q8fsr6ssvrkhhcyeze0pu2glmcr",
			"reward": {
				"denom": "ucmdx",
				"amount": "2015"
			}
		},
		{
			"address": "comdex18adyt689n6ry58pf238fqddlun0xwmarxugemj",
			"reward": {
				"denom": "ucmdx",
				"amount": "6403"
			}
		},
		{
			"address": "comdex18adxx2qcg9nf9uw48v462yys0wqvrq7kk54uf6",
			"reward": {
				"denom": "ucmdx",
				"amount": "703"
			}
		},
		{
			"address": "comdex18a35908a7sglgypgtjfw6nagvnknuw6mwcdf3h",
			"reward": {
				"denom": "ucmdx",
				"amount": "144"
			}
		},
		{
			"address": "comdex18a3579c6dcjlsf4dcj3qp7dna3xga2gjfzycc4",
			"reward": {
				"denom": "ucmdx",
				"amount": "103016"
			}
		},
		{
			"address": "comdex18ajc7euypp0f79eak6q3yutw33ge5y8suk925m",
			"reward": {
				"denom": "ucmdx",
				"amount": "13939"
			}
		},
		{
			"address": "comdex18an5azytvketyt6xyq8wn352egyxlshfdcm8ln",
			"reward": {
				"denom": "ucmdx",
				"amount": "5989"
			}
		},
		{
			"address": "comdex18ank35aua8ftetpmvcsxwdqy4phelteecpua9x",
			"reward": {
				"denom": "ucmdx",
				"amount": "5259"
			}
		},
		{
			"address": "comdex18a5wwnd93gpwelf5emtews80y3q8gs8dajtz0h",
			"reward": {
				"denom": "ucmdx",
				"amount": "3587"
			}
		},
		{
			"address": "comdex18ahhlmgmtmh5gu74p0h5m2xawmqsskarkwpads",
			"reward": {
				"denom": "ucmdx",
				"amount": "1475"
			}
		},
		{
			"address": "comdex18ahlvkf59susr6padt3n5r23mntuqyuvrkxmg7",
			"reward": {
				"denom": "ucmdx",
				"amount": "2817"
			}
		},
		{
			"address": "comdex18aek36rf0cv273x2mu4gg7zv96yjxqzg932x0c",
			"reward": {
				"denom": "ucmdx",
				"amount": "14205"
			}
		},
		{
			"address": "comdex18a6g84sqjj7eaq4g9jrf8d7g0ksc4qnkaqu5ym",
			"reward": {
				"denom": "ucmdx",
				"amount": "14137"
			}
		},
		{
			"address": "comdex18a60kqzkl27vgpgxdu225g0w7gzthf7sp7nxdv",
			"reward": {
				"denom": "ucmdx",
				"amount": "1229"
			}
		},
		{
			"address": "comdex18aav27875wy0zlr34vt88nnwlr2prx09nu0dx5",
			"reward": {
				"denom": "ucmdx",
				"amount": "438"
			}
		},
		{
			"address": "comdex1879jdpn6ee86g30c76khqyaq6vtdn5u6wyfqwa",
			"reward": {
				"denom": "ucmdx",
				"amount": "2057"
			}
		},
		{
			"address": "comdex1879k48hqxgy30amszz66yhpahh0nyvr53wu8xw",
			"reward": {
				"denom": "ucmdx",
				"amount": "7563"
			}
		},
		{
			"address": "comdex1878cs7vj3734tlgcyl9k5ys6lydxdvfexmnms7",
			"reward": {
				"denom": "ucmdx",
				"amount": "87"
			}
		},
		{
			"address": "comdex187d5a7yutv4rr8f3lhaq04cvdgfg5frp0ur775",
			"reward": {
				"denom": "ucmdx",
				"amount": "1689"
			}
		},
		{
			"address": "comdex187wy2grtqlxpp6j3gspyaej3c58efargu6qnhf",
			"reward": {
				"denom": "ucmdx",
				"amount": "13619"
			}
		},
		{
			"address": "comdex187wtyvmx9t4vzgceafrwxhd83udstpelsdzzgq",
			"reward": {
				"denom": "ucmdx",
				"amount": "1248"
			}
		},
		{
			"address": "comdex187j8gffp4rqnult2warg3w9zxkz5g0jxfadcz8",
			"reward": {
				"denom": "ucmdx",
				"amount": "50470"
			}
		},
		{
			"address": "comdex18759k5tz85n576ecyptxune2cucc3cg707uugj",
			"reward": {
				"denom": "ucmdx",
				"amount": "86"
			}
		},
		{
			"address": "comdex1875d64ct7x4zrrspk98lry2v6h5lk6cxsrgxj9",
			"reward": {
				"denom": "ucmdx",
				"amount": "1835"
			}
		},
		{
			"address": "comdex1874j0f9c3cy442cca49dr54lyc9kf4ukeuv8gx",
			"reward": {
				"denom": "ucmdx",
				"amount": "242153"
			}
		},
		{
			"address": "comdex187kuqm2vd4mxc6k9l75dxnf9zuua96d7nt59xc",
			"reward": {
				"denom": "ucmdx",
				"amount": "6876"
			}
		},
		{
			"address": "comdex187ml0m564desz6c97qr2kjws80r4ygp6nypzq3",
			"reward": {
				"denom": "ucmdx",
				"amount": "8831"
			}
		},
		{
			"address": "comdex187u838u9gdyfrtqnz3m8jmhcr5ttrw4q2f4fk6",
			"reward": {
				"denom": "ucmdx",
				"amount": "30922"
			}
		},
		{
			"address": "comdex187a7cnj2z22pyqvnp62aq6kfwn6ky2eeek9lss",
			"reward": {
				"denom": "ucmdx",
				"amount": "5244"
			}
		},
		{
			"address": "comdex187llah7hvkdq4atwgkqqjwe52dcs2amweq8uut",
			"reward": {
				"denom": "ucmdx",
				"amount": "7118"
			}
		},
		{
			"address": "comdex18lqc83qlhthlc0u0j4ux3n0rhtehy9az2c0uye",
			"reward": {
				"denom": "ucmdx",
				"amount": "3448"
			}
		},
		{
			"address": "comdex18lz45r3jznse0rxrfn79mxeepwt8kt0y4h0e6e",
			"reward": {
				"denom": "ucmdx",
				"amount": "7103"
			}
		},
		{
			"address": "comdex18lrvwhtsrx2l28my6cugwx2t7u7sxuafy3434u",
			"reward": {
				"denom": "ucmdx",
				"amount": "11666"
			}
		},
		{
			"address": "comdex18lxfswz35rp6ky09q4n8frcqthsz7x0dqa6p8t",
			"reward": {
				"denom": "ucmdx",
				"amount": "1773"
			}
		},
		{
			"address": "comdex18l8kpncz0gpdr6446mtjt9kpxs0s0mx3epnmga",
			"reward": {
				"denom": "ucmdx",
				"amount": "788"
			}
		},
		{
			"address": "comdex18lg8shcdrpym7ly70m8le9z85sufhav0mlhqks",
			"reward": {
				"denom": "ucmdx",
				"amount": "149"
			}
		},
		{
			"address": "comdex18lvemx64k28jzy537rcywka3nc4gfsquxr0qvn",
			"reward": {
				"denom": "ucmdx",
				"amount": "323"
			}
		},
		{
			"address": "comdex18lv7d7593yuasrhprv3jk2wvcw29llvtvz0c8y",
			"reward": {
				"denom": "ucmdx",
				"amount": "6909"
			}
		},
		{
			"address": "comdex18l07sx5mu76ntk7assa0uy4kz8ucaw82ae68ly",
			"reward": {
				"denom": "ucmdx",
				"amount": "22117"
			}
		},
		{
			"address": "comdex18lsk29nhrj802pjt9nl83npr56teacmr98s39l",
			"reward": {
				"denom": "ucmdx",
				"amount": "17690"
			}
		},
		{
			"address": "comdex18l3qxr2dpxl0cgguu6h36fck5azevhlvcefj37",
			"reward": {
				"denom": "ucmdx",
				"amount": "2072"
			}
		},
		{
			"address": "comdex18lnxwzu4dtxlam4sxm44jmhpyalleeyszdzz2u",
			"reward": {
				"denom": "ucmdx",
				"amount": "1412"
			}
		},
		{
			"address": "comdex18l4f32zejq30u08l808c2y2lu69znzykw6xk6k",
			"reward": {
				"denom": "ucmdx",
				"amount": "17704"
			}
		},
		{
			"address": "comdex18l4uw8s9rzadwg4ccr9k3anqlfqtm4zk970g43",
			"reward": {
				"denom": "ucmdx",
				"amount": "5949"
			}
		},
		{
			"address": "comdex18lcfmmupqjysfg4stvhp7f24hxv96nuwlmwsea",
			"reward": {
				"denom": "ucmdx",
				"amount": "33011"
			}
		},
		{
			"address": "comdex18levssru80xmxalhml45zytl3vm9jlq9plr4ay",
			"reward": {
				"denom": "ucmdx",
				"amount": "204"
			}
		},
		{
			"address": "comdex18lau6n0f8jzz6cr9k2mcm4lyylkwwunw80svkv",
			"reward": {
				"denom": "ucmdx",
				"amount": "59976"
			}
		},
		{
			"address": "comdex18l7yc3wdrdt2lx66lxmx30jlv5hcfj50at0y7u",
			"reward": {
				"denom": "ucmdx",
				"amount": "141"
			}
		},
		{
			"address": "comdex1gqqd6eeedgjtrp83v8zn3c8wgdm29med5s55ay",
			"reward": {
				"denom": "ucmdx",
				"amount": "52"
			}
		},
		{
			"address": "comdex1gqrymq6942cnur56hpf7txvk2ec7njwqfusljr",
			"reward": {
				"denom": "ucmdx",
				"amount": "179"
			}
		},
		{
			"address": "comdex1gq8uv8s94gls578x5hgzq4gq4gwhe6tds2vrm8",
			"reward": {
				"denom": "ucmdx",
				"amount": "179"
			}
		},
		{
			"address": "comdex1gqgnmdtup77z692rmga0sxply654d8qnnd49pp",
			"reward": {
				"denom": "ucmdx",
				"amount": "173"
			}
		},
		{
			"address": "comdex1gq2r0w7hlwa3efl8kkg66708rkfsamfappypjj",
			"reward": {
				"denom": "ucmdx",
				"amount": "0"
			}
		},
		{
			"address": "comdex1gqvt9k75yvuyrxyply3xqz2afn2xymaauqdexf",
			"reward": {
				"denom": "ucmdx",
				"amount": "71930"
			}
		},
		{
			"address": "comdex1gqv4ly0xgzd7w3hkt037ksrwfhezleva6klzk8",
			"reward": {
				"denom": "ucmdx",
				"amount": "167"
			}
		},
		{
			"address": "comdex1gq0srrtqtc3l79jtcde5p74hgmd336jmd0f9h3",
			"reward": {
				"denom": "ucmdx",
				"amount": "544"
			}
		},
		{
			"address": "comdex1gqs0efwlh97rx5kdxtypkykvvr0fqqemgl8g03",
			"reward": {
				"denom": "ucmdx",
				"amount": "13535"
			}
		},
		{
			"address": "comdex1gqseup6q22pj5f9p2w3uxz2uw7758r7ex5vygc",
			"reward": {
				"denom": "ucmdx",
				"amount": "202"
			}
		},
		{
			"address": "comdex1gq5y0ayyrzdvqflv8ff39v8vtfznesugfver52",
			"reward": {
				"denom": "ucmdx",
				"amount": "1762"
			}
		},
		{
			"address": "comdex1gq5vp7dwg2h6s9ekhkgjnly30wfpjq5nm4e22m",
			"reward": {
				"denom": "ucmdx",
				"amount": "1414"
			}
		},
		{
			"address": "comdex1gqcsp5u2j3uew074en3jv8ekcn9y0xssph8cp3",
			"reward": {
				"denom": "ucmdx",
				"amount": "841"
			}
		},
		{
			"address": "comdex1gqe93zn0w39q2pudgcyg3kexq96kxdh2agkg4y",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1gqewyaj0zqeysh0jck5tgemcugv3j45a5f6nlg",
			"reward": {
				"denom": "ucmdx",
				"amount": "7132"
			}
		},
		{
			"address": "comdex1gqew8e5jkrfe5cahzpyscx6led549nrj2r0hq7",
			"reward": {
				"denom": "ucmdx",
				"amount": "83608"
			}
		},
		{
			"address": "comdex1gqe4hffdvjtxzq4553a0l0tgm8g0qyg47p2l6m",
			"reward": {
				"denom": "ucmdx",
				"amount": "395"
			}
		},
		{
			"address": "comdex1gqemfhrysns3n7k3qer5pcl9ynp5zjvf3h8y2t",
			"reward": {
				"denom": "ucmdx",
				"amount": "334"
			}
		},
		{
			"address": "comdex1gq6pd9q8px97720arytdngc7zhyrxp4xjv23kz",
			"reward": {
				"denom": "ucmdx",
				"amount": "27740"
			}
		},
		{
			"address": "comdex1gq65gjqeyy4alxnnn6y2spc7hjj2mhcjm3ptqn",
			"reward": {
				"denom": "ucmdx",
				"amount": "503"
			}
		},
		{
			"address": "comdex1gq7wl8l2yf7qrt2qv2c2um22atdsxdgpykx3cp",
			"reward": {
				"denom": "ucmdx",
				"amount": "143"
			}
		},
		{
			"address": "comdex1gqlvl2qqectkgz7uz6gkjust00m7q63t6sptsn",
			"reward": {
				"denom": "ucmdx",
				"amount": "302"
			}
		},
		{
			"address": "comdex1gppsq2vsvm8a95u8gknav3z5wkge5cu4xzke63",
			"reward": {
				"denom": "ucmdx",
				"amount": "66"
			}
		},
		{
			"address": "comdex1gp94e2st9x758dc996vlcw4uaflphsgk2mqhg8",
			"reward": {
				"denom": "ucmdx",
				"amount": "88"
			}
		},
		{
			"address": "comdex1gpxqaju250ze5w73v968pjut8eup6gpsgwl35u",
			"reward": {
				"denom": "ucmdx",
				"amount": "1698"
			}
		},
		{
			"address": "comdex1gpxwpy9k64pjvgn2gv6axj9exumzq98uq3yknp",
			"reward": {
				"denom": "ucmdx",
				"amount": "5364"
			}
		},
		{
			"address": "comdex1gp8cv3rzz0h96nfv7svn56qysu6vgphluff207",
			"reward": {
				"denom": "ucmdx",
				"amount": "11"
			}
		},
		{
			"address": "comdex1gp0xkm4ud68h90lcq3amjrfnsnzykmp32c3cyw",
			"reward": {
				"denom": "ucmdx",
				"amount": "317"
			}
		},
		{
			"address": "comdex1gpsltwleqlexmvmnp4sd6e6pdygqhethxtv36z",
			"reward": {
				"denom": "ucmdx",
				"amount": "7542"
			}
		},
		{
			"address": "comdex1gpnh7c09r26p3298wnet23ltsennvv8635wrwa",
			"reward": {
				"denom": "ucmdx",
				"amount": "21251"
			}
		},
		{
			"address": "comdex1gp5gcxwjkna8crxcgz9g9yeza7ljyaw3dknqxt",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1gp5w2lk3g6zun22n9mpfpt0hatvgmrm0x2s7sf",
			"reward": {
				"denom": "ucmdx",
				"amount": "1301"
			}
		},
		{
			"address": "comdex1gpc77p04dasu0l4scc4mrqpmlk62rll353p2gn",
			"reward": {
				"denom": "ucmdx",
				"amount": "16744"
			}
		},
		{
			"address": "comdex1gp60pf2zhmh052ry9hdvwjmsy80xa4fw2q6ryt",
			"reward": {
				"denom": "ucmdx",
				"amount": "333"
			}
		},
		{
			"address": "comdex1gp674enk6ldw3728qqavn0dse5uxsun3ucxgj0",
			"reward": {
				"denom": "ucmdx",
				"amount": "12726"
			}
		},
		{
			"address": "comdex1gpm6jmudpml40sysdxeryphy30d9s8zr48t3tw",
			"reward": {
				"denom": "ucmdx",
				"amount": "144037"
			}
		},
		{
			"address": "comdex1gpucwsmsuyxpj2sz06tdgv9npx0jeuawgzzhe3",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1gpl9tmh354m99kq9h8yzr7yunlspxdgd06spv0",
			"reward": {
				"denom": "ucmdx",
				"amount": "61926"
			}
		},
		{
			"address": "comdex1gpla8pke6xcrq06d7rfdyzyx5sjs6jqpaesu0p",
			"reward": {
				"denom": "ucmdx",
				"amount": "71025"
			}
		},
		{
			"address": "comdex1gzpq6fs7eujkl8p5ez2cmu0g2wk0xzy0r7rpql",
			"reward": {
				"denom": "ucmdx",
				"amount": "2825"
			}
		},
		{
			"address": "comdex1gzr8qv7dp3dzx0a6mya78z5jr3qfeyt9zf4qza",
			"reward": {
				"denom": "ucmdx",
				"amount": "57728"
			}
		},
		{
			"address": "comdex1gzrg66th9zt89q9l36xgg3vpdjztr4s8vvytxc",
			"reward": {
				"denom": "ucmdx",
				"amount": "163646"
			}
		},
		{
			"address": "comdex1gzyz2hayj3t4hs6qwlvz8qph3sqjwkerhu8wvx",
			"reward": {
				"denom": "ucmdx",
				"amount": "771"
			}
		},
		{
			"address": "comdex1gzy0f2ggfjk6mrlv96nlur6awqdm2hd5cdamd7",
			"reward": {
				"denom": "ucmdx",
				"amount": "18288"
			}
		},
		{
			"address": "comdex1gzy3rpek8elan5vp3e3vtt2mc35dg8lw7d295f",
			"reward": {
				"denom": "ucmdx",
				"amount": "12535"
			}
		},
		{
			"address": "comdex1gzgzlj4hndwfk6qdzm3pmqhzjmwupcn7qe8gpe",
			"reward": {
				"denom": "ucmdx",
				"amount": "13895"
			}
		},
		{
			"address": "comdex1gzttw349nu9ftjke3a5gs8ydstekcjykx5udjn",
			"reward": {
				"denom": "ucmdx",
				"amount": "15"
			}
		},
		{
			"address": "comdex1gzd30x4zwxyt50xzhnk67akd0dlgmx9x7r69z0",
			"reward": {
				"denom": "ucmdx",
				"amount": "5278"
			}
		},
		{
			"address": "comdex1gzd5mjktsh94d6sg6c0sl8yrux6vgmme6mwgkd",
			"reward": {
				"denom": "ucmdx",
				"amount": "1138"
			}
		},
		{
			"address": "comdex1gzw93kttk9nl7x9qc6d437h7c2fzv9wkedaml6",
			"reward": {
				"denom": "ucmdx",
				"amount": "34971"
			}
		},
		{
			"address": "comdex1gzsc7t6pwsxk2a44jxql8namldq6y6hxc23c4y",
			"reward": {
				"denom": "ucmdx",
				"amount": "19859"
			}
		},
		{
			"address": "comdex1gznv7wjyg8uxvq08pym4r3elnh6yj9nupv58qd",
			"reward": {
				"denom": "ucmdx",
				"amount": "38660"
			}
		},
		{
			"address": "comdex1gz457d2h69khqkjs2npdm9j7am756evxxeu05j",
			"reward": {
				"denom": "ucmdx",
				"amount": "139"
			}
		},
		{
			"address": "comdex1gzkyp0vxsplu8e3dv7fpgmdzrc67ya539wrz0a",
			"reward": {
				"denom": "ucmdx",
				"amount": "31228"
			}
		},
		{
			"address": "comdex1gzk54axu4y3l2vy9434z2grfn86qsd35j0dw2e",
			"reward": {
				"denom": "ucmdx",
				"amount": "182412"
			}
		},
		{
			"address": "comdex1gzhykp96cgyysfs86gzjycgr7vk0cxgn8cn6sw",
			"reward": {
				"denom": "ucmdx",
				"amount": "307"
			}
		},
		{
			"address": "comdex1gz6r9d278t2ppcljltpgfrp03pdmxjhnj2vqdt",
			"reward": {
				"denom": "ucmdx",
				"amount": "3689"
			}
		},
		{
			"address": "comdex1gz6y8kg22crcg7d3va3q936llmufjvhpuddl66",
			"reward": {
				"denom": "ucmdx",
				"amount": "251"
			}
		},
		{
			"address": "comdex1gz6tlyxsf03hxv492xtztmsmu4dssxnmzu56xn",
			"reward": {
				"denom": "ucmdx",
				"amount": "2470"
			}
		},
		{
			"address": "comdex1gzar22al0ku0mvggy4r5jaza62lksv33z9vak8",
			"reward": {
				"denom": "ucmdx",
				"amount": "2622932"
			}
		},
		{
			"address": "comdex1gz7r387xq3r9t5mczqx9wz96kfa0yqpf4ufeyn",
			"reward": {
				"denom": "ucmdx",
				"amount": "70293"
			}
		},
		{
			"address": "comdex1grqyasrtq7aks8xnja6trnc22tyj5zdsk54jkg",
			"reward": {
				"denom": "ucmdx",
				"amount": "2010"
			}
		},
		{
			"address": "comdex1grqalstfxyfm0f8eqnre6plfxt9jmyjma2eduh",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1gr9lft9zsnury0mun7y0e6ve9afn0u2cj5d6tm",
			"reward": {
				"denom": "ucmdx",
				"amount": "8885"
			}
		},
		{
			"address": "comdex1gr2f8ma0fcd9e65fcrnmqsxphjsuvmypf8ex0t",
			"reward": {
				"denom": "ucmdx",
				"amount": "71110"
			}
		},
		{
			"address": "comdex1grtpn36qrfmrfexu6vj3s5pgd4x5hcvhzz3yvy",
			"reward": {
				"denom": "ucmdx",
				"amount": "34821"
			}
		},
		{
			"address": "comdex1grt4aycsq5c7zmcke0wku8sr8xkdyyxzlxces9",
			"reward": {
				"denom": "ucmdx",
				"amount": "1704"
			}
		},
		{
			"address": "comdex1gr0zryw5h0kmxdj24yh3t0y8c4xfel27ed0r95",
			"reward": {
				"denom": "ucmdx",
				"amount": "1777"
			}
		},
		{
			"address": "comdex1gr4xqrhmn44a7uu4ddhdulyrngxs0w2n0rl4wh",
			"reward": {
				"denom": "ucmdx",
				"amount": "31791"
			}
		},
		{
			"address": "comdex1grk94kszu20as74n5dkqrnn9pks88t55xghzda",
			"reward": {
				"denom": "ucmdx",
				"amount": "1422"
			}
		},
		{
			"address": "comdex1grkwn5pz3dcmdeejk2me2zre9r4264hunh7vh3",
			"reward": {
				"denom": "ucmdx",
				"amount": "13359"
			}
		},
		{
			"address": "comdex1gr68jugpza2sd3guk963vl2kj5we7z590tdvhz",
			"reward": {
				"denom": "ucmdx",
				"amount": "5289"
			}
		},
		{
			"address": "comdex1grmp0wktm02hw7ec6xx7k53wsrhm05chlhtj9x",
			"reward": {
				"denom": "ucmdx",
				"amount": "7098"
			}
		},
		{
			"address": "comdex1grmvphmzyzk56mnwjec9y3dhtdvzaftg7ddwvg",
			"reward": {
				"denom": "ucmdx",
				"amount": "36924"
			}
		},
		{
			"address": "comdex1grmn0csvrr57ka5mnfj8a2wvn9eaf6peyv7gqj",
			"reward": {
				"denom": "ucmdx",
				"amount": "3794"
			}
		},
		{
			"address": "comdex1grmapu0sq608utmpeshnklhe5hjt3gsnp9f82q",
			"reward": {
				"denom": "ucmdx",
				"amount": "25059"
			}
		},
		{
			"address": "comdex1gypyfsd4pegh85uys96g6m60t4v3crfvfsz030",
			"reward": {
				"denom": "ucmdx",
				"amount": "12006"
			}
		},
		{
			"address": "comdex1gyrgg9w52x4u9yqkvy4ps7tqs3fr0q8ny32qep",
			"reward": {
				"denom": "ucmdx",
				"amount": "65"
			}
		},
		{
			"address": "comdex1gy2ne7m62uer4h5s4e7xlfq7aeem5zpwyg7gnu",
			"reward": {
				"denom": "ucmdx",
				"amount": "8816"
			}
		},
		{
			"address": "comdex1gywt266ry22c5n0p9q8pk4c4wd7r487ew8v8ar",
			"reward": {
				"denom": "ucmdx",
				"amount": "9765"
			}
		},
		{
			"address": "comdex1gysjqxzuwjre7ut2kzj55dl4d75n03m6pnhrlc",
			"reward": {
				"denom": "ucmdx",
				"amount": "145"
			}
		},
		{
			"address": "comdex1gynfacsqdmrdqx7e04w8sqyfsukkg7jgasel2u",
			"reward": {
				"denom": "ucmdx",
				"amount": "407"
			}
		},
		{
			"address": "comdex1gy59ekywem4qq5a24lmkqytqk724yn92ma8spp",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1gycr4gn9kf90yuwqhedd396yt6e5usc9fqp2nq",
			"reward": {
				"denom": "ucmdx",
				"amount": "1855"
			}
		},
		{
			"address": "comdex1gym8xffd5msakjeqdtzd5f70lpxp2aj0krmu9u",
			"reward": {
				"denom": "ucmdx",
				"amount": "14026"
			}
		},
		{
			"address": "comdex1gymcz07zk8hjfkfkyzmj4x4s4nl50vj6qtn3hn",
			"reward": {
				"denom": "ucmdx",
				"amount": "1445"
			}
		},
		{
			"address": "comdex1gyuencw7ukwz6a6k6302egs32ue6mcwv6n5845",
			"reward": {
				"denom": "ucmdx",
				"amount": "23420"
			}
		},
		{
			"address": "comdex1gyldf5duvm4rn8d04r9n7t5xfzhk5d898jfxng",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1g9q69rgpjq7qke95mqg4fwg6qnede0jfrmmqdq",
			"reward": {
				"denom": "ucmdx",
				"amount": "4120"
			}
		},
		{
			"address": "comdex1g9z2dwcvtw2kd7uux0y6rl7mrmzmy9crmaghxz",
			"reward": {
				"denom": "ucmdx",
				"amount": "284"
			}
		},
		{
			"address": "comdex1g9z6x9e75falw6m04mw3p2u6n6z5f4s7m7pv5r",
			"reward": {
				"denom": "ucmdx",
				"amount": "19347"
			}
		},
		{
			"address": "comdex1g9zluy5ye6vfy9vknr4vrtv4ysqsxq5zx586w7",
			"reward": {
				"denom": "ucmdx",
				"amount": "86711"
			}
		},
		{
			"address": "comdex1g9rm3jmntm006yw7c44443gzg7d0lwh69pc45s",
			"reward": {
				"denom": "ucmdx",
				"amount": "18229"
			}
		},
		{
			"address": "comdex1g9y3g9x9jf0xcjtue6hhelg4age6d0u8ntn2kh",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1g98xsjmgkzztd4mrshlkhh32prx6p7yxafx403",
			"reward": {
				"denom": "ucmdx",
				"amount": "8863"
			}
		},
		{
			"address": "comdex1g98jg0rcfngj6a8kcusn4s9p0pmh6l83avv487",
			"reward": {
				"denom": "ucmdx",
				"amount": "173"
			}
		},
		{
			"address": "comdex1g9gwc0pww5qwjdewx7jmkt9mfdch68ecxqf4va",
			"reward": {
				"denom": "ucmdx",
				"amount": "125929"
			}
		},
		{
			"address": "comdex1g9t95arlf4cmsxt6mq7haec9thkx3ar7h6fgvq",
			"reward": {
				"denom": "ucmdx",
				"amount": "454"
			}
		},
		{
			"address": "comdex1g9ceq73t0s9d2vspy9qdljak2c5lqhp5f7t4n6",
			"reward": {
				"denom": "ucmdx",
				"amount": "21611"
			}
		},
		{
			"address": "comdex1g96zjr90sew2tsewp8q7fx8utg00efm3pzjyht",
			"reward": {
				"denom": "ucmdx",
				"amount": "6952"
			}
		},
		{
			"address": "comdex1g97we496ef0mt39jse82r67lmy2l7hw6falhet",
			"reward": {
				"denom": "ucmdx",
				"amount": "35027"
			}
		},
		{
			"address": "comdex1g9lrzxyhgqslkdvtxvrgkw0sunckmde5v7x443",
			"reward": {
				"denom": "ucmdx",
				"amount": "36"
			}
		},
		{
			"address": "comdex1gxq9up6vf8na5ddau9py380dgfjecgc04h7p9l",
			"reward": {
				"denom": "ucmdx",
				"amount": "1025"
			}
		},
		{
			"address": "comdex1gxq9azf4kjrn5ljd9rz2pml482mdkttrzsfr4p",
			"reward": {
				"denom": "ucmdx",
				"amount": "1217"
			}
		},
		{
			"address": "comdex1gxq5thehfmzzt20cduqf2lzwxarfzec5qrpe92",
			"reward": {
				"denom": "ucmdx",
				"amount": "22769"
			}
		},
		{
			"address": "comdex1gxp9nt44e6mzm6f02kjhd3689mzewuv08jevxe",
			"reward": {
				"denom": "ucmdx",
				"amount": "1769"
			}
		},
		{
			"address": "comdex1gxzzgvt266t4yec44nvm32crzujzasv4pf5qam",
			"reward": {
				"denom": "ucmdx",
				"amount": "1905"
			}
		},
		{
			"address": "comdex1gxy8aqzhmncu5j89t99qtdj3av8la27lsxdqyh",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1gx9ukgjwnn36jeyfhrvmskm4dm55n0cfv9e9u3",
			"reward": {
				"denom": "ucmdx",
				"amount": "268"
			}
		},
		{
			"address": "comdex1gxxss8umpwvqtyaw4043fdft38u3v9mw2ykmsf",
			"reward": {
				"denom": "ucmdx",
				"amount": "171"
			}
		},
		{
			"address": "comdex1gxgcmpes4amm98490wtddv45trzrhfsnnm6npu",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1gxfju0h3xzmsqn6kwcrj7vgmfjlmu0n03ckd5h",
			"reward": {
				"denom": "ucmdx",
				"amount": "6233"
			}
		},
		{
			"address": "comdex1gxsmetwlkn2c2p9cp9d0s4gepqnj7pd6r6aed0",
			"reward": {
				"denom": "ucmdx",
				"amount": "1516"
			}
		},
		{
			"address": "comdex1gxng8fvjfslsrz75qglc7lx6nk9dmlgsx2mgks",
			"reward": {
				"denom": "ucmdx",
				"amount": "1217"
			}
		},
		{
			"address": "comdex1gx5cx6prh66wkqtpmssd84q5dmvuw8cvz056u3",
			"reward": {
				"denom": "ucmdx",
				"amount": "347"
			}
		},
		{
			"address": "comdex1gx4jc7rafx68hsxuc24jar6d9e4ht9pmyecsar",
			"reward": {
				"denom": "ucmdx",
				"amount": "189"
			}
		},
		{
			"address": "comdex1gxk5m6yywrpwx86a62v5geqgclg3n2rra9kqq8",
			"reward": {
				"denom": "ucmdx",
				"amount": "440"
			}
		},
		{
			"address": "comdex1gxh99pjlga34d5qkalyvpt8jhqlwceuxwxfslg",
			"reward": {
				"denom": "ucmdx",
				"amount": "12336"
			}
		},
		{
			"address": "comdex1gxhekedqyc863grt2rahwp9sldh8v5g4v7us9v",
			"reward": {
				"denom": "ucmdx",
				"amount": "4704"
			}
		},
		{
			"address": "comdex1gxe0aq7wd3c0ufqy8gmw9sg29ed4q9vtmx2n4f",
			"reward": {
				"denom": "ucmdx",
				"amount": "26922"
			}
		},
		{
			"address": "comdex1gxmr2c8al68h8az2qpdx4ss9efs97f3xtznfn0",
			"reward": {
				"denom": "ucmdx",
				"amount": "6110"
			}
		},
		{
			"address": "comdex1gxu3sufn9tacyhv67sljglkeqm7zfuverncu8m",
			"reward": {
				"denom": "ucmdx",
				"amount": "8636"
			}
		},
		{
			"address": "comdex1g8qynt25y7g87ah6dnwjmwt7t9fx6mjppz52e3",
			"reward": {
				"denom": "ucmdx",
				"amount": "2184"
			}
		},
		{
			"address": "comdex1g8qtcdmk3ncq3v9y42pwpljcpg83kuwch92erz",
			"reward": {
				"denom": "ucmdx",
				"amount": "15146"
			}
		},
		{
			"address": "comdex1g8ph905lfyqzxclycwhumylywtfqyl0r62vlfs",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1g8rhwatvkjr3dcplu5397r92hehkapcflp3f7r",
			"reward": {
				"denom": "ucmdx",
				"amount": "1993"
			}
		},
		{
			"address": "comdex1g8yq39gcysxam0hzw5qqwj67p8fce0864yyrem",
			"reward": {
				"denom": "ucmdx",
				"amount": "2817"
			}
		},
		{
			"address": "comdex1g8tausx6qfmd03gxcl8cdejtzrq2pgw90uacp5",
			"reward": {
				"denom": "ucmdx",
				"amount": "1751"
			}
		},
		{
			"address": "comdex1g8da45wq3pzrm6eunl3enlklk4hnvn82229qz4",
			"reward": {
				"denom": "ucmdx",
				"amount": "5127"
			}
		},
		{
			"address": "comdex1g80l3p523v0qak88k39u8rxc3fq0jadr2cwaaf",
			"reward": {
				"denom": "ucmdx",
				"amount": "47888"
			}
		},
		{
			"address": "comdex1g8s2yyv8yajjuthtfjkx3yfknt4tdt2zwcy9jv",
			"reward": {
				"denom": "ucmdx",
				"amount": "2130"
			}
		},
		{
			"address": "comdex1g8ng9zygf9em3r0wslds5e9jy5f383ezkk4677",
			"reward": {
				"denom": "ucmdx",
				"amount": "7282"
			}
		},
		{
			"address": "comdex1g8n34e7h6llfwweuppwr8jw23zx2lkdqw9gg5q",
			"reward": {
				"denom": "ucmdx",
				"amount": "17442"
			}
		},
		{
			"address": "comdex1g84y0v3qj4wg7879dpya4g0mdgwvj6r89qn5fx",
			"reward": {
				"denom": "ucmdx",
				"amount": "294"
			}
		},
		{
			"address": "comdex1g8kedd7f4c3vltulyy6fvmhe5gf576sc60e482",
			"reward": {
				"denom": "ucmdx",
				"amount": "1413"
			}
		},
		{
			"address": "comdex1g8hwq8mz9mpjf8a2d6rhl37ue2kuefrs5y65e8",
			"reward": {
				"denom": "ucmdx",
				"amount": "26902"
			}
		},
		{
			"address": "comdex1g8hslcd5uruyn5lyp28rlxedausx3a8rryd8gj",
			"reward": {
				"denom": "ucmdx",
				"amount": "80640"
			}
		},
		{
			"address": "comdex1g8ctqa5jfp253knh4nnjwy4aw0g4h5wwyrt7kr",
			"reward": {
				"denom": "ucmdx",
				"amount": "53010"
			}
		},
		{
			"address": "comdex1g8c5qkdzec4ehlfdleer7yw660k2yc2wst6jmv",
			"reward": {
				"denom": "ucmdx",
				"amount": "1768"
			}
		},
		{
			"address": "comdex1g8652xl5kdeevc4eaepqgw76pdrz0f92v87t9f",
			"reward": {
				"denom": "ucmdx",
				"amount": "30"
			}
		},
		{
			"address": "comdex1g8uc95dtgp937w04zc3ktwgm7yzvu6s5q697wh",
			"reward": {
				"denom": "ucmdx",
				"amount": "3386"
			}
		},
		{
			"address": "comdex1g873jd5uen9vf08ql6zlkfgg9udq4n603753jn",
			"reward": {
				"denom": "ucmdx",
				"amount": "14394"
			}
		},
		{
			"address": "comdex1g8l25h7yn9ty6p0wfqp8vznzavxcupd5tc7pfe",
			"reward": {
				"denom": "ucmdx",
				"amount": "410"
			}
		},
		{
			"address": "comdex1ggqnctdleprqskjsvuyhdszqaqqd9w00fnuvt9",
			"reward": {
				"denom": "ucmdx",
				"amount": "2988"
			}
		},
		{
			"address": "comdex1gg9f8wekkyh7nu54wnlvh0ryk8sqdfn0wgsm5x",
			"reward": {
				"denom": "ucmdx",
				"amount": "13175"
			}
		},
		{
			"address": "comdex1gggfd6sgg30m7mhsq7hwwm7nd8udzhvwgxd5qq",
			"reward": {
				"denom": "ucmdx",
				"amount": "176"
			}
		},
		{
			"address": "comdex1ggfgy669y3mlq7p9ujycen3qfk7gkfgjgmkkhx",
			"reward": {
				"denom": "ucmdx",
				"amount": "28606"
			}
		},
		{
			"address": "comdex1ggfgmwrtsug3vqv23dhgs25tg4nx02fexlv5jl",
			"reward": {
				"denom": "ucmdx",
				"amount": "1755"
			}
		},
		{
			"address": "comdex1gg2z2y3ncjwh0wd0q6e7tj2ecma04my0jv5lxk",
			"reward": {
				"denom": "ucmdx",
				"amount": "1426"
			}
		},
		{
			"address": "comdex1gg2fy45qtfm5gd4tahlryvkpm0gjktr3vu9g7t",
			"reward": {
				"denom": "ucmdx",
				"amount": "18774"
			}
		},
		{
			"address": "comdex1ggtx8flu47wt7n0wd8aa09zkyetlsjfvem9rvw",
			"reward": {
				"denom": "ucmdx",
				"amount": "4104"
			}
		},
		{
			"address": "comdex1ggvylcrejyt26kqc9z5k4k7nkq8cxnac0cx49a",
			"reward": {
				"denom": "ucmdx",
				"amount": "2536"
			}
		},
		{
			"address": "comdex1ggv0ve355lzlfwjr0gg5rwwgnce8qehsu9jete",
			"reward": {
				"denom": "ucmdx",
				"amount": "8581"
			}
		},
		{
			"address": "comdex1gg0v6eqev76kvmc350c0qtamydth0pnlfs3k4a",
			"reward": {
				"denom": "ucmdx",
				"amount": "326"
			}
		},
		{
			"address": "comdex1gg3e4prdf44q4r0dsxmjefrpmf60vas0ys3xvs",
			"reward": {
				"denom": "ucmdx",
				"amount": "136"
			}
		},
		{
			"address": "comdex1gg3ljec83ds4suutndxvezy3pgn7c82pw76nav",
			"reward": {
				"denom": "ucmdx",
				"amount": "71417"
			}
		},
		{
			"address": "comdex1gg5xw45tkzy2maepr4uf3kafat9ue45dxkcqm7",
			"reward": {
				"denom": "ucmdx",
				"amount": "592281"
			}
		},
		{
			"address": "comdex1gg4n6f36n8z59a2uahlt9lajzjelceyzg0n0hc",
			"reward": {
				"denom": "ucmdx",
				"amount": "42536"
			}
		},
		{
			"address": "comdex1gg4kml45cdcw4x7edp832all6fdepzr48cfjru",
			"reward": {
				"denom": "ucmdx",
				"amount": "197"
			}
		},
		{
			"address": "comdex1gg4mdc35djsvfpjry6ay6j6ez84rkdwpau8qga",
			"reward": {
				"denom": "ucmdx",
				"amount": "9033"
			}
		},
		{
			"address": "comdex1ggkzk0j99k8lmthehke0thv2jzygszwtzf0szl",
			"reward": {
				"denom": "ucmdx",
				"amount": "29362"
			}
		},
		{
			"address": "comdex1gghe9mgk703cur0egksda6fdsr7dl58mjrkdpc",
			"reward": {
				"denom": "ucmdx",
				"amount": "156"
			}
		},
		{
			"address": "comdex1ggcxyngdy2zk5ef3spxkyk0lavfxx4yhhuzs83",
			"reward": {
				"denom": "ucmdx",
				"amount": "276"
			}
		},
		{
			"address": "comdex1ggm2unxmnz9jc599gx5f5634htrd2tsped6smk",
			"reward": {
				"denom": "ucmdx",
				"amount": "275"
			}
		},
		{
			"address": "comdex1gfpct7cnsdq0xmpmw8sqjv5fvcemwj35u5ft9a",
			"reward": {
				"denom": "ucmdx",
				"amount": "68343"
			}
		},
		{
			"address": "comdex1gfr63tk0dqnyvfe25lz2hve42fgkpg34xqln4g",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1gf9yy3e3pa6qzhkq09j5g9v927mtpaldgm2fmx",
			"reward": {
				"denom": "ucmdx",
				"amount": "5328"
			}
		},
		{
			"address": "comdex1gffnjakl5gntmcq9yxlv77rqa3heg4v7rdcjd6",
			"reward": {
				"denom": "ucmdx",
				"amount": "14081"
			}
		},
		{
			"address": "comdex1gftgm5zjc90qdt863g3378j9ef82smzamnvep4",
			"reward": {
				"denom": "ucmdx",
				"amount": "169"
			}
		},
		{
			"address": "comdex1gft78kg6h8cc4uxe9vc9sr4dkerpy027c2aegk",
			"reward": {
				"denom": "ucmdx",
				"amount": "34"
			}
		},
		{
			"address": "comdex1gfwlesmnyv7mwj24w7kj9dy0ud6lvder0wxttu",
			"reward": {
				"denom": "ucmdx",
				"amount": "276"
			}
		},
		{
			"address": "comdex1gfsavmxauy8dmfanyfcnhy6lwnj2tqlh2y5we5",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1gfj5e5zp4dqzn5ffemmkd5tx68vxgffu268ez3",
			"reward": {
				"denom": "ucmdx",
				"amount": "2180"
			}
		},
		{
			"address": "comdex1gfje5jjtz37psn067ddyv7nhgkuu7g9lntxzkr",
			"reward": {
				"denom": "ucmdx",
				"amount": "7091"
			}
		},
		{
			"address": "comdex1gfnnpxnn8l9f3zagaawqpsv9f7umyfrkelxjhw",
			"reward": {
				"denom": "ucmdx",
				"amount": "132"
			}
		},
		{
			"address": "comdex1gf582nhn0ef0lw5xpfzf7nluu693464v0ddhxz",
			"reward": {
				"denom": "ucmdx",
				"amount": "19798"
			}
		},
		{
			"address": "comdex1gfc6ae8ldfnjluhdyeqxu5zv4wq4f9hfsptd64",
			"reward": {
				"denom": "ucmdx",
				"amount": "34"
			}
		},
		{
			"address": "comdex1gfesdtyd94k9nw2arzqy9kzkqwqwf0a3grvc27",
			"reward": {
				"denom": "ucmdx",
				"amount": "6621"
			}
		},
		{
			"address": "comdex1gf745ghez80h7j2v3zks3aayzw8nys8mdg2hgh",
			"reward": {
				"denom": "ucmdx",
				"amount": "848"
			}
		},
		{
			"address": "comdex1g2qg8uujk0276ns25q9lrnvknp67sclkqmjykj",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1g2pjpzsq3azt86emu6y383y2tskdewvkhey6hl",
			"reward": {
				"denom": "ucmdx",
				"amount": "3022"
			}
		},
		{
			"address": "comdex1g2pc3ujx0zeq3hp02spuwqeunw24hjaa35yn46",
			"reward": {
				"denom": "ucmdx",
				"amount": "1437"
			}
		},
		{
			"address": "comdex1g2z8ul7xq2754zwesess8x6d5xaay0upv37yd6",
			"reward": {
				"denom": "ucmdx",
				"amount": "11102"
			}
		},
		{
			"address": "comdex1g2zwsrtg50tq4ryrjfth62z76kpwwunvymufs9",
			"reward": {
				"denom": "ucmdx",
				"amount": "1308"
			}
		},
		{
			"address": "comdex1g2rrx3w5rvz6glmuy0trm7n608d0lm3k7k69jf",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1g2ylr0ht53mzgln0fnewrmgcxntmhvp96kpfj3",
			"reward": {
				"denom": "ucmdx",
				"amount": "7107"
			}
		},
		{
			"address": "comdex1g29tygfw3j5vptdzfarsrwuy0jw7kw2uezgkpw",
			"reward": {
				"denom": "ucmdx",
				"amount": "3034"
			}
		},
		{
			"address": "comdex1g29slg8le5ggf7m8760eytmrkwmmgy22gaxq0k",
			"reward": {
				"denom": "ucmdx",
				"amount": "1901"
			}
		},
		{
			"address": "comdex1g2xwrvr5zd37vcmydc5y4dz6s8qw04ctl3ehws",
			"reward": {
				"denom": "ucmdx",
				"amount": "16363"
			}
		},
		{
			"address": "comdex1g22etkwztpkgdmj6kgw8ghjen0f3kru3hh6gx6",
			"reward": {
				"denom": "ucmdx",
				"amount": "15"
			}
		},
		{
			"address": "comdex1g2vzecyrrvwc7l6cvnlh9ltv53alme3s0nx8q9",
			"reward": {
				"denom": "ucmdx",
				"amount": "120"
			}
		},
		{
			"address": "comdex1g2dtv9egd52ea2ahrecf4mjzwr3e5nssev8x26",
			"reward": {
				"denom": "ucmdx",
				"amount": "12626"
			}
		},
		{
			"address": "comdex1g2s2za5lapx9qkr5qneh8t0vr20alnv3vmr8du",
			"reward": {
				"denom": "ucmdx",
				"amount": "1427"
			}
		},
		{
			"address": "comdex1g238t79lec906hgnvpnyu4mywltl0kd6nn6qm8",
			"reward": {
				"denom": "ucmdx",
				"amount": "352"
			}
		},
		{
			"address": "comdex1g2ndtva5g8jtklch4nc36ksenqpuzlcrxqw9wr",
			"reward": {
				"denom": "ucmdx",
				"amount": "5704"
			}
		},
		{
			"address": "comdex1g25prkp9j9f7rngpppt680yjg04hrga5xewhfd",
			"reward": {
				"denom": "ucmdx",
				"amount": "35107"
			}
		},
		{
			"address": "comdex1g25yxdzthma53ew30lmns6nvxp0302gkqpytma",
			"reward": {
				"denom": "ucmdx",
				"amount": "198"
			}
		},
		{
			"address": "comdex1g25mm28ck57tkmvn90ayppd4rqzgn3z0afj2er",
			"reward": {
				"denom": "ucmdx",
				"amount": "61203"
			}
		},
		{
			"address": "comdex1g2kr88dzpwszku2t27ueducezuwgs7dw8h8f37",
			"reward": {
				"denom": "ucmdx",
				"amount": "1803"
			}
		},
		{
			"address": "comdex1g2mjkrugdjr0h75cwg9ml5r8ynltd4z9jgtjf6",
			"reward": {
				"denom": "ucmdx",
				"amount": "2901"
			}
		},
		{
			"address": "comdex1g2mhmq4terefkpt7hquk3d68dlqyexd575chd2",
			"reward": {
				"denom": "ucmdx",
				"amount": "33487"
			}
		},
		{
			"address": "comdex1g2ueljln9qpg29ngppylkadx4ql35xl3tv08wl",
			"reward": {
				"denom": "ucmdx",
				"amount": "4859"
			}
		},
		{
			"address": "comdex1g27wqpsu9uw3j7xjxzvl8fvj0x48rjccum7ppm",
			"reward": {
				"denom": "ucmdx",
				"amount": "101"
			}
		},
		{
			"address": "comdex1g2l5fatn57ddqr94mwzs8xlkvy0xf84ageuy6c",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1gtqlqh9ajc386k9deftddj5k5703hlfp4qh759",
			"reward": {
				"denom": "ucmdx",
				"amount": "1574"
			}
		},
		{
			"address": "comdex1gtpvmlk704wu43h9v8gh87e2wgnjyc2f9ul2s3",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1gtzjc0t4cs22gjks29fc8g4e55cmxj36nyzlha",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1gtrl9uf98cvu3suxyanz09nzjpfj0ewduym4d4",
			"reward": {
				"denom": "ucmdx",
				"amount": "70331"
			}
		},
		{
			"address": "comdex1gty5x97z63mxvdnmfh45d43z2fcenfnyzu9dzq",
			"reward": {
				"denom": "ucmdx",
				"amount": "87171"
			}
		},
		{
			"address": "comdex1gt9g0t7mmkhavcunqw03r8j5t8ugunwql2ce9x",
			"reward": {
				"denom": "ucmdx",
				"amount": "174"
			}
		},
		{
			"address": "comdex1gt9snm2p8q3hvkk77xva3j3tjrltnukdkh94l4",
			"reward": {
				"denom": "ucmdx",
				"amount": "1052"
			}
		},
		{
			"address": "comdex1gt8k590v9phyt7ucasf5djks5g2lwl47eptaxf",
			"reward": {
				"denom": "ucmdx",
				"amount": "19794"
			}
		},
		{
			"address": "comdex1gtfalq27vsk8fu6g7zhhvex76pwaev5wqjj43a",
			"reward": {
				"denom": "ucmdx",
				"amount": "2895"
			}
		},
		{
			"address": "comdex1gtdcphd2l35rq8mfrsvyz69mv07wsckpkeqhnu",
			"reward": {
				"denom": "ucmdx",
				"amount": "1114"
			}
		},
		{
			"address": "comdex1gt077c53an0yv4mxlzjpycuqhgls0kl5vku3jf",
			"reward": {
				"denom": "ucmdx",
				"amount": "5307"
			}
		},
		{
			"address": "comdex1gt3jhrfwnqlpzry0qavmhzwvq4d2wmpuw3z9x3",
			"reward": {
				"denom": "ucmdx",
				"amount": "6154"
			}
		},
		{
			"address": "comdex1gtnrsw32twpwy5rc6sq59gndvwj26ke7y86gf0",
			"reward": {
				"denom": "ucmdx",
				"amount": "3923"
			}
		},
		{
			"address": "comdex1gt4wkpyfpzzxpj2hpks94g6nrqwkavfhamhw5v",
			"reward": {
				"denom": "ucmdx",
				"amount": "60914"
			}
		},
		{
			"address": "comdex1gtmxvyatn509ndnnu3hc4tml6ck50x0h962phy",
			"reward": {
				"denom": "ucmdx",
				"amount": "6850"
			}
		},
		{
			"address": "comdex1gtuj8jeh0kcwjkvh8k5uuk0xsfsh2vxs5uvtx0",
			"reward": {
				"denom": "ucmdx",
				"amount": "1735"
			}
		},
		{
			"address": "comdex1gtap7t69y9tggyrrzc00sfpp2h4vgdxw3fcftv",
			"reward": {
				"denom": "ucmdx",
				"amount": "1717"
			}
		},
		{
			"address": "comdex1gtant7jp8fyfncrmpgm2dgnwe5lskz324pwyvn",
			"reward": {
				"denom": "ucmdx",
				"amount": "554"
			}
		},
		{
			"address": "comdex1gtlrpw0atvc3ydn2s2fpcja0n2stw4djgz7ewp",
			"reward": {
				"denom": "ucmdx",
				"amount": "1463"
			}
		},
		{
			"address": "comdex1gvzyger38xnj9mg7w2ecs2yprflcc4qk27w44r",
			"reward": {
				"denom": "ucmdx",
				"amount": "166"
			}
		},
		{
			"address": "comdex1gv9v4ngjeymlmts5279vsugu46g9tq0vxt0sa5",
			"reward": {
				"denom": "ucmdx",
				"amount": "6970"
			}
		},
		{
			"address": "comdex1gv9naxvkskmm533knsdujz9653dzcmf87kppdg",
			"reward": {
				"denom": "ucmdx",
				"amount": "4562"
			}
		},
		{
			"address": "comdex1gvgy6sh8qvqwj6x3mqpqwz9lca03960f3hcrhw",
			"reward": {
				"denom": "ucmdx",
				"amount": "747"
			}
		},
		{
			"address": "comdex1gv23034suxwklhr4q7fsw9dygph55ykmgfr8uu",
			"reward": {
				"denom": "ucmdx",
				"amount": "29975"
			}
		},
		{
			"address": "comdex1gv03gh7c205ar60khgtmkerfafetavhc67dpjc",
			"reward": {
				"denom": "ucmdx",
				"amount": "16"
			}
		},
		{
			"address": "comdex1gvjfwh63czaj30kqzk73vv34qegcvrcf6rgs0p",
			"reward": {
				"denom": "ucmdx",
				"amount": "997"
			}
		},
		{
			"address": "comdex1gvnm2f8cp8wmqs9lzvyrmud0z5ktptzd00qpqr",
			"reward": {
				"denom": "ucmdx",
				"amount": "151"
			}
		},
		{
			"address": "comdex1gv5u4hkjtu06a3veezvs2903af2r5pq5wem2hr",
			"reward": {
				"denom": "ucmdx",
				"amount": "1749"
			}
		},
		{
			"address": "comdex1gvkzfsggjsl4h5sg4npnchg7tj8zuvnm6xletx",
			"reward": {
				"denom": "ucmdx",
				"amount": "1246"
			}
		},
		{
			"address": "comdex1gvkxsljg77xm2qef2ccdpyw0uc3vgsph9c3erd",
			"reward": {
				"denom": "ucmdx",
				"amount": "692"
			}
		},
		{
			"address": "comdex1gvhta6rldu0s2jxgmvn8hvzg0rer0camscr7kn",
			"reward": {
				"denom": "ucmdx",
				"amount": "2566"
			}
		},
		{
			"address": "comdex1gv6fjfphfr8rdcsd26vcmxt2k55jv6f557992n",
			"reward": {
				"denom": "ucmdx",
				"amount": "65295"
			}
		},
		{
			"address": "comdex1gvm84pp3ayethwajz9fn06drlznjcyzygfa4hr",
			"reward": {
				"denom": "ucmdx",
				"amount": "284"
			}
		},
		{
			"address": "comdex1gdq3mdyv3twcgx7knr4kwrs8hdhvsg7p76e305",
			"reward": {
				"denom": "ucmdx",
				"amount": "1414"
			}
		},
		{
			"address": "comdex1gdrudpp2q8ce3hxrs425a6naa93a95p980lx5f",
			"reward": {
				"denom": "ucmdx",
				"amount": "114"
			}
		},
		{
			"address": "comdex1gdyauq823ckaqnmdenspwwxugc9c4pve39w24e",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1gd9qamw2jkq6vfv0de075vmtgwegqwl8lnpg5z",
			"reward": {
				"denom": "ucmdx",
				"amount": "19"
			}
		},
		{
			"address": "comdex1gd9k7q039ffhkeatknzg9wpgf4swkrkfq5espk",
			"reward": {
				"denom": "ucmdx",
				"amount": "5989"
			}
		},
		{
			"address": "comdex1gdt9d7pnxus3kqq2j3e8x45vxga4yvcwg65ly6",
			"reward": {
				"denom": "ucmdx",
				"amount": "173"
			}
		},
		{
			"address": "comdex1gdv7kzg2lvx6k5w6nxq84q8r9z2ujpvegep40v",
			"reward": {
				"denom": "ucmdx",
				"amount": "5351766"
			}
		},
		{
			"address": "comdex1gdjely4w9arkn9c8usdnwnqltjxwaek3zswvwd",
			"reward": {
				"denom": "ucmdx",
				"amount": "13240"
			}
		},
		{
			"address": "comdex1gdnp4axmu7z5kc0kre3yfsvmmhz8ynww5msugw",
			"reward": {
				"denom": "ucmdx",
				"amount": "183"
			}
		},
		{
			"address": "comdex1gdnnfzt2mdega5wkkp9uvfw98pqkfrlllq08xt",
			"reward": {
				"denom": "ucmdx",
				"amount": "7644"
			}
		},
		{
			"address": "comdex1gd5d4zp8rkrt9f33metha0ymah858zljtqlxug",
			"reward": {
				"denom": "ucmdx",
				"amount": "4516"
			}
		},
		{
			"address": "comdex1gd4cy4548vual7hmzlgfjean8d9xjrnzxjgldv",
			"reward": {
				"denom": "ucmdx",
				"amount": "10493"
			}
		},
		{
			"address": "comdex1gd47zmmcyl7fa2xy0gr9mnyk0u3vkc7zz8clju",
			"reward": {
				"denom": "ucmdx",
				"amount": "3545"
			}
		},
		{
			"address": "comdex1gdkxypp057g2mcjtyae9jezd8epqae9l3a4u8q",
			"reward": {
				"denom": "ucmdx",
				"amount": "716"
			}
		},
		{
			"address": "comdex1gdcjmjvmedu4zujycgl365vcf957s0pm5ynurn",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1gdmvcq07dk04ljcxkymz205hq6td60fktxw9zz",
			"reward": {
				"denom": "ucmdx",
				"amount": "1986"
			}
		},
		{
			"address": "comdex1gd7y8w9xzy4sjut5574cffu57ytsj5rgtf6hrz",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1gd76ykk669f253ju86uf63q0dpu8f6fx628vm4",
			"reward": {
				"denom": "ucmdx",
				"amount": "2897"
			}
		},
		{
			"address": "comdex1gdlg5u0acwme8y472c28mf6sf77w4vr7hsjq0e",
			"reward": {
				"denom": "ucmdx",
				"amount": "150"
			}
		},
		{
			"address": "comdex1gdlf9wh06zz8agjny5v87ggvzdhgwdhyqqdsh4",
			"reward": {
				"denom": "ucmdx",
				"amount": "4620"
			}
		},
		{
			"address": "comdex1gdlapchg4www4ueyq53fvzc3g2keh55szt3e0a",
			"reward": {
				"denom": "ucmdx",
				"amount": "1784"
			}
		},
		{
			"address": "comdex1gwqkt8rka8kg67x4f0jnsvw4hgv6rejg363n4w",
			"reward": {
				"denom": "ucmdx",
				"amount": "28847"
			}
		},
		{
			"address": "comdex1gwqe0u7tfws2zr9qr4jx9le2wn6v6wyy2jszps",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1gwr9ss8mqr2z6dc0madrufdkrg338akvpr26sm",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1gwx6zqk55z2llrxjez9zg78593hz4c7pevfdrt",
			"reward": {
				"denom": "ucmdx",
				"amount": "902"
			}
		},
		{
			"address": "comdex1gw2s88qc0q29ld8zy49eah0w3uw5gwtqu4usjz",
			"reward": {
				"denom": "ucmdx",
				"amount": "149"
			}
		},
		{
			"address": "comdex1gwtqpgfvyx8kz9a9zqmrze0tyt8ls5a0ar00d0",
			"reward": {
				"denom": "ucmdx",
				"amount": "3014"
			}
		},
		{
			"address": "comdex1gwtk204wuv674q7amr4822p2l37dxsyn3y52wr",
			"reward": {
				"denom": "ucmdx",
				"amount": "1515"
			}
		},
		{
			"address": "comdex1gwjfm0m7f80p2gxwsly8jdvxvqnnxejl8kcgfx",
			"reward": {
				"denom": "ucmdx",
				"amount": "304"
			}
		},
		{
			"address": "comdex1gwjtn7xggta9k4s6e8allyds0s3syex3kxu0jh",
			"reward": {
				"denom": "ucmdx",
				"amount": "17627"
			}
		},
		{
			"address": "comdex1gwnfjj65zhxmtnpxzy625cpx692lfjuwlqeg8p",
			"reward": {
				"denom": "ucmdx",
				"amount": "44413"
			}
		},
		{
			"address": "comdex1gwn2ac08ntz4qz9vxsae79y6gxvkyxpuqkgfz7",
			"reward": {
				"denom": "ucmdx",
				"amount": "2806"
			}
		},
		{
			"address": "comdex1gw6yu9jetj2vrht8qrf9wmpjdq4weuee7r35yd",
			"reward": {
				"denom": "ucmdx",
				"amount": "1693"
			}
		},
		{
			"address": "comdex1gwmw3zp3ju0ty2hu4td94jsgwzl4hc5s4cmnyz",
			"reward": {
				"denom": "ucmdx",
				"amount": "107518"
			}
		},
		{
			"address": "comdex1gwa7535t5uezsj0xu438vu5n9kaa7l5q7yelwz",
			"reward": {
				"denom": "ucmdx",
				"amount": "128815"
			}
		},
		{
			"address": "comdex1gwlqmqc9nrnuu8fanls76exyrzum9xjmrmh5hx",
			"reward": {
				"denom": "ucmdx",
				"amount": "11925"
			}
		},
		{
			"address": "comdex1g0ycd8ccup76kvj0wvnx2uhn5xc75w323a0dn9",
			"reward": {
				"denom": "ucmdx",
				"amount": "142"
			}
		},
		{
			"address": "comdex1g0xluyamawjq5g097qtzm840dcf3unhntalk67",
			"reward": {
				"denom": "ucmdx",
				"amount": "2010"
			}
		},
		{
			"address": "comdex1g0gkq4ad8g8elcmz9gputfwg786hanwa0xfpm4",
			"reward": {
				"denom": "ucmdx",
				"amount": "3772"
			}
		},
		{
			"address": "comdex1g00rfafry5k9ymjh6vpr677p274z3pxp5skyyn",
			"reward": {
				"denom": "ucmdx",
				"amount": "88"
			}
		},
		{
			"address": "comdex1g0jt2smscmjh8q6f85y5x3jf6w8wfgn7lmuw0p",
			"reward": {
				"denom": "ucmdx",
				"amount": "169"
			}
		},
		{
			"address": "comdex1g05p0vqqfpucng4gmutq6cpg626h7vwauy2u3l",
			"reward": {
				"denom": "ucmdx",
				"amount": "8923"
			}
		},
		{
			"address": "comdex1g04smgaf6s89kwy59ny9dqp4fhnrjul08jmtnz",
			"reward": {
				"denom": "ucmdx",
				"amount": "1424"
			}
		},
		{
			"address": "comdex1g0kpt93nu5th4raugwglprehapfyzsjksjw9ct",
			"reward": {
				"denom": "ucmdx",
				"amount": "281"
			}
		},
		{
			"address": "comdex1g0kzplgmuatet2p3f2w4z8f8y3skhshyasyhsa",
			"reward": {
				"denom": "ucmdx",
				"amount": "3534"
			}
		},
		{
			"address": "comdex1g0k6jwps8nlqlms97vka0p8qk8espaz9m8g4ns",
			"reward": {
				"denom": "ucmdx",
				"amount": "19529"
			}
		},
		{
			"address": "comdex1g0c9de4zsw2jsj9ys96r0qryuak05fvc2av5jf",
			"reward": {
				"denom": "ucmdx",
				"amount": "141261"
			}
		},
		{
			"address": "comdex1g0el0r82gglvc54sknauzqcp03nxsvzcas52xa",
			"reward": {
				"denom": "ucmdx",
				"amount": "180"
			}
		},
		{
			"address": "comdex1g0mn967fcq3gclrd6ffjfp2le5l89x58cr09dn",
			"reward": {
				"denom": "ucmdx",
				"amount": "8844"
			}
		},
		{
			"address": "comdex1g0mu9dj3m0s2u7m5nwt2z3jwaue37y6qjy3v52",
			"reward": {
				"denom": "ucmdx",
				"amount": "4971"
			}
		},
		{
			"address": "comdex1g0ap9fvex93yzzj69uzl74942vjvq8e4dms6rn",
			"reward": {
				"denom": "ucmdx",
				"amount": "600"
			}
		},
		{
			"address": "comdex1g072la5zq8zjl5z3j23qxrmj36u7q9a8fs4988",
			"reward": {
				"denom": "ucmdx",
				"amount": "2845"
			}
		},
		{
			"address": "comdex1gszxjrh55ycdat3eydug9nauhc77fh80f3qhed",
			"reward": {
				"denom": "ucmdx",
				"amount": "7935"
			}
		},
		{
			"address": "comdex1gszkaw734xpy9hn7xu9gc7kelzwtrf04rxufpg",
			"reward": {
				"denom": "ucmdx",
				"amount": "2882"
			}
		},
		{
			"address": "comdex1gs9szp3upnats03xtp7h29szzx0qcdz8qysu5n",
			"reward": {
				"denom": "ucmdx",
				"amount": "28"
			}
		},
		{
			"address": "comdex1gsxfx9d7vng654drfqe2cv9ku7ljk9p0w8jazk",
			"reward": {
				"denom": "ucmdx",
				"amount": "144"
			}
		},
		{
			"address": "comdex1gstgm2vp8mr66204xw5503pcyyps9qcxfkxr5u",
			"reward": {
				"denom": "ucmdx",
				"amount": "1994"
			}
		},
		{
			"address": "comdex1gst5ysqm8k9rvn88z7lr3qmvxn7jx37dku2e72",
			"reward": {
				"denom": "ucmdx",
				"amount": "13324"
			}
		},
		{
			"address": "comdex1gsdddlk569m46sxu7t838uj4q2w8v0vxxfwp79",
			"reward": {
				"denom": "ucmdx",
				"amount": "9914"
			}
		},
		{
			"address": "comdex1gs0khngjemyn9gn9lstuhhl6t54jkge2mdqf2m",
			"reward": {
				"denom": "ucmdx",
				"amount": "15077"
			}
		},
		{
			"address": "comdex1gsswrxjswvha7erz0e7pwv3ems7sdxd5t4angm",
			"reward": {
				"denom": "ucmdx",
				"amount": "173"
			}
		},
		{
			"address": "comdex1gssmktdt45mgaehncuuchcv089pc8p4jlxau9n",
			"reward": {
				"denom": "ucmdx",
				"amount": "4207"
			}
		},
		{
			"address": "comdex1gs38g7yh9c6srr7hyp9knln4jcm4ckcqyqzaqa",
			"reward": {
				"denom": "ucmdx",
				"amount": "7030"
			}
		},
		{
			"address": "comdex1gs3wkuqvmna5esv8j8qfk7lrav6l6qpxqu0873",
			"reward": {
				"denom": "ucmdx",
				"amount": "1920"
			}
		},
		{
			"address": "comdex1gs3s8ma9hhvhjnpk9q7xvwe7wmmnxe5nkknmz8",
			"reward": {
				"denom": "ucmdx",
				"amount": "5269"
			}
		},
		{
			"address": "comdex1gs3u9qnexx90qq73qfcg320k6nyg464x28xycc",
			"reward": {
				"denom": "ucmdx",
				"amount": "81163"
			}
		},
		{
			"address": "comdex1gs5tqrtv2ltpes98k64cfvhepj0vnvvq6e7352",
			"reward": {
				"denom": "ucmdx",
				"amount": "1054"
			}
		},
		{
			"address": "comdex1gs5wfmq09yku2vd5hchdmvy8v5k9n3n6c9n74m",
			"reward": {
				"denom": "ucmdx",
				"amount": "1764"
			}
		},
		{
			"address": "comdex1gs53kp8d25n6gq87wrh0wm2sku32a7st3c29qa",
			"reward": {
				"denom": "ucmdx",
				"amount": "21851"
			}
		},
		{
			"address": "comdex1gskch2j483gvn2rnvs84nfzknc3g3dpp8a2y0s",
			"reward": {
				"denom": "ucmdx",
				"amount": "3550"
			}
		},
		{
			"address": "comdex1gsmp9alrgwnwde3qzjwckkfk6ysss8j6vv2d5d",
			"reward": {
				"denom": "ucmdx",
				"amount": "180"
			}
		},
		{
			"address": "comdex1g3q827ug0k5g83tj67cldypczlp9w4cygpkye8",
			"reward": {
				"denom": "ucmdx",
				"amount": "1258"
			}
		},
		{
			"address": "comdex1g3rkshwrxhax2nq2lgq3n6tznktdumsu60fk60",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1g39mxhyz6mgq24xh2w2hw8l80vxf8tg0h6n6rp",
			"reward": {
				"denom": "ucmdx",
				"amount": "430"
			}
		},
		{
			"address": "comdex1g3xgtk4a44p95e5ugguy2cfatnzwrng5m7me6m",
			"reward": {
				"denom": "ucmdx",
				"amount": "17846"
			}
		},
		{
			"address": "comdex1g3v3whqk4n34d70qz5yz0c9tjqyr5wesjqaeaf",
			"reward": {
				"denom": "ucmdx",
				"amount": "1400"
			}
		},
		{
			"address": "comdex1g3dhll22talg0dekvyduangkpgrceqe0czwz5n",
			"reward": {
				"denom": "ucmdx",
				"amount": "6893"
			}
		},
		{
			"address": "comdex1g3wzl8ultshmtn7lmpun0hzy0llvyttc7wnrcj",
			"reward": {
				"denom": "ucmdx",
				"amount": "58450"
			}
		},
		{
			"address": "comdex1g334fkt5kdqlu27w3p4zskp8fwvgjvd66qzljt",
			"reward": {
				"denom": "ucmdx",
				"amount": "301"
			}
		},
		{
			"address": "comdex1g3jg90f4sw7q3e0djpqanly0f698q6t2zcj5l9",
			"reward": {
				"denom": "ucmdx",
				"amount": "8954"
			}
		},
		{
			"address": "comdex1g35f0y329ljejqvnx3mquwcsvdly9m7x0umnt2",
			"reward": {
				"denom": "ucmdx",
				"amount": "8639"
			}
		},
		{
			"address": "comdex1g3herkv4gkdnl664ddvhjew0zut74ap7d4qu56",
			"reward": {
				"denom": "ucmdx",
				"amount": "7146"
			}
		},
		{
			"address": "comdex1g3c500uq9f5t2aenwpsgwpeykmjj2alfjlr3s0",
			"reward": {
				"denom": "ucmdx",
				"amount": "44133"
			}
		},
		{
			"address": "comdex1g3un9pws0w5a0c98naq0csndlz0utrjqjen93n",
			"reward": {
				"denom": "ucmdx",
				"amount": "4911"
			}
		},
		{
			"address": "comdex1g3anhn5ek68d89l6kpfvtqnmf2texn9cy8k90t",
			"reward": {
				"denom": "ucmdx",
				"amount": "3295"
			}
		},
		{
			"address": "comdex1g3l5shnqy77kwehaznn28a5ycjc6ys0ygr2zy6",
			"reward": {
				"denom": "ucmdx",
				"amount": "1997"
			}
		},
		{
			"address": "comdex1gjzyn2qlha9xrl88g52rzeyd73v2twa7tht89t",
			"reward": {
				"denom": "ucmdx",
				"amount": "50089"
			}
		},
		{
			"address": "comdex1gjzjr6kkdfkqu5d35ku7a3q0mlza6m8ruwjn6s",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1gjrrd4c4euxp4fg59n3lq8jmna0qp04j2fgstp",
			"reward": {
				"denom": "ucmdx",
				"amount": "1419"
			}
		},
		{
			"address": "comdex1gjy00fv4faz90kwzv729z8ut7jkmzkswr4k0rf",
			"reward": {
				"denom": "ucmdx",
				"amount": "1190"
			}
		},
		{
			"address": "comdex1gjgjsh74w5trmfaaauq4qt2mwvh8p7gsd73g5f",
			"reward": {
				"denom": "ucmdx",
				"amount": "1432"
			}
		},
		{
			"address": "comdex1gjtznf7qkgw3wg6ldv07dpskge6nn2tdjl0qjf",
			"reward": {
				"denom": "ucmdx",
				"amount": "14167"
			}
		},
		{
			"address": "comdex1gj0jmm37wt6y52qdlnnzh9e7h23fvavylnwjs8",
			"reward": {
				"denom": "ucmdx",
				"amount": "29309"
			}
		},
		{
			"address": "comdex1gj39cf7mafc58aj7gj2vqdrjmu8sex9s5u8yuh",
			"reward": {
				"denom": "ucmdx",
				"amount": "725"
			}
		},
		{
			"address": "comdex1gj4nfkvnlxc60r95xhj7z4y97xzcjx7w9f3d54",
			"reward": {
				"denom": "ucmdx",
				"amount": "69585"
			}
		},
		{
			"address": "comdex1gjhc4ezy9x9zlzzs0ejk94vmx7mqxcyk5vrgdr",
			"reward": {
				"denom": "ucmdx",
				"amount": "38162"
			}
		},
		{
			"address": "comdex1gjh6atlq7xypygs7ga32r75kjpey09m9ucrcaj",
			"reward": {
				"denom": "ucmdx",
				"amount": "3560"
			}
		},
		{
			"address": "comdex1gjeefzf0gk8ttkctv945jg27gpupwar8z6nq82",
			"reward": {
				"denom": "ucmdx",
				"amount": "994"
			}
		},
		{
			"address": "comdex1gjakrtl8z5zq4k2upqg4jax022pa6yxvmp5q9p",
			"reward": {
				"denom": "ucmdx",
				"amount": "66600"
			}
		},
		{
			"address": "comdex1gj7tsashvzntr6zpulug7wtlrzp2k9dch88kjp",
			"reward": {
				"denom": "ucmdx",
				"amount": "494"
			}
		},
		{
			"address": "comdex1gnz30c42clfyh975y8ummegc7sdcnp8c7r8szu",
			"reward": {
				"denom": "ucmdx",
				"amount": "17773"
			}
		},
		{
			"address": "comdex1gnx6wfuzrh8tgu29dwml03fda76fc2juv9kygc",
			"reward": {
				"denom": "ucmdx",
				"amount": "30690"
			}
		},
		{
			"address": "comdex1gngva6d909dy7l5gzkthyswepawc0esdt226ha",
			"reward": {
				"denom": "ucmdx",
				"amount": "61186"
			}
		},
		{
			"address": "comdex1gnfqundmu8pzxh52g6ula38pdmghn6mmgnsny2",
			"reward": {
				"denom": "ucmdx",
				"amount": "19658"
			}
		},
		{
			"address": "comdex1gn2du82ucz2y60lxcvux0dzvmvgfkj4n3pjhf7",
			"reward": {
				"denom": "ucmdx",
				"amount": "91775"
			}
		},
		{
			"address": "comdex1gnset028l8x8v78z9njxjj8faj99f7exgyqtrx",
			"reward": {
				"denom": "ucmdx",
				"amount": "144"
			}
		},
		{
			"address": "comdex1gnsusr8lv4qvrrfgrnqz2lv7aud8z0wguw97uc",
			"reward": {
				"denom": "ucmdx",
				"amount": "685"
			}
		},
		{
			"address": "comdex1gn4pqx74udp45lcnxxkclud7r84hracg6eyk55",
			"reward": {
				"denom": "ucmdx",
				"amount": "3604"
			}
		},
		{
			"address": "comdex1gnh40zh065enk7y3pawlfnrg0nut5pgmyg5yv7",
			"reward": {
				"denom": "ucmdx",
				"amount": "288"
			}
		},
		{
			"address": "comdex1gnakutykpfuu5g2e0mlfdfpalua9e77cxqnrm9",
			"reward": {
				"denom": "ucmdx",
				"amount": "140546"
			}
		},
		{
			"address": "comdex1g5qhu97ex9m5feelprg6c745p4f8euqs26vuzn",
			"reward": {
				"denom": "ucmdx",
				"amount": "9964"
			}
		},
		{
			"address": "comdex1g5ysrg2ytk4m6t8ng2g2sh3ul2rd5d4nkqg292",
			"reward": {
				"denom": "ucmdx",
				"amount": "284"
			}
		},
		{
			"address": "comdex1g5xzrq3wnh09r9wh6agfhfzd00573gslvz5p6x",
			"reward": {
				"denom": "ucmdx",
				"amount": "1237"
			}
		},
		{
			"address": "comdex1g5tn2xq4m0ch0c7nxdr4usf08aqfnu2uvjucx5",
			"reward": {
				"denom": "ucmdx",
				"amount": "1878"
			}
		},
		{
			"address": "comdex1g5wrgeztmwnuak9dkxn5803pq6y8yt3q8etmts",
			"reward": {
				"denom": "ucmdx",
				"amount": "159205"
			}
		},
		{
			"address": "comdex1g50z7098eh8rzj08r9hr3gmj9g05gd66h6jfej",
			"reward": {
				"denom": "ucmdx",
				"amount": "34793"
			}
		},
		{
			"address": "comdex1g50gd0yd7v9ue0dfuh65tqatwmjq8mkw4p0v3c",
			"reward": {
				"denom": "ucmdx",
				"amount": "3353"
			}
		},
		{
			"address": "comdex1g53mgtymzz8d3naw0rwl0zrp3g7teyp4g98a92",
			"reward": {
				"denom": "ucmdx",
				"amount": "17703"
			}
		},
		{
			"address": "comdex1g55f83d5l4dc428c3l85x7a7f7vukkr7g4wzee",
			"reward": {
				"denom": "ucmdx",
				"amount": "1234"
			}
		},
		{
			"address": "comdex1g5hpt6hl6jaynjts8zwzkymsn2mwl7hh2hx902",
			"reward": {
				"denom": "ucmdx",
				"amount": "587"
			}
		},
		{
			"address": "comdex1g5h3qa7rm6wshacsxprza8ns09v8na7qda6uy6",
			"reward": {
				"denom": "ucmdx",
				"amount": "13048"
			}
		},
		{
			"address": "comdex1g5eme48375ace4f6xkqhy6djad6q3g4ksh86k4",
			"reward": {
				"denom": "ucmdx",
				"amount": "18734"
			}
		},
		{
			"address": "comdex1g5m4ehvz3pv34w8mm8kegqg0epf8hdz2zktzuf",
			"reward": {
				"denom": "ucmdx",
				"amount": "16"
			}
		},
		{
			"address": "comdex1g5urq7wqff2asddv2jwefadhhnsre2ssclhced",
			"reward": {
				"denom": "ucmdx",
				"amount": "7343"
			}
		},
		{
			"address": "comdex1g5u9puwmauuttq5juke55m4yjh40w94zrxgmfs",
			"reward": {
				"denom": "ucmdx",
				"amount": "126635"
			}
		},
		{
			"address": "comdex1g5ue7nrutnk0vh0xflahv0vt2mm03y2jctdmr7",
			"reward": {
				"denom": "ucmdx",
				"amount": "2729"
			}
		},
		{
			"address": "comdex1g4yq99n3ken7vmn9sxgh98cugw0cu2x3krd3m6",
			"reward": {
				"denom": "ucmdx",
				"amount": "1190"
			}
		},
		{
			"address": "comdex1g4y8rlyjhqz7msh2wspshrguuzzyq96j97jlmp",
			"reward": {
				"denom": "ucmdx",
				"amount": "28021"
			}
		},
		{
			"address": "comdex1g493uk9rr53enlrj85lljnds9g3k32sjny6zzt",
			"reward": {
				"denom": "ucmdx",
				"amount": "1627"
			}
		},
		{
			"address": "comdex1g4g9rfhvuyjnx9w0a3refv6y6v2k0lq53jsheh",
			"reward": {
				"denom": "ucmdx",
				"amount": "299"
			}
		},
		{
			"address": "comdex1g4gtaw3s6hhg352563sedrms2vummj5eklc3ue",
			"reward": {
				"denom": "ucmdx",
				"amount": "577388"
			}
		},
		{
			"address": "comdex1g4ftty8lf3uju0agqs4crdcdznq04t0j8803ld",
			"reward": {
				"denom": "ucmdx",
				"amount": "939"
			}
		},
		{
			"address": "comdex1g4f0du5aeetn22tlr3lzlzpqyn64w4a0q5qxkh",
			"reward": {
				"denom": "ucmdx",
				"amount": "598"
			}
		},
		{
			"address": "comdex1g4fcq507f2zrare449yllsyylukfysvfu3l002",
			"reward": {
				"denom": "ucmdx",
				"amount": "181"
			}
		},
		{
			"address": "comdex1g4wj2k82456ygult3e7efwtqfgxsajqem0nyg4",
			"reward": {
				"denom": "ucmdx",
				"amount": "527"
			}
		},
		{
			"address": "comdex1g40szkqzr63jca0h0usysdc02ugm0sscsvhg4c",
			"reward": {
				"denom": "ucmdx",
				"amount": "28165"
			}
		},
		{
			"address": "comdex1g43k0rj7hwwwghj9pzy4a2yrgkjqe9xn7x20qd",
			"reward": {
				"denom": "ucmdx",
				"amount": "7531"
			}
		},
		{
			"address": "comdex1g4nv8uhwcy8szfaseyvqd4xa38xv73z33rtpq0",
			"reward": {
				"denom": "ucmdx",
				"amount": "1190"
			}
		},
		{
			"address": "comdex1g44jqday548frwdskpx4gf3plth325ucdaghs3",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1g4lg403wat5mzeqr0u9h7034kch2n6cxc5u7q4",
			"reward": {
				"denom": "ucmdx",
				"amount": "1022"
			}
		},
		{
			"address": "comdex1gkqqcg8qhl2h0re4cmz3mqqnm8e6999jkdlgcd",
			"reward": {
				"denom": "ucmdx",
				"amount": "10001"
			}
		},
		{
			"address": "comdex1gkp85m23cag6s5cjk3s6sa4evr39nadkazft0c",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1gkp6aj8dasplgau9rdu6nk58pznlf9uqd63vd9",
			"reward": {
				"denom": "ucmdx",
				"amount": "138"
			}
		},
		{
			"address": "comdex1gkyaffl4t63lafgmquvygsanh9xh8d6jqukswa",
			"reward": {
				"denom": "ucmdx",
				"amount": "1998"
			}
		},
		{
			"address": "comdex1gkgwamswhd842truw5ch6tlvkmu77dsg903ue4",
			"reward": {
				"denom": "ucmdx",
				"amount": "56699"
			}
		},
		{
			"address": "comdex1gkgndes9rtzpc3shgtyqr62vm99dnxquzl6m3d",
			"reward": {
				"denom": "ucmdx",
				"amount": "11212"
			}
		},
		{
			"address": "comdex1gkfqttdzzptu7f9l73t84nyy34jzr6m5t7f496",
			"reward": {
				"denom": "ucmdx",
				"amount": "37545"
			}
		},
		{
			"address": "comdex1gkf0h8akjwplvda024af5l9dzg06tpuq939wrd",
			"reward": {
				"denom": "ucmdx",
				"amount": "3040"
			}
		},
		{
			"address": "comdex1gkfm69ulp493c4srhet0gm8pgpzm7j7u5a99tz",
			"reward": {
				"denom": "ucmdx",
				"amount": "18743"
			}
		},
		{
			"address": "comdex1gk2x0tldm79y46ps5p3uydgm0kyac9pgdrcep9",
			"reward": {
				"denom": "ucmdx",
				"amount": "266272"
			}
		},
		{
			"address": "comdex1gk20pskfuc04yaw2wpgxv2d7kmq52ktfrm5jr3",
			"reward": {
				"denom": "ucmdx",
				"amount": "1168"
			}
		},
		{
			"address": "comdex1gktfpxakn32qwaal8q6606syql9tr3svhj95ul",
			"reward": {
				"denom": "ucmdx",
				"amount": "6045"
			}
		},
		{
			"address": "comdex1gkjx3l4n84uu2hpl6sqv6vfrwxr8qvsvemjlkv",
			"reward": {
				"denom": "ucmdx",
				"amount": "11840"
			}
		},
		{
			"address": "comdex1gkjw5s90nh7mxxa35x96awszz4kqfna9c346d5",
			"reward": {
				"denom": "ucmdx",
				"amount": "185"
			}
		},
		{
			"address": "comdex1gkn9gl58m3zwx32ftle0tm9xzgmzumy999qynd",
			"reward": {
				"denom": "ucmdx",
				"amount": "173"
			}
		},
		{
			"address": "comdex1gk5yez27hkhnrpfs7ych4jyhfst4h0pzjjpp9h",
			"reward": {
				"denom": "ucmdx",
				"amount": "7951"
			}
		},
		{
			"address": "comdex1gk4ss97lzs672wwe72eh6y6jrh56r8lk37zsxu",
			"reward": {
				"denom": "ucmdx",
				"amount": "178"
			}
		},
		{
			"address": "comdex1gke5qyd0uv5yn8kg4ev5fvv9alt42swdssqug7",
			"reward": {
				"denom": "ucmdx",
				"amount": "138490"
			}
		},
		{
			"address": "comdex1gkekgzqseujvlunqk3jx77k6wd3re48v7n8qec",
			"reward": {
				"denom": "ucmdx",
				"amount": "22167"
			}
		},
		{
			"address": "comdex1gk6h6lswwqer8ulfgfvy64nj0p42xd9d6jme04",
			"reward": {
				"denom": "ucmdx",
				"amount": "1753"
			}
		},
		{
			"address": "comdex1gkmha9lm4zha2h9jv2xpzy3wjz7vc0qqsffne9",
			"reward": {
				"denom": "ucmdx",
				"amount": "28"
			}
		},
		{
			"address": "comdex1gkats6cusuernkf34gkwudvsf2havv0d68u7sd",
			"reward": {
				"denom": "ucmdx",
				"amount": "71"
			}
		},
		{
			"address": "comdex1ghqg52xuc3qr5xx06p3fkrkhrvt3ajv2rxgdp5",
			"reward": {
				"denom": "ucmdx",
				"amount": "538"
			}
		},
		{
			"address": "comdex1ghfnm8fe0e6tzp50940qsxkzqsdeqpv47serk9",
			"reward": {
				"denom": "ucmdx",
				"amount": "4956"
			}
		},
		{
			"address": "comdex1ght427rvazvpnm29jzayvmqqvxeef96tj4uhd7",
			"reward": {
				"denom": "ucmdx",
				"amount": "89"
			}
		},
		{
			"address": "comdex1ghvx8xc2krec4mvdkana5rg073qdg30v9wajm6",
			"reward": {
				"denom": "ucmdx",
				"amount": "6784"
			}
		},
		{
			"address": "comdex1ghvdfm50yds2nsjqcy7en884d0403mnu9l24mf",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1ghw2gq6eawg007g2pc3fmqwgktm2jvcyuatypn",
			"reward": {
				"denom": "ucmdx",
				"amount": "1283"
			}
		},
		{
			"address": "comdex1gh3pfrqy4sga7dy9rk3xsmc04u8xyepfvq20kw",
			"reward": {
				"denom": "ucmdx",
				"amount": "6370"
			}
		},
		{
			"address": "comdex1gh33e7j5tcf9zdfcay62f4sy2k8tsn85pjk60k",
			"reward": {
				"denom": "ucmdx",
				"amount": "57"
			}
		},
		{
			"address": "comdex1ghndpjxgnuzzpaxct9c9szwugr2505xqqjvfe6",
			"reward": {
				"denom": "ucmdx",
				"amount": "1406"
			}
		},
		{
			"address": "comdex1ghnhj4ahjakacevgfrgdnx9mjvdxxv5ac8akgv",
			"reward": {
				"denom": "ucmdx",
				"amount": "19211"
			}
		},
		{
			"address": "comdex1ghkc00svrwra2h88wx0q7623lwnja9s6q6n0zs",
			"reward": {
				"denom": "ucmdx",
				"amount": "1091"
			}
		},
		{
			"address": "comdex1ghmte3r8xyqfngw8cpzgtgzzmr7wjxrx4njk96",
			"reward": {
				"denom": "ucmdx",
				"amount": "165"
			}
		},
		{
			"address": "comdex1ghmcern54qv2qcn5lj6v0ss06jha606sux9uu5",
			"reward": {
				"denom": "ucmdx",
				"amount": "13"
			}
		},
		{
			"address": "comdex1gh756ppj2h6guymek06xkulhmahk4tp4f9wlft",
			"reward": {
				"denom": "ucmdx",
				"amount": "58"
			}
		},
		{
			"address": "comdex1ghlgft9awa0ke2ye8ky4mad0r5vqk0fgda0622",
			"reward": {
				"denom": "ucmdx",
				"amount": "14170"
			}
		},
		{
			"address": "comdex1gcq2rks9qmdxh2vvtkcccwtt4pdq8w0lp5xy68",
			"reward": {
				"denom": "ucmdx",
				"amount": "1410"
			}
		},
		{
			"address": "comdex1gcz8rr2tqrhuhf85f9n3pw32ak7rd2aywycq86",
			"reward": {
				"denom": "ucmdx",
				"amount": "3756"
			}
		},
		{
			"address": "comdex1gcr8dfp0qurxuddvvy98jw9gn5ku7005d2k3qr",
			"reward": {
				"denom": "ucmdx",
				"amount": "11808"
			}
		},
		{
			"address": "comdex1gcxkh4z0xhg447ng39u7lfsfhg6cl946mfgrvl",
			"reward": {
				"denom": "ucmdx",
				"amount": "3549"
			}
		},
		{
			"address": "comdex1gcfntcpck5wg23g6xngqr6aavcrhq54vlkp9gw",
			"reward": {
				"denom": "ucmdx",
				"amount": "204"
			}
		},
		{
			"address": "comdex1gcty5saafc957fdx582adzu5k3zsw5qga9x9tf",
			"reward": {
				"denom": "ucmdx",
				"amount": "14530"
			}
		},
		{
			"address": "comdex1gcv5k9cszdz4lkx4kzvsvz5ymnnyeu622yl5f4",
			"reward": {
				"denom": "ucmdx",
				"amount": "302"
			}
		},
		{
			"address": "comdex1gcd3uyeg5xxu95edsywv6rwm8pkwt7xn8757um",
			"reward": {
				"denom": "ucmdx",
				"amount": "1190"
			}
		},
		{
			"address": "comdex1gc0509w3ee4g2fm7w80nxc05zrld695hjh7csr",
			"reward": {
				"denom": "ucmdx",
				"amount": "21714"
			}
		},
		{
			"address": "comdex1gcsfr6ywan84nnsve78nx4ajzjw6qqqlc9zuak",
			"reward": {
				"denom": "ucmdx",
				"amount": "17575"
			}
		},
		{
			"address": "comdex1gc323lmewm3gsxq8rm9swe6crppyx30sh73fp2",
			"reward": {
				"denom": "ucmdx",
				"amount": "12315"
			}
		},
		{
			"address": "comdex1gcjtxsa5dj9vxtnmf3ld6jg63av9gkj7zf3lyd",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1gcj6dduz4nde39xc28v3ftrwujxje3nan34lcy",
			"reward": {
				"denom": "ucmdx",
				"amount": "39038"
			}
		},
		{
			"address": "comdex1gc55cr8hyeu7wffwq0fw4ta0xpf39w0xcuu58t",
			"reward": {
				"denom": "ucmdx",
				"amount": "15482"
			}
		},
		{
			"address": "comdex1gc4vc76uzta0uar2kw2kxygcqpp86nuefwwhwp",
			"reward": {
				"denom": "ucmdx",
				"amount": "2633"
			}
		},
		{
			"address": "comdex1gckju3fmw9lfrfsdky8hrz9kylg2a7hlyfrdt6",
			"reward": {
				"denom": "ucmdx",
				"amount": "1887"
			}
		},
		{
			"address": "comdex1gce08edkdalhnr8yc0h3754edz7x082yn3qk7v",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1gcm0c04t7z7y43awch2tyj4629na3avz0d00fc",
			"reward": {
				"denom": "ucmdx",
				"amount": "1115"
			}
		},
		{
			"address": "comdex1gcafwtsrr3596x7jcslfm7376en8p59u0d49vk",
			"reward": {
				"denom": "ucmdx",
				"amount": "88"
			}
		},
		{
			"address": "comdex1gcatcjwteap6zzmu4ryq2g235jmrncltcwyu02",
			"reward": {
				"denom": "ucmdx",
				"amount": "1466"
			}
		},
		{
			"address": "comdex1gc7hznznsm5t2mldwwlsg02uv6r3z2makhfu4m",
			"reward": {
				"denom": "ucmdx",
				"amount": "203"
			}
		},
		{
			"address": "comdex1gclzhncr3rh97z739zt76kds3ugk7flg347z08",
			"reward": {
				"denom": "ucmdx",
				"amount": "19"
			}
		},
		{
			"address": "comdex1ger0r97emr39eh0kz03afgnqj7t38xth5dmtkr",
			"reward": {
				"denom": "ucmdx",
				"amount": "8704"
			}
		},
		{
			"address": "comdex1ger4sp8lvc2ena7f9xl7nk987ynsy8dpsdmc8x",
			"reward": {
				"denom": "ucmdx",
				"amount": "14051"
			}
		},
		{
			"address": "comdex1ge9a7mxqm50mnqgl0mu7ghd6rlx4sfql535xfn",
			"reward": {
				"denom": "ucmdx",
				"amount": "276"
			}
		},
		{
			"address": "comdex1ge8ue02qdtajllu8lmtx3jhht2fe497y4npvuv",
			"reward": {
				"denom": "ucmdx",
				"amount": "2049"
			}
		},
		{
			"address": "comdex1gegd6ufdzp0xw2jlrvehnyk7mac0mtewqy9a0m",
			"reward": {
				"denom": "ucmdx",
				"amount": "133"
			}
		},
		{
			"address": "comdex1gedynzufda2cwarehjs7dhv8pujr3nlpwnwsq9",
			"reward": {
				"denom": "ucmdx",
				"amount": "3982"
			}
		},
		{
			"address": "comdex1gewdx3f0s40hcss4ygchumqmtark2l80gwhedv",
			"reward": {
				"denom": "ucmdx",
				"amount": "61"
			}
		},
		{
			"address": "comdex1gewl0p63x9r8lunjvwtvjq26qy8q9x7vw3e3s8",
			"reward": {
				"denom": "ucmdx",
				"amount": "1439"
			}
		},
		{
			"address": "comdex1ge0xdghslkg309zka9qc9qrz2gehjw6sx3umnq",
			"reward": {
				"denom": "ucmdx",
				"amount": "414"
			}
		},
		{
			"address": "comdex1ge0f5v38rz6p296t6csxy7hv2nw2fgy0nghzzm",
			"reward": {
				"denom": "ucmdx",
				"amount": "443"
			}
		},
		{
			"address": "comdex1gesrz7ydeluyr3c75du2qxp76peq02a9musrcr",
			"reward": {
				"denom": "ucmdx",
				"amount": "12672"
			}
		},
		{
			"address": "comdex1gest3cqved5en6ulyqw5pj7tluty5rl3v0d7k7",
			"reward": {
				"denom": "ucmdx",
				"amount": "72585"
			}
		},
		{
			"address": "comdex1ge39gkjtvwnh9za8266md0pp306yptn8qhx2ds",
			"reward": {
				"denom": "ucmdx",
				"amount": "8699"
			}
		},
		{
			"address": "comdex1gecvewsuktd6v2djs08j3vvt4z0myeqahqxm5n",
			"reward": {
				"denom": "ucmdx",
				"amount": "357"
			}
		},
		{
			"address": "comdex1gem0vjkg6zg0v3t4elnfyame6n2027cyy97q33",
			"reward": {
				"denom": "ucmdx",
				"amount": "204"
			}
		},
		{
			"address": "comdex1gemss7jcmvpjdnkqkhfp56jgwvfe6t88ylkmk8",
			"reward": {
				"denom": "ucmdx",
				"amount": "28366"
			}
		},
		{
			"address": "comdex1gemekeu337rg38l964drz8ncllg363p90lg6km",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1g6psvrjfns0n9k6m3v3hs6sl05eurjlnjmq4uh",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1g6zszf3qluf4g22fhenvz02ad43kaz5ecuy8e4",
			"reward": {
				"denom": "ucmdx",
				"amount": "3810"
			}
		},
		{
			"address": "comdex1g6t0n77l4lygs8xavdlj82wy4ntpg2jtelu8k7",
			"reward": {
				"denom": "ucmdx",
				"amount": "1399"
			}
		},
		{
			"address": "comdex1g6t76w23v5vzpsyvtxx8slpfhmsuemw4fvcezl",
			"reward": {
				"denom": "ucmdx",
				"amount": "71294"
			}
		},
		{
			"address": "comdex1g6d349vkh6lqtx8z29rnredqe6pfmsn0f89lgx",
			"reward": {
				"denom": "ucmdx",
				"amount": "5835"
			}
		},
		{
			"address": "comdex1g6wtcxry6dmw6ucdmjezvqja5f5yrukcjqdhhh",
			"reward": {
				"denom": "ucmdx",
				"amount": "181"
			}
		},
		{
			"address": "comdex1g6ww3963a4z6vg2npzm29luv084dx0graphlek",
			"reward": {
				"denom": "ucmdx",
				"amount": "1765"
			}
		},
		{
			"address": "comdex1g60teezwmfdj8xxpnd5kehvp25zfzt25hr4n5k",
			"reward": {
				"denom": "ucmdx",
				"amount": "1530"
			}
		},
		{
			"address": "comdex1g603l6h7ly0z2fnzhuzynhhrcftjg5xhq8hgch",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1g6sctg4nyqrq0czx7mtdkdkwx65lc4g03lynav",
			"reward": {
				"denom": "ucmdx",
				"amount": "12598"
			}
		},
		{
			"address": "comdex1g6nyswzhnlw26ygqe42k9k5ffymywzy6ganjqh",
			"reward": {
				"denom": "ucmdx",
				"amount": "442"
			}
		},
		{
			"address": "comdex1g6nn3agajvs772lp780z7hvdecrdwmq75wq7fe",
			"reward": {
				"denom": "ucmdx",
				"amount": "275"
			}
		},
		{
			"address": "comdex1g6neus5gewktat578478p0c5rh4mcqwf9jjwd5",
			"reward": {
				"denom": "ucmdx",
				"amount": "3842"
			}
		},
		{
			"address": "comdex1g6kmdundgtnzpjv48lt4xsj9tvvs2ga0qaggjt",
			"reward": {
				"denom": "ucmdx",
				"amount": "3650"
			}
		},
		{
			"address": "comdex1g6hp8n8mefwmg8x3fhg4zc33d68rgfp8jlqww3",
			"reward": {
				"denom": "ucmdx",
				"amount": "266194"
			}
		},
		{
			"address": "comdex1g6ce6dd8n7k7ej958rvjjtzewpgrkmqmg46mts",
			"reward": {
				"denom": "ucmdx",
				"amount": "1885"
			}
		},
		{
			"address": "comdex1g6euuelqhdfqkqsvp9kx343t5nnax8jeqc2223",
			"reward": {
				"denom": "ucmdx",
				"amount": "3563"
			}
		},
		{
			"address": "comdex1g6mecapy68aguql2eadpszdap8yggsk6gplvrp",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1gmq0zg0jdnqagr3tkmyxc7gl457rdry0mwmzdw",
			"reward": {
				"denom": "ucmdx",
				"amount": "30816"
			}
		},
		{
			"address": "comdex1gmp3gw6gezzufrg83dmapjcemghawv778adwgy",
			"reward": {
				"denom": "ucmdx",
				"amount": "729"
			}
		},
		{
			"address": "comdex1gmpjplqdrszhwqndklwk7y7xf42s8ssevvyn82",
			"reward": {
				"denom": "ucmdx",
				"amount": "14"
			}
		},
		{
			"address": "comdex1gmz8ulswu8ezgt5qrylzqzlkre08j7gqmhlnrd",
			"reward": {
				"denom": "ucmdx",
				"amount": "339"
			}
		},
		{
			"address": "comdex1gm9e0kyyz35gvz9t8eedldarx4k0pnxks35u5m",
			"reward": {
				"denom": "ucmdx",
				"amount": "42257"
			}
		},
		{
			"address": "comdex1gmgauupu4gs3u8fd75sjjauadcymg84nvvtgz4",
			"reward": {
				"denom": "ucmdx",
				"amount": "1506"
			}
		},
		{
			"address": "comdex1gmflv0w82udkfaw3vkyu7j3v4jdckvd4r6ujum",
			"reward": {
				"denom": "ucmdx",
				"amount": "16"
			}
		},
		{
			"address": "comdex1gmvdcu4fga4kvqdm822pfux7fqdcagxxqzll4p",
			"reward": {
				"denom": "ucmdx",
				"amount": "14878"
			}
		},
		{
			"address": "comdex1gmjyf9t8z8p0mhp26fxej67l56rgvgxjen8575",
			"reward": {
				"denom": "ucmdx",
				"amount": "14231"
			}
		},
		{
			"address": "comdex1gmnw7lw5lfrzdsz07hzhvx9x8y79sddwce9qf2",
			"reward": {
				"denom": "ucmdx",
				"amount": "167"
			}
		},
		{
			"address": "comdex1gm5wv6tupd0fyartgld58gm569y3ydjvtytast",
			"reward": {
				"denom": "ucmdx",
				"amount": "14234"
			}
		},
		{
			"address": "comdex1gmh9fyu3s5dtrj0sunkr3uxwadww636wlra9e5",
			"reward": {
				"denom": "ucmdx",
				"amount": "37024"
			}
		},
		{
			"address": "comdex1gmhvppz0nsckw4n2e3q0q2nn9ylvz0jt6m5lf6",
			"reward": {
				"denom": "ucmdx",
				"amount": "77"
			}
		},
		{
			"address": "comdex1gmeegfecqhmc5ysvpfhq5u9840e0wmg45ha8ld",
			"reward": {
				"denom": "ucmdx",
				"amount": "2876"
			}
		},
		{
			"address": "comdex1gmml7u4jfmzd90rhwcsw0flzxsw44scueftu9h",
			"reward": {
				"denom": "ucmdx",
				"amount": "1873"
			}
		},
		{
			"address": "comdex1gmuuzwt4sqlzk9vppchlw73fu3akj2s3gp8jfw",
			"reward": {
				"denom": "ucmdx",
				"amount": "17675"
			}
		},
		{
			"address": "comdex1gmam4ycuccmfevaftzxc5qg5kuwetfkacp08ne",
			"reward": {
				"denom": "ucmdx",
				"amount": "171"
			}
		},
		{
			"address": "comdex1gmldc4msum7fer0mp3dq4u9k3wnemw8qazthwm",
			"reward": {
				"denom": "ucmdx",
				"amount": "358"
			}
		},
		{
			"address": "comdex1gmlnh5kxfl6sylj8lzs6w240d2vun7gc8ddx73",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1guqxh6jrshhvdvgen2zjf9dmet8drh6ay52j27",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1guztpm035q9t3ggyk3p6nzd5557lnk4f9pq70j",
			"reward": {
				"denom": "ucmdx",
				"amount": "22205"
			}
		},
		{
			"address": "comdex1guys8wjhhx4ptfm48e8t8ssx2s8pjalqpfp5px",
			"reward": {
				"denom": "ucmdx",
				"amount": "7065"
			}
		},
		{
			"address": "comdex1guy3nmq4zawfvnzvztajhxt9m2ka0g28lcvjcg",
			"reward": {
				"denom": "ucmdx",
				"amount": "290"
			}
		},
		{
			"address": "comdex1guy40596xzhps4wqlddvmnqv80yml6n72kzjj3",
			"reward": {
				"denom": "ucmdx",
				"amount": "665"
			}
		},
		{
			"address": "comdex1guyhw5fyfhpey83px33wfajuryr7sysxec5u5q",
			"reward": {
				"denom": "ucmdx",
				"amount": "71699"
			}
		},
		{
			"address": "comdex1gu92w3wkwawq6q8epdsyce6x7q5gzvv59nwv6a",
			"reward": {
				"denom": "ucmdx",
				"amount": "49418"
			}
		},
		{
			"address": "comdex1guger5vlqy5r8yhy29ay62aqspslex5erjz84r",
			"reward": {
				"denom": "ucmdx",
				"amount": "820"
			}
		},
		{
			"address": "comdex1gu2kkx2ap4p8c4wpmswq2yxns4t0phhyyn7r8r",
			"reward": {
				"denom": "ucmdx",
				"amount": "166"
			}
		},
		{
			"address": "comdex1gu263hytgnh4psxh2erkvjrrqhlr9uzpja8zzv",
			"reward": {
				"denom": "ucmdx",
				"amount": "82067"
			}
		},
		{
			"address": "comdex1gutppfxgmwcrm4ws796ma467reu4cj8q8mc58k",
			"reward": {
				"denom": "ucmdx",
				"amount": "9240"
			}
		},
		{
			"address": "comdex1gustqeccxqzwvx5096cdsx6087tt7cavm7m5aj",
			"reward": {
				"denom": "ucmdx",
				"amount": "3151"
			}
		},
		{
			"address": "comdex1gu39lh8w2x6aee4t26xm7ez85hmcujzcffwqus",
			"reward": {
				"denom": "ucmdx",
				"amount": "12627"
			}
		},
		{
			"address": "comdex1gu3dkrngqp4mtgpk5a8s3h4yx5uptycw225lj0",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1gujfrkacufgczyayflw84vc0q3gapldkpr6wvc",
			"reward": {
				"denom": "ucmdx",
				"amount": "174"
			}
		},
		{
			"address": "comdex1gu5f6xgkswnlerefk4s86txdsxsety9yk30977",
			"reward": {
				"denom": "ucmdx",
				"amount": "704"
			}
		},
		{
			"address": "comdex1guujdx9samvwwekz7dcrq77e4ggpjkdv54vxha",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1guagrm7cz9ppjt3mjxtan757ur7z7pmzwhy27p",
			"reward": {
				"denom": "ucmdx",
				"amount": "22785"
			}
		},
		{
			"address": "comdex1guagvmgph2kvtvakqk8d96x05x2aung96md5e3",
			"reward": {
				"denom": "ucmdx",
				"amount": "10321"
			}
		},
		{
			"address": "comdex1gu7ludkpuaaul3mkztj4fs2drgpxatpp6dalgy",
			"reward": {
				"denom": "ucmdx",
				"amount": "18200"
			}
		},
		{
			"address": "comdex1gulgrgax2akt7t6jwgf74axwgrlyp07e9e3l9e",
			"reward": {
				"denom": "ucmdx",
				"amount": "1686"
			}
		},
		{
			"address": "comdex1gulc3yp4k4nruux8dwtd6q5h0x5wj44hy39ndw",
			"reward": {
				"denom": "ucmdx",
				"amount": "119"
			}
		},
		{
			"address": "comdex1gaqgs8d823ls6dqsy9g3am6cvhhr8clh4fajnf",
			"reward": {
				"denom": "ucmdx",
				"amount": "369"
			}
		},
		{
			"address": "comdex1gapduxxmcv6q54xehqq3mq6qs3jfrgsxy8d5n8",
			"reward": {
				"denom": "ucmdx",
				"amount": "173"
			}
		},
		{
			"address": "comdex1gapau25gzcjshra6yzd478af9aqeh6a37eclxa",
			"reward": {
				"denom": "ucmdx",
				"amount": "6760"
			}
		},
		{
			"address": "comdex1gazfeg04x85xj58q9wherpl0kswad8thy8yrvm",
			"reward": {
				"denom": "ucmdx",
				"amount": "7128"
			}
		},
		{
			"address": "comdex1gayj3vhg0tud5z4rpujr2lwvglkqzgzhwe5yne",
			"reward": {
				"denom": "ucmdx",
				"amount": "1763"
			}
		},
		{
			"address": "comdex1gax90ly6gezt9yffxnhplux0xhtzgt4v0rmgz7",
			"reward": {
				"denom": "ucmdx",
				"amount": "143"
			}
		},
		{
			"address": "comdex1ga8tq2jv9cks6scpmrjfhst55dqvtmhqv3r235",
			"reward": {
				"denom": "ucmdx",
				"amount": "1981"
			}
		},
		{
			"address": "comdex1ga8nv7ng4955t2y4anxwrhp7rkmj42l6t2ylcx",
			"reward": {
				"denom": "ucmdx",
				"amount": "2841"
			}
		},
		{
			"address": "comdex1gagvd7xj76eflkn2lzxnckv0swxjgfk9vt23dc",
			"reward": {
				"denom": "ucmdx",
				"amount": "20364"
			}
		},
		{
			"address": "comdex1gags48w3gxhzlcln35l9jqxlzypzp7mjr67gkr",
			"reward": {
				"denom": "ucmdx",
				"amount": "1756"
			}
		},
		{
			"address": "comdex1ga2qvmz6zlckc9wxfnmqelrvpcrdpp8lgxpr68",
			"reward": {
				"denom": "ucmdx",
				"amount": "20515"
			}
		},
		{
			"address": "comdex1gatzuvz53pzyt6557ax56nysfz6tqvrx5e4ska",
			"reward": {
				"denom": "ucmdx",
				"amount": "70"
			}
		},
		{
			"address": "comdex1gav4004znpaja3dj25lujp4ujjlf2fmlkuwczc",
			"reward": {
				"denom": "ucmdx",
				"amount": "10002"
			}
		},
		{
			"address": "comdex1gaw9yhkra063c8gnzujgyfkjjldfkvu69dyw8a",
			"reward": {
				"denom": "ucmdx",
				"amount": "119691"
			}
		},
		{
			"address": "comdex1gaw0tl76jvs0gs5lyc2frjktcm07ewfmguqzv0",
			"reward": {
				"denom": "ucmdx",
				"amount": "148"
			}
		},
		{
			"address": "comdex1gasyx4hrk8vnsfrkvzfck79hxggta5dz6mddcn",
			"reward": {
				"denom": "ucmdx",
				"amount": "608118"
			}
		},
		{
			"address": "comdex1ganvmzhxjt60c2udxvzs2xle52khzasq0r0quh",
			"reward": {
				"denom": "ucmdx",
				"amount": "26792"
			}
		},
		{
			"address": "comdex1gan3mx6x5pzqsuttr6m4vsp4znf8fxe7gstq36",
			"reward": {
				"denom": "ucmdx",
				"amount": "14"
			}
		},
		{
			"address": "comdex1gakrtr8p4jvlgpgwsex5sfuphnmjnw365nr2pm",
			"reward": {
				"denom": "ucmdx",
				"amount": "1442"
			}
		},
		{
			"address": "comdex1gakeuauc2akv9rs0c3ezws3qa46hhuj4nu5n7a",
			"reward": {
				"denom": "ucmdx",
				"amount": "85325"
			}
		},
		{
			"address": "comdex1gac5q72ef24z49pfmvwq25pjjaujq6zxqpgwu0",
			"reward": {
				"denom": "ucmdx",
				"amount": "14886"
			}
		},
		{
			"address": "comdex1gac5fdf7kz5ke9vtctzujh3cuxhgt6njsz6u78",
			"reward": {
				"denom": "ucmdx",
				"amount": "1560"
			}
		},
		{
			"address": "comdex1gaatuaed2ljv6afm5kxuxcuefyh7t7ez804ltm",
			"reward": {
				"denom": "ucmdx",
				"amount": "166"
			}
		},
		{
			"address": "comdex1gaa4x38k3ppwnxz2zg0cmxnhausj7udtcqx3cs",
			"reward": {
				"denom": "ucmdx",
				"amount": "1384"
			}
		},
		{
			"address": "comdex1ga7nyuhlgxvtd7aryqtf8mhl3udep4d4ur9c5x",
			"reward": {
				"denom": "ucmdx",
				"amount": "7640"
			}
		},
		{
			"address": "comdex1ga7nl9xzetu6mm6ykns0muhrk7l7meazc7nlqy",
			"reward": {
				"denom": "ucmdx",
				"amount": "1981"
			}
		},
		{
			"address": "comdex1ga75k6vs39t2ln4larxq5yyd0jv84u29m2j0v3",
			"reward": {
				"denom": "ucmdx",
				"amount": "419"
			}
		},
		{
			"address": "comdex1gal22l4farrzhh4r2ux33dc6ks2c8gyjlvm6vy",
			"reward": {
				"denom": "ucmdx",
				"amount": "11525"
			}
		},
		{
			"address": "comdex1g7qpzz8pxaknw0w8xe9zuclwtx6yhlq6ghr20q",
			"reward": {
				"denom": "ucmdx",
				"amount": "1427"
			}
		},
		{
			"address": "comdex1g7q8z9dx5vqz6kxlgmxlua5unj7t6uuwpal26m",
			"reward": {
				"denom": "ucmdx",
				"amount": "2160"
			}
		},
		{
			"address": "comdex1g795397pzayqx6gtml396z5rfk4fqg839pp8lz",
			"reward": {
				"denom": "ucmdx",
				"amount": "12626"
			}
		},
		{
			"address": "comdex1g7gqgwyd42p2v9y8zm7cxxth6gtzmhehn7uzme",
			"reward": {
				"denom": "ucmdx",
				"amount": "124"
			}
		},
		{
			"address": "comdex1g7gf0azwx0wp39rc0gszdcgdqe9tnlhhy9aphh",
			"reward": {
				"denom": "ucmdx",
				"amount": "13497"
			}
		},
		{
			"address": "comdex1g7fpvuz2q0y5t6lch9d6w56fx86cyqm0xwjexz",
			"reward": {
				"denom": "ucmdx",
				"amount": "1762"
			}
		},
		{
			"address": "comdex1g72j0zvhww3hra7m3v9vnmjdudjzw43ugyuc2d",
			"reward": {
				"denom": "ucmdx",
				"amount": "4097"
			}
		},
		{
			"address": "comdex1g72hk7c5sc8rq8qk6c7jaldfr8zcdv27w3tvap",
			"reward": {
				"denom": "ucmdx",
				"amount": "615386"
			}
		},
		{
			"address": "comdex1g7vf8e00y69tf5ej3zygtn2367qfzcuhlcnpr8",
			"reward": {
				"denom": "ucmdx",
				"amount": "142"
			}
		},
		{
			"address": "comdex1g7shdh7exxrfq8av8aums83sgln8h5wfh0vetg",
			"reward": {
				"denom": "ucmdx",
				"amount": "8440"
			}
		},
		{
			"address": "comdex1g73xcmyfmn5tp6632l2t0fmm4a6fd8nhsm3kdk",
			"reward": {
				"denom": "ucmdx",
				"amount": "128837"
			}
		},
		{
			"address": "comdex1g7nwuxwyl6claxknsxjg79a3ylpsu5m2v4fagl",
			"reward": {
				"denom": "ucmdx",
				"amount": "1252"
			}
		},
		{
			"address": "comdex1g75xt0hd4sx7fhf8ewpvnvhavq3yun6wrvs8zn",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1g75e76d2dj3l8p5gxk27flgd9p2v2ewpwgeh9h",
			"reward": {
				"denom": "ucmdx",
				"amount": "3017"
			}
		},
		{
			"address": "comdex1g7knemsy0rs9cy0uagu9ar3azqzqucdkdt3nqa",
			"reward": {
				"denom": "ucmdx",
				"amount": "123"
			}
		},
		{
			"address": "comdex1g7c7vstu060hjngsgmj20fnkveweejwukyk79s",
			"reward": {
				"denom": "ucmdx",
				"amount": "697"
			}
		},
		{
			"address": "comdex1g7ms2jjncwzep5fr4e0q0yksnazs6xxapt3r9a",
			"reward": {
				"denom": "ucmdx",
				"amount": "151"
			}
		},
		{
			"address": "comdex1glqjyvjda2jc8xkgcgfqeldegejw4crwnsutfc",
			"reward": {
				"denom": "ucmdx",
				"amount": "91868"
			}
		},
		{
			"address": "comdex1glpzwnvh8eug6d3c0g0kepp9af90a3gyw6aa3m",
			"reward": {
				"denom": "ucmdx",
				"amount": "697"
			}
		},
		{
			"address": "comdex1glw68t9raq4xd7hgad5zl2r2ndvuwtukhdk9ac",
			"reward": {
				"denom": "ucmdx",
				"amount": "405"
			}
		},
		{
			"address": "comdex1glsz0mk4cy0qppmlpgts6r8z0slw7xt8teycey",
			"reward": {
				"denom": "ucmdx",
				"amount": "1250"
			}
		},
		{
			"address": "comdex1glnzp427yskmxk8kurxten52mvgeqpwdze7qrn",
			"reward": {
				"denom": "ucmdx",
				"amount": "291"
			}
		},
		{
			"address": "comdex1glnh90y92axym9x7p7cl4n7ltd4d77jp0fqmtr",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1gl4vuf5t5h2azj7khrvs8m0gfz0yh49zlxs9qh",
			"reward": {
				"denom": "ucmdx",
				"amount": "950"
			}
		},
		{
			"address": "comdex1glc0l8mfdh93kjctk8g56e5uqnkhwe95gswv6a",
			"reward": {
				"denom": "ucmdx",
				"amount": "211"
			}
		},
		{
			"address": "comdex1gluq7t59ed9nv4e0jf9xz5s5xz3pvgrunuqnx5",
			"reward": {
				"denom": "ucmdx",
				"amount": "15218"
			}
		},
		{
			"address": "comdex1fqqehm5wef04lvnywsdysrguyez5l9tkjhkh60",
			"reward": {
				"denom": "ucmdx",
				"amount": "339"
			}
		},
		{
			"address": "comdex1fqrxu494h6sn76atlmdcujvw2jqn4rs4d4wsun",
			"reward": {
				"denom": "ucmdx",
				"amount": "394"
			}
		},
		{
			"address": "comdex1fqrc4n78yn3ej957yfl4t3t9amdr7fd02mq76s",
			"reward": {
				"denom": "ucmdx",
				"amount": "4498"
			}
		},
		{
			"address": "comdex1fqyejfl4ryzl6fdr9g7wvzgjs98l6yvf3lqjvg",
			"reward": {
				"denom": "ucmdx",
				"amount": "1190"
			}
		},
		{
			"address": "comdex1fq9qrrwldtczqnkh0lnddtgmeyj55tqlam52cl",
			"reward": {
				"denom": "ucmdx",
				"amount": "227"
			}
		},
		{
			"address": "comdex1fq9yfde2guss52hmleju36v7dztqg56kyk0fvn",
			"reward": {
				"denom": "ucmdx",
				"amount": "574"
			}
		},
		{
			"address": "comdex1fqx7felm86l56r38t3thzq3w83mdqqrmqzt90w",
			"reward": {
				"denom": "ucmdx",
				"amount": "201"
			}
		},
		{
			"address": "comdex1fq2ky4ru3s5lh39tqrd6mrpgunydr9rna322w0",
			"reward": {
				"denom": "ucmdx",
				"amount": "1457"
			}
		},
		{
			"address": "comdex1fqt79f5enr27wtlrzfr0ksn852qa25957pd2rv",
			"reward": {
				"denom": "ucmdx",
				"amount": "14226"
			}
		},
		{
			"address": "comdex1fqvhg9ct4ap7j67djrx3v48qjullar7k0raa8l",
			"reward": {
				"denom": "ucmdx",
				"amount": "75776"
			}
		},
		{
			"address": "comdex1fqw7zrxwly0947pqx6s4c0zu2a578lvd58t43e",
			"reward": {
				"denom": "ucmdx",
				"amount": "16"
			}
		},
		{
			"address": "comdex1fq00pc2wd8netcrfxh6434gvqz390u7ag250el",
			"reward": {
				"denom": "ucmdx",
				"amount": "364"
			}
		},
		{
			"address": "comdex1fq0klaec7s59drpl93w6yycgqqtd3d38d34zn0",
			"reward": {
				"denom": "ucmdx",
				"amount": "21242"
			}
		},
		{
			"address": "comdex1fq0778cyp9f7ffhk9yyg56v7jz9ylaxpwhfxfe",
			"reward": {
				"denom": "ucmdx",
				"amount": "515"
			}
		},
		{
			"address": "comdex1fqs77nt3k4fnssxjj2qrwm3vls30kktn2fwkh2",
			"reward": {
				"denom": "ucmdx",
				"amount": "884271"
			}
		},
		{
			"address": "comdex1fq3m8eek6n9c5phjvmnfggquqqm38hn8gxt4eq",
			"reward": {
				"denom": "ucmdx",
				"amount": "101"
			}
		},
		{
			"address": "comdex1fqj5rd8gmr35tck9cfqg5d395rumk87j5yhfkt",
			"reward": {
				"denom": "ucmdx",
				"amount": "21991"
			}
		},
		{
			"address": "comdex1fqnr0dujf6whemur2tzpvlr5jxmrws5wfye8pt",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1fq5twlmmmqhahjpy0h7sleskv4vldmdzry9x5d",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1fq5nkrgdvk9n9gmgy2heawqg0chzxnps2v8u2v",
			"reward": {
				"denom": "ucmdx",
				"amount": "17921"
			}
		},
		{
			"address": "comdex1fq4j7y5jfsw5nw76ndhsmyzmsach04menq0j99",
			"reward": {
				"denom": "ucmdx",
				"amount": "44"
			}
		},
		{
			"address": "comdex1fqkas8hvmcgcxh3rcptz3h26t4lzraraq7hgua",
			"reward": {
				"denom": "ucmdx",
				"amount": "43435"
			}
		},
		{
			"address": "comdex1fqczm6tv2k4cjxfakfxgjm9acstnc2jvzkmtv0",
			"reward": {
				"denom": "ucmdx",
				"amount": "1773"
			}
		},
		{
			"address": "comdex1fqc2xv645j9xldgnfnhtzcslrw5vm784nc8wjx",
			"reward": {
				"denom": "ucmdx",
				"amount": "145"
			}
		},
		{
			"address": "comdex1fqc3g06f0klxft0uv8pw28nlln5z3z4m4m88dn",
			"reward": {
				"denom": "ucmdx",
				"amount": "615"
			}
		},
		{
			"address": "comdex1fqcekpkjele8hvupu63e0hzfknetpp9hynjhnx",
			"reward": {
				"denom": "ucmdx",
				"amount": "4022"
			}
		},
		{
			"address": "comdex1fqe279p2zkj9nme5p7c7k3sk8h5q7n00rz79pd",
			"reward": {
				"denom": "ucmdx",
				"amount": "44962"
			}
		},
		{
			"address": "comdex1fprwur639dp73cfcprt6dt9r43yzfzxljapm6f",
			"reward": {
				"denom": "ucmdx",
				"amount": "11763"
			}
		},
		{
			"address": "comdex1fpflnzvt207znktmyywgwqs2p3axtvdgc4egcs",
			"reward": {
				"denom": "ucmdx",
				"amount": "5809"
			}
		},
		{
			"address": "comdex1fp2ph05q34ed4q5tppurenh3as6r6dfkhy4fma",
			"reward": {
				"denom": "ucmdx",
				"amount": "3393"
			}
		},
		{
			"address": "comdex1fpwhmczw349xjlwwmxhwggdnjs949tscxt84fc",
			"reward": {
				"denom": "ucmdx",
				"amount": "6526"
			}
		},
		{
			"address": "comdex1fp0w3n5mjzp4aqc3y6pwnjd62zepzvk6n7ptvm",
			"reward": {
				"denom": "ucmdx",
				"amount": "15"
			}
		},
		{
			"address": "comdex1fp37ep4m3eaa887zxt5z6wfg8nrm2td3fk8vu4",
			"reward": {
				"denom": "ucmdx",
				"amount": "122"
			}
		},
		{
			"address": "comdex1fpjpcg4t4zsrgsjmllrvdz809vy0nmkqyu0sz3",
			"reward": {
				"denom": "ucmdx",
				"amount": "879"
			}
		},
		{
			"address": "comdex1fpu6987zgdjf65yeyr2uuut82talr59ggvmnq6",
			"reward": {
				"denom": "ucmdx",
				"amount": "3541"
			}
		},
		{
			"address": "comdex1fpasqaze27syq96kp986kjrug35vshg2gl3hyy",
			"reward": {
				"denom": "ucmdx",
				"amount": "6325"
			}
		},
		{
			"address": "comdex1fp7z8d7n6zts2gmnrvstltq9ndcdenef4acrz0",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1fp7lar7valvpajdxealzhaf5julcpvrlzxj7cf",
			"reward": {
				"denom": "ucmdx",
				"amount": "35"
			}
		},
		{
			"address": "comdex1fzq595eyhdnstdm7g2stvqkj7fj99l5h0jy73r",
			"reward": {
				"denom": "ucmdx",
				"amount": "1914"
			}
		},
		{
			"address": "comdex1fzpej8tyu5g02kpx3usqc4s7q29t42nvyw8tx7",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1fzzs0n252zx3j87e3jl73d68ujd0xfrzhj5y99",
			"reward": {
				"denom": "ucmdx",
				"amount": "14"
			}
		},
		{
			"address": "comdex1fzzutjw4cmw80ja8kpy6we4xleknmm2erpsxpy",
			"reward": {
				"denom": "ucmdx",
				"amount": "598089"
			}
		},
		{
			"address": "comdex1fzzacxdzy0ndq3nz2mvx2rzzujvzu4gexkvs43",
			"reward": {
				"denom": "ucmdx",
				"amount": "16014"
			}
		},
		{
			"address": "comdex1fzrgh2uhwffqatdhxezqgqktn08nrzjxadnwkf",
			"reward": {
				"denom": "ucmdx",
				"amount": "298"
			}
		},
		{
			"address": "comdex1fz9qjhg50ptqvk5y0dp4fdgdgf7f6twdmf8umr",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1fzgrznluvwcze0ahg6a4hj8un855d5lx42pytz",
			"reward": {
				"denom": "ucmdx",
				"amount": "5799"
			}
		},
		{
			"address": "comdex1fzt88qytm5dlpr0u9av69ldp2tmfc8djxhctgd",
			"reward": {
				"denom": "ucmdx",
				"amount": "1494"
			}
		},
		{
			"address": "comdex1fzs6yt4f2fnu0j2r0xsdwf38km4faqh74gfl0m",
			"reward": {
				"denom": "ucmdx",
				"amount": "24585"
			}
		},
		{
			"address": "comdex1fz3rfa04y22dy2snemdu3336y8w3gr0r7vj5za",
			"reward": {
				"denom": "ucmdx",
				"amount": "11110"
			}
		},
		{
			"address": "comdex1fzjckfkh0mua4g94lluqqv89zqkj4466tgqw6n",
			"reward": {
				"denom": "ucmdx",
				"amount": "188"
			}
		},
		{
			"address": "comdex1fzn4k96d5yclqewj2whr5lp8dj6tw08hyv5gzw",
			"reward": {
				"denom": "ucmdx",
				"amount": "127901"
			}
		},
		{
			"address": "comdex1fzkxv7k8hlcphdcy0su38w3yylzvs8mh55rpnd",
			"reward": {
				"denom": "ucmdx",
				"amount": "21354"
			}
		},
		{
			"address": "comdex1fz6t2eqz78ycn7jphd3pttttwmvrj0ghx9dtc8",
			"reward": {
				"denom": "ucmdx",
				"amount": "442"
			}
		},
		{
			"address": "comdex1frzs8kdcdnzhamdgu64xp745xndaqfytq3k874",
			"reward": {
				"denom": "ucmdx",
				"amount": "14978"
			}
		},
		{
			"address": "comdex1fry7ag2hdmdlc5unpw9yg4jr3vyuc0smpgm4l3",
			"reward": {
				"denom": "ucmdx",
				"amount": "179"
			}
		},
		{
			"address": "comdex1frgqk7hu8u0awtnst9z3v7vv32px3wvk9ayg4z",
			"reward": {
				"denom": "ucmdx",
				"amount": "19988"
			}
		},
		{
			"address": "comdex1frgzpm634c6k5dwgc22yszfc8zp3j2jcgpehtt",
			"reward": {
				"denom": "ucmdx",
				"amount": "14444"
			}
		},
		{
			"address": "comdex1frvck0gl6e6skayvkl7ez65yqu73e3cc8wljaa",
			"reward": {
				"denom": "ucmdx",
				"amount": "173"
			}
		},
		{
			"address": "comdex1fr3sd8rzfehjv86ukcdscp8q5rlgvyn359u43e",
			"reward": {
				"denom": "ucmdx",
				"amount": "166"
			}
		},
		{
			"address": "comdex1frh2rc78xlyvfa2yfcnxgnxpcd2ssc0u03yfm6",
			"reward": {
				"denom": "ucmdx",
				"amount": "1669"
			}
		},
		{
			"address": "comdex1frhwljnzzdx0x6z00jdpzewv2vwgduverye8d0",
			"reward": {
				"denom": "ucmdx",
				"amount": "16"
			}
		},
		{
			"address": "comdex1frcxh6xdxn0szxqsz6pswzcf9rryzcycgz2t82",
			"reward": {
				"denom": "ucmdx",
				"amount": "6892"
			}
		},
		{
			"address": "comdex1frmqz0096pjfhqk32pa35lwhmzyzwelgpl5gug",
			"reward": {
				"denom": "ucmdx",
				"amount": "1277"
			}
		},
		{
			"address": "comdex1fr7n0jf9jq9ldcxhvuka39n4kd40ue9yv60zep",
			"reward": {
				"denom": "ucmdx",
				"amount": "112"
			}
		},
		{
			"address": "comdex1frl5yg39fxp4u3jtex3hwulh42futg7h9rfjlu",
			"reward": {
				"denom": "ucmdx",
				"amount": "1485"
			}
		},
		{
			"address": "comdex1fyqfehvyh7xfky6drnmmtq2rc3cd8nku76weh9",
			"reward": {
				"denom": "ucmdx",
				"amount": "167"
			}
		},
		{
			"address": "comdex1fyq4lf4cvfc8ytfn93jmw4lktekhahqqhslvcp",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1fypru5z563pesx7gs7q7xejxa0qswr2d7ltgwk",
			"reward": {
				"denom": "ucmdx",
				"amount": "2829"
			}
		},
		{
			"address": "comdex1fyp0lzq9l6rc85cu2udaaj3mda7up4zh5n5dl9",
			"reward": {
				"denom": "ucmdx",
				"amount": "1190"
			}
		},
		{
			"address": "comdex1fy9w6fu9axem3cmp3u8903ptuxnnjd6h69vy2k",
			"reward": {
				"denom": "ucmdx",
				"amount": "3749"
			}
		},
		{
			"address": "comdex1fyxkvm49cmg0al7tn7k2a2zjxfpfxg8anjle4d",
			"reward": {
				"denom": "ucmdx",
				"amount": "1435"
			}
		},
		{
			"address": "comdex1fy2dpfhs2halvy2ukptf00ah5untplvjfjncuv",
			"reward": {
				"denom": "ucmdx",
				"amount": "10160"
			}
		},
		{
			"address": "comdex1fysj4cadldlzcpjkfm2apdxex48wlvrgm5a9ju",
			"reward": {
				"denom": "ucmdx",
				"amount": "550"
			}
		},
		{
			"address": "comdex1fy3tq59uqycc55ywwhpay03yglwfhra9djvy3s",
			"reward": {
				"denom": "ucmdx",
				"amount": "4192"
			}
		},
		{
			"address": "comdex1fyja53h45tls4yrf0mwka4wzgn46q5avuu45kk",
			"reward": {
				"denom": "ucmdx",
				"amount": "10013"
			}
		},
		{
			"address": "comdex1fynvtycdzysvv3lf4snzx6ctsxqf9qapr8sz9a",
			"reward": {
				"denom": "ucmdx",
				"amount": "4137"
			}
		},
		{
			"address": "comdex1fy4cxtyuwh8dzzweawneudhrdqmaaj4x6lxg8g",
			"reward": {
				"denom": "ucmdx",
				"amount": "1506"
			}
		},
		{
			"address": "comdex1fye2m0gm5llkpetucdjefq6hxu2tjvvwvps3vw",
			"reward": {
				"denom": "ucmdx",
				"amount": "909"
			}
		},
		{
			"address": "comdex1fy6n0tg7s7rxaduw09h4wyuu2fp7ldy9vgrdc7",
			"reward": {
				"denom": "ucmdx",
				"amount": "1663"
			}
		},
		{
			"address": "comdex1fy7gk4uvulm06p673qx3tes6g3zz423vksn8g4",
			"reward": {
				"denom": "ucmdx",
				"amount": "1437"
			}
		},
		{
			"address": "comdex1fylfk82mfd80qkyjvzwearuhd6fp226acxj56l",
			"reward": {
				"denom": "ucmdx",
				"amount": "276"
			}
		},
		{
			"address": "comdex1f9qdhy4pdj89tjfhvglu7rpl8c2w28eajqm5xe",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1f9pqtdpn8puanen76tr7x3574tsw9l87l9uzha",
			"reward": {
				"denom": "ucmdx",
				"amount": "2628"
			}
		},
		{
			"address": "comdex1f9r4xgc5xpwq654hrm9jf5q7uj589mlt8wc8gm",
			"reward": {
				"denom": "ucmdx",
				"amount": "10163"
			}
		},
		{
			"address": "comdex1f9ra3zyj3qflzt4jhg6qmva5sdgz77pztwuqu3",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1f9g0paz9yedq363ww97v3d43e6gze6wq2l2wp0",
			"reward": {
				"denom": "ucmdx",
				"amount": "151"
			}
		},
		{
			"address": "comdex1f9fq06z9687ffnundg62tjfv24r2redm38l24w",
			"reward": {
				"denom": "ucmdx",
				"amount": "1393"
			}
		},
		{
			"address": "comdex1f92lgeqllnej7aznxmxjmfw2fn0auh3ky5ltpa",
			"reward": {
				"denom": "ucmdx",
				"amount": "4310"
			}
		},
		{
			"address": "comdex1f90q2peljkhxneye7k3slg99e7hupf58gh46xz",
			"reward": {
				"denom": "ucmdx",
				"amount": "151"
			}
		},
		{
			"address": "comdex1f9sglztnflh3yk73lejztnazljzp85cdrx0a7w",
			"reward": {
				"denom": "ucmdx",
				"amount": "28348"
			}
		},
		{
			"address": "comdex1f9jtg7zyesxsfvsnzs0uth2auxz36ffdydqltq",
			"reward": {
				"denom": "ucmdx",
				"amount": "15700"
			}
		},
		{
			"address": "comdex1f95dmun7ldnyzymrja6qjy063nfug0kzhf07xt",
			"reward": {
				"denom": "ucmdx",
				"amount": "1234"
			}
		},
		{
			"address": "comdex1f9k272cgejjdt0exugldcgzlf46d00ttwdjcrs",
			"reward": {
				"denom": "ucmdx",
				"amount": "1996"
			}
		},
		{
			"address": "comdex1f9ktx3f6p8cdtmmtx35k3lfexxt0x8yte4utmg",
			"reward": {
				"denom": "ucmdx",
				"amount": "140"
			}
		},
		{
			"address": "comdex1f9h82d5x93zkw8c09twrhp9famrs07hdzyta6a",
			"reward": {
				"denom": "ucmdx",
				"amount": "1582"
			}
		},
		{
			"address": "comdex1f9cz2nefyfr2c9wjvkyzkuddn7x5fjlhelm0rt",
			"reward": {
				"denom": "ucmdx",
				"amount": "35648"
			}
		},
		{
			"address": "comdex1f9cmjlkwqrklu6hlhk8a4atw9wnawj2juejgzh",
			"reward": {
				"denom": "ucmdx",
				"amount": "142"
			}
		},
		{
			"address": "comdex1f9u3yfk43kan4d0k4zwr6cezclusrerfgjt7ty",
			"reward": {
				"denom": "ucmdx",
				"amount": "4041"
			}
		},
		{
			"address": "comdex1f9luj7hx5hl9dfh45c6u0fnv0gflwff2w50h2n",
			"reward": {
				"denom": "ucmdx",
				"amount": "64721"
			}
		},
		{
			"address": "comdex1fxqqp5xaxeg0ywy5q3z42r6j2gn3m9g9vq5cgz",
			"reward": {
				"denom": "ucmdx",
				"amount": "1976"
			}
		},
		{
			"address": "comdex1fxpxw6lha2gguctz9s4uqv30rvm7d3e20ju9rd",
			"reward": {
				"denom": "ucmdx",
				"amount": "181"
			}
		},
		{
			"address": "comdex1fxrqjv557r9yy0r5fx7gzatyp3j4kl9llw05pm",
			"reward": {
				"denom": "ucmdx",
				"amount": "576"
			}
		},
		{
			"address": "comdex1fx99zxenxtzaa44k8n7049c3umaud6s7gddytc",
			"reward": {
				"denom": "ucmdx",
				"amount": "30"
			}
		},
		{
			"address": "comdex1fxxy0keh80cxmwl886usl4085d3gx444hlgqer",
			"reward": {
				"denom": "ucmdx",
				"amount": "345"
			}
		},
		{
			"address": "comdex1fxx4hsm9zxw6cctu9mtdt4apcj9xs00tt47wwp",
			"reward": {
				"denom": "ucmdx",
				"amount": "713"
			}
		},
		{
			"address": "comdex1fx8kl43rkcmkk98p42fjefmkr9kpg2nce7f2vd",
			"reward": {
				"denom": "ucmdx",
				"amount": "525"
			}
		},
		{
			"address": "comdex1fxtg53uexmvh47wppasxjpee4u0erv2lqw5fjh",
			"reward": {
				"denom": "ucmdx",
				"amount": "1746"
			}
		},
		{
			"address": "comdex1fxsw204nrg4l0qql7a3s27r0jt82kzrly7d9rf",
			"reward": {
				"denom": "ucmdx",
				"amount": "1335"
			}
		},
		{
			"address": "comdex1fx3fedyn4jwmkyglpe0p2xravc5wjvdr2t4sg6",
			"reward": {
				"denom": "ucmdx",
				"amount": "6514"
			}
		},
		{
			"address": "comdex1fxj0s3n8d53v6xhzp8wrf2fpnupte7v3gunlxr",
			"reward": {
				"denom": "ucmdx",
				"amount": "30297"
			}
		},
		{
			"address": "comdex1fxnvs9kmva2e4ynv9fkmn8e5ysu8z268ln2n7x",
			"reward": {
				"denom": "ucmdx",
				"amount": "28340"
			}
		},
		{
			"address": "comdex1fxn0dw8g34m6y43qskneyyw2f2hartmzcgqn2c",
			"reward": {
				"denom": "ucmdx",
				"amount": "47"
			}
		},
		{
			"address": "comdex1fxnkzpsggexfneqnucf2wmfy0zmgj7s224pacp",
			"reward": {
				"denom": "ucmdx",
				"amount": "88"
			}
		},
		{
			"address": "comdex1fxh8e0e0nl72xkn6sxtxe5qlky3mfjryegr426",
			"reward": {
				"denom": "ucmdx",
				"amount": "5314"
			}
		},
		{
			"address": "comdex1fxc2ktxjjstns53qd8v8llzm6gt0d9fct6ywst",
			"reward": {
				"denom": "ucmdx",
				"amount": "70"
			}
		},
		{
			"address": "comdex1fxexqgzsk9sjw9t4uee4vppcanr72lh2r5c36k",
			"reward": {
				"denom": "ucmdx",
				"amount": "14384"
			}
		},
		{
			"address": "comdex1fx6r88pzy2yc9ssgekwtumpkcjxlkkx5cjcm3z",
			"reward": {
				"denom": "ucmdx",
				"amount": "860"
			}
		},
		{
			"address": "comdex1fxma5wj52f7mhsrjquq3vw8gyjj9h34mw9svg3",
			"reward": {
				"denom": "ucmdx",
				"amount": "1252"
			}
		},
		{
			"address": "comdex1fxua2h6u6p9sd80j2svy7zu9fhhdxgfju4yj4y",
			"reward": {
				"denom": "ucmdx",
				"amount": "1808"
			}
		},
		{
			"address": "comdex1fx7e4msgtdmjtvxtfakq64qrec30a6wkhcxrw8",
			"reward": {
				"denom": "ucmdx",
				"amount": "179"
			}
		},
		{
			"address": "comdex1fx778w26vrrlwpjg2yc2c2xt53x4n8vxhzfpvd",
			"reward": {
				"denom": "ucmdx",
				"amount": "6457"
			}
		},
		{
			"address": "comdex1f8pwg36m887e73v7rsn8w5wkdjfxr8x82alsc4",
			"reward": {
				"denom": "ucmdx",
				"amount": "1486"
			}
		},
		{
			"address": "comdex1f8zvcwcnn259vhm7a7ru95te4g7rxgs4ad8dsh",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1f8z54w7mkzgqu3aq75hn7jt8t6pcmuuqftt865",
			"reward": {
				"denom": "ucmdx",
				"amount": "28"
			}
		},
		{
			"address": "comdex1f8rvqzhw5xh7qp7nhw9n6605njpy025rmhw9ky",
			"reward": {
				"denom": "ucmdx",
				"amount": "5127"
			}
		},
		{
			"address": "comdex1f8rdzyn7uu0yuv5kngwqlahp56y38j4g2rcpwz",
			"reward": {
				"denom": "ucmdx",
				"amount": "1783"
			}
		},
		{
			"address": "comdex1f88jmmzht272elcyuej5cmlthsuydhj59wsdwx",
			"reward": {
				"denom": "ucmdx",
				"amount": "6825"
			}
		},
		{
			"address": "comdex1f8w9duslevxu5qgnldpcy6nv6ke8suq20yqv6w",
			"reward": {
				"denom": "ucmdx",
				"amount": "287"
			}
		},
		{
			"address": "comdex1f8j4jc4qtnz86sgnxwxr0p9s24r6jencxr8393",
			"reward": {
				"denom": "ucmdx",
				"amount": "9379"
			}
		},
		{
			"address": "comdex1f8n7q6e9ty0zdl9a3vm4qv7k744xef82t5yz9v",
			"reward": {
				"denom": "ucmdx",
				"amount": "1621"
			}
		},
		{
			"address": "comdex1f84zfjwcn2sdfp084qwag5qkygwr5x07e5577f",
			"reward": {
				"denom": "ucmdx",
				"amount": "11052"
			}
		},
		{
			"address": "comdex1f8h2uhwk5apr9p23glmp3kwtdmsftxd77jrqna",
			"reward": {
				"denom": "ucmdx",
				"amount": "1404"
			}
		},
		{
			"address": "comdex1f8hm7vhy6le6d6zsd7tde7cfhy5q9ykkcddmjj",
			"reward": {
				"denom": "ucmdx",
				"amount": "4265"
			}
		},
		{
			"address": "comdex1f863jgzxv03ytndt4c4x7cq2xnucj274625uem",
			"reward": {
				"denom": "ucmdx",
				"amount": "1739"
			}
		},
		{
			"address": "comdex1f8l97avty8h3q0krygez8clly3drfkq98d0kpt",
			"reward": {
				"denom": "ucmdx",
				"amount": "9841"
			}
		},
		{
			"address": "comdex1f8l5u9qfppvg6dzw243pf4kxtuzzwxhcenh90l",
			"reward": {
				"denom": "ucmdx",
				"amount": "865"
			}
		},
		{
			"address": "comdex1fgptl3cuatgz3evqf0357q7suh00f9qrrlu9ns",
			"reward": {
				"denom": "ucmdx",
				"amount": "12861"
			}
		},
		{
			"address": "comdex1fgryajvu00554w32dgxeyjtvzmjn4yu8rlcjhr",
			"reward": {
				"denom": "ucmdx",
				"amount": "373"
			}
		},
		{
			"address": "comdex1fg8yyu8jxcskf7frkfkh63tn5wrerfzfn7edtt",
			"reward": {
				"denom": "ucmdx",
				"amount": "482"
			}
		},
		{
			"address": "comdex1fgfljze6gtzlvzq92ek95a25wt24rsnu4j46qf",
			"reward": {
				"denom": "ucmdx",
				"amount": "6877"
			}
		},
		{
			"address": "comdex1fg2d6dsf8n9rhfq559ehvlq5fskdczl202lyax",
			"reward": {
				"denom": "ucmdx",
				"amount": "1024"
			}
		},
		{
			"address": "comdex1fgwcxv057fd3hgsw49tueuu7n2qh4ew4lkmpye",
			"reward": {
				"denom": "ucmdx",
				"amount": "4338"
			}
		},
		{
			"address": "comdex1fgwaqndwpl2hmk6vv98vt07fkefvx46w233xec",
			"reward": {
				"denom": "ucmdx",
				"amount": "10303"
			}
		},
		{
			"address": "comdex1fg3d56zj48xjxflmuy5zf7xqdfcmdnwfxrfhps",
			"reward": {
				"denom": "ucmdx",
				"amount": "1472"
			}
		},
		{
			"address": "comdex1fgngf5pka7m6u68zmdllds7gna7jttw0m9qutd",
			"reward": {
				"denom": "ucmdx",
				"amount": "359"
			}
		},
		{
			"address": "comdex1fgnfg2gnxhkmcc8tgpph6jm6xvn236jeyg3jmt",
			"reward": {
				"denom": "ucmdx",
				"amount": "142524"
			}
		},
		{
			"address": "comdex1fg5y7vq8cukksxqlvpu28uy77pvete0c480dw6",
			"reward": {
				"denom": "ucmdx",
				"amount": "8251"
			}
		},
		{
			"address": "comdex1fg54f87w3pha68ncqzvshat828rt5autau2r9u",
			"reward": {
				"denom": "ucmdx",
				"amount": "293288"
			}
		},
		{
			"address": "comdex1fg54luaem7966d6tnrk4vtuu6uvzw6ra696wrc",
			"reward": {
				"denom": "ucmdx",
				"amount": "1003"
			}
		},
		{
			"address": "comdex1fg4cjjc4sw0c2xyvkzfph6yrwcnncpz3e64pv2",
			"reward": {
				"denom": "ucmdx",
				"amount": "1036"
			}
		},
		{
			"address": "comdex1fgml4dthl876mglampyjgnrf35umee0c8lhpx0",
			"reward": {
				"denom": "ucmdx",
				"amount": "3923"
			}
		},
		{
			"address": "comdex1fgurqdwm38hpsl3y4l4upvsgt46jarsnramazq",
			"reward": {
				"denom": "ucmdx",
				"amount": "1943"
			}
		},
		{
			"address": "comdex1fg7zc653w4d6tssaw4y4vvmkc7jtlvvsnc8lkr",
			"reward": {
				"denom": "ucmdx",
				"amount": "25175"
			}
		},
		{
			"address": "comdex1fg7fcs0vupc5lyy67hatkq2nak995w3mpz6gdt",
			"reward": {
				"denom": "ucmdx",
				"amount": "440"
			}
		},
		{
			"address": "comdex1fglyz77v8yuaq2saewazwm8qvwz5pv0txl74w3",
			"reward": {
				"denom": "ucmdx",
				"amount": "1758"
			}
		},
		{
			"address": "comdex1ffp257p3zymh32qtk2l0ymfpv4pgxeyfzdd0he",
			"reward": {
				"denom": "ucmdx",
				"amount": "69391"
			}
		},
		{
			"address": "comdex1ffpdl72gvhwm82anjlzg6r0pjuj6y5snkvcjc7",
			"reward": {
				"denom": "ucmdx",
				"amount": "1630"
			}
		},
		{
			"address": "comdex1ffzfm7m8lutgqg265mvuymuqxhyjlzyclm8hng",
			"reward": {
				"denom": "ucmdx",
				"amount": "1972"
			}
		},
		{
			"address": "comdex1ffz6sey8ansjw4qcmwfhtzqf28lpe7rnzvl9es",
			"reward": {
				"denom": "ucmdx",
				"amount": "166"
			}
		},
		{
			"address": "comdex1ff9tskaz9faad5ajg2sqhwmuk9nht3hy366wp7",
			"reward": {
				"denom": "ucmdx",
				"amount": "133"
			}
		},
		{
			"address": "comdex1ffxpzvmkh8nhl6v67wgpkm5zqs2p3d75hhhztg",
			"reward": {
				"denom": "ucmdx",
				"amount": "72679"
			}
		},
		{
			"address": "comdex1ff8pcu7yfxx8256r4agtz5z83tccfn0nxxuvam",
			"reward": {
				"denom": "ucmdx",
				"amount": "602"
			}
		},
		{
			"address": "comdex1ffgztpzwm0c472nlpd3mf7kny98xfr7m8hwecp",
			"reward": {
				"denom": "ucmdx",
				"amount": "1390"
			}
		},
		{
			"address": "comdex1fftgce74e9qg7cenj42s4h452m8830uyxrrv6q",
			"reward": {
				"denom": "ucmdx",
				"amount": "25060"
			}
		},
		{
			"address": "comdex1ffkpmjxz7wqy5qdg7qafzwgpdw8qn3wg07d0cr",
			"reward": {
				"denom": "ucmdx",
				"amount": "9800"
			}
		},
		{
			"address": "comdex1ffhy04uwhsvyd0cxnlsdr9r55yc54s5hg8w74h",
			"reward": {
				"denom": "ucmdx",
				"amount": "4013"
			}
		},
		{
			"address": "comdex1ff6saslfvcx223zdpadv2hqwu8efuxflhr7t9g",
			"reward": {
				"denom": "ucmdx",
				"amount": "1007"
			}
		},
		{
			"address": "comdex1ffmtjltr4lfmlgdj0uq9e7zjjwc06h59y8ukz5",
			"reward": {
				"denom": "ucmdx",
				"amount": "16582"
			}
		},
		{
			"address": "comdex1ffuq59z3z9632ml963awvf67geut60dgtcsq4j",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1ff7qd4f83dcqpy4r68vwtu5sf7d2df2da5geza",
			"reward": {
				"denom": "ucmdx",
				"amount": "35841"
			}
		},
		{
			"address": "comdex1ff78n2kt88m465enu080lnnqcg5x8c3zhhuyas",
			"reward": {
				"denom": "ucmdx",
				"amount": "613428"
			}
		},
		{
			"address": "comdex1f2qch7rj6gf3dag79ncwqehgn5axjr70c5a44k",
			"reward": {
				"denom": "ucmdx",
				"amount": "2203"
			}
		},
		{
			"address": "comdex1f2qalvnfu9mqaqwh5exjnvwc3f2ngu9ycc7ng9",
			"reward": {
				"denom": "ucmdx",
				"amount": "12327"
			}
		},
		{
			"address": "comdex1f2px0lzs46n9mf98ug7lzv6kd6d0zv85nwpgmp",
			"reward": {
				"denom": "ucmdx",
				"amount": "21854"
			}
		},
		{
			"address": "comdex1f2yt92vc836gy5dwst7pjsfp3ultajfxdegaap",
			"reward": {
				"denom": "ucmdx",
				"amount": "14215"
			}
		},
		{
			"address": "comdex1f292pju8lmezdw252asa94dxu0yrtg2vqqh0w5",
			"reward": {
				"denom": "ucmdx",
				"amount": "271"
			}
		},
		{
			"address": "comdex1f20dskvw2t42l4m9qxqzexgf37kud4glc48066",
			"reward": {
				"denom": "ucmdx",
				"amount": "700"
			}
		},
		{
			"address": "comdex1f2047lv8dt84vmwnls35cdwa2rk23ktk0yxf8n",
			"reward": {
				"denom": "ucmdx",
				"amount": "10210"
			}
		},
		{
			"address": "comdex1f23quu6d3306pezxk7dqaqz4twdtctludhfyxp",
			"reward": {
				"denom": "ucmdx",
				"amount": "1754"
			}
		},
		{
			"address": "comdex1f2ha75vnuqq58uknkvx5dnznncl7z5u6u068n6",
			"reward": {
				"denom": "ucmdx",
				"amount": "2456"
			}
		},
		{
			"address": "comdex1f2cwy3dux5w4u06245lr7f3wpskcgunv08hanc",
			"reward": {
				"denom": "ucmdx",
				"amount": "71709"
			}
		},
		{
			"address": "comdex1f2ewqncg2t0fax00clrvwa9nx8d2hd8ds87nza",
			"reward": {
				"denom": "ucmdx",
				"amount": "1494"
			}
		},
		{
			"address": "comdex1f2638sl2w9utnr020rzlv9h84aeqa700hwpzwv",
			"reward": {
				"denom": "ucmdx",
				"amount": "1416"
			}
		},
		{
			"address": "comdex1f277jsrcx5tfc5n2xp5rncflr9yj225tqw3fcn",
			"reward": {
				"denom": "ucmdx",
				"amount": "179"
			}
		},
		{
			"address": "comdex1ftp0qsw50ksa0ax7sq5l4m5qgfgvh66ppm7krg",
			"reward": {
				"denom": "ucmdx",
				"amount": "6964"
			}
		},
		{
			"address": "comdex1ftyjxtlg4hehasj584dul70ctvwlcwsyvj7feg",
			"reward": {
				"denom": "ucmdx",
				"amount": "2144"
			}
		},
		{
			"address": "comdex1ftg95ck5ltm639s3x5w9n4kqjgc68d0z3ad3xl",
			"reward": {
				"denom": "ucmdx",
				"amount": "36553"
			}
		},
		{
			"address": "comdex1ftvd054ujth0jk2m9t53qz5a08e99atdkeh6xj",
			"reward": {
				"denom": "ucmdx",
				"amount": "60"
			}
		},
		{
			"address": "comdex1ftw38q80aacjyjuz5xhxcc5t6m3rvp798dslyk",
			"reward": {
				"denom": "ucmdx",
				"amount": "252"
			}
		},
		{
			"address": "comdex1ft06k8yzxyxnte9s0mh0fmravyrev3pas2e5zg",
			"reward": {
				"denom": "ucmdx",
				"amount": "2863"
			}
		},
		{
			"address": "comdex1ftnnzaqgd77epwcemmqf6l4nah44595teckww9",
			"reward": {
				"denom": "ucmdx",
				"amount": "1781"
			}
		},
		{
			"address": "comdex1ft5fgsj2r6mjj7mvl7feqxcejcz2rn8qqj59j9",
			"reward": {
				"denom": "ucmdx",
				"amount": "86"
			}
		},
		{
			"address": "comdex1ftkg0tzfw54k4862a8ja4637xy2jerg34w7fcj",
			"reward": {
				"denom": "ucmdx",
				"amount": "1503"
			}
		},
		{
			"address": "comdex1ftk222zvd5aecyxazxvj3zukcj9ptcqjcmttpl",
			"reward": {
				"denom": "ucmdx",
				"amount": "1243"
			}
		},
		{
			"address": "comdex1ftk4p9a2yg0hxja9mrlahrlrzeatd9jwy40aka",
			"reward": {
				"denom": "ucmdx",
				"amount": "19"
			}
		},
		{
			"address": "comdex1fthr2r7avaztg4fpgjmccft29u7y9rds7dsze3",
			"reward": {
				"denom": "ucmdx",
				"amount": "890"
			}
		},
		{
			"address": "comdex1ftesdfme6gnkgvyty9c2en2620r94p298jdam2",
			"reward": {
				"denom": "ucmdx",
				"amount": "3421"
			}
		},
		{
			"address": "comdex1ft6gdwcedwdqe5gsp4uhcdvze6ualqg30jlyx5",
			"reward": {
				"denom": "ucmdx",
				"amount": "373970"
			}
		},
		{
			"address": "comdex1ftmqcym4qk3wx0dtvzccyjxhphmyjqc3u7fp7f",
			"reward": {
				"denom": "ucmdx",
				"amount": "16158"
			}
		},
		{
			"address": "comdex1ftmtmxwm7ynhcj0s5ud3hktwntgy3pktpvyh9j",
			"reward": {
				"denom": "ucmdx",
				"amount": "25515"
			}
		},
		{
			"address": "comdex1ftmlhq26yc2yadh3ejfuqm2njgcm6l2ayc3g64",
			"reward": {
				"denom": "ucmdx",
				"amount": "254"
			}
		},
		{
			"address": "comdex1ftl7n3pr0qruh04m7d2jd2k7hrl7vn3awavzku",
			"reward": {
				"denom": "ucmdx",
				"amount": "6090"
			}
		},
		{
			"address": "comdex1fvpysd380759nkyfzrc5d0vuduuaukwuw3yjg5",
			"reward": {
				"denom": "ucmdx",
				"amount": "199"
			}
		},
		{
			"address": "comdex1fvpl8sct7jgfmlyu85gsxt8qxvaff8cdtnp4d4",
			"reward": {
				"denom": "ucmdx",
				"amount": "142"
			}
		},
		{
			"address": "comdex1fvzu0l7vv7d43mnd5qjv236r0wznyqx7jk4ekj",
			"reward": {
				"denom": "ucmdx",
				"amount": "75046"
			}
		},
		{
			"address": "comdex1fvylw0ue3mpl2zv09w3djm4a33ay2hlxa4kjxk",
			"reward": {
				"denom": "ucmdx",
				"amount": "284"
			}
		},
		{
			"address": "comdex1fv9trv880n6eglqdetd23y9v829plnm5p28vcu",
			"reward": {
				"denom": "ucmdx",
				"amount": "176"
			}
		},
		{
			"address": "comdex1fv2kkjqt2hv6cn2c0du5qqxkeyvtzjxuqt46hd",
			"reward": {
				"denom": "ucmdx",
				"amount": "2046"
			}
		},
		{
			"address": "comdex1fvwyg9sv8lqsgn59zcptfumhaztu527pr9dnm4",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1fvwhg2yq5n493geq3498mptu09r9u6vf3l8l9d",
			"reward": {
				"denom": "ucmdx",
				"amount": "41713"
			}
		},
		{
			"address": "comdex1fvjf3l0vdccl6a6hshsde8ykdjl2x3ddy0jgxa",
			"reward": {
				"denom": "ucmdx",
				"amount": "853357"
			}
		},
		{
			"address": "comdex1fv54jjs9h8qeuaqj9x56e3d2rrcmht8ystrug7",
			"reward": {
				"denom": "ucmdx",
				"amount": "207"
			}
		},
		{
			"address": "comdex1fv56nhej78tj58nc43pcr3f4ejy6j4u39z8g4a",
			"reward": {
				"denom": "ucmdx",
				"amount": "12706"
			}
		},
		{
			"address": "comdex1fv6cv8er98dm7ctt8y6aqm0trqmdfghfjsu5xt",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1fvmj7j4akqtrfcw6njszu4dvvjkvtfa92ad4x2",
			"reward": {
				"denom": "ucmdx",
				"amount": "6847"
			}
		},
		{
			"address": "comdex1fvucvruyhg5tgquk8p6nxgp0vsvsrzru74plzc",
			"reward": {
				"denom": "ucmdx",
				"amount": "2743"
			}
		},
		{
			"address": "comdex1fdpxvwzrekvpwtx3fdhqdkcxjzgw9vl2yywhfr",
			"reward": {
				"denom": "ucmdx",
				"amount": "25230"
			}
		},
		{
			"address": "comdex1fdzswakqnsd79pu8czw9h56mrnl7w9jtk3pe8c",
			"reward": {
				"denom": "ucmdx",
				"amount": "17739"
			}
		},
		{
			"address": "comdex1fd8kd7qkneye4y96slq8992r34lt0wfez3egkm",
			"reward": {
				"denom": "ucmdx",
				"amount": "14019"
			}
		},
		{
			"address": "comdex1fd8mska6vaksen0n2dfwecqupt2f89xqj29uzl",
			"reward": {
				"denom": "ucmdx",
				"amount": "19899"
			}
		},
		{
			"address": "comdex1fd2hlkk8mz4j3stjd0vv894evxhmq0fkd4adld",
			"reward": {
				"denom": "ucmdx",
				"amount": "24976"
			}
		},
		{
			"address": "comdex1fd2ayfu6kmt66jwr62eazet5y5hp8e3ataf5jj",
			"reward": {
				"denom": "ucmdx",
				"amount": "10103"
			}
		},
		{
			"address": "comdex1fdtgua2gam7ujv63c7trlg2nxwqscadmsp7kcl",
			"reward": {
				"denom": "ucmdx",
				"amount": "13310"
			}
		},
		{
			"address": "comdex1fdd3jefvjhk70acygqdx2ptwdmzkp5tm9wx952",
			"reward": {
				"denom": "ucmdx",
				"amount": "28203"
			}
		},
		{
			"address": "comdex1fdsqnm8w60p3ln0572628zd5seqkapptd3vs29",
			"reward": {
				"denom": "ucmdx",
				"amount": "2025"
			}
		},
		{
			"address": "comdex1fdjphfhuqquctyj8vtgf4pw75rh3hdppwqca72",
			"reward": {
				"denom": "ucmdx",
				"amount": "3498"
			}
		},
		{
			"address": "comdex1fdnqfc656lz5dpscg6fg70fep6m5q7xrh2h9xk",
			"reward": {
				"denom": "ucmdx",
				"amount": "1242"
			}
		},
		{
			"address": "comdex1fd52gmwrug0f9cz508clq8ncg78zd7nfwfc5gs",
			"reward": {
				"denom": "ucmdx",
				"amount": "1420"
			}
		},
		{
			"address": "comdex1fdhmjx0h6pndfajpa7tj39xc95qmmd8urgg088",
			"reward": {
				"denom": "ucmdx",
				"amount": "169"
			}
		},
		{
			"address": "comdex1fdcdapw0060kfu8zqq7j5jeussclw48eefqau9",
			"reward": {
				"denom": "ucmdx",
				"amount": "529"
			}
		},
		{
			"address": "comdex1fd6vgzze8hnfkrc6fstqkm2t92nsyemuxch6a9",
			"reward": {
				"denom": "ucmdx",
				"amount": "124984"
			}
		},
		{
			"address": "comdex1fd6m2t4q8mz20xgwrwz2l0ptk4hjaef9gnq327",
			"reward": {
				"denom": "ucmdx",
				"amount": "1761"
			}
		},
		{
			"address": "comdex1fdm6prm6nn9snm7vgj943fz0fgj70xdmlfds7t",
			"reward": {
				"denom": "ucmdx",
				"amount": "29813"
			}
		},
		{
			"address": "comdex1fdau7p0ueyxzlvs6wxam7vxc7d27ncz96lfkm5",
			"reward": {
				"denom": "ucmdx",
				"amount": "3038"
			}
		},
		{
			"address": "comdex1fd7dsk7kqfkj7ttdsp7hnctq060sh0yhsvshcp",
			"reward": {
				"denom": "ucmdx",
				"amount": "272"
			}
		},
		{
			"address": "comdex1fdlt8lgz2qph9n4uyuah0fw2f76lsh03l5x44c",
			"reward": {
				"denom": "ucmdx",
				"amount": "1978"
			}
		},
		{
			"address": "comdex1fdlva58qvktn4a90jfxwmn2padjnaffcmwzv5y",
			"reward": {
				"denom": "ucmdx",
				"amount": "7159"
			}
		},
		{
			"address": "comdex1fwqj9dtlzpsj7fpa6skc4r6dhfdu2gx03qahyd",
			"reward": {
				"denom": "ucmdx",
				"amount": "1403"
			}
		},
		{
			"address": "comdex1fwzqwujy2xpn99vh3k7d066r3wrtjct85njfq8",
			"reward": {
				"denom": "ucmdx",
				"amount": "57915"
			}
		},
		{
			"address": "comdex1fwy26wxc2qk54huhdtaka69p8w4kctrntce8tt",
			"reward": {
				"denom": "ucmdx",
				"amount": "554"
			}
		},
		{
			"address": "comdex1fwtc02xtyeruwuvhd6kh35y9qy7wzpc2ahqv9h",
			"reward": {
				"denom": "ucmdx",
				"amount": "25129"
			}
		},
		{
			"address": "comdex1fwvhp4jp4fc2fwqf6h8zr4gd04at033hjy7ufd",
			"reward": {
				"denom": "ucmdx",
				"amount": "16"
			}
		},
		{
			"address": "comdex1fw0qlvjr68pqy5nrckqspklkqshd4mp58rg39a",
			"reward": {
				"denom": "ucmdx",
				"amount": "14309"
			}
		},
		{
			"address": "comdex1fwsxl0z3eapqzftsa2vd5lu537qreyadzsl2lu",
			"reward": {
				"denom": "ucmdx",
				"amount": "8671"
			}
		},
		{
			"address": "comdex1fwkvttw5dlnte44nx8935ll5s2pzx8rjhc25pu",
			"reward": {
				"denom": "ucmdx",
				"amount": "1791"
			}
		},
		{
			"address": "comdex1fwkj4hfqr8mv2yjqlev803psnjwhamre66ug2g",
			"reward": {
				"denom": "ucmdx",
				"amount": "179"
			}
		},
		{
			"address": "comdex1fw6z9xmglarau8vhyxxj44nwzz0j9tc8ar65jn",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1fw75ufvfh4k7rtzcemfnsr0qur23x574ecg2k9",
			"reward": {
				"denom": "ucmdx",
				"amount": "7884"
			}
		},
		{
			"address": "comdex1fwlsckj782qs806u85rvz445uva74y5xwk6wgk",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1f0pztnf80zkcks6vp8l0az7tkup0afefzsntxs",
			"reward": {
				"denom": "ucmdx",
				"amount": "2057"
			}
		},
		{
			"address": "comdex1f0pdnvurpgu7hx3qgks4u5ct27sh2j8ds60rhk",
			"reward": {
				"denom": "ucmdx",
				"amount": "284"
			}
		},
		{
			"address": "comdex1f0yg9z4dmell7tm9yl8yqx2nc7q5n53cgrj5du",
			"reward": {
				"denom": "ucmdx",
				"amount": "271"
			}
		},
		{
			"address": "comdex1f0ywa7thm4he6ggc3vfwnec3gfc85jpp24vsm3",
			"reward": {
				"denom": "ucmdx",
				"amount": "6796"
			}
		},
		{
			"address": "comdex1f09hfkxe5nfduhqhg408yvtlv27p25p8umkazg",
			"reward": {
				"denom": "ucmdx",
				"amount": "611"
			}
		},
		{
			"address": "comdex1f0fkttthl9aqfd6y2tt86ue0ujup5sft6jqcgn",
			"reward": {
				"denom": "ucmdx",
				"amount": "5363"
			}
		},
		{
			"address": "comdex1f0d6js085gea2qxtyf5eaa5pz67xgy5sz287zd",
			"reward": {
				"denom": "ucmdx",
				"amount": "77"
			}
		},
		{
			"address": "comdex1f0jdnzdjg2670hl0vusnja9kv3ph503jwtkuhn",
			"reward": {
				"denom": "ucmdx",
				"amount": "23636"
			}
		},
		{
			"address": "comdex1f05f676hzfx826u6rh0el8ez84ql0lc6jve8m4",
			"reward": {
				"denom": "ucmdx",
				"amount": "28"
			}
		},
		{
			"address": "comdex1f0kg3546jrguzu63w76zlprzf2a7mv676rlrt6",
			"reward": {
				"denom": "ucmdx",
				"amount": "1931"
			}
		},
		{
			"address": "comdex1f0kvh8jg4z9yue02t4xwdhu5kr50jw60h5fulj",
			"reward": {
				"denom": "ucmdx",
				"amount": "6007"
			}
		},
		{
			"address": "comdex1f0kcpeuys5mhdcpq87xmgh4vny8vwkl068205s",
			"reward": {
				"denom": "ucmdx",
				"amount": "4862"
			}
		},
		{
			"address": "comdex1f0cxsgdflgnxysjvk5wlvhay85qm6w3g0s8ssa",
			"reward": {
				"denom": "ucmdx",
				"amount": "317"
			}
		},
		{
			"address": "comdex1f0cj0v86ewlhp3htg29ac58yxs2uh0tpzsnuzr",
			"reward": {
				"denom": "ucmdx",
				"amount": "37172"
			}
		},
		{
			"address": "comdex1f0cc4hemajysckm7cttm2373jg024ep35s6zqq",
			"reward": {
				"denom": "ucmdx",
				"amount": "2751"
			}
		},
		{
			"address": "comdex1f0696mzr336933ytvms5m0u57mkcvjhnmdsn2e",
			"reward": {
				"denom": "ucmdx",
				"amount": "5341"
			}
		},
		{
			"address": "comdex1f0urc2u0dd3ne6yl28armthxkfhz39szcey9jm",
			"reward": {
				"denom": "ucmdx",
				"amount": "1025"
			}
		},
		{
			"address": "comdex1fszjwzpu7xz2yl5sxxrcrelt4m6y27aks9fu2g",
			"reward": {
				"denom": "ucmdx",
				"amount": "110907"
			}
		},
		{
			"address": "comdex1fsylvlrjfhsqs2sndcg0vy6wcmgv3pj05g4kxy",
			"reward": {
				"denom": "ucmdx",
				"amount": "8220"
			}
		},
		{
			"address": "comdex1fs99pquu4055y754e65909k0z3h5aruplrfe7z",
			"reward": {
				"denom": "ucmdx",
				"amount": "5319"
			}
		},
		{
			"address": "comdex1fs99djlgwcl886jfx73gzt5kfk8v3qx8tlsera",
			"reward": {
				"denom": "ucmdx",
				"amount": "19392"
			}
		},
		{
			"address": "comdex1fs945ssvyq449wgpg8a2csgsp8qvu8pfcnngss",
			"reward": {
				"denom": "ucmdx",
				"amount": "169"
			}
		},
		{
			"address": "comdex1fsxxtzc7g8srase4w8lsl5sp5r4f2xelk95n4m",
			"reward": {
				"denom": "ucmdx",
				"amount": "2054"
			}
		},
		{
			"address": "comdex1fs89lh2y8glap0u55nn5jy3mqzxsgfj4pv3x9y",
			"reward": {
				"denom": "ucmdx",
				"amount": "2014"
			}
		},
		{
			"address": "comdex1fs8m7dz6jnzvyujrns4l0gd4u47lny99qpfzdn",
			"reward": {
				"denom": "ucmdx",
				"amount": "1970"
			}
		},
		{
			"address": "comdex1fstg6g2qahdvr6lmq0ldrkrnq2ss9c4qynpyju",
			"reward": {
				"denom": "ucmdx",
				"amount": "7167"
			}
		},
		{
			"address": "comdex1fsdpx8k5kaktxhcx47gdgcnnslv3r9xkvwmm25",
			"reward": {
				"denom": "ucmdx",
				"amount": "1899"
			}
		},
		{
			"address": "comdex1fs0kn7rfgvcv9htwxhze58fvymqfyvvmsqvgg7",
			"reward": {
				"denom": "ucmdx",
				"amount": "8855"
			}
		},
		{
			"address": "comdex1fsjpq8ef9g5txgezz7slqp39y7ltcxs8szh40y",
			"reward": {
				"denom": "ucmdx",
				"amount": "438047"
			}
		},
		{
			"address": "comdex1fsj28pxj9wm8rr8ua2pmy7p28lumx3qquc4wc2",
			"reward": {
				"denom": "ucmdx",
				"amount": "6359"
			}
		},
		{
			"address": "comdex1fshxx89xamu2t3u2s7lrut9etzgcwy8h3v49fd",
			"reward": {
				"denom": "ucmdx",
				"amount": "187"
			}
		},
		{
			"address": "comdex1fsczk3dj3eqp62qeeg4jcwvg939zwa67te5cve",
			"reward": {
				"denom": "ucmdx",
				"amount": "57574"
			}
		},
		{
			"address": "comdex1fsclv2hq9vf2jw8dv5wrffprh5ga623mfh9nhp",
			"reward": {
				"denom": "ucmdx",
				"amount": "13303"
			}
		},
		{
			"address": "comdex1fsehze0al6cf38q8njdwv82t4aruawnxhpnmml",
			"reward": {
				"denom": "ucmdx",
				"amount": "25569"
			}
		},
		{
			"address": "comdex1fs67g6a7jdf4u6lkhyuwwlyyf9dw8ljtd02nj7",
			"reward": {
				"denom": "ucmdx",
				"amount": "65506"
			}
		},
		{
			"address": "comdex1fsms30sr0ywkrc6djh630rerptagfucwzfy7k6",
			"reward": {
				"denom": "ucmdx",
				"amount": "367"
			}
		},
		{
			"address": "comdex1f3q9s9r52cd5rxvtpna8fqqzhp9ajr6jfs5gqd",
			"reward": {
				"denom": "ucmdx",
				"amount": "148"
			}
		},
		{
			"address": "comdex1f3pp8ejl5fqtf3utucz8e8a933gge0x9nkf5eg",
			"reward": {
				"denom": "ucmdx",
				"amount": "56637"
			}
		},
		{
			"address": "comdex1f3pg6m6ur33v6ykgrszqt456lk8fmaqzs4766q",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1f3zww5dlt5lmym3cdgmrk3lfg0q97nuxplgv7x",
			"reward": {
				"denom": "ucmdx",
				"amount": "14580"
			}
		},
		{
			"address": "comdex1f3z4u44ds4wp88ca9g24yshe6tyyvpmrhm87vx",
			"reward": {
				"denom": "ucmdx",
				"amount": "151"
			}
		},
		{
			"address": "comdex1f3rmay2df9mfx2ct3c3e0p3t7rkrxqs628dpk2",
			"reward": {
				"denom": "ucmdx",
				"amount": "12425"
			}
		},
		{
			"address": "comdex1f39ux5vduek66zstqr6pzzx3sxvd0zx5nhpaht",
			"reward": {
				"denom": "ucmdx",
				"amount": "8966"
			}
		},
		{
			"address": "comdex1f3xefzvz9znpal39prwvz23e5am23d9xqam0r0",
			"reward": {
				"denom": "ucmdx",
				"amount": "6824"
			}
		},
		{
			"address": "comdex1f38sqwx5llttc2j3e3jj8zd92hl4c8xx5gcqyv",
			"reward": {
				"denom": "ucmdx",
				"amount": "11690"
			}
		},
		{
			"address": "comdex1f3gxw0vqtk2gk4scd220g0ex78u36cqauhx20j",
			"reward": {
				"denom": "ucmdx",
				"amount": "7883"
			}
		},
		{
			"address": "comdex1f3thrzymraht66ggkm82tup3ztcpkqga6j6frl",
			"reward": {
				"denom": "ucmdx",
				"amount": "179"
			}
		},
		{
			"address": "comdex1f3tlzjx5rgxkp535kmfxs8xxg7ddg2gzapuqtg",
			"reward": {
				"denom": "ucmdx",
				"amount": "142"
			}
		},
		{
			"address": "comdex1f3w37sk8dss7jj9vqqnysmcnfmmyt3mtnnt6s0",
			"reward": {
				"denom": "ucmdx",
				"amount": "5931"
			}
		},
		{
			"address": "comdex1f30gmjqxnxyentfdg40fcnfppn55x5ffhr3ac7",
			"reward": {
				"denom": "ucmdx",
				"amount": "64591"
			}
		},
		{
			"address": "comdex1f33lfsu54jdter9y09rhlncedawpfun93x6dv0",
			"reward": {
				"denom": "ucmdx",
				"amount": "2905"
			}
		},
		{
			"address": "comdex1f3jzwlferps5qa2ghp663azfk263gyvqxzg8zc",
			"reward": {
				"denom": "ucmdx",
				"amount": "705"
			}
		},
		{
			"address": "comdex1f3ndu0gyjzctcsqm5mfzgu7hejd86pcpmfd65z",
			"reward": {
				"denom": "ucmdx",
				"amount": "52"
			}
		},
		{
			"address": "comdex1f34g6wqdxf7fym9cyltcvnhkhhd6mwy5ds4k3a",
			"reward": {
				"denom": "ucmdx",
				"amount": "145280"
			}
		},
		{
			"address": "comdex1f34d88rdaqsqqvet6texgugt3vd0nsgzpg976y",
			"reward": {
				"denom": "ucmdx",
				"amount": "57"
			}
		},
		{
			"address": "comdex1f3m8cmymdy4wqcnan2uf9kc5jsexwhdryncw37",
			"reward": {
				"denom": "ucmdx",
				"amount": "3759"
			}
		},
		{
			"address": "comdex1f3um6rd7tzuh3gu77mt7ru9rk5020862qac909",
			"reward": {
				"denom": "ucmdx",
				"amount": "9066"
			}
		},
		{
			"address": "comdex1fjrpuenrvkq8j5zyrmzrwv62g5pw9xz55sec2r",
			"reward": {
				"denom": "ucmdx",
				"amount": "149"
			}
		},
		{
			"address": "comdex1fjr37hkuw5l6wfdke767mfjpef9hxfs9amcdj0",
			"reward": {
				"denom": "ucmdx",
				"amount": "2049428"
			}
		},
		{
			"address": "comdex1fjrjdpw38khuc0nyf02vu78dtyle70u3vq602s",
			"reward": {
				"denom": "ucmdx",
				"amount": "27951"
			}
		},
		{
			"address": "comdex1fjxeyrzc428z6gwpu2gpt0vs4q5hve5s4l2rm9",
			"reward": {
				"denom": "ucmdx",
				"amount": "466"
			}
		},
		{
			"address": "comdex1fj8reqvck66jme4keh0naxfas63x340u9awku3",
			"reward": {
				"denom": "ucmdx",
				"amount": "90"
			}
		},
		{
			"address": "comdex1fj3h3c5ncu6qphm3an6y92ulre2eg27u8j0fu0",
			"reward": {
				"denom": "ucmdx",
				"amount": "1996"
			}
		},
		{
			"address": "comdex1fjjps3u3t4swuyu74epfa2l3sn04fuquvvrl8z",
			"reward": {
				"denom": "ucmdx",
				"amount": "66976"
			}
		},
		{
			"address": "comdex1fjhtga2hg3tke378fn772r8zyt08wvsjzmckn6",
			"reward": {
				"denom": "ucmdx",
				"amount": "2054"
			}
		},
		{
			"address": "comdex1fjh798fyn9vqldc8ru09c99r89v5mt3cdgyqxy",
			"reward": {
				"denom": "ucmdx",
				"amount": "62073"
			}
		},
		{
			"address": "comdex1fje95rhe9te60pmvwqt436kf8my6ty6cd2g796",
			"reward": {
				"denom": "ucmdx",
				"amount": "168"
			}
		},
		{
			"address": "comdex1fj6j8fv0rdhy3m428mzj9ca48m8j46u6twn3fx",
			"reward": {
				"denom": "ucmdx",
				"amount": "2930"
			}
		},
		{
			"address": "comdex1fjm36lrdueszd94374y940vs9yz7twlcd368wx",
			"reward": {
				"denom": "ucmdx",
				"amount": "1686"
			}
		},
		{
			"address": "comdex1fj73m8k4f0kt45a2lda9368v2cvn8c0fxm27y8",
			"reward": {
				"denom": "ucmdx",
				"amount": "1469"
			}
		},
		{
			"address": "comdex1fjlgzd6lxgsj9yyq9ud329a6seh4dqezrtxm22",
			"reward": {
				"denom": "ucmdx",
				"amount": "12367"
			}
		},
		{
			"address": "comdex1fn83fg7f0m4xxlv9e6ljk3m7f44caujh6xq5u2",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1fn85ljcu4gdql3t2mg2yn2e4aw0483qnc8zgsa",
			"reward": {
				"denom": "ucmdx",
				"amount": "6838"
			}
		},
		{
			"address": "comdex1fn87m2aylzphndf9nnweq8vmu64k90angrfg2m",
			"reward": {
				"denom": "ucmdx",
				"amount": "1740"
			}
		},
		{
			"address": "comdex1fngj4qzze0stezv0uecrxkgnjtqf6x8kjchzsf",
			"reward": {
				"denom": "ucmdx",
				"amount": "92"
			}
		},
		{
			"address": "comdex1fnf668gjrkwk7x9pvnr2q0t2nznt75xtlwqd2g",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1fn2hfehwn54h0s43k5e4yzevkdgvskyypzsf77",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1fntvksgf7cl25h3gx77m0n360s2x9wu3djmyh2",
			"reward": {
				"denom": "ucmdx",
				"amount": "2512"
			}
		},
		{
			"address": "comdex1fntv79pzyazld30welfyus6m8adjylf659r6hv",
			"reward": {
				"denom": "ucmdx",
				"amount": "137"
			}
		},
		{
			"address": "comdex1fndt9p7lqgchzl7v447m6dwedmh3p708taf7eu",
			"reward": {
				"denom": "ucmdx",
				"amount": "41"
			}
		},
		{
			"address": "comdex1fnwxrt3gr25zhq7ek5j67qsz3qsrtadr3sthzl",
			"reward": {
				"denom": "ucmdx",
				"amount": "141"
			}
		},
		{
			"address": "comdex1fn0yx4wz50y9qyx782l2p46l3atk5fmdeqwsaa",
			"reward": {
				"denom": "ucmdx",
				"amount": "1410"
			}
		},
		{
			"address": "comdex1fns7zsn2l0tanhmgrxr7p3ep00pmfja5zhhyz6",
			"reward": {
				"denom": "ucmdx",
				"amount": "7152"
			}
		},
		{
			"address": "comdex1fnnpt069wks84d2afhvj8l6xk2jemzzav0n4en",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1fnn7pqramw6hvmy3g8e6krx9jenmmpsavgzyd8",
			"reward": {
				"denom": "ucmdx",
				"amount": "115"
			}
		},
		{
			"address": "comdex1fn5fkex9fsydfdcd3h26gn95j4cwytrqtyyz79",
			"reward": {
				"denom": "ucmdx",
				"amount": "1406"
			}
		},
		{
			"address": "comdex1fn4enehw2t2myewc6uhyz8n8pj2g2l6w9de8n9",
			"reward": {
				"denom": "ucmdx",
				"amount": "389"
			}
		},
		{
			"address": "comdex1fnerv4upel02mpmsxmfz2zhtds8p70zwfc4ntm",
			"reward": {
				"denom": "ucmdx",
				"amount": "1690"
			}
		},
		{
			"address": "comdex1fnmeud9camgvl09fg7nxfkf73r3zp44mlqg504",
			"reward": {
				"denom": "ucmdx",
				"amount": "144"
			}
		},
		{
			"address": "comdex1f5z2ypdztnj5nj7ayprj7hc2knw5z68vgufhlc",
			"reward": {
				"denom": "ucmdx",
				"amount": "1352"
			}
		},
		{
			"address": "comdex1f5243uwxcdqeq8l3agejlyjxrrk2tevxze47e6",
			"reward": {
				"denom": "ucmdx",
				"amount": "73"
			}
		},
		{
			"address": "comdex1f5vgyqx0hvpngvx2dqpeqeer4wgzhld2q8psw6",
			"reward": {
				"denom": "ucmdx",
				"amount": "442"
			}
		},
		{
			"address": "comdex1f5dshmmjynfnlfwjew4g0zdxr37z6r8pzh66ku",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1f5d6kf6jpngj3pgx04axll3cxvvjtqkxpq8pcc",
			"reward": {
				"denom": "ucmdx",
				"amount": "121"
			}
		},
		{
			"address": "comdex1f5jr64takt6rdl4aa6smygafmcvfgts6a2xayx",
			"reward": {
				"denom": "ucmdx",
				"amount": "440132"
			}
		},
		{
			"address": "comdex1f5nvnmwzpf05zymm89wck45gc4yqheyvc3atgz",
			"reward": {
				"denom": "ucmdx",
				"amount": "71"
			}
		},
		{
			"address": "comdex1f5kvg5qekwqq9e7kpk4k8mqcl94d3eyrhfwcje",
			"reward": {
				"denom": "ucmdx",
				"amount": "6835"
			}
		},
		{
			"address": "comdex1f5u7ktfzwdtx58rlnlma5plctcxvk34ww4xwm0",
			"reward": {
				"denom": "ucmdx",
				"amount": "1832"
			}
		},
		{
			"address": "comdex1f5aa0xaz4jkwun4pdxadkf98sdc4fhk0j29q3e",
			"reward": {
				"denom": "ucmdx",
				"amount": "6665"
			}
		},
		{
			"address": "comdex1f4pyghxn647qvaqwcj6adcywetu8tn6z3w6j42",
			"reward": {
				"denom": "ucmdx",
				"amount": "16296"
			}
		},
		{
			"address": "comdex1f4pwqhgzkt6fk70q4tmmlcyexk3nxtgvkspyx8",
			"reward": {
				"denom": "ucmdx",
				"amount": "40557"
			}
		},
		{
			"address": "comdex1f4zjc0zt0dudes8rlnzrsxxw4uvd2s86aljcg4",
			"reward": {
				"denom": "ucmdx",
				"amount": "7877"
			}
		},
		{
			"address": "comdex1f49j7dpljlq63mv4cywpv8p66jx7kun2ugkts2",
			"reward": {
				"denom": "ucmdx",
				"amount": "398"
			}
		},
		{
			"address": "comdex1f497lfypmszxd4dnfqc46ude5v3sercryf62hm",
			"reward": {
				"denom": "ucmdx",
				"amount": "90253"
			}
		},
		{
			"address": "comdex1f4xv8gjynhnl7sn5gw4zpa6uqycgdw4ptlveaj",
			"reward": {
				"denom": "ucmdx",
				"amount": "1780"
			}
		},
		{
			"address": "comdex1f4gpmgmh2jhp2jvezpe4mlevavflhecavdrvnh",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1f4ffzycapwnvyy23p73djcnxa668eydhwmqfna",
			"reward": {
				"denom": "ucmdx",
				"amount": "176"
			}
		},
		{
			"address": "comdex1f4fu6mpf663lh4q9rvtfld6k7xdzy077jjj4kf",
			"reward": {
				"denom": "ucmdx",
				"amount": "14419"
			}
		},
		{
			"address": "comdex1f4tpr5s83xflx95er7pspm0qe48kpajvh9wm4x",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1f4tae57tz9qqqdxm94ujsdzumant76ca8xd8c7",
			"reward": {
				"denom": "ucmdx",
				"amount": "1505"
			}
		},
		{
			"address": "comdex1f4dkghkdn6vts32kf2jed5d3uevw6angeyzynf",
			"reward": {
				"denom": "ucmdx",
				"amount": "15394"
			}
		},
		{
			"address": "comdex1f4w0hhfm6c6e9q80f3u8zf4fht2vm54v65q8vr",
			"reward": {
				"denom": "ucmdx",
				"amount": "24993"
			}
		},
		{
			"address": "comdex1f4j5sl4dafrw4r9u43jg7nr4k9lmm4f73n7p09",
			"reward": {
				"denom": "ucmdx",
				"amount": "2462"
			}
		},
		{
			"address": "comdex1f4jmpk8qqehxxqg6pze5ytgg4xj68yeylccgzv",
			"reward": {
				"denom": "ucmdx",
				"amount": "355"
			}
		},
		{
			"address": "comdex1f4nnees22jlm0p0rgvya7l6eq7gjmxf84lg7le",
			"reward": {
				"denom": "ucmdx",
				"amount": "717"
			}
		},
		{
			"address": "comdex1f4nusr3a4pfkgc4saw5yxe55405tmpn37srlyj",
			"reward": {
				"denom": "ucmdx",
				"amount": "901"
			}
		},
		{
			"address": "comdex1f444aetv2l3w2wscd7pkrrgg47c90zz3gge6ah",
			"reward": {
				"denom": "ucmdx",
				"amount": "3542"
			}
		},
		{
			"address": "comdex1f4k7v9wgrum243j9eqc9vurtxh8ugskna3kvrf",
			"reward": {
				"denom": "ucmdx",
				"amount": "140"
			}
		},
		{
			"address": "comdex1f4c2p8lfjtfxa9m50gpj22665xg45hpm0sfxr0",
			"reward": {
				"denom": "ucmdx",
				"amount": "1931"
			}
		},
		{
			"address": "comdex1f4cc9mu7gc0r4s9nkum4wfsv0dsdyzf9y8gv0g",
			"reward": {
				"denom": "ucmdx",
				"amount": "10582"
			}
		},
		{
			"address": "comdex1f4ed72jjupaqawhzmn5w3rsuyccu7d3zec6rqn",
			"reward": {
				"denom": "ucmdx",
				"amount": "7161"
			}
		},
		{
			"address": "comdex1f4lra52avyw35wqetv6pemcn70eqh5pmq0ltmj",
			"reward": {
				"denom": "ucmdx",
				"amount": "28840"
			}
		},
		{
			"address": "comdex1f4ly5clf93tu84g2tnm0fualwdavp5hqutcyhg",
			"reward": {
				"denom": "ucmdx",
				"amount": "6213"
			}
		},
		{
			"address": "comdex1fky45lvsttf3n4x40hseyuav0wu4vscuwapns7",
			"reward": {
				"denom": "ucmdx",
				"amount": "193"
			}
		},
		{
			"address": "comdex1fk95hrxtxz48ncx9ak5zgqna8zp5f6eym3cfau",
			"reward": {
				"denom": "ucmdx",
				"amount": "1789"
			}
		},
		{
			"address": "comdex1fk29y2hpd79zp5eg25q2nhavk97tqnrp4q45ws",
			"reward": {
				"denom": "ucmdx",
				"amount": "7242"
			}
		},
		{
			"address": "comdex1fkv55duypw3enzu8y4urd2ltmg8mlk4znmw73v",
			"reward": {
				"denom": "ucmdx",
				"amount": "6797"
			}
		},
		{
			"address": "comdex1fkw5n6upztctlkykf5vsgkn0krvtfaw6swjg2a",
			"reward": {
				"denom": "ucmdx",
				"amount": "5472"
			}
		},
		{
			"address": "comdex1fk0cmu5kgqafaq7srcv57c4yvsjnne6jtc68jd",
			"reward": {
				"denom": "ucmdx",
				"amount": "14330"
			}
		},
		{
			"address": "comdex1fknffatm6zuxrnrvqvj8s06uht2vr9mxwyyjc4",
			"reward": {
				"denom": "ucmdx",
				"amount": "14403"
			}
		},
		{
			"address": "comdex1fknhdqr0q6h98nn28gymamdz04ehks3v0nuaf5",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1fk5pnpd9y52zwaf254tvpqj2dcljtvsnwhkmmr",
			"reward": {
				"denom": "ucmdx",
				"amount": "6155"
			}
		},
		{
			"address": "comdex1fkh53jzjltugmxc5y50pm727a85nscrnp7z9s9",
			"reward": {
				"denom": "ucmdx",
				"amount": "1346"
			}
		},
		{
			"address": "comdex1fkcxc7e9zyvq3e9245uh4er33fkntx0askl5ep",
			"reward": {
				"denom": "ucmdx",
				"amount": "430"
			}
		},
		{
			"address": "comdex1fkcd8dtx7yemjdpxf2e6tr2y0z7nmzaqmj7d2a",
			"reward": {
				"denom": "ucmdx",
				"amount": "73777"
			}
		},
		{
			"address": "comdex1fk6g5dk6yljq5jlgtszyguw6erqt5z0s8ewjwd",
			"reward": {
				"denom": "ucmdx",
				"amount": "97"
			}
		},
		{
			"address": "comdex1fku9gl93dy3z4d2y58gza06un72ulmd8vvp7le",
			"reward": {
				"denom": "ucmdx",
				"amount": "1227"
			}
		},
		{
			"address": "comdex1fhz6yklrgv79fq6stensp94p0ykqr0uj4uwyw0",
			"reward": {
				"denom": "ucmdx",
				"amount": "5838"
			}
		},
		{
			"address": "comdex1fhx3scs5czaf22juc6mjas5ellkx5shpxrhzls",
			"reward": {
				"denom": "ucmdx",
				"amount": "10038"
			}
		},
		{
			"address": "comdex1fhxnyzpn0w3jcv6wwmxzz6fs0ep7nk57pq0mc3",
			"reward": {
				"denom": "ucmdx",
				"amount": "271"
			}
		},
		{
			"address": "comdex1fh2c6rvmdt5gahaeqyhf7dc2lxxmcu8kwsp45p",
			"reward": {
				"denom": "ucmdx",
				"amount": "56726"
			}
		},
		{
			"address": "comdex1fhd5nzm0r0j65azpmt69dv04zmznwykmtsv245",
			"reward": {
				"denom": "ucmdx",
				"amount": "8298"
			}
		},
		{
			"address": "comdex1fhwcyxx2lgefnslf0prm53k93qxsef6ursm5f4",
			"reward": {
				"denom": "ucmdx",
				"amount": "4926"
			}
		},
		{
			"address": "comdex1fhnc5vtxeclsqpkk8puyx4vmcuxypps0nzc6pm",
			"reward": {
				"denom": "ucmdx",
				"amount": "280768"
			}
		},
		{
			"address": "comdex1fhkga75da3j3fd0er3m76gj9uwtqqngk5feky0",
			"reward": {
				"denom": "ucmdx",
				"amount": "141"
			}
		},
		{
			"address": "comdex1fhk3snnmqaxjnrrt0xjhwfrrd4y88nqgfx093z",
			"reward": {
				"denom": "ucmdx",
				"amount": "1199"
			}
		},
		{
			"address": "comdex1fhh9p8cs6lej5wsjrk2066ep58xt73ksdfu20r",
			"reward": {
				"denom": "ucmdx",
				"amount": "1158"
			}
		},
		{
			"address": "comdex1fhhl6rh3x7j68ldcgpmxwdz2um4r3ede2pg5s7",
			"reward": {
				"denom": "ucmdx",
				"amount": "3131"
			}
		},
		{
			"address": "comdex1fhce0cf6ct0fd96s9g7y4lq28606sqvqsrujzr",
			"reward": {
				"denom": "ucmdx",
				"amount": "1250"
			}
		},
		{
			"address": "comdex1fhe0rmrdzpnw78s5uej9znvjkskcs3m3e324fp",
			"reward": {
				"denom": "ucmdx",
				"amount": "7187"
			}
		},
		{
			"address": "comdex1fhmvtp2yh7n740xp8rut74sar4u3jrerlyh9kl",
			"reward": {
				"denom": "ucmdx",
				"amount": "62401"
			}
		},
		{
			"address": "comdex1fh77s9hw3rxjsqgszymf92uvzlg67945lf9nrz",
			"reward": {
				"denom": "ucmdx",
				"amount": "35208"
			}
		},
		{
			"address": "comdex1fcqz87qknp66tgqkhl8gadjscta7x3d9947gdr",
			"reward": {
				"denom": "ucmdx",
				"amount": "197203"
			}
		},
		{
			"address": "comdex1fcqt0sl2sr3dq7euu6dduhv3at0q9nupnrvngx",
			"reward": {
				"denom": "ucmdx",
				"amount": "1985"
			}
		},
		{
			"address": "comdex1fcr0xups4l967mtrg8euvrd0v28fprta3vqw8q",
			"reward": {
				"denom": "ucmdx",
				"amount": "148"
			}
		},
		{
			"address": "comdex1fcrsh7rt9exfd3tyt6y2m49x970tya2am9auds",
			"reward": {
				"denom": "ucmdx",
				"amount": "6817"
			}
		},
		{
			"address": "comdex1fcygwdcfwdnt8gpuv0rr5j6llnexsw66amvs98",
			"reward": {
				"denom": "ucmdx",
				"amount": "353"
			}
		},
		{
			"address": "comdex1fcyac5lr5djc9623d0en9r9m0a34da26j4p24z",
			"reward": {
				"denom": "ucmdx",
				"amount": "35002"
			}
		},
		{
			"address": "comdex1fcxywvkt3t5y34vayjy7r9zyqnczu06sn2tudr",
			"reward": {
				"denom": "ucmdx",
				"amount": "29"
			}
		},
		{
			"address": "comdex1fcxac3wxcfxfmp0z577g5awrnyq6xl959dej77",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1fcfus4u2zt8nr4srvm0g9xzh39r6m6qv6sqjkp",
			"reward": {
				"denom": "ucmdx",
				"amount": "28302"
			}
		},
		{
			"address": "comdex1fc2w5avsmcd4xk4rxquhv3rzdxz6kkhwqkglge",
			"reward": {
				"denom": "ucmdx",
				"amount": "30288"
			}
		},
		{
			"address": "comdex1fcwuuev5nnu3jeckykjnn6tp0vre3uu7l68afu",
			"reward": {
				"denom": "ucmdx",
				"amount": "15149"
			}
		},
		{
			"address": "comdex1fc5p3n0mlm2tk9dar55n9astmxs4p5jkgu70yg",
			"reward": {
				"denom": "ucmdx",
				"amount": "129"
			}
		},
		{
			"address": "comdex1fcamvvfm3s6hvvnvkj8tkteacdyg3xce247mq5",
			"reward": {
				"denom": "ucmdx",
				"amount": "6887"
			}
		},
		{
			"address": "comdex1feqtd56sy2c4829sw4jr2dlkc0lgma6f64epxg",
			"reward": {
				"denom": "ucmdx",
				"amount": "122354"
			}
		},
		{
			"address": "comdex1feqk6ugwyx0pt38c9m6xh3exq6u58sdqw4tnkl",
			"reward": {
				"denom": "ucmdx",
				"amount": "143"
			}
		},
		{
			"address": "comdex1fer9yxd3ec9qftvshw8ldsjwgdrsdu0cl2g2qw",
			"reward": {
				"denom": "ucmdx",
				"amount": "7093"
			}
		},
		{
			"address": "comdex1ferfl4vl894w54yrpa5kcr8xeg3z2ye5te0n38",
			"reward": {
				"denom": "ucmdx",
				"amount": "301"
			}
		},
		{
			"address": "comdex1fexgfg4vxy3dlds0dtw0tchkv2cp9l5mkfnpje",
			"reward": {
				"denom": "ucmdx",
				"amount": "17650"
			}
		},
		{
			"address": "comdex1fexcy9jrpplcuqwavt0aj6at8cxjgdc7k0a5s7",
			"reward": {
				"denom": "ucmdx",
				"amount": "12563"
			}
		},
		{
			"address": "comdex1fe8sj6ldpjg5yukwk6s7mhg5yqdhd2ktkhhx66",
			"reward": {
				"denom": "ucmdx",
				"amount": "7063"
			}
		},
		{
			"address": "comdex1fegrduzj7xl8hq6w7uv2kvdne47w98jrcaq66r",
			"reward": {
				"denom": "ucmdx",
				"amount": "301"
			}
		},
		{
			"address": "comdex1fetmzd77a4j2rqysepsmksk782j06s4m49greh",
			"reward": {
				"denom": "ucmdx",
				"amount": "19431"
			}
		},
		{
			"address": "comdex1fewcdj7mgurq90n7drld0dfg8xc6f9jx6d2rtg",
			"reward": {
				"denom": "ucmdx",
				"amount": "12637"
			}
		},
		{
			"address": "comdex1fe0m4s7d2dyje4ene0xsmyrxmymsnu4yjnm9ly",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1fe3ddd3fg99tpt6ex86u4h42uxc9kldxm6twed",
			"reward": {
				"denom": "ucmdx",
				"amount": "1765"
			}
		},
		{
			"address": "comdex1fenp5wf2x0l7usvca990m97echgcmgcs862wjv",
			"reward": {
				"denom": "ucmdx",
				"amount": "19930"
			}
		},
		{
			"address": "comdex1fe453vnvzmf5ktu3tyq7kd68gymcfzuuj0zy5d",
			"reward": {
				"denom": "ucmdx",
				"amount": "1290"
			}
		},
		{
			"address": "comdex1fek89wmzgnfmsmmhfhs8n6mffgkcsg444ntr6d",
			"reward": {
				"denom": "ucmdx",
				"amount": "341"
			}
		},
		{
			"address": "comdex1fek2ezdctakw3y9g69xhtjpxx8078eyslmlr3d",
			"reward": {
				"denom": "ucmdx",
				"amount": "61613"
			}
		},
		{
			"address": "comdex1femh4tsz8c7snqnjah0mdrrlskxfs60tqsswdx",
			"reward": {
				"denom": "ucmdx",
				"amount": "14325"
			}
		},
		{
			"address": "comdex1feu0ca2ed8y268jj3nq46uz4wamfetf5sdsyaj",
			"reward": {
				"denom": "ucmdx",
				"amount": "6870"
			}
		},
		{
			"address": "comdex1fea2dr4uapzkhnwdrg9mfkhnn80q4e878xpuw6",
			"reward": {
				"denom": "ucmdx",
				"amount": "9898"
			}
		},
		{
			"address": "comdex1f6zmpc4m9f80n3lznhvrdkcxcqvke3w0288y6p",
			"reward": {
				"denom": "ucmdx",
				"amount": "78302"
			}
		},
		{
			"address": "comdex1f6tdca8x3nwer4ufwtzgnxl2p0c58cyl4m2kpk",
			"reward": {
				"denom": "ucmdx",
				"amount": "45222"
			}
		},
		{
			"address": "comdex1f6vpsgwwnpanyd6x6jkx6jzmw5lrv7dh9t8ajj",
			"reward": {
				"denom": "ucmdx",
				"amount": "167"
			}
		},
		{
			"address": "comdex1f6df3njeuruzuss2eyg2mdkhkr2m02cu5es7uw",
			"reward": {
				"denom": "ucmdx",
				"amount": "10719"
			}
		},
		{
			"address": "comdex1f6wtm6399g0pzpf5mkva4t0rxfh0a7dn84eep5",
			"reward": {
				"denom": "ucmdx",
				"amount": "2011"
			}
		},
		{
			"address": "comdex1f6s3m5vny8t29qphjsf7z6yg76gn62uve00cwp",
			"reward": {
				"denom": "ucmdx",
				"amount": "166"
			}
		},
		{
			"address": "comdex1f6jlesjkzx42nvg2w2jfm3ypdl5c88rl5uarsv",
			"reward": {
				"denom": "ucmdx",
				"amount": "11313"
			}
		},
		{
			"address": "comdex1f6n9xf4qk6hhrlds358vgyfmlx26j7wmu8p8uj",
			"reward": {
				"denom": "ucmdx",
				"amount": "2672"
			}
		},
		{
			"address": "comdex1f6502e5q5s40pgtvhgc209s2x2q0c9ur8l39t0",
			"reward": {
				"denom": "ucmdx",
				"amount": "81"
			}
		},
		{
			"address": "comdex1f655mfwmg9e307n3pqfrdln05qty2d0edns5j5",
			"reward": {
				"denom": "ucmdx",
				"amount": "2172"
			}
		},
		{
			"address": "comdex1f6klmvl7ftqntv2mpcvnhujf3vs2xl2zvzrtj3",
			"reward": {
				"denom": "ucmdx",
				"amount": "12364"
			}
		},
		{
			"address": "comdex1f6u85enm6dfq8y7rqyfq2550eztcjh9uvztls4",
			"reward": {
				"denom": "ucmdx",
				"amount": "35421"
			}
		},
		{
			"address": "comdex1f6aml5wwj0qny0m57gjentwhujp8zhwy4pt5kf",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1f6aug5r6fflyxp3sd34f6ejklv6n0fc23sz2u2",
			"reward": {
				"denom": "ucmdx",
				"amount": "166"
			}
		},
		{
			"address": "comdex1f674e3sqxm6k4fj7rp4l328xm5dvylh947fg0v",
			"reward": {
				"denom": "ucmdx",
				"amount": "15155"
			}
		},
		{
			"address": "comdex1f6lfrehfsrtj2zz4tyg2n85z49wydle2y0ltx4",
			"reward": {
				"denom": "ucmdx",
				"amount": "7104"
			}
		},
		{
			"address": "comdex1fmpjpqryddaqte5tuv8vnk20xlr9x4syvuqn4v",
			"reward": {
				"denom": "ucmdx",
				"amount": "12748"
			}
		},
		{
			"address": "comdex1fmglx3geseuzkdw0ctj784nlwf0m64k0qerj6e",
			"reward": {
				"denom": "ucmdx",
				"amount": "2052"
			}
		},
		{
			"address": "comdex1fmf7yrtqsjncnfmwgxaztwdhpv9uhpgsgk08f4",
			"reward": {
				"denom": "ucmdx",
				"amount": "180"
			}
		},
		{
			"address": "comdex1fm2376w8560wjn0c97nn92tctquzyepnsfqstf",
			"reward": {
				"denom": "ucmdx",
				"amount": "1014"
			}
		},
		{
			"address": "comdex1fmtvfxezpkt6n7tk34t9h8m8j46fztfyxwzmsl",
			"reward": {
				"denom": "ucmdx",
				"amount": "9885"
			}
		},
		{
			"address": "comdex1fm38dlya37z3u2c5mqyj3c8d6kc5ghy35r5436",
			"reward": {
				"denom": "ucmdx",
				"amount": "1"
			}
		},
		{
			"address": "comdex1fm3utlrj0yjftdt808uz4j4sjt09k0yx9y7pjc",
			"reward": {
				"denom": "ucmdx",
				"amount": "62460"
			}
		},
		{
			"address": "comdex1fmjlzv7h5s9v540r6j2vuy7fw3de9hnt6n8wt2",
			"reward": {
				"denom": "ucmdx",
				"amount": "478"
			}
		},
		{
			"address": "comdex1fmh6ger5d7skj8wj7nwae3eh2vnj4rl2knu85f",
			"reward": {
				"denom": "ucmdx",
				"amount": "14524"
			}
		},
		{
			"address": "comdex1fmc5nw9mtgs8xn42u5jcp50s6y2pvl3rfxsexc",
			"reward": {
				"denom": "ucmdx",
				"amount": "1716"
			}
		},
		{
			"address": "comdex1fmmykyxr0nsnrje3jwteqsjqtk93wfgj3f6wkw",
			"reward": {
				"denom": "ucmdx",
				"amount": "246866"
			}
		},
		{
			"address": "comdex1fmmkme8sa5a9lvf5z4dqv8slk7je26cp2da250",
			"reward": {
				"denom": "ucmdx",
				"amount": "14207"
			}
		},
		{
			"address": "comdex1fmmm23qktlv3dee8zvawlus3skh6yp7jg9zxfe",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1fm7dl6n4yx7pasm3shp2t06c7xnmwl0vt9ua2x",
			"reward": {
				"denom": "ucmdx",
				"amount": "9512"
			}
		},
		{
			"address": "comdex1fmldk8fm6mnhv9kkppdwrh0vfmrz7tang5rk5w",
			"reward": {
				"denom": "ucmdx",
				"amount": "105"
			}
		},
		{
			"address": "comdex1fuxqh43my9zy7v0h03ckk23n48rflmh0avfkes",
			"reward": {
				"denom": "ucmdx",
				"amount": "2828"
			}
		},
		{
			"address": "comdex1fuxuj86rkf7ep2724jz60xtxsle7vtdkgwlfkh",
			"reward": {
				"denom": "ucmdx",
				"amount": "166"
			}
		},
		{
			"address": "comdex1fu2nh0c2lj5kz6yekaavsynceyk0kkz97na2fl",
			"reward": {
				"denom": "ucmdx",
				"amount": "830"
			}
		},
		{
			"address": "comdex1fuvcqj8uw8mkysptlfgjjshvn60jvmhr8rt0e0",
			"reward": {
				"denom": "ucmdx",
				"amount": "17477"
			}
		},
		{
			"address": "comdex1fuw09hfg2gpk5tzcueh0nx054k87jsepnezz9w",
			"reward": {
				"denom": "ucmdx",
				"amount": "18336"
			}
		},
		{
			"address": "comdex1fu0rwu868e0uqg47kr6wg789yz538n02jan3rm",
			"reward": {
				"denom": "ucmdx",
				"amount": "1524"
			}
		},
		{
			"address": "comdex1fu3c5fxww3uv3dlfk5yfq0dutgpk4l74mz8r57",
			"reward": {
				"denom": "ucmdx",
				"amount": "16031"
			}
		},
		{
			"address": "comdex1fujwh4m62l4ppjl4ezg5t5t38s8gxs3afzqsf8",
			"reward": {
				"denom": "ucmdx",
				"amount": "99"
			}
		},
		{
			"address": "comdex1fuk2urnwdv88hjkusdpxkpj9g953w7uj306s53",
			"reward": {
				"denom": "ucmdx",
				"amount": "90"
			}
		},
		{
			"address": "comdex1fue2x62cklfhy07x28hpg63u2j2tz7yyd4uxk6",
			"reward": {
				"denom": "ucmdx",
				"amount": "283"
			}
		},
		{
			"address": "comdex1fuurkyutptm5tv7htyv8xz04sst4ddh9mpnxyw",
			"reward": {
				"denom": "ucmdx",
				"amount": "6774"
			}
		},
		{
			"address": "comdex1fuujyw2ahw4knzstjre9v2u2026le7rlfvz2q4",
			"reward": {
				"denom": "ucmdx",
				"amount": "1727"
			}
		},
		{
			"address": "comdex1fuaa8erz92lma7x403jwffr7awsf5p45dw4yhx",
			"reward": {
				"denom": "ucmdx",
				"amount": "5037"
			}
		},
		{
			"address": "comdex1fapqgamx3c8za8d3z56sydjtjygyu0ly6gurqc",
			"reward": {
				"denom": "ucmdx",
				"amount": "34915"
			}
		},
		{
			"address": "comdex1fa9flyde09y3tmzajaq9qphwzf4krh8y8hfqjf",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1faxfaesrq5mfag3nahtzudtld6ctsudn49ddpa",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1fa88nzanw9uap5tlmn4mez4zgr2srg0hl6p8sw",
			"reward": {
				"denom": "ucmdx",
				"amount": "3514"
			}
		},
		{
			"address": "comdex1fagj78yrkzspc49azwrfg9jhnpfhuwc2eaq6el",
			"reward": {
				"denom": "ucmdx",
				"amount": "14961"
			}
		},
		{
			"address": "comdex1fa2lpqjxxadzpjychz63sq0z7s4vv3278du989",
			"reward": {
				"denom": "ucmdx",
				"amount": "2915"
			}
		},
		{
			"address": "comdex1fav9hrm33dw7hq5hm6fsjlx2flugfyws20v4ww",
			"reward": {
				"denom": "ucmdx",
				"amount": "5059"
			}
		},
		{
			"address": "comdex1fad5eyp5gpauxng5suw6z3qzy7was2uwj07gc3",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1fawxsw55x9wktsfsdlwhmtngq5t4xgmph5f5uf",
			"reward": {
				"denom": "ucmdx",
				"amount": "16994"
			}
		},
		{
			"address": "comdex1fa3f64un84txq0d86negges5l5eg8zavjecdjq",
			"reward": {
				"denom": "ucmdx",
				"amount": "604"
			}
		},
		{
			"address": "comdex1fa33lsekgud66uep8r7v0rgfrxd96u3apjdg5h",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1faj9uaqgu39dtt2awhw90yr7m8yl38dlrnzwe2",
			"reward": {
				"denom": "ucmdx",
				"amount": "343"
			}
		},
		{
			"address": "comdex1fanzkqmw7k4dv2gpeelfs7d7sg54r2zmav954n",
			"reward": {
				"denom": "ucmdx",
				"amount": "880"
			}
		},
		{
			"address": "comdex1fa5yu5jkjm4xdp9jmdx4gsrchsjtysrzyz7ur2",
			"reward": {
				"denom": "ucmdx",
				"amount": "10548"
			}
		},
		{
			"address": "comdex1fa59jlnx2aljljp7v2n9thc8ul9y74pwf99zvn",
			"reward": {
				"denom": "ucmdx",
				"amount": "1723"
			}
		},
		{
			"address": "comdex1fa5lht92g47anhu39c0jptwy2s2zery7kcy9vv",
			"reward": {
				"denom": "ucmdx",
				"amount": "7051"
			}
		},
		{
			"address": "comdex1fa5l76adc78hflllce2t4a4gujj60hfaa4mxx9",
			"reward": {
				"denom": "ucmdx",
				"amount": "204"
			}
		},
		{
			"address": "comdex1fac2uc5eeptn486nnu72ue2seaygmk9j22cyvn",
			"reward": {
				"denom": "ucmdx",
				"amount": "13016"
			}
		},
		{
			"address": "comdex1famxjp7f0csf9kzmjc3arnayfj8el9ystf7l0u",
			"reward": {
				"denom": "ucmdx",
				"amount": "16"
			}
		},
		{
			"address": "comdex1fauuc09464dkrl0mtxcldf0tctm2jheju9czkd",
			"reward": {
				"denom": "ucmdx",
				"amount": "43159"
			}
		},
		{
			"address": "comdex1faa270um7c4qn9d52v66cul4cjzg982pxdmyx6",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1falpfwjl90j6cj27uyzmjdqm8vsc2qu8jallwj",
			"reward": {
				"denom": "ucmdx",
				"amount": "9883"
			}
		},
		{
			"address": "comdex1f7qvypfaee56km5rp8lvay4f0saq9xn45kmzvk",
			"reward": {
				"denom": "ucmdx",
				"amount": "12718"
			}
		},
		{
			"address": "comdex1f7p8xr8vyladxexna8xe72j3wwz6ru7rmnsp3x",
			"reward": {
				"denom": "ucmdx",
				"amount": "453"
			}
		},
		{
			"address": "comdex1f7r6sja4x5hjvp6a7g9lypd94x2nvf9fqj8am8",
			"reward": {
				"denom": "ucmdx",
				"amount": "8844"
			}
		},
		{
			"address": "comdex1f79v6wqh5p6akq4pgp4gyft800usuxe8rvc3jk",
			"reward": {
				"denom": "ucmdx",
				"amount": "6973"
			}
		},
		{
			"address": "comdex1f72xaapvxy782vcxx9ydd59y4lxg809zlq6pqd",
			"reward": {
				"denom": "ucmdx",
				"amount": "147"
			}
		},
		{
			"address": "comdex1f7sa4nr5nhlt57vy5553tqzs5dlz7x6jp47yhr",
			"reward": {
				"denom": "ucmdx",
				"amount": "151571"
			}
		},
		{
			"address": "comdex1f7jt2p2qyr7kq76yn7l0ykvhfzgvcunukagnet",
			"reward": {
				"denom": "ucmdx",
				"amount": "166"
			}
		},
		{
			"address": "comdex1f7e4vgu6qfzl33rqr97dkkmg3wu3nqk33tllsw",
			"reward": {
				"denom": "ucmdx",
				"amount": "11"
			}
		},
		{
			"address": "comdex1f76z0jzt5ky0c9gttn9u28gv5yrgmfqw9frpmh",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1f76yagr9srhj3ufwwjvlfk8dh2cppe6m0st695",
			"reward": {
				"denom": "ucmdx",
				"amount": "58054"
			}
		},
		{
			"address": "comdex1f76n5et4z66cxm96e505kf3608fwzg004hsvat",
			"reward": {
				"denom": "ucmdx",
				"amount": "14231"
			}
		},
		{
			"address": "comdex1f7m4h7us7a2tgnrqyn32dx4nheg7guwqder8qj",
			"reward": {
				"denom": "ucmdx",
				"amount": "165"
			}
		},
		{
			"address": "comdex1f7a0ktuuu9tcqguzy4ydvhfuzzfv847nnryr7a",
			"reward": {
				"denom": "ucmdx",
				"amount": "83776"
			}
		},
		{
			"address": "comdex1f77j8ecfxgrg2nkezgxa7e5k8xvn0lr76rrd80",
			"reward": {
				"denom": "ucmdx",
				"amount": "18934"
			}
		},
		{
			"address": "comdex1f77lh80qqg96elvlzmrmgpau0u9gpqvd5ceusa",
			"reward": {
				"denom": "ucmdx",
				"amount": "89617"
			}
		},
		{
			"address": "comdex1fl9003j9gh9kj357qu7ye40vd568p6g202upa3",
			"reward": {
				"denom": "ucmdx",
				"amount": "173"
			}
		},
		{
			"address": "comdex1fl9uxwh97kjfyw4kmft3cq0w0j50jd3uq88uf9",
			"reward": {
				"denom": "ucmdx",
				"amount": "6126"
			}
		},
		{
			"address": "comdex1flxq5js537z5s3wcrqqr8wp2c5svad79w4zy4k",
			"reward": {
				"denom": "ucmdx",
				"amount": "1403"
			}
		},
		{
			"address": "comdex1fl89jfdx4pvud0qx2gtw8gf3k5mguxlhc0djla",
			"reward": {
				"denom": "ucmdx",
				"amount": "7583"
			}
		},
		{
			"address": "comdex1fl8jaw4uq2p95qu398xt2zlccv0cg3ezq2y37q",
			"reward": {
				"denom": "ucmdx",
				"amount": "215"
			}
		},
		{
			"address": "comdex1fl8uhk0q9c4cpzkz82g9mvwhqnwm78crh48w38",
			"reward": {
				"denom": "ucmdx",
				"amount": "292"
			}
		},
		{
			"address": "comdex1flgw2s7qw4772h0kl40uksrvc26k8w5n9crpg3",
			"reward": {
				"denom": "ucmdx",
				"amount": "6140"
			}
		},
		{
			"address": "comdex1fl2yjvfvjn59g5amucudmvstufq08t3442tw3p",
			"reward": {
				"denom": "ucmdx",
				"amount": "597"
			}
		},
		{
			"address": "comdex1fl2xa7k2wp6pmysur6r5gndgztfdr33c8v9hmd",
			"reward": {
				"denom": "ucmdx",
				"amount": "325406"
			}
		},
		{
			"address": "comdex1fldnr97u6g3ykg5alfcvfcme4e0xjgxfv83dr6",
			"reward": {
				"denom": "ucmdx",
				"amount": "527"
			}
		},
		{
			"address": "comdex1fl0pj4fpy6grcqej098aj8uh8w0nwscmraj9rw",
			"reward": {
				"denom": "ucmdx",
				"amount": "1611"
			}
		},
		{
			"address": "comdex1fl0l475262l36s22dulkm5mex9n37rxqwlttlp",
			"reward": {
				"denom": "ucmdx",
				"amount": "312788"
			}
		},
		{
			"address": "comdex1flswm9qxz09cljf2axw9a0zh7kkj544v9n6txk",
			"reward": {
				"denom": "ucmdx",
				"amount": "72080"
			}
		},
		{
			"address": "comdex1fl3lz3w0zs3tqgf0vuux8suuv440w5gr3nvkyv",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1fln4j0n8alk5yg4fdh3edlwga4pcgua3fdd05p",
			"reward": {
				"denom": "ucmdx",
				"amount": "6531"
			}
		},
		{
			"address": "comdex1fl54eu6pjk6ekwg3rrck7pxv6k3sssdgdvmlak",
			"reward": {
				"denom": "ucmdx",
				"amount": "1501"
			}
		},
		{
			"address": "comdex1flhwgftlterxmgtdxzfyupmnc7hhzsr9z02ct5",
			"reward": {
				"denom": "ucmdx",
				"amount": "271"
			}
		},
		{
			"address": "comdex1flealsydj29hkcpxfdjucjuu4257rq2dqrjsmx",
			"reward": {
				"denom": "ucmdx",
				"amount": "4688"
			}
		},
		{
			"address": "comdex1fl6yrf0sk80f4cg85svjss480z0jg9d7qaeule",
			"reward": {
				"denom": "ucmdx",
				"amount": "142"
			}
		},
		{
			"address": "comdex1fl6td57k7dhjsf5w8jfzqq4lmky95zuyamz6hy",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex12qqhw7fztp8xyp5jcs38ad8037fw9yrwkn82wa",
			"reward": {
				"denom": "ucmdx",
				"amount": "34540"
			}
		},
		{
			"address": "comdex12qt2tm4eh50yfnq2qlkpgzj2zassldkg4gnjm2",
			"reward": {
				"denom": "ucmdx",
				"amount": "1454"
			}
		},
		{
			"address": "comdex12qtnc3hc7a8upx0ld80xqyh8smxqc29j4zvu9j",
			"reward": {
				"denom": "ucmdx",
				"amount": "8783"
			}
		},
		{
			"address": "comdex12qvfy0nkxsps8mgh67cprvej40tss4x69crcc9",
			"reward": {
				"denom": "ucmdx",
				"amount": "37776"
			}
		},
		{
			"address": "comdex12qdnuy9a4uzm6g7a3n09k79g5t4kvejtkeupv8",
			"reward": {
				"denom": "ucmdx",
				"amount": "176"
			}
		},
		{
			"address": "comdex12qwjdjfcwrfwu2kew8wmj4qqmxa0x29369evt0",
			"reward": {
				"denom": "ucmdx",
				"amount": "9238"
			}
		},
		{
			"address": "comdex12q0nmknsp5wpc6zgav34jnsn7fmyqjuf7eczv7",
			"reward": {
				"denom": "ucmdx",
				"amount": "9983"
			}
		},
		{
			"address": "comdex12q3y4jgmctgel5hu5t5llqfs6yalgrp9n4pf9z",
			"reward": {
				"denom": "ucmdx",
				"amount": "2696"
			}
		},
		{
			"address": "comdex12q3359j6shmfmfnf4g6dke48qxlgk0zhxqqn7w",
			"reward": {
				"denom": "ucmdx",
				"amount": "1487"
			}
		},
		{
			"address": "comdex12qjg4ucp3s0fqalyspjwpvpl3dvleg75pxlwvt",
			"reward": {
				"denom": "ucmdx",
				"amount": "215"
			}
		},
		{
			"address": "comdex12q5pxqz8s6yxk05du6crfu5aqdphvup652h0ar",
			"reward": {
				"denom": "ucmdx",
				"amount": "1185"
			}
		},
		{
			"address": "comdex12q567kam3tclejefg8jggkxfxues5xlhmzgy6v",
			"reward": {
				"denom": "ucmdx",
				"amount": "3569008"
			}
		},
		{
			"address": "comdex12q4a5qvuzzn2kspe6dekkhnl77ct52dr8mmuwh",
			"reward": {
				"denom": "ucmdx",
				"amount": "1414"
			}
		},
		{
			"address": "comdex12qhqqdxv6v46w7a7mtwkrrl8uzggw2e64vt5h4",
			"reward": {
				"denom": "ucmdx",
				"amount": "4026"
			}
		},
		{
			"address": "comdex12qug6p62m7n8xju979zlphvz6ul55sf5cgxxy0",
			"reward": {
				"denom": "ucmdx",
				"amount": "1957"
			}
		},
		{
			"address": "comdex12qlf374xswuz8c7uuskfvunqnn9sqqs909ulmn",
			"reward": {
				"denom": "ucmdx",
				"amount": "1190"
			}
		},
		{
			"address": "comdex12qltw26swa52wfdw2qw3l7gv4lu4qwxn8wem9f",
			"reward": {
				"denom": "ucmdx",
				"amount": "528"
			}
		},
		{
			"address": "comdex12pxhu6ndd3zaer7e795s5qdd7xkrpvq60vy5fe",
			"reward": {
				"denom": "ucmdx",
				"amount": "4290"
			}
		},
		{
			"address": "comdex12pgyxa8nhmtwc3n8mkpcsamyjfw7sywffctl56",
			"reward": {
				"denom": "ucmdx",
				"amount": "1760"
			}
		},
		{
			"address": "comdex12pvzmpq85072wgdfealz9v823snxphgrs5sy9p",
			"reward": {
				"denom": "ucmdx",
				"amount": "190"
			}
		},
		{
			"address": "comdex12pw9vnhmv4t7ujqu2fa6ewamyw66ecr6gfeksf",
			"reward": {
				"denom": "ucmdx",
				"amount": "312"
			}
		},
		{
			"address": "comdex12pwlc2xzmx9k8l7nknkdrjk5vz730sjc6vpplv",
			"reward": {
				"denom": "ucmdx",
				"amount": "4484"
			}
		},
		{
			"address": "comdex12psvx0lk2rn6gya6jc0ktxx5c4hu3mmhd0xqy8",
			"reward": {
				"denom": "ucmdx",
				"amount": "2488"
			}
		},
		{
			"address": "comdex12pj656a09rggwy2u5umyunnagts3dwy0h6xrfx",
			"reward": {
				"denom": "ucmdx",
				"amount": "38594"
			}
		},
		{
			"address": "comdex12pnv3ejmu7a9hgws5qkewt545kut5cwrf4plf2",
			"reward": {
				"denom": "ucmdx",
				"amount": "1441"
			}
		},
		{
			"address": "comdex12pnmwj5drqmnr8cvsemfyjdj9e9cs26x3387yr",
			"reward": {
				"denom": "ucmdx",
				"amount": "18046"
			}
		},
		{
			"address": "comdex12phtrywkmtytukd3ujt5qvk2uf8k0sz5em0lq2",
			"reward": {
				"denom": "ucmdx",
				"amount": "82873"
			}
		},
		{
			"address": "comdex12pc8dyq4a5td8xhtg33pau4830rg8dm4sd3y6p",
			"reward": {
				"denom": "ucmdx",
				"amount": "1220"
			}
		},
		{
			"address": "comdex12pek06px76y2mr2fp7ew4wz8asxv3rpw8cuqpt",
			"reward": {
				"denom": "ucmdx",
				"amount": "3695"
			}
		},
		{
			"address": "comdex12p62w5am4tppy98e7ccpypsh7u0lq5c8jsjzk7",
			"reward": {
				"denom": "ucmdx",
				"amount": "570"
			}
		},
		{
			"address": "comdex12paseawz7e36l79q2as7ha8qhazh5uhm3cxshe",
			"reward": {
				"denom": "ucmdx",
				"amount": "18037"
			}
		},
		{
			"address": "comdex12zqxkytn9qaajvvrsg90zkrr7pka9mquhhet4k",
			"reward": {
				"denom": "ucmdx",
				"amount": "746"
			}
		},
		{
			"address": "comdex12zqvs2lu2wsvp9uyzre5vc39zrxrshdfz2evv3",
			"reward": {
				"denom": "ucmdx",
				"amount": "2013"
			}
		},
		{
			"address": "comdex12zq7crlxsd5zak5rjp0k5mlycjy5netmfycua5",
			"reward": {
				"denom": "ucmdx",
				"amount": "1751"
			}
		},
		{
			"address": "comdex12zzxfkxz4mmgsah8rxehw666kt2gfj7fuhgcn4",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex12zxfvxvyxkhfnjcfqkf8dxf8xvfcaynvdprgfz",
			"reward": {
				"denom": "ucmdx",
				"amount": "1208"
			}
		},
		{
			"address": "comdex12z8zwn6x36unxejflk2vfk2t5hpryp82h3pny3",
			"reward": {
				"denom": "ucmdx",
				"amount": "13134"
			}
		},
		{
			"address": "comdex12z8v6h3wtk7udfn6rmwkte22vhxhfptktdxl8r",
			"reward": {
				"denom": "ucmdx",
				"amount": "446"
			}
		},
		{
			"address": "comdex12zf52lux9kqrqqpyawrufl67js4525xcvy2q42",
			"reward": {
				"denom": "ucmdx",
				"amount": "4019"
			}
		},
		{
			"address": "comdex12zvtxlmjzqd6aw2mfycmrckl7jmj8wnwk7ugzu",
			"reward": {
				"denom": "ucmdx",
				"amount": "35431"
			}
		},
		{
			"address": "comdex12zd6wadq245ln77yx47l6e79u38trr2p6n4hw0",
			"reward": {
				"denom": "ucmdx",
				"amount": "2856"
			}
		},
		{
			"address": "comdex12zwgsrt3nj66wvykuw6tdajkszmh6hy4jzjwdl",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex12zwsxe28th39tkq9yrlcytqdy6mdrx24grdxkm",
			"reward": {
				"denom": "ucmdx",
				"amount": "7173"
			}
		},
		{
			"address": "comdex12z02thm8h9dj2h6kywgwe42u0u7uq4f8vqxpha",
			"reward": {
				"denom": "ucmdx",
				"amount": "55473"
			}
		},
		{
			"address": "comdex12z00s8nsrkqrxhnahg8jsln9yftu9vdn4tewk5",
			"reward": {
				"denom": "ucmdx",
				"amount": "11501"
			}
		},
		{
			"address": "comdex12z0s35n6d88u7udacp90qhx6f0m4sqxz03kkdc",
			"reward": {
				"denom": "ucmdx",
				"amount": "7169"
			}
		},
		{
			"address": "comdex12z3k7qmv93endnf08k0c2a5av948lyvjsmy28f",
			"reward": {
				"denom": "ucmdx",
				"amount": "142"
			}
		},
		{
			"address": "comdex12zjcmust9y5nukfatuzh2zlavjwpsg9lzvlq08",
			"reward": {
				"denom": "ucmdx",
				"amount": "143"
			}
		},
		{
			"address": "comdex12zj672nek2fjrgsrje6dfje84lt34ycw06u77d",
			"reward": {
				"denom": "ucmdx",
				"amount": "26347"
			}
		},
		{
			"address": "comdex12z5p547fu3k6n3jcwv0u8znrqrnjph0knzrl6f",
			"reward": {
				"denom": "ucmdx",
				"amount": "111930"
			}
		},
		{
			"address": "comdex12z40f0mgt73ljzzntlspp803fdyu9ejgyfz7wq",
			"reward": {
				"denom": "ucmdx",
				"amount": "1309"
			}
		},
		{
			"address": "comdex12zh9pzvewah35qwkftydmtjt7f3senk4ztqjy2",
			"reward": {
				"denom": "ucmdx",
				"amount": "1424"
			}
		},
		{
			"address": "comdex12zezz4hllguxudufp2vyc2ftewjdjl27f7xck9",
			"reward": {
				"denom": "ucmdx",
				"amount": "7629"
			}
		},
		{
			"address": "comdex12rp58mjhv09nvt0drlr6v2psc4vstlj50ggm37",
			"reward": {
				"denom": "ucmdx",
				"amount": "6126"
			}
		},
		{
			"address": "comdex12rztp7yuy7fgdjku7vxqx8dxdhvlwn97ztasnw",
			"reward": {
				"denom": "ucmdx",
				"amount": "7152"
			}
		},
		{
			"address": "comdex12r8gk60quy0v0thrms5mdlegdw6fwxmcc4s2xd",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex12rgsk2tujy3s6lf27cmnnx3te6qg02pgryg3ry",
			"reward": {
				"denom": "ucmdx",
				"amount": "140779"
			}
		},
		{
			"address": "comdex12r22qmma5yh60tpw0c3t3wzhg8tfp2dagly6y4",
			"reward": {
				"denom": "ucmdx",
				"amount": "1"
			}
		},
		{
			"address": "comdex12r25h0dqmfw5fw4rw6a9vw3zukl7t34lqw6sx8",
			"reward": {
				"denom": "ucmdx",
				"amount": "6760"
			}
		},
		{
			"address": "comdex12r2hrwts2fkpr5txv645q962k2puymtpx7jpqa",
			"reward": {
				"denom": "ucmdx",
				"amount": "21388"
			}
		},
		{
			"address": "comdex12rvv6hx977xcfufht34r5nxxz2xgfyeenduryl",
			"reward": {
				"denom": "ucmdx",
				"amount": "14448"
			}
		},
		{
			"address": "comdex12rs3sr9fqqcwcl4k6rf3sd90agy5t9ggahd9r7",
			"reward": {
				"denom": "ucmdx",
				"amount": "2138"
			}
		},
		{
			"address": "comdex12rjfj4pqleah6mtyy7ukcztz7xd4xex7zxzwyf",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex12rnp33rpj5d8kkty48em8fx4q74p0j8alhmwpq",
			"reward": {
				"denom": "ucmdx",
				"amount": "7115"
			}
		},
		{
			"address": "comdex12rege0n9rdfhg83lp7rn60ws9ezjfzplhqx535",
			"reward": {
				"denom": "ucmdx",
				"amount": "1263"
			}
		},
		{
			"address": "comdex12red9d7rwkawxg9fj62hxsn8njkl736t8028l3",
			"reward": {
				"denom": "ucmdx",
				"amount": "1409"
			}
		},
		{
			"address": "comdex12rm99fruusch08m7wayywe4uarwl0gahuv3k5e",
			"reward": {
				"denom": "ucmdx",
				"amount": "121137"
			}
		},
		{
			"address": "comdex12rlq7jf7hg2dmgkdeuezhaz3n0qktudq2x9e4s",
			"reward": {
				"denom": "ucmdx",
				"amount": "137"
			}
		},
		{
			"address": "comdex12rlmpsen59qh0tduwlw0p6vcv09hc940t7umae",
			"reward": {
				"denom": "ucmdx",
				"amount": "502"
			}
		},
		{
			"address": "comdex12ypwevfacc7amucqzej2pzrsp2pun96dqvxs5r",
			"reward": {
				"denom": "ucmdx",
				"amount": "6205"
			}
		},
		{
			"address": "comdex12yzuctsw03j3j7na7k5qaqadaevt3hljtt72wg",
			"reward": {
				"denom": "ucmdx",
				"amount": "852"
			}
		},
		{
			"address": "comdex12yyqyrwmwaqhsj6jtd8lqp834lu5hlg42xd7nm",
			"reward": {
				"denom": "ucmdx",
				"amount": "12667"
			}
		},
		{
			"address": "comdex12yyye6xqmxyem85yjrthu4dhadldvu6h29dcul",
			"reward": {
				"denom": "ucmdx",
				"amount": "6600"
			}
		},
		{
			"address": "comdex12yycuxxyavwwfw3n6fm3mmn3udu70p24eqkqyh",
			"reward": {
				"denom": "ucmdx",
				"amount": "9713"
			}
		},
		{
			"address": "comdex12y9ynga89chw4meuufppsgxxa6f0qqqc47472d",
			"reward": {
				"denom": "ucmdx",
				"amount": "36900"
			}
		},
		{
			"address": "comdex12yxku84ws7lclkuv38qrtqz4tjcmnn5sjwnwfl",
			"reward": {
				"denom": "ucmdx",
				"amount": "3550"
			}
		},
		{
			"address": "comdex12yxlpwjv0apser40dhn49hrtxr3ekt42vfsaex",
			"reward": {
				"denom": "ucmdx",
				"amount": "5642"
			}
		},
		{
			"address": "comdex12y8390a0ggcsecvga8xkp65q647p0zzkpqlkgy",
			"reward": {
				"denom": "ucmdx",
				"amount": "183"
			}
		},
		{
			"address": "comdex12yf2juf23xe6x4x7cyxu450n52r8ygetd7sm0s",
			"reward": {
				"denom": "ucmdx",
				"amount": "14242"
			}
		},
		{
			"address": "comdex12yfj9nl5r4tf0a309vp8l25hhrr9et6uq73g5l",
			"reward": {
				"denom": "ucmdx",
				"amount": "7027"
			}
		},
		{
			"address": "comdex12yt82yx4x6x8skleglxf08pt3vght99uwp43k6",
			"reward": {
				"denom": "ucmdx",
				"amount": "1900"
			}
		},
		{
			"address": "comdex12yd6ln39t26umvmms6llp4vmez59607vv4nk3q",
			"reward": {
				"denom": "ucmdx",
				"amount": "17692"
			}
		},
		{
			"address": "comdex12y0a0nttq3a2vp3kf0hvt2xw0kd8entkrdx57s",
			"reward": {
				"denom": "ucmdx",
				"amount": "199"
			}
		},
		{
			"address": "comdex12y58mrgm4j2tu8jscr0vhm5sdrusrlcx4fmmk9",
			"reward": {
				"denom": "ucmdx",
				"amount": "45328"
			}
		},
		{
			"address": "comdex12y4s4ncrpyvwrmr8m84wkynjpaah2u5s0gupar",
			"reward": {
				"denom": "ucmdx",
				"amount": "90"
			}
		},
		{
			"address": "comdex12yk6aw3mckad6lkx5u5p863gfh9jgnymkl2a4c",
			"reward": {
				"denom": "ucmdx",
				"amount": "7760"
			}
		},
		{
			"address": "comdex12ycjmqqeqa0rvzxhltf5vq5dp34qpexlqpalv8",
			"reward": {
				"denom": "ucmdx",
				"amount": "9542"
			}
		},
		{
			"address": "comdex12yc58lmcnctards59ep8ujme8c725a8vq80f30",
			"reward": {
				"denom": "ucmdx",
				"amount": "1740"
			}
		},
		{
			"address": "comdex12yevrw0vg0w5vjczhn5ev4zn72dp6nrx8kmc66",
			"reward": {
				"denom": "ucmdx",
				"amount": "1930"
			}
		},
		{
			"address": "comdex12y6reu85m88st7mnt0t8v5w4s5652054u6wlfs",
			"reward": {
				"denom": "ucmdx",
				"amount": "10142"
			}
		},
		{
			"address": "comdex12yu98usj20k36z0x02l00q2ds34ch2dc9798pf",
			"reward": {
				"denom": "ucmdx",
				"amount": "73403"
			}
		},
		{
			"address": "comdex12yu4uujayxzurzqgvwre2pfnj9pxd8289etytm",
			"reward": {
				"denom": "ucmdx",
				"amount": "1889"
			}
		},
		{
			"address": "comdex12yawcplsamryq9zg88ls9970vffrruc94z62z4",
			"reward": {
				"denom": "ucmdx",
				"amount": "765"
			}
		},
		{
			"address": "comdex1299t2cvhtamhpcfqdwqymcgww476s94zpfvm83",
			"reward": {
				"denom": "ucmdx",
				"amount": "519"
			}
		},
		{
			"address": "comdex129x8nrg0uk57985hp699dna333qwkhtfecpz2l",
			"reward": {
				"denom": "ucmdx",
				"amount": "144"
			}
		},
		{
			"address": "comdex129845g373yewyz4a4t78jxt98a6gxwp72d49ua",
			"reward": {
				"denom": "ucmdx",
				"amount": "1775"
			}
		},
		{
			"address": "comdex129g5qtypudndjjav9vpnv90um8dlddfw3n67xu",
			"reward": {
				"denom": "ucmdx",
				"amount": "271"
			}
		},
		{
			"address": "comdex129tzkmzawxkn3lxdes5pzq436nx25lr0wa56jy",
			"reward": {
				"denom": "ucmdx",
				"amount": "3192"
			}
		},
		{
			"address": "comdex129vx4cz4ql92j275eejgrfsp2qjn6n3ng7m6nv",
			"reward": {
				"denom": "ucmdx",
				"amount": "169"
			}
		},
		{
			"address": "comdex129s2zc5r845srym3jlmgddk7sjhw2y5272dpuz",
			"reward": {
				"denom": "ucmdx",
				"amount": "3569"
			}
		},
		{
			"address": "comdex1293m2slr69ewpr75m007rslxwlpy547wwdadcy",
			"reward": {
				"denom": "ucmdx",
				"amount": "173"
			}
		},
		{
			"address": "comdex129ju7mcz0nhue2wzqx6t7xy8lp8qhg3ntk6xnf",
			"reward": {
				"denom": "ucmdx",
				"amount": "169"
			}
		},
		{
			"address": "comdex129h7krwnecvgls9ccee4nndrvmsl40rmgkyf8z",
			"reward": {
				"denom": "ucmdx",
				"amount": "67480"
			}
		},
		{
			"address": "comdex1296d8q7d7wyz8904z6nhgwz9mxela6af94aalg",
			"reward": {
				"denom": "ucmdx",
				"amount": "64517"
			}
		},
		{
			"address": "comdex129u2rcqz4pvtwuhl4zupy6937kka6d7crlzldw",
			"reward": {
				"denom": "ucmdx",
				"amount": "1305"
			}
		},
		{
			"address": "comdex129avvjtjr5d8mlwzl2qaqwe0msygejkhwtk7ad",
			"reward": {
				"denom": "ucmdx",
				"amount": "6553"
			}
		},
		{
			"address": "comdex12xrqd6r5c862rg0jjvnn0ve2rekr8zhxt9n67p",
			"reward": {
				"denom": "ucmdx",
				"amount": "12900"
			}
		},
		{
			"address": "comdex12x94acgkf2kjsue0z3n4n0wfdkgf0geg9hl08r",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex12x83csf7kpr3emumrn6kw3ukavqaphxknuhn0p",
			"reward": {
				"denom": "ucmdx",
				"amount": "5288"
			}
		},
		{
			"address": "comdex12xgathfdw98j9h0snrkhuzdq4l5lpl9hfmx5l5",
			"reward": {
				"denom": "ucmdx",
				"amount": "150"
			}
		},
		{
			"address": "comdex12xfgfc9h23ptg5h7pxr8g8yyvgcn4sxnx9kfhx",
			"reward": {
				"denom": "ucmdx",
				"amount": "28435"
			}
		},
		{
			"address": "comdex12xd976xvfzpyxtued8r0aszdarm6q2ppzh2yy0",
			"reward": {
				"denom": "ucmdx",
				"amount": "37510"
			}
		},
		{
			"address": "comdex12xdlq75cqrct9awyp9wg9lp5zvt43pju2qd7v9",
			"reward": {
				"denom": "ucmdx",
				"amount": "1777"
			}
		},
		{
			"address": "comdex12xsd2tr7x5f4l84en0g5jk4t567ppddv99lcfv",
			"reward": {
				"denom": "ucmdx",
				"amount": "173"
			}
		},
		{
			"address": "comdex12xsep54vkhspullaqrt8mcl8cegvmj6yk0d263",
			"reward": {
				"denom": "ucmdx",
				"amount": "13063"
			}
		},
		{
			"address": "comdex12xnp8cquazp93uycq0se4dvac5qt88gl7qq44r",
			"reward": {
				"denom": "ucmdx",
				"amount": "178"
			}
		},
		{
			"address": "comdex12xn3s7sf555k0650fsnzwqqjcmrg4zj8q4vmwd",
			"reward": {
				"denom": "ucmdx",
				"amount": "28"
			}
		},
		{
			"address": "comdex12xnuq9ugyddsrx2s2aq3mr4q4uvckyhj8kg8mx",
			"reward": {
				"denom": "ucmdx",
				"amount": "169"
			}
		},
		{
			"address": "comdex12x5nn3shnuagejufvyauzc58e0jmpqehr87063",
			"reward": {
				"denom": "ucmdx",
				"amount": "105"
			}
		},
		{
			"address": "comdex12x4ycf2fd0p2v989x82tpu6treygdl4fchx4h3",
			"reward": {
				"denom": "ucmdx",
				"amount": "3558"
			}
		},
		{
			"address": "comdex12xmkx67gpnyk36ras4epfuy2qx2ajcuzn4uagr",
			"reward": {
				"denom": "ucmdx",
				"amount": "7201"
			}
		},
		{
			"address": "comdex12xuz204q2n3ulpx5px40a5z0uchqglnyx2r6tw",
			"reward": {
				"denom": "ucmdx",
				"amount": "12137"
			}
		},
		{
			"address": "comdex12xaxj4m5tsuncv6seytym93zyr6kkltjwfnry0",
			"reward": {
				"denom": "ucmdx",
				"amount": "1555"
			}
		},
		{
			"address": "comdex12xa46268l09fq7f3a9nv27w4y7u8ghwthgdwql",
			"reward": {
				"denom": "ucmdx",
				"amount": "6555"
			}
		},
		{
			"address": "comdex128zjnjs237ku8uqn33wyar9lz5uhyfck2p3wce",
			"reward": {
				"denom": "ucmdx",
				"amount": "91847"
			}
		},
		{
			"address": "comdex1289x8v9txvzxpmqwexn07u4jddvvlwv3ql3u3q",
			"reward": {
				"denom": "ucmdx",
				"amount": "5289"
			}
		},
		{
			"address": "comdex128g908zusk8tdskul5su664drlqz3k8y6nkcv7",
			"reward": {
				"denom": "ucmdx",
				"amount": "628494"
			}
		},
		{
			"address": "comdex128t9dmtvw0u62schyxhf56zx96vg0qn4ea3kqf",
			"reward": {
				"denom": "ucmdx",
				"amount": "6926"
			}
		},
		{
			"address": "comdex128d8xp8f70dtcjqfuu3eps989c9krn58ldl4mx",
			"reward": {
				"denom": "ucmdx",
				"amount": "123"
			}
		},
		{
			"address": "comdex128shzuyl2lmwl42ec0029d7cffela6eqvkl5et",
			"reward": {
				"denom": "ucmdx",
				"amount": "710"
			}
		},
		{
			"address": "comdex1283rurwpcl24pxa6dsnlmxhwyt0fy5yxql6mvl",
			"reward": {
				"denom": "ucmdx",
				"amount": "16"
			}
		},
		{
			"address": "comdex128j7htcfx6pna3djcgllj88hq8dg3x9wej9jtk",
			"reward": {
				"denom": "ucmdx",
				"amount": "180"
			}
		},
		{
			"address": "comdex12845qa79cwlvf3jdcnfq2jy2jfmzslcgsayr3p",
			"reward": {
				"denom": "ucmdx",
				"amount": "2200832"
			}
		},
		{
			"address": "comdex1284hyyh2vqgnq2mfegn2u8nysrmzr4q4677n3y",
			"reward": {
				"denom": "ucmdx",
				"amount": "8772"
			}
		},
		{
			"address": "comdex128h0estu7k2swf7et264v7upqcgkmgfws0zlsa",
			"reward": {
				"denom": "ucmdx",
				"amount": "323"
			}
		},
		{
			"address": "comdex128hemf7s2ntxfawe270unsa6h8gvv0qgtg7jzm",
			"reward": {
				"denom": "ucmdx",
				"amount": "17686"
			}
		},
		{
			"address": "comdex128h6rn4v64g80p436y0fn4m53svhhvenz99gm6",
			"reward": {
				"denom": "ucmdx",
				"amount": "932276"
			}
		},
		{
			"address": "comdex1286235gc7dzrunepvq55ze7kv8qym9u4ep2g7g",
			"reward": {
				"denom": "ucmdx",
				"amount": "145"
			}
		},
		{
			"address": "comdex128mgp8aflyxwxgq57wucvn3pecw0vm4dcz4ffj",
			"reward": {
				"denom": "ucmdx",
				"amount": "5967"
			}
		},
		{
			"address": "comdex12gyadyzlar4pp5nsgmfadf28td7tjkwkp27wf9",
			"reward": {
				"denom": "ucmdx",
				"amount": "14820"
			}
		},
		{
			"address": "comdex12gxdplnx7ndpa0paawtara4qaf3na7za3gc474",
			"reward": {
				"denom": "ucmdx",
				"amount": "68"
			}
		},
		{
			"address": "comdex12g8t77m7p00weya8p3cmztvurft5r5nyz6sq4k",
			"reward": {
				"denom": "ucmdx",
				"amount": "303"
			}
		},
		{
			"address": "comdex12ggnqgtpvnhgy2msu8szkc0ps5lhf3h99xej4p",
			"reward": {
				"denom": "ucmdx",
				"amount": "18290"
			}
		},
		{
			"address": "comdex12gth6nadf2qtlxcftf049yqrhfyvfa9j7z0px6",
			"reward": {
				"denom": "ucmdx",
				"amount": "89"
			}
		},
		{
			"address": "comdex12gtmz5kerg3qdavqqnk5e8rfg73fhlxuxt7nyf",
			"reward": {
				"denom": "ucmdx",
				"amount": "12397"
			}
		},
		{
			"address": "comdex12gvvucdlvkylvyttsjlcs6rful3jc9cmzrrgdv",
			"reward": {
				"denom": "ucmdx",
				"amount": "14428"
			}
		},
		{
			"address": "comdex12gnkn36sk0wfnrhsfr74ep8qa68pa4gjatu32a",
			"reward": {
				"denom": "ucmdx",
				"amount": "14248"
			}
		},
		{
			"address": "comdex12gkpztajtaflpeh3nj9k5fae7hfnflexm06kf3",
			"reward": {
				"denom": "ucmdx",
				"amount": "35"
			}
		},
		{
			"address": "comdex12g6ka3z86wrv8a8p3tkrzpd0jfuemgsrg6vgrp",
			"reward": {
				"denom": "ucmdx",
				"amount": "14253"
			}
		},
		{
			"address": "comdex12g66yxljn50ahu3nnt4txqmzadn6xta4al7nk7",
			"reward": {
				"denom": "ucmdx",
				"amount": "299"
			}
		},
		{
			"address": "comdex12guhhvz72kgh5vldpvwpfdm7u060h7u4k5uudn",
			"reward": {
				"denom": "ucmdx",
				"amount": "156348"
			}
		},
		{
			"address": "comdex12g790zmq7fjaw6e2l0g8hazwhsdyn6jekfm2xm",
			"reward": {
				"denom": "ucmdx",
				"amount": "323"
			}
		},
		{
			"address": "comdex12glnqrdcnzcddeutrxl5m0fm0dc6d0lla50mwr",
			"reward": {
				"denom": "ucmdx",
				"amount": "178"
			}
		},
		{
			"address": "comdex12fpfqec03k9862ry0pfp77nact60ftatutg4w8",
			"reward": {
				"denom": "ucmdx",
				"amount": "2380"
			}
		},
		{
			"address": "comdex12fy33rwu2m4yacr9cys07d4nfnzuxd6a9h68xt",
			"reward": {
				"denom": "ucmdx",
				"amount": "548"
			}
		},
		{
			"address": "comdex12f9qtt9c2na9e0qdqr9nhv278ewsrurhpmafn4",
			"reward": {
				"denom": "ucmdx",
				"amount": "1141"
			}
		},
		{
			"address": "comdex12fxs5fdrwzqhca0vq46ea9gvu8g0snkap5cmmu",
			"reward": {
				"denom": "ucmdx",
				"amount": "2072"
			}
		},
		{
			"address": "comdex12fgg0mcvcxhemumzzw77cxxfqhewv8fju9u2xv",
			"reward": {
				"denom": "ucmdx",
				"amount": "868"
			}
		},
		{
			"address": "comdex12ffu7le2lqspscdfhhqu7p2x8xfettgyyu3hk0",
			"reward": {
				"denom": "ucmdx",
				"amount": "14246"
			}
		},
		{
			"address": "comdex12f2t2pjgdf0ken0z3zmj668ll2p6snlaajgv93",
			"reward": {
				"denom": "ucmdx",
				"amount": "4251"
			}
		},
		{
			"address": "comdex12fjv4krxcea5zukmrddgzcu63yhx9hgppxlazl",
			"reward": {
				"denom": "ucmdx",
				"amount": "16126"
			}
		},
		{
			"address": "comdex12fnqzqy5403c4l3aahcs7kz9dwhj92n3fhu5mr",
			"reward": {
				"denom": "ucmdx",
				"amount": "106894"
			}
		},
		{
			"address": "comdex12f5e4ghlfcygn0v99zrxkeaqnna6g0pz59a9hy",
			"reward": {
				"denom": "ucmdx",
				"amount": "284"
			}
		},
		{
			"address": "comdex12f5alqtgqcmkn5ptzg4m3dp7p3e0n2u8sv4uy8",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex12f4rr3623mqderujs88uzt2dw9jx23dk4ktuu9",
			"reward": {
				"denom": "ucmdx",
				"amount": "179"
			}
		},
		{
			"address": "comdex12fc44r7ede3cdth5glwgusgj6erjyesd8aj5c4",
			"reward": {
				"denom": "ucmdx",
				"amount": "32920"
			}
		},
		{
			"address": "comdex12fmcjn8lq9x69wf48lw3nu7am6nzzehute3aan",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex12fagylmwcsnfcrzv4rm4sksk5rt9lah6ycwkws",
			"reward": {
				"denom": "ucmdx",
				"amount": "34206"
			}
		},
		{
			"address": "comdex12fl64v27w9z03fmzu9agvnvr5rxgyenflav7y5",
			"reward": {
				"denom": "ucmdx",
				"amount": "8788"
			}
		},
		{
			"address": "comdex12flun6nedtqktkttpmd8y44m05sjk6wlkwgdfc",
			"reward": {
				"denom": "ucmdx",
				"amount": "60"
			}
		},
		{
			"address": "comdex12fl7863nlusq5t43ft5mcnaf9tvwwk32vsz94w",
			"reward": {
				"denom": "ucmdx",
				"amount": "39585"
			}
		},
		{
			"address": "comdex122z0zs6ztnflpejaey5j6flt6xrly8rr8kcv26",
			"reward": {
				"denom": "ucmdx",
				"amount": "2861"
			}
		},
		{
			"address": "comdex122x75332ssh28n9qk0ng3mwrugp93nea7f3285",
			"reward": {
				"denom": "ucmdx",
				"amount": "5284"
			}
		},
		{
			"address": "comdex1228e57y86kgzp7z764u76qe9x7a76djn5p5xmg",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex122gz47xmwn2qqr8kfpcujjxlcdss85p8khhm94",
			"reward": {
				"denom": "ucmdx",
				"amount": "401"
			}
		},
		{
			"address": "comdex122fu20frv5jmtj59gr7g3yxx8azl664yhglqxw",
			"reward": {
				"denom": "ucmdx",
				"amount": "16"
			}
		},
		{
			"address": "comdex122taa46qrdv22ryk6v9wvs285y24q9z2rzth0z",
			"reward": {
				"denom": "ucmdx",
				"amount": "8250"
			}
		},
		{
			"address": "comdex122vs79vdvukvsqdwjgz49vrej66jyaj39p4zj4",
			"reward": {
				"denom": "ucmdx",
				"amount": "173"
			}
		},
		{
			"address": "comdex122jzj095fl42q0ejua96p34es8djyxtcv4y24j",
			"reward": {
				"denom": "ucmdx",
				"amount": "71646"
			}
		},
		{
			"address": "comdex122n35y2s425gffyj8pu4c5fs7lf2ka5hvwsd37",
			"reward": {
				"denom": "ucmdx",
				"amount": "4291"
			}
		},
		{
			"address": "comdex122537m0dm0s39n8v03h370z3jmmg0neawrwj8q",
			"reward": {
				"denom": "ucmdx",
				"amount": "1736"
			}
		},
		{
			"address": "comdex12247yy7a2azc9vr256wgrc3rza5vuehdzyxc4e",
			"reward": {
				"denom": "ucmdx",
				"amount": "1252"
			}
		},
		{
			"address": "comdex122c594h9aht2532djcwwlddcvw3xhgcdpcydlq",
			"reward": {
				"denom": "ucmdx",
				"amount": "13296"
			}
		},
		{
			"address": "comdex122e38cpluw24nrvx9f08gpuhfehkxk7pw7zpp2",
			"reward": {
				"denom": "ucmdx",
				"amount": "9666"
			}
		},
		{
			"address": "comdex1226q9glnrx4vvkk68050l82dazcpk4uvn82apc",
			"reward": {
				"denom": "ucmdx",
				"amount": "1708"
			}
		},
		{
			"address": "comdex1226fdauu42ler0leesnmmcedv6sna0ddyyjq23",
			"reward": {
				"denom": "ucmdx",
				"amount": "166"
			}
		},
		{
			"address": "comdex1226c8dcjapuk8jag3qp2fkwp2ruzj24h7u70nw",
			"reward": {
				"denom": "ucmdx",
				"amount": "206"
			}
		},
		{
			"address": "comdex122me5pqz754tw9z96xe4magqyeay99frjdpm85",
			"reward": {
				"denom": "ucmdx",
				"amount": "653"
			}
		},
		{
			"address": "comdex1227u7ms8eskhhy6r3wn56mkqr3v68ehwdkta2x",
			"reward": {
				"denom": "ucmdx",
				"amount": "23009"
			}
		},
		{
			"address": "comdex122l8dtcxdn27lqz5phhnafgsfehw5avmga9clt",
			"reward": {
				"denom": "ucmdx",
				"amount": "10524"
			}
		},
		{
			"address": "comdex122layu2dfd23xgyy9y5nfmhk84zw3v62049la4",
			"reward": {
				"denom": "ucmdx",
				"amount": "355"
			}
		},
		{
			"address": "comdex12tpe7t0j2yqmk8wpzg5h4vw6ze0cld87nd3weq",
			"reward": {
				"denom": "ucmdx",
				"amount": "147444"
			}
		},
		{
			"address": "comdex12tzzs9zhst95q4ynjukt65urscwc2r00pj322k",
			"reward": {
				"denom": "ucmdx",
				"amount": "17657"
			}
		},
		{
			"address": "comdex12tryytyky89napugx75u0lzd3mvp5n9nwspncq",
			"reward": {
				"denom": "ucmdx",
				"amount": "1237"
			}
		},
		{
			"address": "comdex12trcwfd0e82g6904ntgka7y72efjp82hfv92ld",
			"reward": {
				"denom": "ucmdx",
				"amount": "1584"
			}
		},
		{
			"address": "comdex12tycshmw6vmn4mhq6m4xqzs6gfg5czgs9sxfdg",
			"reward": {
				"denom": "ucmdx",
				"amount": "92904"
			}
		},
		{
			"address": "comdex12tfpxen3acrf7d03mh8a0ej92uhfjl4xgz6d8p",
			"reward": {
				"denom": "ucmdx",
				"amount": "44998"
			}
		},
		{
			"address": "comdex12tf3qe9xjcu6neuhw5ay29gknmjzs32cvvckqh",
			"reward": {
				"denom": "ucmdx",
				"amount": "6043"
			}
		},
		{
			"address": "comdex12tvsrn3c4m4fs8sux27nnmvcnkjflgw9pejd57",
			"reward": {
				"denom": "ucmdx",
				"amount": "19435"
			}
		},
		{
			"address": "comdex12td4g6m0q9n0ax7dfkzmd8jgty7nfcnk8en5cj",
			"reward": {
				"denom": "ucmdx",
				"amount": "3516"
			}
		},
		{
			"address": "comdex12tdulmvjmsmpaqrpshz0emu0h9sqz5x5rt6qaf",
			"reward": {
				"denom": "ucmdx",
				"amount": "176"
			}
		},
		{
			"address": "comdex12tw0x9995wpu0ag5h4wj36wjjnrm4qk2z90fcc",
			"reward": {
				"denom": "ucmdx",
				"amount": "105781"
			}
		},
		{
			"address": "comdex12tssy3qf2s2lsqzcrwgye0a0pgxnsfj0kqnwq9",
			"reward": {
				"denom": "ucmdx",
				"amount": "30209"
			}
		},
		{
			"address": "comdex12t3um6vptf3tskzhhlnx3zzvnweqn9txj2yzan",
			"reward": {
				"denom": "ucmdx",
				"amount": "7132"
			}
		},
		{
			"address": "comdex12t5p8z5ukgx5kcwyg9gxlmwzpzulrfq4tvvw7y",
			"reward": {
				"denom": "ucmdx",
				"amount": "57828"
			}
		},
		{
			"address": "comdex12t56xng3fls2ya3nt7ejn0suy5hdcuxxkvs86f",
			"reward": {
				"denom": "ucmdx",
				"amount": "2116"
			}
		},
		{
			"address": "comdex12t4480q2zp6tfx9u7d5ymlzutme20lmw2vmmk8",
			"reward": {
				"denom": "ucmdx",
				"amount": "31354"
			}
		},
		{
			"address": "comdex12tkp0de6zwq8p4rvvqaknwexxtx9jxvh4906wv",
			"reward": {
				"denom": "ucmdx",
				"amount": "142"
			}
		},
		{
			"address": "comdex12tkt3l96an2dwemtk4apel3fl47x4l9cgmp5yk",
			"reward": {
				"denom": "ucmdx",
				"amount": "7058"
			}
		},
		{
			"address": "comdex12tkcj5lkkrh62qe96fw4zxsm0nehukh5u70u6k",
			"reward": {
				"denom": "ucmdx",
				"amount": "140"
			}
		},
		{
			"address": "comdex12thdz9jtwy5xueqnmcg6ydv8srtdjrutcz7mzl",
			"reward": {
				"denom": "ucmdx",
				"amount": "1868"
			}
		},
		{
			"address": "comdex12thj900u8y6qqvutpk95fdyc3rsftsn2xx9a3h",
			"reward": {
				"denom": "ucmdx",
				"amount": "876"
			}
		},
		{
			"address": "comdex12tmhxcccly3nmkgapp3qwg4v6xef95vqlspf7j",
			"reward": {
				"denom": "ucmdx",
				"amount": "2380"
			}
		},
		{
			"address": "comdex12ta0cp7tj9mk6l4wyf9vk079mu008jyx59jtx6",
			"reward": {
				"denom": "ucmdx",
				"amount": "15"
			}
		},
		{
			"address": "comdex12taaz5737vgy362k92y64xspwme7l5kt944cqu",
			"reward": {
				"denom": "ucmdx",
				"amount": "1678"
			}
		},
		{
			"address": "comdex12t7s938g97g57sm3muwvtgy7s9cjtj0xs4l3k4",
			"reward": {
				"denom": "ucmdx",
				"amount": "1476"
			}
		},
		{
			"address": "comdex12vppervree7djk5rqqk0cd6sst5j9q2kjmlxze",
			"reward": {
				"denom": "ucmdx",
				"amount": "176"
			}
		},
		{
			"address": "comdex12vr8evg6tj9xyq5eln7aqnd3q9atxrv7wvhkhm",
			"reward": {
				"denom": "ucmdx",
				"amount": "1779"
			}
		},
		{
			"address": "comdex12vyndr6ludx2cu7urzm9hzagclkz0ysukmeje5",
			"reward": {
				"denom": "ucmdx",
				"amount": "71"
			}
		},
		{
			"address": "comdex12vgmdunnmgsx6xtp8x8f5gkqm98nqmzt28zzdh",
			"reward": {
				"denom": "ucmdx",
				"amount": "3414"
			}
		},
		{
			"address": "comdex12vv8t82y066x94ytvq4qsmzcc4gtfkftgjalrc",
			"reward": {
				"denom": "ucmdx",
				"amount": "630121"
			}
		},
		{
			"address": "comdex12vevntkk8nv8qfyqr6c2hnndztt05p4rlkuy3k",
			"reward": {
				"denom": "ucmdx",
				"amount": "298"
			}
		},
		{
			"address": "comdex12dz8n723vf046jyshulzh0gemn004glnskh5ue",
			"reward": {
				"denom": "ucmdx",
				"amount": "3202"
			}
		},
		{
			"address": "comdex12dy2cwv5m85ml0vzck2gp2d5x8cxflk7k75xsw",
			"reward": {
				"denom": "ucmdx",
				"amount": "1360"
			}
		},
		{
			"address": "comdex12d9k2zwuzq92v93c6whc0zdueq8qq7d6g8052q",
			"reward": {
				"denom": "ucmdx",
				"amount": "463"
			}
		},
		{
			"address": "comdex12dx30wpw2sh289e5cgeuuxqlrdwsl8re8nxwd5",
			"reward": {
				"denom": "ucmdx",
				"amount": "12371"
			}
		},
		{
			"address": "comdex12d8gk9hdnvnrryvnr25xqc27cshalh5qz2e9mn",
			"reward": {
				"denom": "ucmdx",
				"amount": "105448"
			}
		},
		{
			"address": "comdex12dgaqqtn6len74s9392y99rdq70sygy8ugtyzl",
			"reward": {
				"denom": "ucmdx",
				"amount": "90"
			}
		},
		{
			"address": "comdex12dtjn8sa290dcmp4kah8hqqg8c373d4pz8sswu",
			"reward": {
				"denom": "ucmdx",
				"amount": "284"
			}
		},
		{
			"address": "comdex12dt486p2vcsgnjr2x508xjec3sk2t2h0ug96hl",
			"reward": {
				"denom": "ucmdx",
				"amount": "17392"
			}
		},
		{
			"address": "comdex12dwst80l6lq2tzyzhcd77vqye7xnqjfe0us8l2",
			"reward": {
				"denom": "ucmdx",
				"amount": "169"
			}
		},
		{
			"address": "comdex12d5uwzcxmm5xwq9hujvu83gsp2r6ve5fle56kv",
			"reward": {
				"denom": "ucmdx",
				"amount": "2621"
			}
		},
		{
			"address": "comdex12d425eddl7nql8zydt8ep6xy02sujjmwv5dkyl",
			"reward": {
				"denom": "ucmdx",
				"amount": "4592235"
			}
		},
		{
			"address": "comdex12dhdxqkmegdfjku0zxkj32yt765jpf0kqmu36y",
			"reward": {
				"denom": "ucmdx",
				"amount": "10510"
			}
		},
		{
			"address": "comdex12dcy8t4te5kjxs62gjg9u0m43awra7njkw0k87",
			"reward": {
				"denom": "ucmdx",
				"amount": "7187"
			}
		},
		{
			"address": "comdex12d6r4sywmfg4ymw4mv2l7nrgmyryt7t3t29w5p",
			"reward": {
				"denom": "ucmdx",
				"amount": "5269"
			}
		},
		{
			"address": "comdex12dupxagq9h0cej0gss0vks3w2c3yq566g0trka",
			"reward": {
				"denom": "ucmdx",
				"amount": "14"
			}
		},
		{
			"address": "comdex12dux7zlnj2lwws0qedxetxeude9tyqrpnucd9f",
			"reward": {
				"denom": "ucmdx",
				"amount": "5218"
			}
		},
		{
			"address": "comdex12duth72wczxgzp8kdhsu6ulzu2msshe56tuvpa",
			"reward": {
				"denom": "ucmdx",
				"amount": "14316"
			}
		},
		{
			"address": "comdex12du5h2vx78k3gxvq6gfqmhxa4h9yaa0pu2kxnf",
			"reward": {
				"denom": "ucmdx",
				"amount": "11424"
			}
		},
		{
			"address": "comdex12damfath23dc6j6zpc7y97tjyehr9eqaze0352",
			"reward": {
				"denom": "ucmdx",
				"amount": "2132"
			}
		},
		{
			"address": "comdex12d7s8zdgdpgjh4g9zs63tu5rp9zmp6d6vkcqfv",
			"reward": {
				"denom": "ucmdx",
				"amount": "2137"
			}
		},
		{
			"address": "comdex12dlnhj90mngw9d6gjkdexaw0vgeusuk9h200w8",
			"reward": {
				"denom": "ucmdx",
				"amount": "10276"
			}
		},
		{
			"address": "comdex12wzsm83nd2rx7s2kfa6npk50hjgaa5xe0p9w4s",
			"reward": {
				"denom": "ucmdx",
				"amount": "1438"
			}
		},
		{
			"address": "comdex12wrphn3j6vkmx6f4asvpsd2kmp0p4utpk3gac0",
			"reward": {
				"denom": "ucmdx",
				"amount": "2679"
			}
		},
		{
			"address": "comdex12wr8uuyaf576wugenvdk8ngt2jytn4j5f47e5v",
			"reward": {
				"denom": "ucmdx",
				"amount": "201"
			}
		},
		{
			"address": "comdex12wyg2vm7cmkrgqggmptu0jluptsv2ndkr38pej",
			"reward": {
				"denom": "ucmdx",
				"amount": "10905"
			}
		},
		{
			"address": "comdex12w9gse7lgptef2s336ksjwkuczh3w8ntdws6pa",
			"reward": {
				"denom": "ucmdx",
				"amount": "4144"
			}
		},
		{
			"address": "comdex12w932cswf7lp9q2lfmwz2nuhjy4uawqwnrrq4j",
			"reward": {
				"denom": "ucmdx",
				"amount": "143"
			}
		},
		{
			"address": "comdex12w9kf8xqe52uvgavtgvjgegasvmlqueeah5sjs",
			"reward": {
				"denom": "ucmdx",
				"amount": "6232"
			}
		},
		{
			"address": "comdex12w83w79ehv0lwnczahgw6p8hyagdnhwe6hujue",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex12w8kvk6eh0mh2x8fd3h6m0rucl8xd6k26mweqf",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex12w29ctjca49yfj3wgr45w64k29xs7692yr7md3",
			"reward": {
				"denom": "ucmdx",
				"amount": "1495"
			}
		},
		{
			"address": "comdex12w2j2ztmzls5xkd8zggzxhxv28gylzt049fa4k",
			"reward": {
				"denom": "ucmdx",
				"amount": "194"
			}
		},
		{
			"address": "comdex12w26k34qah2nu9xzwymfyzcnknyjqlp22ljdtu",
			"reward": {
				"denom": "ucmdx",
				"amount": "1429"
			}
		},
		{
			"address": "comdex12w034xkc0ep3jqcyymnhgesctujg64s7er7n2q",
			"reward": {
				"denom": "ucmdx",
				"amount": "26884"
			}
		},
		{
			"address": "comdex12wn6juw74mal6jhc9nlw7ykdjcs5qtfh9pr78s",
			"reward": {
				"denom": "ucmdx",
				"amount": "7040"
			}
		},
		{
			"address": "comdex12wkaczrr8v0hc00yz5u5c7395wyhc5fpleat5n",
			"reward": {
				"denom": "ucmdx",
				"amount": "2657"
			}
		},
		{
			"address": "comdex12we5lvzn48mj3na76jpr7tzrh4q6yv8y2alrt5",
			"reward": {
				"denom": "ucmdx",
				"amount": "705"
			}
		},
		{
			"address": "comdex12w6zfkjd0w62nylha62t5yz7enpstp50uhsgn3",
			"reward": {
				"denom": "ucmdx",
				"amount": "2138"
			}
		},
		{
			"address": "comdex12wu58k5qtw3uxagn0l8fh0z2r0pzku2z20qp7n",
			"reward": {
				"denom": "ucmdx",
				"amount": "23002"
			}
		},
		{
			"address": "comdex120qs4hqvr4jhrke3fymytcxyem8cchxrclvmq9",
			"reward": {
				"denom": "ucmdx",
				"amount": "1006"
			}
		},
		{
			"address": "comdex120r8r64f00l3tlvpy39w0e66drwdaa0e95rznc",
			"reward": {
				"denom": "ucmdx",
				"amount": "3575"
			}
		},
		{
			"address": "comdex120ynu7jqzennwqmhzuxzqhsencnpwa532lwmzu",
			"reward": {
				"denom": "ucmdx",
				"amount": "1403"
			}
		},
		{
			"address": "comdex120xsk7dqc7m50jd8cpvjslcxu3r0s8mcnf8f47",
			"reward": {
				"denom": "ucmdx",
				"amount": "4308"
			}
		},
		{
			"address": "comdex1208r23krkdkl5ytdn3zkwf5k7s35nv3n8553wk",
			"reward": {
				"denom": "ucmdx",
				"amount": "206175"
			}
		},
		{
			"address": "comdex12088f4cr8lxrlyqmgynt68mj57ejgg58hgqnjm",
			"reward": {
				"denom": "ucmdx",
				"amount": "8990"
			}
		},
		{
			"address": "comdex1208ww034gaf8rz6a5632luyqts7ygrf85j20jr",
			"reward": {
				"denom": "ucmdx",
				"amount": "6997"
			}
		},
		{
			"address": "comdex1202n4t3vj9dv5yh4xhy6cfgfmjlc0ejhxwqakg",
			"reward": {
				"denom": "ucmdx",
				"amount": "7138"
			}
		},
		{
			"address": "comdex120wkyfh5pggl2z0yfw6nva5mptv9a2vg7w9ffl",
			"reward": {
				"denom": "ucmdx",
				"amount": "79"
			}
		},
		{
			"address": "comdex120whhgkr4rruxu8vmev3a0tfee5x5rpjux5349",
			"reward": {
				"denom": "ucmdx",
				"amount": "40746"
			}
		},
		{
			"address": "comdex120nvf96d8uq6jn3ehnylmejwk2lqty82s9n4u4",
			"reward": {
				"denom": "ucmdx",
				"amount": "75"
			}
		},
		{
			"address": "comdex120n7rckzunmkr8tjq5sendwpcvxrw66ztusr6y",
			"reward": {
				"denom": "ucmdx",
				"amount": "1036"
			}
		},
		{
			"address": "comdex1205sdk56keedtawwdtvtaeujcg6azu6xq09egp",
			"reward": {
				"denom": "ucmdx",
				"amount": "11534"
			}
		},
		{
			"address": "comdex120kecw8c763rm5xq4fm5dn4cugnm064cw635pa",
			"reward": {
				"denom": "ucmdx",
				"amount": "1711"
			}
		},
		{
			"address": "comdex120uznw52cwt76ec6tw7utjnk8qypnlvtc4fkx2",
			"reward": {
				"denom": "ucmdx",
				"amount": "17711"
			}
		},
		{
			"address": "comdex12sqnxw4n8kqlchey5qj3z7tdlztej56um8udgg",
			"reward": {
				"denom": "ucmdx",
				"amount": "178"
			}
		},
		{
			"address": "comdex12szrq92s6unv0hxtqvl6ffu076qv5l0dw6jk39",
			"reward": {
				"denom": "ucmdx",
				"amount": "1802"
			}
		},
		{
			"address": "comdex12szmj4nezpxgvqs53yt74mwla5wdzpx6gghup6",
			"reward": {
				"denom": "ucmdx",
				"amount": "6917"
			}
		},
		{
			"address": "comdex12swcvuwrtkl799gr9kz6kfuvqncj39rcjl3uky",
			"reward": {
				"denom": "ucmdx",
				"amount": "881"
			}
		},
		{
			"address": "comdex12s0wzp3atpc73djz0nwprvde0cej5q2wfek4n8",
			"reward": {
				"denom": "ucmdx",
				"amount": "45"
			}
		},
		{
			"address": "comdex12s3pvtsk664sa50kw7te2rh735c026l67h489n",
			"reward": {
				"denom": "ucmdx",
				"amount": "14066"
			}
		},
		{
			"address": "comdex12s3mr0hsfsexw75l6sp8mcjfnyw3rl4nnxukgf",
			"reward": {
				"denom": "ucmdx",
				"amount": "6630"
			}
		},
		{
			"address": "comdex12sj6emevn267lm5xmjjfktdlmrkdw96ddrnf0l",
			"reward": {
				"denom": "ucmdx",
				"amount": "271"
			}
		},
		{
			"address": "comdex12sn3e4lpv3pader8edkk758n9stgll5wvgqg90",
			"reward": {
				"denom": "ucmdx",
				"amount": "176"
			}
		},
		{
			"address": "comdex12skywkea2p5hppvl8h9qaxyv6k8funsv58vryy",
			"reward": {
				"denom": "ucmdx",
				"amount": "6534"
			}
		},
		{
			"address": "comdex12sk0xlr5krg5q94v8ryfatv4jjruy8ftm0kkn5",
			"reward": {
				"denom": "ucmdx",
				"amount": "1896"
			}
		},
		{
			"address": "comdex12sm9vfj4m95nr0974kjhjlfx4eyum37h279fk8",
			"reward": {
				"denom": "ucmdx",
				"amount": "20084"
			}
		},
		{
			"address": "comdex12sm48v3llrcfrs6zdjtgxrv2vd87ds7yk7yxzj",
			"reward": {
				"denom": "ucmdx",
				"amount": "1940"
			}
		},
		{
			"address": "comdex12smltf87hx94nyur0ygfg3yq63w2ddsnhude4s",
			"reward": {
				"denom": "ucmdx",
				"amount": "1059"
			}
		},
		{
			"address": "comdex12sulya6ftehkgu4gn28t3yevgu3l0qf0pptzp2",
			"reward": {
				"denom": "ucmdx",
				"amount": "14359"
			}
		},
		{
			"address": "comdex12saxy20x0x3n7efuxpfeq4jqp0vqpazyrkhm0v",
			"reward": {
				"denom": "ucmdx",
				"amount": "14471"
			}
		},
		{
			"address": "comdex12sakst7l24er29068mpngzkdq70qfc0u3m8k84",
			"reward": {
				"denom": "ucmdx",
				"amount": "9396"
			}
		},
		{
			"address": "comdex12slt96ddw9wtppg7ce2hgxvcax4wsxfp4xugdu",
			"reward": {
				"denom": "ucmdx",
				"amount": "75201"
			}
		},
		{
			"address": "comdex123q27a4y4uw9qqnv0j8fml4hht63kpt9m074g6",
			"reward": {
				"denom": "ucmdx",
				"amount": "1733"
			}
		},
		{
			"address": "comdex123pjqngfj0zalpwwmfx99qzz49k9j0ncuggaw2",
			"reward": {
				"denom": "ucmdx",
				"amount": "944"
			}
		},
		{
			"address": "comdex123r7mnlg2c46nsltvd0vdfxukjrc8mwu9v6a9s",
			"reward": {
				"denom": "ucmdx",
				"amount": "521"
			}
		},
		{
			"address": "comdex123ysrrpufxp40vwpusayrhuxcdpgh9gzfncqjl",
			"reward": {
				"denom": "ucmdx",
				"amount": "126830"
			}
		},
		{
			"address": "comdex1239hwtdnlwekguzcuum57m3xwkzcpmgy65vkez",
			"reward": {
				"denom": "ucmdx",
				"amount": "137"
			}
		},
		{
			"address": "comdex12389qehnyt8t35svf3s77f6zkquqvtzzscfaqx",
			"reward": {
				"denom": "ucmdx",
				"amount": "180293"
			}
		},
		{
			"address": "comdex123fpwvrw5cjqjv7a72jsxj09e72epz72wwma9k",
			"reward": {
				"denom": "ucmdx",
				"amount": "681"
			}
		},
		{
			"address": "comdex123vke4yzvy87yzauj6t722jcqjmnapy27p9s7d",
			"reward": {
				"denom": "ucmdx",
				"amount": "13325"
			}
		},
		{
			"address": "comdex12308w073hvsw3kj3ustzphngc2ggeltnwnj33k",
			"reward": {
				"denom": "ucmdx",
				"amount": "2116"
			}
		},
		{
			"address": "comdex1230mxz4tcynq3np7frwsfvzs7hcsczhcfwl3mj",
			"reward": {
				"denom": "ucmdx",
				"amount": "7109"
			}
		},
		{
			"address": "comdex123s58tgcv6ca6e5607gd44e66lt6aqayhu6h6w",
			"reward": {
				"denom": "ucmdx",
				"amount": "3533"
			}
		},
		{
			"address": "comdex12347820q7vag6j4vrqa7623jkcq3shzqwq7k95",
			"reward": {
				"denom": "ucmdx",
				"amount": "12367"
			}
		},
		{
			"address": "comdex123kjc4dr4clqt3xvkle3at56zasfcu8m6wpp9m",
			"reward": {
				"denom": "ucmdx",
				"amount": "1518"
			}
		},
		{
			"address": "comdex123u98xsk32ra8xz3fzh9g7yu4ktr4xe7uz744u",
			"reward": {
				"denom": "ucmdx",
				"amount": "355"
			}
		},
		{
			"address": "comdex123u6jvpwku3f2x5mw68rwywgrj7640xj53ywql",
			"reward": {
				"denom": "ucmdx",
				"amount": "819"
			}
		},
		{
			"address": "comdex1237dyvth38gr8zkkqpylvxchvzqqkehc04wpjf",
			"reward": {
				"denom": "ucmdx",
				"amount": "18003"
			}
		},
		{
			"address": "comdex123lptx825enll3s3fvp9vkc0gsue5esqe9m44n",
			"reward": {
				"denom": "ucmdx",
				"amount": "19"
			}
		},
		{
			"address": "comdex12jpda635udcez2nchdnjpp9q6fq8e6aqh9wreg",
			"reward": {
				"denom": "ucmdx",
				"amount": "1342"
			}
		},
		{
			"address": "comdex12jzvhl7edt2ymepsrqjtgkkty6s63k8mhjle3w",
			"reward": {
				"denom": "ucmdx",
				"amount": "14"
			}
		},
		{
			"address": "comdex12jv2k0zyf0xf98w0lcy3uug8upw3xrq8elphsw",
			"reward": {
				"denom": "ucmdx",
				"amount": "5690"
			}
		},
		{
			"address": "comdex12jd96c9mwvarc0n4vesh3xjjqu0kjvackmd9re",
			"reward": {
				"denom": "ucmdx",
				"amount": "5201"
			}
		},
		{
			"address": "comdex12j07sxxvjdtc4yz2k52szzfuq0644uhex0fh5x",
			"reward": {
				"denom": "ucmdx",
				"amount": "5863"
			}
		},
		{
			"address": "comdex12j5pe3p4tnrcd4k6fzwgfdrfaq43uwca6gvaza",
			"reward": {
				"denom": "ucmdx",
				"amount": "1490"
			}
		},
		{
			"address": "comdex12j5g7exnrk89gwemxgpsu35t6t5llj85c3q2wq",
			"reward": {
				"denom": "ucmdx",
				"amount": "16"
			}
		},
		{
			"address": "comdex12jes88juvmycew3a4kmlq4rp5a9ax077qc6s40",
			"reward": {
				"denom": "ucmdx",
				"amount": "14197"
			}
		},
		{
			"address": "comdex12jmerz5kqgrpeqv6wgdct887c586qfx6ztkp7m",
			"reward": {
				"denom": "ucmdx",
				"amount": "1240"
			}
		},
		{
			"address": "comdex12ju6erx80tfljh7rxrgywxts8dhf0c6m73dqyx",
			"reward": {
				"denom": "ucmdx",
				"amount": "917"
			}
		},
		{
			"address": "comdex12ja0ykm78j5aykyxs8qc0yavzhx0x8299tsa8h",
			"reward": {
				"denom": "ucmdx",
				"amount": "585"
			}
		},
		{
			"address": "comdex12nzlm09nv5gpxc0v6f3p70m00jzmfeak2mqx95",
			"reward": {
				"denom": "ucmdx",
				"amount": "2065"
			}
		},
		{
			"address": "comdex12n8q72xus3ul4g6uqnffn6uu8w0gpe962zf66y",
			"reward": {
				"denom": "ucmdx",
				"amount": "889"
			}
		},
		{
			"address": "comdex12n25myjsw2dlgq50sjc8jf5yymfv4g5dpczgwn",
			"reward": {
				"denom": "ucmdx",
				"amount": "3022"
			}
		},
		{
			"address": "comdex12n2utzy6v4x8z5e53nlec6e76ujxxnmq400fej",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex12nty2dv8pmffsf04ur0qf0n2qnvqr9uvn489hs",
			"reward": {
				"denom": "ucmdx",
				"amount": "2753"
			}
		},
		{
			"address": "comdex12nd60d3pkh739n4lpf2kuz5kaxgfxdr9xpuu9v",
			"reward": {
				"denom": "ucmdx",
				"amount": "183"
			}
		},
		{
			"address": "comdex12nwxdgnxketc8dsz94cz2jd5cd9w2nvxjq5fmq",
			"reward": {
				"denom": "ucmdx",
				"amount": "16041"
			}
		},
		{
			"address": "comdex12nn0z0v8v94rp9w8r7x58p36nucnp25xn6n9sn",
			"reward": {
				"denom": "ucmdx",
				"amount": "124165"
			}
		},
		{
			"address": "comdex12n5vnq35gmzcpw0pztvwfuncjy5z6cu26rvjyx",
			"reward": {
				"denom": "ucmdx",
				"amount": "144128"
			}
		},
		{
			"address": "comdex12n4a3l5c8yx4nazyzucu36hnwkqgmeyneeayd3",
			"reward": {
				"denom": "ucmdx",
				"amount": "1472"
			}
		},
		{
			"address": "comdex12ncw76jwtfny7jwrsgwtq97uyap40zv20egwxh",
			"reward": {
				"denom": "ucmdx",
				"amount": "17601"
			}
		},
		{
			"address": "comdex12n7enymltefqmsu3khackgtus88lr4n93qj8a8",
			"reward": {
				"denom": "ucmdx",
				"amount": "37848"
			}
		},
		{
			"address": "comdex12nlk4d2pp4uyrc3zad6h926zsu9kwly84uexau",
			"reward": {
				"denom": "ucmdx",
				"amount": "463"
			}
		},
		{
			"address": "comdex125q9fc42k6kzzdw8yhxhu2cmdfuca99q896ulk",
			"reward": {
				"denom": "ucmdx",
				"amount": "1051"
			}
		},
		{
			"address": "comdex125pfnjzk52ndn5zdau784e33ryzx594fd2eu9s",
			"reward": {
				"denom": "ucmdx",
				"amount": "376"
			}
		},
		{
			"address": "comdex125z34zt5guhqff548hk4pcvrnpzfckjdkqlp47",
			"reward": {
				"denom": "ucmdx",
				"amount": "2471"
			}
		},
		{
			"address": "comdex125r79n444n58hte2hy0cnp7ljsg4xe5dwrqvf0",
			"reward": {
				"denom": "ucmdx",
				"amount": "70842"
			}
		},
		{
			"address": "comdex125yeedv990x048p8azacga80dhudd98vhs87a6",
			"reward": {
				"denom": "ucmdx",
				"amount": "166"
			}
		},
		{
			"address": "comdex125s8mg37cj8ht5u0k2dmhyruzh0ae6xgmq6s8e",
			"reward": {
				"denom": "ucmdx",
				"amount": "73"
			}
		},
		{
			"address": "comdex125sj605dc43r92gksqklx4yt56zxxwpac0qw2e",
			"reward": {
				"denom": "ucmdx",
				"amount": "5424"
			}
		},
		{
			"address": "comdex125j0q8ch4w42gta4zr4uxnyx4yat0ev92t2r4n",
			"reward": {
				"denom": "ucmdx",
				"amount": "5395"
			}
		},
		{
			"address": "comdex125nkhdvfwrt9qfleypcggkdwq2fmfgz5u8q0ke",
			"reward": {
				"denom": "ucmdx",
				"amount": "23452"
			}
		},
		{
			"address": "comdex1255kd9glzk68ltavvya72m9fz03wte5p5zwujz",
			"reward": {
				"denom": "ucmdx",
				"amount": "1774"
			}
		},
		{
			"address": "comdex1254988hlhrpakkxjxd3umsyp85fwufpytpn6zs",
			"reward": {
				"denom": "ucmdx",
				"amount": "1509"
			}
		},
		{
			"address": "comdex125kk4urdck6x9x3lg2jkfjj8g9lz2lvff2x6rp",
			"reward": {
				"denom": "ucmdx",
				"amount": "230"
			}
		},
		{
			"address": "comdex125etrnpnsv8xmtnem4xg8g6zz89lxsxznyqtce",
			"reward": {
				"denom": "ucmdx",
				"amount": "1234"
			}
		},
		{
			"address": "comdex1256vpzzcz4k2ve35tekml7mdv5j50ruluwsrdj",
			"reward": {
				"denom": "ucmdx",
				"amount": "7614"
			}
		},
		{
			"address": "comdex125u8z4kcjn3puuhsgzu8amdks834ysvg05ar3t",
			"reward": {
				"denom": "ucmdx",
				"amount": "2"
			}
		},
		{
			"address": "comdex125a58w4r4af5swp0n7pvdun70ftndhqg8qpfau",
			"reward": {
				"denom": "ucmdx",
				"amount": "10121"
			}
		},
		{
			"address": "comdex125lzh0s82aqcjp4lmmdyjk98gmf7q6cqxhk578",
			"reward": {
				"denom": "ucmdx",
				"amount": "6860"
			}
		},
		{
			"address": "comdex124quae48dtus9jsf50edxlxj53th3pdgezh4q3",
			"reward": {
				"denom": "ucmdx",
				"amount": "90"
			}
		},
		{
			"address": "comdex124pc84ddvn82hea3wj784pk9pv92de0lr357xp",
			"reward": {
				"denom": "ucmdx",
				"amount": "179197"
			}
		},
		{
			"address": "comdex124z22vznc4znxcv96w3xv958trj9z0w7cqdwym",
			"reward": {
				"denom": "ucmdx",
				"amount": "2757"
			}
		},
		{
			"address": "comdex124z3naqypgkxn6zlrhqqfygmhz6ky78fsg926h",
			"reward": {
				"denom": "ucmdx",
				"amount": "261"
			}
		},
		{
			"address": "comdex124rmfu4x0hmaaxdmxkcw6ttdvd4v26pnjfard7",
			"reward": {
				"denom": "ucmdx",
				"amount": "394"
			}
		},
		{
			"address": "comdex124ysmvvr8rprzhq8ksec5mjnkf56hujgqqudcz",
			"reward": {
				"denom": "ucmdx",
				"amount": "6298"
			}
		},
		{
			"address": "comdex124xp72fgkzd68ukccrevcu3yvd54vafa2hr4ge",
			"reward": {
				"denom": "ucmdx",
				"amount": "1430"
			}
		},
		{
			"address": "comdex124g9k42lfhrgl0e64c4rrnqw0gf4x8nw0azuc3",
			"reward": {
				"denom": "ucmdx",
				"amount": "5914"
			}
		},
		{
			"address": "comdex124f3grfjycuxhdheps9ynqpmkunsudwmvxvwc2",
			"reward": {
				"denom": "ucmdx",
				"amount": "180"
			}
		},
		{
			"address": "comdex1242tjd58cedqy9ju6wwrccv88n92an9lexaz4z",
			"reward": {
				"denom": "ucmdx",
				"amount": "14336"
			}
		},
		{
			"address": "comdex124tztqwj5lkf434pq9jjcxykda2mn964gr4972",
			"reward": {
				"denom": "ucmdx",
				"amount": "201"
			}
		},
		{
			"address": "comdex124tehz6l0wpc6nq0d60umwltepd8v6ec3cqlc4",
			"reward": {
				"denom": "ucmdx",
				"amount": "65836"
			}
		},
		{
			"address": "comdex124dmwg49fpm8sjfdw52gr0cjn4u2fls3md5aqc",
			"reward": {
				"denom": "ucmdx",
				"amount": "12216"
			}
		},
		{
			"address": "comdex124wganqs64h0cwpen9hjqna9pl4tg3cpx9cgzs",
			"reward": {
				"denom": "ucmdx",
				"amount": "1492"
			}
		},
		{
			"address": "comdex1240llqnxw55zers7ercvmgcw5mdzav5yf6m53l",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex124sajx7h2z22wx4ux6g6jfa7wgsqxz80jrtha8",
			"reward": {
				"denom": "ucmdx",
				"amount": "11217"
			}
		},
		{
			"address": "comdex1243yhe2plcqqkhwkl8dwv53nzvxrdye2n20mrd",
			"reward": {
				"denom": "ucmdx",
				"amount": "1239"
			}
		},
		{
			"address": "comdex12457le2n07ndr7ecrlmpqg468rc8w9833eqgvz",
			"reward": {
				"denom": "ucmdx",
				"amount": "79588"
			}
		},
		{
			"address": "comdex12444l9ysjhr3tyxjq8xzl04qe20945wh3gs75v",
			"reward": {
				"denom": "ucmdx",
				"amount": "352571"
			}
		},
		{
			"address": "comdex124kp99ktauryst7w8peuzztyc026l2utkp6766",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex124cprh76gavkyddmxc9fe6edyk7w7mgdcghms7",
			"reward": {
				"denom": "ucmdx",
				"amount": "1267"
			}
		},
		{
			"address": "comdex124esrkjggn9tj7r98s85px3d6u3yq3dtytl5yt",
			"reward": {
				"denom": "ucmdx",
				"amount": "42500"
			}
		},
		{
			"address": "comdex124up4umc5258ran4mqx7c8rxka3hxg8artj583",
			"reward": {
				"denom": "ucmdx",
				"amount": "908"
			}
		},
		{
			"address": "comdex124afp727hec9673hhjf328vaxyqyr7ge0cwwya",
			"reward": {
				"denom": "ucmdx",
				"amount": "179"
			}
		},
		{
			"address": "comdex124asthmmf33w0vxw06agxxntvyhuw4fknzuvrs",
			"reward": {
				"denom": "ucmdx",
				"amount": "1960"
			}
		},
		{
			"address": "comdex12kpds4umq4vugtcy0huwsvdfhu8f4q06sfs5yd",
			"reward": {
				"denom": "ucmdx",
				"amount": "21188"
			}
		},
		{
			"address": "comdex12kpc2mfp48wf58kwc9t89rrceexqd5ky2ta5cx",
			"reward": {
				"denom": "ucmdx",
				"amount": "140333"
			}
		},
		{
			"address": "comdex12k84zagyh5tuk8je2ytps57wzw3l5cm2h0ufkm",
			"reward": {
				"denom": "ucmdx",
				"amount": "13311"
			}
		},
		{
			"address": "comdex12kgaa3ctcv80nna0t8asypdqazyn98y5330s9a",
			"reward": {
				"denom": "ucmdx",
				"amount": "6261"
			}
		},
		{
			"address": "comdex12kfrsqky3glf2v8xl23wraxrunnwgprwz8hvy3",
			"reward": {
				"denom": "ucmdx",
				"amount": "6650"
			}
		},
		{
			"address": "comdex12k2w7esq48uhtqguw3egsagmyg44w4286lk6e6",
			"reward": {
				"denom": "ucmdx",
				"amount": "1738"
			}
		},
		{
			"address": "comdex12kdnd6kjvl6e5th42fnpqpdq2qz49s8u75wxm8",
			"reward": {
				"denom": "ucmdx",
				"amount": "35481"
			}
		},
		{
			"address": "comdex12kw9sl2vjp9wdpsexzeyew2uqnpcjpmj2re7q0",
			"reward": {
				"denom": "ucmdx",
				"amount": "13689"
			}
		},
		{
			"address": "comdex12ks357k4ngmd5aarg7h9s056mse06cx5vwclk4",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex12kct4cfuwh7vevec7c4y9slmee949qrjj8c0cq",
			"reward": {
				"denom": "ucmdx",
				"amount": "714"
			}
		},
		{
			"address": "comdex12kc5jtn0p3lcpk2ydsu4tp2h054suygya404u7",
			"reward": {
				"denom": "ucmdx",
				"amount": "2180"
			}
		},
		{
			"address": "comdex12kc4naxs5anh40tyjphkn2g8ypvvfzvtpr56qy",
			"reward": {
				"denom": "ucmdx",
				"amount": "1243"
			}
		},
		{
			"address": "comdex12k6dxm5v66txevqx0fg5dcx3cd9gz5x3n8la75",
			"reward": {
				"denom": "ucmdx",
				"amount": "677"
			}
		},
		{
			"address": "comdex12kmma4mjtnmnqhcky6fppg4g94zc5mvkwy7pqt",
			"reward": {
				"denom": "ucmdx",
				"amount": "4138744"
			}
		},
		{
			"address": "comdex12kurx8fljhtp500gqux9wa4savme9vm8qvp7zm",
			"reward": {
				"denom": "ucmdx",
				"amount": "1190"
			}
		},
		{
			"address": "comdex12k797kt6m3rrfewalzap2c5yu8crt2eqlc23fm",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex12hrm3qpflqn90ddzvwcdajkezfuwz7vcax68zn",
			"reward": {
				"denom": "ucmdx",
				"amount": "1092"
			}
		},
		{
			"address": "comdex12h8m6zd7d390mhhcyqp3zgc9epd4dk97l336e4",
			"reward": {
				"denom": "ucmdx",
				"amount": "1602"
			}
		},
		{
			"address": "comdex12ht53ghfrtd8792e62xta6z37ma8aapx2mtn3d",
			"reward": {
				"denom": "ucmdx",
				"amount": "142"
			}
		},
		{
			"address": "comdex12hdshxqjndfqmvu7uaua0sdsaqp3u5g2jsexrt",
			"reward": {
				"denom": "ucmdx",
				"amount": "143"
			}
		},
		{
			"address": "comdex12hwz6v0xfhatxelgrq2nyru3euzrnl7gvacns6",
			"reward": {
				"denom": "ucmdx",
				"amount": "50"
			}
		},
		{
			"address": "comdex12hj6c9szp8dcahlef8p2sq2ldng3tyrtepfudv",
			"reward": {
				"denom": "ucmdx",
				"amount": "3377"
			}
		},
		{
			"address": "comdex12hnhnxtej4qqq9072lc2rxva405wpgsr3lucq4",
			"reward": {
				"denom": "ucmdx",
				"amount": "28747"
			}
		},
		{
			"address": "comdex12h4krve2t7fwcxlte625u8fsasttvrqdxmlvjt",
			"reward": {
				"denom": "ucmdx",
				"amount": "170"
			}
		},
		{
			"address": "comdex12hkfxs9am3v73ckm2euylf9p6l28efpa8gydp2",
			"reward": {
				"denom": "ucmdx",
				"amount": "886"
			}
		},
		{
			"address": "comdex12hmrnuecynkmzrpj77fw0nj8g44cy8t32lrtuf",
			"reward": {
				"denom": "ucmdx",
				"amount": "8660"
			}
		},
		{
			"address": "comdex12hmg2uu36d4ga3rrpzeycccceyrw6cd7lxlr25",
			"reward": {
				"denom": "ucmdx",
				"amount": "42828"
			}
		},
		{
			"address": "comdex12h7jly5qacf8nswgzlh8nkm626kw9t7fmdt5jl",
			"reward": {
				"denom": "ucmdx",
				"amount": "2853"
			}
		},
		{
			"address": "comdex12cq90xeg0828dtdqgtkgkl60xfgvwa3jana090",
			"reward": {
				"denom": "ucmdx",
				"amount": "1058"
			}
		},
		{
			"address": "comdex12cpylquksrpmtr4zga4yfquj7l8xzjcvguen06",
			"reward": {
				"denom": "ucmdx",
				"amount": "24"
			}
		},
		{
			"address": "comdex12cra4p7vvfxmedt487weyx9dwhnsgeelacut3y",
			"reward": {
				"denom": "ucmdx",
				"amount": "1978"
			}
		},
		{
			"address": "comdex12c9pcgtwtvz43nhf62kzu48tvpjr7u8jffj3l7",
			"reward": {
				"denom": "ucmdx",
				"amount": "1772"
			}
		},
		{
			"address": "comdex12c9fmlj4a37ke6ca65kgw0d9vgt7p9flexe7wy",
			"reward": {
				"denom": "ucmdx",
				"amount": "1795"
			}
		},
		{
			"address": "comdex12cxrdzxnetxcc30epnygytzzhtarkp2g944xen",
			"reward": {
				"denom": "ucmdx",
				"amount": "3591"
			}
		},
		{
			"address": "comdex12cx4rfdv9cz2rspsrqpss8xtgnf9pywphtc3u9",
			"reward": {
				"denom": "ucmdx",
				"amount": "35243"
			}
		},
		{
			"address": "comdex12cxhx7h7qcancqkemyzxqv7myn6r06f4u6l6wn",
			"reward": {
				"denom": "ucmdx",
				"amount": "2649"
			}
		},
		{
			"address": "comdex12ctg3ztj92qtea7na5nx4cgw7qvgm3l4hc37nz",
			"reward": {
				"denom": "ucmdx",
				"amount": "1780"
			}
		},
		{
			"address": "comdex12cwplrtjfvt0q7qwkyzun5pudcwdluzzcth2fv",
			"reward": {
				"denom": "ucmdx",
				"amount": "174"
			}
		},
		{
			"address": "comdex12cw2jz7pc0vjqzdcjvrgsjydl590zprmatapsd",
			"reward": {
				"denom": "ucmdx",
				"amount": "61990"
			}
		},
		{
			"address": "comdex12c063z2cwq7737r2mw3jx5nv7yz55fpllt7ye3",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex12c34ueuzagjduz5cn4m683kwlt5kq8d28mvh0w",
			"reward": {
				"denom": "ucmdx",
				"amount": "403"
			}
		},
		{
			"address": "comdex12cjzzqum5f2fvqgd47p377003l9xhksac5lhsk",
			"reward": {
				"denom": "ucmdx",
				"amount": "15095"
			}
		},
		{
			"address": "comdex12c4xfww70fe4h37hncdea8lt80dny8aeml5547",
			"reward": {
				"denom": "ucmdx",
				"amount": "35721"
			}
		},
		{
			"address": "comdex12chq77z5yzw6mtncrh9kewzxcehakpe2feh8qj",
			"reward": {
				"denom": "ucmdx",
				"amount": "886"
			}
		},
		{
			"address": "comdex12chuqek35shscd68jr59awfwr2uxed5q3y5hu9",
			"reward": {
				"denom": "ucmdx",
				"amount": "3306"
			}
		},
		{
			"address": "comdex12c6y9dm97r05cu2ewjjw0gm4f6y0pnsmcuscn3",
			"reward": {
				"denom": "ucmdx",
				"amount": "8840"
			}
		},
		{
			"address": "comdex12eqgmrepnd8404ugkefv9djn7hzvhy8vnens60",
			"reward": {
				"denom": "ucmdx",
				"amount": "12470"
			}
		},
		{
			"address": "comdex12ez2trpkk4crl7gtmldsdaf6xyla6dywe8ayjm",
			"reward": {
				"denom": "ucmdx",
				"amount": "5795"
			}
		},
		{
			"address": "comdex12ey4ejg2zzpht4arheehxf5j7plu0h2x4ejdax",
			"reward": {
				"denom": "ucmdx",
				"amount": "69587"
			}
		},
		{
			"address": "comdex12e92d84r0nrpzr3wpqnd3kce7lnqdwtx938nsg",
			"reward": {
				"denom": "ucmdx",
				"amount": "5330"
			}
		},
		{
			"address": "comdex12expemzay2gewstc7ajq5652yyvjwqmqnsy0pk",
			"reward": {
				"denom": "ucmdx",
				"amount": "481"
			}
		},
		{
			"address": "comdex12exytstv8ewykz2wgqhnd7cszzhahtu0q326h8",
			"reward": {
				"denom": "ucmdx",
				"amount": "503"
			}
		},
		{
			"address": "comdex12efzmu9ttklrvfsm385hgek7llh0892cfqz2ds",
			"reward": {
				"denom": "ucmdx",
				"amount": "145"
			}
		},
		{
			"address": "comdex12et4g4dfctr8qy9pqk9ns57u6cf9y0xetpz9ql",
			"reward": {
				"denom": "ucmdx",
				"amount": "7658"
			}
		},
		{
			"address": "comdex12edyx3taya7887ml7uv6v48flhdu3hkehayq2y",
			"reward": {
				"denom": "ucmdx",
				"amount": "24793"
			}
		},
		{
			"address": "comdex12edhat8859a0m5p86twpn5caax0f3jpc4nrtuc",
			"reward": {
				"denom": "ucmdx",
				"amount": "3528"
			}
		},
		{
			"address": "comdex12ew4jq3thxauq52yug4c7larpkhgxvk74rrp5q",
			"reward": {
				"denom": "ucmdx",
				"amount": "61740"
			}
		},
		{
			"address": "comdex12e0yagc4ya6cmns6pgwhhht3pkumgr87mh4wzp",
			"reward": {
				"denom": "ucmdx",
				"amount": "40017"
			}
		},
		{
			"address": "comdex12e3u77y4caujhckrkf6w9x2qx53he248jsw06v",
			"reward": {
				"denom": "ucmdx",
				"amount": "64"
			}
		},
		{
			"address": "comdex12e7ptnnn9cp48axz0s2wh5suumh8zfu4fjyuu2",
			"reward": {
				"denom": "ucmdx",
				"amount": "1990"
			}
		},
		{
			"address": "comdex12e70xta4kqv2qh8zudqtvwwlfycwexs25ffmks",
			"reward": {
				"denom": "ucmdx",
				"amount": "173"
			}
		},
		{
			"address": "comdex12elwumqx6vquat403pkplzjnldlze0ep8sd6qe",
			"reward": {
				"denom": "ucmdx",
				"amount": "1189"
			}
		},
		{
			"address": "comdex126pq3ajt59jxlyumnhw0rzsqa0lnakhmhz5hjm",
			"reward": {
				"denom": "ucmdx",
				"amount": "3425"
			}
		},
		{
			"address": "comdex126zvv3fmvpst8v9jy7qgqjs9975qvgal8fhd30",
			"reward": {
				"denom": "ucmdx",
				"amount": "11632"
			}
		},
		{
			"address": "comdex126ygdxx4uk25pvmhnzjq20clla6heymwx305m3",
			"reward": {
				"denom": "ucmdx",
				"amount": "35127"
			}
		},
		{
			"address": "comdex126t4zpdh4leq2suxeud0fjhuy5xxvzsrm9x40x",
			"reward": {
				"denom": "ucmdx",
				"amount": "74921"
			}
		},
		{
			"address": "comdex126vrwh2fmmd0lje32dp36pureq0xyytyz6sx2d",
			"reward": {
				"denom": "ucmdx",
				"amount": "28224"
			}
		},
		{
			"address": "comdex126wfrll0d5lz2yslxqd55ua2g0hkv8zsvnsn6t",
			"reward": {
				"denom": "ucmdx",
				"amount": "1"
			}
		},
		{
			"address": "comdex126w0w3x5w0fe8cfw6wf02p44f3ey5skgrc2rey",
			"reward": {
				"denom": "ucmdx",
				"amount": "1516"
			}
		},
		{
			"address": "comdex126s88pyn4ar3cwfx4f3s8h0m0vcsxskzgenqz3",
			"reward": {
				"denom": "ucmdx",
				"amount": "8648"
			}
		},
		{
			"address": "comdex126nshq8sq96vknl0arksys66u93spnm62qnzwc",
			"reward": {
				"denom": "ucmdx",
				"amount": "7647"
			}
		},
		{
			"address": "comdex1265ft802n6jnekcskvymr5tcnkawdzdfv58mp0",
			"reward": {
				"denom": "ucmdx",
				"amount": "3269"
			}
		},
		{
			"address": "comdex1264pndegzvk39v4vjqa9jdmxxy2vkur5tg4252",
			"reward": {
				"denom": "ucmdx",
				"amount": "5665"
			}
		},
		{
			"address": "comdex126hjru9fw7ss000c2zq2w3mcquep6xph2dg6mf",
			"reward": {
				"denom": "ucmdx",
				"amount": "509"
			}
		},
		{
			"address": "comdex126u6y3m9rqpl4gp3klnclws27x3hp2xw765hhp",
			"reward": {
				"denom": "ucmdx",
				"amount": "62"
			}
		},
		{
			"address": "comdex12mqgfxh25tqrl4xz39z0t67r4q024c7q47kdtj",
			"reward": {
				"denom": "ucmdx",
				"amount": "2015"
			}
		},
		{
			"address": "comdex12mxwg07t2r02xwlvkvqv3kqp4ep3q9p4dq8dcx",
			"reward": {
				"denom": "ucmdx",
				"amount": "1602"
			}
		},
		{
			"address": "comdex12mx0yxyrpdxg0h6vtqw8t9k65as6v26xqapgpz",
			"reward": {
				"denom": "ucmdx",
				"amount": "4195"
			}
		},
		{
			"address": "comdex12mt8xhkhjjr0a2prmcdwyca6z58wu5dcgzv05h",
			"reward": {
				"denom": "ucmdx",
				"amount": "140925"
			}
		},
		{
			"address": "comdex12mt3pg5p9g6ur4vfxuaedkj7yjd8vhdgxj4vun",
			"reward": {
				"denom": "ucmdx",
				"amount": "14170"
			}
		},
		{
			"address": "comdex12mvlwequdm23ltcl4utawkfm0ezxt3xhvgfmwf",
			"reward": {
				"denom": "ucmdx",
				"amount": "168"
			}
		},
		{
			"address": "comdex12mdekmt8f9vnrh4ldhmhspsrvlyus68qjnj9kq",
			"reward": {
				"denom": "ucmdx",
				"amount": "27318"
			}
		},
		{
			"address": "comdex12mnrdxfpa5rcrp8tq0q0rm42lphpc5n4vfzhu9",
			"reward": {
				"denom": "ucmdx",
				"amount": "38081"
			}
		},
		{
			"address": "comdex12me94waxeh974rujzndmhnmu407675tn35cyw9",
			"reward": {
				"denom": "ucmdx",
				"amount": "16971"
			}
		},
		{
			"address": "comdex12mmzge4zaxuj29ckgk2ergpmzr0fwew5ep3hyw",
			"reward": {
				"denom": "ucmdx",
				"amount": "699"
			}
		},
		{
			"address": "comdex12mud9wh0te2pjms7xhccwdx0cwmrdqjcuzs6wt",
			"reward": {
				"denom": "ucmdx",
				"amount": "1364"
			}
		},
		{
			"address": "comdex12m7yfznlwvark5ghj3dmchuefnqdr05ngrfn7k",
			"reward": {
				"denom": "ucmdx",
				"amount": "173"
			}
		},
		{
			"address": "comdex12uqvfm7xh86yjxxage05ev9mgl6mnahd8d4zx8",
			"reward": {
				"denom": "ucmdx",
				"amount": "190"
			}
		},
		{
			"address": "comdex12upz8dny2gn93qmjl5zyr9gqwqqlr6e6zcsjxl",
			"reward": {
				"denom": "ucmdx",
				"amount": "5104"
			}
		},
		{
			"address": "comdex12uzy958ffwv3x243d7hv8c038q85a7sw3v5qfv",
			"reward": {
				"denom": "ucmdx",
				"amount": "1758"
			}
		},
		{
			"address": "comdex12ux2p4sa272pywd28gt58mwszj8rf95wk066cc",
			"reward": {
				"denom": "ucmdx",
				"amount": "26181"
			}
		},
		{
			"address": "comdex12ugqg9gpxlfukygjzrhuqlujpupd83hf2ypx3q",
			"reward": {
				"denom": "ucmdx",
				"amount": "582"
			}
		},
		{
			"address": "comdex12ugw06txgx5ssc0x34qr5ecv4l32hpnp4chj35",
			"reward": {
				"denom": "ucmdx",
				"amount": "1526"
			}
		},
		{
			"address": "comdex12ugk0pqj27t6t9agdcwfa906f2xy53jh6emayw",
			"reward": {
				"denom": "ucmdx",
				"amount": "174"
			}
		},
		{
			"address": "comdex12u2sw5ymz2hvgu30fkztylgtyen3fa0yt4k5qf",
			"reward": {
				"denom": "ucmdx",
				"amount": "4699"
			}
		},
		{
			"address": "comdex12utwps9zawwuerkkwrsywd2maef9w577040z24",
			"reward": {
				"denom": "ucmdx",
				"amount": "166"
			}
		},
		{
			"address": "comdex12ud5hpcuv2n2hghay08kqx4pp9prepfmrrp0es",
			"reward": {
				"denom": "ucmdx",
				"amount": "697"
			}
		},
		{
			"address": "comdex12u3d6d7sy8j6z3e8afzayd076psmndv908jm4r",
			"reward": {
				"denom": "ucmdx",
				"amount": "1313"
			}
		},
		{
			"address": "comdex12uk22nzee0hgahzttujcdce78ax627asqp62j4",
			"reward": {
				"denom": "ucmdx",
				"amount": "96455"
			}
		},
		{
			"address": "comdex12uh42d0tqt35cac9g35mflj2smegaly5syx4lp",
			"reward": {
				"denom": "ucmdx",
				"amount": "5810"
			}
		},
		{
			"address": "comdex12ucns0x28kmn0d6nggqhrhlau2urjk48nnn4ds",
			"reward": {
				"denom": "ucmdx",
				"amount": "26631"
			}
		},
		{
			"address": "comdex12u6fq4fpczph2rvwkpa7unjgl7fgcsmmhj62ah",
			"reward": {
				"denom": "ucmdx",
				"amount": "3346"
			}
		},
		{
			"address": "comdex12u677u3z8e2k2tva6k2nz8dmpjz5xw86r2zvrn",
			"reward": {
				"denom": "ucmdx",
				"amount": "204"
			}
		},
		{
			"address": "comdex12um7rctt9kslhg0a09k2agwv2emsqdqkhtqm7f",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex12ulhqyw9awsf4ddzzkn87zhx77y0m6gth8pc4d",
			"reward": {
				"denom": "ucmdx",
				"amount": "17259"
			}
		},
		{
			"address": "comdex12aq6kj6vx5hgkvxm8mpwyhyah9rxjx8mc8hrvz",
			"reward": {
				"denom": "ucmdx",
				"amount": "166"
			}
		},
		{
			"address": "comdex12ap2h4eglahd7qn4ynezqycw6u94uyzufy48x5",
			"reward": {
				"denom": "ucmdx",
				"amount": "2948"
			}
		},
		{
			"address": "comdex12aphhxhxwa083z89sjlg7wsdjwk7rv8ze7uzsd",
			"reward": {
				"denom": "ucmdx",
				"amount": "71681"
			}
		},
		{
			"address": "comdex12ax8kp68k6cw0jq84pch6kzr7u7nf034eswkvq",
			"reward": {
				"denom": "ucmdx",
				"amount": "5255"
			}
		},
		{
			"address": "comdex12ax4nsc37wlme542klnyadnwm8v75px43xpr56",
			"reward": {
				"denom": "ucmdx",
				"amount": "13585"
			}
		},
		{
			"address": "comdex12ag4k3sahlk369y724h75k6jpa8tjpjk9n5c0g",
			"reward": {
				"denom": "ucmdx",
				"amount": "15"
			}
		},
		{
			"address": "comdex12adakf8lcn5fehxutfy52c599h70jzssk3qusg",
			"reward": {
				"denom": "ucmdx",
				"amount": "198"
			}
		},
		{
			"address": "comdex12aj3x6ecfus0n9antc0mx3prqahpc86ud0jvzp",
			"reward": {
				"denom": "ucmdx",
				"amount": "13626"
			}
		},
		{
			"address": "comdex12anlve3z6f45sahheskdedlrhy7knrjf53lwja",
			"reward": {
				"denom": "ucmdx",
				"amount": "1750"
			}
		},
		{
			"address": "comdex12acxelra3ur7x4awh7y2gx3semvsk0dkq7ge0l",
			"reward": {
				"denom": "ucmdx",
				"amount": "39965"
			}
		},
		{
			"address": "comdex12acclawexc9p4s37xezux2kkqxrpf9a0e2pmqn",
			"reward": {
				"denom": "ucmdx",
				"amount": "1438"
			}
		},
		{
			"address": "comdex12aeug8dgr0dhgmf7nafmpm2d9995ncpwcg7h2q",
			"reward": {
				"denom": "ucmdx",
				"amount": "17238"
			}
		},
		{
			"address": "comdex12ael77hffzlc43ur8vpuh838e4zna3srcvdu6p",
			"reward": {
				"denom": "ucmdx",
				"amount": "103"
			}
		},
		{
			"address": "comdex12amy64kkzqdgzpv0khz7f2y9s5ep8gpt55sjhr",
			"reward": {
				"denom": "ucmdx",
				"amount": "30"
			}
		},
		{
			"address": "comdex12auy75s8c9f0rz46glvgtnkdl4yvnllswq09x8",
			"reward": {
				"denom": "ucmdx",
				"amount": "1957"
			}
		},
		{
			"address": "comdex127phrn949370sdlg26vv5hqtjprr87yteauk2l",
			"reward": {
				"denom": "ucmdx",
				"amount": "34805"
			}
		},
		{
			"address": "comdex127yusg0086xkzagduxtlfy5lvep6h2h7vvxw7k",
			"reward": {
				"denom": "ucmdx",
				"amount": "1759"
			}
		},
		{
			"address": "comdex1279zrl88kuzjz5u8huzmulky9e3akqwakun2l2",
			"reward": {
				"denom": "ucmdx",
				"amount": "938"
			}
		},
		{
			"address": "comdex127xf5z5wdpwke0cmqvav4f0paxt6ygllc9gg7h",
			"reward": {
				"denom": "ucmdx",
				"amount": "17557"
			}
		},
		{
			"address": "comdex1278jawm55hux8ekq5evkugzr59vchdwczjf57l",
			"reward": {
				"denom": "ucmdx",
				"amount": "10723"
			}
		},
		{
			"address": "comdex127w97gr28rnzs5rk3mqg63cgnv9ezvjhcj9eew",
			"reward": {
				"denom": "ucmdx",
				"amount": "3343"
			}
		},
		{
			"address": "comdex127we646lz9h9jn3n8yucw7cwrd0ghzwmckgsqn",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1270yutqk8wc0w6ut92djlgcnw3gz7qwg2zpx8r",
			"reward": {
				"denom": "ucmdx",
				"amount": "393"
			}
		},
		{
			"address": "comdex127sf7va0q7awqf48mkfdt560t3yn7xkx9jtx7m",
			"reward": {
				"denom": "ucmdx",
				"amount": "37537"
			}
		},
		{
			"address": "comdex1273vkyljgnuxdehfcdqr33xt93ufk7jk78xhh0",
			"reward": {
				"denom": "ucmdx",
				"amount": "11638"
			}
		},
		{
			"address": "comdex12736ed7frh0uxwv2tm29qz8v30fe03kqevynea",
			"reward": {
				"denom": "ucmdx",
				"amount": "15"
			}
		},
		{
			"address": "comdex1273772k35v4ye7pp3zdzlfwgph68fl694r6xrc",
			"reward": {
				"denom": "ucmdx",
				"amount": "1494"
			}
		},
		{
			"address": "comdex127hgjjrst9mngejd4l4wprnnppwl6e22agnrjh",
			"reward": {
				"denom": "ucmdx",
				"amount": "9282"
			}
		},
		{
			"address": "comdex127mw4n5lp0z6jlyj0r9m6nmu5xht4ktegnlqlh",
			"reward": {
				"denom": "ucmdx",
				"amount": "1746"
			}
		},
		{
			"address": "comdex127uhn5mlkglw94gn8qvu7zv50pdz55vqqdu6ht",
			"reward": {
				"denom": "ucmdx",
				"amount": "13439"
			}
		},
		{
			"address": "comdex12lzyps03x3p58534j3nc223hnyngqhpj33jjkt",
			"reward": {
				"denom": "ucmdx",
				"amount": "1775"
			}
		},
		{
			"address": "comdex12ld097zdqyntcfhxcpepuzwd25dt7pk6agtcm0",
			"reward": {
				"denom": "ucmdx",
				"amount": "1268"
			}
		},
		{
			"address": "comdex12ld6paptu5zyen9gpqlaqx0l2r7m0989gp59au",
			"reward": {
				"denom": "ucmdx",
				"amount": "3857"
			}
		},
		{
			"address": "comdex12lwyk347px4n3mcyqxtafaq0234nc83qur9swh",
			"reward": {
				"denom": "ucmdx",
				"amount": "1252"
			}
		},
		{
			"address": "comdex12l0ns2c834uuq4q7wh4yh748sel8fgamyldk58",
			"reward": {
				"denom": "ucmdx",
				"amount": "1250"
			}
		},
		{
			"address": "comdex12lhfk7apm6ecvnkwzlmjgp62lh8pph85rcnyf5",
			"reward": {
				"denom": "ucmdx",
				"amount": "12361"
			}
		},
		{
			"address": "comdex12lh70wrj9gj5wvuxfq5e6crtmk3jzsmxzj0x6f",
			"reward": {
				"denom": "ucmdx",
				"amount": "14941"
			}
		},
		{
			"address": "comdex12luppuxdlw4qev2mkzv8x6z95xs9z3xfz8wtq4",
			"reward": {
				"denom": "ucmdx",
				"amount": "32798"
			}
		},
		{
			"address": "comdex12lan46dtm0f6lkaaskxz8vqqn0eaxw8sjnug29",
			"reward": {
				"denom": "ucmdx",
				"amount": "3018"
			}
		},
		{
			"address": "comdex12llvmd4tavn54lquxcl4yldw3yvrgmlvw25xfl",
			"reward": {
				"denom": "ucmdx",
				"amount": "14781"
			}
		},
		{
			"address": "comdex12llk842x8c8u34ay29npmpa8cd7z6lkx2x6wz5",
			"reward": {
				"denom": "ucmdx",
				"amount": "36173"
			}
		},
		{
			"address": "comdex1tqq70pyy5j9reeh9jsu60kt9mzc6w6p92e4k7z",
			"reward": {
				"denom": "ucmdx",
				"amount": "150"
			}
		},
		{
			"address": "comdex1tqrfmdymqlr7lkprs0pr8ursesends7nlgsk9s",
			"reward": {
				"denom": "ucmdx",
				"amount": "1191"
			}
		},
		{
			"address": "comdex1tqyz2ny0jyugw224gx99wc7panh4u9j3z0eezy",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1tq9rfnr0pmq005mwt8qdmh5fhfjaf8w5zj8gj5",
			"reward": {
				"denom": "ucmdx",
				"amount": "30464"
			}
		},
		{
			"address": "comdex1tq93vps908y7lvue5dxp7dl73msqpw52vcdmcv",
			"reward": {
				"denom": "ucmdx",
				"amount": "3279"
			}
		},
		{
			"address": "comdex1tq2spklg3ss44cjxamhggru4utt3mw7hf37l30",
			"reward": {
				"denom": "ucmdx",
				"amount": "399"
			}
		},
		{
			"address": "comdex1tqvqdqhmr7sx6087yhfvfkrq3e8c907egpddzs",
			"reward": {
				"denom": "ucmdx",
				"amount": "28765"
			}
		},
		{
			"address": "comdex1tqdvlku97p88z5ndf3jz55vnccmzgs8jsfqvmx",
			"reward": {
				"denom": "ucmdx",
				"amount": "1338"
			}
		},
		{
			"address": "comdex1tq0lyuhxh8qqzxcz8rzywyd3kvdcd3ts9vsv2r",
			"reward": {
				"denom": "ucmdx",
				"amount": "1397"
			}
		},
		{
			"address": "comdex1tqskstpdcptrdhyrnwu8ngdu4vf53m067nxr8t",
			"reward": {
				"denom": "ucmdx",
				"amount": "754"
			}
		},
		{
			"address": "comdex1tqnv9a2w3e6gps8497fwdmc02l3hugjq0sn9pl",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1tq5qwc8zc0j2n9cfsr4gunvhjznjyvj5aw0xpm",
			"reward": {
				"denom": "ucmdx",
				"amount": "123"
			}
		},
		{
			"address": "comdex1tq4wvp6xlager9gdwvkwuyf5s8gup7g8gc9vc2",
			"reward": {
				"denom": "ucmdx",
				"amount": "28"
			}
		},
		{
			"address": "comdex1tq4sva3h5p6aum2ptfz5qjwqjynf7xrnw7vsqj",
			"reward": {
				"denom": "ucmdx",
				"amount": "1408"
			}
		},
		{
			"address": "comdex1tqhfx347aw28jve4nx3tw7w7jl4wyu3n70eklz",
			"reward": {
				"denom": "ucmdx",
				"amount": "598"
			}
		},
		{
			"address": "comdex1tqaezx3ezur5rcecqcqkhfve4seqjlj7xre4x6",
			"reward": {
				"denom": "ucmdx",
				"amount": "6177"
			}
		},
		{
			"address": "comdex1tqa7kt5vkzg3x5m64ky8x2lrv2uvlqmshgn4vu",
			"reward": {
				"denom": "ucmdx",
				"amount": "97004"
			}
		},
		{
			"address": "comdex1tppr3nqdy8yxw4u63pa8gskc4zc4gucrfc3hxg",
			"reward": {
				"denom": "ucmdx",
				"amount": "70559"
			}
		},
		{
			"address": "comdex1tpzw7sdwvdgul9uz80zvey3tp4qx4njavlm775",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1tpydrwlpzkhy8jzaj7yg0qklnvpflu6pwsp8ar",
			"reward": {
				"denom": "ucmdx",
				"amount": "20621"
			}
		},
		{
			"address": "comdex1tp9scplfeheksreg6mcehnl4lccqqakcull6s2",
			"reward": {
				"denom": "ucmdx",
				"amount": "1767"
			}
		},
		{
			"address": "comdex1tp827vsfegmghz3n9dpw6rgrevgesydcxv93y6",
			"reward": {
				"denom": "ucmdx",
				"amount": "705"
			}
		},
		{
			"address": "comdex1tpthjg3qx0ht4nqzcyu662n90qhekxff0kr6yy",
			"reward": {
				"denom": "ucmdx",
				"amount": "8957"
			}
		},
		{
			"address": "comdex1tpkk5sj5sek2h2k0a0mqd8lnk4y96vj4jcevsf",
			"reward": {
				"denom": "ucmdx",
				"amount": "7194"
			}
		},
		{
			"address": "comdex1tph4w98zcv4k6c9cc6x8qug5kfqjphmzvfgald",
			"reward": {
				"denom": "ucmdx",
				"amount": "6895"
			}
		},
		{
			"address": "comdex1tp6fa7pfkpc6m8dyhwm8rx64h0ef8xmvf476ey",
			"reward": {
				"denom": "ucmdx",
				"amount": "195"
			}
		},
		{
			"address": "comdex1tpmrwgmzsxmppjkxfqx3za9sqh7weffxleywah",
			"reward": {
				"denom": "ucmdx",
				"amount": "19439"
			}
		},
		{
			"address": "comdex1tpmxu3qekz6h9mgvulcutkg42p7wnga0j7jj65",
			"reward": {
				"denom": "ucmdx",
				"amount": "4358"
			}
		},
		{
			"address": "comdex1tpm8xez3p9wd3xwpter9387qlqn3tpxremdq0z",
			"reward": {
				"denom": "ucmdx",
				"amount": "28"
			}
		},
		{
			"address": "comdex1tpu3x05006ntwkc0c8wrpvzdjz0mt5wh8694kh",
			"reward": {
				"denom": "ucmdx",
				"amount": "199"
			}
		},
		{
			"address": "comdex1tzqd0wt36l5euw66p356h30xtfcry0lc5dmuks",
			"reward": {
				"denom": "ucmdx",
				"amount": "15"
			}
		},
		{
			"address": "comdex1tzz9mqtfz28yuajgx3jspzx7ctp7y6dggzpag8",
			"reward": {
				"denom": "ucmdx",
				"amount": "1501"
			}
		},
		{
			"address": "comdex1tzfs83qa7k9ycufvepckxjpkgjmcq9wyjmtfjy",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1tz2n9wu5lpvm2c02hfn22s25auel5gckuahtyk",
			"reward": {
				"denom": "ucmdx",
				"amount": "1264"
			}
		},
		{
			"address": "comdex1tzdhtz7k04q49xrprrzxlelfq9w2cvskmqvzmt",
			"reward": {
				"denom": "ucmdx",
				"amount": "16482"
			}
		},
		{
			"address": "comdex1tzs9kr6uq4jxswfrlq4zddyp6k2q8ztad5qxhf",
			"reward": {
				"denom": "ucmdx",
				"amount": "6184"
			}
		},
		{
			"address": "comdex1tzs2mrp9fkwq6ydtlzlrzghm6ktqccyvwvxafa",
			"reward": {
				"denom": "ucmdx",
				"amount": "354"
			}
		},
		{
			"address": "comdex1tz4ssc67mz09pyc5xtvt3zvkm8g6uggkjdtvc5",
			"reward": {
				"denom": "ucmdx",
				"amount": "13"
			}
		},
		{
			"address": "comdex1tza8u8l4mqjyy8pn8ryssl52aqs2e4y3vya5y3",
			"reward": {
				"denom": "ucmdx",
				"amount": "264"
			}
		},
		{
			"address": "comdex1tzamymwq34kgx7jr0q0cghx32uv4tygp73crze",
			"reward": {
				"denom": "ucmdx",
				"amount": "855"
			}
		},
		{
			"address": "comdex1tzll9wxf3w67mghcec5zlmxv52lwhz6ftvrt25",
			"reward": {
				"denom": "ucmdx",
				"amount": "11974"
			}
		},
		{
			"address": "comdex1trqknkspvcwk6f43q6qzxzz37rpqrnyp0v7876",
			"reward": {
				"denom": "ucmdx",
				"amount": "6754"
			}
		},
		{
			"address": "comdex1trqlp27j27seeeux88tde7acal3k9vv8k7dne4",
			"reward": {
				"denom": "ucmdx",
				"amount": "16731"
			}
		},
		{
			"address": "comdex1trpqycqrjm2pzys6v065p35gej7qzap0m6g5xf",
			"reward": {
				"denom": "ucmdx",
				"amount": "15007"
			}
		},
		{
			"address": "comdex1trrgphjtst6glc6vpk4hd8qur70ajersuwqrhr",
			"reward": {
				"denom": "ucmdx",
				"amount": "569"
			}
		},
		{
			"address": "comdex1trygc5hp2djlw96d8arv5kmgt2wtn9xmzp79rt",
			"reward": {
				"denom": "ucmdx",
				"amount": "2108"
			}
		},
		{
			"address": "comdex1trxd92uxfgjpeamgxgwtn627eu0jaanqycly0k",
			"reward": {
				"denom": "ucmdx",
				"amount": "142"
			}
		},
		{
			"address": "comdex1tr8em90ramvug087r9husz704kp9h6vht84due",
			"reward": {
				"denom": "ucmdx",
				"amount": "14362"
			}
		},
		{
			"address": "comdex1trfh22gfvj3dfvr0lvvwexdu94xx67dhm9mk34",
			"reward": {
				"denom": "ucmdx",
				"amount": "4047"
			}
		},
		{
			"address": "comdex1trtadzag7e9zwg96m7pp3e68mp8ceh5rq3qc72",
			"reward": {
				"denom": "ucmdx",
				"amount": "2017"
			}
		},
		{
			"address": "comdex1trdjj8528dn26aj34k3qf2d86ddcsqhk8c32m0",
			"reward": {
				"denom": "ucmdx",
				"amount": "7461"
			}
		},
		{
			"address": "comdex1tr0pcc2p080t4j6u943q0g5ke0n3nlt609hdxz",
			"reward": {
				"denom": "ucmdx",
				"amount": "888"
			}
		},
		{
			"address": "comdex1tr523g053jer0wyg8lqj6nj38zkatcp59h2jaw",
			"reward": {
				"denom": "ucmdx",
				"amount": "8314"
			}
		},
		{
			"address": "comdex1tr55zghf7d6rh3l4qfjhs4p4la2pz79zqztdd2",
			"reward": {
				"denom": "ucmdx",
				"amount": "623280"
			}
		},
		{
			"address": "comdex1trhw34pzc2wrkkxasp8hhwkaly9cau8zdlz6jx",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1trhe5jgrmzpvjgtkx3lvdhfq0v5uls9ma9cfg3",
			"reward": {
				"denom": "ucmdx",
				"amount": "392"
			}
		},
		{
			"address": "comdex1treczrfhsjdvlfvcgldt7d94wj4vhk3r97km9f",
			"reward": {
				"denom": "ucmdx",
				"amount": "661737"
			}
		},
		{
			"address": "comdex1tr6najckswy2krca8ju7tqyamelcrmqzr53rk5",
			"reward": {
				"denom": "ucmdx",
				"amount": "24206"
			}
		},
		{
			"address": "comdex1trmwvl23g2msqvd0w3n9reaxgs3kwcgv9y45f6",
			"reward": {
				"denom": "ucmdx",
				"amount": "4468"
			}
		},
		{
			"address": "comdex1trm0mjpufne93rqrwavxklftdguuxngefs8ayk",
			"reward": {
				"denom": "ucmdx",
				"amount": "18365"
			}
		},
		{
			"address": "comdex1trm32mfzpgd9e696xujpauu45xrg4hdwtnwwvl",
			"reward": {
				"denom": "ucmdx",
				"amount": "143"
			}
		},
		{
			"address": "comdex1truhswn48uzny068ts2vmtsh7pj966mat4lfxz",
			"reward": {
				"denom": "ucmdx",
				"amount": "15361"
			}
		},
		{
			"address": "comdex1trlc2zla7u5r0j9gxqs6zhec6gylclg6xp4fpn",
			"reward": {
				"denom": "ucmdx",
				"amount": "15"
			}
		},
		{
			"address": "comdex1tyqunz6gtkt3hml2x9dg60ue4kk634k3q2v7ae",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1typ45xcrpnlt2j5rzc8pmhfp4u6pglm00g5sje",
			"reward": {
				"denom": "ucmdx",
				"amount": "21301"
			}
		},
		{
			"address": "comdex1tyzhzqfyw6y98zlxstvnwcs5cu6zx8uf398pg8",
			"reward": {
				"denom": "ucmdx",
				"amount": "12331"
			}
		},
		{
			"address": "comdex1tyr8lskch6mc098f5u5t3dqvq5wyplc98amajn",
			"reward": {
				"denom": "ucmdx",
				"amount": "6471"
			}
		},
		{
			"address": "comdex1tyxhgmztx8m9jpp2qsvdz03cvgfqtu785rjl98",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1tygmpjmq0vw4h6m5ht90vemx6eghjvfk4q82a5",
			"reward": {
				"denom": "ucmdx",
				"amount": "9414"
			}
		},
		{
			"address": "comdex1tytuzsk6ztxs00qls3554ut07t8j2kq4x9e6uj",
			"reward": {
				"denom": "ucmdx",
				"amount": "27358"
			}
		},
		{
			"address": "comdex1tydgc75vyyqnzzcc79fufkenvuwjc3hcgnyqe2",
			"reward": {
				"denom": "ucmdx",
				"amount": "1992"
			}
		},
		{
			"address": "comdex1ty0kk5nm5wgc2fpf53ns7wfzqk42lmgmvtkcn3",
			"reward": {
				"denom": "ucmdx",
				"amount": "16880"
			}
		},
		{
			"address": "comdex1tynl4ra53v5upf4de88uqp3ncgyymfea8h6ulu",
			"reward": {
				"denom": "ucmdx",
				"amount": "1627"
			}
		},
		{
			"address": "comdex1tyk33uzuen03du58sh5vq5n2s49deecq6dvzfk",
			"reward": {
				"denom": "ucmdx",
				"amount": "456"
			}
		},
		{
			"address": "comdex1tyk7gn408xx3mscqy777lgnyrx2whuh4yreee6",
			"reward": {
				"denom": "ucmdx",
				"amount": "14666"
			}
		},
		{
			"address": "comdex1tyhpnfk25gkvg3gpw5f9f54s9xpgprn3mlepas",
			"reward": {
				"denom": "ucmdx",
				"amount": "2106"
			}
		},
		{
			"address": "comdex1tyhz04rp93x2uvat3mpd67pkdgsgdrym283mu7",
			"reward": {
				"denom": "ucmdx",
				"amount": "7194"
			}
		},
		{
			"address": "comdex1tycqlmve4m8j0lje6krgfw0nf68qfrkup6dyq9",
			"reward": {
				"denom": "ucmdx",
				"amount": "4285"
			}
		},
		{
			"address": "comdex1ty6q0dlm0awe9pjzukt0e5c49xr0hx85k5c94v",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1tylk7dce32nlqr0uk80hffq9ky6g89xewthvqk",
			"reward": {
				"denom": "ucmdx",
				"amount": "1274"
			}
		},
		{
			"address": "comdex1t9yhdcaagx2l2ftmyfqhz84m859p3z7tljsdd2",
			"reward": {
				"denom": "ucmdx",
				"amount": "123"
			}
		},
		{
			"address": "comdex1t9xcjprzl0gjnkfh9u7d8ypsd8pkqrn3m4p7kz",
			"reward": {
				"denom": "ucmdx",
				"amount": "1353"
			}
		},
		{
			"address": "comdex1t9vd7l9a2yv0ym6mz5plqzcpt2pa8reew68qum",
			"reward": {
				"denom": "ucmdx",
				"amount": "7102"
			}
		},
		{
			"address": "comdex1t9dfvx63de4dey3ffg2vq354677jq2pghy6twq",
			"reward": {
				"denom": "ucmdx",
				"amount": "1393"
			}
		},
		{
			"address": "comdex1t9d6hckvr2q2eghzhtgkyvef8yz8mn8p5ks2av",
			"reward": {
				"denom": "ucmdx",
				"amount": "1774"
			}
		},
		{
			"address": "comdex1t9592v0zyvgqew7lfl3lmn9hhrkl8jrn5swh9q",
			"reward": {
				"denom": "ucmdx",
				"amount": "15"
			}
		},
		{
			"address": "comdex1t95l7vvlhjj23wk5573ktvmkv0v7tr5z3ahwqf",
			"reward": {
				"denom": "ucmdx",
				"amount": "1600"
			}
		},
		{
			"address": "comdex1t9e2s4cftjkqx282egju63fqw8f952kduh3qqt",
			"reward": {
				"denom": "ucmdx",
				"amount": "24117"
			}
		},
		{
			"address": "comdex1t9ewqryptggnmrddw2g56rfuqp8wf0y7lg4chx",
			"reward": {
				"denom": "ucmdx",
				"amount": "200"
			}
		},
		{
			"address": "comdex1t96ds9s985cpwtyayykgez6qnrwec9cxvf2yrt",
			"reward": {
				"denom": "ucmdx",
				"amount": "22275"
			}
		},
		{
			"address": "comdex1t96mzu8wtuqh9tt4c48gx0ayrf5y4xqm4r0dz0",
			"reward": {
				"denom": "ucmdx",
				"amount": "7"
			}
		},
		{
			"address": "comdex1t96ufr2r3vfnzy4xsspjaxsgmneurctcepn3uz",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1t96axfycj9uvth4djck2j3ckzfrhwl3acneyjp",
			"reward": {
				"denom": "ucmdx",
				"amount": "7038"
			}
		},
		{
			"address": "comdex1t9u0hsqx8au3a3jwxnzq2rr87n6adzs03trxls",
			"reward": {
				"denom": "ucmdx",
				"amount": "1740"
			}
		},
		{
			"address": "comdex1t97nkssxvd77wy9a7qhcvhrzjjh32pt9vtwdh3",
			"reward": {
				"denom": "ucmdx",
				"amount": "1489"
			}
		},
		{
			"address": "comdex1t9lyln4j3748jpvgfwdw82sp4tw5gytfnt45e9",
			"reward": {
				"denom": "ucmdx",
				"amount": "122"
			}
		},
		{
			"address": "comdex1txruxj97ettt07czusg76kptc6k4w2575g5f9d",
			"reward": {
				"denom": "ucmdx",
				"amount": "122154"
			}
		},
		{
			"address": "comdex1tx9d4zg90uxwxyu32n2m06sftse92sk59u2wxz",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1txxadmf97vhx6972k7nyspr0al072ncnw5l46z",
			"reward": {
				"denom": "ucmdx",
				"amount": "132"
			}
		},
		{
			"address": "comdex1tx80zqvylpykr90swwkjnnlhtj394gn0zgdfmf",
			"reward": {
				"denom": "ucmdx",
				"amount": "1160"
			}
		},
		{
			"address": "comdex1tx2s8744nanf4wjfdsnj5um0gf0hj9v34jh6ln",
			"reward": {
				"denom": "ucmdx",
				"amount": "169"
			}
		},
		{
			"address": "comdex1txvtlvhj9xgmyghkrjxp2dgzfj53rd4cmh8crz",
			"reward": {
				"denom": "ucmdx",
				"amount": "176"
			}
		},
		{
			"address": "comdex1tx3yt8c8mw20z38eaxjwzs77ajtlrc84g2wqj4",
			"reward": {
				"denom": "ucmdx",
				"amount": "2044"
			}
		},
		{
			"address": "comdex1tx3ladcmy45schan7tp5qgz5kwhdzjvjp3pskr",
			"reward": {
				"denom": "ucmdx",
				"amount": "87"
			}
		},
		{
			"address": "comdex1txns8n6l855kl9rd9a7pqzr9gc4xdumpjrv34n",
			"reward": {
				"denom": "ucmdx",
				"amount": "7018"
			}
		},
		{
			"address": "comdex1txcenjsjk4zn2x7jp29tr0w0a09chtqfnjwtg8",
			"reward": {
				"denom": "ucmdx",
				"amount": "6658"
			}
		},
		{
			"address": "comdex1tx62yqx7qvpzuve3mu9mha2zu52q093dn4eqa2",
			"reward": {
				"denom": "ucmdx",
				"amount": "16788"
			}
		},
		{
			"address": "comdex1tx72zrx86r6eegndpgqqejz62880t9qssgzg00",
			"reward": {
				"denom": "ucmdx",
				"amount": "14997"
			}
		},
		{
			"address": "comdex1t8fl58e48p8tyrcmj479tv6ljuhu8askchg332",
			"reward": {
				"denom": "ucmdx",
				"amount": "7200"
			}
		},
		{
			"address": "comdex1t82haf5c448nqhfdnjhpn4lfv36zrtzyldsa49",
			"reward": {
				"denom": "ucmdx",
				"amount": "173"
			}
		},
		{
			"address": "comdex1t8v354fhc940wuzw8s639aqcx3x7mla0lc9g9x",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1t8w4k3xlmwy36kmlc382z2fdhydwlvc2q6t84t",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1t8sqv6ajetuyukr5v9pl8u6xpslzptadjeshyt",
			"reward": {
				"denom": "ucmdx",
				"amount": "32407"
			}
		},
		{
			"address": "comdex1t84t30dp436z8zgmpupd0qx4rruwmtplw20twj",
			"reward": {
				"denom": "ucmdx",
				"amount": "11647"
			}
		},
		{
			"address": "comdex1t84h2gyq5j25pdr7ktkk8q7dxha23z5ypyz0dk",
			"reward": {
				"denom": "ucmdx",
				"amount": "19"
			}
		},
		{
			"address": "comdex1t8kk4t9rgx8v2n2etclm5pe9ff6kce94efc077",
			"reward": {
				"denom": "ucmdx",
				"amount": "12375"
			}
		},
		{
			"address": "comdex1t8khvdqud2su8vagq647mdshgkt5zdr3q8t2sx",
			"reward": {
				"denom": "ucmdx",
				"amount": "2069"
			}
		},
		{
			"address": "comdex1t8hxac9ahpzslsn4nyk20l0lrvfwcas22swch0",
			"reward": {
				"denom": "ucmdx",
				"amount": "5385"
			}
		},
		{
			"address": "comdex1t8c83d9t44t829pzjq9uxc9aq5kzf64xe5rlcf",
			"reward": {
				"denom": "ucmdx",
				"amount": "2484"
			}
		},
		{
			"address": "comdex1t8cdk3gx6nrq0tn7r5e55ry9mdmgcm4vll72am",
			"reward": {
				"denom": "ucmdx",
				"amount": "47788"
			}
		},
		{
			"address": "comdex1t8ckqpnv5lp42qk2rzduuet6qlxrvmghgfhu3u",
			"reward": {
				"denom": "ucmdx",
				"amount": "6370"
			}
		},
		{
			"address": "comdex1t8clt3lf5n9awdv8m3g4fv7nnczqp6hsjzvpzl",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1t8e5vktvrtwxn98tfr3k3esl424s2jm6jk378z",
			"reward": {
				"denom": "ucmdx",
				"amount": "1770"
			}
		},
		{
			"address": "comdex1t8mmjdp0hxr6c9mgdeanuuet0zer0zeetjvuh6",
			"reward": {
				"denom": "ucmdx",
				"amount": "73103"
			}
		},
		{
			"address": "comdex1t8u5nkyhh5gsx9slm5j66vhahypyp9kwpdldt5",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1t8ale0jj3hwdyajdukazupfltad9eyq7jn3key",
			"reward": {
				"denom": "ucmdx",
				"amount": "152"
			}
		},
		{
			"address": "comdex1t87cmj32xx0up3ect9deevprsdhdnh6n25235g",
			"reward": {
				"denom": "ucmdx",
				"amount": "14517"
			}
		},
		{
			"address": "comdex1t8lvk5hyjukgk4s8yk4tm7d6fgr34xpn9n8p2r",
			"reward": {
				"denom": "ucmdx",
				"amount": "1354"
			}
		},
		{
			"address": "comdex1tgzpekk7h8d7f4dcz6qaq243allt27f0luxmkk",
			"reward": {
				"denom": "ucmdx",
				"amount": "167"
			}
		},
		{
			"address": "comdex1tg9f7xd9e4fw66y2y8srzzxys940ys2xypltkq",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1tg92wv39ry43hy8wuy7h9s63yjlefkymh9mhew",
			"reward": {
				"denom": "ucmdx",
				"amount": "2647"
			}
		},
		{
			"address": "comdex1tggjr342nd0ze0g9vtswg25cqdmsphwwjds2ag",
			"reward": {
				"denom": "ucmdx",
				"amount": "540"
			}
		},
		{
			"address": "comdex1tg20eujl62v226h0404thnu4pu42rd8v4u2g2r",
			"reward": {
				"denom": "ucmdx",
				"amount": "10373"
			}
		},
		{
			"address": "comdex1tgtwlu8l3pjwnyx8smlp3ht6cv54r4qfsq8977",
			"reward": {
				"denom": "ucmdx",
				"amount": "123394"
			}
		},
		{
			"address": "comdex1tgdjd4qzlsqnut9yyfecyde6wmr379a3avm3l8",
			"reward": {
				"denom": "ucmdx",
				"amount": "150"
			}
		},
		{
			"address": "comdex1tgwtefxwkhde6dg7vjrvjwwqs52jdp67ruce0t",
			"reward": {
				"denom": "ucmdx",
				"amount": "7161"
			}
		},
		{
			"address": "comdex1tg0yc4phhws49n48ks473zee3fstsqf67tztyy",
			"reward": {
				"denom": "ucmdx",
				"amount": "5451"
			}
		},
		{
			"address": "comdex1tg0kr6x74ytyjrx5ex74vk50ssszmlhd7fhd43",
			"reward": {
				"denom": "ucmdx",
				"amount": "4031"
			}
		},
		{
			"address": "comdex1tgjh39tkx5wmcemeq3zuql33fwuupqud7ldcqg",
			"reward": {
				"denom": "ucmdx",
				"amount": "1489"
			}
		},
		{
			"address": "comdex1tgje6kmcyfxzc08qzvfpd0404npp963ht426fy",
			"reward": {
				"denom": "ucmdx",
				"amount": "4268"
			}
		},
		{
			"address": "comdex1tgefja4etdwefqheku6leuprh9gq9mdks7pzr6",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1tgm68zns6rn46dav5lmrxx40r9xu7ganeytdu2",
			"reward": {
				"denom": "ucmdx",
				"amount": "1770"
			}
		},
		{
			"address": "comdex1tgugsdy284774w0k3lavusfvyancrgek86axvr",
			"reward": {
				"denom": "ucmdx",
				"amount": "39"
			}
		},
		{
			"address": "comdex1tgasu89e9s20k53y05e40gashu36uc6u7xy2jm",
			"reward": {
				"denom": "ucmdx",
				"amount": "14861"
			}
		},
		{
			"address": "comdex1tgl0383paqttp6k4gq2t8un33t0l66khvazp0h",
			"reward": {
				"denom": "ucmdx",
				"amount": "2641"
			}
		},
		{
			"address": "comdex1tfqmka5ycjgdyfxktr3eyd5lqsgggxx753r6tx",
			"reward": {
				"denom": "ucmdx",
				"amount": "99458"
			}
		},
		{
			"address": "comdex1tfrw29k87a5srgrpv08k99mzr2mqcmd0ws4y2e",
			"reward": {
				"denom": "ucmdx",
				"amount": "25118"
			}
		},
		{
			"address": "comdex1tfr7ttkevwpk3wmxf06u2vztlyac0jw29xrrxf",
			"reward": {
				"denom": "ucmdx",
				"amount": "7026"
			}
		},
		{
			"address": "comdex1tfx39lj79hjt2uzys4zf405x8kqg73037djhs6",
			"reward": {
				"denom": "ucmdx",
				"amount": "13636"
			}
		},
		{
			"address": "comdex1tfxc4mdk684p909hqtx9lvlkzskh5f6etzzzsq",
			"reward": {
				"denom": "ucmdx",
				"amount": "36170"
			}
		},
		{
			"address": "comdex1tfgyag83d00gtqgmqtwx33ns6n596gkthfcrcr",
			"reward": {
				"denom": "ucmdx",
				"amount": "17391"
			}
		},
		{
			"address": "comdex1tff687mc8qkhyayl3xk9trplqkhkf2748g073q",
			"reward": {
				"denom": "ucmdx",
				"amount": "2522"
			}
		},
		{
			"address": "comdex1tfdqm5d3paekqslvexsk3nw506l95ylcwwst24",
			"reward": {
				"denom": "ucmdx",
				"amount": "13547"
			}
		},
		{
			"address": "comdex1tf08p2qqpspl0qtv2w3z4h8d973j6weljjqz4l",
			"reward": {
				"denom": "ucmdx",
				"amount": "1741"
			}
		},
		{
			"address": "comdex1tf0aw7dq4w3vppfqdglefs6wzyz5um2sw7yct0",
			"reward": {
				"denom": "ucmdx",
				"amount": "28968"
			}
		},
		{
			"address": "comdex1tfsfehskwzg0j8zpxtaezzux2n0djxsguq5m2x",
			"reward": {
				"denom": "ucmdx",
				"amount": "5835"
			}
		},
		{
			"address": "comdex1tfjwfmvcacls45zrf7u7q9annl2nu96pwy99dl",
			"reward": {
				"denom": "ucmdx",
				"amount": "184"
			}
		},
		{
			"address": "comdex1tfecf5ha96lht76fx30y9lm9hqed8sjr0p5w5s",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1tfuyg9jjfygkn3a0wha84fafj59wgccmuvxa9y",
			"reward": {
				"denom": "ucmdx",
				"amount": "5308"
			}
		},
		{
			"address": "comdex1tfu9998q8r9gu46guqyywdfrfc2yvg45avjl06",
			"reward": {
				"denom": "ucmdx",
				"amount": "2280"
			}
		},
		{
			"address": "comdex1tfluv8w0aves0tapxez5axdek3kksvjeayv2uw",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1t2xcxhsxp5fa344tk558v5r02ft2y4ywqr4za6",
			"reward": {
				"denom": "ucmdx",
				"amount": "29510"
			}
		},
		{
			"address": "comdex1t2gurshlgr2e47atylfcqpqs6etj22kfz0v2e4",
			"reward": {
				"denom": "ucmdx",
				"amount": "1768"
			}
		},
		{
			"address": "comdex1t2w5dqeddvtptx8xly2e52lk789trw79shke8s",
			"reward": {
				"denom": "ucmdx",
				"amount": "1234"
			}
		},
		{
			"address": "comdex1t20235r28damv23y73f2vvwcwy3ffzt5c9rylx",
			"reward": {
				"denom": "ucmdx",
				"amount": "17745"
			}
		},
		{
			"address": "comdex1t25eagpetdt77tkf8h9vfuc3runh83n393salj",
			"reward": {
				"denom": "ucmdx",
				"amount": "112660"
			}
		},
		{
			"address": "comdex1t2k53kgn7pf2gqc6z9mefajuxr4qu7hrglny7v",
			"reward": {
				"denom": "ucmdx",
				"amount": "4624"
			}
		},
		{
			"address": "comdex1t269x6002mr0tzh3qcc4mjfv7uhxmgxkhkle7y",
			"reward": {
				"denom": "ucmdx",
				"amount": "6177"
			}
		},
		{
			"address": "comdex1t269dt0nfxa7effdg9lu27u7fffp44ke5qks7g",
			"reward": {
				"denom": "ucmdx",
				"amount": "1424"
			}
		},
		{
			"address": "comdex1t2u8vqx0ldckr7njp39zyxc3gtewz6dg9ua7tj",
			"reward": {
				"denom": "ucmdx",
				"amount": "1945"
			}
		},
		{
			"address": "comdex1ttyy7fz35e8m83ypdujvlhrr39h2eesygjcqwz",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1tt9qzhc3zrefeez07jrg5d8gm7u4xa7nazwyp6",
			"reward": {
				"denom": "ucmdx",
				"amount": "3334"
			}
		},
		{
			"address": "comdex1tt83xdje9gx99cp8757wjrxzc29nycnqx4xa8m",
			"reward": {
				"denom": "ucmdx",
				"amount": "1602"
			}
		},
		{
			"address": "comdex1tt2gyvasggkcgxdmugt2zwpev67dlvhtsqc250",
			"reward": {
				"denom": "ucmdx",
				"amount": "1748"
			}
		},
		{
			"address": "comdex1ttvafr9fl98pg496fdmagd7fswqxpy0j936rpu",
			"reward": {
				"denom": "ucmdx",
				"amount": "1022"
			}
		},
		{
			"address": "comdex1ttw88wvlfd3yaxm764zepvxyn2yf2vnalkkq4w",
			"reward": {
				"denom": "ucmdx",
				"amount": "494"
			}
		},
		{
			"address": "comdex1tt0cjj28fpx9c026maca8tz5sredksmgz3063y",
			"reward": {
				"denom": "ucmdx",
				"amount": "2544"
			}
		},
		{
			"address": "comdex1ttnlxu2kl3zz94canzxvcsvupvv8dyu098ahyz",
			"reward": {
				"denom": "ucmdx",
				"amount": "3363"
			}
		},
		{
			"address": "comdex1tteageasxpu3ltv993hzyk24t0zeea4jzffcr7",
			"reward": {
				"denom": "ucmdx",
				"amount": "1764"
			}
		},
		{
			"address": "comdex1ttelgfwnxmxaqhgppg39xdxp35q3v688gyh2wv",
			"reward": {
				"denom": "ucmdx",
				"amount": "170"
			}
		},
		{
			"address": "comdex1tt6qz0hmgmmq2q7ck6vhyx6rf6ae0xnxwfsthl",
			"reward": {
				"denom": "ucmdx",
				"amount": "92934"
			}
		},
		{
			"address": "comdex1ttmnnjw2u84azt2sqcpwk5jsr2elz8s0y0f37n",
			"reward": {
				"denom": "ucmdx",
				"amount": "180"
			}
		},
		{
			"address": "comdex1tvy2nlvl5zfm2a8nmsw3zh6tflp8z9wm4h7jhg",
			"reward": {
				"denom": "ucmdx",
				"amount": "5700"
			}
		},
		{
			"address": "comdex1tvx03r3jykwe6q5d9q5hjmy2trxwngmvaqsmd3",
			"reward": {
				"denom": "ucmdx",
				"amount": "2822"
			}
		},
		{
			"address": "comdex1tv85ggjuxnxden0v44assh854k88t9pct5w3wv",
			"reward": {
				"denom": "ucmdx",
				"amount": "373"
			}
		},
		{
			"address": "comdex1tvg27ngtxs0xqec3lr9srw3svf458j0dvhp6me",
			"reward": {
				"denom": "ucmdx",
				"amount": "4088"
			}
		},
		{
			"address": "comdex1tvtchh7tjzhjhfldd0d7r7dg5a4jv9esjj8uh4",
			"reward": {
				"denom": "ucmdx",
				"amount": "89"
			}
		},
		{
			"address": "comdex1tvvgpullxj8zxs68zjkvxeqa4w00qe7pu8scan",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1tvvtjuejcwnfn9qhfu3f6mnn6g4ffmmedns8t8",
			"reward": {
				"denom": "ucmdx",
				"amount": "2027"
			}
		},
		{
			"address": "comdex1tvdkle3qurvd2y8a87473550kjdx7n5rquchxg",
			"reward": {
				"denom": "ucmdx",
				"amount": "140970"
			}
		},
		{
			"address": "comdex1tvdl4e9lka08ngpwhywt3x2zlm3m68t36449l8",
			"reward": {
				"denom": "ucmdx",
				"amount": "3537"
			}
		},
		{
			"address": "comdex1tvwq82a6s42c0y0yqw05w3uxx2yzzxnanc73ft",
			"reward": {
				"denom": "ucmdx",
				"amount": "2969"
			}
		},
		{
			"address": "comdex1tv06tkmualww2thjfz9nfsu344g5fytj8fwhpp",
			"reward": {
				"denom": "ucmdx",
				"amount": "156"
			}
		},
		{
			"address": "comdex1tv5wxuhf8w38jd479ayv05m2kar9wn7ekclsq3",
			"reward": {
				"denom": "ucmdx",
				"amount": "591"
			}
		},
		{
			"address": "comdex1tv5efzrzwkz0mln5pyaeex9rkta7q2xqt3rees",
			"reward": {
				"denom": "ucmdx",
				"amount": "5202"
			}
		},
		{
			"address": "comdex1tvcd8welups8wwgj9jxqpt9027eyzatljwt8e2",
			"reward": {
				"denom": "ucmdx",
				"amount": "12543"
			}
		},
		{
			"address": "comdex1tvcwgdcm7jez2wryfcrdndetpkd5fvuzgj0vkc",
			"reward": {
				"denom": "ucmdx",
				"amount": "348"
			}
		},
		{
			"address": "comdex1tveyveuc2t8arvj0lxac9mzagxl29tk3jk0znp",
			"reward": {
				"denom": "ucmdx",
				"amount": "22838"
			}
		},
		{
			"address": "comdex1tv6zkzzmutqh3z8pnptgy95tvcmxfvkwx93p98",
			"reward": {
				"denom": "ucmdx",
				"amount": "144"
			}
		},
		{
			"address": "comdex1tdr2930956y3fpdyx45de7f5x35udjhmd0k9tq",
			"reward": {
				"denom": "ucmdx",
				"amount": "114300"
			}
		},
		{
			"address": "comdex1tdrket2n9txvrdsdw3rcqg4awjptrn0respsmg",
			"reward": {
				"denom": "ucmdx",
				"amount": "387"
			}
		},
		{
			"address": "comdex1td9pka9kvzxrv7umm5yhmtsfqk0v5uqqu8ylnw",
			"reward": {
				"denom": "ucmdx",
				"amount": "143"
			}
		},
		{
			"address": "comdex1td9kf3l5vcc6d4zkgdrm90awfyu9v2dnuyyuq4",
			"reward": {
				"denom": "ucmdx",
				"amount": "3494"
			}
		},
		{
			"address": "comdex1td8vpehqf098vjfq4nxurwr02y5wfeu4vz7aqa",
			"reward": {
				"denom": "ucmdx",
				"amount": "61037"
			}
		},
		{
			"address": "comdex1tdg2jgpy83t678vla50dhf5v27n39va3lz4weu",
			"reward": {
				"denom": "ucmdx",
				"amount": "42"
			}
		},
		{
			"address": "comdex1tdv0c2v97tucdrd53lwq6uqdv7pt92hc54ymna",
			"reward": {
				"denom": "ucmdx",
				"amount": "4827"
			}
		},
		{
			"address": "comdex1td3ujdf2dl5c35y6ly4ljamhfn65wy28jnt7ly",
			"reward": {
				"denom": "ucmdx",
				"amount": "13084"
			}
		},
		{
			"address": "comdex1tdnzajzaqnzt3xdcjhyzar73txj0vj8gr8l3gn",
			"reward": {
				"denom": "ucmdx",
				"amount": "1976"
			}
		},
		{
			"address": "comdex1tduxg6mf33um8cuymam7umgyr95g55ylq22mdn",
			"reward": {
				"denom": "ucmdx",
				"amount": "140041"
			}
		},
		{
			"address": "comdex1tdak9awpptac9gtu7ymle0vhkt25wjl7e44vz3",
			"reward": {
				"denom": "ucmdx",
				"amount": "1787"
			}
		},
		{
			"address": "comdex1tdlwkxvvg39rnf4hkc3cve5zxrxd609kd03w5n",
			"reward": {
				"denom": "ucmdx",
				"amount": "1925"
			}
		},
		{
			"address": "comdex1tdlcev8329gnllv7n60fxsm55jmxeta582maep",
			"reward": {
				"denom": "ucmdx",
				"amount": "358"
			}
		},
		{
			"address": "comdex1twqplv0rexm59er6wemj6pdfa53y5rcae4ernn",
			"reward": {
				"denom": "ucmdx",
				"amount": "11214"
			}
		},
		{
			"address": "comdex1twzy9l3zy7rj2xkdc8w30mqwku0hencx3mzwfu",
			"reward": {
				"denom": "ucmdx",
				"amount": "106"
			}
		},
		{
			"address": "comdex1twyjvs8hparn39phfcg8227mme77kx5rvt2kng",
			"reward": {
				"denom": "ucmdx",
				"amount": "203"
			}
		},
		{
			"address": "comdex1tw9xan7tr52gvy2smkxha85ec7c4hunv7tkz62",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1twg49qpylps5802ma0wpyrnfeszu5xhja76f7e",
			"reward": {
				"denom": "ucmdx",
				"amount": "59878"
			}
		},
		{
			"address": "comdex1twfw8gfl2twe0mvdt50q8qa62dlk4dgdcs5me6",
			"reward": {
				"denom": "ucmdx",
				"amount": "4322"
			}
		},
		{
			"address": "comdex1tw2w3nha339gld2qddkwrn72kjw9h63m398urw",
			"reward": {
				"denom": "ucmdx",
				"amount": "29664"
			}
		},
		{
			"address": "comdex1twtlpqqlgzphukdnxedcahlql0x8ptx86rd8q2",
			"reward": {
				"denom": "ucmdx",
				"amount": "98"
			}
		},
		{
			"address": "comdex1twvx3h88e5zvmya2557cz05ncqau9venwvtrfq",
			"reward": {
				"denom": "ucmdx",
				"amount": "562705"
			}
		},
		{
			"address": "comdex1twvw5dm5864uzu84ux2nfpuden7k00x8fqxghq",
			"reward": {
				"denom": "ucmdx",
				"amount": "18550"
			}
		},
		{
			"address": "comdex1twv5acq26qnzv4xj40h60fhlq6axwl9d4n0c95",
			"reward": {
				"denom": "ucmdx",
				"amount": "12718"
			}
		},
		{
			"address": "comdex1twvezk8vvncy4937q23uj0jwyv8phy79judwfa",
			"reward": {
				"denom": "ucmdx",
				"amount": "13203"
			}
		},
		{
			"address": "comdex1tws4arhtp3yxr5vh76v9dkly8gd6tq3t3vssph",
			"reward": {
				"denom": "ucmdx",
				"amount": "6154"
			}
		},
		{
			"address": "comdex1twjpqsa8q6rypq6vdr25wjh3s4hlwu2wvpm4ga",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1tw5s9shprmaqzkkjp0v6ku8dk6af9r08tsuyf9",
			"reward": {
				"denom": "ucmdx",
				"amount": "5306"
			}
		},
		{
			"address": "comdex1twkygksd6ewacsfxfvux0ulnfjvdm5yy0cpqpy",
			"reward": {
				"denom": "ucmdx",
				"amount": "442"
			}
		},
		{
			"address": "comdex1twkn26js963hg48t63lyuaj9u63s468s89rdtf",
			"reward": {
				"denom": "ucmdx",
				"amount": "8765"
			}
		},
		{
			"address": "comdex1twcrrvqex4cy5cfdqhftkm75cvadrn6xgs4g2n",
			"reward": {
				"denom": "ucmdx",
				"amount": "1584"
			}
		},
		{
			"address": "comdex1twaqkw5fje0xj5z4jxvft6hsd2gxwfk6k4l8x4",
			"reward": {
				"denom": "ucmdx",
				"amount": "35051"
			}
		},
		{
			"address": "comdex1tw77mu0slcav5uvanjmxsyl0peyrrlpv767p26",
			"reward": {
				"denom": "ucmdx",
				"amount": "2966"
			}
		},
		{
			"address": "comdex1twlpnp9zn8hpn3c48wadffssr04e70pztrtrua",
			"reward": {
				"denom": "ucmdx",
				"amount": "3647"
			}
		},
		{
			"address": "comdex1t0q6sfhs08erlh4lpqf0xanrtzqvprer9ekdey",
			"reward": {
				"denom": "ucmdx",
				"amount": "5092"
			}
		},
		{
			"address": "comdex1t0pdyze7vfu42krk4dl8zxqg89hphwewcqjqge",
			"reward": {
				"denom": "ucmdx",
				"amount": "1777"
			}
		},
		{
			"address": "comdex1t09pev3z7uxuw4znm49nwqz8fsm4ytrtp2x9ld",
			"reward": {
				"denom": "ucmdx",
				"amount": "58"
			}
		},
		{
			"address": "comdex1t0x4sxjuzyjnekd4jlk6t8s78wk4sgg85l3z8n",
			"reward": {
				"denom": "ucmdx",
				"amount": "1991"
			}
		},
		{
			"address": "comdex1t0gcjt3f2q7etxzj5ku8ynf5pcdzw9tc76uwm2",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1t0fyv5ntem4z6nmk8utn6u0ysgc9pfxvw42utz",
			"reward": {
				"denom": "ucmdx",
				"amount": "285"
			}
		},
		{
			"address": "comdex1t0falr88272thmqrwzdf6xuktvtjzn9ldcpymc",
			"reward": {
				"denom": "ucmdx",
				"amount": "1488"
			}
		},
		{
			"address": "comdex1t0tu0vch5k6v3pcw3ua9zlme75ktzl92ds2q44",
			"reward": {
				"denom": "ucmdx",
				"amount": "10093"
			}
		},
		{
			"address": "comdex1t0tljj3ccmp3l7awn5gkwsm5exytye6ag3pxag",
			"reward": {
				"denom": "ucmdx",
				"amount": "16397"
			}
		},
		{
			"address": "comdex1t0sen2k02t2qxqyxnmcxhhrg2pe447lw0ry0mr",
			"reward": {
				"denom": "ucmdx",
				"amount": "431542"
			}
		},
		{
			"address": "comdex1t033crs7ljcgnu4fhfelgrwxvtemhacapjru5g",
			"reward": {
				"denom": "ucmdx",
				"amount": "1718"
			}
		},
		{
			"address": "comdex1t0jrqq6md64y8zswdla70jz0w86xxzxuc66tk3",
			"reward": {
				"denom": "ucmdx",
				"amount": "1984"
			}
		},
		{
			"address": "comdex1t0ndk7lsy58aedx7lkmlny7tdscqcfkzy2ve6q",
			"reward": {
				"denom": "ucmdx",
				"amount": "624"
			}
		},
		{
			"address": "comdex1t0nhs5hflr067hcxsxz2j5mkdj6td7hjg07jrz",
			"reward": {
				"denom": "ucmdx",
				"amount": "9129"
			}
		},
		{
			"address": "comdex1t05987dvn0hyn5rfa23hav2nxty6eluh68k425",
			"reward": {
				"denom": "ucmdx",
				"amount": "55137"
			}
		},
		{
			"address": "comdex1t053qpny94rt6f7hvy4crgu79rd4gsanr5urdc",
			"reward": {
				"denom": "ucmdx",
				"amount": "4546"
			}
		},
		{
			"address": "comdex1t04ksp5vn7vvh6tuaydyv4qa7rlacvhwau5q6a",
			"reward": {
				"denom": "ucmdx",
				"amount": "89"
			}
		},
		{
			"address": "comdex1t0mxq7edxf09y6we3a0qdqrwrf53kftm6x09yq",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1t0mnv26z967jy72z37x6j39xq9ntwwcdzapv95",
			"reward": {
				"denom": "ucmdx",
				"amount": "13875"
			}
		},
		{
			"address": "comdex1t0u9hawxet4czhfcm3v5k7kqjpcvp8sthwt4ht",
			"reward": {
				"denom": "ucmdx",
				"amount": "1733"
			}
		},
		{
			"address": "comdex1t0ascndjtgznu9mpeaeznznmxv6tt42njkdavy",
			"reward": {
				"denom": "ucmdx",
				"amount": "667"
			}
		},
		{
			"address": "comdex1t0lwp8h3ardsn550fwgpq4mkq5u8nmlz565jzu",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1tsq6ceep8s8cghshpj2mzwungxww4sqr57dgrw",
			"reward": {
				"denom": "ucmdx",
				"amount": "1470"
			}
		},
		{
			"address": "comdex1tspftx9q6zxtepgt2hdw90nd4pnggda9fe2l55",
			"reward": {
				"denom": "ucmdx",
				"amount": "8369"
			}
		},
		{
			"address": "comdex1tsy9gndt4paq4hfvghjaftvsuc55vf6s3speg0",
			"reward": {
				"denom": "ucmdx",
				"amount": "4811"
			}
		},
		{
			"address": "comdex1ts9j9a2acqjd6wu0u3pm6krea4dwdvdv6hdsmj",
			"reward": {
				"denom": "ucmdx",
				"amount": "534"
			}
		},
		{
			"address": "comdex1tsx6mjzmnyxmrm4q6yl2kx6kscuuhpc9q6apac",
			"reward": {
				"denom": "ucmdx",
				"amount": "5292"
			}
		},
		{
			"address": "comdex1tsv2qhh263naaklxlwk80e9c4mcprwrvq46gne",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1tswycvy9e5rvltmwg9z2fe8zmfueh7djpyex4a",
			"reward": {
				"denom": "ucmdx",
				"amount": "5569"
			}
		},
		{
			"address": "comdex1tss72ylng5az3rjq6qzu05ze4rgc9d4m4m0gr6",
			"reward": {
				"denom": "ucmdx",
				"amount": "166"
			}
		},
		{
			"address": "comdex1ts3pncya02duykrjj3tmmkdy6nwvmwzqgmdlj3",
			"reward": {
				"denom": "ucmdx",
				"amount": "1432"
			}
		},
		{
			"address": "comdex1tsjc4r6u38snatq9gxsxchq3q96vcclgpemsg8",
			"reward": {
				"denom": "ucmdx",
				"amount": "90795"
			}
		},
		{
			"address": "comdex1tshwngxafaaetznslcppa8exqc0zhfr4rj2k8e",
			"reward": {
				"denom": "ucmdx",
				"amount": "57673"
			}
		},
		{
			"address": "comdex1tshuryfz6lnkpkpc8dy5sytcrnhjrpffqy4efm",
			"reward": {
				"denom": "ucmdx",
				"amount": "170"
			}
		},
		{
			"address": "comdex1tsc2g4fwvjxjeer4cj3qcj24xj6t57qyajutzm",
			"reward": {
				"denom": "ucmdx",
				"amount": "882"
			}
		},
		{
			"address": "comdex1tsehe2v6rxqzw7875p2770f50n44n8nva6yu9q",
			"reward": {
				"denom": "ucmdx",
				"amount": "183817"
			}
		},
		{
			"address": "comdex1tsutq6zqxm6jp38mj2vle3t9ea9zmqj00mtekw",
			"reward": {
				"denom": "ucmdx",
				"amount": "302"
			}
		},
		{
			"address": "comdex1ts70x8nwej3tnd724fkxc73dtjeaqjks936g9n",
			"reward": {
				"denom": "ucmdx",
				"amount": "1554"
			}
		},
		{
			"address": "comdex1t3r07cp8yqrj69ddcxz0rct94kq56nz9qq3gtx",
			"reward": {
				"denom": "ucmdx",
				"amount": "8392"
			}
		},
		{
			"address": "comdex1t3r6xel6yrfuhty0h3w3004ppsxk4calsvrxd5",
			"reward": {
				"denom": "ucmdx",
				"amount": "62291"
			}
		},
		{
			"address": "comdex1t3g2c7yntsj05yc967hzwhugtcv7l5998js6ky",
			"reward": {
				"denom": "ucmdx",
				"amount": "791857"
			}
		},
		{
			"address": "comdex1t32lv3ec0c79k8e4dz33p8hnq6ph5zypp46xlw",
			"reward": {
				"denom": "ucmdx",
				"amount": "88"
			}
		},
		{
			"address": "comdex1t3v309m87vdc5zc9d365yjwk4y76rq9lwthd4j",
			"reward": {
				"denom": "ucmdx",
				"amount": "61648"
			}
		},
		{
			"address": "comdex1t3whg5h5cwc6rtdygf5q4r8pyhe2grkuhx99uw",
			"reward": {
				"denom": "ucmdx",
				"amount": "145"
			}
		},
		{
			"address": "comdex1t308jd83sy8zcv5jzrh7cd852sjzp9um03h9dz",
			"reward": {
				"denom": "ucmdx",
				"amount": "3587"
			}
		},
		{
			"address": "comdex1t3387esh7x7hhz0d5908e4vgf7f7lyvhkau0cm",
			"reward": {
				"denom": "ucmdx",
				"amount": "1251"
			}
		},
		{
			"address": "comdex1t3jrhn2h3r8t55dn35p095ssfk4czywynjmjrt",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1t3nw4ppfset7u0f0ueau3fkwmc0dd9xmm503f4",
			"reward": {
				"denom": "ucmdx",
				"amount": "502920"
			}
		},
		{
			"address": "comdex1t3n7zrq52snn3kw23my80q7ln3t8l3v47gnlqg",
			"reward": {
				"denom": "ucmdx",
				"amount": "137608"
			}
		},
		{
			"address": "comdex1t3hcmuv42krnuxf0ejwsh75dx6h2mpzyhe29wn",
			"reward": {
				"denom": "ucmdx",
				"amount": "1464"
			}
		},
		{
			"address": "comdex1t3ecu6khpws8jz4wuu770yqahvwnfrydr4sdmw",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1t36cnacr83tdljnnv9qv3skndfdslwfd5ntv7j",
			"reward": {
				"denom": "ucmdx",
				"amount": "136182"
			}
		},
		{
			"address": "comdex1t3myppltl0jznzttf00lva982t8vl7eh4pkyhy",
			"reward": {
				"denom": "ucmdx",
				"amount": "17844"
			}
		},
		{
			"address": "comdex1t3usy5x8xfggzspknnka02ny7u65u6k0kv0ed4",
			"reward": {
				"denom": "ucmdx",
				"amount": "7046"
			}
		},
		{
			"address": "comdex1t3avpnkcscu55dkjng4tdhc8lc2hq3mlhm8cu7",
			"reward": {
				"denom": "ucmdx",
				"amount": "124"
			}
		},
		{
			"address": "comdex1tjql39rx73f4gteu397c0nvg0nsez767kaet4h",
			"reward": {
				"denom": "ucmdx",
				"amount": "1243"
			}
		},
		{
			"address": "comdex1tjp774e79uq7mnfjwuzc4astv2zqc3xxg7t5lh",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1tjr8sc8kduckucu8h094muqcd9gurcewr4rk77",
			"reward": {
				"denom": "ucmdx",
				"amount": "3322"
			}
		},
		{
			"address": "comdex1tjy722hfsmcgke83knd4dlp6jvyycet77lr73q",
			"reward": {
				"denom": "ucmdx",
				"amount": "12404"
			}
		},
		{
			"address": "comdex1tjxkwc0at9v2n8fm5lmp5mzsr67wtlaccslj4u",
			"reward": {
				"denom": "ucmdx",
				"amount": "36098"
			}
		},
		{
			"address": "comdex1tjx7m3lltqkdrhdhy9v5h90lwerl3tkryn2nn4",
			"reward": {
				"denom": "ucmdx",
				"amount": "6013"
			}
		},
		{
			"address": "comdex1tjgphyrve5et5jpw3fata3ysydzntkqcegw7za",
			"reward": {
				"denom": "ucmdx",
				"amount": "3513"
			}
		},
		{
			"address": "comdex1tjg0kxcxdn4t589dzxml255c57fsr9w7mveqg5",
			"reward": {
				"denom": "ucmdx",
				"amount": "212"
			}
		},
		{
			"address": "comdex1tjvczk4nryajvh8v76xeua5jyf0h5dn4lt0gpu",
			"reward": {
				"denom": "ucmdx",
				"amount": "3831"
			}
		},
		{
			"address": "comdex1tjszaerpqh7e2t5myagfxkc6605mgxpa4cm5qc",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1tj3j9dfqvfljtytzvrcs9fdddjy8lkmzmxtkl0",
			"reward": {
				"denom": "ucmdx",
				"amount": "22824"
			}
		},
		{
			"address": "comdex1tj3jhe269ltdu3mnmuyd37r3yakct9s5m306nf",
			"reward": {
				"denom": "ucmdx",
				"amount": "2675"
			}
		},
		{
			"address": "comdex1tjjz2nv9jqn0802dfp0z2lm9walmg7e49t6rg6",
			"reward": {
				"denom": "ucmdx",
				"amount": "740"
			}
		},
		{
			"address": "comdex1tjkl46wwpxtym4f83wwflnyptlcvlvjfqeat9r",
			"reward": {
				"denom": "ucmdx",
				"amount": "1019"
			}
		},
		{
			"address": "comdex1tjhfuchkfdz973at5q7xdvux2mcfna6zeegjla",
			"reward": {
				"denom": "ucmdx",
				"amount": "12700"
			}
		},
		{
			"address": "comdex1tjc6y59tsn9dz2yaypneld4z3nxuqfx2ns7y0z",
			"reward": {
				"denom": "ucmdx",
				"amount": "14"
			}
		},
		{
			"address": "comdex1tjmlsvzc3ejstjwm2889c44ejrsat2ph899cy8",
			"reward": {
				"denom": "ucmdx",
				"amount": "15154"
			}
		},
		{
			"address": "comdex1tjadvuzp4g6r6nq8mqe2gsk23pghaf6ye5vas5",
			"reward": {
				"denom": "ucmdx",
				"amount": "528"
			}
		},
		{
			"address": "comdex1tjajx9c75kguhdl65masvyfrm8z2d48tv5hgge",
			"reward": {
				"denom": "ucmdx",
				"amount": "1723"
			}
		},
		{
			"address": "comdex1tj77h9xdwgvfs77p669eqehlhu8a29zr7678q3",
			"reward": {
				"denom": "ucmdx",
				"amount": "17530"
			}
		},
		{
			"address": "comdex1tnrx7z2s8ztj3m8kfm0pleydfyt69wsshw4sps",
			"reward": {
				"denom": "ucmdx",
				"amount": "6186"
			}
		},
		{
			"address": "comdex1tnr7wr7uku6jm48ad6x3xjef4c3l2fk4zqq5w0",
			"reward": {
				"denom": "ucmdx",
				"amount": "1432"
			}
		},
		{
			"address": "comdex1tnyrxmdk6jtw8vqsq74z8t8y562mgcf7k6hg7z",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1tn9v83gg4lex96ujgmpv8zd3v6xn29wexh4rsw",
			"reward": {
				"denom": "ucmdx",
				"amount": "1804"
			}
		},
		{
			"address": "comdex1tn9hg4j8xj7jjcup33cmr96egtkz8g5hc7p6te",
			"reward": {
				"denom": "ucmdx",
				"amount": "153347"
			}
		},
		{
			"address": "comdex1tnx3z8ke9za65v32qxhu65mav2282jrulgxw5y",
			"reward": {
				"denom": "ucmdx",
				"amount": "14480"
			}
		},
		{
			"address": "comdex1tngrk6fdalfzff4k90uzyy88zfprg9vnu835wc",
			"reward": {
				"denom": "ucmdx",
				"amount": "374226"
			}
		},
		{
			"address": "comdex1tngsl3htr8ltdtcqmjv05207tsvesezke59k7p",
			"reward": {
				"denom": "ucmdx",
				"amount": "27200"
			}
		},
		{
			"address": "comdex1tngl57943azs5pkyawfydc4fhptzxkc8a4yq9s",
			"reward": {
				"denom": "ucmdx",
				"amount": "4949"
			}
		},
		{
			"address": "comdex1tnffg6s9x35eac2kaw4zvkx7lvyhd7spkhkcy8",
			"reward": {
				"denom": "ucmdx",
				"amount": "747"
			}
		},
		{
			"address": "comdex1tntlr6g3f5mzkguxl58v7y9z403x9fzlmz07mm",
			"reward": {
				"denom": "ucmdx",
				"amount": "1485"
			}
		},
		{
			"address": "comdex1tnd24rd9c8kv4xdfw7k9yyqcxntdp05czmds2n",
			"reward": {
				"denom": "ucmdx",
				"amount": "2398"
			}
		},
		{
			"address": "comdex1tnwt696l9jspwlehcafkcdpxzvcpzcw7rcsvh0",
			"reward": {
				"denom": "ucmdx",
				"amount": "141"
			}
		},
		{
			"address": "comdex1tn0xzu0msj6nhs5n43uet87ljzsrxgy59wtyks",
			"reward": {
				"denom": "ucmdx",
				"amount": "23925"
			}
		},
		{
			"address": "comdex1tn08fv8mj5kptv894fyp3j9np72vhpk02n5xem",
			"reward": {
				"denom": "ucmdx",
				"amount": "63299"
			}
		},
		{
			"address": "comdex1tn0np0pt2wvq6e7wqun6amumptcypk5d74tvhj",
			"reward": {
				"denom": "ucmdx",
				"amount": "1075"
			}
		},
		{
			"address": "comdex1tnnhrqpm8rxj7xja5jcu2vmlv2yxz5nadavwn7",
			"reward": {
				"denom": "ucmdx",
				"amount": "8918"
			}
		},
		{
			"address": "comdex1tn5xtk6l0p8pqk2e7k97fnklv5jxdp0fjajtyp",
			"reward": {
				"denom": "ucmdx",
				"amount": "7227"
			}
		},
		{
			"address": "comdex1tnkvcdf645dhg49l32lkfsu3djlnqf4k524r4j",
			"reward": {
				"denom": "ucmdx",
				"amount": "12377"
			}
		},
		{
			"address": "comdex1tnkme3hckzaychu8echwhm93sndgp379gfssyg",
			"reward": {
				"denom": "ucmdx",
				"amount": "16"
			}
		},
		{
			"address": "comdex1tnls6mqv6efye6lw4u9qz6rhkajwxj08ettjct",
			"reward": {
				"denom": "ucmdx",
				"amount": "702"
			}
		},
		{
			"address": "comdex1t5r4lhmd3j0uzzsvdtrraskckwetdkktp935u8",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1t5r7z6ngykst3hutqjjzmnulvk34t87annp2r8",
			"reward": {
				"denom": "ucmdx",
				"amount": "14276"
			}
		},
		{
			"address": "comdex1t595842x24q2s80gqpatcu0hx4w8qpnmq9fpmw",
			"reward": {
				"denom": "ucmdx",
				"amount": "1775"
			}
		},
		{
			"address": "comdex1t5gtn03s6r8x6hyuw57j79a0jzdg4s8zt6tudc",
			"reward": {
				"denom": "ucmdx",
				"amount": "1731"
			}
		},
		{
			"address": "comdex1t5g4606lp5p4n52cx5hxyftq0hq2n3pzuvhfwy",
			"reward": {
				"denom": "ucmdx",
				"amount": "2835"
			}
		},
		{
			"address": "comdex1t5fjkzx9qnwytaghgd0mxv6cmvx03fvydssgjh",
			"reward": {
				"denom": "ucmdx",
				"amount": "2403"
			}
		},
		{
			"address": "comdex1t52vm8yf0u96tckvhhfq54z44cdh956q0mvsxd",
			"reward": {
				"denom": "ucmdx",
				"amount": "12632"
			}
		},
		{
			"address": "comdex1t5wmu8x093ac6xsvqdm2v95a3ekk7fsr4kw6xd",
			"reward": {
				"denom": "ucmdx",
				"amount": "721"
			}
		},
		{
			"address": "comdex1t509x4544hmnawaj4a0uwjy603efe50ee5gklr",
			"reward": {
				"denom": "ucmdx",
				"amount": "11221"
			}
		},
		{
			"address": "comdex1t55g6jh4fgs9zhkjp39dkh4wycttts7yehsev7",
			"reward": {
				"denom": "ucmdx",
				"amount": "40052"
			}
		},
		{
			"address": "comdex1t552ujd5nnc52nye56c0sdlrwc3vvzthumph4f",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1t5h20xuz4v483p32l0x7xjetxhu6dee35j4mpu",
			"reward": {
				"denom": "ucmdx",
				"amount": "434"
			}
		},
		{
			"address": "comdex1t5cmxpum635jjwunc68e9lvghgek6qv864ml5t",
			"reward": {
				"denom": "ucmdx",
				"amount": "14"
			}
		},
		{
			"address": "comdex1t5ajej6t2mzyrjhvp6yqysx0ycf069d2dxk5f7",
			"reward": {
				"denom": "ucmdx",
				"amount": "64"
			}
		},
		{
			"address": "comdex1t5lx9y368e0fsf8yr84hqwdy3rqtut3p98y40d",
			"reward": {
				"denom": "ucmdx",
				"amount": "10930"
			}
		},
		{
			"address": "comdex1t5lgz6tenwd6x3rw7prlz86qeuurufttu55ghq",
			"reward": {
				"denom": "ucmdx",
				"amount": "7641"
			}
		},
		{
			"address": "comdex1t4q2as2p5m39tnsnnmfa4zxndmm2vhrrlzzg64",
			"reward": {
				"denom": "ucmdx",
				"amount": "151"
			}
		},
		{
			"address": "comdex1t4prj5rpgm0h6nppzzyjk7ygkksqxmuwjk08sf",
			"reward": {
				"denom": "ucmdx",
				"amount": "142230"
			}
		},
		{
			"address": "comdex1t4z3kvzsjh7myaxn57cyhm4v8kjhuj2etunqpm",
			"reward": {
				"denom": "ucmdx",
				"amount": "2859"
			}
		},
		{
			"address": "comdex1t4zjnf2x0zczm39a3lks0ekerpd6we22r7p430",
			"reward": {
				"denom": "ucmdx",
				"amount": "40"
			}
		},
		{
			"address": "comdex1t49cnx72fy86kauqk07pht23klq97ufcksyr23",
			"reward": {
				"denom": "ucmdx",
				"amount": "1759"
			}
		},
		{
			"address": "comdex1t4g0tg8s8nysfcsjcg4v8nx6mk3ca322ltxht8",
			"reward": {
				"denom": "ucmdx",
				"amount": "14587"
			}
		},
		{
			"address": "comdex1t42xgztw2ety5gleldmy3xxm4d7yz2e8vehptc",
			"reward": {
				"denom": "ucmdx",
				"amount": "180"
			}
		},
		{
			"address": "comdex1t4v3h355vsvexz46nz6h9repmgxgfml897l7an",
			"reward": {
				"denom": "ucmdx",
				"amount": "13135"
			}
		},
		{
			"address": "comdex1t4drg6nh49pusk2aas8tj7vlwmv0cdgmaasyjl",
			"reward": {
				"denom": "ucmdx",
				"amount": "12477"
			}
		},
		{
			"address": "comdex1t4dy2r53v3vuf22zlueeheg6t3mdgstgcevnec",
			"reward": {
				"denom": "ucmdx",
				"amount": "7252"
			}
		},
		{
			"address": "comdex1t43ysmgx80gae9g6zdu7jwf9qe20k5prsfpgt7",
			"reward": {
				"denom": "ucmdx",
				"amount": "726"
			}
		},
		{
			"address": "comdex1t4js2s63xd8sukqq9ckw5zxgykd6j9m049d34e",
			"reward": {
				"denom": "ucmdx",
				"amount": "553"
			}
		},
		{
			"address": "comdex1t4kz5yk8x8x0u35mdmzxssj6qmpqmxtf2q6llv",
			"reward": {
				"denom": "ucmdx",
				"amount": "324"
			}
		},
		{
			"address": "comdex1t4c7l2v2s8p9d0ukctkhduqj0n85rwaegt8ef8",
			"reward": {
				"denom": "ucmdx",
				"amount": "3495"
			}
		},
		{
			"address": "comdex1t4etz7yktnexalgmt6jr0g93rdfltvfm5ws3sq",
			"reward": {
				"denom": "ucmdx",
				"amount": "2448359"
			}
		},
		{
			"address": "comdex1t4etdtjz63xnlt8n8xlmnfvym4gjrt07q8j63k",
			"reward": {
				"denom": "ucmdx",
				"amount": "2283"
			}
		},
		{
			"address": "comdex1t47fyg46afchqykm8stktfhayunp0axrn62vwl",
			"reward": {
				"denom": "ucmdx",
				"amount": "3670"
			}
		},
		{
			"address": "comdex1t47473lel2mq3f8k0lcl8rkea4rmfvmg473gv8",
			"reward": {
				"denom": "ucmdx",
				"amount": "6097"
			}
		},
		{
			"address": "comdex1t4l2wgy8rl8neshxv5nww4ls3s0qpq5zmqfma8",
			"reward": {
				"denom": "ucmdx",
				"amount": "2507"
			}
		},
		{
			"address": "comdex1t4l23wsja0qcx7cq80lhwzltyasfsfnadszqj9",
			"reward": {
				"denom": "ucmdx",
				"amount": "1529"
			}
		},
		{
			"address": "comdex1tkp5h9s58mgh4d6ecsznrgynfx4cn6l7hpn8l7",
			"reward": {
				"denom": "ucmdx",
				"amount": "26424"
			}
		},
		{
			"address": "comdex1tkzk79faqzthe8yd4mtfw696mr7xfw6z8n8gpn",
			"reward": {
				"denom": "ucmdx",
				"amount": "339"
			}
		},
		{
			"address": "comdex1tkzu4s42aksksclw7eygsqj4rhsvdpenmn8dsn",
			"reward": {
				"denom": "ucmdx",
				"amount": "10215"
			}
		},
		{
			"address": "comdex1tkytcv03k89e9dafnt60mu0s5s3ugjeqzus79q",
			"reward": {
				"denom": "ucmdx",
				"amount": "524"
			}
		},
		{
			"address": "comdex1tkymudh9lhnckh6wd00scu4cya5eyaj9ljkffz",
			"reward": {
				"denom": "ucmdx",
				"amount": "17609"
			}
		},
		{
			"address": "comdex1tkvp6yhem77f5cz39yyxhw422tm4qcu57ujtmt",
			"reward": {
				"denom": "ucmdx",
				"amount": "40811"
			}
		},
		{
			"address": "comdex1tkvg4rhxgk35uuzcwdv8rrvv2ww34y3p0wqt92",
			"reward": {
				"denom": "ucmdx",
				"amount": "680481"
			}
		},
		{
			"address": "comdex1tkvts5ppmjtrschvpk5v0qv904xlda2yvpk0m7",
			"reward": {
				"denom": "ucmdx",
				"amount": "141550"
			}
		},
		{
			"address": "comdex1tkvkxz9fqrcd45zp3pxhtz7fhzmxgm5vg6xr77",
			"reward": {
				"denom": "ucmdx",
				"amount": "358"
			}
		},
		{
			"address": "comdex1tks9cvg9l4qgsza0pgz5gxdc7smh7rljejr5nv",
			"reward": {
				"denom": "ucmdx",
				"amount": "133"
			}
		},
		{
			"address": "comdex1tk3a2f9ty0e02sdatdj6y3xwqlvxk7zhkkqwyj",
			"reward": {
				"denom": "ucmdx",
				"amount": "20938"
			}
		},
		{
			"address": "comdex1tkj0w6jay6uh5j8pytpdwn94tjkn7ck2el8n3g",
			"reward": {
				"denom": "ucmdx",
				"amount": "1190"
			}
		},
		{
			"address": "comdex1tk5lzp7g7j0xtmrsvxvy9gv28whrgkh6hmx3mt",
			"reward": {
				"denom": "ucmdx",
				"amount": "78878"
			}
		},
		{
			"address": "comdex1tkk0qg6p5r9m6ap8gn6sshlj5wxgqkmvmc0v98",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1tkkesjkmck63vaw0h46ttsc7ld9ttrhxray8n7",
			"reward": {
				"denom": "ucmdx",
				"amount": "1799"
			}
		},
		{
			"address": "comdex1tk7enahu86762uvvpvxyqg8nlnpv9qf0nm67jt",
			"reward": {
				"denom": "ucmdx",
				"amount": "268"
			}
		},
		{
			"address": "comdex1thpf0m46pa3305myxh5z3yp5kvr8ys6e5e688k",
			"reward": {
				"denom": "ucmdx",
				"amount": "585"
			}
		},
		{
			"address": "comdex1thrv4n06k5zj30y82nvxl6f9fc89cz06d4mrzx",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1thy7fckmhp2ytwjeygdyhvms578pm2u39zkg74",
			"reward": {
				"denom": "ucmdx",
				"amount": "5354"
			}
		},
		{
			"address": "comdex1th9c8yd3m9d3hf08ft6khg9fumfcdnhf2g7g6y",
			"reward": {
				"denom": "ucmdx",
				"amount": "891"
			}
		},
		{
			"address": "comdex1tht2l85qc5qc2n46aa902hsuxhqf2x7sutwxxc",
			"reward": {
				"denom": "ucmdx",
				"amount": "12718"
			}
		},
		{
			"address": "comdex1tht0mf64wy6lq5axggh6j73uzga4mw5dxntheq",
			"reward": {
				"denom": "ucmdx",
				"amount": "375800"
			}
		},
		{
			"address": "comdex1thvruh96q8j7pr0uxphgfztphdmak02sjudmy8",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1thvgkh8zy46fh7mc357dnm8n9j0v8yanz2akk8",
			"reward": {
				"denom": "ucmdx",
				"amount": "12860"
			}
		},
		{
			"address": "comdex1thvvuw62ztf92yk98a5qdylesaun9luj5zvsjm",
			"reward": {
				"denom": "ucmdx",
				"amount": "4908"
			}
		},
		{
			"address": "comdex1thvncfatxukryyceuqxh5fzvxa57sj7qf8yseh",
			"reward": {
				"denom": "ucmdx",
				"amount": "31112"
			}
		},
		{
			"address": "comdex1th0jfqn8k0vef72yej2yjv303ej0cwc7uj4vfy",
			"reward": {
				"denom": "ucmdx",
				"amount": "1024"
			}
		},
		{
			"address": "comdex1thnx8asy4h0fmjh38vykcc7fsggqqwf2dvk7e4",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1thnjt5wf6v4j3wux7kflgd2plrnf6nu26yd4d4",
			"reward": {
				"denom": "ucmdx",
				"amount": "3242"
			}
		},
		{
			"address": "comdex1thkkpje9aq2tz42vgr8aj0vdtz7mxsc24uvfq8",
			"reward": {
				"denom": "ucmdx",
				"amount": "7985"
			}
		},
		{
			"address": "comdex1thhxjdqy2gujjhukugn7zh6rjh54pfma7hdjg7",
			"reward": {
				"denom": "ucmdx",
				"amount": "7184"
			}
		},
		{
			"address": "comdex1thhfnltq8yjpwmvujzqp8thcypq7ylcx0hehtd",
			"reward": {
				"denom": "ucmdx",
				"amount": "7030"
			}
		},
		{
			"address": "comdex1thcxkkjx0f97q5mne2w7kfwddlw5mg7hft57zh",
			"reward": {
				"denom": "ucmdx",
				"amount": "1741"
			}
		},
		{
			"address": "comdex1thm5qu3dl7x5tzdwzfjelj8m0jr0gsx0ytkf25",
			"reward": {
				"denom": "ucmdx",
				"amount": "4929"
			}
		},
		{
			"address": "comdex1tharcgrfu6j0dcwpe5y6ez3s904rhq2k22tjay",
			"reward": {
				"denom": "ucmdx",
				"amount": "174"
			}
		},
		{
			"address": "comdex1thadzqdd7frcqn3d22nut2f4ylsvaee4n2at3p",
			"reward": {
				"denom": "ucmdx",
				"amount": "122"
			}
		},
		{
			"address": "comdex1tc8d22lxwutsq9ph2rhe37k2ujfllu6z4sa7ua",
			"reward": {
				"denom": "ucmdx",
				"amount": "5995"
			}
		},
		{
			"address": "comdex1tc8mfh0wwpsvdkmykl4nvyeravv6wt2pt0lqha",
			"reward": {
				"denom": "ucmdx",
				"amount": "1482"
			}
		},
		{
			"address": "comdex1tc2hll3wfhu247scwnvu6ppgrwgc6sjzqkk4hy",
			"reward": {
				"denom": "ucmdx",
				"amount": "21617"
			}
		},
		{
			"address": "comdex1tctusz3svjsce6tu82gwt35ang3z74ywyr5w4y",
			"reward": {
				"denom": "ucmdx",
				"amount": "743"
			}
		},
		{
			"address": "comdex1tcd2nhh5pkqhhl3s0schp2enatls4n9cyla2kh",
			"reward": {
				"denom": "ucmdx",
				"amount": "31"
			}
		},
		{
			"address": "comdex1tcwh9xqp973v32ts58za3pp3hur89zsmqpm4eu",
			"reward": {
				"denom": "ucmdx",
				"amount": "2054"
			}
		},
		{
			"address": "comdex1tcjzju5sxh36a9wxcf73crl2z7t6z8lf3rte2r",
			"reward": {
				"denom": "ucmdx",
				"amount": "1409"
			}
		},
		{
			"address": "comdex1tccsy053nry9g9aaf73gucg6vrh8y0wleyq444",
			"reward": {
				"denom": "ucmdx",
				"amount": "16612"
			}
		},
		{
			"address": "comdex1tcc6834rlrpavrxf2a3t967es70trjary3xf0k",
			"reward": {
				"denom": "ucmdx",
				"amount": "2006"
			}
		},
		{
			"address": "comdex1tc6r2ecnu3gf9fsdzaj4kp8zzukse6090nwnpc",
			"reward": {
				"denom": "ucmdx",
				"amount": "1740"
			}
		},
		{
			"address": "comdex1tcmxtpxl2strglzmdkxze60hf53hrev4cl7pcl",
			"reward": {
				"denom": "ucmdx",
				"amount": "27470"
			}
		},
		{
			"address": "comdex1tcuqa8jwar5vcrlxet2du4pj9yzyglk9gn98l2",
			"reward": {
				"denom": "ucmdx",
				"amount": "3351"
			}
		},
		{
			"address": "comdex1tclwugzd0szrafm4rrt8wv3734m3a79rf6k4z7",
			"reward": {
				"denom": "ucmdx",
				"amount": "104786"
			}
		},
		{
			"address": "comdex1tcleqqtpe9jnl9y5qkah9j40e99k5c5mykhq7a",
			"reward": {
				"denom": "ucmdx",
				"amount": "19426"
			}
		},
		{
			"address": "comdex1tezzmqjvh0ql6gls26qh27jfrnn4e3wj2g4c4a",
			"reward": {
				"denom": "ucmdx",
				"amount": "167837"
			}
		},
		{
			"address": "comdex1tervfp4zlfmxky9reur09nlyyka9ezurfl4v6c",
			"reward": {
				"denom": "ucmdx",
				"amount": "35620"
			}
		},
		{
			"address": "comdex1te948j8q4w6slmqsjrnth6376jvhmf28vh6c5s",
			"reward": {
				"denom": "ucmdx",
				"amount": "105205"
			}
		},
		{
			"address": "comdex1tettdt8ftc449pza0hw2gdrnl6pksdjn7l0xg0",
			"reward": {
				"denom": "ucmdx",
				"amount": "3465"
			}
		},
		{
			"address": "comdex1tevtg4ygkr2e0t4tjkvlfk45qq7qpmg8zmlsd7",
			"reward": {
				"denom": "ucmdx",
				"amount": "16751"
			}
		},
		{
			"address": "comdex1tewgmlcgks6z8hullc35tx80nccal9mvq8wlua",
			"reward": {
				"denom": "ucmdx",
				"amount": "818"
			}
		},
		{
			"address": "comdex1tespngs8q3m4a799hxdw37zqde6ge9ne62k9j5",
			"reward": {
				"denom": "ucmdx",
				"amount": "28"
			}
		},
		{
			"address": "comdex1tentvycdxl5v7f6xtlemdywf4j8mwq4evra6c9",
			"reward": {
				"denom": "ucmdx",
				"amount": "4223"
			}
		},
		{
			"address": "comdex1te5lt7thak2ktrqdqfjkvl7ud66e4jaunkx6ul",
			"reward": {
				"denom": "ucmdx",
				"amount": "713"
			}
		},
		{
			"address": "comdex1tekfrly8jyspek7lyattgq0ayqswjzml8gs7dy",
			"reward": {
				"denom": "ucmdx",
				"amount": "840"
			}
		},
		{
			"address": "comdex1tekcy4723dnu0j3zurczxdj9yek4ujd3p2yxug",
			"reward": {
				"denom": "ucmdx",
				"amount": "58937"
			}
		},
		{
			"address": "comdex1tee99rqet3tza2tpygfaqy70dk4syz8ldnhgh0",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1temvcezkp2f5u42u2f69ra8k85c5p8n30yxvz8",
			"reward": {
				"denom": "ucmdx",
				"amount": "167616"
			}
		},
		{
			"address": "comdex1tel4alqtz8y3pkgvezqaj2hpusd8sav5n0enem",
			"reward": {
				"denom": "ucmdx",
				"amount": "8758"
			}
		},
		{
			"address": "comdex1t6qmcazp9kgdx6a7qurvtqhzxscd6sf3kwgn7t",
			"reward": {
				"denom": "ucmdx",
				"amount": "24419"
			}
		},
		{
			"address": "comdex1t6pg6vuj9jp2fw0en7fk7awqxsyzlgq7xlvsrw",
			"reward": {
				"denom": "ucmdx",
				"amount": "25"
			}
		},
		{
			"address": "comdex1t6rr3r5gp94cvgmmspq7gx5tjl7tjhq538g3fx",
			"reward": {
				"denom": "ucmdx",
				"amount": "43947"
			}
		},
		{
			"address": "comdex1t6r78fr392c0ey8stxlffhag0ejxgq7g7gegkk",
			"reward": {
				"denom": "ucmdx",
				"amount": "288"
			}
		},
		{
			"address": "comdex1t6y876a4r25mfercqu2vg9xx09afs3qt8tqrpx",
			"reward": {
				"denom": "ucmdx",
				"amount": "182"
			}
		},
		{
			"address": "comdex1t6x3lajrqnnqzjfrkvp94nf0lx3jlak5ej7ax2",
			"reward": {
				"denom": "ucmdx",
				"amount": "9643"
			}
		},
		{
			"address": "comdex1t6f896n5cd5d605ye6d9fzuvgj5r836vn9d8dr",
			"reward": {
				"denom": "ucmdx",
				"amount": "5289"
			}
		},
		{
			"address": "comdex1t6deppgktctsd5q6cn3txh7uwwatu6j06e9nwh",
			"reward": {
				"denom": "ucmdx",
				"amount": "797"
			}
		},
		{
			"address": "comdex1t64u2enf3t7jrq2sw49n2aqngxhpa4qjp6ssxh",
			"reward": {
				"denom": "ucmdx",
				"amount": "125004"
			}
		},
		{
			"address": "comdex1t6exp5dc29ljwy93vlp3rmlgs0dvxhasnu9d6v",
			"reward": {
				"denom": "ucmdx",
				"amount": "27012"
			}
		},
		{
			"address": "comdex1t66d4thfctrksdkdve4y7r3h06e5llwgc6gjux",
			"reward": {
				"denom": "ucmdx",
				"amount": "126080"
			}
		},
		{
			"address": "comdex1t6ae72nm96rs4tw30q6pclqlxkj0t972d0p6tw",
			"reward": {
				"denom": "ucmdx",
				"amount": "88"
			}
		},
		{
			"address": "comdex1tmp99hgqux43a3hnr29k3kaks4stjd583e9g3c",
			"reward": {
				"denom": "ucmdx",
				"amount": "129"
			}
		},
		{
			"address": "comdex1tmz5r4zy8ef0qqgy737uwsmwm5sv79mgu2nw32",
			"reward": {
				"denom": "ucmdx",
				"amount": "1568"
			}
		},
		{
			"address": "comdex1tmy8wkswuznf80xqcxq8xmd7djqqyqn0lxnf0c",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1tmynkkhlrpu84592n42r75xsyw3d70xhu797zd",
			"reward": {
				"denom": "ucmdx",
				"amount": "21"
			}
		},
		{
			"address": "comdex1tmxpvmqcdv3rwnfurq922xlq873yeteuxq5k9u",
			"reward": {
				"denom": "ucmdx",
				"amount": "239770"
			}
		},
		{
			"address": "comdex1tmtwh7lddse4qt6lpk4d8ze3av4wkw7nqtmhrf",
			"reward": {
				"denom": "ucmdx",
				"amount": "12326"
			}
		},
		{
			"address": "comdex1tmvhxvmcu3krv9rsegl9q8rfehu7reem22ujjj",
			"reward": {
				"denom": "ucmdx",
				"amount": "1559"
			}
		},
		{
			"address": "comdex1tm5lt2y6sw62dptctecq9g9wjv7hwzkf5dt57p",
			"reward": {
				"denom": "ucmdx",
				"amount": "36"
			}
		},
		{
			"address": "comdex1tmc04s7zrlrul96d7glmv65x8p837fmyatm0fe",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1tmm9nfl3ycv70fzxglzgm3vt0yfjmes8jaxdsu",
			"reward": {
				"denom": "ucmdx",
				"amount": "4999179"
			}
		},
		{
			"address": "comdex1tmmdj7svuaek2umnuxtan5rpm4s0nvkxrxyfny",
			"reward": {
				"denom": "ucmdx",
				"amount": "192"
			}
		},
		{
			"address": "comdex1tmarkr6ta2sge625udkl6rfgqnarn4h5d4zmqy",
			"reward": {
				"denom": "ucmdx",
				"amount": "530"
			}
		},
		{
			"address": "comdex1tm7va3plpzw49k6g5myyfe3jqfxqlgun9au496",
			"reward": {
				"denom": "ucmdx",
				"amount": "2527"
			}
		},
		{
			"address": "comdex1tmlp2kn38m9u05n2jddqecdn7pa53hplsjx7aj",
			"reward": {
				"denom": "ucmdx",
				"amount": "6033"
			}
		},
		{
			"address": "comdex1tml7rljfvgmxvgdms5qcm6s9c97ttep6v7axsk",
			"reward": {
				"denom": "ucmdx",
				"amount": "13429"
			}
		},
		{
			"address": "comdex1tuz39ug3wq332u5w8nt4s9k59dp8aq24jf0mul",
			"reward": {
				"denom": "ucmdx",
				"amount": "6154"
			}
		},
		{
			"address": "comdex1tuyzae8ncxneum3c8m47hrv5jmkrfrh2sv6mzh",
			"reward": {
				"denom": "ucmdx",
				"amount": "77876"
			}
		},
		{
			"address": "comdex1tu8ejq8uph3gsdt07ctxnhwdrugggt3g8whqfc",
			"reward": {
				"denom": "ucmdx",
				"amount": "995"
			}
		},
		{
			"address": "comdex1tug7g6hh747m50n7wn297zu7lpkujnkr3vmkyc",
			"reward": {
				"denom": "ucmdx",
				"amount": "24330"
			}
		},
		{
			"address": "comdex1tufvt8xqgjkc09y00dhgwzv7cv0apa6g4ymxwf",
			"reward": {
				"denom": "ucmdx",
				"amount": "181"
			}
		},
		{
			"address": "comdex1tu2yhmefqxe9ay3mxmv0af8kxdq8qe7xwgqk0x",
			"reward": {
				"denom": "ucmdx",
				"amount": "189"
			}
		},
		{
			"address": "comdex1tutqjxfp8eppwrerr40m5my5z5amepfvvn7jhx",
			"reward": {
				"denom": "ucmdx",
				"amount": "521"
			}
		},
		{
			"address": "comdex1tut25u7zmrsjhp80mh2yf3ekh6d376jl8adck5",
			"reward": {
				"denom": "ucmdx",
				"amount": "325094"
			}
		},
		{
			"address": "comdex1tuv0slmdq8p9uehsjn4ugycdpg896tuhnmk0qt",
			"reward": {
				"denom": "ucmdx",
				"amount": "28"
			}
		},
		{
			"address": "comdex1tuwtl5zth5udv4k5kfykkpf28dqkx27gldyf6a",
			"reward": {
				"denom": "ucmdx",
				"amount": "12336"
			}
		},
		{
			"address": "comdex1tu5l7494unhjw738774ma3dp62pw3el2ed8sw5",
			"reward": {
				"denom": "ucmdx",
				"amount": "5934"
			}
		},
		{
			"address": "comdex1tukxzdp0cu8wf30nwzcwhu6xu32gqaqee0d2mm",
			"reward": {
				"denom": "ucmdx",
				"amount": "13962"
			}
		},
		{
			"address": "comdex1tucrzvxlhtqvk8cvwn46uummx54r6ytxfs88j5",
			"reward": {
				"denom": "ucmdx",
				"amount": "88"
			}
		},
		{
			"address": "comdex1tulnw68jdjjlp2sw3ncj6wj930md4cpz5eccze",
			"reward": {
				"denom": "ucmdx",
				"amount": "67730"
			}
		},
		{
			"address": "comdex1taqdzsm4kt6h2chkdt6a8yp9v826w3rq4l3vs5",
			"reward": {
				"denom": "ucmdx",
				"amount": "1745"
			}
		},
		{
			"address": "comdex1taqs55wngnr24fdxxj8kr2qyyal02nvyp7y3e4",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1tapjtvaas6h7aludtf59zfq5s8u6y9njp8pyqd",
			"reward": {
				"denom": "ucmdx",
				"amount": "41787"
			}
		},
		{
			"address": "comdex1tar8jwydn4wyxxtefr8jrrzemcp2l07ywr7pug",
			"reward": {
				"denom": "ucmdx",
				"amount": "7371"
			}
		},
		{
			"address": "comdex1tayy0stl4sg6c8fka7cksg469jufkdt74x9r8p",
			"reward": {
				"denom": "ucmdx",
				"amount": "179221"
			}
		},
		{
			"address": "comdex1ta9p7lwkcc2myfmn8neq4vvns403wwwru43s7s",
			"reward": {
				"denom": "ucmdx",
				"amount": "17617"
			}
		},
		{
			"address": "comdex1ta9d38getyj2c8my5yd52whqdxgtlhtauctvd3",
			"reward": {
				"denom": "ucmdx",
				"amount": "1236"
			}
		},
		{
			"address": "comdex1tax9uxdd42vanvd9vxs2carfpvjpu5sv4un4ny",
			"reward": {
				"denom": "ucmdx",
				"amount": "8861"
			}
		},
		{
			"address": "comdex1ta8nyt27xuv856xslp504878etvndwqvccp0vc",
			"reward": {
				"denom": "ucmdx",
				"amount": "175"
			}
		},
		{
			"address": "comdex1ta2zej0dm6av2mjxps9rsyqm6jru5jljtwrs0d",
			"reward": {
				"denom": "ucmdx",
				"amount": "1736"
			}
		},
		{
			"address": "comdex1tavf9pe3fl94hpl37u08mwt6p5l3fn4wwnzzck",
			"reward": {
				"denom": "ucmdx",
				"amount": "8724"
			}
		},
		{
			"address": "comdex1tav74fcugcz89cw3z3shekn6tct2hxcvh0kgcc",
			"reward": {
				"denom": "ucmdx",
				"amount": "12946"
			}
		},
		{
			"address": "comdex1tady5gfzzwlzcjkmn8suzxp0uxf54vvx87hghh",
			"reward": {
				"denom": "ucmdx",
				"amount": "353"
			}
		},
		{
			"address": "comdex1tascup98m4hr58vsqmeswvqej42gjrqce2feeq",
			"reward": {
				"denom": "ucmdx",
				"amount": "20444"
			}
		},
		{
			"address": "comdex1ta3kpc5uwa358hul6nsqrhcegqj4qnm9lmxw65",
			"reward": {
				"denom": "ucmdx",
				"amount": "8625"
			}
		},
		{
			"address": "comdex1ta5vh9x7qjy98r909rm4amp97ukegm5futscca",
			"reward": {
				"denom": "ucmdx",
				"amount": "6380"
			}
		},
		{
			"address": "comdex1ta4fju044uz2zldn7yc0njmh96gtweymrzaefj",
			"reward": {
				"denom": "ucmdx",
				"amount": "1429"
			}
		},
		{
			"address": "comdex1tackp7pmk7adr5gcu43h5pev8qq4p0k6wx2js7",
			"reward": {
				"denom": "ucmdx",
				"amount": "10891"
			}
		},
		{
			"address": "comdex1tacm3u9dmd84jw46y4mzx8tjaeu0tvtkxm57ax",
			"reward": {
				"denom": "ucmdx",
				"amount": "1582"
			}
		},
		{
			"address": "comdex1taeesxxhehgputdy65ulfnyv6eghaqcqkrpmkw",
			"reward": {
				"denom": "ucmdx",
				"amount": "2024"
			}
		},
		{
			"address": "comdex1taeue9h65cn5ch35y4v2rrwenm2f7wdfk5yg0u",
			"reward": {
				"denom": "ucmdx",
				"amount": "2194"
			}
		},
		{
			"address": "comdex1ta749htc8gwtm05mkcqkaz5cpgx7pzepsnvdmj",
			"reward": {
				"denom": "ucmdx",
				"amount": "381"
			}
		},
		{
			"address": "comdex1ta7hjdrsxr97u85y9gf00efwwzncwqp7wuw973",
			"reward": {
				"denom": "ucmdx",
				"amount": "7191"
			}
		},
		{
			"address": "comdex1tale992h403jrputvgdpedl0ajwmazedj42e6x",
			"reward": {
				"denom": "ucmdx",
				"amount": "354"
			}
		},
		{
			"address": "comdex1t7ps40sztmg7t3lldppuhff2j6656r7xp0h6jr",
			"reward": {
				"denom": "ucmdx",
				"amount": "9132"
			}
		},
		{
			"address": "comdex1t78aljjpjlejj7z6ter7e0xrezjfdaz4du3s6c",
			"reward": {
				"denom": "ucmdx",
				"amount": "1750"
			}
		},
		{
			"address": "comdex1t72lu76qan2t28zdljrnxnrjkewey2awpv3g0f",
			"reward": {
				"denom": "ucmdx",
				"amount": "2058"
			}
		},
		{
			"address": "comdex1t7dwlnpaj9fg0t8mpaaruuyphnas654r782sp0",
			"reward": {
				"denom": "ucmdx",
				"amount": "1658"
			}
		},
		{
			"address": "comdex1t7s23l8854eveqe9petll4drwmvxy2la0l3z40",
			"reward": {
				"denom": "ucmdx",
				"amount": "203"
			}
		},
		{
			"address": "comdex1t73tt28s7tvjqxxy8avv247ukh6wynqffevfdx",
			"reward": {
				"denom": "ucmdx",
				"amount": "24708"
			}
		},
		{
			"address": "comdex1t759h7sanngfdflzrtjgr7p3rxyn7km9keyv2y",
			"reward": {
				"denom": "ucmdx",
				"amount": "38821"
			}
		},
		{
			"address": "comdex1t752lygwmfr3zf7ufmzmwdez8aphmslq9l6y8g",
			"reward": {
				"denom": "ucmdx",
				"amount": "1763"
			}
		},
		{
			"address": "comdex1t750tywhzkucn5l23esx5l2y50fph3040hyzqc",
			"reward": {
				"denom": "ucmdx",
				"amount": "177"
			}
		},
		{
			"address": "comdex1t747p8ktusmwne7nkvfxs9spg7elt72duheysp",
			"reward": {
				"denom": "ucmdx",
				"amount": "833"
			}
		},
		{
			"address": "comdex1t7kggne7kpsrctfkau00cnkp4lu34u7f0rh8tv",
			"reward": {
				"denom": "ucmdx",
				"amount": "286"
			}
		},
		{
			"address": "comdex1t76qq26z5waqlxcv67n5f8l873vv709jyvz2vc",
			"reward": {
				"denom": "ucmdx",
				"amount": "539"
			}
		},
		{
			"address": "comdex1t762uenjgjcz2ys8f3kz3ewjyzf5lfg3na9veg",
			"reward": {
				"denom": "ucmdx",
				"amount": "0"
			}
		},
		{
			"address": "comdex1tl9g9r7feqjvwe4y8ejeqqnw4hrr2h9l3y7m0h",
			"reward": {
				"denom": "ucmdx",
				"amount": "10215"
			}
		},
		{
			"address": "comdex1tl9fwdu2r6vflvz4pae7nnd373qn5ymhgr07t4",
			"reward": {
				"denom": "ucmdx",
				"amount": "6307"
			}
		},
		{
			"address": "comdex1tlgp5rmt74qpcpwhlgfqh39grzx7cyevfzjgnh",
			"reward": {
				"denom": "ucmdx",
				"amount": "166"
			}
		},
		{
			"address": "comdex1tl2p5pz70zjp4s60gr5he0cljwjezpmkghyr0s",
			"reward": {
				"denom": "ucmdx",
				"amount": "9021"
			}
		},
		{
			"address": "comdex1tl2hdfwr3wcgymz9xyw7hzpwv9ad3n3386n0ke",
			"reward": {
				"denom": "ucmdx",
				"amount": "206"
			}
		},
		{
			"address": "comdex1tltykqs6j4zvjeuj8eywfmajhxzvwpx2r0aw94",
			"reward": {
				"denom": "ucmdx",
				"amount": "176"
			}
		},
		{
			"address": "comdex1tlt6amykuu5vf87gwk29y0034td5fuzycc6hk2",
			"reward": {
				"denom": "ucmdx",
				"amount": "372"
			}
		},
		{
			"address": "comdex1tlw3up46dz8mgemvqqcvflgkkwa2qm4rghd8jw",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1tl30hruvh7gsrhsz8p884pfnav8twlpclg0qwl",
			"reward": {
				"denom": "ucmdx",
				"amount": "28"
			}
		},
		{
			"address": "comdex1tl3u4nltmzflkfedymtus2a4y4t0cdstsfwjr0",
			"reward": {
				"denom": "ucmdx",
				"amount": "3254"
			}
		},
		{
			"address": "comdex1tlk9d3w3tkpustelgzlej7yaya3csjc0k5h76e",
			"reward": {
				"denom": "ucmdx",
				"amount": "8868"
			}
		},
		{
			"address": "comdex1tlkf0fddrz4tur69uz8p94ccg07mwha0dmgauq",
			"reward": {
				"denom": "ucmdx",
				"amount": "1421"
			}
		},
		{
			"address": "comdex1tlmduy8lamnvptxyvc2r7uenmc70mque6yw38a",
			"reward": {
				"denom": "ucmdx",
				"amount": "1793"
			}
		},
		{
			"address": "comdex1tlax65f27j2svx40xgep2kzpj5agrepmgq077w",
			"reward": {
				"denom": "ucmdx",
				"amount": "356"
			}
		},
		{
			"address": "comdex1tl7jrxre8l9my0vfjuwqzjy8qlvcx7mspayqst",
			"reward": {
				"denom": "ucmdx",
				"amount": "633"
			}
		},
		{
			"address": "comdex1tllfl5nk07ac9zmaux5sqtv4wn9wcapazgvv55",
			"reward": {
				"denom": "ucmdx",
				"amount": "1758"
			}
		},
		{
			"address": "comdex1vqp88kzuvllwr57m6r38dghy6w64qvv7rtsjw6",
			"reward": {
				"denom": "ucmdx",
				"amount": "2653"
			}
		},
		{
			"address": "comdex1vqrs82ru0d8qrzkhd8wmwt4uasuuc4kc0yrml4",
			"reward": {
				"denom": "ucmdx",
				"amount": "1949"
			}
		},
		{
			"address": "comdex1vq97vj6hsrvg8569eclpnpsvs6nk7jll97de4u",
			"reward": {
				"denom": "ucmdx",
				"amount": "6325"
			}
		},
		{
			"address": "comdex1vqx88lldhh8guvn6jn80zqe07w2s7ayzu3ng0v",
			"reward": {
				"denom": "ucmdx",
				"amount": "275"
			}
		},
		{
			"address": "comdex1vq8w8c757w5x08ptudvzd8vmd06xsuhmk09fa8",
			"reward": {
				"denom": "ucmdx",
				"amount": "7214"
			}
		},
		{
			"address": "comdex1vqwyzsnfkqrg4pyq4jazgr3n76vwuzupw4ql4j",
			"reward": {
				"denom": "ucmdx",
				"amount": "16928"
			}
		},
		{
			"address": "comdex1vqwx49ymv2czkuk4vrtdy0rkz8laec9luz45x0",
			"reward": {
				"denom": "ucmdx",
				"amount": "100"
			}
		},
		{
			"address": "comdex1vqw37ku8he8xv635kvjt2kf0hm27z9muhy2a73",
			"reward": {
				"denom": "ucmdx",
				"amount": "6833"
			}
		},
		{
			"address": "comdex1vq37p20utgz0lh7rcv6t2d7tdep4kpjc5wcqga",
			"reward": {
				"denom": "ucmdx",
				"amount": "18"
			}
		},
		{
			"address": "comdex1vq5yzvkh3h222sqrfvua5hzswq0x9cm6wwcxv0",
			"reward": {
				"denom": "ucmdx",
				"amount": "14"
			}
		},
		{
			"address": "comdex1vq45jx2trk3wfqfyff5wsd6zpzpq8tehszln9l",
			"reward": {
				"denom": "ucmdx",
				"amount": "5056"
			}
		},
		{
			"address": "comdex1vqkdaalu3edekdadvnx4qewpzjpf9wmf75hzwq",
			"reward": {
				"denom": "ucmdx",
				"amount": "1749"
			}
		},
		{
			"address": "comdex1vqkjg003xs34srgcr0hnpuz3rcxd8l6e2suagn",
			"reward": {
				"denom": "ucmdx",
				"amount": "6076"
			}
		},
		{
			"address": "comdex1vqcudc0fk3l5uzwazj5wa8gvkszlsqtn68kwrj",
			"reward": {
				"denom": "ucmdx",
				"amount": "531"
			}
		},
		{
			"address": "comdex1vqem22tk25epwzu0zcjdmd8famwezy3nc73lzh",
			"reward": {
				"denom": "ucmdx",
				"amount": "19"
			}
		},
		{
			"address": "comdex1vq6tmauue4at8yy3av2l7w8mcvm9k3eslupzep",
			"reward": {
				"denom": "ucmdx",
				"amount": "41711"
			}
		},
		{
			"address": "comdex1vq6danhktg0eyqchjx008raf3jnslawlwrf4lr",
			"reward": {
				"denom": "ucmdx",
				"amount": "16"
			}
		},
		{
			"address": "comdex1vqmcvycmg202clnx8ev9rjgz7hwtw2ms6tns3p",
			"reward": {
				"denom": "ucmdx",
				"amount": "176"
			}
		},
		{
			"address": "comdex1vquuncpxjru7nk7c2qpv7qnagk008ppsyv5x2d",
			"reward": {
				"denom": "ucmdx",
				"amount": "1259"
			}
		},
		{
			"address": "comdex1vqlsp77euem8nnhky8qjr5mlkcpclp8as3wc9f",
			"reward": {
				"denom": "ucmdx",
				"amount": "3895"
			}
		},
		{
			"address": "comdex1vqlj2l6r5qc7a2zt2dmsl3wwfuu85ergdxrq4h",
			"reward": {
				"denom": "ucmdx",
				"amount": "118917"
			}
		},
		{
			"address": "comdex1vpzlwskzcvy60lasrwa5xepztfp95vsdfxeq78",
			"reward": {
				"denom": "ucmdx",
				"amount": "17131"
			}
		},
		{
			"address": "comdex1vp934uw50u8erdf5c65v3unacecg9rcdj7jq9x",
			"reward": {
				"denom": "ucmdx",
				"amount": "5359"
			}
		},
		{
			"address": "comdex1vpftx25dymd3s0282ln36waqnc2jags2gz35tp",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1vp2mtfv6hpnauku5vchfhlkmvqwz529z3tf2fy",
			"reward": {
				"denom": "ucmdx",
				"amount": "168"
			}
		},
		{
			"address": "comdex1vpvdtmtl3wa2czqj7clz484yjhdvwh3rhsqvx5",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1vpwgqgda7n3lqppz3rhf8cjxz2m6zh8lc2nh56",
			"reward": {
				"denom": "ucmdx",
				"amount": "142295"
			}
		},
		{
			"address": "comdex1vpwsfc0gxc94xkylkq3dttxjuak2uejzgkwlxk",
			"reward": {
				"denom": "ucmdx",
				"amount": "14008"
			}
		},
		{
			"address": "comdex1vpwnq5sfxwgl94tmxvshjvwmgcnl3ykujghrlv",
			"reward": {
				"denom": "ucmdx",
				"amount": "71807"
			}
		},
		{
			"address": "comdex1vp0grydm2sp80wskl9r37ze6n9y8x3v26enhvq",
			"reward": {
				"denom": "ucmdx",
				"amount": "1417"
			}
		},
		{
			"address": "comdex1vpszg27hrmxklyr0wc3huv0ypw9lxl9pak242a",
			"reward": {
				"denom": "ucmdx",
				"amount": "35217"
			}
		},
		{
			"address": "comdex1vpshzex6pt7n3mjm3fpyzkh6q05s4ef0y68f0y",
			"reward": {
				"denom": "ucmdx",
				"amount": "398"
			}
		},
		{
			"address": "comdex1vps63smcfl37z0w8psqsgue73gnu74cmde0tpz",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1vp4qvumd3993mse4pllq5hmxdqjm2m6m48xx43",
			"reward": {
				"denom": "ucmdx",
				"amount": "71"
			}
		},
		{
			"address": "comdex1vphre7fmvkjknxlxxzmzaaz6rye50rx6wl42wm",
			"reward": {
				"denom": "ucmdx",
				"amount": "31500"
			}
		},
		{
			"address": "comdex1vphsn9aglpff3h70xsuk4k5d9y9up5yqevvflp",
			"reward": {
				"denom": "ucmdx",
				"amount": "20"
			}
		},
		{
			"address": "comdex1vphh0qq88y37cj6pza0qjydymm9ac3wugqwfff",
			"reward": {
				"denom": "ucmdx",
				"amount": "143"
			}
		},
		{
			"address": "comdex1vpc39jmmf9p37dwvcgxy4gks9weqpzqaxjg8c0",
			"reward": {
				"denom": "ucmdx",
				"amount": "180"
			}
		},
		{
			"address": "comdex1vpujs4zc4ut02d5pjwf83qswz7m843ssvmw9pp",
			"reward": {
				"denom": "ucmdx",
				"amount": "1270"
			}
		},
		{
			"address": "comdex1vzqc64g9w4rak0cu34nwhxg6yeftkckzunfmly",
			"reward": {
				"denom": "ucmdx",
				"amount": "10612"
			}
		},
		{
			"address": "comdex1vzxt2akmr76ddg7epxklnwmdjz5m04l3m65hz9",
			"reward": {
				"denom": "ucmdx",
				"amount": "223"
			}
		},
		{
			"address": "comdex1vzx49qugge8fqg2cmasn0w9wzyr275rxr4wwd7",
			"reward": {
				"denom": "ucmdx",
				"amount": "149"
			}
		},
		{
			"address": "comdex1vz2fcvuarywcp7cej7rx3pt8v5gaz5r993ydz9",
			"reward": {
				"denom": "ucmdx",
				"amount": "35838"
			}
		},
		{
			"address": "comdex1vzv4ukqclpvljs7keqhng232dz5yrhwtyq4kd8",
			"reward": {
				"denom": "ucmdx",
				"amount": "20106"
			}
		},
		{
			"address": "comdex1vz06pvcl6d3cetfgz5642ms9uc08as8csl0gwm",
			"reward": {
				"denom": "ucmdx",
				"amount": "1090"
			}
		},
		{
			"address": "comdex1vzsdagxzjxephp74t38u48auvd7mlvv0cn830p",
			"reward": {
				"denom": "ucmdx",
				"amount": "6306"
			}
		},
		{
			"address": "comdex1vz3tsw3099pheq9cedztvahdg5xn2chwu4frk5",
			"reward": {
				"denom": "ucmdx",
				"amount": "62422"
			}
		},
		{
			"address": "comdex1vz4eqy3q6fudvwghnn2s4zjuza403qjtn4q4t9",
			"reward": {
				"denom": "ucmdx",
				"amount": "12903"
			}
		},
		{
			"address": "comdex1vzcgznld3uhaq3elc638kjda4ecy9dupz8yvr5",
			"reward": {
				"denom": "ucmdx",
				"amount": "4047"
			}
		},
		{
			"address": "comdex1vzmd5m40pucsp0ura6gfq7cy664rzn33klk4wf",
			"reward": {
				"denom": "ucmdx",
				"amount": "96368"
			}
		},
		{
			"address": "comdex1vzazv307v8l68yrnz00u4rc7t7dmma4qw00uv7",
			"reward": {
				"denom": "ucmdx",
				"amount": "4383"
			}
		},
		{
			"address": "comdex1vza94wc9ds8j45rxc8huhvgnv6n3vpmcs2rr6c",
			"reward": {
				"denom": "ucmdx",
				"amount": "83939"
			}
		},
		{
			"address": "comdex1vz73yxk594xlxeveu7muafzkynejjcgr782s4q",
			"reward": {
				"denom": "ucmdx",
				"amount": "890"
			}
		},
		{
			"address": "comdex1vrpfls8ynj4fxz535pqyyektl55w06zj5sat78",
			"reward": {
				"denom": "ucmdx",
				"amount": "20455"
			}
		},
		{
			"address": "comdex1vrrj4kw52kchqzgzpn5zw46dha2cdp6h46qu0e",
			"reward": {
				"denom": "ucmdx",
				"amount": "17"
			}
		},
		{
			"address": "comdex1vr98m37kchf5ffucq7l8p7wcmhncv3yplwg5d6",
			"reward": {
				"denom": "ucmdx",
				"amount": "166"
			}
		},
		{
			"address": "comdex1vrxagmlfkk29sxh3ynr5dgcmkvgd0s9fkal0cg",
			"reward": {
				"denom": "ucmdx",
				"amount": "587"
			}
		},
		{
			"address": "comdex1vr8vusfdf6wkvhheeq85kl75h2mjz6j9ne8c9x",
			"reward": {
				"denom": "ucmdx",
				"amount": "1799"
			}
		},
		{
			"address": "comdex1vrgqdr95mra6ud99spqk6a4k8jzu43n0hryeaa",
			"reward": {
				"denom": "ucmdx",
				"amount": "17128"
			}
		},
		{
			"address": "comdex1vrg6g8jncpyjn03f6gqnks27adjd2c2x99mrz6",
			"reward": {
				"denom": "ucmdx",
				"amount": "114842"
			}
		},
		{
			"address": "comdex1vr294smddrvwntu2f65emyt7q08naaesh269z6",
			"reward": {
				"denom": "ucmdx",
				"amount": "12402"
			}
		},
		{
			"address": "comdex1vrv22325xw757vmlvqx65k24cgzkavxkrgzzuq",
			"reward": {
				"denom": "ucmdx",
				"amount": "57"
			}
		},
		{
			"address": "comdex1vrwme6ss0q440x7erj3xgtpyzns3cva08dlh3m",
			"reward": {
				"denom": "ucmdx",
				"amount": "12290"
			}
		},
		{
			"address": "comdex1vr0hndzfeumnhavv0rl6mcgpw44zfyc80hu2jl",
			"reward": {
				"denom": "ucmdx",
				"amount": "8691"
			}
		},
		{
			"address": "comdex1vrsdsrl2gnz795jx46f6x7a6fhz7k7697s2ngh",
			"reward": {
				"denom": "ucmdx",
				"amount": "19785"
			}
		},
		{
			"address": "comdex1vrjdd20jz2sk0gsuak57unhgrhqdv806e0xzu8",
			"reward": {
				"denom": "ucmdx",
				"amount": "2221"
			}
		},
		{
			"address": "comdex1vrjl8pvxajmc5efa464g64xkgmwzfaztkyv9fs",
			"reward": {
				"denom": "ucmdx",
				"amount": "1438"
			}
		},
		{
			"address": "comdex1vrhs5ezreer74nhy0kdgctyxvl7vkvta9uu07d",
			"reward": {
				"denom": "ucmdx",
				"amount": "163"
			}
		},
		{
			"address": "comdex1vrhlzf7g7xflrzuvxuj67897npnphkwt8jh9wx",
			"reward": {
				"denom": "ucmdx",
				"amount": "2899"
			}
		},
		{
			"address": "comdex1vrml8rkxcstrneqv208wmd9f3yd5uwkwhny08c",
			"reward": {
				"denom": "ucmdx",
				"amount": "45"
			}
		},
		{
			"address": "comdex1vr7yduk6ygqnn547wdt3444pq96actjspdlgs4",
			"reward": {
				"denom": "ucmdx",
				"amount": "3137"
			}
		},
		{
			"address": "comdex1vr7afd4zyk8rfpl5q7ekjvpp0gqeee8dk4sacj",
			"reward": {
				"denom": "ucmdx",
				"amount": "1190"
			}
		},
		{
			"address": "comdex1vyqnpw690s96qujnvt5ug38qvmj5zyw7f6zdqh",
			"reward": {
				"denom": "ucmdx",
				"amount": "6982"
			}
		},
		{
			"address": "comdex1vypctkdnn52z9cxq5wn7v2uchmp46v3ucy3ckc",
			"reward": {
				"denom": "ucmdx",
				"amount": "271"
			}
		},
		{
			"address": "comdex1vyz48nxrg3eauksq69602pzggtsytexme8mkcz",
			"reward": {
				"denom": "ucmdx",
				"amount": "1757"
			}
		},
		{
			"address": "comdex1vyrp5ma6d2kzjww4crl4juefw79kstu4ls9aqx",
			"reward": {
				"denom": "ucmdx",
				"amount": "3879"
			}
		},
		{
			"address": "comdex1vyy4furd7qndtevu2x0qufcen58yypers4dev2",
			"reward": {
				"denom": "ucmdx",
				"amount": "3385"
			}
		},
		{
			"address": "comdex1vyxzsrklv9r48c9dykzdtyrchccc5anu7097c6",
			"reward": {
				"denom": "ucmdx",
				"amount": "3084"
			}
		},
		{
			"address": "comdex1vy2ppv6md42l6arhhg63dm89gdl4mh26swgx58",
			"reward": {
				"denom": "ucmdx",
				"amount": "2238"
			}
		},
		{
			"address": "comdex1vyt3mzyv4r5ljjkr4q9yagkvcxmsj9hjrcnhrw",
			"reward": {
				"denom": "ucmdx",
				"amount": "8907"
			}
		},
		{
			"address": "comdex1vytjusuzs8rftxc69v4uepvjk39j7mh508weea",
			"reward": {
				"denom": "ucmdx",
				"amount": "18553"
			}
		},
		{
			"address": "comdex1vyw9yft6nqm8gj4ggchjrwnvjckz05rqruzm25",
			"reward": {
				"denom": "ucmdx",
				"amount": "14571"
			}
		},
		{
			"address": "comdex1vy0j22htwcke8uz8vpyd5arefl0tawpxdm3a8q",
			"reward": {
				"denom": "ucmdx",
				"amount": "96"
			}
		},
		{
			"address": "comdex1vysg8pgmlnxrx9knlm0rn76mtedvn6q46nvy83",
			"reward": {
				"denom": "ucmdx",
				"amount": "19406"
			}
		},
		{
			"address": "comdex1vys66e04d6r4rtu7vrkrquzrd5gt23enj0clv3",
			"reward": {
				"denom": "ucmdx",
				"amount": "9701"
			}
		},
		{
			"address": "comdex1vy4c3rh44jp9jeuuhahgecunc339u8n549jkhn",
			"reward": {
				"denom": "ucmdx",
				"amount": "5328"
			}
		},
		{
			"address": "comdex1vy472j5kvkzeqzwp7vzylcch0pjk6227tpze6y",
			"reward": {
				"denom": "ucmdx",
				"amount": "48145"
			}
		},
		{
			"address": "comdex1vye0qm6l35k55qm8g5q503ahre8w9zm0vstfee",
			"reward": {
				"denom": "ucmdx",
				"amount": "171"
			}
		},
		{
			"address": "comdex1vy6yv8vkmk0ye45u5uymhg29nrp793nxv9ng56",
			"reward": {
				"denom": "ucmdx",
				"amount": "847"
			}
		},
		{
			"address": "comdex1vy6w3cna8fsx6jfl83e8cg4x4uvfhxmq6y37jn",
			"reward": {
				"denom": "ucmdx",
				"amount": "3800"
			}
		},
		{
			"address": "comdex1vym5524ekqgzy7dcqhcksyl4v8hsul0cq2m4y2",
			"reward": {
				"denom": "ucmdx",
				"amount": "1409"
			}
		},
		{
			"address": "comdex1v9p6kvze04ns55qmv6s0xpp8qm05d2ajxnrlhg",
			"reward": {
				"denom": "ucmdx",
				"amount": "284"
			}
		},
		{
			"address": "comdex1v9x3qq7zdh8t4vt8dzx4q0wh67ty3dtaql9h80",
			"reward": {
				"denom": "ucmdx",
				"amount": "14262"
			}
		},
		{
			"address": "comdex1v9xa0d7hadv6t8nazzqvge4ycpqty4jw45xnjk",
			"reward": {
				"denom": "ucmdx",
				"amount": "9790"
			}
		},
		{
			"address": "comdex1v9ffgwhulrgrq2dpuggpwmjhngqe2hgdfzmpc3",
			"reward": {
				"denom": "ucmdx",
				"amount": "14165"
			}
		},
		{
			"address": "comdex1v9fda4ymnkhwqy947gl4rnvtm8nme538ne8h46",
			"reward": {
				"denom": "ucmdx",
				"amount": "368"
			}
		},
		{
			"address": "comdex1v92lczqymdmq5qaxptuswjkwykrrenjm9huprx",
			"reward": {
				"denom": "ucmdx",
				"amount": "10976"
			}
		},
		{
			"address": "comdex1v9v7dtfk64vz95597dh5e48zla9v650nua8kkt",
			"reward": {
				"denom": "ucmdx",
				"amount": "209"
			}
		},
		{
			"address": "comdex1v9vlktls7h2fdequpsxk49qfz9u9jxq07w625x",
			"reward": {
				"denom": "ucmdx",
				"amount": "1723"
			}
		},
		{
			"address": "comdex1v9w49rtdpkwxeje8xnx9pnf9z9hyr2e0tskur8",
			"reward": {
				"denom": "ucmdx",
				"amount": "1746"
			}
		},
		{
			"address": "comdex1v94wupxyjt293tghfrj8cn69s67guznxet72mr",
			"reward": {
				"denom": "ucmdx",
				"amount": "2015"
			}
		},
		{
			"address": "comdex1v9ktpsraqcw2ta943m82ukjlgak8vdd2nkk62p",
			"reward": {
				"denom": "ucmdx",
				"amount": "167"
			}
		},
		{
			"address": "comdex1v9epyje7vcfvw7sc5qh49xjupaw046jglwu4xy",
			"reward": {
				"denom": "ucmdx",
				"amount": "7058"
			}
		},
		{
			"address": "comdex1v9eup5s52sc0p3hh0kcc5sq24lkhyypwvjpukj",
			"reward": {
			}