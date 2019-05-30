<template>
  <div>
    <section class="bg-white">
      <el-row>

        <el-col :span="8" :offset="4">
          <div style="margin-bottom: 100px;">
            <span class="selectTitle"><b>Select NFT Type to Upload</b></span>
          </div>
          <el-select v-model="selectedType" class="nftSelect"
          >
            <el-option
              v-for="op in uploadOptions"
              :value="op.type"
              :label="op.label"
              >
            </el-option>
          </el-select>
        </el-col>


        <el-col :span="6" :offset="1" v-if="selectedType==='dat'">
          <el-row type="flex" align="middle">
            <el-col :span="12">
              <div class="uploadTitle">
                <span>Music NFT</span>
              </div>
            </el-col>
            <el-col :span="12">
              <el-upload
                class="avatar-uploader"
                :show-file-list="false"
                :auto-upload="false"
                action=""
                :on-change="previewMusicAvatar">
                <img v-if="imageUrl" :src="imageUrl" class="avatar" style="width: 150px;height: 150px;">
                <el-row v-else class="el-upload" type="flex" align="middle">
                  <el-col style="height: 100%">
                    <span class="el-icon-plus avatar-uploader-icon" style="margin-top: 42%"></span>
                  </el-col>
                </el-row>
              </el-upload>
            </el-col>
          </el-row>
          <div class="description">
            <span>Name:</span>
            <el-input placeholder="Name" v-model="datName" label="Name"
                      @change="datNameChange"
            ></el-input>
          </div>
          <div class="description">
            <span>Price:</span>
            <el-input placeholder="Price" v-model.number="datPrice" label="Price"
                      @change="datPriceChange"
            ></el-input>
          </div>
          <div class="description">
            <span>Number:</span>
            <el-input placeholder="Number" v-model.number="datNumber" label="Number"
                      @change="datNumberChange"
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

        <el-col :span="6" :offset="1" v-if="selectedType==='avatar'">
          <el-row type="flex" align="middle">
            <el-col>
              <div class="uploadTitle">
                <span>Photo NFT</span>
              </div>
            </el-col>
          </el-row>
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

        <el-col :span="6" :offset="1" v-if="selectedType==='other'">
          <div class="uploadTitle">
            <span>Other NFT</span>
          </div>
          <div class="description">
            <span>Name:</span>
            <el-input placeholder="Name" v-model="otherName" label="Name" @change="otherNameChange"></el-input>
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
            list-type="picture"
            :action="uploadOtherPath"
            :on-success="uploadOtherSuccessHook"
            :data="uploadOtherAdditionalData"
            :auto-upload="false">
            <el-button slot="trigger" size="small" type="primary">select data file</el-button>
            <el-button style="margin-left: 10px;" size="small" type="success" @click="submitOther">upload to server
            </el-button>
            <div class="el-upload__tip" slot="tip">Upload image <b>one at a time</b></div>
          </el-upload>
        </el-col>
      </el-row>
    </section>
    <section class="bg-light">
      <el-row style="margin-bottom: 50px;">
        <!--          <img src="../assets/images/avatar.jpg">-->
        <el-col :span="7" :offset="2">
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
        <!--          <img src="../assets/images/avatar.jpg">-->
        <el-col :span="7" :offset="2">
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
      <el-row class="marketPlaceTitle">
        <el-col :offset="1">
          Market Place Active Listings
        </el-col>
      </el-row>
      <el-row :gutter="20">
        <el-col :span=24 :offset="1">
          <div class="upload-file-container">
            <el-table
              :data="tableData.slice((currentPage-1)*pagesize,currentPage*pagesize)"
              stripe
              style="width: 100%">
              <el-table-column :min-width="60"
                               prop="nftLdefIndex"
                               label="Def Index">
                <template slot-scope="scope">
                  <router-link :to="{name:'Child',params:{nftLdefIndex:scope.row.nftLdefIndex}}"
                          class="buttonText"><a>{{scope.row.nftLdefIndex}}</a></router-link>
                </template>
              </el-table-column>
              <el-table-column :min-width="40"
                               prop="nftType"
                               label="Type">
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

