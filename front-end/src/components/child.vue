<script src="../router.js"></script>
<template>
  <div>
    <body id="page-top">
    <nav class="navbar navbar-expand-lg navbar-dark fixed-top navbar-shrink" id="mainNav">
      <div class="container">
        <a class="navbar-brand  js-scroll-trigger" href="#" v-scroll-to="'#page-top'">Music Hotpot</a>
        <p id="user-addr" class="js-scroll-trigger user-addr" href="#" v-scroll-to="'#page-top'"></p>
        <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarResponsive"
                aria-controls="navbarResponsive" aria-expanded="false" aria-label="Toggle navigation">
          Menu
          <i class="fas fa-bars"></i>
        </button>
        <div class="collapse navbar-collapse" id="navbarResponsive">
          <ul class="navbar-nav  ml-auto">
            <li class="nav-item text-uppercase">
              <router-link to="/"><a class="nav-link">Home</a></router-link>
            </li>
            <li class="nav-item text-uppercase">
              <a class="nav-link" href="#" v-scroll-to="'#upload'">Child</a>
            </li>
            <li class="nav-item text-uppercase">
              <router-link :to="{name:'Profile'}"><a><img style="width: 50px;height: 50px;" :src="avatarUrl"/></a></router-link>
            </li>
            <li class="nav-item text-uppercase">
              <el-button round style="margin-left:20px;margin-top:2px;height: 45px; width: 75px;" @click="logout"><span style="margin-left:-10px;">Logout</span></el-button>
            </li>
          </ul>
        </div>
      </div>
    </nav>
    <header class="masthead">
      <div class="container">
        <div class="intro-text">
          <div class="intro-lead-in"></div>
        </div>
      </div>
    </header>

      <a id="upload">
        <section class="bg-light">


          <el-row :gutter="20" justify="space-around" align="middle">


            <el-col :span="14" :offset="5">
              <div class="uploadTitle">
                <span>Children NFT</span>
              </div>
              <div class="description">
                <span>Name:</span>
                <el-input placeholder="Name" v-model="otherName" label="Name"
                          @change="otherNameChange"
                ></el-input>
              </div>
              <div class="description">
                <span>Short Description:</span>
                <el-input placeholder="Short Description" type="textarea" v-model="otherShortDesc" label="Name"
                          @change="otherShortDescChange"
                ></el-input>
              </div>
              <div class="description">
                <span>Long Description:</span>
                <el-input placeholder="Long Description" rows=4 type="textarea" v-model="otherLongDesc" label="Name"
                          @change="otherLongDescChange"
                ></el-input>
              </div>
              <el-upload
                ref="uploadOther"
                name="file"
                :action="uploadOtherPath"
                :data="uploadOtherAdditionalData"
                :on-success="uploadOtherSuccessHook"
                :on-error="uploadOtherErrorHook"
                :auto-upload="false">
                <el-button slot="trigger" size="small" type="primary">select data file</el-button>
                <el-button style="margin-left: 10px;" size="small" type="success" @click="submitOther">upload to server
                </el-button>
                <div class="el-upload__tip" slot="tip">Upload music file <b>one at a time</b></div>
              </el-upload>
            </el-col>

          </el-row>

        </section>
      </a>

    <section class="bg-white">
      <el-row style="margin-bottom: 50px;">
        <!--          <img src="../assets/images/avatar.jpg">-->
        <el-col :span="6" :offset="2">
          <el-row type="flex" align="middle">
            <el-col :span="12">
              <b>User Name: </b>
            </el-col>
            <el-col :span="6">
              <span style="font-size: 0.8rem">{{username}}</span>
            </el-col>
            <el-col :span="2" :offset="4">
              <img :src="avatarUrl" style="width: 100px;"/>
            </el-col>
          </el-row>
        </el-col>
      </el-row>
      <el-row style="margin-bottom: 100px;">
        <el-col :span="6" :offset="2">
          <el-row>
            <el-col :span="12">
              <b>NftLdefIndex: </b>
            </el-col>
            <el-col :span="12">
              <span style="font-size: 0.8rem">{{parent}}</span>
            </el-col>
          </el-row>
        </el-col>
        <el-col :span="6" :offset="6">
          <el-row>
            <el-col :span="12">
              <b>Total Child NFT: </b>
            </el-col>
            <el-col :span="12" class="text-center">
              {{totalNFT}}
            </el-col>
          </el-row>
        </el-col>
      </el-row>
      <el-row class="marketPlaceTitle">
        <el-col :offset="1">
          Children NFT
        </el-col>
      </el-row>
      <el-row :gutter="20">
        <el-col :span=24 :offset="1">
          <div class="upload-file-container">
            <el-table
              :data="tableData.slice((currentPage-1)*pagesize,currentPage*pagesize)"
              stripe
              style="width: 100%"
            >
              <el-table-column :min-width="60"
                               prop="nftLdefIndex"
                               label="Def Index">
              </el-table-column>
              <el-table-column :min-width="60"
                               prop="nftName"
                               label="Name">
              </el-table-column>
              <el-table-column :min-width="25"
                               prop="price"
                               label="Price">
              </el-table-column>
              <el-table-column :min-width="50"
                               prop="shortDesc"
                               label="Short desc">
              </el-table-column>
              <el-table-column
                prop="longDesc"
                label="Long desc">
              </el-table-column>
            </el-table>
            <div style="text-align: center;margin-top: 30px;">
              <el-pagination
                background
                :page-size=pagesize
                layout="prev, pager, next"
                :total="total"
                @current-change="current_change">
              </el-pagination>
            </div>

          </div>
        </el-col>
      </el-row>
    </section>
    </body>
  </div>
