package repository

type UserImagesRepository interface {
    SaveUserImage(userId string, imageUrl string) error
    GetUserImageUrl(userId string) (string, error)
}

type FilesRepository interface {

}
