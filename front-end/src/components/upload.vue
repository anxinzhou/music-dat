<script src="../store.js"></script>
<script src="../router.js"></script>
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
              :key="op.type"
            >
            </el-option>
          </el-select>
          <div class="airDropTitle">
            <span><b>Allow Airdrop</b></span>
          </div>
          <el-radio-group v-model="allowAirdrop" @input="airDropOptionChange">
            <el-radio :label=true>Yes</el-radio>
            <el-radio :label=false>No</el-radio>
          </el-radio-group>
        </el-col>


        <el-col :span="10" :offset="1" v-if="selectedType==='dat'">
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
          <el-row type="flex">
            <el-col :span="16">
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
                :on-error="uploadDatErrorHook"
                :auto-upload="false">
                <el-button slot="trigger" size="small" type="primary">select music file</el-button>
                <el-button style="margin-left: 10px;" size="small" type="success" @click="submitDat">upload to server
                </el-button>
                <div class="el-upload__tip" slot="tip">Upload music file <b>one at a time</b></div>
              </el-upload>
            </el-col>
            <el-col :span="12" :offset="2">
              <div class="description">
                <span>Creator:</span>
                <el-input placeholder="Creator Percent" v-model.number="creatorPercent" label="Number"
                          @change="creatorPercentChange"
                ></el-input>
              </div>
              <div class="description">
                <span>Lyrics Writer:</span>
                <el-input placeholder="Lyrics Writer Percent" v-model.number="lyricsWriterPercent" label="Number"
                          @change="lyricsWriterPercentChange"
                ></el-input>
              </div>
              <div class="description">
                <span>Song Composer</span>
                <el-input placeholder="Song Composer Percent" v-model.number="songComposerPercent" label="Number"
                          @change="songComposerPercentChange"
                ></el-input>
              </div>
              <div class="description">
                <span>Publisher:</span>
                <el-input placeholder="Publisher Percent" v-model.number="publisherPercent" label="Number"
                          @change="publisherPercentChange"
                ></el-input>
              </div>
              <div class="description">
                <span>User:</span>
                <el-input placeholder="User Percent" v-model.number="userPercent" label="Number"
                          @change="userPercentChange"
                ></el-input>
              </div>
            </el-col>
          </el-row>
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
            :on-error="uploadAvatarErrorHook"
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
            :on-error="uploadOtherErrorHook"
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
          <el-row>
            <el-col :span="12">
              <b>User: </b>
            </el-col>
            <el-col :span="6">
              <span style="font-size: 0.8rem">{{nickname}}</span>
            </el-col>
            <el-col :span="2" :offset="4">
              <img :src="avatarUrl" style="margin-top: -20px;width: 100px;"/>
            </el-col>
          </el-row>
        </el-col>
        <el-col :span="7" :offset="3">
          <el-row type="flex" align="middle">
            <el-col :span="12">
              <b>Intro: </b>
            </el-col>
            <el-col :span="12" >
              <span style="font-size: 0.8rem">{{intro}}</span>
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
        <el-col :span="7" :offset="3">
          <el-row>
            <el-col :span="12">
              <b>Total NFT: </b>
            </el-col>
            <el-col :span="12" >
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
      <el-row class="marketPlaceListTitle">
        <el-col :offset="1">
          Music List
        </el-col>
      </el-row>
      <el-row :gutter="20">
        <el-col :span=24 :offset="1">
          <div class="upload-file-container">
            <el-table
              :data="datTableData.slice((datCurrentPage-1)*datPagesize,datCurrentPage*datPagesize)"
              stripe
              style="width: 100%">
              <el-table-column :min-width="40"
                               prop="nftLdefIndex"
                               label="Def Index">
                <template slot-scope="scope">
                  <router-link :to="{name:'Child',params:{nftLdefIndex:scope.row.nftLdefIndex}}"
                               class="buttonText"><a>{{scope.row.nftLdefIndex}}</a></router-link>
                </template>
              </el-table-column>
              <el-table-column :min-width="25"
                               prop="nftName"
                               label="Name">
              </el-table-column>
              <el-table-column :min-width="20"
                               prop="price"
                               label="Price">
              </el-table-column>
              <el-table-column :min-width="20"
                               prop="qty"
                               label="Qty">
              </el-table-column>
              <el-table-column :min-width="40"
                               prop="shortDesc"
                               label="Short desc">
              </el-table-column>
              <el-table-column :min-width="40"
                               prop="longDesc"
                               label="Long desc">
              </el-table-column>
              <el-table-column :min-width="25"
                prop="creatorPercent"
                label="Creator">
              </el-table-column>
              <el-table-column :min-width="25"
                prop="lyricsWriterPercent"
                label="Lyrics Writer">
              </el-table-column>
              <el-table-column :min-width="25"
                prop="songComposerPercent"
                label="Song Composer">
              </el-table-column>
              <el-table-column :min-width="25"
                               prop="publisherPercent"
                               label="Publisher">
              </el-table-column>
              <el-table-column :min-width="45"
                               prop="userPercent"
                               label="User">
              </el-table-column>
            </el-table>
            <div style="text-align: center;margin-top: 30px;">
              <el-pagination
                background
                :page-size=datPagesize
                layout="prev, pager, next"
                :total="datTotal"
                @current-change="datCurrentChange">
              </el-pagination>
            </div>

          </div>
        </el-col>
      </el-row>
