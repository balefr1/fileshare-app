import React, {Component} from 'react';
import {Table,Label,Input,FormGroup,Button, Modal, ModalHeader, ModalBody, ModalFooter} from 'reactstrap'
import axios from 'axios'

class UserPage extends Component{
  state = {
    users: [],
    newUserModal:false,
    newUserData:{
      name:"",
      lastname:"",
      email:"",
      username:""
    }
  }
  componentWillMount() {
    axios.get('/user-api/users').then((response)=>{
      this.setState({
        users:response.data
      })
    }).catch( (error) =>{
      console.log(error.response)
      alert("Error - " + error.response.data.error)
    });
  }
  _RefreshUsers(){
    axios.get('/user-api/users').then((response)=>{
      this.setState({
        users:response.data
      })
    }).catch( (error) =>{
      console.log(error.response)
      alert("Error - " + error.response.data.error)
    });
  }
  AddUser() {
    axios.post('/user-api/user',this.state.newUserData).then((response)=>{
      let {users} = this.state;
      users.push(response.data);
      this.setState({
        users,
        newUserModal:false,
        newUserData:{
          name:"",
          lastname:"",
          email:"",
          username:""
        }
      })
      console.log(response.data);
    }).catch( (error) =>{
      console.log(error.response)
      alert("Error - " + error.response.data.error)
    });
  }
  toggleNewUserModal(){
    this.setState({
      newUserModal:!this.state.newUserModal
    });
  }
  DeleteUser(id){
    axios.delete('/user-api/user/'+id).then((response=>{
      console.log(response)
      this._RefreshUsers()
    })).catch( (error) =>{
      console.log(error.response)
      alert("Error - " + error.response.data.error)
    })
  }
  render(){
    let users =this.state.users.map((user) =>{
      return(
        <tr key={user.id}>
          <td>{user.id}</td>
          <td><a  href={`/user/${user.username}/attachments`}>{user.username}</a></td>
          <td>{user.name}</td>
          <td>{user.lastname}</td>
          <td>{user.email}</td>
          <td> <Button color="primary" onClick={this.DeleteUser.bind(this,user.id)}>X</Button></td>
        </tr>
      )
    });
    return (
      <div className="App">
      <Button color="primary" onClick={this.toggleNewUserModal.bind(this)}>Add User</Button>
      <a href='/'><Button color="secondary">Home</Button></a>
      <Modal isOpen={this.state.newUserModal} toggle={this.toggleNewUserModal.bind(this)}>
        <ModalHeader toggle={this.toggleNewUserModal.bind(this)}>Create new user</ModalHeader>
        <ModalBody>
         <FormGroup>
          <Label for="name">First Name</Label>
          <Input type="text" name="firstname" id="firstname" value={this.state.newUserData.name} placeholder="First Name" onChange={
            (e)=>{
              let {newUserData} = this.state;
              newUserData.name=e.target.value;
              this.setState({newUserData})
            }
          }></Input>
         </FormGroup>
         <FormGroup>
          <Label for="lastname">Last Name</Label>
          <Input type="text" name="lastname" id="lastname" placeholder="Last Name" value={this.state.newUserData.lastname} onChange={
            (e)=>{
              let {newUserData} = this.state;
              newUserData.lastname=e.target.value;
              this.setState({newUserData})
            }
          }></Input>
         </FormGroup>
         <FormGroup>
          <Label for="email">Email</Label>
          <Input type="email" name="email" id="email" placeholder="Email" value={this.state.newUserData.email} onChange={
            (e)=>{
              let {newUserData} = this.state;
              newUserData.email=e.target.value;
              this.setState({newUserData})
            }
          }></Input>
         </FormGroup>
         <FormGroup>
          <Label for="username">Username</Label>
          <Input type="text" name="username" id="username" placeholder="Username" value={this.state.newUserData.username} onChange={
            (e)=>{
              let {newUserData} = this.state;
              newUserData.username=e.target.value;
              this.setState({newUserData})
            }
          }></Input>   
         </FormGroup>
        </ModalBody>
        <ModalFooter>
          <Button color="primary" onClick={this.AddUser.bind(this)}>Create User</Button>{' '}
          <Button color="secondary" onClick={this.toggleNewUserModal.bind(this)}>Cancel</Button>
        </ModalFooter>
      </Modal>
        <Table>
          <thead>
            <tr>
              <th>#</th>
              <th>Username</th>
              <th>FirstName</th>
              <th>LastName</th>
              <th>Email</th>
            </tr>
          </thead>
          <tbody>
          {users.length ? users : <div><tr><td colspan={4}>No Users found!</td></tr></div>}
          </tbody>
        </Table>
      </div>
    );
  }
}

export default UserPage;
