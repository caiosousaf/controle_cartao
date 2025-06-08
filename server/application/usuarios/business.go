package usuarios

import (
	"controle_cartao/config/database"
	"controle_cartao/domain/cadastros/usuarios"
	infra "controle_cartao/infrastructure/cadastros/usuarios"
	"controle_cartao/middlewares"
	"controle_cartao/utils"
	"github.com/google/uuid"
)

// CadastrarUsuario contém a regra de negócio para cadastro de um novo usuário
func CadastrarUsuario(req *ReqUsuario) (res *ResCadastroUsuario, err error) {
	const msgErrPadrao = "Erro ao cadastrar novo usuário"

	var reqInfra = new(infra.Usuario)

	db, err := database.Conectar()
	if err != nil {
		return res, utils.Wrap(err, msgErrPadrao)
	}
	defer db.Close()

	repo := usuarios.NovoRepo(db)

	req.Senha, err = middlewares.HashSenha(req.Senha)
	if err != nil {
		return res, utils.Wrap(err, msgErrPadrao)
	}

	if err = utils.ConvertStructByAlias(req, reqInfra); err != nil {
		return res, utils.Wrap(err, msgErrPadrao)
	}

	if err = repo.CadastrarUsuario(reqInfra); err != nil {
		return res, utils.Wrap(err, msgErrPadrao)
	}

	token, err := middlewares.GerarJWT(*reqInfra.Nome, reqInfra.ID)

	res = &ResCadastroUsuario{
		ID:    reqInfra.ID,
		Token: token,
	}

	return
}

// LoginUsuario contém a regra de negócio para realizar o login do usuário
func LoginUsuario(req *ReqUsuarioLogin) (res *Res, err error) {
	const (
		msgErrPadrao             = "Erro ao realizar login"
		msgErrIdentificarUsuario = "Erro ao identificar usuário"
		msgErrCredenciais        = "Credenciais inválidas"
	)

	db, err := database.Conectar()
	if err != nil {
		return res, utils.Wrap(err, msgErrPadrao)
	}
	defer db.Close()

	repo := usuarios.NovoRepo(db)

	usuario, err := repo.BuscarUsuarioLogin(req.Email)
	if err != nil {
		return res, utils.NewErr(msgErrIdentificarUsuario)
	}

	if ok := middlewares.VerificarSenha(*req.Senha, *usuario.Senha); ok != true {
		return res, utils.NewErr(msgErrCredenciais)
	}

	token, err := middlewares.GerarJWT(*usuario.Nome, usuario.ID)
	if err != nil {
		return res, utils.Wrap(err, msgErrPadrao)
	}

	res = &Res{
		Token: token,
	}

	return
}

// BuscarUsuario contém a regra de negócio para buscar os dados de usuário
func BuscarUsuario(usuarioID *uuid.UUID) (res *ResUsuario, err error) {
	const msgErrPadrao = "Erro ao buscar dados de usuário"

	db, err := database.Conectar()
	if err != nil {
		return res, utils.Wrap(err, msgErrPadrao)
	}
	defer db.Close()

	repo := usuarios.NovoRepo(db)

	dados, err := repo.BuscarUsuario(usuarioID)
	if err != nil {
		return res, utils.NewErr(msgErrPadrao)
	}

	res = &ResUsuario{
		Nome:        dados.Nome,
		Email:       dados.Email,
		DataCriacao: dados.DataCriacao,
	}

	return
}

// AtualizarSenhaUsuario contém a regra de negócio para atualizar a senha do usuárioo
func AtualizarSenhaUsuario(req *ReqAlterarSenhaUsuario, usuarioID *uuid.UUID) (err error) {
	const (
		msgErrPadrao      = "Erro ao atualizar senha"
		msgErrCredenciais = "Credenciais inválidas"
		msgErrSenhaAtual  = "Senha atual inválida"
	)

	db, err := database.Conectar()
	if err != nil {
		return utils.Wrap(err, msgErrPadrao)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return utils.Wrap(err, msgErrPadrao)
	}

	defer tx.Rollback()

	repo := usuarios.NovoRepo(db)

	usuario, err := repo.BuscarUsuarioLogin(req.Email)
	if err != nil {
		return utils.NewErr(msgErrCredenciais)
	}

	if ok := middlewares.VerificarSenha(*req.SenhaAtual, *usuario.Senha); ok != true {
		return utils.NewErr(msgErrSenhaAtual)
	}

	req.SenhaNova, err = middlewares.HashSenha(req.SenhaNova)
	if err != nil {
		return utils.Wrap(err, msgErrPadrao)
	}

	if err := repo.AtualizarSenhaUsuario(req.SenhaNova, req.EmailNovo, usuarioID); err != nil {
		return utils.NewErr(msgErrPadrao)
	}

	if err = tx.Commit(); err != nil {
		return utils.Wrap(err, msgErrPadrao)
	}

	return
}
