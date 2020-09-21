import React from 'react'

import Drawer from '@material-ui/core/Drawer';
import Divider from '@material-ui/core/Divider';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import { SvgIconComponent } from '@material-ui/icons';
import { Link } from "react-router-dom";

interface MenuItem {
    name: string;
    text: string;
    icon: SvgIconComponent;
    onClick(name: string): void;
    link: string;
}

interface IProps {
    drawerClassName: string;
    drawerWidth: string;
    drawerOpen: boolean;
    drawerHeader: string;
    items: MenuItem[];
}

export function Menu(props: IProps) {
    const items = props.items.map((i) => {
        const Icon = i.icon;
        return (
            <ListItem button key={i.name} to={i.link} component={Link}>
                <ListItemIcon><Icon/></ListItemIcon>
                <ListItemText primary={i.text}>
                </ListItemText>
            </ListItem>
        )
    });
    return (
        <Drawer
        className={props.drawerClassName}
        variant="persistent"
        anchor="left"
        open={props.drawerOpen}
        classes={{
          paper: props.drawerWidth,
        }}
      >
        <div className={props.drawerHeader}/>
        <Divider />
        <List>
            {items}
        </List>
      </Drawer>
    )
}