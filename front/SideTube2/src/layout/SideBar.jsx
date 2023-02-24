/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-14 13:43:02
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-01-14 20:43:49
 * @FilePath: /sidetube/src/layout/SideBar.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import * as React from 'react';
import Box from '@mui/material/Box';
import Drawer from '@mui/material/Drawer';
import List from '@mui/material/List';
import Divider from '@mui/material/Divider';
import ListItem from '@mui/material/ListItem';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import InboxIcon from '@mui/icons-material/MoveToInbox';
import MailIcon from '@mui/icons-material/Mail';
import MenuIcon from '@mui/icons-material/Menu';
import Button from '@mui/material/Button'
import { YouTube } from '@mui/icons-material';
import VideoLibraryIcon from '@mui/icons-material/VideoLibrary';

export default function SideBar() {
  const [state, setState] = React.useState({
    show: false,
    position: 'left',

  });

  const toggleDrawer = ( open) => (event) => {
    if (event.type === 'keydown' && (event.key === 'Tab' || event.key === 'Shift')) {
      return;
    }

    setState({ ...state, show: open });
  };

  const menuList = [
    { name: "videoList",
      icon :<VideoLibraryIcon />,
    },
  ]

  const list = (anchor) => (
    <Box
      sx={{ width:  250 }}
      role="presentation"
      onClick={toggleDrawer(anchor, false)}
      onKeyDown={toggleDrawer(anchor, false)}
    >
      <ListItem  style={{paddingLeft:0 }}>
       <ListItemButton>
              <ListItemIcon>
              <MenuIcon/>
              </ListItemIcon>
              <ListItemText >
                <span style={{display:"flex", alignItems: "center"}}>
                <YouTube/>
                SideTube
                </span>
              </ListItemText>
      </ListItemButton>
      </ListItem>
      <Divider />
      <List>
        {['Inbox', 'Starred', 'Send email', 'Drafts'].map((text, index) => (
          <ListItem key={text} disablePadding>
            <ListItemButton>
              <ListItemIcon>
                {index % 2 === 0 ? <InboxIcon /> : <MailIcon />}
              </ListItemIcon>
              <ListItemText primary={text} />
            </ListItemButton>
          </ListItem>
        ))}
      </List>
      <Divider />
      <List>
        {menuList.map((obj, index) => (
          <ListItem key={obj.name} disablePadding>
            <ListItemButton>
              <ListItemIcon>
                {obj.icon}
              </ListItemIcon>
              <ListItemText primary={obj.name} />
            </ListItemButton>
          </ListItem>
        ))}
      </List>
    </Box>
  );

  return (
    <div>
        <React.Fragment>
        <Button onClick={toggleDrawer(true)}  size="large"  color="inherit" startIcon={<MenuIcon/>}></Button>
          <Drawer
            anchor={state['position']}
            open={state['show']}
            onClose={toggleDrawer(false)}
          >
            {list(state['position'])}
          </Drawer>
        </React.Fragment>
    </div>
  );
}