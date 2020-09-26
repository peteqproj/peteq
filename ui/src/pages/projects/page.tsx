import React, { useState, useEffect } from 'react'
import { makeStyles } from '@material-ui/core/styles';
import Paper from '@material-ui/core/Paper';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableContainer from '@material-ui/core/TableContainer';
import TableHead from '@material-ui/core/TableHead';
import TablePagination from '@material-ui/core/TablePagination';
import TableRow from '@material-ui/core/TableRow';
import Fab from '@material-ui/core/Fab';
import Tooltip from '@material-ui/core/Tooltip';
import IconButton from '@material-ui/core/IconButton';
import Backdrop from '@material-ui/core/Backdrop';
import Dialog from '@material-ui/core/Dialog';
import Fade from '@material-ui/core/Fade';
import AddIcon from '@material-ui/icons/Add';
import OpenInNewIcon from '@material-ui/icons/OpenInNew';
import { Link } from "react-router-dom";
import { ProjectModal } from './new-project-modal/modal';
import { ProjectAPI, Project } from './../../services/project'
import { ProjectsViewAPI, ProjectView } from '../../services/views/projects'

const useStyles = makeStyles((theme) => ({
  root: {
    width: '100%',
  },
  container: {
    maxHeight: '50%',
  },
  formControl: {
    margin: theme.spacing(1),
    minWidth: 120,
  },
  modal: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
  },
  paper: {
    backgroundColor: theme.palette.background.paper,
    border: '2px solid #000',
    boxShadow: theme.shadows[5],
    padding: theme.spacing(2, 4, 3),
  },
}));

interface Column {
  title: string;
  minWidth: number;
}
interface Row extends ProjectView {}
function makeColumn(title: string): Column {
  return {
    title,
    minWidth: 30
  }
}

interface IProps  {
  ProjectAPI: ProjectAPI;
  ProjectsViewAPI: ProjectsViewAPI;
}

export function ProjectsPage(props: IProps) {
    const classes = useStyles();
    const [page, setPage] = useState(0);
    const [rowsPerPage, setRowsPerPage] = useState(10);
    const [rows, updateRows] = useState<Row[]>([]);
    useEffect(() => {
      const fetch = async () => {
        const res = await props.ProjectsViewAPI.get()
        console.log(res)
        updateRows(res.projects);
      };

      fetch();
    }, []);
    const handleChangePage = (event: unknown, newPage: number) => {
      setPage(newPage);
    };
    const handleChangeRowsPerPage = (event: React.ChangeEvent<HTMLInputElement>) => {
      setRowsPerPage(+event.target.value);
      setPage(0);
    };
    const [projectModal, setProjectModalOpen] = useState(false);
    const handleProjectModalClose = (data: any) => {
      setProjectModalOpen(false);
    };
    const handleProjectModalOpen = (data: any) => {
      setProjectModalOpen(true);
    };
    const columns = [
      makeColumn('Actions'),
      makeColumn('Name'),
      makeColumn('Description'),
      makeColumn('Tasks'),
    ]
    const addProject = async (row: Row): Promise<void> => {
      return new Promise((resolve) => {
          setTimeout(() => {
            resolve();
            updateRows((prevState) => {
              const data = [...prevState];
              data.splice(0, 0, row)
              return data;
            });
          }, 600);
      });
    }
    return (
      <Paper className={classes.root}>
      <TableContainer className={classes.container}>
        <Table stickyHeader aria-label="sticky table" size={'small'}>
          <TableHead>
              <TableRow>
                <TableCell key={"color"}/>
                {columns.map((column, index) => (
                  <TableCell
                    key={index}
                  
                    style={{ minWidth: column.minWidth }}
                  >
                    {column.title}
                  </TableCell>
                ))}
              </TableRow>
            </TableHead>
            <TableBody>
            {rows.slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage).map((row, index) => {
              return (
                <TableRow hover role="checkbox" tabIndex={-1} key={index}>
                  <TableCell key={'c-00'} style={{
                      backgroundColor: row.metadata.color || 'gray',
                      width: '2px'
                    }}>
                  </TableCell>
                  <TableCell key={'c-0'}>
                  <Tooltip title={"Open"} aria-label={"open"}>
                      <Link to={`/projects/${row.metadata.id}`}>
                        <IconButton aria-label="toggleProjetView" color="primary" id={row.metadata.id}>
                          <OpenInNewIcon/>
                        </IconButton>
                      </Link>
                    </Tooltip>
                  </TableCell>
                  <TableCell key={'c-1'}>
                    {row.metadata.name}
                  </TableCell>
                  <TableCell key={'c-2'}>
                    {row.metadata.description}
                  </TableCell>
                  <TableCell key={'c-3'}>
                    {row.tasks.length}
                  </TableCell>
                </TableRow>
              );
            })}
          </TableBody>
        </Table>
      </TableContainer>
      <TablePagination
        rowsPerPageOptions={[10, 25, 100]}
        component="div"
        count={1}
        rowsPerPage={rowsPerPage}
        page={page}
        onChangePage={handleChangePage}
        onChangeRowsPerPage={handleChangeRowsPerPage}
      />
      <Fab
          color="primary"
          style={{ position: 'fixed', right: '15px', bottom: '15px' }}
          onClick={handleProjectModalOpen}
          >
            <AddIcon />
      </Fab>
      <Dialog
        aria-labelledby="transition-modal-title"
        aria-describedby="transition-modal-description"
        className={classes.modal}
        open={projectModal}
        onClose={handleProjectModalClose}
        closeAfterTransition
        BackdropComponent={Backdrop}
        BackdropProps={{
          timeout: 500,
        }}
      >
        <Fade in={projectModal}>
          <ProjectModal Callback={(project: Project) => {
            setProjectModalOpen(false)
            const row: Row = {
              ...project,
              objects: [],
            }
            addProject(row)
          }} ProjectAPI={props.ProjectAPI}></ProjectModal>
        </Fade>
      </Dialog>
    </Paper>
    )
}