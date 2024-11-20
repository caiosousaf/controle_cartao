package usuarios

import (
	"controle_cartao/config/database"
	"controle_cartao/domain/cadastros/usuarios"
	infra "controle_cartao/infrastructure/cadastros/usuarios"
	"controle_cartao/middlewares"
	"controle_cartao/utils"
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

	usuario, err := repo.BuscarUsuario(req.Email)
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
