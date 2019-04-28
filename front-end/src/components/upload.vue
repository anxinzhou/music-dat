<template>
  <div>
    <section class="bg-light">


      <el-row :gutter="20" justify="space-around" align="middle">


        <el-col :span="8" :offset="2">
          <div class="description">
            <span>Name:</span>
            <el-input placeholder="Name" v-model="datName" label="Name"
                      @change="datNameChange"
            ></el-input>
          </div>
          <div class="description">
            <span>Short Description:</span>
            <el-input placeholder="Short Description" type="textarea" v-model="datShortDesc" label="Name"
                      @change="datShortDescChange"
            ></el-input>
          </div>
          <div class="description">
            <span>Long Description:</span>
            <el-input placeholder="Long Description" rows=4 type="textarea" v-model="datLongDesc" label="Name"
                      @change="datLongDescChange"
            ></el-input>
          </div>
          <el-upload
            ref="uploadDat"
            name="file"
            :action="uploadDatPath"
            :data="uploadDatAdditionalData"
            :on-success="uploadDatSuccessHook"
            :auto-upload="false">
            <el-button slot="trigger" size="small" type="primary">select music file</el-button>
            <el-button style="margin-left: 10px;" size="small" type="success" @click="submitDat">upload to server
            </el-button>
            <div class="el-upload__tip" slot="tip">Upload music file <b>one at a time</b></div>
          </el-upload>
        </el-col>

        <el-col :span="8" :offset="2">
          <div class="description">
            <span>Name:</span>
            <el-input placeholder="Name" v-model="avatarName" label="Name" @change="avatarNameChange"></el-input>
          </div>
          <div class="description">
            <span>Short Description:</span>
            <el-input placeholder="Short Description" type="textarea" v-model="avatarShortDesc" label="Name"
                      @change="avatarShortDescChange"
            ></el-input>
          </div>
          <div class="description">
            <span>Long Description:</span>
            <el-input placeholder="Long Description" rows=4 type="textarea" v-model="avatarLongDesc" label="Name"
                      @change="avatarLongDescChange"
            ></el-input>
          </div>
          <el-upload
            ref="uploadAvatar"
            name="file"
            list-type="picture"
            :action="uploadAvatarPath"
            :on-success="uploadAvatarSuccessHook"
            :data="uploadAvatarAdditionalData"
            :auto-upload="false">
            <el-button slot="trigger" size="small" type="primary">select image file</el-button>
            <el-button style="margin-left: 10px;" size="small" type="success" @click="submitAvatar">upload to server
            </el-button>
            <div class="el-upload__tip" slot="tip">Upload image <b>one at a time</b></div>
          </el-upload>
        </el-col>
      </el-row>

    </section>
    <section class="bg-white">

      <el-row style="margin-bottom: 100px;">
        <!--          <img src="../assets/images/avatar.jpg">-->
        <el-col :span="6" :offset="2">
          <el-row>
            <el-col :span="12">
              <b>Wallet Address: </b>
            </el-col>
            <el-col :span="12">
              <span style="font-size: 0.8rem">{{address}}</span>
            </el-col>
          </el-row>
        </el-col>
        <el-col :span="6" :offset="6">
          <el-row>
            <el-col :span="12">
              <b>Total NFT: </b>
            </el-col>
            <el-col :span="12" class="text-center">
              {{totalNFT}}
            </el-col>
          </el-row>
        </el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :span=24>
          <div class="upload-file-container">
            <el-table
              :data="tableData.slice((currentPage-1)*pagesize,currentPage*pagesize)"
              stripe
              style="width: 100%">
              <el-table-column :min-width="40"
                               prop="nftType"
                               label="Type">
              </el-table-column>
              <el-table-column :min-width="60"
                               prop="nftLdefIndex"
                               label="Def Index">
              </el-table-column>
              <el-table-column :min-width="60"
                               prop="nftName"
                               label="Name">
              </el-table-column>
              <el-table-column
                prop="activeTicker" :min-width="40"
                label="Active Ticker">
              </el-table-column>
              <el-table-column :min-width="25"
                               prop="nftValue"
                               label="Price">
              </el-table-column>
              <el-table-column :min-width="20"
                               prop="qty"
                               label="Qty">
              </el-table-column>
              <el-table-column :min-width="25"
                               prop="nftPowerIndex"
                               label="Power">
              </el-table-column>
              <el-table-column :min-width="25"
                               prop="nftLifeIndex"
                               label="Life">
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
  </div>