<!--      avatar list-->
      <el-row class="marketPlaceListTitle">
        <el-col :offset="1">
          Photo NFT List
        </el-col>
      </el-row>
      <el-row :gutter="20">
        <el-col :span=24 :offset="1">
          <div class="upload-file-container">
            <el-table
              :data="avatarTableData.slice((avatarCurrentPage-1)*avatarPagesize,avatarCurrentPage*avatarPagesize)"
              stripe
              style="width: 100%">
              <el-table-column :min-width="40"
                               prop="nftLdefIndex"
                               label="Def Index">
                <template slot-scope="scope">
                  <router-link :to="{name:'Child',params:{nftLdefIndex:scope.row.nftLdefIndex}}"
                               class="buttonText"><a>{{scope.row.nftLdefIndex}}</a></router-link>
                </template>
              </el-table-column>
              <el-table-column :min-width="25"
                               prop="nftName"
                               label="Name">
              </el-table-column>
              <el-table-column :min-width="20"
                               prop="price"
                               label="Price">
              </el-table-column>
              <el-table-column :min-width="25"
                               prop="nftPowerIndex"
                               label="Power">
              </el-table-column>
              <el-table-column :min-width="20"
                               prop="nftLifeIndex"
                               label="Life">
              </el-table-column>
              <el-table-column :min-width="40"
                               prop="shortDesc"
                               label="Short desc">
              </el-table-column>
              <el-table-column :min-width="40"
                               prop="longDesc"
                               label="Long desc">
              </el-table-column>
            </el-table>
            <div style="text-align: center;margin-top: 30px;">
              <el-pagination
                background
                :page-size=avatarPagesize
                layout="prev, pager, next"
                :total="avatarTotal"
                @current-change="avatarCurrentChange">
              </el-pagination>
            </div>

          </div>
        </el-col>
      </el-row>
