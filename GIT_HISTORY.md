History of the original repo

commit ca9d341f11731ffeabf0bce1036a7d0e3d942b88
Author: shipperizer <alexcabb@gmail.com>
Date:   Mon Feb 13 15:41:22 2023 +0000

    ci: use only one kafka with multiple listeners

commit 7efbe9c2b29794a2aae7886958086991238d3353
Author: shipperizer <alexcabb@gmail.com>
Date:   Thu Feb 9 08:51:30 2023 +0000

    ci: add kafka sasl cnx string env var

commit a430294449387a07adb2fd467b2e8ab7e986a94c
Author: shipperizer <alexcabb@gmail.com>
Date:   Thu Feb 9 08:51:11 2023 +0000

    fix: drop commented code

commit 0be3f3211b86468f762dc985b0bf70eae757679a
Author: shipperizer <alexcabb@gmail.com>
Date:   Wed Feb 8 17:06:30 2023 +0000

    feat: implement a sasl config interface

commit 6fd37fc907d93ea05b0201e9143a73fcd0a37a21
Author: shipperizer <alexcabb@gmail.com>
Date:   Wed Feb 8 17:05:56 2023 +0000

    feat: reader and writer deal with sasl config during connection renewal

commit b76bc86ea2b952006c16d7516a46f0d950d96eb7
Author: shipperizer <alexcabb@gmail.com>
Date:   Wed Feb 8 17:05:33 2023 +0000

    feat: autorefresher passes sasl as well as tls config now

commit 1aa944c8adb7362b943adcda23bd0702863754ba
Author: shipperizer <alexcabb@gmail.com>
Date:   Wed Feb 8 17:05:01 2023 +0000

    feat: add sasl config

commit 7f658839485f16e398ce197c78f44c9155788a55
Author: shipperizer <alexcabb@gmail.com>
Date:   Mon Feb 6 21:43:43 2023 +0000

    ci: update deps and auxiliary files

commit 3d435f3ae1c61ab1c40f997cc7a95b21bfcbeaf3
Author: shipperizer <alexcabb@gmail.com>
Date:   Mon Feb 6 21:41:47 2023 +0000

    feat: tls package updates and introduce vault package

commit 51225db15f928e8e5055386f45f410c083722884
Author: shipperizer <alexcabb@gmail.com>
Date:   Mon Feb 6 21:41:19 2023 +0000

    feat: update to follow Renew method changes and interfaces

commit 32f0e1ad455977ce48a67f77f39a99f43b71404d
Author: shipperizer <alexcabb@gmail.com>
Date:   Mon Feb 6 21:40:41 2023 +0000

    feat: drop Renew method returning and instead modify in place

commit 6134cd7d92942333abbe870e5492cfb1394faf03
Author: shipperizer <alexcabb@gmail.com>
Date:   Mon Feb 6 21:40:15 2023 +0000

    feat: shuffle interfaces around

commit 74a4c8456676ed637b0949def9306cd1be97158b
Author: shipperizer <alexcabb@gmail.com>
Date:   Thu Jan 26 10:48:55 2023 +0000

    fix: swap zap for loggerInterface

commit 627acafcca108ac854ee978998684d53dfbb923e
Merge: 13a05d2 deefee3
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Wed Jan 4 11:59:52 2023 +0000

    Merge pull request #21 from wbd-streaming/release-please--branches--main
    
    chore(main): release 2.0.1

commit deefee3e207fd990daf8cae78c4b52a3c5e1f958
Author: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>
Date:   Wed Jan 4 11:53:09 2023 +0000

    chore(main): release 2.0.1

commit 13a05d25b94b029a966eb9cf62232df2184d21db
Merge: bda66e3 be1db5d
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Wed Jan 4 11:44:39 2023 +0000

    Merge pull request #19 from wbd-streaming/dev-v2
    
    dev v2

commit be1db5d182d11cbf469dda7e90537904e616a095
Author: shipperizer <alexcabb@gmail.com>
Date:   Wed Jan 4 11:10:01 2023 +0000

    fix: upgrade paths to use v2

commit dab4ab4e7f062143a45f00f8dbffb10275f9b71f
Author: shipperizer <alexcabb@gmail.com>
Date:   Wed Jan 4 11:08:15 2023 +0000

    build: update to v2, following golang practices

commit bda66e3d81be2e2af46a8d1d78fef003c8a68425
Merge: 8bf8f88 a488435
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Wed Jan 4 10:33:54 2023 +0000

    Merge pull request #17 from wbd-streaming/release-please--branches--main
    
    chore(main): release 2.0.0

commit a488435b9e58724ffca6519d4ab79f8e7664ace9
Author: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>
Date:   Wed Jan 4 10:22:16 2023 +0000

    chore(main): release 2.0.0

commit 8bf8f88a8693ab46121b4f39de6f764f65004f0e
Merge: 0213626 894d323
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Wed Jan 4 10:21:46 2023 +0000

    Merge pull request #14 from wbd-streaming/lib/crypto
    
    fix: change crypto lib to use SSL Mate due to frozen crypto/pkcs12 and lack of sha256 support

commit 894d32371c1492c12a439b2ed89f4aa1b0d2f4c6
Author: shipperizer <alexcabb@gmail.com>
Date:   Tue Jan 3 22:15:19 2023 +0000

    fix!: use SSLMate pkcs12 library for better support of sha256 signing alg

commit 0a1da0022aee241b2aee98750bd2f9ec2239548d
Author: shipperizer <alexcabb@gmail.com>
Date:   Tue Jan 3 22:13:39 2023 +0000

    ci: add p12 cert with sha256 signing

commit 0213626b74e64abe427721875bf9ed997043487a
Merge: 160b536 1a6cf99
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Fri Dec 23 10:32:27 2022 +0000

    Merge pull request #13 from wbd-streaming/dependabot/github_actions/ossf/scorecard-action-2.1.0
    
    build(deps): bump ossf/scorecard-action from 1.1.1 to 2.1.0

