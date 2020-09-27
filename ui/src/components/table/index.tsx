import React, { useState  } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import { map } from 'lodash';
import MTable from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableContainer from '@material-ui/core/TableContainer';
import TableHead from '@material-ui/core/TableHead';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import TableRow from '@material-ui/core/TableRow';
import { SvgIconComponent } from '@material-ui/icons';
import MoreVertIcon from '@material-ui/icons/MoreVert';



const useStyles = makeStyles((theme) => ({
    root: {
      width: '100%',
      height: '100%',
    },
  }));
  

interface IProps {
    columns: Column[];
    data: Row[];
}

export interface Column {
    title: string;
    width?: string;
}
export interface Cell {
    content?: string;
};
export interface Row {
    data: Cell[];
    menu: RowMenuItem[];
}

export interface RowMenuItem {
    title: string;
    icon: SvgIconComponent;
    onClick(row: number): void;
};

export function Table(props: IProps) {
    const classes = useStyles()
    return (
        <TableContainer className={classes.root}>
            <MTable stickyHeader aria-label="sticky table" size={'small'}>
            <TableHead>
              <TableRow>
                {makeColumns(props.columns)}
              </TableRow>
            </TableHead>
            <TableBody>
                {makeRows(props.data)}
            </TableBody>
            </MTable>
        </TableContainer>
    )
}

function makeColumns(columns: Column[]) {
    const actions = { title: '', width: '2px' };
    return map(columns.concat(actions), (c, i) => {
        return (
            <TableCell key={i} style={{width: c.width}}>{c.title}</TableCell>
        )
    })
}

function makeRows(rows: Row[]) {
    return map(rows, (r, i) => {
        return (
            <TableRow key={i}>
                {map(r.data, (c, i) => {
                    return (
                        <TableCell key={i}>
                            {c.content}
                        </TableCell>
                    )
                })}
            <RowMenu
                items={r.menu.map(m => ({
                    onClick: () => m.onClick(i),
                    icon: m.icon,
                    title: m.title,
                }))}/>
            </TableRow>
        )
    });
}

interface IRowMenuProps {
    items: {
        onClick(): void;
        title: string;
        icon: SvgIconComponent;
    }[]
}

function RowMenu(props: IRowMenuProps) {
    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

    const handleClick = (event: any) => {
      setAnchorEl(event.currentTarget);
    };
  
    const handleClose = () => {
      setAnchorEl(null);
    };
    return (
        <TableCell>
            <MoreVertIcon onClick={handleClick}/>
            <Menu
                id="simple-menu"
                anchorEl={anchorEl}
                keepMounted
                open={Boolean(anchorEl)}
                onClose={handleClose}
            >
                {props.items.map((item, i) => 
                    <MenuItem key={i} onClick={ () => {
                            handleClose();
                            return item.onClick();
                        }}>
                        {item.icon && <item.icon/>}
                        {item.title}
                    </MenuItem>
                )}
            </Menu>
        </TableCell>
    )
}