<!--      other list-->
      <el-row class="marketPlaceListTitle">
        <el-col :offset="1">
          Other NFT List
        </el-col>
      </el-row>
      <el-row :gutter="20">
        <el-col :span=24 :offset="1">
          <div class="upload-file-container">
            <el-table
              :data="otherTableData.slice((otherCurrentPage-1)*otherPagesize,otherCurrentPage*otherPagesize)"
              stripe
              style="width: 100%">
              <el-table-column :min-width="40"
                               prop="nftLdefIndex"
                               label="Def Index">
                <template slot-scope="scope">
                  <router-link :to="{name:'Child',params:{nftLdefIndex:scope.row.nftLdefIndex}}"
                               class="buttonText"><a>{{scope.row.nftLdefIndex}}</a></router-link>
                </template>
              </el-table-column>
              <el-table-column :min-width="25"
                               prop="nftName"
                               label="Name">
              </el-table-column>
              <el-table-column :min-width="20"
                               prop="price"
                               label="Price">
              </el-table-column>
              <el-table-column :min-width="40"
                               prop="shortDesc"
                               label="Short desc">
              </el-table-column>
              <el-table-column :min-width="40"
                               prop="longDesc"
                               label="Long desc">
              </el-table-column>
            </el-table>
            <div style="text-align: center;margin-top: 30px;">
              <el-pagination
                background
                :page-size=otherPagesize
                layout="prev, pager, next"
                :total="otherTotal"
                @current-change="otherCurrentChange">
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
                  <!--                  <a class="buttonText"> :href="'https://kovan.etherscan.io/tx/'+scope.row.transactionAddress"-->
                  {{truncateTxAddress(scope.row.transactionAddress)}}...
                  <!--                  </a>-->
                </template>
              </el-table-column>
              <el-table-column :min-width="60"
                               prop="nftLdefIndex"
                               label="Def Index">
              </el-table-column>
              <el-table-column :min-width="60"
                               prop="buyerNickname"
                               label="Buyer">
              </el-table-column>
              <el-table-column :min-width="60"
                               prop="sellerNickname"
                               label="Seller">
              </el-table-column>
              <el-table-column :min-width="60"
                               prop="timestamp"
                               label="Date">
              </el-table-column>
            </el-table>
            <div style="text-align: center;margin-top: 30px;">
              <el-pagination
                background
                :page-size=mkTxHistoryPagesize
                layout="prev, pager, next"
                :total="mkTxHistoryTotal"
                @current-change="mkTxCurrentChange">
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
        datPrice: 0,
        datNumber: 0,
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
        datTableData: [],
        datTotal: 0,
        datPagesize: 10,
        datCurrentPage: 1,
        avatarTableData: [],
        avatarTotal: 0,
        avatarPagesize: 10,
        avatarCurrentPage: 1,
        otherTableData: [],
        otherTotal: 0,
        otherPagesize: 10,
        otherCurrentPage: 1,
        mkTxHistoryTableData: [],
        mkTxHistoryPagesize: 10,
        mkTxHistoryCurrentPage: 1,
        mkTxHistoryTotal: 0,
        nftList: undefined,
        nickname: undefined,
        avatarUrl: undefined,
        imageUrl: undefined,
        selectedType: undefined,
        allowAirdrop: true,
        intro: '',
        creatorPercent: undefined,
        lyricsWriterPercent: undefined,
        songComposerPercent: undefined,
        publisherPercent: undefined,
        userPercent: undefined,
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
      truncateTxAddress: function (txAddress) {
        return txAddress.slice(0, 20)
      },
      submitDat: function () {
        if(this.creatorPercent + this.lyricsWriterPercent + this.songComposerPercent + this.publisherPercent + this.userPercent !==100) {
          this.$store.state.notifyError("The sum of interest percent should be 100");
          return
        }
        this.$refs.uploadDat.submit();
      },
      submitAvatar: function () {
        this.$refs.uploadAvatar.submit();
      },
      submitOther: function () {
        this.$refs.uploadOther.submit();
      },
      uploadDatSuccessHook: function (res, file, fileList) {
        this.$store.state.notifySuccess("Upload music NFT success");
        this.datTableData = [res].concat(this.datTableData);
        this.$refs.uploadDat.clearFiles();
        this.totalNFT += 1;
        this.datTotal +=1;
      },
      uploadDatErrorHook: function(err,file,fileList) {
        this.$store.state.notifyError("Fail to Upload music NFT");
        console.log(err)
      },
      uploadAvatarSuccessHook: function (res, file, fileList) {
        this.$store.state.notifySuccess("Upload photo NFT success");
        this.totalNFT += 1;
        this.avatarTableData = [res].concat(this.avatarTableData);
        this.$refs.uploadAvatar.clearFiles();
        this.avatarTotal+=1
      },
      uploadAvatarErrorHook: function(err,file,fileList) {
        this.$store.state.notifyError("Fail to Upload photo NFT");
        console.log(err)
      },
      uploadOtherSuccessHook: function (res, file, fileList) {
        this.$store.state.notifySuccess("Upload other NFT success");
        console.log(res);
        this.totalNFT += 1;
        this.otherTableData = [res].concat(this.otherTableData);
        this.$refs.uploadOther.clearFiles();
        this.otherTotal+=1;
      },
      uploadOtherErrorHook: function(err,file,fileList) {
        this.$store.state.notifyError("Fail to Upload other NFT");
        console.log(err)
      },
      datCurrentChange: function (currentPage) {
        this.datCurrentPage = currentPage
      },
      avatarCurrentChange: function (currentPage) {
        this.avatarCurrentPage = currentPage
      },
      otherCurrentChange: function(currentPage) {
        this.otherCurrentPage = currentPage
      },
      mkTxCurrentChange: function (currentPage) {
        this.mkTxHistoryCurrentPage = currentPage
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
      datNumberChange: function () {
        this.$set(this.uploadDatAdditionalData, 'number', this.datNumber);
      },
      datPriceChange: function () {
        this.$set(this.uploadDatAdditionalData, 'price', this.datPrice);
      },
      creatorPercentChange: function() {
        this.$set(this.uploadDatAdditionalData, 'creatorPercent', this.creatorPercent);
      },
      lyricsWriterPercentChange: function() {
        this.$set(this.uploadDatAdditionalData, 'lyricsWriterPercent', this.lyricsWriterPercent);
      },
      songComposerPercentChange: function() {
        this.$set(this.uploadDatAdditionalData, 'songComposerPercent', this.songComposerPercent);
      },
      publisherPercentChange: function() {
        this.$set(this.uploadDatAdditionalData, 'publisherPercent', this.publisherPercent);
      },
      userPercentChange: function() {
        this.$set(this.uploadDatAdditionalData, 'userPercent', this.userPercent);
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
      txInfoHandler: function (row, col, e) {
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
      airDropOptionChange: function (value) {
        console.log("allow airdrop?", value);
        this.allowAirdrop = value;
        this.$set(this.uploadAvatarAdditionalData, 'allowAirdrop', value);
        this.$set(this.uploadDatAdditionalData, 'allowAirdrop', value);
        this.$set(this.uploadOtherAdditionalData, 'allowAirdrop', value);
      },
      getIntro: function(uuid) {
        let httpPath = this.$store.state.config.httpPath;
        this.axios.get(`${httpPath}/profile/${uuid}/intro`).then(res=>{
          this.intro = res.data.intro;
        }).catch(err=>{
          console.log(err.response.data.reason)
        })
      },
      getNickname: function(uuid) {
        let httpPath = this.$store.state.config.httpPath;
        this.axios.get(`${httpPath}/profile/${uuid}/nickname`).then(res=>{
          this.nickname = res.data.nickname;
        }).catch(err=>{
          console.log(err.response.data.reason)
        })
      },
      getWallet: function(uuid) {
        let httpPath = this.$store.state.config.httpPath;
        this.axios.get(`${httpPath}/profile/${uuid}/wallet`).then(res=>{
          this.address = res.data.wallet;
          this.totalNFT = res.data.count;
        }).catch(err=>{
          console.log(err.response.data.reason)
        });
      },
      getAvatarUrl: function(uuid) {
        let httpPath = this.$store.state.config.httpPath;
        this.axios.get(`${httpPath}/profile/${uuid}/avatar`).then(res=>{
          this.avatarUrl = res.data.avatarUrl;
        }).catch(err=>{
          console.log(err.response.data.reason)
        })
      },
      getNFTList: function (uuid) {
        this.getDatList(uuid);
        this.getAvatarList(uuid);
        this.getOtherList(uuid);
      },
      getDatList: function(uuid) {
        let httpPath = this.$store.state.config.httpPath;
        this.axios.get(`${httpPath}/nftList/dat/${uuid}`).then(res => {
           this.datTableData = res.data.nftTranData
           this.datTotal = this.datTableData.length;
          if(this.datTotal===0) {
             this.datCurrentPage = 1;
          } else {
              this.datCurrentPage = Math.floor((this.datTotal-1)/ this.datPagesize)+1;
          }
        }).catch(err=>{
          console.log(err.response.data.reason)
        })
      },
      getAvatarList: function(uuid) {
        let httpPath = this.$store.state.config.httpPath;
        this.axios.get(`${httpPath}/nftList/avatar/${uuid}`).then(res => {
          this.avatarTableData = res.data.nftTranData
          this.avatarTotal = this.avatarTableData.length;
          if(this.datTotal===0) {
            this.avatarCurrentPage = 1;
          } else {
            this.avatarCurrentPage = Math.floor((this.avatarTotal-1)/ this.avatarPagesize)+1;
          }
        }).catch(err=>{
          console.log(err.response.data.reason)
        })
      },
      getOtherList: function(uuid) {
        let httpPath = this.$store.state.config.httpPath;
        this.axios.get(`${httpPath}/nftList/other/${uuid}`).then(res => {
          this.otherTableData = res.data.nftTranData
          this.otherTotal = this.otherTableData.length;
          if(this.otherTotal===0) {
            this.otherCurrentPage = 1;
          } else {
            this.otherCurrentPage = Math.floor((this.otherTotal-1)/ this.otherPagesize)+1;
          }
        }).catch(err=>{
          console.log(err.response.data.reason)
        })
      },
      getMarketHistoryList: function (uuid) {
        let httpPath = this.$store.state.config.httpPath;
        this.axios.get(`${httpPath}/market/transactionHistory/${uuid}`).then(res => {
          console.log(res.data.nftTranData);
          this.mkTxHistoryTableData = res.data.nftTranData
          this.mkTxHistoryTotal = this.mkTxHistoryTableData.length;
          if(this.mkTxHistoryTotal===0) {
            this.mkTxHistoryCurrentPage = 1;
          } else {
            this.mkTxHistoryCurrentPage = Math.floor((this.mkTxHistoryTotal-1)/ this.mkTxHistoryPagesize)+1;
          }
        }).catch(err=>{
          console.log(err.response.data.reason)
        })
      },
    },
    created: function () {
      // console.log(this.$store.state.account)
      // init variables
      let httpPath = this.$store.state.config.httpPath;
      let uuid = this.$cookies.get('uuid');
      this.uploadDatPath = httpPath + "/file/dat";
      this.uploadAvatarPath = httpPath + "/file/avatar";
      this.uploadOtherPath = httpPath + "/file/other";
      let uploadBaseObject = {
        uuid: uuid,
        number: 0,
        price: 0,
      };

      this.uploadDatAdditionalData = Object.assign({}, uploadBaseObject);
      this.uploadAvatarAdditionalData = Object.assign({}, uploadBaseObject);
      this.uploadOtherAdditionalData = Object.assign({}, uploadBaseObject);
      this.$set(this.uploadDatAdditionalData, 'allowAirdrop', true);
      // get nft list of user from market place
      this.getNFTList(uuid);
      this.getMarketHistoryList(uuid);

      // set default select item
      this.selectedType = this.uploadOptions[0].type;
      this.getIntro(uuid);
      this.getNickname(uuid);
      this.getWallet(uuid);
      this.getAvatarUrl(uuid);
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

  .airDropTitle {
    margin-top: 100px;
    font-size: 1.5rem;
    margin-bottom: 100px;
  }

  .marketPlaceListTitle {
    font-size: 1.2rem;
    margin-bottom: 20px;
  }
</style>