</template>
<script>
  export default {
    props: {},
    data() {
      return {
        avatarName: '',
        datName: '',
        avatarShortDesc: '',
        datShortDesc: '',
        avatarLongDesc: '',
        datLongDesc: '',
        fileList: [],
        uploadDatPath: undefined,
        uploadDatAdditionalData: undefined,
        uploadAvatarPath: undefined,
        uploadAvatarAdditionalData: undefined,
        httpPath: undefined,
        address: undefined,
        totalNFT: undefined,
        tableData: [],
        total: 0,
        pagesize: 10,
        currentPage: 1,
        nftList: undefined,
      }
    },
    methods: {
      submitDat: function () {
        this.$nextTick(() => {
          this.$refs.uploadDat.submit();
        })
      },
      submitAvatar: function () {
        this.$nextTick(() => {
          this.$refs.uploadAvatar.submit();
        });
      },
      uploadDatSuccessHook: function (res, file, fileList) {
        this.totalNFT += 1;
        let el = this.setNFTFromResponse(res);
        this.tableData = [el].concat(this.tableData);
        this.$refs.uploadDat.clearFiles();
      },
      uploadAvatarSuccessHook: function (res, file, fileList) {
        this.totalNFT += 1;
        let el = this.setNFTFromResponse(res);
        this.tableData = [el].concat(this.tableData);
        this.$refs.uploadAvatar.clearFiles();
      },
      setNFTFromResponse: function (nftData) {
        let el = {}
        el.nftLdefIndex = nftData.nftLdefIndex;
        el.nftName = nftData.nftName;
        el.nftCharacId = nftData.nftCharacId;
        el.activeTicker = nftData.activeTicker;
        el.nftValue = nftData.nftValue;
        el.qty = nftData.qty;
        el.shortDesc = nftData.shortDesc;
        el.longDesc = nftData.longDesc;
        if (nftData.supportedType === '721-04') {
          el.nftType = "Dat"
          el.nftPowerIndex = "/"
          el.nftLifeIndex = "/"
        } else if (nftData.supportedType === '721-02') {
          el.nftType = "Avatar"
          el.nftPowerIndex = nftData.nftPowerIndex;
          el.nftLifeIndex = nftData.nftLifeIndex;
        }
        return el;
      },
      getTotalNFT: function (address) {
        // get total nft balance
        this.axios.get(`${this.httpPath}/balance/${address}`).then(res => {
          this.totalNFT = res.data.count;
        }).catch(console.log);
      },
      getNFTList: function (address) {
        this.axios.get(`${this.httpPath}/nftList/${address}`).then(res => {
          for (let i = res.data.nftTranData.length - 1; i >= 0; --i) {
            let nftData = res.data.nftTranData[i];
            let el = this.setNFTFromResponse(nftData);
            this.tableData.push(el);
          }
          this.total = this.tableData.length;
          console.log(res.data.nftTranData);
        }).catch(console.log);
      },
      current_change: function (currentPage) {
        this.currentPage = currentPage
      },
      avatarNameChange: function () {
        this.$set(this.uploadAvatarAdditionalData, 'nftName', this.avatarName);
      },
      avatarShortDescChange: function () {
        this.$set(this.uploadAvatarAdditionalData, 'shortDesc', this.avatarShortDesc);
      },
      avatarLongDescChange: function () {
        this.$set(this.uploadAvatarAdditionalData, 'longDesc', this.avatarLongDesc);
      },
      datNameChange: function () {
        this.$set(this.uploadDatAdditionalData, 'nftName', this.datName);
      },
      datShortDescChange: function () {
        this.$set(this.uploadDatAdditionalData, 'shortDesc', this.datShortDesc);
      },
      datLongDescChange: function () {
        this.$set(this.uploadDatAdditionalData, 'longDesc', this.datLongDesc);
      }
    },
    created: function () {
      // console.log(this.$store.state.account)
      // init variables
      let address = this.$cookies.get('account').address;
      console.log("address:", address);
      this.address = address;
      this.httpPath = this.$store.state.config.httpPath;
      this.uploadDatPath = this.httpPath + "/file/dat";
      this.uploadAvatarPath = this.httpPath + "/file/avatar";
      this.total = this.tableData.length / this.pagesize * 10;
      this.uploadDatAdditionalData = {
        address: this.address,
      };
      this.uploadAvatarAdditionalData = {
        address: this.address,
      };
      // get total nft
      this.getTotalNFT(address);
      // get nft list of user from market place
      this.getNFTList(address);
    }
  }
</script>
<style scoped>
  .description {
    margin-bottom: 20px;
  }
</style>