commit 1a6cf99ab0402769de54637333a49119e22a83a0
Author: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>
Date:   Sun Dec 18 15:01:11 2022 +0000

    build(deps): bump ossf/scorecard-action from 1.1.1 to 2.1.0
    
    Bumps [ossf/scorecard-action](https://github.com/ossf/scorecard-action) from 1.1.1 to 2.1.0.
    - [Release notes](https://github.com/ossf/scorecard-action/releases)
    - [Changelog](https://github.com/ossf/scorecard-action/blob/main/RELEASE.md)
    - [Commits](https://github.com/ossf/scorecard-action/compare/v1.1.1...v2.1.0)
    
    ---
    updated-dependencies:
    - dependency-name: ossf/scorecard-action
      dependency-type: direct:production
      update-type: version-update:semver-major
    ...
    
    Signed-off-by: dependabot[bot] <support@github.com>

commit 160b5369b8ba0120dd031723f0adcda3aa75ede9
Merge: fa1ab63 4cf46b8
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Thu Dec 1 09:08:47 2022 +0000

    Merge pull request #12 from wbd-streaming/release-please--branches--main
    
    chore(main): release 1.0.1

commit 4cf46b85f13c1b0b1a90de387cce6f913af120d4
Author: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>
Date:   Wed Nov 30 16:28:01 2022 +0000

    chore(main): release 1.0.1

commit fa1ab63ad9fba0b4ca5715d2f81eba8577543d94
Merge: 67c276b 5dcc42d
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Wed Nov 30 16:27:27 2022 +0000

    Merge pull request #11 from wbd-streaming/wbd/pkg
    
    wbd/pkg

commit 5dcc42dfa28eaed410972abb8a514b4f28926e67
Author: shipperizer <alexcabb@gmail.com>
Date:   Mon Nov 28 12:07:08 2022 +0000

    fix: use wbd-streaming monitoring lib

commit 80501a5d273b55bece9f160aa7b9f6851503037e
Author: shipperizer <alexcabb@gmail.com>
Date:   Mon Nov 28 11:36:05 2022 +0000

    ci: fix sonarqube

commit 67c276bd554c534a5832f5d6bc9219701692d975
Merge: 5b7d6bf 6374036
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Mon Nov 28 11:33:55 2022 +0000

    Merge pull request #10 from wbd-streaming/changelog
    
    docs: separate old changelog history

commit 63740364d1711aaabffe76a30aff124e10533ea6
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Mon Nov 28 10:50:49 2022 +0000

    docs: separate old changelog history

commit 5b7d6bfabfd13698315434a40c5de6ba43632e1e
Merge: be1dc89 ccbbdc3
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Fri Nov 25 11:33:54 2022 +0000

    Merge pull request #9 from wbd-streaming/release-please--branches--main
    
    chore(main): release 1.0.0

commit ccbbdc345d89b64516436914539faa4c8f163666
Author: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>
Date:   Fri Nov 25 10:25:47 2022 +0000

    chore(main): release 1.0.0

commit be1dc89c274d152aed3e666c5ee4a91ba66de6cd
Merge: 5b2303f f3ced19
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Fri Nov 25 10:24:30 2022 +0000

    Merge pull request #8 from wbd-streaming/wbd/pkg
    
    Update package to use newer wbd-streaming path

commit f3ced198e18ba922dfab3886081f36e655336f0e
Author: shipperizer <alexcabb@gmail.com>
Date:   Fri Nov 25 09:50:17 2022 +0000

    docs: update refs on readme

commit fb02385d0eda48e9ba7263a200752c3755637b2c
Author: shipperizer <alexcabb@gmail.com>
Date:   Fri Nov 25 09:27:16 2022 +0000

    fix: drop old code

commit 7a3d0b412d1e2d480ef3eee79ea25ad7fa24c2d1
Author: shipperizer <alexcabb@gmail.com>
Date:   Fri Nov 25 09:22:43 2022 +0000

    ci: use wbd-streaming in GOPRIVATE

commit 93d5252a091a7535e4ebf674a5678a59724edb6d
Author: shipperizer <alexcabb@gmail.com>
Date:   Wed Nov 23 16:59:29 2022 +0000

    ci: adjust team settings

commit e8ad669dddda0e1a4e6459778d55892e29091734
Author: shipperizer <alexcabb@gmail.com>
Date:   Wed Nov 23 16:58:55 2022 +0000

    fix: use newer repo path for imports

commit d4d9c6b558f7b1c60b400f7342b80d52afd6eced
Author: shipperizer <alexcabb@gmail.com>
Date:   Wed Nov 23 16:58:00 2022 +0000

    ci: fix branch triggers to use main

commit f821728750d7774d4cec66e549ea50b283dd7a65
Author: shipperizer <alexcabb@gmail.com>
Date:   Wed Nov 23 16:57:39 2022 +0000

    ci: fix omd folder

commit 5b2303fc5162bcc14a292d2387767e0d77fca8fc
Author: Terraform User <terraform@example.com>
Date:   Tue Oct 25 08:24:03 2022 +0000

    auto updated file

commit 011ebf06055a1cf5221b60a3beab8d3080324005
Merge: 2636838 47475ba
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Tue Nov 15 12:46:58 2022 +0000

    Merge pull request #143 from sbs-discovery-sweden/release-please--branches--master
    
    chore(master): release 1.5.0

commit 47475ba44a7e53dc8a5b915759ed41a7bb1b2464
Author: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>
Date:   Tue Nov 15 11:45:53 2022 +0000

    chore(master): release 1.5.0

commit 263683838e13466e70c0ef850b16b4198551f156
Merge: 24ab203 dc144b6
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Tue Nov 15 11:45:15 2022 +0000

    Merge pull request #142 from sbs-discovery-sweden/tls/p12
    
    tls/p12

commit dc144b67471ecb1096720ccd707b363678449bd6
Author: shipperizer <alexcabb@gmail.com>
Date:   Tue Nov 15 11:10:33 2022 +0000

    fix: dont b64 decode secretBinary

commit 195e8859642b7c33197c70993bde285eb04e7ac4
Author: shipperizer <alexcabb@gmail.com>
Date:   Tue Nov 15 10:46:11 2022 +0000

    feat: support p12 format

commit ba0364f755f32457634c9d1c319f6830499bcd4a
Author: shipperizer <alexcabb@gmail.com>
Date:   Tue Nov 15 10:45:49 2022 +0000

    ci: add p12 secret on makefile target

commit 73cd9098b2dcbe8614cd7309891fa908fac1c418
Author: shipperizer <alexcabb@gmail.com>
Date:   Tue Nov 15 10:45:32 2022 +0000

    ci: downgrade to use go1.17 as minimum version

commit 24ab2035c2e23c45231d467f65ee310232cc0592
Merge: 837d91c af79da2
Author: Ian Rankin <60607903+irdiscovery@users.noreply.github.com>
Date:   Fri Sep 16 11:39:42 2022 +0100

    Merge pull request #113 from sbs-discovery-sweden/release-please--branches--master
    
    chore(master): release 1.4.0

commit af79da2ccb01ca9fa4d006f4870a33bbef55287f
Author: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>
Date:   Fri Sep 16 10:31:12 2022 +0000

    chore(master): release 1.4.0

commit 837d91c3670aa974ae7c46801869384a3b435a05
Merge: c1367b5 38989cb
Author: Ian Rankin <60607903+irdiscovery@users.noreply.github.com>
Date:   Fri Sep 16 11:30:42 2022 +0100

    Merge pull request #115 from sbs-discovery-sweden/LABS-2195-lb
    
    feat: LABS-2195 update dep library

commit 38989cb0de2fff8d18e15f636def801684d48c48
Author: irdiscovery <ian_rankin@discovery.com>
Date:   Thu Sep 15 20:17:17 2022 +0100

    feat: LABS-2195 update dep library

commit c1367b5cec79b6af9a52e286e150bbe1a1475d3d
Merge: bf53bc2 c54ff15
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Thu Sep 15 15:55:49 2022 +0100

    Merge pull request #109 from sbs-discovery-sweden/LABS-2165
    
    feat: LABS-2165 update golang version to 1.19

commit c54ff15181560bce8d69dcc05172195250071cff
Author: irdiscovery <ian_rankin@discovery.com>
Date:   Wed Sep 14 12:07:36 2022 +0100

    feat: LABS-2165 update golang version to 1.19

commit bf53bc2ffa84c985a9ae7e09ca9cec2bf07ecfdc
Merge: e158f60 d980987
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Tue Sep 6 15:20:25 2022 +0100

    Merge pull request #104 from sbs-discovery-sweden/release-please--branches--master
    
    chore(master): release 1.3.1

commit d9809877bfbc8f1955eab02bfa1eaa9dda8eb3ce
Author: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>
Date:   Tue Sep 6 14:13:24 2022 +0000

    chore(master): release 1.3.1

commit e158f60dfaa8ac9e1d2d440f488bc607ce00436c
Merge: d276851 bac6640
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Tue Sep 6 15:12:49 2022 +0100

    Merge pull request #103 from sbs-discovery-sweden/error/handling
    
    fix: add error logging

commit bac6640bf5ef6c51a8aea0736ad5374b7a114f7e
Author: shipperizer <alexcabb@gmail.com>
Date:   Tue Sep 6 15:05:32 2022 +0100

    fix: add error logging

commit d276851c905d6b3b6c694325f0aabaff57caa276
Merge: 41acf57 f3155a7
Author: Mariama <33070292+Mariamakbbh@users.noreply.github.com>
Date:   Tue Aug 9 15:40:48 2022 +0100

    Merge pull request #69 from sbs-discovery-sweden/now-yaml
    
    chore: add operational metadata now.yaml file

commit f3155a7ce87ecda46a46db160267f9cb5cbbed2c
Author: Mariama <33070292+Mariamakbbh@users.noreply.github.com>
Date:   Tue Aug 9 14:47:09 2022 +0100

    Update now.yaml

commit 02b07aa4579f5489352f3f15c65f178d52da0dcd
Author: John Veldboom <john_veldboom@discovery.com>
Date:   Mon Aug 1 21:36:59 2022 +0000

    chore: add operational metadata now.yaml file

commit 41acf576cb5ece7ef456b2ebcf682907f5e5f650
Merge: 1794b91 f952858
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Mon Jun 27 17:31:48 2022 +0100

    Merge pull request #44 from sbs-discovery-sweden/dependabot/github_actions/actions/stale-5
    
    build(deps): bump actions/stale from 3 to 5

commit f952858a3b3215a93719a7c7c9e6ca69aa17690e
Author: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>
Date:   Mon Jun 27 16:24:46 2022 +0000

    build(deps): bump actions/stale from 3 to 5
    
    Bumps [actions/stale](https://github.com/actions/stale) from 3 to 5.
    - [Release notes](https://github.com/actions/stale/releases)
    - [Changelog](https://github.com/actions/stale/blob/main/CHANGELOG.md)
    - [Commits](https://github.com/actions/stale/compare/v3...v5)
    
    ---
    updated-dependencies:
    - dependency-name: actions/stale
      dependency-type: direct:production
      update-type: version-update:semver-major
    ...
    
    Signed-off-by: dependabot[bot] <support@github.com>

commit 1794b91a278e85bbdddd39c991137fea8a2c2a90
Merge: 8cbca7a 978cedb
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Mon Jun 27 17:23:59 2022 +0100

    Merge pull request #42 from sbs-discovery-sweden/dependabot/github_actions/ossf/scorecard-action-1.1.1
    
    build(deps): bump ossf/scorecard-action from 1.0.4 to 1.1.1

commit 978cedb113a53980636bc1a828f0f9d746f1573e
Author: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>
Date:   Mon Jun 27 16:17:09 2022 +0000

    build(deps): bump ossf/scorecard-action from 1.0.4 to 1.1.1
    
    Bumps [ossf/scorecard-action](https://github.com/ossf/scorecard-action) from 1.0.4 to 1.1.1.
    - [Release notes](https://github.com/ossf/scorecard-action/releases)
    - [Changelog](https://github.com/ossf/scorecard-action/blob/main/RELEASE.md)
    - [Commits](https://github.com/ossf/scorecard-action/compare/v1.0.4...v1.1.1)
    
    ---
    updated-dependencies:
    - dependency-name: ossf/scorecard-action
      dependency-type: direct:production
      update-type: version-update:semver-minor
    ...
    
    Signed-off-by: dependabot[bot] <support@github.com>

commit 8cbca7a4575bfcde7c29d5a990a93eb7170cc632
Merge: cb0747f e261bfa
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Mon Jun 27 17:07:38 2022 +0100

    Merge pull request #47 from sbs-discovery-sweden/release-please--branches--master
    
    chore(master): release 1.3.0

commit e261bfab00e9c1afc15d475a8f207edbbecb1fc1
Author: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>
Date:   Mon Jun 27 15:44:14 2022 +0000

    chore(master): release 1.3.0

commit cb0747f1b067220f2e8a81c7c71532d350738f6e
Merge: e26892e a027084
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Mon Jun 27 16:43:40 2022 +0100

    Merge pull request #41 from sbs-discovery-sweden/monitoring/v2
    
    monitoring/v2

commit a027084190fcf8dc26ace436cc32244f1562a521
Author: shipperizer <alexcabb@gmail.com>
Date:   Mon Jun 27 16:23:45 2022 +0100

    fix: drop satori go.uuid for google.uuid

commit ce8bfe5c080686d1b9ff63812d508a3554f6bcee
Author: shipperizer <alexcabb@gmail.com>
Date:   Mon Jun 27 14:40:05 2022 +0100

    build: update deps

commit 61bf614afd23ef877a74735ab1eec61c0f9453d6
Author: shipperizer <alexcabb@gmail.com>
Date:   Mon Jun 27 14:17:42 2022 +0100

    ci: standardize ci release

commit 61358efe0104f1a8a91a43c4250e256db7874301
Author: shipperizer <alexcabb@gmail.com>
Date:   Fri Jun 24 14:43:30 2022 +0100

    ci: adjust vendoring

commit 04d96a74640a3792e4265028f878d2893470b431
Author: shipperizer <alexcabb@gmail.com>
Date:   Fri Jun 24 14:43:22 2022 +0100

    feat: update to monitoring lib v0.4.4

commit 5ecc75f2cc74f5ab91ec5b52de278ce2aaf2abdf
Author: shipperizer <alexcabb@gmail.com>
Date:   Thu Sep 2 15:52:27 2021 +0100

    ci: fix missing release check

commit e26892e4f9db66944dfa84900a969347da6b1594
Merge: b909c9a a6c54ba
Author: Ian Rankin <60607903+irdiscovery@users.noreply.github.com>
Date:   Mon Feb 21 13:18:21 2022 -0500

    Merge pull request #40 from sbs-discovery-sweden/LABS-1574
    
    test: LABS-1574 bump up test coverage

commit a6c54badb439b5f8edc1b1f6bd68a5ae84f5c048
Author: irdiscovery <ian_rankin@discovery.com>
Date:   Mon Feb 21 13:05:44 2022 -0500

    test: LABS-1574 change test error format

commit a2250f1de78351cf399cf221ae9077cc6bde3806
Author: irdiscovery <ian_rankin@discovery.com>
Date:   Fri Feb 18 05:20:47 2022 -0500

    test: LABS-1574 update release script dependency name

commit 8c63ab232071f10c34240dc26f74267ec73c57cb
Author: irdiscovery <ian_rankin@discovery.com>
Date:   Fri Feb 18 05:16:50 2022 -0500

    test: LABS-1574 bump up test coverage

commit b909c9a22a19b43117b3e169a98493372e9bb8d2
Merge: e48678e 9b1998d
Author: Ian Rankin <60607903+irdiscovery@users.noreply.github.com>
Date:   Thu Dec 2 17:40:16 2021 +0000

    Merge pull request #39 from sbs-discovery-sweden/LABS-1290-tee-cmd-fix
    
    ci: LABS-1290 fix blank code coverage file output

commit 9b1998d8652e5cef2315249ed2a2a6b15c9453b7
Author: irdiscovery <ian_rankin@discovery.com>
Date:   Thu Dec 2 17:27:28 2021 +0000

    ci: LABS-1290 fix makefile test commands

commit cc72af1a7675f1d5b96afe2e36566f0419a7a06f
Author: irdiscovery <ian_rankin@discovery.com>
Date:   Thu Dec 2 16:21:03 2021 +0000

    ci: LABS-1290 fix blank code coverage file output

commit e48678e46d0fb0466fe5aea28a16466bcf0491bb
Merge: 6352d21 1ec7d09
Author: Ian Rankin <60607903+irdiscovery@users.noreply.github.com>
Date:   Thu Dec 2 15:04:19 2021 +0000

    Merge pull request #38 from sbs-discovery-sweden/LABS-1290
    
    ci: #LABS-1290 switch sonarcube for sonarcloud

commit 1ec7d098f95bf85531bbadb0fa596f1403566291
Author: irdiscovery <ian_rankin@discovery.com>
Date:   Thu Dec 2 14:51:40 2021 +0000

    ci: LABS-1290 update sonarcube token and config

commit 1c4b11ba12bec9f36263b857ec58920550df3a08
Author: irdiscovery <ian_rankin@discovery.com>
Date:   Thu Nov 25 14:27:09 2021 +0000

    ci: #LABS-1290 exclude test json from test coverage

commit 0e6ebeaee9fc270dcce170b479e82a3d46d5ef56
Author: irdiscovery <ian_rankin@discovery.com>
Date:   Thu Nov 25 14:06:05 2021 +0000

    ci: #LABS-1290 add mock exclusion to test coverage json

commit 7f4994a2c3363c385012609063cb89728a5b1935
Author: irdiscovery <ian_rankin@discovery.com>
Date:   Thu Nov 25 10:22:30 2021 +0000

    ci: #LABS-1290 switch sonarcube for sonarcloud

commit 6352d21ac96c3848ec980879521384668699bee7
Merge: cc7eb1b 2bf3518
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Tue Aug 31 13:15:51 2021 +0100

    Merge pull request #35 from sbs-discovery-sweden/release-v1.2.6
    
    chore: release 1.2.6

commit 2bf35180089b89b452f924af2969d9a71b736b1b
Author: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>
Date:   Tue Aug 31 11:19:05 2021 +0000

    chore: release 1.2.6

commit cc7eb1b60f569e39002f67e4b46d90abdf7fcc5d
Merge: ea4fb5c 00f3eba
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Tue Aug 31 12:18:40 2021 +0100

    Merge pull request #36 from sbs-discovery-sweden/fix/empty/msg
    
    fix/empty/msg

commit 00f3eba91a65167a386921cbf76be1e74091ff22
Author: shipperizer <alexcabb@gmail.com>
Date:   Tue Aug 31 12:11:55 2021 +0100

    fix: move stream event metric call, adjust tests

commit 91ba33f5ce44e3db3689b673e88390aa4631867e
Author: shipperizer <alexcabb@gmail.com>
Date:   Tue Aug 31 12:11:31 2021 +0100

    fix: extra mock call

commit ea4fb5cd3bb74ac869d9d19ace6aa48488e4d946
Merge: 6deb877 29d5dbb
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Tue Aug 31 10:33:16 2021 +0100

    Merge pull request #34 from sbs-discovery-sweden/fix/empty/msg
    
    fix: don't feed empty messages to the consumer flow function if error

commit 29d5dbb7df33bb52cd6d7df4a3e9c20e950f1c99
Author: shipperizer <alexcabb@gmail.com>
Date:   Tue Aug 31 10:25:49 2021 +0100

    fix: don't feed empty messages to the consumer flow function if error

commit 6deb877ef30e74c4a7fa4759619e6a7d01c70b5b
Merge: c6756b1 9e53064
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Thu Aug 12 15:41:20 2021 +0200

    Merge pull request #32 from sbs-discovery-sweden/release-v1.2.5
    
    chore: release 1.2.5

commit 9e530646fdbe3a71cf938a548631e786f3e6dd2f
Author: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>
Date:   Thu Aug 12 13:39:45 2021 +0000

    chore: release 1.2.5

commit c6756b19289ac5d033cc5787b4dca22362081875
Merge: 9a8c548 99375f9
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Thu Aug 12 15:39:08 2021 +0200

    Merge pull request #31 from sbs-discovery-sweden/fix/sm
    
    fix: remove HostnameImmutable from sm config

commit 99375f9507e7135f2c72671e03411358e3e82211
Author: shipperizer <alexcabb@gmail.com>
Date:   Thu Aug 12 15:38:01 2021 +0200

    fix: remove HostnameImmutable from sm config

commit 9a8c548b6d03fa5e1a48e57e413988fff97b4d34
Merge: 182b3fb c368fb8
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Wed Aug 4 18:09:42 2021 +0200

    Merge pull request #30 from sbs-discovery-sweden/release-v1.2.4
    
    chore: release 1.2.4

commit c368fb8f774f3dd937dfaccbc473cd524ba5aee7
Author: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>
Date:   Wed Aug 4 15:54:54 2021 +0000

    chore: release 1.2.4

commit 182b3fbe0ba0a958f1658ef36e29d895fce7fe39
Merge: 7cea3ce 62e9ecc
Author: Xin Jin <xin_jin@discovery.com>
Date:   Wed Aug 4 16:54:29 2021 +0100

    Merge pull request #29 from sbs-discovery-sweden/fix/refresh_issue
    
    fix: use previous encoder when it's not nil

commit 62e9ecc8fb7d42c958028744c3def8d34badee57
Author: xjin <xinjin214@gmail.com>
Date:   Wed Aug 4 16:48:07 2021 +0100

    fix: use previous encoder when it's not nil

commit 565958329f16f5b4fa9cf5521dd701347b3e9201
Author: xjin <xinjin214@gmail.com>
Date:   Wed Aug 4 16:44:29 2021 +0100

    fix: use previous encoder when it's not nil

commit 7cea3ce8840f2f6b51bf4c1716362d69169e466b
Merge: acc954e 41f3ef3
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Fri Jul 30 14:44:23 2021 +0100

    Merge pull request #28 from sbs-discovery-sweden/release-v1.2.3
    
    chore: release 1.2.3

commit 41f3ef3fd791c60282186e3cae28b9bbd31e7c5a
Author: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>
Date:   Fri Jul 30 13:44:08 2021 +0000

    chore: release 1.2.3

commit acc954ed13f00986fe9a89afd5c46b9b375aa9d8
Merge: 4d38d5f 977c1f6
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Fri Jul 30 14:43:42 2021 +0100

    Merge pull request #27 from sbs-discovery-sweden/fix/close
    
    fix: allow a producer to be closed multiple times without panicking

commit 977c1f6d95d24207314575e39091f675b32634d2
Author: shipperizer <alexcabb@gmail.com>
Date:   Fri Jul 30 14:06:14 2021 +0100

    fix: allow a producer to be closed multiple times without panicking

commit 4d38d5fd66cc9483d04353e8234488258afb1d09
Merge: 2149447 1c0e404
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Thu Jul 22 18:44:31 2021 +0100

    Merge pull request #26 from sbs-discovery-sweden/release-v1.2.2
    
    chore: release 1.2.2

commit 1c0e40414f127450abe7a18af0db0a7110992ef6
Author: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>
Date:   Thu Jul 22 17:19:27 2021 +0000

    chore: release 1.2.2

commit 21494477c8b2764d36fc82155be6d01ee3bc9a3e
Merge: 524e0a3 aac0cb7
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Thu Jul 22 18:19:04 2021 +0100

    Merge pull request #25 from sbs-discovery-sweden/lock/autorefresh
    
    lock/autorefresh

commit aac0cb7ad3666bcb066bd213b54fe784c2c223d0
Author: shipperizer <alexcabb@gmail.com>
Date:   Thu Jul 22 17:11:10 2021 +0100

    fix: close reader first, then consumer goroutines

commit 78d8d17575754abc79a6b58d74ba86944dad29a9
Author: shipperizer <alexcabb@gmail.com>
Date:   Thu Jul 22 15:54:44 2021 +0100

    fix: adjust ci triggers

commit 74da928118d2ef262b7d363797ed9dc49d27bd55
Author: shipperizer <alexcabb@gmail.com>
Date:   Thu Jul 22 15:51:01 2021 +0100

    fix: use context WithTimeout to stop the blocking ReadMessage and unlock the mutex on the object for the autorefresh

commit e52a38d024581212626d33db0348027e1744520a
Author: shipperizer <alexcabb@gmail.com>
Date:   Thu Jul 22 15:49:57 2021 +0100

    fix: adjust tests for new config

commit fc78c5a99eec4e02a8d069b9a2be01b2aefe0c41
Author: shipperizer <alexcabb@gmail.com>
Date:   Thu Jul 22 15:48:02 2021 +0100

    fix: add ReadTimeout to ReaderConfig

commit 524e0a3995fcdd082d88cc4d7a54b447d8c8b9d7
Merge: bec68fb db14169
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Wed Jul 21 13:32:42 2021 +0100

    Merge pull request #24 from sbs-discovery-sweden/release-v1.2.1
    
    chore: release 1.2.1

commit db141696a30092aa689b418bd8f7be2b5a89071c
Author: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>
Date:   Wed Jul 21 12:25:08 2021 +0000

    chore: release 1.2.1

commit bec68fb42362d4f0a5ae1cf78d1b44daa8caf58f
Merge: 85adc3b ac90a09
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Wed Jul 21 13:24:35 2021 +0100

    Merge pull request #23 from sbs-discovery-sweden/LABS-709
    
    LABS 709

commit ac90a090dce35ccdf438e79f74f1264907e51dae
Author: shipperizer <alexcabb@gmail.com>
Date:   Wed Jul 21 13:20:08 2021 +0100

    fix: adjust tests in subscriber module

commit 6be358d73b85e27721eff71fafb2add7760a93cd
Author: shipperizer <alexcabb@gmail.com>
Date:   Wed Jul 21 13:19:50 2021 +0100

    fix: adjust message interface to accept interfaces, not bytes

commit 85adc3b1df4f6fdbb80cc76632600fd706620dc5
Merge: a3a559e b34223c
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Tue Jul 20 15:03:29 2021 +0100

    Merge pull request #21 from sbs-discovery-sweden/release-v1.2.0
    
    chore: release 1.2.0

commit b34223cd484d2255fe4171836941a73cdc16a976
Author: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>
Date:   Tue Jul 20 13:13:03 2021 +0000

    chore: release 1.2.0

commit a3a559ec64f75a9c0346f36f633c18456f936260
Merge: 60774a9 3ce571b
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Tue Jul 20 14:12:41 2021 +0100

    Merge pull request #22 from sbs-discovery-sweden/LABS-709
    
    docs: adding comments and func helper

commit 3ce571bf1bac01080c8b0fd91ed076ea0500fa3d
Author: shipperizer <alexcabb@gmail.com>
Date:   Tue Jul 20 13:57:10 2021 +0100

    fix: remove wrong refresh metric

commit 1ba3da82f72a9dc0296a068bd97d457137eae033
Author: shipperizer <alexcabb@gmail.com>
Date:   Tue Jul 20 13:56:48 2021 +0100

    docs: adding comments and func helper

commit 60774a9db2ac90424dee8ba4a53528c0a15d8668
Merge: ac054ae df8705b
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Fri Jul 16 17:40:52 2021 +0100

    Merge pull request #20 from sbs-discovery-sweden/LABS-709
    
    LABS-70(8|9): autorefresh for channeled & standard consumer and publishers

commit df8705bcb7c8bdf11b29ae80264cd7bfe1290a70
Author: shipperizer <alexcabb@gmail.com>
Date:   Wed Jul 14 15:42:26 2021 +0100

    fix: exclude mock_*.go files

commit 7b629012857acf1f563ed2e7e1ef6113f89c9ddf
Author: shipperizer <alexcabb@gmail.com>
Date:   Tue Jul 13 21:09:16 2021 +0100

    feat: modularize everything

commit 5bd438a5cb95f5b74ef8c4b9481e08c4392bf20f
Author: shipperizer <alexcabb@gmail.com>
Date:   Tue Jul 13 14:07:51 2021 +0100

    fix: update deps and ci

commit c46fa80d01b0717fb24603cc11a6ed45e8402294
Author: shipperizer <alexcabb@gmail.com>
Date:   Tue Jul 13 14:06:29 2021 +0100

    fix: minor syntax changes

commit 97064bc732c34d6e77fce2e0b9bf83eafaba37ca
Author: shipperizer <alexcabb@gmail.com>
Date:   Fri Jul 9 16:48:35 2021 +0100

    feat: publisher package

commit 8841d630db12738cdd7edd6825d65442b609434c
Author: shipperizer <alexcabb@gmail.com>
Date:   Fri Jul 9 14:06:53 2021 +0100

    fix: rename konsumer package to subscriber

commit 4f8818ace73afc700d0fd15f441bc37d750240ca
Author: shipperizer <alexcabb@gmail.com>
Date:   Thu Jul 8 12:23:21 2021 +0100

    fix: update deps

commit eb38713b64511fb5781348a52fb2b6d8474cf22e
Author: shipperizer <alexcabb@gmail.com>
Date:   Thu Jul 8 12:22:37 2021 +0100

    feat: use zap logger

commit 485893eb74121902cc9668cd5dc44863c3657fef
Author: shipperizer <alexcabb@gmail.com>
Date:   Fri Jul 2 16:49:52 2021 +0100

    feat: standard consumer code, no channels

commit 7a2c8288294770bdf257b34df6e8ed660517a1df
Author: shipperizer <alexcabb@gmail.com>
Date:   Fri Apr 30 20:17:58 2021 +0100

    fix: adjust ci

commit e57b18cbd1cd8ea7234a10e6d84417df8d18ead0
Author: shipperizer <alexcabb@gmail.com>
Date:   Fri Apr 30 19:25:47 2021 +0100

    fix: update dependencies, mainly labs monitoring lib

commit 8ef59fc4e28a0cdaec85460c910f39f013de6181
Author: shipperizer <alexcabb@gmail.com>
Date:   Fri Apr 30 19:25:25 2021 +0100

    feat: channel consumer implementation

commit 6ef986d338b2ceb0c6ed47887ca214b27fa7b508
Author: shipperizer <alexcabb@gmail.com>
Date:   Fri Apr 30 19:25:05 2021 +0100

    fix: move reader into konsumer module

commit cf38d77c61702e99df40314868bea55a7d341d30
Author: shipperizer <alexcabb@gmail.com>
Date:   Fri Apr 30 19:24:49 2021 +0100

    feat: create autorefresh struct

commit 51d60d0b5ef207ffdb689d809fdeaaa36906abc4
Author: shipperizer <alexcabb@gmail.com>
Date:   Fri Apr 30 19:23:18 2021 +0100

    feat: move tls handling into separate module

commit ac054ae7244cbe4e506461798a6f1eaf38c6e128
Merge: 8d17a6e c4cd1a1
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Fri Apr 23 15:23:26 2021 +0100

    Merge pull request #18 from sbs-discovery-sweden/release-v1.1.2
    
    chore: release 1.1.2

commit c4cd1a1dc180e7e0a90589db1debe6977a5d1e5a
Author: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>
Date:   Fri Apr 23 14:21:41 2021 +0000

    chore: release 1.1.2

commit 8d17a6ed9e7ad687c5976c8dac58025ad04653a0
Merge: cc10392 3c5c139
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Fri Apr 23 15:21:20 2021 +0100

    Merge pull request #19 from sbs-discovery-sweden/LABS-723
    
    fix: LABS-723 test commit

commit 3c5c139e5ea91166666f652f67bbe0396e9d7353
Author: shipperizer <alexcabb@gmail.com>
Date:   Fri Apr 23 15:18:30 2021 +0100

    fix: LABS-723 test commit

commit cc103929081069837be1a84886cf0e6fc1adf4bb
Merge: 9c4b165 759b99d
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Fri Apr 23 15:19:14 2021 +0100

    Merge pull request #17 from sbs-discovery-sweden/jira/release
    
    fix: LABS-723 add extra gotcha to docs

commit 759b99dfa612444fd76f88d9497f95fc99e75f9c
Author: shipperizer <alexcabb@gmail.com>
Date:   Fri Apr 23 15:17:33 2021 +0100

    fix: LABS-723 add extra gotcha to docs

commit 9c4b1657a47aa04af36395743816e4e59f6b74f5
Merge: f5a0fdc 8ff9d92
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Fri Apr 23 15:11:46 2021 +0100

    Merge pull request #16 from sbs-discovery-sweden/release-v1.1.1
    
    chore: release 1.1.1

commit 8ff9d921895e25ba2c5d8c275b6f9960074f8b3b
Author: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>
Date:   Fri Apr 23 14:11:25 2021 +0000

    chore: release 1.1.1

commit f5a0fdc3114c482e5efa33e07142e59ca17c4cb8
Merge: 4ac2369 8f41045
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Fri Apr 23 15:10:57 2021 +0100

    Merge pull request #15 from sbs-discovery-sweden/jira/release
    
    [LABS-723] fix docs

commit 8f41045b0ebc95c633de03de9e39b3d7d72449b3
Author: shipperizer <alexcabb@gmail.com>
Date:   Fri Apr 23 15:10:06 2021 +0100

    fix: LABS-723 fix docs

commit 4ac2369d205f7cf5db80159d5ad9ee3cba198104
Merge: f5c56ea f5c4908
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Fri Apr 23 15:03:38 2021 +0100

    Merge pull request #14 from sbs-discovery-sweden/release-v1.1.0
    
    chore: release 1.1.0

commit f5c4908fbe5045a7fd0e2705cd15bde6063054ba
Author: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>
Date:   Fri Apr 23 14:02:10 2021 +0000

    chore: release 1.1.0

commit f5c56eaf49ad8661367f9271d54ff9286b0b3095
Merge: 74fe098 9b10c3b
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Fri Apr 23 15:01:45 2021 +0100

    Merge pull request #13 from sbs-discovery-sweden/test/release
    
    Test/release

commit 9b10c3b326d7d8efe15090f5e994868f850212bb
Author: shipperizer <alexcabb@gmail.com>
Date:   Fri Apr 23 14:55:01 2021 +0100

    fix: add gotcha about release process

commit 05d419d201923920b73d36a8a346cabc8bc48c7f
Author: shipperizer <alexcabb@gmail.com>
Date:   Fri Apr 23 14:53:51 2021 +0100

    feat: add extras to release process in README

commit dfb1514623bd042872afb7c4dac6da9ba7895519
Author: shipperizer <alexcabb@gmail.com>
Date:   Fri Apr 23 14:51:42 2021 +0100

    fix: remove obsolete docs and generate extra commit for testing release-please-action

commit 04f3c522ee6f17c84bcb2c5f45627417fafceaa2
Author: shipperizer <alexcabb@gmail.com>
Date:   Fri Apr 23 14:51:03 2021 +0100

    feat: update documentation for relase

commit 74fe098a6654e725cc1021bf20cf476e1e64835f
Merge: da73240 df1e1c7
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Fri Apr 23 14:35:39 2021 +0100

    Merge pull request #12 from sbs-discovery-sweden/stale/action
    
    fix: add stale action to close stale issues/prs

commit df1e1c77b1e610a60848c94ce05e53f7193c8e54
Author: shipperizer <alexcabb@gmail.com>
Date:   Fri Apr 23 14:32:32 2021 +0100

    fix: add stale action to close stale issues/prs

commit da732402f9734a9a166eb9236417a66f9ebe7eb7
Merge: e76b825 8e8a9d4
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Fri Apr 23 14:22:43 2021 +0100

    Merge pull request #10 from sbs-discovery-sweden/release/action
    
    feat: add release-please-action for release management

commit 8e8a9d4d7fa60bac056b53c187ff671d74923fb3
Author: shipperizer <alexcabb@gmail.com>
Date:   Fri Apr 23 14:19:44 2021 +0100

    feat: add release-please-action for release management

commit e76b8253b7d791625a235870fe66052576311315
Merge: 3656d84 7fc00bf
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Tue Apr 13 11:01:27 2021 +0100

    Merge pull request #9 from sbs-discovery-sweden/LABS-713
    
    [release] v1.0.0: upgraded to GOv1.16 and AWS SDK v2

commit 7fc00bfb71e53f15a475f45f65dc77ecce3e2806
Author: Athar Khan <iamatharkhan@gmail.com>
Date:   Tue Apr 13 10:53:07 2021 +0100

    [release] removed v1 from module

commit 8b4142e099a22e44da28b0a2fc7dc96d4c8ea848
Author: Athar Khan <iamatharkhan@gmail.com>
Date:   Tue Apr 13 10:42:29 2021 +0100

    [release] reinstated go.sum

commit eb64d326bf3910fd4c1e44274d615678b83b8a89
Author: Athar Khan <iamatharkhan@gmail.com>
Date:   Tue Apr 13 10:34:38 2021 +0100

    [release] upgraded to GOv1.16 and AWS SDK v2

commit 3656d849b77f2d5ba5d6efc6bb54a1ec29964a6b
Merge: 63a6f09 571f8ae
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Mon Feb 8 11:56:09 2021 +0000

    Merge pull request #8 from sbs-discovery-sweden/release
    
    [release] v0.3.16: release management

commit 571f8ae9ce860fe0e5271cb02ec95909a7c35780
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Mon Feb 8 11:51:02 2021 +0000

    [release] v0.3.16: release management

commit 63a6f093a6887a0935416533974d6d72389d2e24
Merge: 404df8e fd647be
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Mon Feb 8 10:44:13 2021 +0000

    Merge pull request #7 from sbs-discovery-sweden/release
    
    Release

commit fd647be5f9fc26b0492559a433902d3d05c59917
Author: shipperizer <alexcabb@gmail.com>
Date:   Mon Feb 8 10:40:43 2021 +0000

    [release] v0.3.16: new release management

commit 681720e73b0f64c19b14217bcf089a8069572ca6
Author: shipperizer <alexcabb@gmail.com>
Date:   Mon Feb 8 10:40:12 2021 +0000

    [fix] adjust release process

commit 404df8eb69f89b06f872bd3472cc350e23e190df
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Mon Jan 11 11:20:38 2021 +0000

    Update CHANGELOG.md

commit 232ac5a5b848aa5dd30fbab9c7cf74874113f0e8
Merge: 3eab135 a83d954
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Mon Jan 11 11:19:12 2021 +0000

    Merge pull request #5 from sbs-discovery-sweden/kafka/0.4
    
    Upgrade to segmentio/kafka-go@v0.4.8

commit a83d9542312b9b154a984f3a28a41acc887a05e0
Author: shipperizer <alexcabb@gmail.com>
Date:   Mon Jan 11 11:03:27 2021 +0000

    [feature] bump version and add changelog

commit cb1a4e65f5d860d60c47a040ca16b1a3ef0c7a75
Author: shipperizer <alexcabb@gmail.com>
Date:   Mon Jan 11 11:03:02 2021 +0000

    [feature] upgrade to kafka-go v0.4.8

commit 3eab13528c16e7b2f4d32ec1278f21f9ad51b64f
Merge: 096273c d5f958c
Author: Ian Rankin <60607903+irdiscovery@users.noreply.github.com>
Date:   Wed Nov 25 12:01:16 2020 +0000

    Merge pull request #4 from sbs-discovery-sweden/get-tls-error
    
    Panic when missing secret

commit d5f958ca22a939633cabf454177ba502d4fc1677
Author: irdiscovery <ian_rankin@discovery.com>
Date:   Wed Nov 25 11:26:00 2020 +0000

    Use sirupsen logging to panic

commit 709d2e4ec300a067574e1169ecf61ed44efeeb96
Author: irdiscovery <ian_rankin@discovery.com>
Date:   Tue Nov 24 21:04:18 2020 +0000

    Remove error from panic argument

commit 2b7a7f4bb7ab1d44b360b247f21d9e998bb0d827
Author: irdiscovery <ian_rankin@discovery.com>
Date:   Tue Nov 24 17:09:46 2020 +0000

    Add change log

commit 99e2b068c7d2184068b21ff05c792e1b854e9f96
Author: irdiscovery <ian_rankin@discovery.com>
Date:   Tue Nov 24 16:57:16 2020 +0000

    Revert SecretManagerConfig SMClient field type

commit dc11720dd322bfba77bbb78ca0451adfa9d6581a
Author: irdiscovery <ian_rankin@discovery.com>
Date:   Tue Nov 24 16:54:40 2020 +0000

    Update version

commit cb6c8040b09533a8f98d50c8228927764a949d5b
Author: irdiscovery <ian_rankin@discovery.com>
Date:   Tue Nov 24 16:49:52 2020 +0000

    Add panic when unable to retrieve secrets

commit a9167bcf84c4ff22391bd60af92d0a114a4685f3
Author: irdiscovery <ian_rankin@discovery.com>
Date:   Tue Nov 24 15:47:30 2020 +0000

    Return error from GetTLS

commit 096273c33c7b3d93dae0e84942330ccaaffcaa45
Merge: a1b6009 5443219
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Thu Nov 12 10:22:50 2020 +0000

    Merge pull request #3 from sbs-discovery-sweden/LABS-350
    
    [LABS-350] Allow to pass region and debug args to SMClient

commit 5443219e5f35f270b3697cc01458a91c5c0333bd
Author: shipperizer <alexcabb@gmail.com>
Date:   Thu Nov 12 09:59:20 2020 +0000

    [docs]

commit 9496611be20c2383d903a822d45adb998e2b34e4
Author: shipperizer <alexcabb@gmail.com>
Date:   Thu Nov 12 09:58:37 2020 +0000

    [feature] pass region and debug args to SM client

commit a1b6009dbcd09822370e17147b688036177f5347
Merge: bba4add a6c0c70
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Mon Nov 2 12:20:03 2020 +0000

    Merge pull request #2 from sbs-discovery-sweden/docs
    
    [DOCS]

commit a6c0c7034b4ca55579859736bfb012a6b9944556
Author: shipperizer <alexcabb@gmail.com>
Date:   Mon Nov 2 11:59:53 2020 +0000

    [docs]

commit a475ef7577de7ed35b0fd04d3a88b4dab71a1930
Author: shipperizer <alexcabb@gmail.com>
Date:   Mon Nov 2 11:59:46 2020 +0000

    [ci] add sonarqube

commit bba4add80c12a5800dcf93ab1b6018a91a067f8f
Merge: 42cfa10 99ff92b
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Mon Oct 19 12:25:11 2020 +0100

    Merge pull request #1 from sbs-discovery-sweden/test
    
    Adding Tests

commit 99ff92bc804b627c8a5ef7c2c594792ecf3c6857
Author: shipperizer <alexcabb@gmail.com>
Date:   Mon Oct 19 12:15:27 2020 +0100

    [fix] add tests to tls and sm machinery

commit b463d9331b8acfcac07a0fb6c023d61f0632fe95
Author: shipperizer <alexcabb@gmail.com>
Date:   Mon Oct 19 12:15:09 2020 +0100

    [ci] add files and change CI pipeline to include SM tests and TLS

commit 42cfa1014546065a3b355e14e09a38d2c650e8d1
Author: Alessandro Cabbia <alexcabb@gmail.com>
Date:   Thu Oct 15 11:14:10 2020 +0100

    Update README.md

commit d8ffa67b804f00a2aabcc9c2334ef28a06be141b
Author: shipperizer <alexcabb@gmail.com>
Date:   Mon Oct 12 12:32:29 2020 +0100

    [feature] new release process

commit e34df004e48fd4b8647688d66366c91459891d94
Author: shipperizer <alexcabb@gmail.com>
Date:   Mon Oct 12 12:32:11 2020 +0100

    [fix] downgrade to kafka-go v0.3.10

commit ddf3a8c5f33178527da6964596dc0dad62496efd
Author: shipperizer <alexcabb@gmail.com>
Date:   Mon Oct 12 12:31:50 2020 +0100

    [fix] use latest changes from events-api

commit a589c695847e602c9eb2de3cc781b55d4293deb1
Author: shipperizer <alexcabb@gmail.com>
Date:   Thu Oct 8 07:59:17 2020 +0100

    [wip] move go files on root folder

commit b23a094677c04b6ff78bce86b73b7324cf60e5d4
Author: shipperizer <alexcabb@gmail.com>
Date:   Wed Oct 7 15:00:23 2020 +0100

    [fix] use config SMClient

commit 1ee11bbf707fa2fba7b03918b9efc64b44800a10
Author: shipperizer <alexcabb@gmail.com>
Date:   Wed Oct 7 14:05:24 2020 +0100

    [docs]

commit df2076df28b6131453dc7a3908a61cd8ca91565e
Author: shipperizer <alexcabb@gmail.com>
Date:   Wed Oct 7 14:05:12 2020 +0100

    [feature] first draft of library

commit e329688608c7e6f1162bcc8b894745a1aaa355c3
Author: shipperizer <alexcabb@gmail.com>
Date:   Wed Oct 7 14:05:01 2020 +0100

    [build] setup makefile and deps

commit 7d5534bf103b552f8217e404b780f551ad3ac854
Author: shipperizer <alexcabb@gmail.com>
Date:   Wed Oct 7 14:04:43 2020 +0100

    [ci] boilerplate
