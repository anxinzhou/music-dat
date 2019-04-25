<template>
    <section class="bg-light">


      <el-row :gutter="20">

        <el-col :span="8">
          <div class="intro text-center">
            <img src="../assets/images/avatar.jpg">
            <div>
              <!--<h3> Market </h3>-->
              <!--<br>-->
              <!--<br>-->
              <!--<b> With bridge</b>-->
              <br>
              <b> Create and sell your dat and avatar</b>
              <!--<br>-->
              <!--<b>between public and private chain</b>-->
            </div>
          </div>
        </el-col>

        <el-col :span="8">
          <div class="wallet-container">
            <div class=" bridge-item-container">
              <div class="shop-avatar-info-container">
                <div class="role-info-header">
                  Your current number of avatar
                </div>
                <div>
                  1
                </div>
              </div>
            </div>
            <div class="bridge-item-container">
              <div class="token-container">
                <el-upload
                  class="upload-demo upload-dat"
                  ref="upload"
                  name="file"
                  :action="upLoadDatPath"
                  :data="upLoadAdditionalData"
                  :auto-upload="false">
                  <el-button slot="trigger" size="small" type="primary">select file</el-button>
                  <el-button style="margin-left: 10px;" size="small" type="success" @click="submitUpload">upload to server</el-button>
                  <div class="el-upload__tip" slot="tip">jpg/png files with a size less than 500kb</div>
                </el-upload>
              </div>
            </div>
          </div>

        </el-col>

        <el-col :span="8">
            <div class="upload-file-container">
              <el-table
                :data="tableData.slice((currentPage-1)*pagesize,currentPage*pagesize)"
                stripe
                style="width: 100%">
                <el-table-column
                  prop="fileName"
                  label="File Name">
                </el-table-column>
                <el-table-column
                  prop="uploadTime"
                  label="Upload Time">
                </el-table-column>
                <el-table-column
                  prop="fileSize"
                  label="File Size">
                </el-table-column>
              </el-table>
              <div style="text-align: center;margin-top: 30px;">
                <el-pagination
                  background
                  layout="prev, pager, next"
                  :total="total"
                  @current-change="current_change">
                </el-pagination>
              </div>

            </div>
          </el-col>
      </el-row>

    </section>
</template>
<script>
  export default {
    props: {
    },
    data () {
      return {
        fileList: [],
        upLoadDatPath: undefined,
        upLoadAdditionalData:undefined,
        httpPath:undefined,
        wallet: {
          publicBalance: undefined,
          balance: undefined,
          publicEther: undefined
        },
        tableData: [
          {
            fileName: 'a.jpg',
            uploadTime: '2017-05-03 17:30',
            fileSize: '3 MB'
          }, {
            fileName: 'df.jpg',
            uploadTime: '2007-05-03 17:30',
            fileSize: '4 MB'
          }, {
            fileName: 'dfs.mp3',
            uploadTime: '2017-12-03 17:30',
            fileSize: '45 MB'
          }, {
            fileName: 'this.jpg',
            uploadTime: '2017-05-03 04:24',
            fileSize: '567 KB'
          }, {
            fileName: 'c.png',
            uploadTime: '2017-05-03 17:30',
            fileSize: '3 MB'
          }],
        total: 0,
        pagesize:3,
        currentPage:1
      }
    },
    methods: {
      submitUpload: function () {
        this.$refs.upload.submit();
      },
      current_change: function(currentPage) {
        this.currentPage = currentPage
      }
    },
    beforeMount: function() {
      // console.log(this.$store.state.account)
      this.total = this.tableData.length / this.pagesize * 10
      this.upLoadAdditionalData = {
        address: this.$cookies.get('account').address
      }
      this.httpPath = this.$store.state.config.httpPath
      this.upLoadDatPath = this.httpPath + "/file/avatar"
      console.log("upLoadDatPath:",upLoadDatPath)
    }
  }
</script>