</template>

<script>
  export default {
    name: "child",
    data() {
      return {
        otherName: '',
        otherShortDesc: '',
        otherLongDesc: '',
        uploadOtherPath: undefined,
        uploadOtherAdditionalData: undefined,
        httpPath: undefined,
        parent: '',
        totalNFT: undefined,
        tableData: [],
        total: 0,
        pagesize: 10,
        currentPage: 1,
        nickName: undefined,
        avatarUrl: undefined,
        username: undefined,
      }
    },
    methods: {
      submitOther: function () {
        this.$refs.uploadOther.submit();
      },
      uploadOtherSuccessHook: function (res, file, fileList) {
        this.totalNFT += 1;
        this.tableData = [res].concat(this.tableData);
        this.$refs.uploadOther.clearFiles();
        this.total +=1;
      },
      uploadOtherErrorHook: function(err, file, fileList) {
        this.$store.state.notifyError("Fail to Upload photo NFT");
        console.log(err)
      },
      getTotalChildNFT: function (parent) {
        this.axios.get(`${this.httpPath}/nfts/${parent}/balance`).then(res => {
          this.totalNFT = res.data.count;
        }).catch(console.log);
      },
      getChildNFTList: function (parent) {
        this.axios.get(`${this.httpPath}/nfts/${parent}/children`).then(res => {
          this.tableData = res.data.nftTranData;
          console.log(res.data);
          this.total = this.tableData.length;
          if (this.total===0) {
            this.currentPage = 1;
          } else {
            this.currentPage = Math.floor((this.total-1)/this.pagesize) +1;
          }
        }).catch(console.log);
      },
      current_change: function (currentPage) {
        this.currentPage = currentPage
      },
      otherNameChange: function () {
        this.$set(this.uploadOtherAdditionalData, 'nftName', this.otherName);
      },
      otherShortDescChange: function () {
        this.$set(this.uploadOtherAdditionalData, 'shortDesc', this.otherShortDesc);
      },
      otherLongDescChange: function () {
        this.$set(this.uploadOtherAdditionalData, 'longDesc', this.otherLongDesc);
      },
      getNickname: function(uuid) {
        let httpPath = this.$store.state.config.httpPath;
        this.axios.get(`${httpPath}/profile/${uuid}/nickname`).then(res=>{
          this.nickname = res.data.nickname;
        }).catch(err=>{
          console.log(err.response.data.reason)
        })
      },
      getAvatarUrl: function(uuid) {
        let httpPath = this.$store.state.config.httpPath;
        this.axios.get(`${httpPath}/profile/${uuid}/avatar`).then(res=>{
          this.avatarUrl = res.data.avatarUrl;
        }).catch(err=>{
          console.log(err.response.data.reason)
        })
      },
      logout: function () {
        console.log("logout")
        this.$cookies.remove("access-token");
        this.$cookies.remove("uuid");
        this.$router.replace('/login');
      },
    },
    created: function () {
      let uuid = this.$cookies.get('uuid');
      this.httpPath = this.$store.state.config.httpPath;
      this.uploadOtherPath = this.httpPath + "/file/other";
      this.parent = this.$route.params.nftLdefIndex;
      console.log("this nft index",this.parent);
      this.uploadOtherAdditionalData = {
        uuid: uuid,
        parent: this.parent,
      }
      this.getTotalChildNFT(this.parent);
      this.getChildNFTList(this.parent);
      this.getNickname(uuid);
      this.getAvatarUrl(uuid);
    },
    mounted: function () {
      // var el = document.getElementById('mainNav')
      // el.classList.add('navbar-shrink');
      // function scrollHandle() {
      //   if (document.documentElement.scrollTop > 100) {
      //     el.classList.add('navbar-shrink')
      //   } else {
      //     el.classList.remove('navbar-shrink')
      //   }
      // }
      //
      // window.addEventListener('scroll', scrollHandle)
    },
  }
</script>

<style scoped>
  .description {
    margin-bottom: 20px;
  }

  .uploadTitle {
    font-size: 1.5rem;
    margin-bottom: 50px;
  }

  .marketPlaceTitle {
    font-size: 1.5rem;
    margin-bottom: 50px;
  }
</style>

