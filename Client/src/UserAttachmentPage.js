import React, {Component} from 'react'
import {Table,Label,Input,FormGroup,Button, Modal, ModalHeader, ModalBody, ModalFooter} from 'reactstrap'
import axios from 'axios'

class UserAttachmentPage extends Component{
    state = {
        attachments: [],
        newAttachmentModal:false,

      }
    componentWillMount() {
        axios.get('/user-files/attachments/'+this.props.match.params.username).then((response)=>{
          this.setState({
            attachments:response.data
          })
        }).catch( (error) =>{
            console.log(error.response)
            alert("Error - " + error.response.data.error)
          });
      }
    toggleNewAttachmentModal(){
        this.setState({
            newAttachmentModal:!this.state.newAttachmentModal
        });
    }
      _RefreshAttachments(){
        axios.get('/user-files/attachments/'+this.props.match.params.username).then((response)=>{
          this.setState({
            attachments:response.data
          })
        }).catch( (error) =>{
            console.log(error.response)
            alert("Error - " + error.response.data.error)
          });
      }

      AddAttachment(upload_type) {
        var formData = new FormData();
        var imagefile = document.querySelector('#file');
        formData.append("file", imagefile.files[0]);
        var is_s3 = ""
        if (upload_type != "") {
            is_s3="/"+upload_type
        }
        axios.post('/user-files/attachment/'+this.props.match.params.username+is_s3,formData,{
            headers: {
                'Content-Type': 'multipart/form-data'
              }
        }).then((response)=>{
          console.log(response);
          this._RefreshAttachments();
          this.setState({
              newAttachmentModal:false
          })
        }).catch( (error) =>{
            console.log(error.response)
            alert("Error - " + error.response.data.error)
          });
      }
      DeleteAttachment(id){
        axios.delete('/user-files/attachment/'+id).then((response=>{
          console.log(response.data)
          this._RefreshAttachments()
        })).catch( (error) =>{
            console.log(error.response)
            alert("Error - " + error.response.data.error)
          })
      }
      render(){
        let attachments =this.state.attachments.map((attachment) =>{
          return(
            <tr key={attachment.id}>
              <td>{attachment.id}</td>
              <td><a target='_blank' href={'/user-files/attachment/' + attachment.id + "/get" }>{attachment.file_name}</a></td>
              <td>{attachment.date}</td>
              <td>{attachment.upload_type}</td>
              <td> <Button color="primary" onClick={this.DeleteAttachment.bind(this,attachment.id)}>X</Button></td>
            </tr>
          )
        });
        return (
          <div className="App">
          <Button color="primary" onClick={this.toggleNewAttachmentModal.bind(this)}>Upload file</Button>
          <a href='/'><Button color="secondary">Home</Button></a>
          <a href='/users'><Button color="secondary">Users</Button></a>
          <Modal isOpen={this.state.newAttachmentModal} toggle={this.toggleNewAttachmentModal.bind(this)}>
            <ModalHeader toggle={this.toggleNewAttachmentModal.bind(this)}>Upload new file</ModalHeader>
            <ModalBody>
             <FormGroup>
              <Label for="name">Choose File</Label>
              <Input type="file" name="file" id="file" ></Input>
             </FormGroup>
            </ModalBody>
            <ModalFooter>
              <Button color="primary" onClick={this.AddAttachment.bind(this,"")}>Upload File (FS)</Button>{' '}
              <Button color="primary" onClick={this.AddAttachment.bind(this,"s3")}>Upload File (S3)</Button>{' '}
              <Button color="secondary" onClick={this.toggleNewAttachmentModal.bind(this)}>Cancel</Button>
            </ModalFooter>
          </Modal>
            <Table>
              <thead>
                <tr>
                  <th>#</th>
                  <th>Filename</th>
                  <th>Last Updated</th>
                  <th> Upload Type</th>
                </tr>
              </thead>
              <tbody>
              {attachments.length ? attachments : <div><tr><td colspan={4}>No Attachments found!</td></tr></div>}
              </tbody>
            </Table>
          </div>
        );
      }
}

export default UserAttachmentPage;