<!--    transaction history-->
      <el-row class="marketPlaceTitle">
        <el-col :offset="1">
          Market Place Transaction History
        </el-col>
      </el-row>
      <el-row :gutter="20">
        <el-col :span=24 :offset="1">
          <div class="upload-file-container">
            <el-table
              :data="mkTxHistoryTableData.slice((mkTxHistoryCurrentPage-1)*mkTxHistoryPagesize,mkTxHistoryCurrentPage*mkTxHistoryPagesize)"
              stripe
              style="width: 100%"
              @cell-click="txInfoHandler">
              <el-table-column :min-width="40"
                               prop="transactionAddress"
                               label="Tx Address">
                <template slot-scope="scope">
                  <a :href="'https://kovan.etherscan.io/tx/'+scope.row.transactionAddress"
                     class="buttonText">{{truncateTxAddress(scope.row.transactionAddress)}}...</a>
                </template>
              </el-table-column>
              <el-table-column :min-width="60"
                               prop="nftLdefIndex"
                               label="Def Index">
              </el-table-column>
              <el-table-column :min-width="60"
                               prop="buyer"
                               label="Buyer">
              </el-table-column>
              <el-table-column :min-width="60"
                               prop="seller"
                               label="Seller">
              </el-table-column>
            </el-table>
            <div style="text-align: center;margin-top: 30px;">
              <el-pagination
                background
                :page-size=mkTxHistoryPagesize
                layout="prev, pager, next"
                :total="mkTxHistoryTotal"
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
        otherName: '',
        avatarShortDesc: '',
        datShortDesc: '',
        otherShortDesc: '',
        avatarLongDesc: '',
        datLongDesc: '',
        otherLongDesc: '',
        datPrice:0,
        datNumber:0,
        fileList: [],
        uploadDatPath: undefined,
        uploadDatAdditionalData: undefined,
        uploadOtherPath: undefined,
        uploadOtherAdditionalData: undefined,
        uploadAvatarPath: undefined,
        uploadAvatarAdditionalData: undefined,
        httpPath: undefined,
        address: undefined,
        totalNFT: undefined,
        tableData: [],
        total: 0,
        pagesize: 10,
        currentPage: 1,
        mkTxHistoryTableData: [],
        mkTxHistoryPagesize: 10,
        mkTxHistoryCurrentPage: 1,
        mkTxHistoryTotal:0,
        nftList: undefined,
        nickName: undefined,
        avatarUrl: undefined,
        username: undefined,
        imageUrl: undefined,
        selectedType: undefined,
        uploadOptions: [
          {
            "type": "dat",
            "label": "Music NFT"
          },
          {
            "type": "avatar",
            "label": "Photo NFT"
          },
          {
            "type": "other",
            "label": "Other NFT"
          }
        ]
      }
    },
    methods: {
      truncateTxAddress: function(txAddress) {
          return txAddress.slice(0,20)
      },
      submitDat: function () {
        this.$refs.uploadDat.submit();
      },
      submitAvatar: function () {
        this.$refs.uploadAvatar.submit();
      },
      submitOther: function () {
        this.$refs.uploadOther.submit();
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
      uploadOtherSuccessHook: function (res, file, fileList) {
        this.totalNFT += 1;
        let el = this.setNFTFromResponse(res);
        this.tableData = [el].concat(this.tableData);
        this.$refs.uploadOther.clearFiles();
      },
      setNFTFromResponse: function (nftData) {
        let el = {};
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
        } else if (nftData.supportedType === "721-05") {
          el.nftType = "Other"
          el.nftPowerIndex = "/"
          el.nftLifeIndex = "/"
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
      getMarketHistoryList: function(address) {
        this.axios.get(`${this.httpPath}/market/transactionHistory/${address}`).then(res => {
          for (let i = res.data.nftPurchaseInfo.length - 1; i >= 0; --i) {
            let purchaseInfo = res.data.nftPurchaseInfo[i];
            // let el = this.setMarketTableFromResponse(purchaseInfo);
            this.mkTxHistoryTableData.push(purchaseInfo);
          }
          this.mkTxHistoryTotal = this.mkTxHistoryTableData.length;
          console.log(res.data.mkTxHistoryTableData);
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
      },
      datNumberChange: function() {
        this.$set(this.uploadDatAdditionalData, 'number', this.datNumber);
      },
      datPriceChange: function() {
        this.$set(this.uploadDatAdditionalData, 'price', this.datPrice);
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
      txInfoHandler: function(row, col ,e) {
        // https://kovan.etherscan.io/tx/
      },
      previewMusicAvatar: function (file, fileList) {
        if (file !== undefined) {
          this.$set(this.uploadDatAdditionalData, 'icon', file.raw);
          this.imageUrl = URL.createObjectURL(file.raw);
        } else {
          this.$delete(this.uploadDatAdditionalData, 'icon');
          this.imageUrl = undefined;
        }
      },
    },
    created: function () {
      // console.log(this.$store.state.account)
      // init variables
      this.username = this.$cookies.get('username');
      this.nickName = this.$cookies.get('nickName');
      this.avatarUrl = this.$cookies.get('avatarUrl');
      let address = this.$cookies.get('account').address;
      console.log("address:", address);
      this.address = address;
      this.httpPath = this.$store.state.config.httpPath;
      this.uploadDatPath = this.httpPath + "/file/dat";
      this.uploadAvatarPath = this.httpPath + "/file/avatar";
      this.uploadOtherPath = this.httpPath + "/file/other";
      this.total = this.tableData.length / this.pagesize * 10;
      this.mkTxHistoryTotal = this.mkTxHistoryTableData.length / this.mkTxHistoryPagesize * 10;

      let uploadBaseObject = {
        address: this.address,
        username: this.username,
      };

      this.uploadDatAdditionalData = Object.assign({}, uploadBaseObject);
      this.uploadAvatarAdditionalData = Object.assign({}, uploadBaseObject);
      this.uploadOtherAdditionalData = Object.assign({}, uploadBaseObject);
      // get total nft
      this.getTotalNFT(address);
      // get nft list of user from market place
      this.getNFTList(address);
      this.getMarketHistoryList(address);

      // set default select item
      this.selectedType = this.uploadOptions[0].type;
    }
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

  .avatar-uploader .el-upload {
    border: 1px dashed #d9d9d9;
    border-radius: 6px;
    cursor: pointer;
    position: relative;
    overflow: hidden;
    width: 150px;
    height: 150px;
    justify-content: center;
  }

  .avatar-uploader .el-upload:hover {
    border-color: #409EFF;
  }

  .selectTitle {
    font-size: 1.5rem;
    margin-bottom: 100px;
  }

  .nftSelect {

  }
</style>
